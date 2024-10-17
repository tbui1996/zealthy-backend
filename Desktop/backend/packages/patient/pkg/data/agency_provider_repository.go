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

type AgencyProviderRepository struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewAgencyProviderRepository(logger *zap.Logger) (*AgencyProviderRepository, error) {
	db, err := dao.OpenConnectionToDoppler()

	if err != nil {
		return nil, err
	}

	return &AgencyProviderRepository{
		Logger: logger,
		DB:     db,
	}, nil
}

func (repo *AgencyProviderRepository) FindAll() (*[]model.AgencyProvider, *appointmenterror.AppointmentRepositoryError) {
	agencyProviders := []model.AgencyProvider{}
	result := repo.DB.Find(&agencyProviders)

	if result.Error != nil {
		return nil, appointmenterror.New(result.Error.Error(), appointmenterror.UNKNOWN)
	}

	return &agencyProviders, nil
}

func (repo *AgencyProviderRepository) Save(agencyProvider *model.AgencyProvider) *appointmenterror.AppointmentRepositoryError {
	now := time.Now()
	var result *gorm.DB
	if agencyProvider.IsNew() {
		agencyProvider.LastModifiedTimestamp = &now
		agencyProvider.CreatedTimestamp = &now

		result = repo.DB.Omit("AgencyProviderId", "CreatedTimestamp", "LastModifiedTimestamp").Create(agencyProvider)
	} else {
		result = repo.DB.Omit("AgencyProviderId", "CreatedTimestamp").Save(&agencyProvider)
	}

	if result.Error == nil {
		return nil
	}

	repo.Logger.Debug(result.Error.Error())
	if dao.IsUniqueConstraintViolation(result.Error, "dodd_number") {
		return appointmenterror.New("DoDD id already in use", appointmenterror.DODD_NUMBER_CONFLICT)
	} else if dao.IsUniqueConstraintViolation(result.Error, "national_provider_id") {
		return appointmenterror.New("National Provider Id already in use", appointmenterror.NATIONAL_PROVIDER_ID_CONFLICT)
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return appointmenterror.New(fmt.Sprintf("Could not find agency provider with Name: %s", agencyProvider.FirstName), appointmenterror.NOT_FOUND)
	}

	return appointmenterror.New(result.Error.Error(), appointmenterror.UNKNOWN)
}

func (repo *AgencyProviderRepository) Find(agencyProviderId string) (*model.AgencyProvider, *appointmenterror.AppointmentRepositoryError) {
	agencyprovider := model.AgencyProvider{}

	result := repo.DB.First(&agencyprovider, "agency_provider_id = ?", agencyProviderId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, appointmenterror.New(fmt.Sprintf("Could not find agency provider with id: %s", agencyProviderId), appointmenterror.NOT_FOUND)
		}
		repo.Logger.Error(fmt.Sprintf("Failed to retrieve patient by id: %s, error: %s", agencyprovider.AgencyProviderId, result.Error.Error()))
		return nil, appointmenterror.New(result.Error.Error(), appointmenterror.UNKNOWN)
	}

	return &agencyprovider, nil
}
