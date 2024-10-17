package authorizer

import (
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func StripBearer(bearerToken string) string {
	// given 'Bearer eyydjkfh...' return 'eyydjkfh...'
	tokenTrimmed := strings.TrimSpace(bearerToken)
	tokenWithoutBearer := strings.TrimPrefix(tokenTrimmed, "Bearer")
	token := strings.TrimSpace(tokenWithoutBearer)
	return token
}

func GetAuthorizationToken(headers map[string]string) string {
	// https://github.com/aws/aws-lambda-go/issues/117
	token, ok := headers["authorization"]
	if ok && token != "" {
		return StripBearer(token)
	}

	token, ok = headers["Authorization"]
	if ok && token != "" {
		return StripBearer(token)
	}

	return ""
}

func GetToken(event events.APIGatewayCustomAuthorizerRequestTypeRequest) string {
	token := GetAuthorizationToken(event.Headers)

	if token != "" {
		log.Println("found authorization token in header")
		return token
	}

	log.Println("did not receive token in headers, checking query parameters")
	token = event.QueryStringParameters["authorization"]
	if token != "" {
		log.Println("found authorization token in query parameters")
	}

	return token
}
