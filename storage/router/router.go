package router

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/garyburd/redigo/redis"

	"github.com/pkg/errors"
)

/*
Router defines a router engine. The router will try the pools in the same
order defined by the table. And write to cache if wasn't retrieved from the
cache pool
*/
type Router struct {
	pools map[string]Pool

	lookup []string
	cache  []string
}

func (r *Router) get(key string) (string, []byte, error) {
	for _, poolName := range r.lookup {
		pool, ok := r.pools[poolName]
		if !ok {
			return "", nil, ErrPoolNotFound
		}

		data, err := pool.Get(key)
		//only try next entry if entry is not found in this pool, or not routable
		//otherwise return (nil, or other errors)
		if err == ErrNotRoutable || err == redis.ErrNil {
			continue
		}

		return poolName, data, err
	}

	return "", nil, errors.Wrap(ErrNotRoutable, "no pools matches key")
}

//Get gets key from table
func (r *Router) Get(key string) (io.ReadCloser, error) {
	_, data, err := r.get(key)
	if err != nil {
		return nil, err
	}

	//TODO: replicate to cache

	//TODO CRC check (gonna be dropped)
	if len(data) <= 16 {
		return nil, fmt.Errorf("wrong data size")
	}

	buf := bytes.NewBuffer(data[16:])

	return ioutil.NopCloser(buf), nil
}

//Merge merge multiple routers
func Merge(routers ...*Router) *Router {
	merged := Router{
		pools: make(map[string]Pool),
	}

	if len(routers) == 0 {
		panic("invalid call to merge")
	}

	for i, router := range routers {
		for name, pool := range router.pools {
			name = fmt.Sprintf("%d.%s", i, name)
			merged.pools[name] = pool
		}

		for _, name := range router.lookup {
			name = fmt.Sprintf("%d.%s", i, name)
			merged.lookup = append(merged.lookup, name)
		}

		for _, name := range router.cache {
			name = fmt.Sprintf("%d.%s", i, name)
			merged.lookup = append(merged.cache, name)
		}
	}

	return &merged
}
