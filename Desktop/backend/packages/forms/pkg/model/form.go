package model

import "time"

type Form struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Created     time.Time  `json:"created"`
	Creator     string     `json:"creator"`
	CreatorId   string     `json:"creatorId"`
	DeletedAt   *time.Time `json:"DeletedAt"`
	DateClosed  *time.Time `json:"dateClosed"`
}
