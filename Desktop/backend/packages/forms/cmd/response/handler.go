package main

import (
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/model"
	"gorm.io/gorm"
)

type DiscardAndSubmitInput struct {
	Db        *gorm.DB
	FormSents []int
	Submit    []model.FormSubmission
}

func getDiscardAndSubmitValues(in *DiscardAndSubmitInput) ([]model.FormDiscard, [][]model.InputSubmission, error) {
	tx := in.Db.Begin()
	var discard []model.FormDiscard
	result := tx.Where("form_sent_id IN ?", in.FormSents).Find(&discard)

	if result.Error != nil {
		tx.Rollback()
		return nil, nil, result.Error
	}

	var inputs [][]model.InputSubmission
	for _, v := range in.Submit {
		var input []model.InputSubmission
		result = tx.Find(&input, "form_submission_id = ?", v.ID)

		if result.Error != nil {
			tx.Rollback()
			return nil, nil, result.Error
		}

		if input != nil {
			inputs = append(inputs, input)
		}
	}

	tx.Commit()

	return discard, inputs, nil
}

func findFormSent(id string, db *gorm.DB) ([]model.FormSent, error) {
	var sents []model.FormSent
	result := db.Find(&sents, "form_id = ?", id)

	return sents, result.Error
}

func findSubmitByFormSent(formSents []int, db *gorm.DB) ([]model.FormSubmission, error) {
	var submit []model.FormSubmission
	result := db.Where("form_sent_id IN ?", formSents).Find(&submit)

	return submit, result.Error
}
