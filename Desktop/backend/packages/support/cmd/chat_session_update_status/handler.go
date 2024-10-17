package main

import (
	"fmt"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

type UpdateStatusRequest struct {
	SessionId     string
	Open          bool
	RideScheduled bool
	Logger        *zap.Logger
	DB            *gorm.DB
}

func Handler(req UpdateStatusRequest) error {

	status := model.CLOSED
	if req.Open {
		status = model.OPEN
	}

	id, err := strconv.Atoi(req.SessionId)
	if err != nil {
		errMsg := fmt.Errorf("expected an int when getting sessionId: (%s)", err.Error())
		req.Logger.Error(errMsg.Error())
		return errMsg
	}

	descriptor := model.SessionDescriptor{
		SessionId: id,
		Name:      model.RIDE_SCHEDULED,
		Value:     strconv.FormatBool(req.RideScheduled),
	}

	err = req.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "session_id"}, {Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{"value"}),
		}).Create(&descriptor).Error; err != nil {
			req.Logger.Error(err.Error())
			return err
		}

		t := time.Now()
		update := map[string]interface{}{"status": status}

		if status == model.CLOSED {
			update["closed_at"] = t
		} else {
			update["opened_at"] = t
		}

		if err := tx.Model(&model.SessionStatus{}).Where("session_id = ?", id).Updates(update).Error; err != nil {
			req.Logger.Error(err.Error())
			return err
		}

		return nil
	})

	return err
}
