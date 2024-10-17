package request

type PendingChatSessionRequestInternal struct {
	UserID string `json:"userID"`
}

type AssignPendingChatSessionRequestInternal struct {
	SessionID      string `json:"sessionID"`
	InternalUserID string `json:"internalUserID"`
}
