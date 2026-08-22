package match

import (
	"context"
	"sync"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/judge"

	"github.com/jackc/pgx/v5/pgtype"
)

type Manager struct {
	mu      sync.RWMutex
	rooms   map[string]*Room
	queries *dbgen.Queries
	pool    *judge.Pool
}

func NewManager(queries *dbgen.Queries, pool *judge.Pool) *Manager {
	return &Manager{
		rooms:   make(map[string]*Room),
		queries: queries,
		pool:    pool,
	}
}

func (m *Manager) CreateRoom(ctx context.Context, dbMatch dbgen.Match) (*Room, error) {
	key := dbMatch.ID.String()

	m.mu.Lock()
	defer m.mu.Unlock()

	if existing, ok := m.rooms[key]; ok {
		return existing, nil
	}

	room, err := NewRoom(ctx, dbMatch, m.queries, m.pool)
	if err != nil {
		return nil, err
	}
	room.manager = m

	m.rooms[key] = room
	go room.Run()

	return room, nil
}

func (m *Manager) GetRoom(matchID pgtype.UUID) (*Room, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	room, ok := m.rooms[matchID.String()]
	return room, ok
}

func (m *Manager) removeRoom(matchID pgtype.UUID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.rooms, matchID.String())
}
