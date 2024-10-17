package main

import (
	"fmt"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type SubmitChatSessionStarRequest struct {
	DB        *gorm.DB
	SessionID string
	OnStar    bool
	Logger    *zap.Logger
}

func Handler(req SubmitChatSessionStarRequest) error {
	onStarStr := strconv.FormatBool(req.OnStar)
	sessionId, err := strconv.Atoi(req.SessionID)

	if err != nil {
		errMsg := fmt.Errorf("expected an int. (%s)", err)
		req.Logger.Error(errMsg.Error())
		return errMsg
	}

	descriptor := model.SessionDescriptor{
		SessionId: sessionId,
		Name:      model.STARRED,
		Value:     onStarStr,
	}

	return req.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "session_id"}, {Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&descriptor).Error
}
