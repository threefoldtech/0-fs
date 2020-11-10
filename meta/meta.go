package meta

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("meta")

	//ErrNotFound in case of an entry miss
	ErrNotFound = fmt.Errorf("not found")
)

// NodeType is the enum for all different file types
type NodeType uint32

// NodeType enum values
const (
	UnknownType     = NodeType(0)
	DirType         = NodeType(syscall.S_IFDIR)
	RegularType     = NodeType(syscall.S_IFREG)
	BlockDeviceType = NodeType(syscall.S_IFBLK)
	CharDeviceType  = NodeType(syscall.S_IFCHR)
	SocketType      = NodeType(syscall.S_IFSOCK)
	FIFOType        = NodeType(syscall.S_IFIFO)
	LinkType        = NodeType(syscall.S_IFLNK)
)

// String implements fmt.Stringer interface
func (nt NodeType) String() string {
	switch nt {
	case DirType:
		return "dir type"
	case RegularType:
		return "file type"
	case BlockDeviceType:
		return "block device type"
	case CharDeviceType:
		return "char device type"
	case SocketType:
		return "socket type"
	case FIFOType:
		return "fifo type"
	case LinkType:
		return "link type"
	default:
		return "unknown type"
	}
}

// Access is the ACL of a file
type Access struct {
	UID  uint32
	GID  uint32
	Mode uint32
}

// Info is the metadata of a file
type Info struct {
	//Common
	CreationTime     uint32
	ModificationTime uint32
	Access           Access
	Type             NodeType
	Size             uint64

	//Specific Attr

	//Link
	LinkTarget string

	//File
	FileBlockSize uint64

	//Special
	SpecialData string
}

// BlockInfo is the information needed to retrieve and decrypt a data block
type BlockInfo struct {
	Key      []byte
	Decipher []byte
}

// Meta is an interface that can be implemented by any type that needs to be used as metadata store for the filesystem
type Meta interface {
	fmt.Stringer
	//base name
	ID() string
	Name() string
	IsDir() bool
	Blocks() []BlockInfo

	Info() Info

	Children() []Meta
}

// WalkFn walk function
type WalkFn func(path string, meta Meta) error

// ErrSkipDir if returned by the WalkFn the directory is skipped
var ErrSkipDir = errors.New("skip this directory")

// Store is the interface to implement to read filesystem metadata from an flist
type Store interface {
	// Populate(entry Entry) error
	Get(name string) (Meta, bool)
	Close() error
}

// Walker interface, some stores can implement this interface
type Walker interface {
	Store
	// Walk walks over the given path and call fn for each entry
	// that is found under this given path, including `path` itself
	// Note the following:
	// - if fn return ErrSkipDir and meta.IsDir() is true
	//   then this directory is not scanned.
	// - if fn return ErrSkipDir and meta.IsDir() is false
	//   then the remaining items of the current directory are skipped
	// - Walk does not follow symlinks, so if the meta.IsDir() is false
	//   extra checks can be done on the "file" type from the meta.Info()
	//   if Type is SymLink the caller can decide to read the Target and
	//   follow the link
	Walk(path string, fn WalkFn) error
}
