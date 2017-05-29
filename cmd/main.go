package main

import (
	"flag"
	"fmt"
	"github.com/zero-os/0-fs"
	"github.com/zero-os/0-fs/meta"
	"github.com/zero-os/0-fs/storage"
	"github.com/op/go-logging"
	"net/url"
	"os"
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
	var errors []error
	if c.MetaDB == "" {
		errors = append(errors,
			fmt.Errorf("meta is required"),
		)
	}

	return errors
}

func mount(cmd *Cmd, target string) error {
	u, err := url.Parse(cmd.URL)
	if err != nil {
		return err
	}

	store, err := meta.NewRocksMeta("", cmd.MetaDB)
	if err != nil {
		return fmt.Errorf("failed to initialize meta store: %s", err)
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

func main() {
	var cmd Cmd
	var version bool
	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.BoolVar(&cmd.Reset, "reset", false, "Reset filesystem on mount")
	flag.StringVar(&cmd.MetaDB, "meta", "", "Path to metadata database (rocksdb)")
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
