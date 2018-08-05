package router

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	crcHeader = "xxxxxxxxxxxxxxxx" //16 bytes of CRC
)

type TestPool struct {
	mock.Mock
}

func (t *TestPool) In(h string) bool {
	args := t.Called(h)
	return args.Bool(0)
}
func (t *TestPool) Route(h string) Destination {
	args := t.Called(h)
	return args.Get(0).(Destination)
}

func (t *TestPool) Get(key string) ([]byte, error) {
	args := t.Called(key)
	if data := args.Get(0); data != nil {
		return data.([]byte), args.Error(1)
	}

	return nil, args.Error(1)
}

func (t *TestPool) Set(key string, data []byte) error {
	args := t.Called(key, data)
	return args.Error(0)
}

func newTestPool(rules ...Rule) Pool {
	return &TestPool{}
}

func TestRouterGetSuccess(t *testing.T) {
	config := Config{
		Pools: map[string]PoolConfig{
			"local": PoolConfig{
				"00:FF": "ardb://destination.local:1234",
			},
		},
		Lookup: []string{"local"},
	}

	router, err := config.Router(newTestPool)

	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	key := "abcdef"
	value := "result value"
	pool := router.pools["local"].(*TestPool)
	pool.On("Get", key).Return([]byte(crcHeader+value), nil)
	ret, err := router.Get(key)

	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	result, _ := ioutil.ReadAll(ret)

	if ok := assert.Equal(t, value, string(result)); !ok {
		t.Error()
	}

}

func TestRouterGetLocalMiss(t *testing.T) {
	config := Config{
		Pools: map[string]PoolConfig{
			"local": PoolConfig{
				"00:FF": "ardb://destination.local:1234",
			},
			"remote": PoolConfig{
				"00:FF": "ardb://destination.remote:1234",
			},
		},
		Lookup: []string{"local", "remote"},
	}

	router, err := config.Router(newTestPool)

	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	key := "abcdef"
	value := "result value"
	pool := router.pools["local"].(*TestPool)
	pool.On("Get", key).Return(nil, ErrNotRoutable)
	pool = router.pools["remote"].(*TestPool)
	pool.On("Get", key).Return([]byte(crcHeader+value), nil)

	ret, err := router.Get(key)

	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	result, _ := ioutil.ReadAll(ret)

	if ok := assert.Equal(t, value, string(result)); !ok {
		t.Error()
	}

}

func TestRouterGetError(t *testing.T) {
	config := Config{
		Pools: map[string]PoolConfig{
			"local": PoolConfig{
				"00:FF": "ardb://destination.local:1234",
			},
		},
		Lookup: []string{"local"},
	}

	router, err := config.Router(newTestPool)

	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	key := "abcdef"
	pool := router.pools["local"].(*TestPool)
	pool.On("Get", key).Return(nil, ErrNotRoutable)
	_, err = router.Get(key)
	if ok := assert.Error(t, err); !ok {
		t.Error()
	}
}
