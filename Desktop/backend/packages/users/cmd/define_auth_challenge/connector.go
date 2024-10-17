//go:build !test
// +build !test

package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func connector(event *events.CognitoEventUserPoolsDefineAuthChallenge) (*events.CognitoEventUserPoolsDefineAuthChallenge, error) {
	return handler(event)
}

func main() {
	lambda.Start(connector)
}
