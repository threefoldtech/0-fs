package meta

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path"
)

// Unpack decrompress and unpack a tgz (flist) archive from r to dest folder
// dest is created is it doesn't exist
func Unpack(r io.Reader, dest string) error {
	err := os.MkdirAll(dest, 0770)
	if err != nil {
		return err
	}

	zr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	tr := tar.NewReader(zr)
	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return err
		}
		if hdr.Name == "/" || hdr.Name == "./" {
			continue
		}

		f, err := os.OpenFile(path.Join(dest, hdr.Name), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(hdr.Mode))
		if err != nil {
			log.Errorf("%s", err)
			return err
		}
		if _, err := io.Copy(f, tr); err != nil {
			return err
		}

		f.Close()
	}

	return err
}
