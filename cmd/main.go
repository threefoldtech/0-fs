package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/op/go-logging"
	"github.com/threefoldtech/0-fs"
	"github.com/threefoldtech/0-fs/meta"
	"github.com/threefoldtech/0-fs/storage"
	"github.com/threefoldtech/0-fs/storage/router"
)

var log = logging.MustGetLogger("main")

type Cmd struct {
	MetaDB  string
	Backend string
	Cache   string
	URL     string
	Router  string
	Reset   bool
	Debug   bool
}

func (c *Cmd) Validate() []error {
	return nil
}

func layerLocalStore(local string, store *router.Router) (*router.Router, error) {
	if len(local) == 0 {
		//no local router
		return store, nil
	}

	config, err := router.NewConfigFromFile(local)
	if err != nil {
		return nil, err
	}

	localRouter, err := config.Router(nil)
	if err != nil {
		return nil, err
	}

	return router.Merge(localRouter, store), nil
}

func mount(cmd *Cmd, target string) error {
	// Test if the meta path is a directory
	// if not, it's maybe a flist/tar.gz

	var metaStore meta.MetaStore
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
				log.Errorf("%s", err)
			} else {
				cmd.MetaDB = cmd.MetaDB + ".d"
			}
		}

		f.Close()

		metaStore, err = meta.NewStore(cmd.MetaDB)
		if err != nil {
			return fmt.Errorf("failed to initialize meta store: %s", err)
		}
	}

	var store *router.Router
	var err error
	router := path.Join(cmd.MetaDB, "router.yaml")
	if file, e := os.Open(router); e == nil {
		log.Debugf("loading router config from: %s", router)
		store, err = storage.NewStorage(file)
		file.Close()
	} else if os.IsNotExist(e) {
		log.Debugf("no router config in meta, fallback to storage-url")
		store, err = storage.NewSimpleStorage(cmd.URL)
	} else {
		return e
	}

	if err != nil {
		return err
	}

	store, err = layerLocalStore(cmd.Router, store)
	if err != nil {
		return err
	}

	log.Debug("router\n", store)

	fs, err := g8ufs.Mount(&g8ufs.Options{
		MetaStore: metaStore,
		Backend:   cmd.Backend,
		Cache:     cmd.Cache,
		Target:    target,
		Storage:   store,
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

func main() {
	var cmd Cmd
	var version bool
	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.BoolVar(&cmd.Reset, "reset", false, "Reset filesystem on mount")
	flag.StringVar(&cmd.MetaDB, "meta", "", "Path to metadata database (optional)")
	flag.StringVar(&cmd.Backend, "backend", "/tmp/backend", "Working directory of the filesystem (cache and others)")
	flag.StringVar(&cmd.Cache, "cache", "", "Optional external (common) cache directory, if not provided a temporary cache location will be created under `backend`")
	flag.StringVar(&cmd.URL, "storage-url", "ardb://hub.gig.tech:16379", "Fallback storage url in case no router.yaml available in flist")
	flag.StringVar(&cmd.Router, "local-router", "", "Path to local router.yaml to merge with the router.yaml from the flist. This will allow adding some caching layers")
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
