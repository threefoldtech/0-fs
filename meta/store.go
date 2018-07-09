package meta

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"strconv"
	"time"

	"github.com/codahale/blake2"
	"github.com/patrickmn/go-cache"
	rocksdb "github.com/tecbot/gorocksdb"
	np "github.com/zero-os/0-fs/cap.np"
	"zombiezen.com/go/capnproto2"
)

var (
	//DefaultAccess fallback
	DefaultAccess = Access{
		Mode: 0400,
		UID:  1000,
		GID:  1000,
	}

	errNoACI = errors.New("no ACI attached with this object")
)

const (
	//TraverseLimit capnp message traverse limit
	TraverseLimit = ^uint64(0)
)

//NewRocksStore creates a new rocks store with namespace ns, and path p
func NewRocksStore(ns, p string) (MetaStore, error) {
	opt := rocksdb.NewDefaultOptions()
	db, err := rocksdb.OpenDbForReadOnly(opt, p, false)
	if err != nil {
		return nil, err
	}

	ro := rocksdb.NewDefaultReadOptions()

	return &rocksStore{
		db:    db,
		ro:    ro,
		ns:    ns,
		cache: cache.New(60*time.Second, 20*time.Second),
	}, nil
}

type rocksStore struct {
	db *rocksdb.DB
	ro *rocksdb.ReadOptions
	ns string

	cache *cache.Cache
}

func (s *rocksStore) hash(path string) string {
	bl2b := blake2.New(&blake2.Config{
		Size: 32,
	})
	io.WriteString(bl2b, fmt.Sprintf("%s%s", s.ns, path))

	hash := fmt.Sprintf("%x", bl2b.Sum(nil))
	if s.ns != "" {
		return fmt.Sprintf("%s:%s", s.ns, hash)
	}

	return hash
}

//getACI gets aci object with key from db
func (s *rocksStore) getACI(key string) (*np.ACI, error) {
	slice, err := s.db.Get(s.ro, []byte(key))
	if err != nil {
		log.Debugf("failed to get slice for aci %s: %s", key, err)
		return nil, err
	}

	if slice.Size() == 0 {
		// no ACI attached with the object
		return nil, errNoACI
	}

	msg, err := capnp.NewDecoder(bytes.NewBuffer(slice.Data())).Decode()
	if err != nil {
		log.Debugf("failed to get msg for aci %s: %s", key, err)
		return nil, err
	}
	msg.TraverseLimit = TraverseLimit
	aci, err := np.ReadRootACI(msg)
	if err != nil {
		return nil, err
	}

	return &aci, nil
}

//getAccess gets access object from db
func (s *rocksStore) getAccess(key string) (Access, error) {
	aci, err := s.getACI(key)
	if err != nil {
		log.Debugf("failed to get aci for key %s: %s", key, err)
		return DefaultAccess, err
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
	}, nil
}

//getDir gets dir entry from db
func (s *rocksStore) getDir(path string) (*Dir, error) {
	if path == "." {
		path = ""
	}

	hash := s.hash(path)
	log.Debugf("hash(%s) == %s", path, hash)
	return s.getDirWithHash(hash)
}

//getDir gets dir entry from db
func (s *rocksStore) getDirWithHash(hash string) (*Dir, error) {
	slice, err := s.db.Get(s.ro, []byte(hash))
	if err != nil {
		log.Debugf("failed to get slice for hash %s: %s", hash, err)
		return nil, err
	}

	defer slice.Free()

	//we need to load this slice as a capnpn dir
	msg, err := capnp.NewDecoder(bytes.NewBuffer(slice.Data())).Decode()
	if err != nil {
		log.Debugf("failed to get msg from slice %s: %s", hash, err)
		return nil, err
	}

	dir, err := np.ReadRootDir(msg)
	if err != nil {
		log.Debugf("failed to get dir from msg %s: %s", hash, err)
		return nil, err
	}

	key, err := dir.Aclkey()
	if err != nil {
		log.Debugf("failed to get acl from dir %s: %s", hash, err)
		return nil, err
	}

	access, err := s.getAccess(key)
	if err != nil && err != errNoACI {
		log.Debugf("failed to get access from key %s: %s", key, err)
		return nil, err
	}

	return &Dir{Dir: dir, store: s, access: access}, nil
}

func (s *rocksStore) get(p string) (Meta, error) {
	if m, ok := s.cache.Get(p); ok {
		return m.(Meta), nil
	}

	dir, err := s.getDir(p)
	if err == nil {
		s.cache.Set(p, dir, cache.DefaultExpiration)
		return dir, nil
	}

	if p == "" {
		//Should not reach here, unless flist is broken
		//avoid inifinte recursion
		return nil, ErrNotFound
	}

	parentPath := path.Dir(p)
	if parentPath == "." {
		parentPath = ""
	}

	parent, err := s.get(parentPath)
	if err != nil {
		return nil, err
	}

	name := path.Base(p)
	var meta Meta
	for _, child := range parent.Children() {
		if child.Name() == name {
			meta = child
		}
		s.cache.Set(path.Join(parentPath, child.Name()), child, cache.DefaultExpiration)
	}

	if meta != nil {
		return meta, nil
	}

	return nil, ErrNotFound
}

func (s *rocksStore) Get(path string) (Meta, bool) {
	meta, err := s.get(path)
	if err != nil {
		log.Debugf("cannot resolve %s: %s", path, err)
		return nil, false
	}
	return meta, true
}
