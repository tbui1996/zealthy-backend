package request

type ExternalSignUpRequest struct {
	FullName         string `json:"fullName"`
	OrganizationName string `json:"organizationName"`
}
