package router

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

var (
	rangeMatchRegex = regexp.MustCompile(`^(?P<start>[0-9a-fA-F]+)(?::(?P<end>[0-9a-fA-F]+))?$`)
)

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

	start := strings.ToUpper(match[1])
	end := strings.ToUpper(match[2])

	var startBytes []byte
	fmt.Sscanf(start, "%x", &startBytes)
	if len(end) == 0 {
		return exactMatch(startBytes), nil
	}

	var endBytes []byte
	fmt.Sscanf(end, "%x", &endBytes)

	//if range has an end (not exact match) the start and end must be of equal length
	if len(startBytes) != len(endBytes) {
		return nil, ErrInvalidRange
	}

	return rangeMatch{startBytes, endBytes}, nil
}
