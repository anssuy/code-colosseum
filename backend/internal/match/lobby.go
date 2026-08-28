package match

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/judge"
	"github.com/anssuy/code-colosseum/backend/internal/ws"

	"github.com/jackc/pgx/v5/pgtype"
)

type Lobby struct {
	mu      sync.RWMutex
	byUser  map[string]*Room
	hub     *ws.Hub
	queries *dbgen.Queries
	pool    *judge.Pool
	queue   *Queue
}

func NewLobby(hub *ws.Hub, queries *dbgen.Queries, pool *judge.Pool) *Lobby {
	l := &Lobby{
		byUser:  make(map[string]*Room),
		hub:     hub,
		queries: queries,
		pool:    pool,
		queue:   NewQueue(),
	}

	hub.OnMessage(l.HandleMessage)
	hub.OnConnect(l.HandleConnect)
	hub.OnDisconnect(l.HandleDisconnect)
	l.startMatchmaker()
	l.startForfeitSweep()

	return l
}

func (l *Lobby) HandleMessage(userID string, data []byte) {
	var msg InboundMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		l.sendError(userID, "invalid message format")
		return
	}

	switch msg.Type {
	case MsgQueueJoin:
		l.handleQueueJoin(userID)
	case MsgQueueLeft:
		l.handleQueueLeft(userID)
	case MsgReady, MsgSubmit:
		l.routeToRoom(userID, msg)
	default:
		l.sendError(userID, "unknown message type")
	}
}

func (l *Lobby) HandleDisconnect(userID string) {
	l.mu.RLock()
	room, ok := l.byUser[userID]
	l.mu.RUnlock()

	if ok {
		room.MarkDisconnected(userID)
		return
	}

	l.queue.Leave(userID)
}

func (l *Lobby) HandleConnect(userID string) {
	l.mu.RLock()
	room, ok := l.byUser[userID]
	l.mu.RUnlock()

	if ok {
		room.MarkReconnected(userID)
	}
}

func (l *Lobby) sendError(userID string, message string) {
	out, _ := json.Marshal(OutboundMessage{
		Type:    MsgError,
		Payload: ErrorPayload{Message: message},
	})
	l.hub.SendTo(userID, out)
}

func (l *Lobby) handleQueueJoin(userID string) {
	l.mu.RLock()
	_, inMatch := l.byUser[userID]
	l.mu.RUnlock()

	if inMatch {
		l.sendError(userID, "already in a match")
		return
	}

	var uid pgtype.UUID
	if err := uid.Scan(userID); err != nil {
		l.sendError(userID, "invalid user")
		return
	}

	user, err := l.queries.GetUserByID(context.Background(), uid)
	if err != nil {
		l.sendError(userID, "could not join queue")
		return
	}

	joined := l.queue.Join(QueueEntry{
		UserID: userID,
		Rating: user.Rating,
		Joined: time.Now(),
	})

	if !joined {
		l.sendError(userID, "already in a queue")
		return
	}

	out, _ := json.Marshal(OutboundMessage{Type: MsgQueueJoined})
	l.hub.SendTo(userID, out)
}

func (l *Lobby) handleQueueLeft(userID string) {
	left := l.queue.Leave(userID)
	if !left {
		l.sendError(userID, "not in queue")
		return
	}

	out, _ := json.Marshal(OutboundMessage{Type: MsgQueueLeft})
	l.hub.SendTo(userID, out)
}

func (l *Lobby) pair(a, b QueueEntry) {
	l.queue.Leave(a.UserID)
	l.queue.Leave(b.UserID)

	avgRating := (a.Rating + b.Rating) / 2
	difficulty := difficultyForRating(avgRating)

	ctx := context.Background()

	problem, err := l.queries.GetRandomProblemByDifficulty(ctx, difficulty)
	if err != nil {
		l.sendError(a.UserID, "could not find a match")
		l.sendError(b.UserID, "could not find a match")
		return
	}

	var playerOne, playerTwo pgtype.UUID
	playerOne.Scan(a.UserID)
	playerTwo.Scan(b.UserID)

	dbMatch, err := l.queries.CreateMatch(ctx, dbgen.CreateMatchParams{
		ProblemID:   problem.ID,
		PlayerOneID: playerOne,
		PlayerTwoID: playerTwo,
	})
	if err != nil {
		l.sendError(a.UserID, "could not create match")
		l.sendError(b.UserID, "could not create match")
		return
	}

	room, err := NewRoom(ctx, dbMatch, l.queries, l.pool, l.hub, func(p1, p2 string) {
		l.mu.Lock()
		delete(l.byUser, p1)
		delete(l.byUser, p2)
		l.mu.Unlock()
	})
	if err != nil {
		l.sendError(a.UserID, "could not start match")
		l.sendError(b.UserID, "could not start match")
		return
	}

	l.mu.Lock()
	l.byUser[a.UserID] = room
	l.byUser[b.UserID] = room
	l.mu.Unlock()

	out, _ := json.Marshal(OutboundMessage{
		Type:    MsgMatchFound,
		Payload: MatchFoundPayload{MatchID: dbMatch.ID.String()},
	})
	l.hub.SendTo(a.UserID, out)
	l.hub.SendTo(b.UserID, out)
}

func difficultyForRating(avgRating int32) dbgen.Difficulty {
	switch {
	case avgRating < 900:
		return dbgen.DifficultyEasy
	case avgRating < 1300:
		return dbgen.DifficultyMedium
	default:
		return dbgen.DifficultyHard
	}
}

func (l *Lobby) routeToRoom(userID string, msg InboundMessage) {
	l.mu.RLock()
	room, ok := l.byUser[userID]
	l.mu.RUnlock()

	if !ok {
		l.sendError(userID, "not in a match")
		return
	}

	switch msg.Type {
	case MsgReady:
		allReady := room.HandleReady(userID)
		if allReady {
			if err := room.Start(context.Background()); err == nil {
				out, _ := json.Marshal(OutboundMessage{Type: MsgMatchStarted})
				room.broadcast(out)
			}
		}
	case MsgSubmit:
		var payload SubmitPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			l.sendError(userID, "invalid submit payload")
			return
		}
		if err := room.HandleSubmit(userID, payload.Language, payload.SourceCode); err != nil {
			l.sendError(userID, err.Error())
		}
	}
}
