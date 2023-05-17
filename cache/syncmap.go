package cache

import (
	"fmt"
	"strings"
	"sync"
)

type Cache interface {
	Put(key string, value []byte)
	// returns value and timestamp unix
	Get(key string) []byte
}

type cache struct {
	store sync.Map
}

func NewCache() *cache {
	return &cache{
		store: sync.Map{},
	}
}

func (r *cache) Put(key string, value []byte) {
	r.store.Store(key, value)
}
func (r *cache) Get(key string) []byte {
	val, ok := r.store.Load(key)
	if !ok {
		return nil
	}
	var value []byte
	value, ok = val.([]byte)
	if !ok {
		return nil
	}
	return value
}

func MakeKey(country, id string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s-%s", country, id)
	return b.String()
}

func MakeKeyStr(country, id string) string {
	s := fmt.Sprintf("%s-%s", country, id)
	return s
}
