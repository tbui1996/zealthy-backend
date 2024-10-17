package mapper

// fields have to be exported for Gorm to place values into them
// however the struct can still be unexported from the package
type externalUserSQLRecord struct {
	ID string

	// null indicates that there is no relationship
	ExternalUserOrganizationID *int
}

// Abstracts database specifics through a struct
type externalUserSQLRecordUpdater struct {
	id string

	externalUserOrganizationID        *int
	externalUserOrganizationIDChanged bool
}

func (externalUserSQLRecord) TableName() string {
	return "users.external_users"
}
