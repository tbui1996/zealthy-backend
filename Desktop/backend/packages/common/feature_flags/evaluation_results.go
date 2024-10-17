package featureflags

type EvaluationResults struct {
	result map[string]bool
}

func (results *EvaluationResults) FlagOrDefault(flagKey string, defaultValue bool) bool {
	isEnabled, ok := results.result[flagKey]

	if !ok {
		return defaultValue
	}

	return isEnabled
}

func (results *EvaluationResults) Map() map[string]bool {
	return results.result
}
