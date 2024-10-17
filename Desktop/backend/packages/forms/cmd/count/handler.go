package main

import (
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/model"
	"gorm.io/gorm"
)

func countFormsSent(db *gorm.DB) (int64, error) {
	var count int64
	result := db.Model(&model.FormSent{}).Count(&count)

	return count, result.Error
}
