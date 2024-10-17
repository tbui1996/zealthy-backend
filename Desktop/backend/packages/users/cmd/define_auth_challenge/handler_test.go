package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		input    *events.CognitoEventUserPoolsDefineAuthChallenge
		expected *events.CognitoEventUserPoolsDefineAuthChallenge
	}{
		{
			input: &events.CognitoEventUserPoolsDefineAuthChallenge{},
			expected: &events.CognitoEventUserPoolsDefineAuthChallenge{
				Response: events.CognitoEventUserPoolsDefineAuthChallengeResponse{
					IssueTokens:        true,
					FailAuthentication: false,
				},
			},
		},
	}

	for _, test := range tests {
		actual, err := handler(test.input)
		assert.Equal(t, test.expected, actual)
		assert.Nil(t, err)
	}
}
