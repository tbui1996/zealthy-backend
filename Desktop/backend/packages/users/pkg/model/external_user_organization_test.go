package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildExternalUserOrganization(t *testing.T) {
	baseInput := &BuildExternalUserOrganizationInput{
		ID: 0,
	}

	tests := []struct {
		builder  *ExternalUserOrganizationBuilder
		expected *ExternalUserOrganization
	}{
		{
			builder: BuildExternalUserOrganization(baseInput),
			expected: &ExternalUserOrganization{
				ID: 0,
				name: StringValueHolder{
					valueHolder: valueHolder{
						value:    nil,
						hasValue: false,
						changed:  false,
					},
				},
			},
		},
		{
			builder: BuildExternalUserOrganization(baseInput).WithName("N"),
			expected: &ExternalUserOrganization{
				ID: 0,
				name: StringValueHolder{
					valueHolder: valueHolder{
						value:    "N",
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

func TestModifyExternalUserOrganization(t *testing.T) {
	baseInput := &BuildExternalUserOrganizationInput{
		ID: 0,
	}

	tests := []struct {
		executor func() *ExternalUserOrganization
		expected *ExternalUserOrganization
	}{
		{
			executor: func() *ExternalUserOrganization {
				base := BuildExternalUserOrganization(baseInput).Value()
				base.SetName("T")
				return base
			},
			expected: &ExternalUserOrganization{
				ID: 0,
				name: StringValueHolder{
					valueHolder: valueHolder{
						value:    "T",
						hasValue: true,
						changed:  true,
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
