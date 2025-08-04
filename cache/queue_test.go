package cache

import (
	"testing"
	"time"
)

func TestQueueAdd(t *testing.T) {
	q := newQueue()

	want := q.Add("KEY", time.Millisecond)
	got := q.l.Front()

	if want != got {
		t.Errorf("queue Add not ok, want %v got %v", want, got)
	}
}

func TestQueueRemove(t *testing.T) {
	q := newQueue()

	e := q.Add("KEY", time.Millisecond)
	q.Remove(e)

	got := q.l.Front()

	if got != nil {
		t.Errorf("queue Remove not ok, want %v got %v", nil, got)
	}
}

func TestQueueRefresh(t *testing.T) {
	q := newQueue()

	e := q.Add("KEY", time.Millisecond)
	T1 := e.Value.(qElement).T

	time.Sleep(time.Millisecond * 5)

	q.Refresh(e, time.Millisecond)
	T2 := e.Value.(qElement).T

	if T1.After(T2) {
		t.Errorf("queue Refresh time not updated, initial time %v, refreshed time %v", T1, T2)
	}
}

func TestQueueIsStale(t *testing.T) {
	tests := []struct {
		ttl   time.Duration
		sleep time.Duration
		want  bool
	}{
		{
			ttl:   time.Millisecond * 10,
			sleep: time.Millisecond * 2,
			want:  false,
		},
		{
			ttl:   time.Millisecond * 2,
			sleep: time.Millisecond * 5,
			want:  true,
		},
	}

	q := newQueue()
	for _, test := range tests {
		e := q.Add("antigravity", test.ttl)
		time.Sleep(test.sleep)
		got := q.IsStale()

		if test.want != got {
			t.Errorf("queue IsStale not ok for ttl %v, want %v got %v", test.ttl, test.want, got)
		}

		q.Remove(e)
	}
}
