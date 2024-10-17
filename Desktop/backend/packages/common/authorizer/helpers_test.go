package authorizer

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestStripBearer(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Bearer asdfasdfasdf.asdfasdfasdf.asdfasdfasdf",
			expected: "asdfasdfasdf.asdfasdfasdf.asdfasdfasdf",
		},
		{
			input:    "   Bearer     asdfasdfasdf.asdfasdfasdf.asdfasdfasdf       ",
			expected: "asdfasdfasdf.asdfasdfasdf.asdfasdfasdf",
		},
		{
			input:    "   asdfasdfasdf.asdfasdfasdf.asdfasdfasdf   ",
			expected: "asdfasdfasdf.asdfasdfasdf.asdfasdfasdf",
		},
	}

	for _, test := range tests {
		actualStr := StripBearer(test.input)
		assert.Equal(t, test.expected, actualStr)
	}
}

func TestGetAuthorizationToken(t *testing.T) {
	tests := []struct {
		input    map[string]string
		expected string
	}{
		{
			input: map[string]string{
				"authorization": "token",
			},
			expected: "token",
		},
		{
			input: map[string]string{
				"Authorization": "token",
			},
			expected: "token",
		},
		{
			input:    map[string]string{},
			expected: "",
		},
	}

	for _, test := range tests {
		actualStr := GetAuthorizationToken(test.input)
		assert.Equal(t, test.expected, actualStr)
	}
}

func TestGetToken(t *testing.T) {
	tests := []struct {
		input    events.APIGatewayCustomAuthorizerRequestTypeRequest
		expected string
	}{
		{
			input: events.APIGatewayCustomAuthorizerRequestTypeRequest{
				Headers: map[string]string{
					"authorization": "token",
				},
			},
			expected: "token",
		},
		{
			input: events.APIGatewayCustomAuthorizerRequestTypeRequest{
				Headers: map[string]string{
					"Authorization": "token",
				},
			},
			expected: "token",
		},
		{
			input: events.APIGatewayCustomAuthorizerRequestTypeRequest{
				QueryStringParameters: map[string]string{
					"authorization": "token",
				},
			},
			expected: "token",
		},
		{
			input:    events.APIGatewayCustomAuthorizerRequestTypeRequest{},
			expected: "",
		},
	}

	for _, test := range tests {
		actualStr := GetToken(test.input)
		assert.Equal(t, test.expected, actualStr)
	}
}
