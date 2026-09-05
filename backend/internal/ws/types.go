package ws

type InboundEvent struct {
	UserID string
	Data   []byte
}

type Handler func(userID string, data []byte)
