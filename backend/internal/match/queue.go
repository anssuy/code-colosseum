package match

import (
	"sync"
	"time"
)

type QueueEntry struct {
	UserID string
	Rating int32
	Joined time.Time
}

type Queue struct {
	mu      sync.Mutex
	waiting []QueueEntry
}

func NewQueue() *Queue {
	return &Queue{
		waiting: make([]QueueEntry, 0),
	}
}

func (q *Queue) Join(entry QueueEntry) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, e := range q.waiting {
		if e.UserID == entry.UserID {
			return false
		}
	}

	q.waiting = append(q.waiting, entry)
	return true
}

func (q *Queue) Leave(userID string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	for i, e := range q.waiting {
		if e.UserID == userID {
			q.waiting = append(q.waiting[:i], q.waiting[i+1:]...)
			return true
		}
	}

	return false
}

func (q *Queue) Contains(userID string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, e := range q.waiting {
		if e.UserID == userID {
			return true
		}
	}

	return false
}

func (q *Queue) Snapshot() []QueueEntry {
	q.mu.Lock()
	defer q.mu.Unlock()

	out := make([]QueueEntry, len(q.waiting))
	copy(out, q.waiting)
	return out
}
