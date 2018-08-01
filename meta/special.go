package meta

import np "github.com/threefoldtech/0-fs/cap.np"

type Special struct {
	np.Inode
	special np.Special
	access  Access
}

//ID link id
func (s *Special) ID() string {
	return ""
}

//Name link name
func (s *Special) Name() string {
	name, _ := s.Inode.Name()
	return name
}

//IsDir returns false
func (s *Special) IsDir() bool {
	return false
}

//Blocks returns empty list
func (s *Special) Blocks() []BlockInfo {
	return nil
}

//Children returns empty list
func (s *Special) Children() []Meta {
	return nil
}

//Info returns emtpty list
func (s *Special) Info() MetaInfo {
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
	return MetaInfo{
		CreationTime:     s.CreationTime(),
		ModificationTime: s.ModificationTime(),
		Size:             s.Size(),
		Type:             t,
		Access:           s.access,
		SpecialData:      string(data),
	}
}
