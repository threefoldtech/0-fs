package main

import (
	"os"
	"path"

	"github.com/threefoldtech/0-fs/meta"
	"github.com/threefoldtech/0-fs/storage/router"
)

func getMetaStore(dbs []string) (meta.MetaStore, error) {
	var stores []meta.MetaStore

	for i, db := range dbs {
		f, err := os.Open(db)
		if err != nil {
			return nil, err
		}
		info, err := f.Stat()
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			db += ".d"
			dbs[i] = db //update the entry in the list
			if err := unpack(f, db); err != nil {
				return nil, err
			}
		}

		f.Close()
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
