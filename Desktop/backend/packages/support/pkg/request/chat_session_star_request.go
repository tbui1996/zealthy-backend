package request

type ChatSessionStarRequest struct {
	SessionID string `json:"sessionID"`
	OnStar    bool   `json:"onStar"`
}
