package request

type ChatSessionRequestInternal struct {
	InternalUserID string `json:"internalUserID"`
	UserID         string `json:"userID"`
}

type ChatSessionRequestExternal struct {
	UserID string `json:"userID"`
}
