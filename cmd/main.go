package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/op/go-logging"
	"github.com/threefoldtech/0-fs"
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
			cli.StringFlag{
				Name:  "log",
				Usage: "write logs to file (default to stdout)",
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

			if log := ctx.GlobalString("log"); len(log) != 0 {
				file, err := os.Create(log)
				if err != nil {
					return err
				}

				logging.SetBackend(logging.NewLogBackend(file, "", 0))
			}

			return nil
		},
		Action: action,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
