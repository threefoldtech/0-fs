package router

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	rangeMatchRegex = regexp.MustCompile(`^(?P<start>[0-9a-fA-F]+)(?::(?P<end>[0-9a-fA-F]+))?$`)
)

//Range defines a hash range matcher
type Range interface {
	In(h string) bool
}

type exactMatch string

func (e exactMatch) In(h string) bool {
	if len(h) < len(e) {
		return false
	}

	head := strings.ToUpper(h[0:len(e)])

	return strings.Compare(string(e), head) == 0
}

type rangeMatch [2]string

func (r rangeMatch) String() string {
	return fmt.Sprintf("%s:%s", r[0], r[1])
}

func (r rangeMatch) In(h string) bool {
	//range match
	head := strings.ToUpper(h[0:len(r[0])])

	if strings.Compare(r[0], head) > 0 || strings.Compare(r[1], head) < 0 {
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

	if len(end) == 0 {
		return exactMatch(start), nil
	}

	//if range has an end (not exact match) the start and end must be of equal length
	if len(end) != len(start) {
		return nil, ErrInvalidRange
	}

	return rangeMatch{start, end}, nil
}
