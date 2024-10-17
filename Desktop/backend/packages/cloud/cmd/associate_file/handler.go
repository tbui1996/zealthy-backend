package main

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/circulohealth/sonar-backend/packages/cloud/pkg/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func parseRequest(event events.APIGatewayV2HTTPRequest) (*model.File, error) {
	var file model.File
	err := json.Unmarshal([]byte(event.Body), &file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

type Patient struct {
	MedicaidId               string
	SignedCirculoConsentForm bool
	CirculoConsentFormLink   string
}

type AssociateFileHandler struct {
	DopplerDb *gorm.DB
	SonarDb   *gorm.DB
	Logger    *zap.Logger
	File      *model.File
}

func associateFile(req *AssociateFileHandler) error {
	var wg sync.WaitGroup
	wg.Add(2)                   // nolint
	errs := make(chan error, 2) // nolint
	sonarTx := req.SonarDb.Begin()
	dopplerTx := req.DopplerDb.Begin()

	go func() {
		defer wg.Done()

		result := dopplerTx.Model(&Patient{}).Where("medicaid_id = ?", req.File.MemberID).Updates(map[string]interface{}{
			"signed_circulo_consent_form": true,
			"circulo_consent_form_link":   req.File.FilePath,
			"last_modified_timestamp":     time.Now(),
		})

		if result.Error != nil {
			errs <- result.Error
		} else if result.RowsAffected == 0 {
			errs <- errors.New("0 rows were updated. Member not found with that ID")
		}
	}()

	go func() {
		defer wg.Done()

		result := sonarTx.Model(&model.File{}).Where("file_id = ?", req.File.FileID).Updates(map[string]interface{}{
			"member_id":          req.File.MemberID,
			"date_last_accessed": time.Now(),
		})

		if result.Error != nil {
			errs <- result.Error
		} else if result.RowsAffected == 0 {
			errs <- errors.New("0 rows were updated. File not found with that ID")
		}
	}()

	wg.Wait()
	close(errs)

	dbErrors := make([]string, 0, len(errs)) // no lint
	for err := range errs {
		dbErrors = append(dbErrors, err.Error())
	}

	if len(dbErrors) > 0 {
		sonarTx.Rollback()
		dopplerTx.Rollback()

		errMessage := strings.Join(dbErrors, ", ")
		req.Logger.Error(errMessage)
		return errors.New(errMessage)
	}

	err := commit(sonarTx, dopplerTx, req.Logger)
	if err != nil {
		return err
	}

	return nil
}

func commit(sonarTx *gorm.DB, dopplerTx *gorm.DB, logger *zap.Logger) error {
	var wg sync.WaitGroup
	wg.Add(2)                   // nolint
	errs := make(chan error, 2) // nolint

	go func() {
		defer wg.Done()
		err := sonarTx.Commit().Error

		if err != nil {
			errs <- err
		}
	}()

	go func() {
		defer wg.Done()
		err := dopplerTx.Commit().Error

		if err != nil {
			errs <- err
		}
	}()

	wg.Wait()
	close(errs)

	commitErrors := make([]string, 0, len(errs)) // nolint
	for err := range errs {
		commitErrors = append(commitErrors, err.Error())
	}

	if len(commitErrors) > 0 {
		errMessage := strings.Join(commitErrors, ", ")

		logger.Error(errMessage)
		return errors.New(errMessage)
	}

	return nil
}
