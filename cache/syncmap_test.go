package cache

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_syn_put_get(t *testing.T) {
	cache := NewCache()
	td := []struct {
		name     string
		key      string
		value    []byte
		expected []byte
	}{
		{"Put-Get data", "us-SFO", []byte("EWR"), []byte("EWR")},
		{"Put-Get empty", "us-SFO", nil, nil},
	}
	for _, td := range td {
		t.Run(td.name, func(t *testing.T) {
			cache.Put(td.key, td.value)
			val := cache.Get(td.key)
			assert.Equal(t, td.expected, val)
		})
	}
}

func Test_syn_Access_Put(t *testing.T) {
	cache := NewCache()
	var wg *sync.WaitGroup = new(sync.WaitGroup)
	wg.Add(100)
	for i := 0; i < 100; i++ {
		key := strconv.Itoa(i)
		go syncput(wg, cache, key, []byte("Random stuff "+key))
	}
	wg.Wait()
	val := cache.Get("77")
	assert.Equal(t, []byte("Random stuff 77"), val)
	val = cache.Get("277")
	assert.Nil(t, val)
}

func syncput(w *sync.WaitGroup, cache Cache, key string, val []byte) {
	defer w.Done()
	cache.Put(key, val)
}
