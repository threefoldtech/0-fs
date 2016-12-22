package main

import (
	"github.com/g8os/g8ufs/meta"
	"github.com/g8os/g8ufs/rofs"
	"github.com/g8os/g8ufs/storage"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"net/url"
)

func main() {
	store := meta.NewMemoryMetaStore()

	if err := meta.Populate(store, "ubuntu.flist", ""); err != nil {
		panic(err)
	}

	u, _ := url.Parse("https://stor.jumpscale.org/stor2/store/ubuntu-g8os-flist/")

	stor, err := storage.NewAydoStorage(u)
	if err != nil {
		panic(err)
	}

	fs := rofs.New(stor, store, "/tmp")

	server, err := fuse.NewServer(
		nodefs.NewFileSystemConnector(
			pathfs.NewPathNodeFs(fs, nil).Root(),
			nil,
		).RawFS(), "test", &fuse.MountOptions{
			AllowOther: true,
		})

	if err != nil {
		panic(err)
	}

	server.Serve()
}
