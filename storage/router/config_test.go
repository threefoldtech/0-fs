package router

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigSimple(t *testing.T) {
	str := `
pools:
  hub:
    00:ff: ardb://destination

lookup:
  - hub
`
	buf := bytes.NewBufferString(str)

	config, err := NewConfig(buf)

	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	if ok := assert.Equal(t, []string{"hub"}, config.Lookup); !ok {
		t.Error()
	}

	if ok := assert.Len(t, config.Pools, 1); !ok {
		t.Error()
	}

	if ok := assert.Equal(t, "ardb://destination", config.Pools["hub"]["00:ff"]); !ok {
		t.Error()
	}
}

func TestNewConfig(t *testing.T) {
	str := `
pools:
  hub:
    00:ff: ardb://destination.remote
  local:
    00:ff: zdb://destination.local

lookup:
  - local
  - hub

cache:
  - local
`
	buf := bytes.NewBufferString(str)

	config, err := NewConfig(buf)

	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	if ok := assert.Equal(t, []string{"local", "hub"}, config.Lookup); !ok {
		t.Error()
	}

	if ok := assert.Equal(t, []string{"local"}, config.Cache); !ok {
		t.Error()
	}

	if ok := assert.Len(t, config.Pools, 2); !ok {
		t.Error()
	}

	if ok := assert.Equal(t, "ardb://destination.remote", config.Pools["hub"]["00:ff"]); !ok {
		t.Error()
	}

	if ok := assert.Equal(t, "zdb://destination.local", config.Pools["local"]["00:ff"]); !ok {
		t.Error()
	}
}

func TestConfigRouter(t *testing.T) {
	str := `
pools:
  hub:
    00:ff: zdb://destination.remote
  local:
    00:ff: zdb://destination.local

lookup:
  - local
  - hub

cache:
  - local
`

	buf := bytes.NewBufferString(str)

	config, err := NewConfig(buf)

	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	table, err := config.Router(NewScanPool)

	if ok := assert.NoError(t, err); !ok {
		t.Fatal()
	}

	if ok := assert.Equal(t, []string{"local", "hub"}, table.lookup); !ok {
		t.Error()
	}

	if ok := assert.Equal(t, map[string]struct{}{"local": struct{}{}}, table.cache); !ok {
		t.Error()
	}

	if ok := assert.Len(t, table.pools, 2); !ok {
		t.Error()
	}

	hub, ok := table.pools["hub"].(*ScanPool)

	if ok := assert.True(t, ok); !ok {
		t.Fatal()
	}

	if ok := assert.Len(t, hub.Rules, 1); !ok {
		t.Fatal()
	}

	rule := hub.Rules[0]

	if ok := assert.Equal(t, fmt.Sprint(rule.Range), "00:FF"); !ok {
		t.Error()
	}
}
