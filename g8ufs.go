package g8ufs

import (
	"fmt"
	"github.com/g8os/g8ufs/meta"
	"github.com/g8os/g8ufs/rofs"
	"github.com/g8os/g8ufs/storage"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"os"
	"os/exec"
	"path"
	"syscall"
)

type Options struct {
	//PList (required) to mount
	PList     string
	//PListTrim (optional) trim prefix of file paths
	PListTrim string
	//Backend (required) working directory where the filesystem keeps it's cache and others
	//will be created if doesn't exist
	Backend   string
	//Mount (required) is the mount point
	Target    string
	//MetaStore (optional) will use meta.NewMemoryMeta if not provided
	MetaStore meta.MetaStore
	//Storage (required) storage to download files from
	Storage   storage.Storage
	//Reset if set, will wipe up the backend clean before mounting.
	Reset     bool
}

//Mount mounts fuse with given options, it blocks forever until unmount is called on the given mount point
func Mount(opt *Options) error {
	backend := opt.Backend
	ro := path.Join(backend, "ro")
	rw := path.Join(backend, "rw")
	ca := path.Join(backend, "ca")

	for _, name := range []string{ro, rw, ca} {
		if opt.Reset {
			os.RemoveAll(name)
		}
		os.MkdirAll(name, 0755)
	}

	metaStore := opt.MetaStore
	if metaStore == nil {
		metaStore = meta.NewMemoryMetaStore()
	}

	if err := meta.Populate(metaStore, opt.PList, rw, opt.PListTrim); err != nil {
		return err
	}

	fs := rofs.New(opt.Storage, metaStore, ca)

	server, err := fuse.NewServer(
		nodefs.NewFileSystemConnector(
			pathfs.NewPathNodeFs(fs, nil).Root(),
			nil,
		).RawFS(), ro, &fuse.MountOptions{
			AllowOther: true,
			Options:    []string{"ro"},
		})

	if err != nil {
		return err
	}

	//TODO: the following code should be moved to the library
	//so we can reuse it.
	go server.Serve()
	server.WaitMount()

	defer server.Unmount()

	branch := fmt.Sprintf("%s=RW:%s=RO", rw, ro)
	cmd := exec.Command("unionfs", "-f",
		"-o", "cow",
		"-o", "allow_other",
		branch, opt.Target)

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error staring union fs: %s\n%s", err, string(out))
	}

	return nil
}

func Unmount(mount string) error {
	return syscall.Unmount(mount, syscall.MNT_FORCE|syscall.MNT_DETACH)
}
