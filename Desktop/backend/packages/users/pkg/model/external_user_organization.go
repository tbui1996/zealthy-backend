package model

type ExternalUserOrganization struct {
	ID   int
	name StringValueHolder
}

type ExternalUserOrganizationBuilder struct {
	value *ExternalUserOrganization
}

type BuildExternalUserOrganizationInput struct {
	ID int
}

func BuildExternalUserOrganization(input *BuildExternalUserOrganizationInput) *ExternalUserOrganizationBuilder {
	return &ExternalUserOrganizationBuilder{
		value: &ExternalUserOrganization{
			ID: input.ID,
		},
	}
}

func (b *ExternalUserOrganizationBuilder) Value() *ExternalUserOrganization {
	return b.value
}

func (m *ExternalUserOrganization) Name() string {
	return m.name.GetValue()
}

func (m *ExternalUserOrganization) SetName(value string) {
	m.name.SetValue(value)
}

func (m *ExternalUserOrganization) NameChanged() bool {
	return m.name.Changed()
}

func (b *ExternalUserOrganizationBuilder) WithName(value string) *ExternalUserOrganizationBuilder {
	b.value.name.SetInitialValue(value)
	return b
}

func (m *ExternalUserOrganization) IsDeepEqual(other *ExternalUserOrganization) bool {
	return m.ID == other.ID && m.name.IsEqual(&other.name.valueHolder)
}

// Deep Clones properties and references
// If resetChanged is true, it sets all properties on the this struct
// and any nested struct to not changed
func (m *ExternalUserOrganization) Clone(resetChanged bool) *ExternalUserOrganization {
	organization := &ExternalUserOrganization{
		ID:   m.ID,
		name: m.name,
	}

	if resetChanged {
		organization.name.changed = false
	}

	return organization
}
