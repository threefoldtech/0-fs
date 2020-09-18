package meta

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"strconv"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	// import sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	np "github.com/threefoldtech/0-fs/cap.np"
	"golang.org/x/crypto/blake2b"
	capnp "zombiezen.com/go/capnproto2"
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

	//SQLiteDBName is the name of the sqlite3 database stored in an flist
	SQLiteDBName = "flistdb.sqlite3"

	// DirCacheSize defines the size of the LRU cache for the directories
	DirCacheSize = 1024
	// AccessCacheSize defines the size of the LRU cache for ACL
	AccessCacheSize = 64
)

//NewStore creates a new meta store with path p
func NewStore(p string) (Store, error) {
	p = path.Join(p, SQLiteDBName)
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?mode=ro", p))
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare("select value from entries where key = ?")
	if err != nil {
		return nil, err
	}

	cache, err := lru.New(DirCacheSize)
	if err != nil {
		return nil, err
	}
	aclCache, err := lru.New(AccessCacheSize)
	if err != nil {
		return nil, err
	}

	return &sqlStore{
		db:    db,
		stmt:  stmt,
		cache: cache,
		acl:   aclCache,

		users:  make(map[string]int),
		groups: make(map[string]int),
	}, nil
}

type sqlStore struct {
	stmt *sql.Stmt
	db   *sql.DB

	//cache  *cache.Cache
	cache *lru.Cache
	acl   *lru.Cache

	users  map[string]int
	groups map[string]int

	usersM  sync.Mutex
	groupsM sync.Mutex
}

func (s *sqlStore) hash(path string) string {
	hasher, _ := blake2b.New(16, nil)
	io.WriteString(hasher, path)

	return fmt.Sprintf("%x", hasher.Sum(nil))
}

//getACI gets aci object with key from db
func (s *sqlStore) getACI(key string) (*np.ACI, error) {
	if aci, ok := s.acl.Get(key); ok {
		return aci.(*np.ACI), nil
	}

	row := s.stmt.QueryRow(key)
	var data []byte

	if err := row.Scan(&data); err != nil {
		if err == sql.ErrNoRows {
			return nil, errNoACI
		}
		return nil, err
	}

	msg, err := capnp.NewDecoder(bytes.NewBuffer(data)).Decode()
	if err != nil {
		log.Debugf("failed to get msg for aci %s: %s", key, err)
		return nil, err
	}
	msg.TraverseLimit = TraverseLimit
	aci, err := np.ReadRootACI(msg)
	if err != nil {
		return nil, err
	}

	s.acl.Add(key, &aci)
	return &aci, nil
}

func (s *sqlStore) lookUpUser(name string) int {
	s.usersM.Lock()
	defer s.usersM.Unlock()

	if id, ok := s.users[name]; ok {
		return id
	}
	uid := 1000
	if u, err := user.Lookup(name); err == nil {
		if id, err := strconv.ParseInt(u.Uid, 10, 32); err == nil {
			uid = int(id)
		}
	}

	s.users[name] = uid
	return uid
}

func (s *sqlStore) lookUpGroup(name string) int {
	s.groupsM.Lock()
	defer s.groupsM.Unlock()

	if id, ok := s.groups[name]; ok {
		return id
	}
	gid := 1000
	if g, err := user.LookupGroup(name); err == nil {
		if id, err := strconv.ParseInt(g.Gid, 10, 32); err == nil {
			gid = int(id)
		}
	}

	s.groups[name] = gid
	return gid
}

//getAccess gets access object from db
func (s *sqlStore) getAccess(key string) (Access, error) {
	aci, err := s.getACI(key)
	if err != nil {
		log.Debugf("failed to get aci for key %s: %s", key, err)
		return DefaultAccess, err
	}

	uid := aci.Uid()
	gid := aci.Gid()

	if uid == -1 {
		uname, _ := aci.Uname()
		uid = int64(s.lookUpUser(uname))
	}

	if gid == -1 {
		gname, _ := aci.Gname()
		gid = int64(s.lookUpGroup(gname))
	}

	mode := uint32(aci.Mode())
	return Access{
		Mode: uint32(os.ModePerm) & mode,
		UID:  uint32(uid),
		GID:  uint32(gid),
	}, nil
}

//getDir gets dir entry from db
func (s *sqlStore) getDir(path string) (*Dir, error) {
	if path == "." {
		path = ""
	}

	hash := s.hash(path)
	log.Debugf("hash(%s) == %s", path, hash)
	return s.getDirWithHash(hash)
}

//getDir gets dir entry from db
func (s *sqlStore) getDirWithHash(hash string) (*Dir, error) {
	row := s.stmt.QueryRow(hash)
	var data []byte
	if err := row.Scan(&data); err != nil {
		return nil, err
	}

	//we need to load this slice as a capnpn dir
	msg, err := capnp.NewDecoder(bytes.NewBuffer(data)).Decode()
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

func (s *sqlStore) get(p string) (Meta, error) {
	if m, ok := s.cache.Get(p); ok {
		return m.(Meta), nil
	}

	dir, err := s.getDir(p)
	if err == nil {
		s.cache.Add(p, dir)
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
			break
		}
	}

	if meta != nil {
		return meta, nil
	}

	return nil, ErrNotFound
}

func (s *sqlStore) Get(path string) (Meta, bool) {
	meta, err := s.get(path)
	if err != nil {
		return nil, false
	}
	return meta, true
}
