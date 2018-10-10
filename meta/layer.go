package meta

import "sync"

type stores []MetaStore

type mergedDir struct {
	Meta
	lower []Meta

	merged []Meta
	o      sync.Once
}

func (m *mergedDir) children() {
	set := make(map[string]Meta)
	for _, child := range m.Meta.Children() {
		set[child.Name()] = child
	}

	for _, l := range m.lower {
		for _, child := range l.Children() {
			if _, ok := set[child.Name()]; !ok {
				set[child.Name()] = child
			}
		}
	}
	m.merged = make([]Meta, 0, len(set))
	for _, child := range set {
		m.merged = append(m.merged, child)
	}
}

func (m *mergedDir) Children() []Meta {
	m.o.Do(m.children)
	return m.merged
}

func (s stores) getMerge(p string, top Meta, under []MetaStore) Meta {
	var lower []Meta
	for _, store := range under {
		if m, ok := store.Get(p); ok {
			lower = append(lower, m)
		}
	}

	return &mergedDir{Meta: top, lower: lower}
}

func (s stores) Get(p string) (Meta, bool) {
	for i, store := range s {
		m, ok := store.Get(p)
		if !ok {
			continue
		}

		if !m.IsDir() {
			//we hit a file, then we should return
			return m, true
		}

		if i+1 == len(s) {
			//no lower layers
			return m, true
		}

		//a directory
		return s.getMerge(p, m, s[i+1:]), true
	}

	return nil, false
}

// Layered return a meta store that layer the given stores in a way that last store is on top
// Example:
//  store = Layered(s1, s2)
//  store.Get(p) will search s2 first, then s1
func Layered(store ...MetaStore) MetaStore {
	if len(store) == 1 {
		return store[0]
	}

	var s stores
	//reverse order
	for i := len(store) - 1; i >= 0; i-- {
		if store[i] == nil {
			continue
		}

		s = append(s, store[i])
	}
	return s
}
