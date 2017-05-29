package rofs

import (
	"github.com/zero-os/0-fs/meta"
	"os"
	"path"
	"syscall"
)

func (fs *filesystem) path(hash string) string {
	return path.Join(fs.cache, hash)
}

func (fs *filesystem) checkAndGet(m meta.Meta) (*os.File, error) {
	//atomic check and download a file
	name := fs.path(m.ID())
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, os.ModePerm&os.FileMode(0755))
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

	f.Seek(0, os.SEEK_SET)
	return f, nil
}

// download file from storage
func (fs *filesystem) download(file *os.File, m meta.Meta) error {
	downloader := Downloader{
		Storage:   fs.storage,
		BlockSize: m.Info().FileBlockSize,
		Blocks:    m.Blocks(),
	}

	return downloader.Download(file)
}
