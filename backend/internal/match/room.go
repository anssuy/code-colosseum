package match

import (
	"context"
	"encoding/json"
	"log"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/judge"
	"github.com/anssuy/code-colosseum/backend/internal/sandbox"

	"github.com/jackc/pgx/v5/pgtype"
)

type Event struct {
	Conn   *Conn
	UserID string
	Data   []byte
}

type submissionResult struct {
	conn   *Conn
	userID string
	result judge.Result
	code   string
	lang   string
}

type Room struct {
	match      dbgen.Match
	players    map[string]*Player
	incoming   chan Event
	unregister chan *Conn
	results    chan submissionResult
	done       chan struct{}
	queries    *dbgen.Queries
	judgePool  *judge.Pool
	manager    *Manager
	testCases  []judge.TestCase
}

func NewRoom(ctx context.Context, m dbgen.Match, queries *dbgen.Queries, pool *judge.Pool) (*Room, error) {
	dbTestCases, err := queries.ListTestCasesForProblem(ctx, m.ProblemID)
	if err != nil {
		return nil, err
	}

	testCases := make([]judge.TestCase, len(dbTestCases))
	for i, tc := range dbTestCases {
		testCases[i] = judge.TestCase{
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
		}
	}

	return &Room{
		match:      m,
		players:    make(map[string]*Player),
		incoming:   make(chan Event, 32),
		unregister: make(chan *Conn),
		results:    make(chan submissionResult, 8),
		done:       make(chan struct{}),
		queries:    queries,
		judgePool:  pool,
		testCases:  testCases,
	}, nil
}

func (r *Room) Run() {
	for {
		select {
		case event := <-r.incoming:
			r.handleEvent(event)
		case res := <-r.results:
			r.handleSubmissionResult(res)
		case conn := <-r.unregister:
			r.handleDisconnect(conn)
		case <-r.done:
			return
		}
	}
}

func (r *Room) Stop() {
	close(r.done)
	if r.manager != nil {
		r.manager.removeRoom(r.match.ID)
	}
}

func (r *Room) AddPlayer(userID string, conn *Conn) {
	r.players[userID] = &Player{
		Conn: conn,
	}

	out, _ := json.Marshal(OutboundMessage{
		Type:    MsgOpponentJoined,
		Payload: OpponentEventPayload{UserID: userID},
	})
	r.broadcastExcept(userID, out)
}

func (r *Room) handleEvent(e Event) {
	var msg InboundMessage
	if err := json.Unmarshal(e.Data, &msg); err != nil {
		r.sendError(e.Conn, "invalid message format")
		return
	}

	switch msg.Type {
	case MsgReady:
		r.handleReady(e.Conn, e.UserID)
	case MsgSubmit:
		var payload SubmitPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			r.sendError(e.Conn, "invalid submit payload")
			return
		}
		r.handleSubmit(e.Conn, e.UserID, payload)
	default:
		r.sendError(e.Conn, "unknown message type")
	}
}

func (r *Room) handleReady(c *Conn, userID string) {
	player, ok := r.players[userID]
	if !ok {
		r.sendError(c, "player not in room")
		return
	}
	player.Ready = true

	if !r.allPlayersReady() {
		return
	}

	updated, err := r.queries.StartMatch(context.Background(), r.match.ID)
	if err != nil {
		log.Printf("StartMatch failed: %v", err)
		r.sendError(c, "could not start match")
		return
	}
	r.match = updated

	out, _ := json.Marshal(OutboundMessage{
		Type: MsgMatchStarted,
	})
	r.broadcast(out)
}

func (r *Room) allPlayersReady() bool {
	if len(r.players) < 2 {
		return false
	}
	for _, p := range r.players {
		if !p.Ready {
			return false
		}
	}
	return true
}

func (r *Room) handleSubmit(c *Conn, userID string, payload SubmitPayload) {
	if r.match.Status != dbgen.MatchStatusActive {
		r.sendError(c, "match is not active")
		return
	}

	if !sandbox.IsSupported(payload.Language) {
		r.sendError(c, "unsupported language")
		return
	}

	resultCh := make(chan judge.Result, 1)
	r.judgePool.Submit(judge.Job{
		Ctx:       context.Background(),
		Language:  payload.Language,
		Code:      payload.SourceCode,
		TestCases: r.testCases,
		ResultCh:  resultCh,
	})

	go func() {
		result := <-resultCh
		r.results <- submissionResult{
			conn:   c,
			userID: userID,
			result: result,
			code:   payload.SourceCode,
			lang:   payload.Language,
		}
	}()
}

func (r *Room) handleSubmissionResult(res submissionResult) {
	var userID pgtype.UUID

	if err := userID.Scan(res.userID); err != nil {
		r.sendError(res.conn, "invalid user id")
		return
	}

	submission, err := r.queries.CreateSubmission(context.Background(), dbgen.CreateSubmissionParams{
		UserID:      userID,
		ProblemID:   r.match.ProblemID,
		MatchID:     r.match.ID,
		Language:    res.lang,
		SourceCode:  res.code,
		Status:      res.result.Status,
		PassedTests: res.result.PassedTests,
		TotalTests:  res.result.TotalTests,
		ExecutionTimeMs: pgtype.Int8{
			Int64: res.result.ExecutionTimeMS,
			Valid: true,
		},
	})
	if err != nil {
		r.sendError(res.conn, "could not save submission")
		return
	}

	out, _ := json.Marshal(OutboundMessage{
		Type: MsgSubmissionResult,
		Payload: SubmissionResultPayload{
			SubmissionID: submission.ID.String(),
			Status:       submission.Status,
			PassedTests:  submission.PassedTests,
			TotalTests:   submission.TotalTests,
		},
	})
	r.broadcast(out)

	if submission.Status == judge.Accepted {
		r.finishMatch(userID)
	}
}

func (r *Room) finishMatch(winnerID pgtype.UUID) {
	updated, err := r.queries.FinishMatch(context.Background(), dbgen.FinishMatchParams{
		ID:       r.match.ID,
		WinnerID: winnerID,
	})
	if err != nil {
		return
	}
	r.match = updated

	out, _ := json.Marshal(OutboundMessage{
		Type: MsgMatchFinished,
		Payload: MatchFinishedPayload{
			WinnerID: winnerID.String(),
		},
	})
	r.broadcast(out)
	r.Stop()
}

func (r *Room) handleDisconnect(c *Conn) {
	for userID, p := range r.players {
		if p.Conn == c {
			delete(r.players, userID)

			out, _ := json.Marshal(OutboundMessage{
				Type:    MsgOpponentLeft,
				Payload: OpponentEventPayload{UserID: userID},
			})
			r.broadcast(out)
			return
		}
	}
}

func (r *Room) sendError(c *Conn, message string) {
	out, _ := json.Marshal(OutboundMessage{
		Type:    MsgError,
		Payload: ErrorPayload{Message: message},
	})
	select {
	case c.send <- out:
	default:
	}
}

func (r *Room) broadcast(msg []byte) {
	for _, p := range r.players {
		if p.Conn != nil {
			select {
			case p.Conn.send <- msg:
			default:
				close(p.Conn.send)
			}
		}
	}
}

func (r *Room) broadcastExcept(excludeUserID string, msg []byte) {
	for userID, p := range r.players {
		if userID == excludeUserID {
			continue
		}
		if p.Conn != nil {
			select {
			case p.Conn.send <- msg:
			default:
				close(p.Conn.send)
			}
		}
	}
}
