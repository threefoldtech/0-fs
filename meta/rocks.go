package meta

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/codahale/blake2"
	"github.com/zero-os/0-fs/cap.np"
	"github.com/patrickmn/go-cache"
	rocksdb "github.com/tecbot/gorocksdb"
	"io"
	"os"
	"os/user"
	"path"
	"strconv"
	"sync"
	"time"
	"zombiezen.com/go/capnproto2"
)

const (
	TraverseLimit = ^uint64(0)
)

func NewRocksMeta(ns string, dbpath string) (MetaStore, error) {
	opt := rocksdb.NewDefaultOptions()
	db, err := rocksdb.OpenDbForReadOnly(opt, dbpath, true)
	if err != nil {
		return nil, err
	}

	return &rocksMetaStore{
		ns:    ns,
		db:    db,
		ro:    rocksdb.NewDefaultReadOptions(),
		cache: cache.New(5*time.Second, 1*time.Second),
	}, nil
}

type rocksMeta struct {
	name  string
	dir   *np.Dir
	inode np.Inode

	store *rocksMetaStore
	o     sync.Once

	blks []BlockInfo
	bo   sync.Once
}

func (rm *rocksMeta) load() {
	//load runs once and only take effect if
	//the meta wasn't already filed in by a previous Meta node.

	rm.o.Do(func() {
		if rm.dir != nil && rm.dir.HasData() {
			return
		}

		if loaded, ok := rm.store.Get(rm.name); ok {
			self := loaded.(*rocksMeta)
			rm.dir = self.dir
			rm.inode = self.inode
		}
	})
}

func (rm *rocksMeta) IsDir() bool {
	rm.load()
	return !rm.inode.HasData()
}

func (rm *rocksMeta) String() string {
	return rm.name
}

func (rm *rocksMeta) ID() string {
	m := md5.New()
	for _, blk := range rm.Blocks() {
		m.Write(blk.Key)
	}
	return fmt.Sprintf("%x", m.Sum(nil))
}

//base name
func (rm *rocksMeta) Name() string {
	return path.Base(rm.name)
}

func (rm *rocksMeta) blocks() {
	rm.load()
	var blocks []BlockInfo
	if !rm.inode.HasData() {
		return
	}

	attrs := rm.inode.Attributes()
	if !attrs.HasFile() {
		return
	}
	file, _ := attrs.File()
	if !file.HasBlocks() {
		return
	}

	cblocks, _ := file.Blocks()
	for i := 0; i < cblocks.Len(); i++ {
		block := cblocks.At(i)

		hash, _ := block.Hash()
		key, _ := block.Key()
		blocks = append(blocks, BlockInfo{
			Key:      hash,
			Decipher: key,
		})
	}

	rm.blks = blocks
}

func (rm *rocksMeta) Blocks() []BlockInfo {
	rm.bo.Do(rm.blocks)
	return rm.blks
}

func (rm *rocksMeta) Children() []Meta {
	rm.load()
	//if that is a file, we must have no children
	var children []Meta
	if rm.inode.HasData() {
		return children
	}

	if !rm.dir.HasContents() {
		return children
	}

	contents, _ := rm.dir.Contents()

	for i := 0; i < contents.Len(); i++ {
		inode := contents.At(i)
		name, _ := inode.Name()
		var child *rocksMeta
		if inode.Attributes().HasDir() {
			//we don't set the dir, to force it to load it's content when
			//it has too.
			child = &rocksMeta{
				store: rm.store,
				name:  path.Join(rm.name, name),
			}
		} else {
			//here we already have everything we know about this node (file)
			//so we just create a complete struct.
			child = &rocksMeta{
				store: rm.store,
				name:  path.Join(rm.name, name),
				dir:   rm.dir,
				inode: inode,
			}
		}

		children = append(children, child)
	}

	return children
}

func (rm *rocksMeta) Info() MetaInfo {
	rm.load()

	if !rm.inode.HasData() {
		//that must be a dir.

		aciKey, _ := rm.dir.Aclkey()
		access := rm.store.getAccess(aciKey)

		return MetaInfo{
			Type:             DirType,
			Size:             rm.dir.Size(),
			CreationTime:     rm.dir.CreationTime(),
			ModificationTime: rm.dir.ModificationTime(),
			Access:           access,
		}
	}

	aciKey, _ := rm.inode.Aclkey()
	access := rm.store.getAccess(aciKey)

	info := MetaInfo{
		CreationTime:     rm.inode.CreationTime(),
		ModificationTime: rm.inode.ModificationTime(),
		Access:           access,
	}

	attrs := rm.inode.Attributes()
	if !attrs.HasData() {
		log.Errorf("'%s' attributes is empty", attrs)
		return info
	}

	if attrs.HasFile() {
		file, _ := attrs.File()
		info.Type = RegularType
		info.Size = rm.inode.Size()
		info.FileBlockSize = uint64(file.BlockSize()) * 4096 //Block size is actually the number of 4K blocks in a file
	} else if attrs.HasLink() {
		link, _ := attrs.Link()
		info.Type = LinkType
		target, _ := link.Target()
		info.LinkTarget = target
	} else if attrs.HasSpecial() {
		special, _ := attrs.Special()
		switch special.Type() {
		case np.Special_Type_block:
			info.Type = BlockDeviceType
		case np.Special_Type_chardev:
			info.Type = CharDeviceType
		case np.Special_Type_fifopipe:
			info.Type = FIFOType
		case np.Special_Type_socket:
			info.Type = SocketType
		}

		if special.HasData() {
			bytes, _ := special.Data()
			info.SpecialData = string(bytes)
		}
	}

	return info
}

type rocksMetaStore struct {
	ns    string
	db    *rocksdb.DB
	ro    *rocksdb.ReadOptions
	cache *cache.Cache
}

func (rs *rocksMetaStore) hash(namespace string, path string) (string, error) {
	bl2b := blake2.New(&blake2.Config{
		Size: 32,
	})

	_, err := io.WriteString(bl2b, fmt.Sprintf("%s%s", namespace, path))
	if err != nil {
		return "", err
	}

	hash := fmt.Sprintf("%x", bl2b.Sum(nil))
	if namespace != "" {
		return fmt.Sprintf("%s:%s", namespace, hash), nil
	} else {
		return hash, nil
	}
}

func (rs *rocksMetaStore) dirFromSlice(slice *rocksdb.Slice) (*np.Dir, error) {
	msg, err := capnp.NewDecoder(bytes.NewBuffer(slice.Data())).Decode()
	if err != nil {
		return nil, err
	}
	msg.TraverseLimit = TraverseLimit
	dir, err := np.ReadRootDir(msg)
	if err != nil {
		return nil, err
	}

	return &dir, nil
}

func (rs *rocksMetaStore) aciFromSlice(slice *rocksdb.Slice) (*np.ACI, error) {
	msg, err := capnp.NewDecoder(bytes.NewBuffer(slice.Data())).Decode()
	if err != nil {
		return nil, err
	}
	msg.TraverseLimit = TraverseLimit
	aci, err := np.ReadRootACI(msg)
	if err != nil {
		return nil, err
	}

	return &aci, nil
}

func (rs *rocksMetaStore) getACI(key string) (*np.ACI, error) {
	cKey := fmt.Sprintf("accesskey.%x", key)
	if aci, ok := rs.cache.Get(cKey); ok {
		return aci.(*np.ACI), nil
	}

	slice, err := rs.db.Get(rs.ro, []byte(key))
	if err != nil {
		return nil, err
	}

	aci, err := rs.aciFromSlice(slice)
	if err != nil {
		return nil, err
	}

	rs.cache.Set(cKey, aci, cache.DefaultExpiration)
	return aci, nil
}

func (rs *rocksMetaStore) getAccess(key string) Access {
	aci, err := rs.getACI(key)
	if err != nil {
		return Access{
			Mode: 0400,
			UID:  1000,
			GID:  1000,
		}
	}

	uname, _ := aci.Uname()
	gname, _ := aci.Gname()
	mode := uint32(aci.Mode())

	uid := 1000
	gid := 1000

	if u, err := user.Lookup(uname); err == nil {
		if id, err := strconv.ParseInt(u.Uid, 10, 32); err != nil {
			uid = int(id)
		}
	}

	if g, err := user.LookupGroup(gname); err == nil {
		if id, err := strconv.ParseInt(g.Gid, 10, 32); err != nil {
			gid = int(id)
		}
	}

	return Access{
		Mode: uint32(os.ModePerm) & mode,
		UID:  uint32(uid),
		GID:  uint32(gid),
	}
}

func (rs *rocksMetaStore) get(name string, level int) (*rocksMeta, bool) {
	if level == 0 {
		return nil, false
	}

	if name == "." {
		name = ""
	}

	if obj, ok := rs.cache.Get(name); ok {
		log.Debugf("cache hit for name: '%s' (%d)", name, level)
		return obj.(*rocksMeta), true
	}

	hash, _ := rs.hash(rs.ns, name)
	slice, err := rs.db.Get(rs.ro, []byte(hash))
	if err != nil {
		log.Errorf("no entry for '%s': %s", name, err)
		return nil, false
	}

	defer slice.Free()

	if slice.Size() != 0 {
		dir, err := rs.dirFromSlice(slice)

		if err != nil {
			log.Errorf("failed to get dir entry: %s", err)
			return nil, false
		}

		return &rocksMeta{
			store: rs,
			name:  name,
			dir:   dir,
		}, true
	}

	return rs.get(path.Dir(name), level-1)
}

func (rs *rocksMetaStore) Get(name string) (Meta, bool) {
	/*
		When we try to hit an object, this object can be either a directory or any other file type.
		The rocks-db has all values as directories, which means if we try to retrieve an object and we
		didn't find it, we should try to retrieve it's parent. Then find the file name in that directory entry.
	*/

	//we search only 2 levels
	meta, ok := rs.get(name, 2)
	if !ok {
		return nil, false
	}

	//direct hit. we return
	if meta.name == name {
		rs.cache.Set(name, meta, cache.DefaultExpiration)
		return meta, true
	}

	log.Debugf("searching children of '%s'", meta.name)
	//searching the directory contents.
	if !meta.dir.HasContents() {
		return nil, false
	}

	contents, _ := meta.dir.Contents()
	base := path.Base(name)

	for i := 0; i < contents.Len(); i++ {
		inode := contents.At(i)
		nodeName, _ := inode.Name()
		if nodeName == base {
			log.Debugf("found %s in children of %s", nodeName, meta.name)
			nodeMeta := &rocksMeta{
				name:  name,
				dir:   meta.dir,
				inode: inode,
				store: rs,
			}
			//cache it
			rs.cache.Set(name, nodeMeta, cache.DefaultExpiration)
			return nodeMeta, true
		}
	}

	return nil, false
}
