package router

import (
	"io/ioutil"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	crcHeader = "xxxxxxxxxxxxxxxx" //16 bytes of CRC
)

type TestPool struct {
	mock.Mock
	wg sync.WaitGroup
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
	defer t.wg.Done()
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

func TestMerget(t *testing.T) {
	local := Config{
		Pools: map[string]PoolConfig{
			"local": PoolConfig{
				"00:FF": "ardb://destination.local:1234",
			},
		},
		Lookup: []string{"local"},
		Cache:  []string{"local"},
	}

	remote := Config{
		Pools: map[string]PoolConfig{
			"remote": PoolConfig{
				"00:FF": "ardb://destination.local:1234",
			},
		},
		Lookup: []string{"remote"},
	}

	localRouter, err := local.Router(newTestPool)
	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}
	remoteRouter, err := remote.Router(newTestPool)
	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	localPool := localRouter.pools["local"].(*TestPool)
	remotePool := remoteRouter.pools["remote"].(*TestPool)

	key := "abcdef"
	value := "result value"
	//The set is expected to be call on localPool with the value retrieved from remote
	localPool.On("Set", key, []byte(crcHeader+value)).Return(nil)
	localPool.On("Get", key).Return(nil, ErrNotRoutable)
	remotePool.On("Get", key).Return([]byte(crcHeader+value), nil)

	localPool.wg.Add(1)

	router := Merge(localRouter, remoteRouter)

	if ok := assert.Equal(t, []string{"0.local", "1.remote"}, router.lookup); !ok {
		t.Error()
	}

	if ok := assert.Equal(t, map[string]struct{}{"0.local": struct{}{}}, router.cache); !ok {
		t.Error()
	}

	ret, err := router.Get(key)

	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	result, _ := ioutil.ReadAll(ret)

	if ok := assert.Equal(t, value, string(result)); !ok {
		t.Error()
	}

	if ok := localPool.AssertCalled(t, "Get", key); !ok {
		t.Error()
	}
	localPool.wg.Wait()
	if ok := localPool.AssertCalled(t, "Set", key, []byte(crcHeader+value)); !ok {
		t.Error()
	}

	if ok := remotePool.AssertCalled(t, "Get", key); !ok {
		t.Error()
	}

}
