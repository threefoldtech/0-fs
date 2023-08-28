package rofs

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/threefoldtech/0-fs/meta"
	"github.com/threefoldtech/0-fs/storage"
)

type Cache struct {
	cache   string
	storage storage.Storage
}

func NewCache(path string, storage storage.Storage) Cache {
	return Cache{
		cache:   path,
		storage: storage,
	}
}

func (c *Cache) path(hash string) string {
	base := c.cache
	// these checks are here to avoid panicing
	// in case a bad name (hash) was provided
	// it will still return a valid filepath
	if len(hash) >= 2 {
		base = filepath.Join(base, hash[0:2])
	}

	if len(hash) >= 4 {
		base = filepath.Join(base, hash[2:4])
	}

	return path.Join(base, hash)
}

// makes sure file exists in cache and return its stat
func (c *Cache) check(m meta.Meta) (os.FileInfo, error) {
	//atomic check and download a file
	name := c.path(m.ID())
	f, err := c.ensure(name)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return f.Stat()
}

func (c *Cache) ensure(name string) (*os.File, error) {
	for {
		file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0444)
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
				return nil, err
			}

			continue
		} else if err != nil {
			return nil, err
		}

		return file, nil
	}
}

// CheckAndGet makes sure the file exists in cache and makes sure the file content is downloaded safely
func (c *Cache) CheckAndGet(m meta.Meta) (*os.File, error) {
	//atomic check and download a file
	name := c.path(m.ID())
	f, err := c.ensure(name)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return nil, err
	}

	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)

	fstat, err := f.Stat()

	if err != nil {
		return nil, err
	}

	info := m.Info()
	if fstat.Size() == int64(info.Size) {
		log.Debug("cache hit for file with hash", m.ID())
		return f, nil
	}

	if err := c.download(f, m); err != nil {
		f.Close()
		os.Remove(name)
		return nil, err
	}

	f.Sync()
	f.Seek(0, io.SeekStart)
	return f, nil
}

// download file from storage
func (c *Cache) download(file *os.File, m meta.Meta) error {
	downloader := Downloader{
		storage:   c.storage,
		blockSize: m.Info().FileBlockSize,
		blocks:    m.Blocks(),
	}

	return downloader.Download(file)
}
