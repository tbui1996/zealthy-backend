package main

import (
	"errors"
	"fmt"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UpdateChatNotesRequest struct {
	DB        *gorm.DB
	SessionID int
	Notes     string
	Logger    *zap.Logger
}

func Handler(req UpdateChatNotesRequest) error {
	err := req.DB.Transaction(func(tx *gorm.DB) error {
		var sessionNote model.SessionNote
		err := tx.Take(&sessionNote, req.SessionID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				req.Logger.Debug("no chat.session_notes record found, creating")
				if cErr := tx.Create(&model.SessionNote{SessionId: req.SessionID, Notes: req.Notes}).Error; cErr != nil {
					req.Logger.Error(fmt.Sprintf("while creating chat.session_notes record (%s)", cErr))
					return cErr
				}

				return nil
			}

			req.Logger.Error(fmt.Sprintf("while finding chat.session_notes record (%s)", err))
			return err
		}

		req.Logger.Debug("chat.session_notes record found, updating")
		if uErr := tx.Model(&sessionNote).Where("session_id = ?", req.SessionID).Update("notes", req.Notes).Error; uErr != nil {
			req.Logger.Error(fmt.Sprintf("while updating chat.session_notes record (%s)", uErr))
			return uErr
		}

		return nil
	})

	if err != nil {
		req.Logger.Error(fmt.Sprintf("during transaction for chat.session_notes record (%s)", err))
		return err
	}

	return nil
}
