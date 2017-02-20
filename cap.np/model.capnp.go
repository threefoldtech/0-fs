package np

// AUTO GENERATED - DO NOT EDIT

import (
	strconv "strconv"
	capnp "zombiezen.com/go/capnproto2"
	text "zombiezen.com/go/capnproto2/encoding/text"
	schemas "zombiezen.com/go/capnproto2/schemas"
)

type FileBlock struct{ capnp.Struct }

// FileBlock_TypeID is the unique identifier for the type FileBlock.
const FileBlock_TypeID = 0xd5a2538369c2f82a

func NewFileBlock(s *capnp.Segment) (FileBlock, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return FileBlock{st}, err
}

func NewRootFileBlock(s *capnp.Segment) (FileBlock, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return FileBlock{st}, err
}

func ReadRootFileBlock(msg *capnp.Message) (FileBlock, error) {
	root, err := msg.RootPtr()
	return FileBlock{root.Struct()}, err
}

func (s FileBlock) String() string {
	str, _ := text.Marshal(0xd5a2538369c2f82a, s.Struct)
	return str
}

func (s FileBlock) Hash() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return []byte(p.Data()), err
}

func (s FileBlock) HasHash() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s FileBlock) SetHash(v []byte) error {
	return s.Struct.SetData(0, v)
}

func (s FileBlock) Key() ([]byte, error) {
	p, err := s.Struct.Ptr(1)
	return []byte(p.Data()), err
}

func (s FileBlock) HasKey() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s FileBlock) SetKey(v []byte) error {
	return s.Struct.SetData(1, v)
}

// FileBlock_List is a list of FileBlock.
type FileBlock_List struct{ capnp.List }

// NewFileBlock creates a new list of FileBlock.
func NewFileBlock_List(s *capnp.Segment, sz int32) (FileBlock_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	return FileBlock_List{l}, err
}

func (s FileBlock_List) At(i int) FileBlock { return FileBlock{s.List.Struct(i)} }

func (s FileBlock_List) Set(i int, v FileBlock) error { return s.List.SetStruct(i, v.Struct) }

// FileBlock_Promise is a wrapper for a FileBlock promised by a client call.
type FileBlock_Promise struct{ *capnp.Pipeline }

func (p FileBlock_Promise) Struct() (FileBlock, error) {
	s, err := p.Pipeline.Struct()
	return FileBlock{s}, err
}

type File struct{ capnp.Struct }

// File_TypeID is the unique identifier for the type File.
const File_TypeID = 0xecfda38634f4591a

func NewFile(s *capnp.Segment) (File, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return File{st}, err
}

func NewRootFile(s *capnp.Segment) (File, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return File{st}, err
}

func ReadRootFile(msg *capnp.Message) (File, error) {
	root, err := msg.RootPtr()
	return File{root.Struct()}, err
}

func (s File) String() string {
	str, _ := text.Marshal(0xecfda38634f4591a, s.Struct)
	return str
}

func (s File) BlockSize() uint16 {
	return s.Struct.Uint16(0)
}

func (s File) SetBlockSize(v uint16) {
	s.Struct.SetUint16(0, v)
}

func (s File) Blocks() (FileBlock_List, error) {
	p, err := s.Struct.Ptr(0)
	return FileBlock_List{List: p.List()}, err
}

func (s File) HasBlocks() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s File) SetBlocks(v FileBlock_List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewBlocks sets the blocks field to a newly
// allocated FileBlock_List, preferring placement in s's segment.
func (s File) NewBlocks(n int32) (FileBlock_List, error) {
	l, err := NewFileBlock_List(s.Struct.Segment(), n)
	if err != nil {
		return FileBlock_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// File_List is a list of File.
type File_List struct{ capnp.List }

// NewFile creates a new list of File.
func NewFile_List(s *capnp.Segment, sz int32) (File_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return File_List{l}, err
}

func (s File_List) At(i int) File { return File{s.List.Struct(i)} }

func (s File_List) Set(i int, v File) error { return s.List.SetStruct(i, v.Struct) }

// File_Promise is a wrapper for a File promised by a client call.
type File_Promise struct{ *capnp.Pipeline }

func (p File_Promise) Struct() (File, error) {
	s, err := p.Pipeline.Struct()
	return File{s}, err
}

type Link struct{ capnp.Struct }

// Link_TypeID is the unique identifier for the type Link.
const Link_TypeID = 0xe419a0e5a661965c

func NewLink(s *capnp.Segment) (Link, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Link{st}, err
}

func NewRootLink(s *capnp.Segment) (Link, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Link{st}, err
}

func ReadRootLink(msg *capnp.Message) (Link, error) {
	root, err := msg.RootPtr()
	return Link{root.Struct()}, err
}

func (s Link) String() string {
	str, _ := text.Marshal(0xe419a0e5a661965c, s.Struct)
	return str
}

func (s Link) Target() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s Link) HasTarget() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Link) TargetBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s Link) SetTarget(v string) error {
	return s.Struct.SetText(0, v)
}

// Link_List is a list of Link.
type Link_List struct{ capnp.List }

// NewLink creates a new list of Link.
func NewLink_List(s *capnp.Segment, sz int32) (Link_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return Link_List{l}, err
}

func (s Link_List) At(i int) Link { return Link{s.List.Struct(i)} }

func (s Link_List) Set(i int, v Link) error { return s.List.SetStruct(i, v.Struct) }

// Link_Promise is a wrapper for a Link promised by a client call.
type Link_Promise struct{ *capnp.Pipeline }

func (p Link_Promise) Struct() (Link, error) {
	s, err := p.Pipeline.Struct()
	return Link{s}, err
}

type Special struct{ capnp.Struct }

// Special_TypeID is the unique identifier for the type Special.
const Special_TypeID = 0xdc74a897ce683c6b

func NewSpecial(s *capnp.Segment) (Special, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Special{st}, err
}

func NewRootSpecial(s *capnp.Segment) (Special, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Special{st}, err
}

func ReadRootSpecial(msg *capnp.Message) (Special, error) {
	root, err := msg.RootPtr()
	return Special{root.Struct()}, err
}

func (s Special) String() string {
	str, _ := text.Marshal(0xdc74a897ce683c6b, s.Struct)
	return str
}

func (s Special) Type() Special_Type {
	return Special_Type(s.Struct.Uint16(0))
}

func (s Special) SetType(v Special_Type) {
	s.Struct.SetUint16(0, uint16(v))
}

func (s Special) Data() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return []byte(p.Data()), err
}

func (s Special) HasData() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Special) SetData(v []byte) error {
	return s.Struct.SetData(0, v)
}

// Special_List is a list of Special.
type Special_List struct{ capnp.List }

// NewSpecial creates a new list of Special.
func NewSpecial_List(s *capnp.Segment, sz int32) (Special_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return Special_List{l}, err
}

func (s Special_List) At(i int) Special { return Special{s.List.Struct(i)} }

func (s Special_List) Set(i int, v Special) error { return s.List.SetStruct(i, v.Struct) }

// Special_Promise is a wrapper for a Special promised by a client call.
type Special_Promise struct{ *capnp.Pipeline }

func (p Special_Promise) Struct() (Special, error) {
	s, err := p.Pipeline.Struct()
	return Special{s}, err
}

type Special_Type uint16

// Values of Special_Type.
const (
	Special_Type_socket   Special_Type = 0
	Special_Type_block    Special_Type = 1
	Special_Type_chardev  Special_Type = 2
	Special_Type_fifopipe Special_Type = 3
	Special_Type_unknown  Special_Type = 4
)

// String returns the enum's constant name.
func (c Special_Type) String() string {
	switch c {
	case Special_Type_socket:
		return "socket"
	case Special_Type_block:
		return "block"
	case Special_Type_chardev:
		return "chardev"
	case Special_Type_fifopipe:
		return "fifopipe"
	case Special_Type_unknown:
		return "unknown"

	default:
		return ""
	}
}

// Special_TypeFromString returns the enum value with a name,
// or the zero value if there's no such value.
func Special_TypeFromString(c string) Special_Type {
	switch c {
	case "socket":
		return Special_Type_socket
	case "block":
		return Special_Type_block
	case "chardev":
		return Special_Type_chardev
	case "fifopipe":
		return Special_Type_fifopipe
	case "unknown":
		return Special_Type_unknown

	default:
		return 0
	}
}

type Special_Type_List struct{ capnp.List }

func NewSpecial_Type_List(s *capnp.Segment, sz int32) (Special_Type_List, error) {
	l, err := capnp.NewUInt16List(s, sz)
	return Special_Type_List{l.List}, err
}

func (l Special_Type_List) At(i int) Special_Type {
	ul := capnp.UInt16List{List: l.List}
	return Special_Type(ul.At(i))
}

func (l Special_Type_List) Set(i int, v Special_Type) {
	ul := capnp.UInt16List{List: l.List}
	ul.Set(i, uint16(v))
}

type SubDir struct{ capnp.Struct }

// SubDir_TypeID is the unique identifier for the type SubDir.
const SubDir_TypeID = 0xa4a421ce00f301dd

func NewSubDir(s *capnp.Segment) (SubDir, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return SubDir{st}, err
}

func NewRootSubDir(s *capnp.Segment) (SubDir, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return SubDir{st}, err
}

func ReadRootSubDir(msg *capnp.Message) (SubDir, error) {
	root, err := msg.RootPtr()
	return SubDir{root.Struct()}, err
}

func (s SubDir) String() string {
	str, _ := text.Marshal(0xa4a421ce00f301dd, s.Struct)
	return str
}

func (s SubDir) Key() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s SubDir) HasKey() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s SubDir) KeyBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s SubDir) SetKey(v string) error {
	return s.Struct.SetText(0, v)
}

// SubDir_List is a list of SubDir.
type SubDir_List struct{ capnp.List }

// NewSubDir creates a new list of SubDir.
func NewSubDir_List(s *capnp.Segment, sz int32) (SubDir_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return SubDir_List{l}, err
}

func (s SubDir_List) At(i int) SubDir { return SubDir{s.List.Struct(i)} }

func (s SubDir_List) Set(i int, v SubDir) error { return s.List.SetStruct(i, v.Struct) }

// SubDir_Promise is a wrapper for a SubDir promised by a client call.
type SubDir_Promise struct{ *capnp.Pipeline }

func (p SubDir_Promise) Struct() (SubDir, error) {
	s, err := p.Pipeline.Struct()
	return SubDir{s}, err
}

type Inode struct{ capnp.Struct }
type Inode_attributes Inode
type Inode_attributes_Which uint16

const (
	Inode_attributes_Which_dir     Inode_attributes_Which = 0
	Inode_attributes_Which_file    Inode_attributes_Which = 1
	Inode_attributes_Which_link    Inode_attributes_Which = 2
	Inode_attributes_Which_special Inode_attributes_Which = 3
)

func (w Inode_attributes_Which) String() string {
	const s = "dirfilelinkspecial"
	switch w {
	case Inode_attributes_Which_dir:
		return s[0:3]
	case Inode_attributes_Which_file:
		return s[3:7]
	case Inode_attributes_Which_link:
		return s[7:11]
	case Inode_attributes_Which_special:
		return s[11:18]

	}
	return "Inode_attributes_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

// Inode_TypeID is the unique identifier for the type Inode.
const Inode_TypeID = 0xc0029f81b3eee594

func NewInode(s *capnp.Segment) (Inode, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3})
	return Inode{st}, err
}

func NewRootInode(s *capnp.Segment) (Inode, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3})
	return Inode{st}, err
}

func ReadRootInode(msg *capnp.Message) (Inode, error) {
	root, err := msg.RootPtr()
	return Inode{root.Struct()}, err
}

func (s Inode) String() string {
	str, _ := text.Marshal(0xc0029f81b3eee594, s.Struct)
	return str
}

func (s Inode) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s Inode) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Inode) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s Inode) SetName(v string) error {
	return s.Struct.SetText(0, v)
}

func (s Inode) Size() uint64 {
	return s.Struct.Uint64(0)
}

func (s Inode) SetSize(v uint64) {
	s.Struct.SetUint64(0, v)
}

func (s Inode) Attributes() Inode_attributes { return Inode_attributes(s) }

func (s Inode_attributes) Which() Inode_attributes_Which {
	return Inode_attributes_Which(s.Struct.Uint16(8))
}
func (s Inode_attributes) Dir() (SubDir, error) {
	p, err := s.Struct.Ptr(1)
	return SubDir{Struct: p.Struct()}, err
}

func (s Inode_attributes) HasDir() bool {
	if s.Struct.Uint16(8) != 0 {
		return false
	}
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Inode_attributes) SetDir(v SubDir) error {
	s.Struct.SetUint16(8, 0)
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewDir sets the dir field to a newly
// allocated SubDir struct, preferring placement in s's segment.
func (s Inode_attributes) NewDir() (SubDir, error) {
	s.Struct.SetUint16(8, 0)
	ss, err := NewSubDir(s.Struct.Segment())
	if err != nil {
		return SubDir{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

func (s Inode_attributes) File() (File, error) {
	p, err := s.Struct.Ptr(1)
	return File{Struct: p.Struct()}, err
}

func (s Inode_attributes) HasFile() bool {
	if s.Struct.Uint16(8) != 1 {
		return false
	}
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Inode_attributes) SetFile(v File) error {
	s.Struct.SetUint16(8, 1)
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewFile sets the file field to a newly
// allocated File struct, preferring placement in s's segment.
func (s Inode_attributes) NewFile() (File, error) {
	s.Struct.SetUint16(8, 1)
	ss, err := NewFile(s.Struct.Segment())
	if err != nil {
		return File{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

func (s Inode_attributes) Link() (Link, error) {
	p, err := s.Struct.Ptr(1)
	return Link{Struct: p.Struct()}, err
}

func (s Inode_attributes) HasLink() bool {
	if s.Struct.Uint16(8) != 2 {
		return false
	}
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Inode_attributes) SetLink(v Link) error {
	s.Struct.SetUint16(8, 2)
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewLink sets the link field to a newly
// allocated Link struct, preferring placement in s's segment.
func (s Inode_attributes) NewLink() (Link, error) {
	s.Struct.SetUint16(8, 2)
	ss, err := NewLink(s.Struct.Segment())
	if err != nil {
		return Link{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

func (s Inode_attributes) Special() (Special, error) {
	p, err := s.Struct.Ptr(1)
	return Special{Struct: p.Struct()}, err
}

func (s Inode_attributes) HasSpecial() bool {
	if s.Struct.Uint16(8) != 3 {
		return false
	}
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Inode_attributes) SetSpecial(v Special) error {
	s.Struct.SetUint16(8, 3)
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewSpecial sets the special field to a newly
// allocated Special struct, preferring placement in s's segment.
func (s Inode_attributes) NewSpecial() (Special, error) {
	s.Struct.SetUint16(8, 3)
	ss, err := NewSpecial(s.Struct.Segment())
	if err != nil {
		return Special{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

func (s Inode) Aclkey() (string, error) {
	p, err := s.Struct.Ptr(2)
	return p.Text(), err
}

func (s Inode) HasAclkey() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
}

func (s Inode) AclkeyBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(2)
	return p.TextBytes(), err
}

func (s Inode) SetAclkey(v string) error {
	return s.Struct.SetText(2, v)
}

func (s Inode) ModificationTime() uint32 {
	return s.Struct.Uint32(12)
}

func (s Inode) SetModificationTime(v uint32) {
	s.Struct.SetUint32(12, v)
}

func (s Inode) CreationTime() uint32 {
	return s.Struct.Uint32(16)
}

func (s Inode) SetCreationTime(v uint32) {
	s.Struct.SetUint32(16, v)
}

// Inode_List is a list of Inode.
type Inode_List struct{ capnp.List }

// NewInode creates a new list of Inode.
func NewInode_List(s *capnp.Segment, sz int32) (Inode_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3}, sz)
	return Inode_List{l}, err
}

func (s Inode_List) At(i int) Inode { return Inode{s.List.Struct(i)} }

func (s Inode_List) Set(i int, v Inode) error { return s.List.SetStruct(i, v.Struct) }

// Inode_Promise is a wrapper for a Inode promised by a client call.
type Inode_Promise struct{ *capnp.Pipeline }

func (p Inode_Promise) Struct() (Inode, error) {
	s, err := p.Pipeline.Struct()
	return Inode{s}, err
}

func (p Inode_Promise) Attributes() Inode_attributes_Promise {
	return Inode_attributes_Promise{p.Pipeline}
}

// Inode_attributes_Promise is a wrapper for a Inode_attributes promised by a client call.
type Inode_attributes_Promise struct{ *capnp.Pipeline }

func (p Inode_attributes_Promise) Struct() (Inode_attributes, error) {
	s, err := p.Pipeline.Struct()
	return Inode_attributes{s}, err
}

func (p Inode_attributes_Promise) Dir() SubDir_Promise {
	return SubDir_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

func (p Inode_attributes_Promise) File() File_Promise {
	return File_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

func (p Inode_attributes_Promise) Link() Link_Promise {
	return Link_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

func (p Inode_attributes_Promise) Special() Special_Promise {
	return Special_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

type Dir struct{ capnp.Struct }

// Dir_TypeID is the unique identifier for the type Dir.
const Dir_TypeID = 0x8a228653b964fd48

func NewDir(s *capnp.Segment) (Dir, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 5})
	return Dir{st}, err
}

func NewRootDir(s *capnp.Segment) (Dir, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 5})
	return Dir{st}, err
}

func ReadRootDir(msg *capnp.Message) (Dir, error) {
	root, err := msg.RootPtr()
	return Dir{root.Struct()}, err
}

func (s Dir) String() string {
	str, _ := text.Marshal(0x8a228653b964fd48, s.Struct)
	return str
}

func (s Dir) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s Dir) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Dir) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s Dir) SetName(v string) error {
	return s.Struct.SetText(0, v)
}

func (s Dir) Location() (string, error) {
	p, err := s.Struct.Ptr(1)
	return p.Text(), err
}

func (s Dir) HasLocation() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Dir) LocationBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(1)
	return p.TextBytes(), err
}

func (s Dir) SetLocation(v string) error {
	return s.Struct.SetText(1, v)
}

func (s Dir) Contents() (Inode_List, error) {
	p, err := s.Struct.Ptr(2)
	return Inode_List{List: p.List()}, err
}

func (s Dir) HasContents() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
}

func (s Dir) SetContents(v Inode_List) error {
	return s.Struct.SetPtr(2, v.List.ToPtr())
}

// NewContents sets the contents field to a newly
// allocated Inode_List, preferring placement in s's segment.
func (s Dir) NewContents(n int32) (Inode_List, error) {
	l, err := NewInode_List(s.Struct.Segment(), n)
	if err != nil {
		return Inode_List{}, err
	}
	err = s.Struct.SetPtr(2, l.List.ToPtr())
	return l, err
}

func (s Dir) Parent() (string, error) {
	p, err := s.Struct.Ptr(3)
	return p.Text(), err
}

func (s Dir) HasParent() bool {
	p, err := s.Struct.Ptr(3)
	return p.IsValid() || err != nil
}

func (s Dir) ParentBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(3)
	return p.TextBytes(), err
}

func (s Dir) SetParent(v string) error {
	return s.Struct.SetText(3, v)
}

func (s Dir) Size() uint64 {
	return s.Struct.Uint64(0)
}

func (s Dir) SetSize(v uint64) {
	s.Struct.SetUint64(0, v)
}

func (s Dir) Aclkey() (string, error) {
	p, err := s.Struct.Ptr(4)
	return p.Text(), err
}

func (s Dir) HasAclkey() bool {
	p, err := s.Struct.Ptr(4)
	return p.IsValid() || err != nil
}

func (s Dir) AclkeyBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(4)
	return p.TextBytes(), err
}

func (s Dir) SetAclkey(v string) error {
	return s.Struct.SetText(4, v)
}

func (s Dir) ModificationTime() uint32 {
	return s.Struct.Uint32(8)
}

func (s Dir) SetModificationTime(v uint32) {
	s.Struct.SetUint32(8, v)
}

func (s Dir) CreationTime() uint32 {
	return s.Struct.Uint32(12)
}

func (s Dir) SetCreationTime(v uint32) {
	s.Struct.SetUint32(12, v)
}

// Dir_List is a list of Dir.
type Dir_List struct{ capnp.List }

// NewDir creates a new list of Dir.
func NewDir_List(s *capnp.Segment, sz int32) (Dir_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 5}, sz)
	return Dir_List{l}, err
}

func (s Dir_List) At(i int) Dir { return Dir{s.List.Struct(i)} }

func (s Dir_List) Set(i int, v Dir) error { return s.List.SetStruct(i, v.Struct) }

// Dir_Promise is a wrapper for a Dir promised by a client call.
type Dir_Promise struct{ *capnp.Pipeline }

func (p Dir_Promise) Struct() (Dir, error) {
	s, err := p.Pipeline.Struct()
	return Dir{s}, err
}

type UserGroup struct{ capnp.Struct }

// UserGroup_TypeID is the unique identifier for the type UserGroup.
const UserGroup_TypeID = 0xee5217621d9cbb3a

func NewUserGroup(s *capnp.Segment) (UserGroup, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	return UserGroup{st}, err
}

func NewRootUserGroup(s *capnp.Segment) (UserGroup, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	return UserGroup{st}, err
}

func ReadRootUserGroup(msg *capnp.Message) (UserGroup, error) {
	root, err := msg.RootPtr()
	return UserGroup{root.Struct()}, err
}

func (s UserGroup) String() string {
	str, _ := text.Marshal(0xee5217621d9cbb3a, s.Struct)
	return str
}

func (s UserGroup) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s UserGroup) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s UserGroup) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s UserGroup) SetName(v string) error {
	return s.Struct.SetText(0, v)
}

func (s UserGroup) IyoId() (string, error) {
	p, err := s.Struct.Ptr(1)
	return p.Text(), err
}

func (s UserGroup) HasIyoId() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s UserGroup) IyoIdBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(1)
	return p.TextBytes(), err
}

func (s UserGroup) SetIyoId(v string) error {
	return s.Struct.SetText(1, v)
}

func (s UserGroup) IyoInt() uint64 {
	return s.Struct.Uint64(0)
}

func (s UserGroup) SetIyoInt(v uint64) {
	s.Struct.SetUint64(0, v)
}

// UserGroup_List is a list of UserGroup.
type UserGroup_List struct{ capnp.List }

// NewUserGroup creates a new list of UserGroup.
func NewUserGroup_List(s *capnp.Segment, sz int32) (UserGroup_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2}, sz)
	return UserGroup_List{l}, err
}

func (s UserGroup_List) At(i int) UserGroup { return UserGroup{s.List.Struct(i)} }

func (s UserGroup_List) Set(i int, v UserGroup) error { return s.List.SetStruct(i, v.Struct) }

// UserGroup_Promise is a wrapper for a UserGroup promised by a client call.
type UserGroup_Promise struct{ *capnp.Pipeline }

func (p UserGroup_Promise) Struct() (UserGroup, error) {
	s, err := p.Pipeline.Struct()
	return UserGroup{s}, err
}

type ACI struct{ capnp.Struct }

// ACI_TypeID is the unique identifier for the type ACI.
const ACI_TypeID = 0xe7b4959415dabf9c

func NewACI(s *capnp.Segment) (ACI, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 3})
	return ACI{st}, err
}

func NewRootACI(s *capnp.Segment) (ACI, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 3})
	return ACI{st}, err
}

func ReadRootACI(msg *capnp.Message) (ACI, error) {
	root, err := msg.RootPtr()
	return ACI{root.Struct()}, err
}

func (s ACI) String() string {
	str, _ := text.Marshal(0xe7b4959415dabf9c, s.Struct)
	return str
}

func (s ACI) Uname() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s ACI) HasUname() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s ACI) UnameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s ACI) SetUname(v string) error {
	return s.Struct.SetText(0, v)
}

func (s ACI) Gname() (string, error) {
	p, err := s.Struct.Ptr(1)
	return p.Text(), err
}

func (s ACI) HasGname() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s ACI) GnameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(1)
	return p.TextBytes(), err
}

func (s ACI) SetGname(v string) error {
	return s.Struct.SetText(1, v)
}

func (s ACI) Mode() uint16 {
	return s.Struct.Uint16(0)
}

func (s ACI) SetMode(v uint16) {
	s.Struct.SetUint16(0, v)
}

func (s ACI) Rights() (ACI_Right_List, error) {
	p, err := s.Struct.Ptr(2)
	return ACI_Right_List{List: p.List()}, err
}

func (s ACI) HasRights() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
}

func (s ACI) SetRights(v ACI_Right_List) error {
	return s.Struct.SetPtr(2, v.List.ToPtr())
}

// NewRights sets the rights field to a newly
// allocated ACI_Right_List, preferring placement in s's segment.
func (s ACI) NewRights(n int32) (ACI_Right_List, error) {
	l, err := NewACI_Right_List(s.Struct.Segment(), n)
	if err != nil {
		return ACI_Right_List{}, err
	}
	err = s.Struct.SetPtr(2, l.List.ToPtr())
	return l, err
}

func (s ACI) Id() uint32 {
	return s.Struct.Uint32(4)
}

func (s ACI) SetId(v uint32) {
	s.Struct.SetUint32(4, v)
}

// ACI_List is a list of ACI.
type ACI_List struct{ capnp.List }

// NewACI creates a new list of ACI.
func NewACI_List(s *capnp.Segment, sz int32) (ACI_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 3}, sz)
	return ACI_List{l}, err
}

func (s ACI_List) At(i int) ACI { return ACI{s.List.Struct(i)} }

func (s ACI_List) Set(i int, v ACI) error { return s.List.SetStruct(i, v.Struct) }

// ACI_Promise is a wrapper for a ACI promised by a client call.
type ACI_Promise struct{ *capnp.Pipeline }

func (p ACI_Promise) Struct() (ACI, error) {
	s, err := p.Pipeline.Struct()
	return ACI{s}, err
}

type ACI_Right struct{ capnp.Struct }

// ACI_Right_TypeID is the unique identifier for the type ACI_Right.
const ACI_Right_TypeID = 0xe615914de76be38f

func NewACI_Right(s *capnp.Segment) (ACI_Right, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return ACI_Right{st}, err
}

func NewRootACI_Right(s *capnp.Segment) (ACI_Right, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return ACI_Right{st}, err
}

func ReadRootACI_Right(msg *capnp.Message) (ACI_Right, error) {
	root, err := msg.RootPtr()
	return ACI_Right{root.Struct()}, err
}

func (s ACI_Right) String() string {
	str, _ := text.Marshal(0xe615914de76be38f, s.Struct)
	return str
}

func (s ACI_Right) Right() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s ACI_Right) HasRight() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s ACI_Right) RightBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s ACI_Right) SetRight(v string) error {
	return s.Struct.SetText(0, v)
}

func (s ACI_Right) Usergroupid() uint16 {
	return s.Struct.Uint16(0)
}

func (s ACI_Right) SetUsergroupid(v uint16) {
	s.Struct.SetUint16(0, v)
}

// ACI_Right_List is a list of ACI_Right.
type ACI_Right_List struct{ capnp.List }

// NewACI_Right creates a new list of ACI_Right.
func NewACI_Right_List(s *capnp.Segment, sz int32) (ACI_Right_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return ACI_Right_List{l}, err
}

func (s ACI_Right_List) At(i int) ACI_Right { return ACI_Right{s.List.Struct(i)} }

func (s ACI_Right_List) Set(i int, v ACI_Right) error { return s.List.SetStruct(i, v.Struct) }

// ACI_Right_Promise is a wrapper for a ACI_Right promised by a client call.
type ACI_Right_Promise struct{ *capnp.Pipeline }

func (p ACI_Right_Promise) Struct() (ACI_Right, error) {
	s, err := p.Pipeline.Struct()
	return ACI_Right{s}, err
}

const schema_ae9223e76351538a = "x\xda\xacV]l\x14U\x14>\xdf\xbd\xb3;\xbb\xc9" +
	"\xd6\xddq\x96\xf0\x13\x9bU\xd4D\x1a\xdb\xd0\x02F\x1a" +
	"M\xa1\xa2R\x02\x91\xdb\x85\x18\x0d\x0fNw\xa7\xed\xb0" +
	"\xdb\xdd\xcd\xee\xacX\xa2\xa9 \"4\x12\x03\xa9?\x18" +
	"\x88\x80\xc4\xc0\x8bF\xd4\x07b\x8c?/F\xa3$\x10" +
	"I\xd0\x801Ql0JL\x84\x08\x09u\xcc\xd9i" +
	"g\xb7\xb5\xbc\xf9v\xe7\x9bs\xcf\xf7\xdd\xef\x9c{f" +
	"\x16\xdf'W\x88\xf6\xd0K!\"\xb5\"\x14\xf6\xb6\x9e" +
	"Y\xf7\xd8\xb93\x1d\xbb\xc9H\x0a/\xf7\xc0\xe0\xa9\xd7" +
	"\x8f\xb9\xe7\x89`\x1e\x17_\x9b\x1f\x09\x9d\xc8|_\x8c" +
	"\x10\xbc\xd5\x13\xd9\x93\xe9\x9d\x0bGI\xc5 \xbc\xd1\xb4" +
	"\xca\x8c\xdf\xb9\xef]\x0a\x858\xe4\x92\xd8n\xfe\xc1\xc1" +
	"K.\x89/A\xf8\xfb\x02\xfe:u\xc7\xd1\xa3F\x0c" +
	"\x0d\xa1\xe0\xd0\xd3\xda\x9b\xe6\xf7\x1a\xaf\xcej]\x04o" +
	"\xec\xe2\xe5\x0f\xb6\xbd%>\xe3\xbc\xb2!Xr\xc8U" +
	"m\x9f9\xc1\xc1K\xaek\x8f\x83\xe0\xb5\\\xfb\xc2y" +
	"!}\xe4,M\xcf\\\xd3\xb9<|\xc2\\\x19\xe6\xd5" +
	"\x83\xe1-\x84\xfaiT\x0c\xff\x91\xb1?|\xc4<\x1c" +
	"\x9eKd\x1e\xaf\x05oz\xcdz\xe7\xe2\xa1\xf9\xbf\xd0" +
	",\x92C\xfa\xa8\xd9\xa4\xf3*\xaa\xb3\xe4W~\xce\x8d" +
	"\xaf\xdb;\xe7WR\x09\xc0;\xf0\xe9\x0fs\xc6^\xfd" +
	"p|2x\x91~\xc2l\xaf\x05\xb7\xea\x9c8x=" +
	"CE\xed|{\xf4\xed\xe6^}.\xd1\x92\xfdz\x8a" +
	"\xcf\xb7\xe0\x89+Kw\xbe=\xf1\xfb\xac\x9aOFF" +
	"\xcd\xcf#\xbc\xfa$\xc2\xa9;?>\xd0\xdc7\xb7\xf7" +
	"\xf2\xcc\xe0\x9a\x1b\xad\xd1\x13\xe6\xb2(\xaf\xda\xa3\xef\x11" +
	"\xbc\xd8\x8f\x9d\xa5\xe5\xd5\xcauR\xb7B\xd6]\xdf(" +
	"uh\xd0\xcc\xaf\xa2\xbf\x11\xcco\xa3\xe3\xd4\xea\x0d\x15" +
	"\xb3v\xbe-c\x89R\xa1\xd4\x99.\xd9\x19\xc7\xca\xb7" +
	"m\x18.\xd9D\xeb\x01\x95\x84 2\x96u\x12\x01F" +
	"k\x07\x11\x84qw7\x11\xa4\xd1\xbc\x86\x08\x9a1\xbf" +
	"\x9b\xa8\xabR\xcc\xe4l7\xd5\x97/fr#\x99A" +
	"\xab\x9c\xb5\x9f\xf6\xfa\x9d\xfeb\xc9\xe1L4R-\xe4" +
	"\x0a\xc5-\x85\x80\x0eL\xb7\xca)\xd7Hn\x97\x1a\x91" +
	"\x06\"\xe3t\x0b\x91\xfaFB\x9d\x130\x80$\x18<" +
	"\xbb\x86H}'\xa1~\x120\x84\xf0%]`\xf0\xbc" +
	"\x84\xba&`H\x99\x84$2\xaev\x12\xa9?%\xd4" +
	"\x0d\x01hIhD\xc6uNyE\xa2\x17\x02FH" +
	"K\"DdLp\xe05\x89\xb4\xc6hX$\x11&" +
	"2\x81Q\xa2\xb4\x06\x89t\x82q]&k\xb5h\xc2" +
	"f\xa2t\x8c\xf1y\x10\x88\x17\xac!\x1b1\x12\x88\x11" +
	"\xbc|1c\xb9N\xb1@D\x01\x96)\x16\\\xbb\xe0" +
	"V\x18\xbb\x85\xb0^\x02\x89z\x19\x08\x0cv\x95\xac\xb2" +
	"]p\xa7\xf6\xc4+\xceV\x1bQ\x12\x88\x12\xba\xacL" +
	">g\x0f\x07\xf9\x86\x8aY\xa7\xdf\xc9X`\xa2\x0d\xce" +
	"\x90M\x84\x08\x09D\x98\xabl\xfb\xfcq~\x11\xc0\xd3" +
	"\x8cNW\xfbVI\xa7\xcc^k\x81\xd7M\x0b\x89T" +
	"DB%\x05\xf4\x19l\xf5\xad=\x85b\x166\xef\x9c" +
	"\x17\xec\xdc\xcf\x96\x8eI\xa8C\x02SE:\xc8\xd8\x1b" +
	"\x12\xea\xa8\x00\x04\x1az\xd08\xfc$\x09C\xfa\x1e\x1b" +
	"{\xd8\xf9]\x12jL\xc0\xd0|\x83\x8d\xbd\xa3\xf5\x84" +
	"\xb5\x1aE8\xe3f\"u@B\x1d\x9ba\xf94\xab" +
	"<\xcbu\xcbN_\xd5%iW\xfeo\xdf\x1eq\xf2" +
	"v\xaa\x9b\x9b\x9a\x0d\x88\x04\x06,\xe2\xc3\xde%\xa1\x16" +
	"7\xb4i+\xfby\x8f\x84Z*\x10\x1f\xb4*\x83h" +
	"\"\x81&\xf2\xcd\x9d\\\xcf\xa8K\xc9\xce\xe8\x8e\x95W" +
	"\x1a\xd00\x9d\xd1\x12\xe7\x1bx\x13\xc2\x80\xaf\xa5\x81\xcf" +
	"\x1d.\xd9\x88\xd7s\x10\x10'\xc4\xb3\x96k\xcdN\xbd" +
	"\xd6)\xe4\xfc+\xde\xd0\x11\x9d\xf5\x8e\xe8r\xad\xf2\x80" +
	"\xed\xce\xde\x14+\x1f\xeaiK\xf5:\x03\x83\xee\x0c_" +
	":f\x91\xd9G\xa4\xee\x95P\xf7\x0b\xa4\xca\xbc'\xc8" +
	"Y\xad\xd8\xe5\x81r\xb1Jz\xc9\xc9B'\x01}\x16" +
	"&\xaa\xb9\x13\xcca\x03\x1d>\xb5J\x06\xbc\xcf1\xef" +
	"3\x12jGC=\xb61\xf8\xac\x84\xda\xc5\x1d\xe9O" +
	"\x8d\x17\xd9\xb3\xe7%\xd4\xcb<5\x84?5v\xf3\xb1" +
	"wLv\x9f\x06\x7fl\x1c\\P\xef\xe7T\xb5\xb1\xfd" +
	"R\x03\xd3\x9a\x91\xe5Ni\xef\xaa\x9d\xafR\xbf\xf3\x81" +
	"j\xff\xceK'{\xf3F\xf3\xab\xd1`fo\xbd\xbe" +
	"Sgj\xef\x9c4s\xb5\x80W\x9b\xb6ig+\xa1" +
	".\xa0\x865\x08\x08\xbe\xa1\xbe\x80\xe9\xac\x1b+v9" +
	"\xf5h\xb9X-1s,`~\x98]Z!\xa1\xd6" +
	"6\xd8\xd9\xc3v\xae\x92P\xeb\xebv\xaec9\xab%" +
	"\xd4\x86\x19W4\xe5\x0c\x17{\xb2SO]\xfcTp" +
	"\x83+;\xed\x9b\xc3\x03\xc6n\xf3oq\xbc\xea\xda\x15" +
	"\x95\x90\xdam\x9e\x07\x9f\xc2\xe2[\xb5IB\x0d\x0a4" +
	"\xe3\x1f\x86\xb9f6K|JB\xe5\x05\x9a\xc5\x847" +
	"Y5\x87\xe1\xac\x84*\x094\xcb\x1b\x0c\xf3\xb8\x1f\xea" +
	"&R\x83\x12\xca\x15\xd0\xb3N\x19\x89\xa9\x9f\x16\x02\x12" +
	"\x84x\xbf\x93\xb7\x91\xa8\x7f\x91'\xe1\xbcS\xc8!Q" +
	"\xff_\xf0\xe1\x91\x8a\xff\x99D\xa2\xf1'\x8a\xdf\xfc\x1b" +
	"\x00\x00\xff\xffR12\x1d"

func init() {
	schemas.Register(schema_ae9223e76351538a,
		0x8932d2d84f4dd27a,
		0x8a228653b964fd48,
		0xa4a421ce00f301dd,
		0xc0029f81b3eee594,
		0xd5a2538369c2f82a,
		0xdc74a897ce683c6b,
		0xe419a0e5a661965c,
		0xe615914de76be38f,
		0xe7b4959415dabf9c,
		0xecfda38634f4591a,
		0xee5217621d9cbb3a,
		0xf9737539703ade0c)
}
