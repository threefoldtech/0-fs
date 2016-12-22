package main

import (
	"fmt"
	"github.com/g8os/g8ufs/meta"
	"github.com/g8os/g8ufs/rofs"
	"github.com/g8os/g8ufs/storage"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
)

func main() {
	backend := "/tmp/backend"
	mount := "/mnt/fuse"

	ro := path.Join(backend, "ro")
	rw := path.Join(backend, "rw")
	ca := path.Join(backend, "ca")
	for _, name := range []string{ro, rw, ca} {
		os.RemoveAll(name)
		os.MkdirAll(name, 0755)
	}

	store := meta.NewMemoryMetaStore()

	if err := meta.Populate(store, "ubuntu.flist", rw, ""); err != nil {
		panic(err)
	}
	u, _ := url.Parse("https://stor.jumpscale.org/stor2/store/ubuntu-g8os-flist/")

	stor, err := storage.NewAydoStorage(u)
	if err != nil {
		panic(err)
	}

	fs := rofs.New(stor, store, ca)

	server, err := fuse.NewServer(
		nodefs.NewFileSystemConnector(
			pathfs.NewPathNodeFs(fs, nil).Root(),
			nil,
		).RawFS(), ro, &fuse.MountOptions{
			AllowOther: true,
			Options:    []string{"ro"},
		})

	if err != nil {
		panic(err)
	}

	//TODO: the following code should be moved to the library
	//so we can reuse it.
	go server.Serve()
	server.WaitMount()

	defer server.Unmount()

	//unionfs -o cow rw=RW:ro=RO /mnt/fus
	branch := fmt.Sprintf("%s=RW:%s=RO", rw, ro)
	cmd := exec.Command("unionfs", "-f",
		"-o", "cow",
		"-o", "allow_other",
		branch, mount)

	if out, err := cmd.CombinedOutput(); err != nil {
		log.Println(string(out))
		panic(err)
	}
}
