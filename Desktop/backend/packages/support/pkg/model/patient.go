package model

import "time"

type Patient struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	LastName    string    `json:"lastName"`
	Address     string    `json:"address"`
	InsuranceID string    `json:"insuranceID"`
	Birthday    time.Time `json:"birthday"`
	ProviderId  string    `json:"providerId"`
}
