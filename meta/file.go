package meta

import (
	"crypto/md5"
	"fmt"
	"sync"

	np "github.com/threefoldtech/0-fs/cap.np"
)

//File represents a file inode
type File struct {
	np.Inode
	file   np.File
	access Access

	name   string
	info   MetaInfo
	blocks []BlockInfo

	nOnce sync.Once
	iOnce sync.Once
	bOnce sync.Once
}

//ID returns file ID
func (f *File) ID() string {
	m := md5.New()
	for _, blk := range f.Blocks() {
		m.Write(blk.Key)
	}
	return fmt.Sprintf("%x", m.Sum(nil))
}

//Name return file name
func (f *File) Name() string {
	f.nOnce.Do(func() {
		f.name, _ = f.Inode.Name()
	})

	return f.name
}

//IsDir false for files
func (f *File) IsDir() bool {
	return false
}

//Children nil for files
func (f *File) Children() []Meta {
	return nil
}

//Info return meta info for this dir
func (f *File) Info() MetaInfo {
	f.iOnce.Do(func() {
		f.info = MetaInfo{
			CreationTime:     f.CreationTime(),
			ModificationTime: f.ModificationTime(),
			Size:             f.Size(),
			Type:             RegularType,
			Access:           f.access,
			FileBlockSize:    uint64(f.file.BlockSize()) * 4096,
		}
	})

	return f.info
}

func (f *File) getBlocks() []BlockInfo {
	var blocks []BlockInfo
	if !f.file.HasBlocks() {
		return blocks
	}

	cblocks, _ := f.file.Blocks()
	for i := 0; i < cblocks.Len(); i++ {
		block := cblocks.At(i)

		hash, _ := block.Hash()
		key, _ := block.Key()
		blocks = append(blocks, BlockInfo{
			Key:      hash,
			Decipher: key,
		})
	}

	return blocks
}

//Blocks loads and return blocks of file
func (f *File) Blocks() []BlockInfo {
	f.bOnce.Do(func() {
		f.blocks = f.getBlocks()
	})

	return f.blocks
}
