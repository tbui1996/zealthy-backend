package main

import (
	"time"

	"github.com/circulohealth/sonar-backend/packages/forms/pkg/model"
	"gorm.io/gorm"
)

func Handler(db *gorm.DB, date time.Time, id string) error {
	result := db.Model(&model.Form{}).Where("id = ?", id).Update("date_closed", &date)
	return result.Error
}
