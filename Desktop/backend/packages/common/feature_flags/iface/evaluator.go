package iface

import featureflags "github.com/circulohealth/sonar-backend/packages/common/feature_flags"

type FlagEvaluator interface {
	Evaluate() (*featureflags.EvaluationResults, error)
}
