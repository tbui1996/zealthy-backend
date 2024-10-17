package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func Test_parse(t *testing.T) {
	tests := []struct {
		event               events.CognitoEventUserPoolsPostAuthentication
		environment         string
		expectedParsedEvent ParsedEvent
	}{
		{ // one group exists
			event: events.CognitoEventUserPoolsPostAuthentication{
				Request: events.CognitoEventUserPoolsPostAuthenticationRequest{
					UserAttributes: map[string]string{"custom:groups": "[\"internals_program_manager.dev-developer\"]"},
				},
			},
			environment: "dev-developer",
			expectedParsedEvent: ParsedEvent{
				OktaSonarGroups: map[string][]string{"dev-developer": {"internals_program_manager"}},
			},
		},
		{ // multiple groups exist, non sonar groups are ignored
			event: events.CognitoEventUserPoolsPostAuthentication{
				Request: events.CognitoEventUserPoolsPostAuthenticationRequest{
					UserAttributes: map[string]string{"custom:groups": "[\"internals_program_manager.dev\", \"some_other_group_that_doesnt_start_with_internals\"]"},
				},
			},
			environment: "dev",
			expectedParsedEvent: ParsedEvent{
				OktaSonarGroups: map[string][]string{"dev": {"internals_program_manager"}},
			},
		},
		{ // dev group is ignored
			event: events.CognitoEventUserPoolsPostAuthentication{
				Request: events.CognitoEventUserPoolsPostAuthenticationRequest{
					UserAttributes: map[string]string{"custom:groups": "[\"internals_program_manager.dev\", \"internals_program_manager.prod\"]"},
				},
			},
			environment: "prod",
			expectedParsedEvent: ParsedEvent{
				OktaSonarGroups: map[string][]string{"prod": {"internals_program_manager"}},
			},
		},
		{ // no groups exist
			event: events.CognitoEventUserPoolsPostAuthentication{
				Request: events.CognitoEventUserPoolsPostAuthenticationRequest{
					UserAttributes: map[string]string{},
				},
			},
			environment: "dev",
			expectedParsedEvent: ParsedEvent{
				OktaSonarGroups: nil,
			},
		},
		{ // invalid json
			event: events.CognitoEventUserPoolsPostAuthentication{
				Request: events.CognitoEventUserPoolsPostAuthenticationRequest{
					UserAttributes: map[string]string{"custom:groups": "\"internals_program_manager.dev\"]"},
				},
			},
			environment: "dev",
			expectedParsedEvent: ParsedEvent{
				OktaSonarGroups: nil,
			},
		},
	}

	for _, test := range tests {
		actualParsedEvent := parse(test.event, test.environment, zaptest.NewLogger(t))
		assert.Equal(t, test.expectedParsedEvent, actualParsedEvent)
	}
}
