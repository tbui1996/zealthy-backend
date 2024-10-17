package main

import (
	"time"

	"gorm.io/gorm"

	"github.com/circulohealth/sonar-backend/packages/common/router"
	"github.com/circulohealth/sonar-backend/packages/forms/pkg/model"
)

type CreateSendRecordInput struct {
	FormID int
	Sent   time.Time
	Db     *gorm.DB
}

func sendThroughRouter(r *router.Session, body string) error {
	return r.Router.Send(&router.RouterSendInput{
		Source:     "forms",
		Action:     "forms",
		Procedure:  "send",
		Body:       body,
		Recipients: []string{},
	})
}

func createSendRecord(Input *CreateSendRecordInput) (model.FormSent, error) {
	sent := model.FormSent{FormId: Input.FormID, Sent: Input.Sent}
	result := Input.Db.Create(&sent)

	return sent, result.Error
}
