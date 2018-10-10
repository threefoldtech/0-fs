package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/op/go-logging"
	"github.com/threefoldtech/0-fs"
	"github.com/threefoldtech/0-fs/storage"
	"github.com/threefoldtech/0-fs/storage/router"
)

var log = logging.MustGetLogger("main")

type Cmd struct {
	Meta    []string
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

	metaStore, err := getMetaStore(cmd.Meta)
	if err != nil {
		return err
	}

	var dataStore *router.Router
	if len(cmd.URL) != 0 {
		//prepare the fallback storage
		dataStore, err = storage.NewSimpleStorage(cmd.URL)
		if err != nil {
			return err
		}
	}

	//get a merged datastore from all flists
	dataStore, err = getDataStore(cmd.Meta, dataStore)
	if err != nil {
		return err
	}

	//finally merge with local router.yaml
	dataStore, err = layerLocalStore(cmd.Router, dataStore)
	if err != nil {
		return err
	}

	log.Debug("router\n", dataStore)

	fs, err := g8ufs.Mount(&g8ufs.Options{
		MetaStore: metaStore,
		Backend:   cmd.Backend,
		Cache:     cmd.Cache,
		Target:    target,
		Storage:   dataStore,
		Reset:     cmd.Reset,
	})

	if err != nil {
		return err
	}

	fmt.Println("mount starts")

	return fs.Wait()
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

func action(ctx *cli.Context) error {
	args := ctx.Args()
	if len(args) != 1 {
		return fmt.Errorf("expecting a single mount point argument")
	}

	cmd := Cmd{
		Meta:    ctx.GlobalStringSlice("meta"),
		Backend: ctx.GlobalString("backend"),
		Cache:   ctx.GlobalString("cache"),
		URL:     ctx.GlobalString("storage-url"),
		Router:  ctx.GlobalString("local-router"),
		Reset:   ctx.GlobalBool("reset"),
	}

	return mount(&cmd, args.First())
}

func main() {
	app := cli.App{
		Name:      "0-fs",
		Usage:     "start a zero-fs instance by mounting one or more flists into mount target",
		UsageText: "0-fs [options] <mount-target>",
		Version:   g8ufs.Version().String(),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "reset",
				Usage: "resets filesystem on mount",
			},
			cli.BoolFlag{
				Name:  "debug",
				Usage: "enable debug logging",
			},
			cli.StringSliceFlag{
				Name:  "meta",
				Usage: "path to meta backend, can appear many times. The meta is layered in order so last meta to be added will be on top",
			},
			cli.StringFlag{
				Name:  "backend",
				Value: "/tmp/backend",
				Usage: "working directory of the filesystem (cache and others)",
			},
			cli.StringFlag{
				Name:  "cache",
				Usage: "external (common) cache directory, if not provided a temporary cache location will be created under `backend`",
			},
			cli.StringFlag{
				Name:  "storage-url",
				Value: "zdb://hub.grid.tf:9900",
				Usage: "fallback storage url in case no router.yaml available in flist",
			},
			cli.StringFlag{
				Name:  "local-router",
				Usage: "path to local router.yaml to merge with the router.yaml from the flist. This will allow adding some caching layers",
			},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.GlobalBool("version") {
				fmt.Println(g8ufs.Version())
				os.Exit(0)
			}

			formatter := logging.MustStringFormatter("%{time}: %{color}%{module} %{level:.1s} > %{message} %{color:reset}")
			logging.SetFormatter(formatter)

			if ctx.GlobalBool("debug") {
				logging.SetLevel(logging.DEBUG, "")
			} else {
				logging.SetLevel(logging.INFO, "")
			}

			return nil
		},
		Action: action,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
