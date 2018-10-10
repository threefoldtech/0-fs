package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	g8ufs "github.com/threefoldtech/0-fs"
)

func start(cmd *Cmd, target string) (*g8ufs.G8ufs, error) {
	// Test if the meta path is a directory
	// if not, it's maybe a flist/tar.gz

	metaStore, dataStore, err := getStoresFromCmd(cmd)

	if err != nil {
		return nil, err
	}

	log.Debug("router\n", dataStore)

	return g8ufs.Mount(&g8ufs.Options{
		MetaStore: metaStore,
		Backend:   cmd.Backend,
		Cache:     cmd.Cache,
		Target:    target,
		Storage:   dataStore,
		Reset:     cmd.Reset,
	})
}

func reload(fs *g8ufs.G8ufs, cmd *Cmd) error {
	// metaStore, dataStore, err := getStoresFromCmd(cmd)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func mount(cmd *Cmd, target string) error {
	fs, err := start(cmd, target)
	if err != nil {
		return err
	}

	fmt.Println("mount starts")

	exit := make(chan error)

	go func() {
		exit <- fs.Wait()
	}()

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	defer signal.Stop(sig)

	for {
		select {
		case err := <-exit:
			return err
		case s := <-sig:
			if s == syscall.SIGTERM || s == syscall.SIGINT {
				log.Info("terminating ...")
				fs.Unmount()
				return nil
			}

			//SIGHUP reload store
		}
	}
}
