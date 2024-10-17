package iface

import (
	flagerror "github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/data/flag_error"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/model"
)

type FeatureFlagRepository interface {
	Save(flag *model.FeatureFlag) *flagerror.FeatureFlagRepositoryError
	FindAll() (*[]model.FeatureFlag, *flagerror.FeatureFlagRepositoryError)
	Find(id int) (*model.FeatureFlag, *flagerror.FeatureFlagRepositoryError)
	Delete(flag *model.FeatureFlag) *flagerror.FeatureFlagRepositoryError
}
