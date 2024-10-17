package model

import "github.com/lib/pq"

type Input struct {
	ID      int            `json:"id"`
	Order   int            `json:"order" validate:"required,gte=0"`
	Type    string         `json:"type" validate:"required,oneof=number password radio select telephone text link email divider checkbox message"`
	FormId  int            `json:"formId"`
	Label   string         `json:"label" validate:"required"`
	Options pq.StringArray `json:"options" gorm:"type:text[]"`
}
