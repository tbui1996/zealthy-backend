package model

type InternalUserMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
