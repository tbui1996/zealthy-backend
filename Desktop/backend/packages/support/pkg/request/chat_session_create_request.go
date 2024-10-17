package request

type ChatSessionCreateRequest struct {
	UserID         string  `json:"userID"`
	Topic          string  `json:"topic"`
	InternalUserID *string `json:"internalUserID"`
	ChatOpen       bool    `json:"chatOpen"`
	Created        int64   `json:"created"`
}
