package router

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

/*
Router defines a router engine. The router will try the pools in the same
order defined by the table. And write to cache if wasn't retrieved from the
cache pool
*/
type Router struct {
	Pools map[string]Pool

	Lookup []string
	Cache  []string
}

func (r *Router) get(key string) (string, []byte, error) {
	for _, poolName := range r.Lookup {
		pool, ok := r.Pools[poolName]
		if !ok {
			return "", nil, ErrPoolNotFound
		}

		data, err := pool.Get(key)
		//only try next entry if entry is not found in this pool
		//otherwise return (nil, or other errors)
		if err == ErrNotRoutable {
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
