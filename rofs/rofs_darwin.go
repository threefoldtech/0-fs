package rofs

import (
	"fmt"
	"math"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/zero-os/0-fs/meta"
)

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
		Rdev: major<<8 | minor,
	}, fuse.OK
}
