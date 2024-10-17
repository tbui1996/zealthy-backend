package model

import (
	"fmt"
	"time"

	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"gorm.io/gorm"
)

type FeatureFlag struct {
	Id        int
	Key       string `gorm:"unique;<-:create"`
	Name      string
	IsEnabled bool
	CreatedAt *time.Time `gorm:"<-:create"`
	CreatedBy *string    `gorm:"<-:create"`
	UpdatedAt *time.Time
	UpdatedBy *string
	DeletedAt gorm.DeletedAt
}

func NewFeatureFlagWithUserId(featureFlag FeatureFlag, userId *string) *FeatureFlag {
	return &FeatureFlag{
		Key:       featureFlag.Key,
		Name:      featureFlag.Name,
		IsEnabled: false,
		CreatedBy: userId,
		UpdatedBy: userId,
	}
}

func (flag *FeatureFlag) IsNew() bool {
	return flag.CreatedAt == nil
}

// gorm tabler interface
func (flag *FeatureFlag) TableName() string {
	return fmt.Sprintf("%sflags", dao.FeatureFlags)
}
