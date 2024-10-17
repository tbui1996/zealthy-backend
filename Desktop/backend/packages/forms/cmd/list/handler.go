package main

import (
	"time"

	"gorm.io/gorm"
)

type Result struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Creator     string    `json:"creator"`
	Sent        time.Time `json:"sent"`
}

func getAllForms(db *gorm.DB) ([]Result, error) {
	var forms []Result

	result := db.Raw("SELECT * FROM form.forms forms LEFT JOIN (SELECT ROW_NUMBER() OVER (PARTITION BY form.form_sents.form_id ORDER BY form.form_sents.sent desc) AS rownum,form.form_sents.sent, form.form_sents.form_id FROM form.form_sents) sents ON sents.form_id = forms.id AND sents.rownum = 1 WHERE forms.deleted_at IS NULL").Scan(&forms)

	if result.Error != nil {
		return nil, result.Error
	}

	return forms, nil
}
