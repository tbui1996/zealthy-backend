package featureflags

import "gorm.io/gorm"

type FlagEvaluator struct {
	db *gorm.DB
}

func NewEvaluatorWithDB(db *gorm.DB) *FlagEvaluator {
	return &FlagEvaluator{
		db,
	}
}

func (evaluator *FlagEvaluator) Evaluate() (*EvaluationResults, error) {
	contexts := []FlagEvaluationContext{}
	result := evaluator.db.Find(&contexts)

	if result.Error != nil {
		return nil, result.Error
	}

	evaluatedFlags := make(map[string]bool)

	for _, context := range contexts {
		evaluatedFlags[context.Key] = context.IsEnabled
	}

	return &EvaluationResults{
		result: evaluatedFlags,
	}, nil
}
