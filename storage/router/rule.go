package router

import (
	"net/url"
)

var (
	SupportedScheme = []string{
		"ardb", "zdb", "redis",
	}
)

//Destination defines a route destination
type Destination *url.URL

//Rule defines a hash routing rule
type Rule struct {
	Range
	Destination Destination
}

//NewDestination parse and validate destination
func NewDestination(dest string) (Destination, error) {
	u, err := url.Parse(dest)
	if err != nil {
		return nil, err
	}

	in := func(s string, l []string) bool {
		for _, a := range l {
			if a == s {
				return true
			}
		}

		return false
	}

	if !in(u.Scheme, SupportedScheme) {
		return nil, ErrUnknownScheme
	}

	return Destination(u), nil
}
