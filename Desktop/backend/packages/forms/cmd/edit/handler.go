package main

import (
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/request"
	"gorm.io/gorm"
)

type EditFormItemInput struct {
	Form model.Form
	Req  request.EditForm
	Db   *gorm.DB
}

func findFormItem(formid int, db *gorm.DB) (form model.Form, err error) {
	var formToEdit model.Form
	result := db.First(&formToEdit, formid)

	if result.Error != nil {
		return formToEdit, result.Error
	}

	return formToEdit, result.Error
}

func editFormItem(input *EditFormItemInput) error {
	var formToEdit = input.Form
	formToEdit.Title = input.Req.Title
	formToEdit.Description = input.Req.Description

	result := input.Db.Save(&formToEdit)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
