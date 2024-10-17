package request

type TypingActionRequest struct {
	UserID    string `json:"userID"`
	SessionID string `json:"sessionID"`
	Action    string `json:"action"`
}
