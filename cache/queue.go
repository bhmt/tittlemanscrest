package cache

import (
	"container/list"
	"time"
)

func now() time.Time {
	return time.Now().UTC()
}

type queue struct {
	l *list.List
}

func newQueue() *queue {
	l := list.New()
	return &queue{l: l}
}

type qElement struct {
	Key any
	T   time.Time
}

func (q *queue) Refresh(e *list.Element, ttl time.Duration) {
	if e == nil || ttl == 0 {
		return
	}

	q.l.MoveToBack(e)

	v := e.Value.(qElement)
	v.T = now()
	v.T = v.T.Add(ttl)

	e.Value = v
}

func (q *queue) Add(key any, ttl time.Duration) *list.Element {
	v := qElement{
		Key: key,
		T:   now(),
	}

	v.T = v.T.Add(ttl)
	return q.l.PushBack(v)
}

func (q *queue) Remove(e *list.Element) {
	if e == nil {
		return
	}

	q.l.Remove(e)
}

func (q *queue) IsStale() bool {
	if e := q.l.Front(); e != nil {
		n := now()
		return !e.Value.(qElement).T.After(n)
	}

	return false
}
