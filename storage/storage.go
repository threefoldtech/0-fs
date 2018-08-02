package storage

import (
	"io"

	"github.com/op/go-logging"
	"github.com/threefoldtech/0-fs/storage/router"
)

var (
	log = logging.MustGetLogger("storage")

	//DefaultConfig falls back to hub.gig.tech in case
	//flist does not define a routing table.
	DefaultConfig = router.Config{
		Pools: map[string]router.PoolConfig{
			"hub.gig.tech": {
				"00:FF": "ardb://hub.gig.tech:16379",
			},
		},
		Lookup: []string{
			"hub.gig.tech",
		},
	}
)

//NewSimpleStorage backward compatible storage for a single endpoint
func NewSimpleStorage(url string) (Storage, error) {
	config := router.Config{
		Pools: map[string]router.PoolConfig{
			"simple": {
				"00:FF": url,
			},
		},
		Lookup: []string{
			"simple",
		},
	}

	return config.Router()
}

/*
NewStorage creates a storage from a router.yaml file the config syntax is

	pools:
	  <pool-name>:
		<hash-range>: <scheme>://host[:port]

	lookup:
	  - <pool-name>
	  - ...

Example:
	pools:
	  hub:
		00:FF: ardb://hub.gig.tech:16379

	lookup:
	 - hub
*/
func NewStorage(c io.Reader) (Storage, error) {
	conf, err := router.NewConfig(c)
	if err != nil {
		return nil, err
	}

	return conf.Router()
}

//Storage interface
type Storage interface {
	Get(key string) (io.ReadCloser, error)
}