package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func handler(event *events.CognitoEventUserPoolsDefineAuthChallenge) (*events.CognitoEventUserPoolsDefineAuthChallenge, error) {
	log.Printf("Define Auth Challenge: %+v\n", event)

	event.Response.IssueTokens = true
	event.Response.FailAuthentication = false

	return event, nil
}
