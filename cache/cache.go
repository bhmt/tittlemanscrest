package cache

import (
	"container/list"
	"sync"
	"time"
)

type cacheElement[V any] struct {
	Value    V
	QElement *list.Element
}

type LRU[K comparable, V any] struct {
	m    map[K]cacheElement[V]
	q    *queue
	size int
	ttl  time.Duration

	mu sync.Mutex
}

func New[K comparable, V any](size int, ttl time.Duration) (*LRU[K, V], error) {
	lru := &LRU[K, V]{
		m:    make(map[K]cacheElement[V]),
		q:    newQueue(),
		size: size,
		ttl:  ttl,
	}

	if ttl != 0 {
		go lru.evict()
	}

	return lru, nil
}

func (lru *LRU[K, V]) Get(k K) (*V, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	i, ok := lru.m[k]
	if !ok {
		return nil, ok
	}

	lru.q.Refresh(i.QElement, lru.ttl)
	return &i.Value, ok
}

func (lru *LRU[K, V]) Add(k K, v V) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	i, ok := lru.m[k]
	if ok {
		lru.q.Remove(i.QElement)
	}

	if !ok && len(lru.m) >= lru.size {
		if e := lru.q.l.Front(); e != nil {
			k := e.Value.(qElement).Key.(K)
			lru.q.Remove(e)
			delete(lru.m, k)
		}
	}

	element := lru.q.Add(k, lru.ttl)
	mapping := cacheElement[V]{
		Value:    v,
		QElement: element,
	}

	lru.m[k] = mapping
}

func (lru *LRU[K, V]) remove(k K) (*V, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	i, ok := lru.m[k]
	if !ok {
		return nil, ok
	}

	value := i.Value
	lru.q.Remove(i.QElement)
	delete(lru.m, k)
	return &value, ok
}

func (lru *LRU[K, V]) evict() {
	for {
		if lru.q.l.Front() == nil {
			time.Sleep(lru.ttl)
			continue
		}

		if !lru.q.IsStale() {
			n := now()
			delay := lru.q.l.Front().Value.(qElement).T.Sub(n)
			time.Sleep(delay)
			continue
		}

		k := lru.q.l.Front().Value.(qElement).Key.(K)
		lru.remove(k)
	}
}
