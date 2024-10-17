package model

type SessionStatus struct {
	SessionId int        `json:"sessionId"`
	Status    ChatStatus `json:"status"`
}

type ChatStatus string

const (
	OPEN    ChatStatus = "OPEN"
	CLOSED  ChatStatus = "CLOSED"
	PENDING ChatStatus = "PENDING"
)
