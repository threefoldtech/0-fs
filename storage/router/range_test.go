package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRage(t *testing.T) {
	_, err := NewRange("00")
	if ok := assert.NoError(t, err); !ok {
		t.Error()
	}

	_, err = NewRange("")

	if ok := assert.Error(t, err); !ok {
		t.Error()
	}

	_, err = NewRange("00:FF")
	if ok := assert.NoError(t, err); !ok {
		t.Error()
	}

	_, err = NewRange("00:FFA")
	if ok := assert.Error(t, err); !ok {
		t.Error()
	}

}

func TestRangeExact(t *testing.T) {
	r, err := NewRange("1C")
	if err != nil {
		t.Fatal(err)
	}

	if ok := assert.True(t, r.In(HexToBytes("1cabc"))); !ok {
		t.Error()
	}

	if ok := assert.False(t, r.In(HexToBytes("0"))); !ok {
		t.Error()
	}

	if ok := assert.False(t, r.In(HexToBytes(""))); !ok {
		t.Error()
	}
}

func TestRange(t *testing.T) {
	r, err := NewRange("01:c1")
	if err != nil {
		t.Fatal(err)
	}

	if ok := assert.True(t, r.In(HexToBytes("01"))); !ok {
		t.Error()
	}

	if ok := assert.True(t, r.In(HexToBytes("11"))); !ok {
		t.Error()
	}

	if ok := assert.True(t, r.In(HexToBytes("01"))); !ok {
		t.Error()
	}

	if ok := assert.True(t, r.In(HexToBytes("ba39055da55fb79da29f23848d3120b220f543dedd9081d0bdf463928eef7491"))); !ok {
		t.Error()
	}

	if ok := assert.False(t, r.In(HexToBytes("c2123"))); !ok {
		t.Error()
	}

	if ok := assert.False(t, r.In(HexToBytes("00"))); !ok {
		t.Error()
	}

}
