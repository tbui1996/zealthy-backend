package response

type FeatureFlagResponse struct {
	Id        int    `json:"id"`
	Key       string `json:"key"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`
	IsEnabled bool   `json:"isEnabled"`
}
