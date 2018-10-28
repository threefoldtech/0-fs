package g8ufs

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/op/go-logging"
	"github.com/threefoldtech/0-fs/meta"
	"github.com/threefoldtech/0-fs/rofs"
	"github.com/threefoldtech/0-fs/storage"
	"golang.org/x/sys/unix"
)

var (
	log = logging.MustGetLogger("g8ufs")
)

type Starter interface {
	Start() error
	Wait() error
}

type Exec func(name string, arg ...string) Starter

type Options struct {
	//Backend (required) working directory where the filesystem keeps it's cache and others
	//will be created if doesn't exist
	Backend string
	//Cache location where downloaded files are gonna be kept (optional). If not provided
	//a cache directly will be created under the backend.
	Cache string
	//Mount (required) is the mount point
	Target string
	//MetaStore (optional), if not provided `Reset` flag will have no effect, and only the backend overlay
	//will be mount at target, allows *full* backups of the backend to be mounted.
	MetaStore meta.MetaStore
	//Storage (required) storage to download files from
	Storage storage.Storage
	//Reset if set, will wipe up the backend clean before mounting.
	Reset bool
}

type G8ufs struct {
	*rofs.Config
	target string
	fuse   string

	w sync.WaitGroup
}

//Mount mounts fuse with given options, it blocks forever until unmount is called on the given mount point
func Mount(opt *Options) (*G8ufs, error) {
	backend := opt.Backend
	ro := path.Join(backend, "ro") //ro lower layer provided by fuse
	rw := path.Join(backend, "rw") //rw upper layer on filyestem
	wd := path.Join(backend, "wd") //wd workdir used by overlayfs
	toSetup := []string{ro, rw, wd}
	ca := path.Join(backend, "ca") //ca cache for downloaded files used by fuse
	if opt.Cache != "" {
		ca = opt.Cache
		os.MkdirAll(ca, 0755)
	} else {
		toSetup = append(toSetup, ca)
	}

	for _, name := range toSetup {
		if opt.MetaStore != nil && opt.Reset {
			os.RemoveAll(name)
		}
		os.MkdirAll(name, 0755)
	}

	cfg := rofs.NewConfig(opt.Storage, opt.MetaStore, ca)
	var server *fuse.Server
	if opt.MetaStore != nil {
		fs := rofs.New(cfg)
		var err error
		server, err = fuse.NewServer(
			nodefs.NewFileSystemConnector(
				pathfs.NewPathNodeFs(fs, nil).Root(),
				nil,
			).RawFS(), ro, &fuse.MountOptions{
				AllowOther: true,
				Options:    []string{"ro"},
			})

		if err != nil {
			return nil, err
		}

		go server.Serve()
		log.Debugf("Waiting for fuse mount")
		server.WaitMount()
	}

	log.Debugf("Fuse mount is complete")

	err := syscall.Mount("overlay",
		opt.Target,
		"overlay",
		syscall.MS_NOATIME,
		fmt.Sprintf(
			"lowerdir=%s,upperdir=%s,workdir=%s",
			ro, rw, wd,
		),
	)

	if err != nil {
		if server != nil {
			server.Unmount()
		}
		return nil, err
	}

	success := false
	for i := 0; i < 5; i++ {
		//wait for mount point
		chk := exec.Command("mountpoint", "-q", opt.Target)
		if err := chk.Run(); err != nil {
			log.Debugf("mount point still not ready: %s", err)
			time.Sleep(time.Second)
			continue
		}
		success = true
		break
	}

	if !success {
		if server != nil {
			server.Unmount()
		}
		return nil, fmt.Errorf("failed to start mount")
	}

	fs := &G8ufs{
		Config: cfg,
		target: opt.Target,
		fuse:   ro,
	}

	fs.w.Add(1)
	go fs.watch()

	return fs, nil
}

func (fs *G8ufs) watch() {
	defer fs.w.Done()

	n, err := unix.InotifyInit()
	if err != nil {
		panic(err)
	}

	defer unix.Close(n)
	_, err = unix.InotifyAddWatch(n, fs.target, unix.IN_IGNORED|unix.IN_UNMOUNT)
	if err != nil {
		panic(err)
	}

	var buffer [4096]byte
	_, err = unix.Read(n, buffer[:])
	if err != nil {
		log.Errorf("failed watching target: %s", err)
		return
	}

	return
}

//Wait filesystem until it's unmounted.
func (fs *G8ufs) Wait() error {
	defer func() {
		fs.umountFuse()
	}()

	fs.w.Wait()

	return nil
}

type errors []interface{}

func (e errors) Error() string {
	return fmt.Sprint(e...)
}

func (fs *G8ufs) umountFuse() error {
	if err := syscall.Unmount(fs.fuse, syscall.MNT_FORCE|syscall.MNT_DETACH); err != nil {
		return err
	}

	return nil
}

func (fs *G8ufs) Unmount() error {
	var errs errors

	if err := syscall.Unmount(fs.target, syscall.MNT_FORCE|syscall.MNT_DETACH); err != nil {
		errs = append(errs, err)
	}

	if err := fs.umountFuse(); err != nil {
		errs = append(errs, err)
	}

	return errs
}
