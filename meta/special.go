package meta

import (
	"sync"

	np "github.com/threefoldtech/0-fs/cap.np"
)

// Special is the information for a special file in the filesystem
type Special struct {
	np.Inode
	special np.Special
	access  Access

	name string
	info Info

	nOnce sync.Once
	iOnce sync.Once
}

// ID link id
func (s *Special) ID() string {
	return ""
}

// Name link name
func (s *Special) Name() string {
	s.nOnce.Do(func() {
		s.name, _ = s.Inode.Name()
	})

	return s.name
}

// IsDir returns false
func (s *Special) IsDir() bool {
	return false
}

// Blocks returns empty list
func (s *Special) Blocks() []BlockInfo {
	return nil
}

// Children returns empty list
func (s *Special) Children() []Meta {
	return nil
}

// Info returns empty list
func (s *Special) Info() Info {
	s.iOnce.Do(func() {
		s.info = s.getInfo()
	})

	return s.info
}

func (s *Special) getInfo() Info {
	t := UnknownType
	switch s.special.Type() {
	case np.Special_Type_socket:
		t = SocketType
	case np.Special_Type_block:
		t = BlockDeviceType
	case np.Special_Type_chardev:
		t = CharDeviceType
	case np.Special_Type_fifopipe:
		t = FIFOType
	}

	data, _ := s.special.Data()
	return Info{
		CreationTime:     s.CreationTime(),
		ModificationTime: s.ModificationTime(),
		Size:             s.Size(),
		Type:             t,
		Access:           s.access,
		SpecialData:      string(data),
	}
}
