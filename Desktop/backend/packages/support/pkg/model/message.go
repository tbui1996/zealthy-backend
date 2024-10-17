package model

type Message struct {
	ID               string
	SessionID        string
	SenderID         string
	Message          string
	CreatedTimestamp int64
	FileID           string
}
