package router

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	blockGetRetries = 3
)

//Pool defines a pool interface
type Pool interface {
	Range
	Route(h []byte) Destination
	Get(key []byte) ([]byte, error)
	Set(key []byte, data []byte) error
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

//NewScanPool initialize a new scan pool
func NewScanPool(rules ...Rule) Pool {
	return &ScanPool{
		Rules: rules,
	}
}

//In checks if hash is in pool
func (p *ScanPool) In(h []byte) bool {
	for _, rule := range p.Rules {
		if rule.In(h) {
			return true
		}
	}

	return false
}

//Route matches hash against the pool and return the first matched destination
func (p *ScanPool) Route(h []byte) Destination {
	for _, rule := range p.Rules {
		if rule.In(h) {
			return rule.Destination
		}
	}

	return nil
}

//Routes returns all possible destinations for hash h.
func (p *ScanPool) Routes(h []byte) []Destination {
	var dest []Destination
	for _, rule := range p.Rules {
		if rule.In(h) {
			dest = append(dest, rule.Destination)
		}
	}

	return dest
}

func (p *ScanPool) newPool(d Destination) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			opts := []redis.DialOption{
				redis.DialNetDial(dial),
			}

			if d.User != nil {
				//assume ardb://password@host.com:port/
				opts = append(opts, redis.DialPassword(d.User.Username()))
			}

			return redis.Dial("tcp", d.Host, opts...)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) > 10*time.Second {
				//only check connection if more than 10 second of inactivity
				_, err := c.Do("PING")
				return err
			}

			return nil
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

func (p *ScanPool) get(pool *redis.Pool, key []byte) ([]byte, error) {
	con := pool.Get()
	defer con.Close()
	trial := 1
	var err error
	var bytes []byte
	for trial <= blockGetRetries {
		log.Debugf("try %x: trial %d/%d", key, trial, blockGetRetries)
		bytes, err = redis.Bytes(con.Do("GET", key))
		if err == nil || err == redis.ErrNil {
			return bytes, err
		}
		trial++
	}

	return bytes, err
}

//Get key from pool
func (p *ScanPool) Get(key []byte) ([]byte, error) {
	dests := p.Routes(key)
	if len(dests) == 0 {
		return nil, ErrNotRoutable
	}

	for _, dest := range dests {
		pool, err := p.getPool(dest)
		if err != nil {
			return nil, err
		}

		data, err := p.get(pool, key)
		if err != nil {
			if err != redis.ErrNil {
				log.Errorf("destination(%s://%s, %x): %s", dest.Scheme, dest.Host, key, err)
			}

			continue
		}

		return data, nil
	}

	return nil, ErrNotFound
}

//Set key to data
func (p *ScanPool) Set(key, data []byte) error {
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

func (p *ScanPool) String() string {
	var buf bytes.Buffer
	buf.WriteString("scan-pool {\n")
	for _, rule := range p.Rules {
		buf.WriteString(
			fmt.Sprintf("%s -> %s://%s\n", rule.Range, rule.Destination.Scheme, rule.Destination.Host),
		)
	}
	buf.WriteString("}")

	return buf.String()
}
