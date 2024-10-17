package mapper

import "github.com/circulohealth/sonar-backend/packages/users/pkg/mapper/iface"

type externalUserOrganizationRecord struct {
	ID   int
	Name string
}

func (externalUserOrganizationRecord) TableName() string {
	return "users.external_user_organizations"
}

// Necessary to keep package scope of records while exposing an api on insert
func externalUserOrganizationRecordFromInsertInput(input *iface.ExternalUserOrganizationInsertInput) *externalUserOrganizationRecord {
	return &externalUserOrganizationRecord{
		Name: input.Name,
	}
}
