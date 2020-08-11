package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/sevlyar/go-daemon"
	"github.com/threefoldtech/0-fs/meta"

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
		Store: metaStore,
		Backend:   cmd.Backend,
		Cache:     cmd.Cache,
		Target:    target,
		Storage:   dataStore,
		Reset:     cmd.Reset,
		ReadOnly:  cmd.ReadOnly,
	})
}

func reload(fs *g8ufs.G8ufs, cmd *Cmd) error {
	log.Info("reload flists")
	//load extra flist from external file /backend/.layered
	content, err := ioutil.ReadFile(path.Join(cmd.Backend, ".layered"))
	if os.IsNotExist(err) {
		return nil //nothing to do
	} else if err != nil {
		return err
	}

	//rebuild the stores
	extra := strings.Split(string(content), "\n")
	extraMeta, err := getMetaStore(extra)
	if err != nil {
		return err
	}

	// - first use the ones passed via command line
	metaStore, _, err := getStoresFromCmd(cmd)
	if err != nil {
		return err
	}

	// - then add the extra on top
	metaStore = meta.Layered(metaStore, extraMeta)
	fs.SetMetaStore(metaStore)

	return nil
}

func mount(cmd *Cmd, target string) error {
	if cmd.LogPath == "" {
		cmd.LogPath = "/var/log/g8ufs.log"
	}
	cntxt := &daemon.Context{
		PidFileName: cmd.PidPath,
		PidFilePerm: 0644,
		LogFileName: cmd.LogPath,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	var (
		fs  *g8ufs.G8ufs
		err error
	)

	if cmd.Daemon {
		child, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if child != nil {
			// parent process stops
			return nil
		}
	}
	defer cntxt.Release()

	fs, err = start(cmd, target)
	if err != nil {
		return err
	}

	// this line is very important because it works as
	// a signal to core0 that the rootfs of the container
	// is ready and then proceed with starting the container
	// init itself. without this print statement core0 will
	// wait for sometime before times out.
	fmt.Println("mount starts")

	exit := make(chan error)

	go func() {
		exit <- fs.Wait()
	}()

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	defer signal.Stop(sig)

	log.Info("mount ready")

	for {
		select {
		case err := <-exit:
			log.Info("filesystem unmounted, terminating")
			return err
		case s := <-sig:
			if s == syscall.SIGTERM || s == syscall.SIGINT {
				log.Info("terminating ...")
				fs.Unmount()
				return nil
			}

			if err := reload(fs, cmd); err != nil {
				log.Errorf("failed to reload flists: %s", err)
			}
		}
	}
}
