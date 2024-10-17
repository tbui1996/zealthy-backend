package request

type ReadReceiptRequest struct {
	UserID    string `json:"userID"`
	SessionID string `json:"sessionID"`
}
