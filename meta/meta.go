package meta

import (
	"fmt"
	"github.com/op/go-logging"
	"syscall"
)

var (
	log         = logging.MustGetLogger("meta")
	ErrNotFound = fmt.Errorf("not found")
)

type NodeType uint32

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
		return "unkown type"
	}
}

type Access struct {
	UID  uint32
	GID  uint32
	Mode uint32
}

type MetaInfo struct {
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

type BlockInfo struct {
	Key      []byte
	Decipher []byte
}

type Meta interface {
	fmt.Stringer
	//base name
	ID() string
	Name() string
	IsDir() bool
	Blocks() []BlockInfo

	Info() MetaInfo

	Children() []Meta
}

type MetaStore interface {
	// Populate(entry Entry) error
	Get(name string) (Meta, bool)
}
