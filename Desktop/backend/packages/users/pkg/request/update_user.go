package request

type UpdateUserRequest struct {
	ID             string `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Group          string `json:"group"`
	OrganizationID int    `json:"organizationId"`
}
