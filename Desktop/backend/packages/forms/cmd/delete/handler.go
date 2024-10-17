package main

import (
	"time"

	"github.com/circulohealth/sonar-backend/packages/forms/pkg/model"
	"gorm.io/gorm"
)

type AddDeleteDateInput struct {
	Date   *time.Time
	Db     *gorm.DB
	FormID string
}

func findFormItem(formid string, db *gorm.DB) (form model.Form, err error) {
	var formToDelete model.Form
	result := db.First(&formToDelete, formid)

	if result.Error != nil {
		return formToDelete, result.Error
	}

	return formToDelete, result.Error
}

func addDeleteDate(input *AddDeleteDateInput) error {
	result := input.Db.Model(&model.Form{}).Where("id = ?", input.FormID).Update("deleted_at", &input.Date)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
