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

type AppointmentRepository struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewAppointmentRepository(logger *zap.Logger) (*AppointmentRepository, error) {
	db, err := dao.OpenConnectionToDoppler()

	if err != nil {
		return nil, err
	}

	return &AppointmentRepository{
		Logger: logger,
		DB:     db,
	}, nil
}

func (repo *AppointmentRepository) FindAll() (*[]model.JoinResult, *appointmenterror.AppointmentRepositoryError) {
	appointments := []model.JoinResult{}
	result := repo.DB.Table("appointments as a").
		Select("a.appointment_id, a.appointment_status, a.appointment_purpose, a.appointment_notes, a.appointment_other_purpose,a.appointment_created, a.appointment_scheduled, a.appointment_status_changed_on, a.circulator_driver_fullname, a.patient_diastolic_blood_pressure, a.patient_systolic_blood_pressure, a.patient_respirations_per_minute, a.patient_pulse_beats_per_minute, a.patient_weight_lbs, a.patient_chief_complaint, b.*, c.agency_provider_id, c.national_provider_id, CONCAT_WS(' ', c.first_name, c.middle_name, c.last_name, ' ', c.suffix) as provider_fullname, c.business_name, c.business_tin, c.business_address_1, c.business_address_2, c.business_city, c.business_state, c.business_zip").
		Joins("JOIN patient as b on a.patient_id = b.patient_id").
		Joins("JOIN agency_provider as c on a.agency_provider_id = c.agency_provider_id").
		Order("a.last_modified_timestamp desc").
		Find(&appointments)
	if result.Error != nil {
		return nil, appointmenterror.New(result.Error.Error(), appointmenterror.UNKNOWN)
	}

	return &appointments, nil
}

func (repo *AppointmentRepository) Save(appointment *model.Appointment) *appointmenterror.AppointmentRepositoryError {
	now := time.Now()
	var result *gorm.DB
	if appointment.IsNew() {
		appointment.LastModifiedTimestamp = &now
		appointment.AppointmentCreated = &now
		appointment.CreatedTimestamp = &now
		appointment.AppointmentStatusChangedOn = &now

		result = repo.DB.Table("appointments").Omit("AppointmentId").Create(appointment)
	} else {
		appointment.LastModifiedTimestamp = &now
		result = repo.DB.Table("appointments").Omit("AppointmentId", "CreatedTimestamp", "AppointmentCreated").Save(&appointment)
	}

	if result.Error == nil {
		return nil
	}

	repo.Logger.Debug(result.Error.Error())
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return appointmenterror.New(fmt.Sprint("Could not find appointment with patient Id: ", appointment.PatientId), appointmenterror.NOT_FOUND)
	}

	return appointmenterror.New(result.Error.Error(), appointmenterror.UNKNOWN)
}

func (repo *AppointmentRepository) Find(appointmentId string) (*model.Appointment, *appointmenterror.AppointmentRepositoryError) {
	appointment := model.Appointment{}

	result := repo.DB.Table("appointments").First(&appointment, "appointment_id = ?", appointmentId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, appointmenterror.New(fmt.Sprintf("Could not find appointment with id: %s", appointmentId), appointmenterror.NOT_FOUND)
		}
		repo.Logger.Error(fmt.Sprintf("Failed to retrieve appointment by id: %s, error: %s", appointment.AppointmentId, result.Error.Error()))
		return nil, appointmenterror.New(result.Error.Error(), appointmenterror.UNKNOWN)
	}
	return &appointment, nil
}

func (repo *AppointmentRepository) Delete(appointment *model.Appointment) *appointmenterror.AppointmentRepositoryError {
	result := repo.DB.Table("appointments").Delete(appointment)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return appointmenterror.New(fmt.Sprintf("Could not find flag with id: %s", appointment.AppointmentId), appointmenterror.NOT_FOUND)
		}
		repo.Logger.Error(fmt.Sprintf("Failed to retrieve appointment by id: %s, error: %s", appointment.AppointmentId, result.Error.Error()))
		return appointmenterror.New(result.Error.Error(), appointmenterror.UNKNOWN)
	}
	return nil
}
