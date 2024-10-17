package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/stretchr/testify/assert"
)

func Test_getAttribute(t *testing.T) {
	tests := []struct {
		inputName       string
		inputAttributes []*cognitoidentityprovider.AttributeType
		expected        *string
	}{
		{
			inputName: "name_that_exists",
			inputAttributes: []*cognitoidentityprovider.AttributeType{
				{
					Name:  aws.String("name_that_exists"),
					Value: aws.String("value"),
				},
			},
			expected: aws.String("value"),
		},
		{
			inputName: "name_that_does_not_exist",
			inputAttributes: []*cognitoidentityprovider.AttributeType{
				{
					Name:  aws.String("name"),
					Value: aws.String("value"),
				},
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		actual := getAttribute(test.inputName, test.inputAttributes)
		assert.Equal(t, test.expected, actual)
	}
}

func Test_getGroupName(t *testing.T) {
	tests := []struct {
		input    *cognitoidentityprovider.AdminListGroupsForUserOutput
		expected *string
	}{
		{
			input: &cognitoidentityprovider.AdminListGroupsForUserOutput{
				Groups: []*cognitoidentityprovider.GroupType{
					{
						GroupName: aws.String("test_group"),
					},
				},
			},
			expected: aws.String("test_group"),
		},
		{
			input:    &cognitoidentityprovider.AdminListGroupsForUserOutput{},
			expected: aws.String("no_group"),
		},
	}

	for _, test := range tests {
		actual := getGroupName(test.input)
		assert.Equal(t, test.expected, actual)
	}
}
