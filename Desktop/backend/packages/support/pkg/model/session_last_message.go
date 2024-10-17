package model

type SessionLastMessage struct {
	SessionUserId int    `json:"sessionUserId"`
	LastMessage   string `json:"lastMessage"`
	LastSent      int64  `json:"lastSent"`
}
