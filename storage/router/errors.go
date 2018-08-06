package router

import (
	"bytes"
	"fmt"
)

var (
	//ErrInvalidRange when range syntax is invalid
	ErrInvalidRange = fmt.Errorf("invalid range")

	//ErrNotRoutable when pool can't find a match for this hash
	ErrNotRoutable = fmt.Errorf("not routable")

	//ErrNotFound when key can't be found in a pool
	ErrNotFound = fmt.Errorf("not found")

	//ErrPoolNotFound when table reference a pool that is not configured
	ErrPoolNotFound = fmt.Errorf("pool not found")

	//ErrInvalidDestination
	ErrUnknownScheme = fmt.Errorf("unknown scheme")
)

//Errors holds many errors at once, suitable for config validation
type Errors []error

//Add add error to errors
func (e Errors) Add(err error) Errors {
	return append(e, err)
}

func (e Errors) HasErrors() bool {
	return len(e) > 0
}

func (e Errors) Error() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("found %d errors", len(e)))
	for _, err := range e {
		buf.WriteByte('\n')
		buf.WriteString(fmt.Sprintf("  - %s", err))
	}

	return buf.String()
}
