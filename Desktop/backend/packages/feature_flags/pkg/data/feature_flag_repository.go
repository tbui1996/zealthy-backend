package data

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/circulohealth/sonar-backend/packages/common/dao"
	flagerror "github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/data/flag_error"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FeatureFlagRepository struct {
	UserID string
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewFeatureFlagRepository(userID string, logger *zap.Logger) (*FeatureFlagRepository, error) {
	db, err := dao.OpenConnectionWithTablePrefix(dao.FeatureFlags)
	if err != nil {
		return nil, err
	}
	return &FeatureFlagRepository{
		UserID: userID,
		Logger: logger,
		DB:     db,
	}, nil
}

func (repo *FeatureFlagRepository) Save(flag *model.FeatureFlag) *flagerror.FeatureFlagRepositoryError {
	now := time.Now()
	flag.UpdatedAt = &now
	flag.UpdatedBy = &repo.UserID
	var result *gorm.DB

	if flag.IsNew() {
		flag.CreatedAt = &now
		flag.CreatedBy = &repo.UserID

		result = repo.DB.Omit("CreatedAt", "UpdatedAt").Create(flag)
	} else {
		result = repo.DB.Save(&flag)
	}

	if result.Error == nil {
		return nil
	}

	repo.Logger.Debug(result.Error.Error())

	if dao.IsUniqueConstraintViolation(result.Error, "key_deleted_at_key") {
		return flagerror.New("flagKey already in use", flagerror.KEY_CONFLICT)
	} else if dao.IsUniqueConstraintViolation(result.Error, "name_deleted_at_key") {
		return flagerror.New("name already in use", flagerror.NAME_CONFLICT)
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return flagerror.New(fmt.Sprintf("Could not find flag with id: %s", strconv.Itoa(flag.Id)), flagerror.NOT_FOUND)
	}

	return flagerror.New(result.Error.Error(), flagerror.KEY_CONFLICT)
}

func (repo *FeatureFlagRepository) FindAll() (*[]model.FeatureFlag, *flagerror.FeatureFlagRepositoryError) {
	flags := []model.FeatureFlag{}
	result := repo.DB.Find(&flags)

	if result.Error != nil {
		return nil, flagerror.New(result.Error.Error(), flagerror.UNKNOWN)
	}

	return &flags, nil
}

func (repo *FeatureFlagRepository) Find(id int) (*model.FeatureFlag, *flagerror.FeatureFlagRepositoryError) {
	flag := model.FeatureFlag{}

	result := repo.DB.First(&flag, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, flagerror.New(fmt.Sprintf("Could not find flag with id: %s", strconv.Itoa(id)), flagerror.NOT_FOUND)
		}
		repo.Logger.Error(fmt.Sprintf("Failed to retrieve flag by id: %d, error: %s", flag.Id, result.Error.Error()))
		return nil, flagerror.New(result.Error.Error(), flagerror.UNKNOWN)
	}

	return &flag, nil
}

func (repo *FeatureFlagRepository) Delete(flag *model.FeatureFlag) *flagerror.FeatureFlagRepositoryError {
	result := repo.DB.Delete(flag)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return flagerror.New(fmt.Sprintf("Could not find flag with id: %s", strconv.Itoa(flag.Id)), flagerror.NOT_FOUND)

		}
		repo.Logger.Error(fmt.Sprintf("Failed to retrieve flag by id: %d, error: %s", flag.Id, result.Error.Error()))
		return flagerror.New(result.Error.Error(), flagerror.UNKNOWN)
	}

	return nil
}
