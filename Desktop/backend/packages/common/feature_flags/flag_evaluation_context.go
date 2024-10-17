package featureflags

import (
	"fmt"

	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"gorm.io/gorm"
)

type FlagEvaluationContext struct {
	Key       string `gorm:"->"`
	IsEnabled bool   `gorm:"->"`
	DeletedAt gorm.DeletedAt
}

// gorm tabler interface
func (ctx *FlagEvaluationContext) TableName() string {
	return fmt.Sprintf("%sflags", dao.FeatureFlags)
}
