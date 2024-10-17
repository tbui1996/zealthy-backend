package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
)

type Logger struct {
	Fields         *LoggerFields
	SkipValidation bool // We're skipping validation on Cognito and Authorizer events
	Error          error
}

type LoggerFields struct {
	SourceIP  string `json:"sourceIP,omitempty"`
	RouteKey  string `json:"routeKey,omitempty"`
	RequestID string `json:"requestID,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
	UserID    string `json:"userID,omitempty"`
	Email     string `json:"email,omitempty"`
}

func (fields *LoggerFields) ToString() (string, error) {
	bytes, err := json.Marshal(fields)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (fields *LoggerFields) FromString(fieldsStr string) error {
	err := json.Unmarshal([]byte(fieldsStr), &fields)
	if err != nil {
		return fmt.Errorf("unable to parse logger fields from message attributes")
	}

	return nil
}

func (fields *LoggerFields) FromSQSMessage(message events.SQSMessage) (*zap.Logger, error) {
	loggerFieldsJson, ok := message.MessageAttributes["LoggerFields"]
	if !ok {
		return nil, fmt.Errorf("couldn't get logger fields from message attributes")
	}

	err := fields.FromString(*loggerFieldsJson.StringValue)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse logger fields from string: " + err.Error())
	}

	logger, err := fields.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("couldn't get logger: " + err.Error())
	}

	return logger, nil
}

func (fields *LoggerFields) Validate() error {
	var missingRequiredFields []string

	// check to make sure that we can identify the requester
	if fields.Email == "" && fields.UserID == "" {
		missingRequiredFields = append(missingRequiredFields, "user id or email")
	}

	if fields.SourceIP == "" {
		missingRequiredFields = append(missingRequiredFields, "sourceIP")
	}

	if fields.RouteKey == "" {
		missingRequiredFields = append(missingRequiredFields, "routeKey")
	}

	if fields.RequestID == "" {
		missingRequiredFields = append(missingRequiredFields, "requestID")
	}

	if fields.UserAgent == "" {
		missingRequiredFields = append(missingRequiredFields, "userAgent")
	}

	if len(missingRequiredFields) > 0 {
		log.Printf("WARNING: missing highly suggested fields to create a logger: %s", strings.Join(missingRequiredFields, ", "))
		// TODO: Once ALL the send and receive packages have been updated, enforce all base logger fields
		// return fmt.Errorf("missing required fields to create a logger: %s", strings.Join(missingRequiredFields, ", "))
	}

	return nil
}

func NewAPIGatewayWebsocketProxyRequestLoggerFields(event events.APIGatewayWebsocketProxyRequest) (*LoggerFields, error) {
	config := &LoggerFields{
		SourceIP:  event.RequestContext.Identity.SourceIP,
		RouteKey:  event.RequestContext.RouteKey,
		RequestID: event.RequestContext.RequestID,
		UserAgent: event.RequestContext.Identity.UserAgent,
	}

	if userID, userIDOk := event.RequestContext.Authorizer.(map[string]interface{})["userID"].(string); userIDOk {
		config.UserID = userID
	}

	if email, emailOk := event.RequestContext.Authorizer.(map[string]interface{})["email"].(string); emailOk {
		config.Email = email
	}

	return config, nil
}

func NewAPIGatewayCustomAuthorizerRequestTypeRequestLoggerFields(event events.APIGatewayCustomAuthorizerRequestTypeRequest) (*LoggerFields, error) {
	config := &LoggerFields{
		SourceIP:  event.RequestContext.Identity.SourceIP,
		RequestID: event.RequestContext.RequestID,
		RouteKey:  event.MethodArn,
		UserAgent: event.Headers["User-Agent"],
	}

	return config, nil
}

func NewAPIGatewayV2HTTPRequestLoggerFields(event events.APIGatewayV2HTTPRequest) (*LoggerFields, error) {
	config := &LoggerFields{
		SourceIP:  event.RequestContext.HTTP.SourceIP,
		RouteKey:  event.RequestContext.RouteKey,
		RequestID: event.RequestContext.RequestID,
		UserAgent: event.RequestContext.HTTP.UserAgent,
	}

	if userID, userIDOk := event.RequestContext.Authorizer.Lambda["userID"].(string); userIDOk {
		config.UserID = userID
	}

	if email, emailOk := event.RequestContext.Authorizer.Lambda["email"].(string); emailOk {
		config.Email = email
	}

	return config, nil
}

func NewCognitoEventUserPoolsPostConfirmationLoggerFields(event events.CognitoEventUserPoolsPostConfirmation) (*LoggerFields, error) {
	config := &LoggerFields{
		UserID: event.UserName,
	}

	if email, emailOk := event.Request.UserAttributes["email"]; emailOk {
		config.Email = email
	}

	return config, nil
}

func NewCognitoEventUserPoolsPreSignupLoggerFields(event events.CognitoEventUserPoolsPreSignup) (*LoggerFields, error) {
	config := &LoggerFields{
		UserID: event.UserName,
	}

	if email, emailOk := event.Request.UserAttributes["email"]; emailOk {
		config.Email = email
	}

	return config, nil
}

func LoggerFieldsFromEvent(event interface{}) Logger {
	var logger Logger
	logger.SkipValidation = false

	switch event := event.(type) {
	case events.APIGatewayCustomAuthorizerRequestTypeRequest:
		logger.Fields, logger.Error = NewAPIGatewayCustomAuthorizerRequestTypeRequestLoggerFields(event)
		logger.SkipValidation = true
	case events.APIGatewayV2HTTPRequest:
		logger.Fields, logger.Error = NewAPIGatewayV2HTTPRequestLoggerFields(event)
	case events.APIGatewayWebsocketProxyRequest:
		logger.Fields, logger.Error = NewAPIGatewayWebsocketProxyRequestLoggerFields(event)
	case events.CognitoEventUserPoolsPostConfirmation:
		logger.Fields, logger.Error = NewCognitoEventUserPoolsPostConfirmationLoggerFields(event)
		logger.SkipValidation = true
	case events.CognitoEventUserPoolsPreSignup:
		logger.Fields, logger.Error = NewCognitoEventUserPoolsPreSignupLoggerFields(event)
		logger.SkipValidation = true
	case events.CloudWatchEvent:
		logger.Fields = &LoggerFields{}
		logger.Error = nil
		logger.SkipValidation = true
	default:
		logger.Fields = nil
		logger.Error = fmt.Errorf("could not get logger from event: unsupported event type %T", event)
	}

	return logger
}
