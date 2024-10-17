package data

import (
	"errors"
	"fmt"
	"time"

	"github.com/circulohealth/sonar-backend/packages/common/dao"
	appointmenterror "github.com/circulohealth/sonar-backend/packages/patient/pkg/data/appointment_error"
	"github.com/circulohealth/sonar-backend/packages/patient/pkg/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PatientRepository struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewPatientRepository(logger *zap.Logger) (*PatientRepository, error) {
	db, err := dao.OpenConnectionToDoppler()

	if err != nil {
		return nil, err
	}

	return &PatientRepository{
		Logger: logger,
		DB:     db,
	}, nil
}

func (repo *PatientRepository) FindAll() (*[]model.Patient, *appointmenterror.AppointmentRepositoryError) {
	patients := []model.Patient{}
	result := repo.DB.Find(&patients)

	if result.Error != nil {
		return nil, appointmenterror.New(result.Error.Error(), appointmenterror.UNKNOWN)
	}

	return &patients, nil
}

func (repo *PatientRepository) Save(patient *model.Patient) *appointmenterror.AppointmentRepositoryError {
	now := time.Now()
	var result *gorm.DB
	if patient.IsNew() {
		patient.LastModifiedTimestamp = &now
		patient.CreatedTimestamp = &now

		result = repo.DB.Omit("PatientId", "CreatedTimestamp", "LastModifiedTimestamp").Create(patient)
	} else {
		patient.LastModifiedTimestamp = &now
		result = repo.DB.Omit("PatientId", "CreatedTimestamp").Save(&patient)
	}

	if result.Error == nil {
		return nil
	}

	repo.Logger.Debug(result.Error.Error())
	if dao.IsUniqueConstraintViolation(result.Error, "insurance_id") {
		return appointmenterror.New("insurance id already in use", appointmenterror.INSURANCE_ID_CONFLICT)
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return appointmenterror.New(fmt.Sprintf("Could not find patient with Name: %s", patient.FirstName), appointmenterror.NOT_FOUND)
	}

	return appointmenterror.New(result.Error.Error(), appointmenterror.UNKNOWN)
}

func (repo *PatientRepository) Find(patientId string) (*model.Patient, *appointmenterror.AppointmentRepositoryError) {
	patient := model.Patient{}

	result := repo.DB.First(&patient, "patient_id = ?", patientId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, appointmenterror.New(fmt.Sprintf("Could not find patient with id: %s", patientId), appointmenterror.NOT_FOUND)
		}
		repo.Logger.Error(fmt.Sprintf("Failed to retrieve patient by id: %s, error: %s", patient.PatientId, result.Error.Error()))
		return nil, appointmenterror.New(result.Error.Error(), appointmenterror.UNKNOWN)
	}

	return &patient, nil
}
