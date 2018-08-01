package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"

	"github.com/op/go-logging"
	"github.com/threefoldtech/0-fs"
	"github.com/threefoldtech/0-fs/meta"
	"github.com/threefoldtech/0-fs/storage"
)

var log = logging.MustGetLogger("main")

type Cmd struct {
	MetaDB  string
	Backend string
	Cache   string
	URL     string
	Reset   bool
	Debug   bool
}

func (c *Cmd) Validate() []error {
	return nil
}

func mount(cmd *Cmd, target string) error {
	u, err := url.Parse(cmd.URL)
	if err != nil {
		return err
	}

	// Test if the meta path is a directory
	// if not, it's maybe a flist/tar.gz

	var store meta.MetaStore
	if len(cmd.MetaDB) != 0 {
		f, err := os.Open(cmd.MetaDB)
		if err != nil {
			return err
		}
		info, err := f.Stat()
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err = unpack(f, cmd.MetaDB+".d")
			if err != nil {
				log.Error(err)
			} else {
				cmd.MetaDB = cmd.MetaDB + ".d"
			}
		}

		store, err = meta.NewRocksStore("", cmd.MetaDB)
		if err != nil {
			return fmt.Errorf("failed to initialize meta store: %s", err)
		}
	}

	aydo, err := storage.NewARDBStorage(u)
	if err != nil {
		return err
	}

	fs, err := g8ufs.Mount(&g8ufs.Options{
		MetaStore: store,
		Backend:   cmd.Backend,
		Cache:     cmd.Cache,
		Target:    target,
		Storage:   aydo,
		Reset:     cmd.Reset,
	})

	if err != nil {
		return err
	}
	fmt.Println("mount starts")
	if err := fs.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "exit with error: %s", err)
		os.Exit(1)
	}

	return nil
}

// unpack decrompress and unpackt a tgz archive from r to dest folder
// dest is created is it doesn't exist
func unpack(r io.Reader, dest string) error {
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
		if hdr.Name == "/" {
			continue
		}

		f, err := os.OpenFile(path.Join(dest, hdr.Name), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(hdr.Mode))
		if err != nil {
			log.Error(err)
			return err
		}
		if _, err := io.Copy(f, tr); err != nil {
			return err
		}

		f.Close()
	}

	return err
}

func main() {
	var cmd Cmd
	var version bool
	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.BoolVar(&cmd.Reset, "reset", false, "Reset filesystem on mount")
	flag.StringVar(&cmd.MetaDB, "meta", "", "Path to metadata database (optional)")
	flag.StringVar(&cmd.Backend, "backend", "/tmp/backend", "Working directory of the filesystem (cache and others)")
	flag.StringVar(&cmd.Cache, "cache", "", "Optional external (common) cache directory, if not provided a temporary cache location will be created under `backend`")
	flag.StringVar(&cmd.URL, "storage-url", "ardb://hub.gig.tech:16379", "Storage url")
	flag.BoolVar(&cmd.Debug, "debug", false, "Print debug messages")

	flag.Parse()

	if version {
		fmt.Println(g8ufs.Version())
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Missing mount point argument\n")
		os.Exit(1)
	}

	formatter := logging.MustStringFormatter("%{time}: %{color}%{module} %{level:.1s} > %{message} %{color:reset}")
	logging.SetFormatter(formatter)

	if cmd.Debug {
		logging.SetLevel(logging.DEBUG, "")
	} else {
		logging.SetLevel(logging.INFO, "")
	}

	if errs := cmd.Validate(); errs != nil {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		os.Exit(1)
	}

	if err := mount(&cmd, flag.Arg(0)); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
