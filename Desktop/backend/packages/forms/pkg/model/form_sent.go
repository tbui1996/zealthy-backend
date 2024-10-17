package model

import "time"

type FormSent struct {
	ID     int       `json:"id"`
	FormId int       `json:"formId"`
	Sent   time.Time `json:"sent"`
}
