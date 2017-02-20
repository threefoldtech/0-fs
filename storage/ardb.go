package storage

import (
	"bytes"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"io/ioutil"
	"net/url"
	"time"
)

func NewARDBStorage(u *url.URL) (Storage, error) {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			var opts []redis.DialOption
			if u.User != nil {
				//assume ardb://password@host.com:port/
				opts = append(opts, redis.DialPassword(u.User.Username()))
			}

			return redis.Dial("tcp", u.Host, opts...)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		MaxActive:   10,
		IdleTimeout: 1 * time.Minute,
		Wait:        true,
	}

	c := pool.Get()
	defer c.Close()
	if _, err := c.Do("PING"); err != nil {
		return nil, err
	}

	return &ardbStor{
		pool: pool,
	}, nil
}

type ardbStor struct {
	pool *redis.Pool
}

func (s *ardbStor) Get(key string) (io.ReadCloser, error) {
	cl := s.pool.Get()
	defer cl.Close()

	data, err := redis.Bytes(cl.Do("GET", key))
	if err != nil {
		return nil, err
	}
	//TODO CRC check
	if len(data) <= 16 {
		return nil, fmt.Errorf("wrong data size")
	}

	return ioutil.NopCloser(bytes.NewBuffer(data[16:])), nil
}
