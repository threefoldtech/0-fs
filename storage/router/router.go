package router

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	"github.com/garyburd/redigo/redis"
	logging "github.com/op/go-logging"

	"github.com/pkg/errors"
)

const (
	commitWorkers = 100
)

var (
	log = logging.MustGetLogger("router")
)

/*
Router defines a router engine. The router will try the pools in the same
order defined by the table. And write to cache if wasn't retrieved from the
cache pool
*/
type Router struct {
	pools map[string]Pool

	lookup []string
	cache  map[string]struct{}
	o      sync.Once
	feed   chan chunk
}

type chunk struct {
	key  []byte
	data []byte
}

func (r *Router) init() {
	r.feed = make(chan chunk)
	for i := 0; i < commitWorkers; i++ {
		go r.worker()
	}
}

func (r *Router) worker() {
	for chunk := range r.feed {
		for name := range r.cache {
			pool := r.pools[name]
			if err := pool.Set(chunk.key, chunk.data); err != nil {
				log.Errorf("failed to update cache pool (%s): %s", name, err)
			}
		}
	}
}

func (r *Router) get(key []byte) (string, []byte, error) {
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
		} else if err != nil {
			log.Errorf("pool(%s, %s) : %s", poolName, key, err)
			continue
		}

		return poolName, data, err
	}

	return "", nil, errors.Wrap(ErrNotRoutable, "no pools matches key")
}

/*
updateCache will update keys that does not exist in (any) of the  cache pools
*/
func (r *Router) updateCache(src string, key []byte, data []byte) {
	r.o.Do(r.init)

	//only update cache if src is not one of the cache pools
	if _, ok := r.cache[src]; ok {
		return
	}

	r.feed <- chunk{key: key, data: data}
}

//Get gets key from table
func (r *Router) Get(key []byte) (io.ReadCloser, error) {
	src, data, err := r.get(key)
	if err != nil {
		return nil, err
	}

	r.updateCache(src, key, data)

	return ioutil.NopCloser(bytes.NewBuffer(data)), nil
}

func (r *Router) String() string {
	var buf bytes.Buffer
	for name, pool := range r.pools {
		buf.WriteString(fmt.Sprintf("%s %s\n\n", name, pool))
	}
	buf.WriteString("lookup:\n")
	for _, name := range r.lookup {
		buf.WriteString(fmt.Sprintf("- %s\n", name))
	}
	buf.WriteString("cache:\n")
	for _, name := range r.cache {
		buf.WriteString(fmt.Sprintf("- %s\n", name))
	}
	return buf.String()
}

//Merge merge multiple routers
func Merge(routers ...*Router) *Router {
	merged := Router{
		pools: make(map[string]Pool),
		cache: make(map[string]struct{}),
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

		for name := range router.cache {
			name = fmt.Sprintf("%d.%s", i, name)
			merged.cache[name] = struct{}{}
		}
	}

	return &merged
}
