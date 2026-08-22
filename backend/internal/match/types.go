package match

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Status string

const (
	StatusWaiting   Status = "waiting"
	StatusActive    Status = "active"
	StatusFinished  Status = "finished"
	StatusAbandoned Status = "abandoned"
)

type Match struct {
	ID          pgtype.UUID
	ProblemID   pgtype.UUID
	PlayerOneID pgtype.UUID
	PlayerTwoID pgtype.UUID
	WinnerID    pgtype.UUID
	Status      Status
	CreatedAt   time.Time
}

type MatchFinishedPayload struct {
	WinnerID string `json:"winnerId"`
}

type Player struct {
	UserID pgtype.UUID
	Conn   *Conn
	Ready  bool
}
