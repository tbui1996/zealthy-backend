package model

type SessionUser struct {
	ID        int    `json:"id"`
	UserId    string `json:"userId"`
	SessionId int    `json:"sessionId"`
}
