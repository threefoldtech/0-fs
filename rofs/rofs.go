package rofs

import (
	"fmt"
	"github.com/zero-os/0-fs/meta"
	"github.com/zero-os/0-fs/storage"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/op/go-logging"
	"math"
)

const (
	blkSize = 4 * 1024
)

var (
	log = logging.MustGetLogger("rofs")
)

type filesystem struct {
	pathfs.FileSystem
	cache   string
	store   meta.MetaStore
	storage storage.Storage
}

func New(storage storage.Storage, store meta.MetaStore, cache string) pathfs.FileSystem {
	fs := &filesystem{
		FileSystem: pathfs.NewDefaultFileSystem(),
		storage:    storage,
		store:      store,
		cache:      cache,
	}

	return pathfs.NewReadonlyFileSystem(fs)
}

func (fs *filesystem) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	log.Debugf("GetAttr %s", name)
	m, ok := fs.store.Get(name)
	if !ok {
		return nil, fuse.ENOENT
	}

	info := m.Info()
	if info.Type == meta.UnknownType {
		return nil, fuse.EIO
	}

	nodeType := uint32(info.Type)

	access := info.Access

	blocks := uint64(math.Ceil(float64(info.Size / blkSize)))

	var major, minor uint32
	if info.SpecialData != "" {
		fmt.Sscanf(info.SpecialData, "%d,%d", &major, &minor)
	}

	size := info.Size
	if info.Type == meta.LinkType {
		size = uint64(len(info.LinkTarget))
	}

	return &fuse.Attr{
		Size:   size,
		Mtime:  uint64(info.ModificationTime),
		Mode:   nodeType | access.Mode,
		Blocks: blocks,
		Owner: fuse.Owner{
			Uid: access.UID,
			Gid: access.GID,
		},
		Rdev:    major<<8 | minor,
		Blksize: blkSize, //4K blocks
	}, fuse.OK
}

func (fs *filesystem) Open(name string, flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	log.Debugf("Open %s", name)
	if flags&fuse.O_ANYWRITE != 0 {
		return nil, fuse.EPERM
	}
	m, ok := fs.store.Get(name)
	if !ok {
		return nil, fuse.ENOENT
	}
	f, err := fs.checkAndGet(m)
	if err != nil {
		log.Errorf("Failed to open/download the file: %s", err)
	}

	return nodefs.NewReadOnlyFile(nodefs.NewLoopbackFile(f)), fuse.OK
}

func (fs *filesystem) OpenDir(name string, context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	log.Debugf("OpenDir %s", name)
	m, ok := fs.store.Get(name)
	if !ok {
		return nil, fuse.ENOENT
	}
	var entries []fuse.DirEntry
	for _, child := range m.Children() {
		info := child.Info()
		log.Debugf("child '%s', type: %s", child.Name(), info.Type)
		entries = append(entries, fuse.DirEntry{
			Mode: uint32(info.Type),
			Name: child.Name(),
		})
	}

	return entries, fuse.OK
}

func (fs *filesystem) String() string {
	return "G8UFS"
}

func (fs *filesystem) Access(name string, mode uint32, context *fuse.Context) fuse.Status {
	return fuse.OK
}

func (fs *filesystem) Readlink(name string, context *fuse.Context) (string, fuse.Status) {
	log.Debugf("Readlink %s", name)
	m, ok := fs.store.Get(name)
	if !ok {
		return "", fuse.ENOENT
	}

	info := m.Info()

	return info.LinkTarget, fuse.OK
}

func (fs *filesystem) GetXAttr(name string, attr string, context *fuse.Context) ([]byte, fuse.Status) {
	log.Debugf("GetXAttr %s", name)
	return nil, fuse.ENOSYS
}

func (fs *filesystem) ListXAttr(name string, context *fuse.Context) ([]string, fuse.Status) {
	log.Debugf("ListXAttr %s", name)
	return nil, fuse.ENOSYS
}

func (fs *filesystem) StatFs(name string) *fuse.StatfsOut {
	return &fuse.StatfsOut{}
}
