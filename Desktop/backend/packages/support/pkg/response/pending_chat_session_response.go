package response

type PendingChatSession struct {
	ID               string `json:"ID"`
	UserID           string `json:"userID"`
	Email            string `json:"email"`
	Topic            string `json:"topic"`
	CreatedTimestamp int64  `json:"createdTimestamp"`
}
