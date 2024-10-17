package main

import (
	"time"

	"github.com/circulohealth/sonar-backend/packages/forms/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/request"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type SubmitInput struct {
	Db                     *gorm.DB
	Err                    error
	InputSubmissionRequest request.InputSubmissionRequest
}

type DiscardInput struct {
	Db             *gorm.DB
	Err            error
	DiscardRequest request.DiscardFormRequest
	Deleted        time.Time
}

func validateRequest(s interface{}) error {
	validate := validator.New()
	e := validate.Struct(s)

	return e
}

func submit(input *SubmitInput) error {
	if input.Err != nil {
		return input.Err
	}

	formSubmit := model.FormSubmission{FormSentId: input.InputSubmissionRequest.FormSentId}
	result := input.Db.Create(&formSubmit)

	if result.Error != nil {
		return result.Error
	}

	inputs := make([]model.InputSubmission, len(input.InputSubmissionRequest.SubmitData))
	for i, v := range input.InputSubmissionRequest.SubmitData {
		inputs[i] = model.InputSubmission{FormSubmissionId: formSubmit.ID, InputId: v.ID, Response: v.Response}
	}

	result = input.Db.Create(&inputs)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func discard(input *DiscardInput) error {
	if input.Err != nil {
		return input.Err
	}

	result := input.Db.Create(&model.FormDiscard{FormSentId: input.DiscardRequest.FormSentId, Deleted: input.Deleted})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
