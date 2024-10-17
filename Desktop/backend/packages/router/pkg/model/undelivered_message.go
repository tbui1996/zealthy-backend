package model

type UndeliveredMessage struct {
	UserID           string
	CreatedTimestamp int64
	DeleteTimestamp  int64
	Message          string
}
