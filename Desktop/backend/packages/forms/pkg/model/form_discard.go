package model

import "time"

type FormDiscard struct {
	ID         int       `json:"id"`
	FormSentId int       `json:"formSentId"`
	Deleted    time.Time `json:"deleted"`
}
