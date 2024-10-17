package model

import "time"

type Session struct {
	ID       int       `json:"id"`
	Created  time.Time `json:"created"`
	ChatType ChatType  `json:"chatType"`
}

type ChatType string

const (
	CIRCULATOR ChatType = "CIRCULATOR"
	GENERAL    ChatType = "GENERAL"
)

func StringToChatType(s string) ChatType {
	switch s {
	case "CIRCULATOR":
		return CIRCULATOR
	default:
		return GENERAL
	}
}

func ChatTypeToString(chatType ChatType) string {
	switch chatType {
	case CIRCULATOR:
		return "CIRCULATOR"
	default:
		return "GENERAL"
	}
}
