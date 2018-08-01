package meta

import np "github.com/threefoldtech/0-fs/cap.np"

type Link struct {
	np.Inode
	link   np.Link
	access Access
}

//ID link id
func (l *Link) ID() string {
	return ""
}

//Name link name
func (l *Link) Name() string {
	name, _ := l.Inode.Name()
	return name
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

//Info returns emtpty list
func (l *Link) Info() MetaInfo {
	target, _ := l.link.Target()
	return MetaInfo{
		CreationTime:     l.CreationTime(),
		ModificationTime: l.ModificationTime(),
		Size:             l.Size(),
		Type:             LinkType,
		Access:           l.access,
		LinkTarget:       target,
	}
}
