package match

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	dbgen "github.com/anssuy/code-colosseum/backend/internal/db/generated"
	"github.com/anssuy/code-colosseum/backend/internal/elo"
	"github.com/anssuy/code-colosseum/backend/internal/judge"
	"github.com/anssuy/code-colosseum/backend/internal/sandbox"
	"github.com/anssuy/code-colosseum/backend/internal/ws"

	"github.com/jackc/pgx/v5/pgtype"
)

type Room struct {
	mu           sync.Mutex
	match        dbgen.Match
	ready        map[string]bool
	queries      *dbgen.Queries
	judgePool    *judge.Pool
	testCases    []judge.TestCase
	hub          *ws.Hub
	onFinish     func(playerOneID, playerTwoID string)
	disconnected map[string]time.Time
	createdAt    time.Time
}

type submissionResult struct {
	userID     string
	result     judge.Result
	sourceCode string
	language   string
}

func NewRoom(ctx context.Context, m dbgen.Match, queries *dbgen.Queries, pool *judge.Pool, hub *ws.Hub, onFinish func(playerOneID, playerTwoID string)) (*Room, error) {
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
		match:        m,
		ready:        make(map[string]bool),
		queries:      queries,
		judgePool:    pool,
		testCases:    testCases,
		hub:          hub,
		onFinish:     onFinish,
		disconnected: make(map[string]time.Time),
		createdAt:    time.Now(),
	}, nil
}

func (r *Room) HandleReady(userID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.ready[userID] = true

	return r.ready[r.match.PlayerOneID.String()] && r.ready[r.match.PlayerTwoID.String()]
}

func (r *Room) HandleSubmit(userID, language, sourceCode string) error {
	if !sandbox.IsSupported(language) {
		return errors.New("unsupported language")
	}

	resultCh := make(chan judge.Result, 1)
	r.judgePool.Submit(judge.Job{
		Ctx:       context.Background(),
		Language:  language,
		Code:      sourceCode,
		TestCases: r.testCases,
		ResultCh:  resultCh,
	})

	go func() {
		result := <-resultCh
		r.onSubmissionResult(submissionResult{
			userID:     userID,
			result:     result,
			sourceCode: sourceCode,
			language:   language,
		})
	}()

	return nil
}

func (r *Room) Start(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	updated, err := r.queries.StartMatch(ctx, r.match.ID)
	if err != nil {
		return err
	}
	r.match = updated
	return nil
}

func (r *Room) MarkDisconnected(userID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.disconnected[userID] = time.Now()
}

func (r *Room) MarkReconnected(userID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.disconnected, userID)
}

func (r *Room) CheckForfeit(gracePeriod time.Duration) {
	r.mu.Lock()

	p1 := r.match.PlayerOneID.String()
	p2 := r.match.PlayerTwoID.String()

	p1Disc, p1Expired := r.disconnected[p1]
	p2Disc, p2Expired := r.disconnected[p2]

	p1Timedout := p1Expired && time.Since(p1Disc) >= gracePeriod
	p2Timedout := p2Expired && time.Since(p2Disc) >= gracePeriod

	r.mu.Unlock()

	switch {
	case p1Timedout && p2Timedout:
		r.abandon()
	case p1Timedout:
		var winnerID pgtype.UUID
		winnerID.Scan(p2)
		r.finish(winnerID)
	case p2Timedout:
		var winnerID pgtype.UUID
		winnerID.Scan(p1)
		r.finish(winnerID)
	}
}

func (r *Room) CheckReadyTimeout(timeout time.Duration) {
	r.mu.Lock()

	if r.match.Status != dbgen.MatchStatusWaiting {
		r.mu.Unlock()
		return
	}

	if time.Since(r.createdAt) < timeout {
		r.mu.Unlock()
		return
	}

	p1 := r.match.PlayerOneID.String()
	p2 := r.match.PlayerTwoID.String()
	p1Ready := r.ready[p1]
	p2Ready := r.ready[p2]

	r.mu.Unlock()

	switch {
	case p1Ready && !p2Ready:
		var winnerID pgtype.UUID
		winnerID.Scan(p1)
		r.finish(winnerID)
	case p2Ready && !p1Ready:
		var winnerID pgtype.UUID
		winnerID.Scan(p2)
		r.finish(winnerID)
	case !p1Ready && !p2Ready:
		r.abandon()
	}
}

func (r *Room) onSubmissionResult(res submissionResult) {
	r.mu.Lock()

	var userID pgtype.UUID
	if err := userID.Scan(res.userID); err != nil {
		r.mu.Unlock()
		return
	}

	submission, err := r.queries.CreateSubmission(context.Background(), dbgen.CreateSubmissionParams{
		UserID:      userID,
		ProblemID:   r.match.ProblemID,
		MatchID:     r.match.ID,
		Language:    res.language,
		SourceCode:  res.sourceCode,
		Status:      res.result.Status,
		PassedTests: res.result.PassedTests,
		TotalTests:  res.result.TotalTests,
		ExecutionTimeMs: pgtype.Int8{
			Int64: res.result.ExecutionTimeMS,
			Valid: true,
		},
	})
	if err != nil {
		r.mu.Unlock()
		return
	}

	accepted := submission.Status == judge.Accepted
	r.mu.Unlock()

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

	if accepted {
		r.finish(userID)
	}
}

func (r *Room) broadcast(msg []byte) {
	r.hub.SendTo(r.match.PlayerOneID.String(), msg)
	r.hub.SendTo(r.match.PlayerTwoID.String(), msg)
}

func (r *Room) abandon() {
	r.mu.Lock()
	updated, err := r.queries.AbandonMatch(context.Background(), r.match.ID)
	if err != nil {
		r.mu.Unlock()
		return
	}
	r.match = updated
	r.mu.Unlock()

	out, _ := json.Marshal(OutboundMessage{Type: MsgMatchAbandoned})
	r.broadcast(out)

	if r.onFinish != nil {
		r.onFinish(r.match.PlayerOneID.String(), r.match.PlayerTwoID.String())
	}
}

func (r *Room) finish(winnerID pgtype.UUID) {
	r.mu.Lock()

	updated, err := r.queries.FinishMatch(context.Background(), dbgen.FinishMatchParams{
		ID:       r.match.ID,
		WinnerID: winnerID,
	})
	if err != nil {
		r.mu.Unlock()
		return
	}
	r.match = updated

	loserID := r.match.PlayerOneID
	if winnerID == r.match.PlayerOneID {
		loserID = r.match.PlayerTwoID
	}

	r.mu.Unlock()

	r.applyRatingChanges(winnerID, loserID)

	out, _ := json.Marshal(OutboundMessage{
		Type:    MsgMatchFinished,
		Payload: MatchFinishedPayload{WinnerID: winnerID.String()},
	})
	r.broadcast(out)

	if r.onFinish != nil {
		r.onFinish(r.match.PlayerOneID.String(), r.match.PlayerTwoID.String())
	}
}

func (r *Room) applyRatingChanges(winnerID, loserID pgtype.UUID) {
	ctx := context.Background()

	winner, err := r.queries.GetUserByID(ctx, winnerID)
	if err != nil {
		return
	}
	loser, err := r.queries.GetUserByID(ctx, loserID)
	if err != nil {
		return
	}

	newWinnerRating, newLoserRating := elo.Update(winner.Rating, loser.Rating)

	r.queries.RecordMatchResult(ctx, dbgen.RecordMatchResultParams{
		ID: winnerID, Rating: newWinnerRating, Wins: 1, Losses: 0,
	})
	r.queries.RecordMatchResult(ctx, dbgen.RecordMatchResultParams{
		ID: loserID, Rating: newLoserRating, Wins: 0, Losses: 1,
	})
}
