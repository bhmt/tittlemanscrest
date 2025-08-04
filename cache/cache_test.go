package cache

import (
	"crypto/rand"
	"math"
	"math/big"
	"slices"
	"testing"
	"time"
)

func getRand(tb testing.TB) int64 {
	out, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		tb.Fatal(err)
	}
	return out.Int64()
}

func BenchmarkCache(b *testing.B) {
	lru, _ := New[int64, int64](8192, time.Millisecond*10)

	trace := make([]int64, b.N*2)
	for i := range b.N * 2 {
		trace[i] = getRand(b) % 32768
	}

	b.ResetTimer()

	var hit, miss int
	for i := range 2 * b.N {
		if i%2 == 0 {
			lru.Add(trace[i], trace[i])
		} else {
			if _, ok := lru.Get(trace[i]); ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(hit+miss))
}

func BenchmarkCacheNoEvict(b *testing.B) {
	lru, _ := New[int64, int64](8192, 0)

	trace := make([]int64, b.N*2)
	for i := range b.N * 2 {
		trace[i] = getRand(b) % 32768
	}

	b.ResetTimer()

	var hit, miss int
	for i := range 2 * b.N {
		if i%2 == 0 {
			lru.Add(trace[i], trace[i])
		} else {
			if _, ok := lru.Get(trace[i]); ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(hit+miss))
}

func BenchmarkCacheParallel(b *testing.B) {
	lru, _ := New[int, struct{}](8192, time.Millisecond*10)
	trace := make([]int, b.N*2)
	for i := range b.N * 2 {
		trace[i] = int(getRand(b) % 32768)
	}

	b.ResetTimer()

	var hit, miss int
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				lru.Add(trace[i], struct{}{})
			} else {
				if _, ok := lru.Get(trace[i]); ok {
					hit++
				} else {
					miss++
				}
			}
			i++
		}
	})
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(hit+miss))
}

func BenchmarkCacheNoEvictParallel(b *testing.B) {
	lru, _ := New[int, struct{}](8192, 0)
	trace := make([]int, b.N*2)
	for i := range b.N * 2 {
		trace[i] = int(getRand(b) % 32768)
	}

	b.ResetTimer()

	var hit, miss int
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				lru.Add(trace[i], struct{}{})
			} else {
				if _, ok := lru.Get(trace[i]); ok {
					hit++
				} else {
					miss++
				}
			}
			i++
		}
	})
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(hit+miss))
}

func TestCacheAdd(t *testing.T) {
	ttl := time.Duration(time.Hour)
	lru, err := New[int, struct{}](5, ttl)
	if err != nil {
		t.Error(err)
	}

	for i := range 10 {
		lru.Add(i+1, struct{}{})
	}

	wantKeys := []int{6, 7, 8, 9, 10}
	gotkeys := make([]int, 0, len(lru.m))
	for k := range lru.m {
		gotkeys = append(gotkeys, k)
	}

	slices.Sort(gotkeys)
	if !slices.Equal(wantKeys, gotkeys) {
		t.Errorf("cache Add keys mismatch, want %v, got %v", wantKeys, gotkeys)
	}
}

func TestCacheGet(t *testing.T) {
	ttl := time.Duration(time.Hour)

	lru, err := New[int, struct{}](5, ttl)
	if err != nil {
		t.Error(err)
	}

	lru.Add(1, struct{}{})

	tests := []struct {
		key  int
		want bool
	}{
		{key: 1, want: true},
		{key: 2, want: false},
	}

	for _, test := range tests {
		_, gotOk := lru.Get(test.key)
		if test.want != gotOk {
			t.Errorf("cache Get not ok for key %v, want %v got %v", test.key, test.want, gotOk)
		}
	}
}

func TestCacheEvict(t *testing.T) {
	lru, err := New[int, struct{}](1, time.Duration(time.Millisecond*100))
	if err != nil {
		t.Error(err)
	}

	lru.Add(1, struct{}{})
	sizeBefore := len(lru.m)
	time.Sleep(time.Duration(time.Millisecond * 150))
	sizeAfter := len(lru.m)

	if sizeBefore != 1 || sizeAfter != 0 {
		t.Error("cache evict not ok")
	}
}

func TestCacheNoEvict(t *testing.T) {
	lru, err := New[int, struct{}](1, time.Duration(0))
	if err != nil {
		t.Error(err)
	}

	lru.Add(1, struct{}{})
	sizeBefore := len(lru.m)
	time.Sleep(time.Duration(time.Millisecond * 150))
	sizeAfter := len(lru.m)

	if sizeBefore != 1 || sizeAfter != 1 {
		t.Error("cache no evict not ok")
	}
}
