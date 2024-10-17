package main

import (
	"github.com/circulohealth/sonar-backend/packages/cloud/pkg/model"
	"gorm.io/gorm"
	"time"
)

type AddDeleteDateInput struct {
	Date   *time.Time
	Db     *gorm.DB
	FileID int
}

func handler(req *AddDeleteDateInput) error {
	result := req.Db.Model(&model.File{}).Where("id = ?", req.FileID).Update("deleted_at", &req.Date)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
