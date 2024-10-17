package model

// Not expected to be returned directly to the client, use ChatSessionResponseDTO instead
// Maps values to DynamoDB table for non pending chat sessions
type ChatSessionDTO struct {
	ID                   string   `json:"ID"`
	InternalUserID       *string  `json:"internalUserID"`
	UserID               string   `json:"userID"`
	CreatedTimestamp     int64    `json:"createdTimestamp"`
	ChatOpen             bool     `json:"chatOpen"`
	Topic                string   `json:"topic"`
	Patient              Patient  `json:"patient"`
	UserLastRead         int64    `json:"userLastRead"`
	InternalUserLastRead int64    `json:"internalUserLastRead"`
	Notes                *string  `json:"notes"`
	ChatType             ChatType `json:"chatType"`
	Starred              bool     `json:"starred"`

	LastMessageTimestamp int64  `json:"lastMessageTimestamp"`
	LastMessagePreview   string `json:"lastMessagePreview"`
	LastMessageSenderID  string `json:"lastMessageSenderID"`
}
