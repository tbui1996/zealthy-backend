package repo

import (
	"sort"

	"github.com/circulohealth/sonar-backend/packages/forms/pkg/model"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) Inputs(formId string) ([]model.Input, error) {
	// get inputs
	var inputs []model.Input

	result := repo.db.Find(&inputs, "form_id = ?", formId)

	if result.Error != nil {
		return nil, result.Error
	}

	sort.Slice(inputs, func(i, j int) bool {
		return inputs[i].Order < inputs[j].Order
	})

	return inputs, nil
}

func (repo *Repository) Form(id string) (*model.Form, error) {
	var form model.Form

	result := repo.db.First(&form, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &form, nil
}
