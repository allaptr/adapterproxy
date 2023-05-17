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

var globalKey string

func BenchmarkMake(b *testing.B) {
	b.Run("buffer", benchmarkMakeKey)
	b.Run("string", benchmarkMakeKeyStr)
}

func benchmarkMakeKey(b *testing.B) {
	var key string
	for i := 0; i < b.N; i++ {
		key = MakeKey("us", "foo1234")
	}
	globalKey = key
}

func benchmarkMakeKeyStr(b *testing.B) {
	var key string
	for i := 0; i < b.N; i++ {
		key = MakeKeyStr("us", "foo1234")
	}
	globalKey = key
}

// $ go test -run none -bench . -benchtime 3s -benchmem
// goos: linux
// goarch: amd64
// pkg: backendify/cache
// cpu: Intel(R) Core(TM) i7-7500U CPU @ 2.70GHz
// BenchmarkMakeKey-4               8084622               407.1 ns/op            80 B/op          4 allocs/op
// BenchmarkMakeKeyStr-4           17240190               191.5 ns/op            48 B/op          3 allocs/op
// PASS
// ok      backendify/cache        9.172s

// go test -run none -bench BenchmarkMake/buffer -benchtime 3s -benchmem
