package g8ufs

import "fmt"

/*
The constants in this file are auto-replaced with the actual values
during the build of both core0 and coreX (only using the make file)
*/

var (
	// Branch checkedout during build
	Branch = "{branch}"
	// Revision during build
	Revision = "{revision}"
	// Dirty report the git repository status when building the binary
	Dirty = "{dirty}"
)

type version struct{}

// String implements the fmt.Stringer interface
func (v *version) String() string {
	s := fmt.Sprintf("'%s' @Revision: %s", Branch, Revision)
	if Dirty != "" {
		s += " (dirty-repo)"
	}

	return s
}

// Version returns the version of the app
func Version() fmt.Stringer {
	return &version{}
}
