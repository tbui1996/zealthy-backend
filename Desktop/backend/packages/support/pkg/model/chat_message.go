package model

type ChatMessage struct {
	ID               string  `json:"id"`
	SessionID        string  `json:"sessionID"`
	SenderID         string  `json:"senderID"`
	Message          string  `json:"message"`
	CreatedTimestamp int64   `json:"createdTimestamp"`
	FileID           *string `json:"fileID"`
}
