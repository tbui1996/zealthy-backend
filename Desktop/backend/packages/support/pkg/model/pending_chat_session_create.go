package model

type PendingChatSessionCreate struct {
	UserID      string    `json:"userID"`
	Email       string    `json:"email"`
	Topic       *string   `json:"topic"`
	Description *ChatType `json:"description"`
	Created     int64     `json:"created"`
	Patient     *Patient  `json:"patient"`
}
