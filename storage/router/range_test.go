package router_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/0-fs/storage/router"
)

func TestNewRage(t *testing.T) {
	_, err := router.NewRange("00")
	if ok := assert.NoError(t, err); !ok {
		t.Error()
	}

	_, err = router.NewRange("")

	if ok := assert.Error(t, err); !ok {
		t.Error()
	}

	_, err = router.NewRange("00:FF")
	if ok := assert.NoError(t, err); !ok {
		t.Error()
	}

	_, err = router.NewRange("00:FFA")
	if ok := assert.Error(t, err); !ok {
		t.Error()
	}

}

func TestRangeExact(t *testing.T) {
	r, err := router.NewRange("1C")
	if err != nil {
		t.Fatal(err)
	}

	if ok := assert.True(t, r.In("1cabc")); !ok {
		t.Error()
	}

	if ok := assert.False(t, r.In("0")); !ok {
		t.Error()
	}

	if ok := assert.False(t, r.In("")); !ok {
		t.Error()
	}
}

func TestRange(t *testing.T) {
	r, err := router.NewRange("01:c1")
	if err != nil {
		t.Fatal(err)
	}

	if ok := assert.True(t, r.In("01")); !ok {
		t.Error()
	}

	if ok := assert.True(t, r.In("11")); !ok {
		t.Error()
	}

	if ok := assert.True(t, r.In("01")); !ok {
		t.Error()
	}

	if ok := assert.True(t, r.In("ba39055da55fb79da29f23848d3120b220f543dedd9081d0bdf463928eef7491")); !ok {
		t.Error()
	}

	if ok := assert.False(t, r.In("c2123")); !ok {
		t.Error()
	}

	if ok := assert.False(t, r.In("00")); !ok {
		t.Error()
	}

}
