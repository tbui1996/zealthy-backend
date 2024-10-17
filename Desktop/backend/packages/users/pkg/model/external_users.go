package model

// -- Example --
// externalUser := BuildExternalUser(&BuildExternalUserInput{
//   ID: "1",
//   Username: "test",
//   Email: "test@test.com",
// }).WithStatus("status").WithEnabled(false).Value()
type ExternalUserBuilder struct {
	value ExternalUser
}

type BuildExternalUserInput struct {
	ID       string
	Username string
	Email    string
	Status   string
}

func BuildExternalUser(input *BuildExternalUserInput) *ExternalUserBuilder {
	return &ExternalUserBuilder{
		value: ExternalUser{
			ID:           input.ID,
			Username:     input.Username,
			Email:        input.Email,
			Status:       input.Status,
			enabled:      BoolValueHolder{},
			firstName:    StringValueHolder{},
			lastName:     StringValueHolder{},
			group:        StringValueHolder{},
			hasGroup:     BoolValueHolder{},
			organization: ExternalUserOrganizationValueHolder{},
		},
	}
}

func (b *ExternalUserBuilder) Value() *ExternalUser {
	return &b.value
}

type ExternalUser struct {
	// These cannot be changed
	ID       string
	Username string
	Email    string
	Status   string

	enabled      BoolValueHolder
	firstName    StringValueHolder
	lastName     StringValueHolder
	group        StringValueHolder
	hasGroup     BoolValueHolder
	organization ExternalUserOrganizationValueHolder
}

func (m *ExternalUser) Enabled() bool {
	return m.enabled.GetValue()
}

func (m *ExternalUser) SetEnabled(value bool) {
	m.enabled.SetValue(value)
}

func (m *ExternalUser) EnabledChanged() bool {
	return m.enabled.Changed()
}

func (b *ExternalUserBuilder) WithEnabled(value bool) *ExternalUserBuilder {
	b.value.enabled.SetInitialValue(value)
	return b
}

func (m *ExternalUser) FirstName() string {
	return m.firstName.GetValue()
}

func (m *ExternalUser) SetFirstName(value string) {
	m.firstName.SetValue(value)
}

func (m *ExternalUser) FirstNameChanged() bool {
	return m.firstName.Changed()
}

func (b *ExternalUserBuilder) WithFirstName(value string) *ExternalUserBuilder {
	b.value.firstName.SetInitialValue(value)
	return b
}

func (m *ExternalUser) LastName() string {
	return m.lastName.GetValue()
}

func (m *ExternalUser) SetLastName(value string) {
	m.lastName.SetValue(value)
}

func (m *ExternalUser) LastNameChanged() bool {
	return m.lastName.Changed()
}

func (b *ExternalUserBuilder) WithLastName(value string) *ExternalUserBuilder {
	b.value.lastName.SetInitialValue(value)
	return b
}

func (m *ExternalUser) Group() string {
	return m.group.GetValue()
}

// Calling SetGroup when there is not a group previously
// automatically set's HasGroup to true
func (m *ExternalUser) SetGroup(value string) {
	// if there isn't a group, set hasGroup to true
	if !m.group.HasValue() {
		m.hasGroup.SetValue(true)
	}

	m.group.SetValue(value)
}

func (m *ExternalUser) GroupChanged() bool {
	return m.group.Changed()
}

// if called, automatically set's HasGroup to true as well
func (b *ExternalUserBuilder) WithGroup(value string) *ExternalUserBuilder {
	b.value.group.SetInitialValue(value)
	b.value.hasGroup.SetInitialValue(true)
	return b
}

func (m *ExternalUser) HasGroup() bool {
	return m.hasGroup.GetValue()
}

// Only need to call if removing group from user
func (m *ExternalUser) SetHasGroup(value bool) {
	m.hasGroup.SetValue(value)
}

func (m *ExternalUser) HasGroupChanged() bool {
	return m.hasGroup.Changed()
}

func (m *ExternalUser) Organization() *ExternalUserOrganization {
	return m.organization.GetValue()
}

func (m *ExternalUser) SetOrganization(value *ExternalUserOrganization) {
	m.organization.SetValue(value)
}

func (m *ExternalUser) OrganizationChanged() bool {
	return m.organization.Changed()
}

func (b *ExternalUserBuilder) WithOrganization(value *ExternalUserOrganization) *ExternalUserBuilder {
	b.value.organization.SetInitialValue(value)
	return b
}

// Due to underlying properties of the domain model, a simple comparison is not possible
// this helper function is provided. This is a deep
func (m *ExternalUser) IsDeepEqual(other *ExternalUser) bool {
	if m.ID != other.ID ||
		m.Status != other.Status ||
		m.Email != other.Email ||
		m.Username != other.Username ||
		!m.enabled.IsEqual(&other.enabled.valueHolder) ||
		!m.firstName.IsEqual(&other.firstName.valueHolder) ||
		!m.lastName.IsEqual(&other.lastName.valueHolder) ||
		!m.group.IsEqual(&other.group.valueHolder) ||
		!m.hasGroup.IsEqual(&other.hasGroup.valueHolder) {
		return false
	}

	if m.organization.hasValue != other.organization.hasValue {
		return false
	}

	if m.organization.changed != other.organization.changed {
		return false
	}

	if m.organization.hasValue && !m.Organization().IsDeepEqual(other.Organization()) {
		return false
	}

	return true
}

// Deep Clones properties and references
// If resetChanged is true, it sets all properties on the this struct
// and any nested struct to not changed
func (m *ExternalUser) Clone(resetChanged bool) *ExternalUser {
	user := &ExternalUser{
		ID:       m.ID,
		Username: m.Username,
		Email:    m.Email,
		Status:   m.Status,
		// structs copy all values
		enabled:      m.enabled,
		firstName:    m.firstName,
		lastName:     m.lastName,
		group:        m.group,
		hasGroup:     m.hasGroup,
		organization: m.organization,
	}

	if m.Organization() != nil {
		user.organization.value = m.Organization().Clone(resetChanged)
	}

	if resetChanged {
		user.enabled.changed = false
		user.firstName.changed = false
		user.lastName.changed = false
		user.group.changed = false
		user.hasGroup.changed = false
		user.organization.changed = false
	}

	return user
}
