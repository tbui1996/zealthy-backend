package model

// Not expected to be returned directly to the client, use ChatSessionResponseDTO instead
// Maps values to DynamoDB table for pending chat sessions
type PendingChatSessionDTO struct {
	ID               string
	UserID           string
	Email            string
	Topic            string
	CreatedTimestamp int64

	LastMessageTimestamp int64  `json:"lastMessageTimestamp"`
	LastMessagePreview   string `json:"lastMessagePreview"`
	LastMessageSenderID  string `json:"lastMessageSenderID"`
}
