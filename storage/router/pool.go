package router

import (
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

//Pool defines a pool interface
type Pool interface {
	Range
	Route(h string) Destination
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
}

/*
ScanPool defines a set of routing rules
This implementation of pool does a sequential scan of the rules. That's not very efficient usually
plus it always returns the first match.

More sophisticated implementation of the pool should balance the routing if more than rule matches
the hash.
*/
type ScanPool struct {
	Rules []Rule
	conn  map[Destination]*redis.Pool

	m sync.Mutex
}

//In checks if hash is in pool
func (p *ScanPool) In(h string) bool {
	for _, rule := range p.Rules {
		if rule.In(h) {
			return true
		}
	}

	return false
}

//Route matches hash against the pool and return the first matched destination
func (p *ScanPool) Route(h string) Destination {
	for _, rule := range p.Rules {
		if rule.In(h) {
			return rule.Destination
		}
	}

	return nil
}

func (p *ScanPool) newPool(d Destination) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			var opts []redis.DialOption
			if d.User != nil {
				//assume ardb://password@host.com:port/
				opts = append(opts, redis.DialPassword(d.User.Username()))
			}

			return redis.Dial("tcp", d.Host, opts...)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		MaxActive:   10,
		IdleTimeout: 1 * time.Minute,
		Wait:        true,
	}
}

func (p *ScanPool) getPool(d Destination) (*redis.Pool, error) {
	p.m.Lock()
	defer p.m.Unlock()

	pool, ok := p.conn[d]
	if ok {
		return pool, nil
	}

	pool = p.newPool(d)
	if p.conn == nil {
		p.conn = make(map[Destination]*redis.Pool)
	}
	p.conn[d] = pool

	return pool, nil
}

//Get key from pool
func (p *ScanPool) Get(key string) ([]byte, error) {
	dest := p.Route(key)
	if dest == nil {
		return nil, ErrNotRoutable
	}

	pool, err := p.getPool(dest)
	if err != nil {
		return nil, err
	}

	con := pool.Get()
	defer con.Close()

	return redis.Bytes(con.Do("GET", key))
}

//Set key to data
func (p *ScanPool) Set(key string, data []byte) error {
	dest := p.Route(key)
	if dest == nil {
		return ErrNotRoutable
	}

	pool, err := p.getPool(dest)
	if err != nil {
		return err
	}

	con := pool.Get()
	defer con.Close()

	_, err = con.Do("SET", key, data)
	return err
}
