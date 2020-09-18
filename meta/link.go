package meta

import (
	"sync"

	np "github.com/threefoldtech/0-fs/cap.np"
)

// Link is a inode struct representing a link in the filesystem
type Link struct {
	np.Inode
	link   np.Link
	access Access

	name string
	info Info

	nOnce sync.Once
	iOnce sync.Once
}

//ID link id
func (l *Link) ID() string {
	return ""
}

//Name link name
func (l *Link) Name() string {
	l.nOnce.Do(func() {
		l.name, _ = l.Inode.Name()
	})

	return l.name
}

//IsDir returns false
func (l *Link) IsDir() bool {
	return false
}

//Blocks returns empty list
func (l *Link) Blocks() []BlockInfo {
	return nil
}

//Children returns empty list
func (l *Link) Children() []Meta {
	return nil
}

//Info returns empty list
func (l *Link) Info() Info {
	l.iOnce.Do(func() {
		target, _ := l.link.Target()
		l.info = Info{
			CreationTime:     l.CreationTime(),
			ModificationTime: l.ModificationTime(),
			Size:             l.Size(),
			Type:             LinkType,
			Access:           l.access,
			LinkTarget:       target,
		}
	})

	return l.info
}
