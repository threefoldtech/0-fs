package g8ufs

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/hanwen/go-fuse/v2/fuse/nodefs"
	"github.com/hanwen/go-fuse/v2/fuse/pathfs"
	"github.com/op/go-logging"
	"github.com/threefoldtech/0-fs/meta"
	"github.com/threefoldtech/0-fs/rofs"
	"github.com/threefoldtech/0-fs/storage"
	"golang.org/x/sys/unix"
)

var (
	log = logging.MustGetLogger("g8ufs")
)

// type Starter interface {
// 	Start() error
// 	Wait() error
// }

// type Exec func(name string, arg ...string) Starter

//Options are mount options
type Options struct {
	//Backend (required) working directory where the filesystem keeps it's cache and others
	//will be created if doesn't exist
	Backend string
	//Cache location where downloaded files are gonna be kept (optional). If not provided
	//a cache directly will be created under the backend.
	Cache string
	//Mount (required) is the mount point
	Target string
	//Store (optional), if not provided `Reset` flag will have no effect, and only the backend overlay
	//will be mount at target, allows *full* backups of the backend to be mounted.
	Store meta.Store
	//Storage (required) storage to download files from
	Storage storage.Storage
	//Reset if set, will wipe up the backend clean before mounting.
	Reset bool
	//Mount fs read-only
	ReadOnly bool
}

//G8ufs struct
type G8ufs struct {
	*rofs.Config
	layers []string
	w      sync.WaitGroup
}

func mountRO(target string, storage storage.Storage, meta meta.Store, cache string) (*G8ufs, error) {
	log.Debugf("ro: '%s' ca: %s", target, cache)

	cfg := rofs.NewConfig(storage, meta, cache)
	fs := rofs.New(cfg)
	// opts := nodefs.Options{Debug: true}
	opts := nodefs.Options{}

	server, err := fuse.NewServer(
		nodefs.NewFileSystemConnector(
			pathfs.NewPathNodeFs(fs, nil).Root(),
			&opts,
		).RawFS(), target, &fuse.MountOptions{
			// Debug:         true,
			AllowOther:    true,
			FsName:        "g8ufs",
			DisableXAttrs: true,
			Options:       []string{"ro", "default_permissions"},
		})

	if err != nil {
		return nil, err
	}

	go server.Serve()

	zfs := &G8ufs{
		Config: cfg,
		layers: []string{target},
	}

	log.Debugf("Waiting for fuse mount")
	server.WaitMount()

	return zfs, nil
}

//Mount mounts fuse with given options, it blocks forever until unmount is called on the given mount point
func Mount(opt *Options) (fs *G8ufs, err error) {
	backend := opt.Backend

	if opt.Reset {
		os.RemoveAll(backend)
	}

	ca := path.Join(backend, "ca") //ca cache for downloaded files used by fuse
	if opt.Cache != "" {
		ca = opt.Cache
	}

	if err = os.MkdirAll(ca, 0755); err != nil && !os.IsExist(err) {
		err = fmt.Errorf("failed to prepare cache directory (%s): %s", ca, err)
		return
	}

	ro := path.Join(backend, "ro") //ro lower layer provided by fuse
	if opt.ReadOnly {
		ro = opt.Target
	}

	if err = os.MkdirAll(ro, 0755); err != nil && !os.IsExist(err) {
		err = fmt.Errorf("failed to create director '%s': %s", ro, err)
		return
	}

	fs, err = mountRO(ro, opt.Storage, opt.Store, ca)
	if err != nil {
		err = fmt.Errorf("failed to do ro layer mount: %s", err)
		return
	}

	log.Debugf("read-only layer mounted")

	defer func() {
		if err != nil {
			fs.Unmount()
			return
		}

		fs.w.Add(1)
		go fs.watch()
	}()

	if opt.ReadOnly {
		return
	}

	rw := path.Join(backend, "rw") //rw upper layer on filyestem
	wd := path.Join(backend, "wd") //wd workdir used by overlayfs

	for _, name := range []string{rw, wd} {
		if err = os.MkdirAll(name, 0755); err != nil && !os.IsExist(err) {
			err = fmt.Errorf("failed to create director '%s': %s", name, err)
			return
		}
	}

	info, err := os.Stat(ro)
	if err == nil {
		//this should not fail ever, because we already make sure that
		//ro exists before we reach here. we still do the check just in
		//case, but no action is needed in case it failed

		//Note, we need to change the `rw` perm to match the `ro` perm
		//so the final mount point has the same permissions as the flist
		os.Chmod(rw, info.Mode())
	}

	err = syscall.Mount("overlay",
		opt.Target,
		"overlay",
		syscall.MS_NOATIME,
		fmt.Sprintf(
			"lowerdir=%s,upperdir=%s,workdir=%s",
			ro, rw, wd,
		),
	)

	if err != nil {
		err = fmt.Errorf("failed to mount overlay: %s", err)
		return
	}

	fs.layers = append(fs.layers, opt.Target)

	mounted := false
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		//wait for mount point
		mounted, err = Mountpoint(opt.Target)
		if err != nil {
			return
		}
		if mounted {
			break
		}
	}

	if !mounted {
		err = fmt.Errorf("failed to start mount")
		return
	}

	return fs, nil
}

func (fs *G8ufs) watch() {
	defer fs.w.Done()

	n, err := unix.InotifyInit()
	if err != nil {
		panic(err)
	}

	defer unix.Close(n)
	// we only watch last mounted layer
	_, err = unix.InotifyAddWatch(n, fs.layers[len(fs.layers)-1], unix.IN_IGNORED|unix.IN_UNMOUNT)
	if err != nil {
		panic(err)
	}

	var buffer [4096]byte
	_, err = unix.Read(n, buffer[:])
	if err != nil {
		log.Errorf("failed watching target: %s", err)
	}
}

//Wait filesystem until it's unmounted.
func (fs *G8ufs) Wait() error {
	defer func() {
		fs.Unmount()
	}()

	fs.w.Wait()

	return nil
}

type errors []interface{}

func (e errors) Error() string {
	return fmt.Sprint(e[:]...)
}

// Unmount make sure 0-fs is unmounted properly
func (fs *G8ufs) Unmount() error {
	var errs errors

	for i := len(fs.layers) - 1; i >= 0; i-- {
		if err := syscall.Unmount(fs.layers[i], syscall.MNT_FORCE|syscall.MNT_DETACH); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// Mountpoint checks if a given path is a mount point
func Mountpoint(path string) (bool, error) {
	if err := exec.Command("mountpoint", "-q", path).Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
