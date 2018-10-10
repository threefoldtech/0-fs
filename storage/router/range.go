package router

import (
	"bytes"
	"fmt"
	"regexp"
)

var (
	rangeMatchRegex = regexp.MustCompile(`^(?P<start>[0-9a-fA-F]+)(?::(?P<end>[0-9a-fA-F]+))?$`)
)

//HexToBytes converts a hexstring to byte array
func HexToBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	if len(s)%2 == 1 {
		//odd hex (012) for example is invalid, we fix this by
		//adding a trailing 0 to become 01 20
		s += "0"
	}
	var data []byte
	if _, err := fmt.Sscanf(s, "%x", &data); err != nil {
		panic(err)
	}

	return data
}

//Range defines a hash range matcher
type Range interface {
	In(h []byte) bool
}

type exactMatch []byte

func (e exactMatch) In(h []byte) bool {
	if len(h) < len(e) {
		return false
	}

	head := h[0:len(e)]

	return bytes.Compare(e, head) == 0
}

func (e exactMatch) String() string {
	return fmt.Sprintf("%x", []byte(e))
}

type rangeMatch [2][]byte

func (r rangeMatch) String() string {
	return fmt.Sprintf("%x:%x", r[0], r[1])
}

func (r rangeMatch) In(h []byte) bool {
	//range match
	head := h[0:len(r[0])]

	if bytes.Compare(r[0], head) > 0 || bytes.Compare(r[1], head) < 0 {
		return false
	}

	return true
}

//NewRange parse range from string
func NewRange(r string) (Range, error) {
	match := rangeMatchRegex.FindStringSubmatch(r)
	if match == nil {
		return nil, ErrInvalidRange
	}

	start := HexToBytes(match[1])
	end := HexToBytes(match[2])

	if len(end) == 0 {
		return exactMatch(start), nil
	}

	//if range has an end (not exact match) the start and end must be of equal length
	if len(start) != len(end) {
		return nil, ErrInvalidRange
	}

	return rangeMatch{start, end}, nil
}
