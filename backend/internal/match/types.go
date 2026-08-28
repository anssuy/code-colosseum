package match

import "encoding/json"

type MessageType string

const (
	MsgQueueJoin        MessageType = "queue_join"
	MsgQueueLeft        MessageType = "queue_left"
	MsgQueueJoined      MessageType = "queue_joined"
	MsgMatchFound       MessageType = "match_found"
	MsgReady            MessageType = "ready"
	MsgSubmit           MessageType = "submit"
	MsgSubmissionResult MessageType = "submission_result"
	MsgOpponentJoined   MessageType = "opponent_joined"
	MsgOpponentLeft     MessageType = "opponent_left"
	MsgMatchStarted     MessageType = "match_started"
	MsgMatchFinished    MessageType = "match_finished"
	MsgMatchAbandoned   MessageType = "match_abandoned"
	MsgError            MessageType = "error"
)

type InboundMessage struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type OutboundMessage struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

type SubmitPayload struct {
	Language   string `json:"language"`
	SourceCode string `json:"sourceCode"`
}

type MatchFoundPayload struct {
	MatchID string `json:"matchId"`
}

type SubmissionResultPayload struct {
	SubmissionID string `json:"submissionId"`
	Status       string `json:"status"`
	PassedTests  int32  `json:"passedTests"`
	TotalTests   int32  `json:"totalTests"`
}

type MatchFinishedPayload struct {
	WinnerID string `json:"winnerId"`
}

type OpponentEventPayload struct {
	UserID string `json:"userId"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}
