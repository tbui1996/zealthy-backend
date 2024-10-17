package main

import (
	"time"

	"github.com/circulohealth/sonar-backend/packages/forms/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/request"

	"gorm.io/gorm"
)

type CreateFormItemInput struct {
	In      request.CreateForm
	Created time.Time
	Db      *gorm.DB
}

type CreateInputItemsInput struct {
	FormId int
	Form   request.CreateForm
	Db     *gorm.DB
}

func createFormItem(Input *CreateFormItemInput) (form model.Form, err error) {
	// translate
	form = model.Form{Title: Input.In.Title, Description: Input.In.Description, Created: Input.Created, Creator: Input.In.Creator, CreatorId: Input.In.CreatorId, DeletedAt: nil, DateClosed: nil}

	result := Input.Db.Create(&form)

	return form, result.Error
}

func createInputItems(Input *CreateInputItemsInput) error {
	inputs := make([]model.Input, len(Input.Form.Inputs))
	for i, v := range Input.Form.Inputs {
		inputs[i] = model.Input{FormId: Input.FormId, Order: v.Order, Type: v.Type, Options: v.Options, Label: v.Label}
	}

	result := Input.Db.Create(&inputs)

	return result.Error
}
