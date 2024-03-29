package main

import (
	"os"
	"path"

	"github.com/threefoldtech/0-fs/meta"
	"github.com/threefoldtech/0-fs/storage"
	"github.com/threefoldtech/0-fs/storage/router"
)

func getDB(db string) (string, error) {
	f, err := os.Open(db)
	if err != nil {
		return db, err
	}

	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return db, err
	}

	if !info.IsDir() {
		db += ".d"
		//dbs[i] = db //update the entry in the list
		if err := meta.Unpack(f, db); err != nil {
			return db, err
		}
	}

	return db, nil
}

func getMetaStore(dbs []string) (meta.Store, error) {
	var stores []meta.Store

	for i, db := range dbs {
		if len(db) == 0 {
			continue //ignore empty lines in file
		}
		var err error
		db, err = getDB(db)
		if err != nil {
			return nil, err
		}

		//update the entry in the list, in case if the db has been extracted
		dbs[i] = db

		store, err := meta.NewStore(db)
		if err != nil {
			return nil, err
		}

		stores = append(stores, store)
	}

	return meta.Layered(stores...), nil
}

func getDataStore(dbs []string, fb *router.Router) (*router.Router, error) {
	var routers []*router.Router
	for _, db := range dbs {
		cfg, err := router.NewConfigFromFile(path.Join(db, "router.yaml"))
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			return nil, err
		}

		r, err := cfg.Router(nil)
		if err != nil {
			return nil, err
		}

		routers = append(routers, r)
	}

	if fb != nil {
		//if fallback router is not nil, add to list
		routers = append(routers, fb)
	}

	return router.Merge(routers...), nil
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

// getStoresFromCmd helper function to initialize stores from cmd line
func getStoresFromCmd(cmd *Cmd) (metaStore meta.Store, dataStore *router.Router, err error) {
	metaStore, err = getMetaStore(cmd.Meta)
	if err != nil {
		return
	}

	if len(cmd.URL) != 0 {
		//prepare the fallback storage
		dataStore, err = storage.NewSimpleStorage(cmd.URL)
		if err != nil {
			return
		}
	}

	//get a merged datastore from all flists
	dataStore, err = getDataStore(cmd.Meta, dataStore)
	if err != nil {
		return
	}

	//finally merge with local router.yaml
	dataStore, err = layerLocalStore(cmd.Router, dataStore)
	if err != nil {
		return
	}

	return
}
