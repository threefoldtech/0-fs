package meta

import (
	"sync"

	np "github.com/threefoldtech/0-fs/cap.np"
)

//Dir represents a dir inode
type Dir struct {
	np.Dir
	store  *sqlStore
	access Access

	name     string
	info     MetaInfo
	children []Meta

	nOnce sync.Once
	iOnce sync.Once
	cOnce sync.Once
}

//ID empty string for a dir
func (d *Dir) ID() string {
	return ""
}

//Name return file name
func (d *Dir) Name() string {
	d.nOnce.Do(func() {
		d.name, _ = d.Dir.Name()
	})

	return d.name
}

//IsDir returns true for a dir
func (d *Dir) IsDir() bool {
	return true
}

//Blocks return file block, nil in case of a dir
func (d *Dir) Blocks() []BlockInfo {
	return nil
}

//Info return meta info for this dir
func (d *Dir) Info() MetaInfo {
	d.iOnce.Do(func() {
		d.info = MetaInfo{
			CreationTime:     d.CreationTime(),
			ModificationTime: d.ModificationTime(),
			Size:             4096,
			Type:             DirType,
			Access:           d.access,
		}
	})

	return d.info
}

//Children return items in this dir
func (d *Dir) Children() []Meta {
	d.cOnce.Do(func() {
		d.children = d.getChildren()
	})

	return d.children
}

func (d *Dir) getChildren() []Meta {
	if !d.HasContents() {
		return nil
	}

	contents, err := d.Contents()
	if err != nil {
		log.Errorf("unable to read directly children: %s", err)
		return nil
	}

	var children []Meta
	for i := 0; i < contents.Len(); i++ {
		inode := contents.At(i)

		var m Meta
		attributes := inode.Attributes()
		switch attributes.Which() {
		case np.Inode_attributes_Which_dir:
			dir, _ := attributes.Dir()
			subkey, _ := dir.Key()
			sub, _ := d.store.getDirWithHash(subkey)
			m = sub
		case np.Inode_attributes_Which_file:
			file, _ := attributes.File()
			key, _ := inode.Aclkey()
			access, _ := d.store.getAccess(key)
			m = &File{Inode: inode, file: file, access: access}
		case np.Inode_attributes_Which_link:
			link, _ := attributes.Link()
			key, _ := inode.Aclkey()
			access, _ := d.store.getAccess(key)
			m = &Link{Inode: inode, link: link, access: access}
		case np.Inode_attributes_Which_special:
		default:
			continue
		}
		if m != nil {
			children = append(children, m)
		}
	}

	return children
}
