package router

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/garyburd/redigo/redis"
	logging "github.com/op/go-logging"

	"github.com/pkg/errors"
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
		} else if err != nil {
			log.Errorf("pool(%s, %s) : %s", poolName, key, err)
			continue
		}

		return poolName, data, err
	}

	return "", nil, errors.Wrap(ErrNotRoutable, "no pools matches key")
}

/*
updateCache will update keys that does not exist in local cache
currently this is done sequentially, which means a GET for a block
will not succeed unless the cache is updated. The errors of cache SET
will only be logged thought and should not cause the GET operation to fail
*/
func (r *Router) updateCache(src, key string, data []byte) error {
	for _, name := range r.cache {
		if name == src {
			//key was already retrieved from this pool, we skip
			continue
		}

		pool := r.pools[name]
		if err := pool.Set(key, data); err != nil {
			log.Errorf("failed to update cache pool (%s): %s", name, err)
		}
	}

	return nil
}

//Get gets key from table
func (r *Router) Get(key string) (io.ReadCloser, error) {
	src, data, err := r.get(key)
	if err != nil {
		return nil, err
	}

	r.updateCache(src, key, data)

	//TODO CRC check (gonna be dropped)
	if len(data) <= 16 {
		return nil, fmt.Errorf("wrong data size")
	}

	buf := bytes.NewBuffer(data[16:])

	return ioutil.NopCloser(buf), nil
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
			merged.cache = append(merged.cache, name)
		}
	}

	return &merged
}
