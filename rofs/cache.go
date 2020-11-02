package rofs

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/threefoldtech/0-fs/meta"
)

func (fs *filesystem) path(hash string) string {
	base := fs.cache
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
func (fs *filesystem) check(m meta.Meta) (os.FileInfo, error) {
	//atomic check and download a file
	name := fs.path(m.ID())
	f, err := fs.ensure(name)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return f.Stat()
}

func (fs *filesystem) ensure(name string) (*os.File, error) {
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

// checkAndGet makes sure the file exists in cache and makes sure the file content is downloaded safely
func (fs *filesystem) checkAndGet(m meta.Meta) (*os.File, error) {
	//atomic check and download a file
	name := fs.path(m.ID())
	f, err := fs.ensure(name)
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
		return f, nil
	}

	if err := fs.download(f, m); err != nil {
		f.Close()
		os.Remove(name)
		return nil, err
	}

	f.Sync()
	f.Seek(0, io.SeekStart)
	return f, nil
}

// download file from storage
func (fs *filesystem) download(file *os.File, m meta.Meta) error {
	downloader := Downloader{
		storage:   fs.storage,
		blockSize: m.Info().FileBlockSize,
		blocks:    m.Blocks(),
	}

	return downloader.Download(file)
}
