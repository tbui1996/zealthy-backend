package model

type SessionReadReceipt struct {
	SessionUserId int   `json:"sessionUserId"`
	LastRead      int64 `json:"lastRead"`
}
