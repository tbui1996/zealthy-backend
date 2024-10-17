package main

import (
	"github.com/circulohealth/sonar-backend/packages/cloud/pkg/model"
	"gorm.io/gorm"
)

func handler(db *gorm.DB) ([]model.File, error) {
	// TODO: Eventually query the db based on chatid/user id
	// These parameters will hopefully be obtainable via Okta SSO
	var files []model.File
	result := db.Where("deleted_at IS NULL").Find(&files)
	if result.Error != nil {
		return nil, result.Error
	}

	return files, nil
}
