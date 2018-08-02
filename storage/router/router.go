package router

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

	return "", nil, ErrNotRoutable
}

//Get gets key from table
func (r *Router) Get(key string) ([]byte, error) {
	_, data, err := r.get(key)
	if err != nil {
		return nil, err
	}

	//TODO: replicate to cache

	return data, nil
}
