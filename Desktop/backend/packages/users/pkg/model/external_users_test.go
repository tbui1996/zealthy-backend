package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildExternalUser(t *testing.T) {
	baseInput := &BuildExternalUserInput{
		ID:       "1",
		Username: "test",
		Email:    "test@gmail.com",
		Status:   "test",
	}

	organization := &ExternalUserOrganization{}

	tests := []struct {
		builder  *ExternalUserBuilder
		expected *ExternalUser
	}{
		{
			// default values
			builder: BuildExternalUser(baseInput),
			expected: &ExternalUser{
				ID:       "1",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
			},
		},
		{
			builder: BuildExternalUser(baseInput).WithEnabled(true),
			expected: &ExternalUser{
				ID:       "1",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    true,
						hasValue: true,
						changed:  false,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
			},
		},
		{
			builder: BuildExternalUser(baseInput).WithEnabled(true).WithFirstName("Charles").WithLastName("Scholle").WithGroup("Test").WithOrganization(organization),
			expected: &ExternalUser{
				ID:       "1",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    true,
						hasValue: true,
						changed:  false,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    "Charles",
						hasValue: true,
						changed:  false,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    "Scholle",
						hasValue: true,
						changed:  false,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    "Test",
						hasValue: true,
						changed:  false,
					},
				},
				// Assert that setting the group set's to true
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    true,
						hasValue: true,
						changed:  false,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value:    organization,
						hasValue: true,
						changed:  false,
					},
				},
			},
		},
	}

	for _, test := range tests {
		actual := test.builder.Value()
		assert.Equal(t, test.expected, actual)
	}
}

func TestModifyExternalUser(t *testing.T) {
	organizationDefault := &ExternalUserOrganization{}
	organizationTester := &ExternalUserOrganization{}

	baseInput := &BuildExternalUserInput{
		ID:       "1",
		Username: "test",
		Email:    "test@gmail.com",
		Status:   "test",
	}

	buildExternalUser := func() *ExternalUser {
		return BuildExternalUser(baseInput).WithEnabled(true).WithFirstName("FN").WithLastName("LN").WithGroup("Test").WithOrganization(organizationDefault).Value()
	}

	tests := []struct {
		executor func() *ExternalUser
		expected *ExternalUser
	}{
		{
			// test changing values sets changed to true without previously set values
			executor: func() *ExternalUser {
				base := BuildExternalUser(baseInput).Value()
				base.SetEnabled(false)
				base.SetFirstName("FN2")
				base.SetLastName("LN2")
				base.SetGroup("Test2")
				// base.SetHasGroup(false)
				base.SetOrganization(organizationTester)
				return base
			},
			expected: &ExternalUser{
				ID:       "1",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    false,
						hasValue: true,
						changed:  true,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    "FN2",
						hasValue: true,
						changed:  true,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    "LN2",
						hasValue: true,
						changed:  true,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    "Test2",
						hasValue: true,
						changed:  true,
					},
				},
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    true,
						hasValue: true,
						// setting the group set's has group to true
						changed: true,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value:    organizationTester,
						hasValue: true,
						changed:  true,
					},
				},
			},
		},
		{
			// test changing values sets changed to true when already previously set
			executor: func() *ExternalUser {
				base := buildExternalUser()
				base.SetEnabled(false)
				base.SetFirstName("FN2")
				base.SetLastName("LN2")
				base.SetGroup("Test2")
				// base.SetHasGroup(false)
				base.SetOrganization(organizationTester)
				return base
			},
			expected: &ExternalUser{
				ID:       "1",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    false,
						hasValue: true,
						changed:  true,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    "FN2",
						hasValue: true,
						changed:  true,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    "LN2",
						hasValue: true,
						changed:  true,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    "Test2",
						hasValue: true,
						changed:  true,
					},
				},
				// Assert that setting the group set's to true
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    true,
						hasValue: true,
						changed:  false,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value:    organizationTester,
						hasValue: true,
						changed:  true,
					},
				},
			},
		},
		{
			// test setting has group to false
			executor: func() *ExternalUser {
				base := buildExternalUser()
				base.SetHasGroup(false)
				return base
			},
			expected: &ExternalUser{
				ID:       "1",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    true,
						hasValue: true,
						changed:  false,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    "FN",
						hasValue: true,
						changed:  false,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    "LN",
						hasValue: true,
						changed:  false,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    "Test",
						hasValue: true,
						changed:  false,
					},
				},
				// Assert that setting the group set's to true
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    false,
						hasValue: true,
						changed:  true,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value:    organizationDefault,
						hasValue: true,
						changed:  false,
					},
				},
			},
		},
	}

	for _, test := range tests {
		actual := test.executor()
		assert.Equal(t, test.expected, actual)
	}
}

func TestIsDeepEqual(t *testing.T) {
	tests := []struct {
		original *ExternalUser
		other    *ExternalUser
		expected bool
	}{
		{
			// defaults return true
			original: &ExternalUser{
				ID:       "1",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
			},
			other: &ExternalUser{
				ID:       "1",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
			},
			expected: true,
		},
		{
			// basic property mismatch fails
			original: &ExternalUser{
				// this is the mismatched property
				ID:       "2",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
			},
			other: &ExternalUser{
				ID:       "1",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
			},
			expected: false,
		},
		{
			// value holder mismatch fails
			original: &ExternalUser{
				ID:       "1",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    "First Name",
						hasValue: true,
						changed:  false,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
			},
			other: &ExternalUser{
				ID:       "1",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
			},
			expected: false,
		},
		{
			// organization matches
			original: &ExternalUser{
				ID:       "1",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value: &ExternalUserOrganization{
							ID: 0,
							name: StringValueHolder{
								valueHolder: valueHolder{
									value:    "test",
									hasValue: true,
									changed:  false,
								},
							},
						},
						hasValue: true,
						changed:  false,
					},
				},
			},
			other: &ExternalUser{
				ID:       "1",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value: &ExternalUserOrganization{
							ID: 0,
							name: StringValueHolder{
								valueHolder: valueHolder{
									value:    "test",
									hasValue: true,
									changed:  false,
								},
							},
						},
						hasValue: true,
						changed:  false,
					},
				},
			},
			expected: true,
		},
		{
			// organization does not match
			original: &ExternalUser{
				ID:       "1",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value: &ExternalUserOrganization{
							ID: 0,
							name: StringValueHolder{
								valueHolder: valueHolder{
									value:    "test",
									hasValue: true,
									changed:  false,
								},
							},
						},
						hasValue: true,
						changed:  false,
					},
				},
			},
			other: &ExternalUser{
				ID:       "1",
				Username: "test",
				Email:    "test@gmail.com",
				Status:   "test",
				enabled: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				firstName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				lastName: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				group: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				hasGroup: BoolValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
				organization: ExternalUserOrganizationValueHolder{
					valueHolder: valueHolder{
						value: &ExternalUserOrganization{
							ID: 1,
							name: StringValueHolder{
								valueHolder: valueHolder{
									value:    "fail",
									hasValue: true,
									changed:  false,
								},
							},
						},
						hasValue: true,
						changed:  false,
					},
				},
			},
			expected: false,
		},
	}

	for _, test := range tests {
		actual := test.original.IsDeepEqual(test.other)
		assert.Equal(t, test.expected, actual)
	}
}

func TestClone_PropertiesAreImmutable(t *testing.T) {
	user := &ExternalUser{
		ID:       "1",
		Username: "test",
		Email:    "test@gmail.com",
		Status:   "test",
		enabled: BoolValueHolder{
			valueHolder: valueHolder{
				value:    true,
				hasValue: true,
				changed:  false,
			},
		},
		firstName: StringValueHolder{
			valueHolder: valueHolder{
				value:    "firstname",
				hasValue: true,
				changed:  false,
			},
		},
		lastName: StringValueHolder{
			valueHolder: valueHolder{
				value:    "lastname",
				hasValue: true,
				changed:  true,
			},
		},
		group: StringValueHolder{
			valueHolder: valueHolder{
				value:    nil,
				hasValue: false,
				changed:  false,
			},
		},
		hasGroup: BoolValueHolder{
			valueHolder: valueHolder{
				value:    false,
				hasValue: false,
				changed:  false,
			},
		},
		organization: ExternalUserOrganizationValueHolder{
			valueHolder: valueHolder{
				value: &ExternalUserOrganization{
					ID: 0,
					name: StringValueHolder{
						valueHolder: valueHolder{
							value:    "organization",
							hasValue: true,
							changed:  false,
						},
					},
				},
				hasValue: true,
				changed:  true,
			},
		},
	}

	clone := user.Clone(false)

	// Assert that the clone is equal
	assert.True(t, user.IsDeepEqual(clone))

	// Ensure that changing one property doesn't affect the other struct
	clone.SetGroup("next group")
	assert.True(t, clone.GroupChanged())
	assert.False(t, user.GroupChanged())
	assert.False(t, user.IsDeepEqual(clone))
}

func TestClone_NestedPropertiesAreImmutable(t *testing.T) {
	user := &ExternalUser{
		ID:       "1",
		Username: "test",
		Email:    "test@gmail.com",
		Status:   "test",
		enabled: BoolValueHolder{
			valueHolder: valueHolder{
				value:    true,
				hasValue: true,
				changed:  false,
			},
		},
		firstName: StringValueHolder{
			valueHolder: valueHolder{
				value:    "firstname",
				hasValue: true,
				changed:  false,
			},
		},
		lastName: StringValueHolder{
			valueHolder: valueHolder{
				value:    "lastname",
				hasValue: true,
				changed:  true,
			},
		},
		group: StringValueHolder{
			valueHolder: valueHolder{
				value:    nil,
				hasValue: false,
				changed:  false,
			},
		},
		hasGroup: BoolValueHolder{
			valueHolder: valueHolder{
				value:    false,
				hasValue: false,
				changed:  false,
			},
		},
		organization: ExternalUserOrganizationValueHolder{
			valueHolder: valueHolder{
				value: &ExternalUserOrganization{
					ID: 0,
					name: StringValueHolder{
						valueHolder: valueHolder{
							value:    "organization",
							hasValue: true,
							changed:  false,
						},
					},
				},
				hasValue: true,
				changed:  true,
			},
		},
	}

	clone := user.Clone(false)

	// Assert that the clone is equal
	assert.True(t, user.IsDeepEqual(clone))

	// Ensure that changing one property doesn't affect the other struct
	assert.NotSame(t, user.Organization(), clone.Organization())
	clone.Organization().SetName("organization2")
	assert.True(t, clone.Organization().NameChanged())
	assert.False(t, user.Organization().NameChanged())
	assert.False(t, user.IsDeepEqual(clone))
}

func TestClone_ResetsChangedProperties(t *testing.T) {
	user := &ExternalUser{
		ID:       "1",
		Username: "test",
		Email:    "test@gmail.com",
		Status:   "test",
		enabled: BoolValueHolder{
			valueHolder: valueHolder{
				value:    true,
				hasValue: true,
				changed:  true,
			},
		},
		firstName: StringValueHolder{
			valueHolder: valueHolder{
				value:    "firstname",
				hasValue: true,
				changed:  false,
			},
		},
		lastName: StringValueHolder{
			valueHolder: valueHolder{
				value:    "lastname",
				hasValue: true,
				changed:  true,
			},
		},
		group: StringValueHolder{
			valueHolder: valueHolder{
				value:    nil,
				hasValue: false,
				changed:  false,
			},
		},
		hasGroup: BoolValueHolder{
			valueHolder: valueHolder{
				value:    false,
				hasValue: false,
				changed:  false,
			},
		},
		organization: ExternalUserOrganizationValueHolder{
			valueHolder: valueHolder{
				value: &ExternalUserOrganization{
					ID: 0,
					name: StringValueHolder{
						valueHolder: valueHolder{
							value:    "organization",
							hasValue: true,
							changed:  true,
						},
					},
				},
				hasValue: true,
				changed:  true,
			},
		},
	}

	clone := user.Clone(true)

	// ensure that changes were reset
	assert.False(t, clone.Organization().NameChanged())
	assert.True(t, user.Organization().NameChanged())
	assert.False(t, clone.OrganizationChanged())
	assert.True(t, user.OrganizationChanged())
	assert.False(t, clone.LastNameChanged())
	assert.True(t, user.LastNameChanged())
	assert.False(t, clone.EnabledChanged())
	assert.True(t, user.EnabledChanged())
}
