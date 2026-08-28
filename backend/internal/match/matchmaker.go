package match

import (
	"sort"
	"time"
)

const (
	baseThreshold   = 50
	thresholdGrowth = 25
)

func (l *Lobby) startMatchmaker() {
	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for range ticker.C {
			l.tryPairPlayers()
		}
	}()
}

func (l *Lobby) startForfeitSweep() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			l.sweepForfeits()
		}
	}()
}

func (l *Lobby) sweepForfeits() {
	l.mu.RLock()
	seen := make(map[*Room]bool)
	rooms := make([]*Room, 0, len(l.byUser)/2)
	for _, room := range l.byUser {
		if !seen[room] {
			seen[room] = true
			rooms = append(rooms, room)
		}
	}
	l.mu.RUnlock()

	for _, room := range rooms {
		room.CheckForfeit(60 * time.Second)
		room.CheckReadyTimeout(30 * time.Second)
	}
}

func (l *Lobby) tryPairPlayers() {
	entries := l.queue.Snapshot()
	if len(entries) < 2 {
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Rating < entries[j].Rating
	})

	paired := make(map[string]bool)

	for i := 0; i < len(entries); i++ {
		a := entries[i]
		if paired[a.UserID] {
			continue
		}

		for j := i + 1; j < len(entries); j++ {
			b := entries[j]
			if paired[b.UserID] {
				continue
			}

			if compatible(a, b) {
				paired[a.UserID] = true
				paired[b.UserID] = true
				l.pair(a, b)
				break
			}
		}
	}
}

func compatible(a, b QueueEntry) bool {
	waited := time.Since(a.Joined)
	if time.Since(b.Joined) > waited {
		waited = time.Since(b.Joined)
	}

	threshold := int32(baseThreshold + int(waited.Seconds())*thresholdGrowth)

	diff := a.Rating - b.Rating
	if diff < 0 {
		diff = -diff
	}

	return diff <= threshold
}
