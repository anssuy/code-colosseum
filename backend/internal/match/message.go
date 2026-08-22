package match

import "encoding/json"

type MessageType string

const (
	MsgReady            MessageType = "ready"
	MsgSubmit           MessageType = "submit"
	MsgSubmissionResult MessageType = "submission_result"
	MsgOpponentJoined   MessageType = "opponent_joined"
	MsgOpponentLeft     MessageType = "opponent_left"
	MsgMatchStarted     MessageType = "match_started"
	MsgMatchFinished    MessageType = "match_finished"
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

type SubmissionResultPayload struct {
	SubmissionID string `json:"submissionId"`
	Status       string `json:"status"`
	PassedTests  int32  `json:"passedTests"`
	TotalTests   int32  `json:"totalTests"`
}

type OpponentEventPayload struct {
	UserID string `json:"userId"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}
