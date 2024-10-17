package externalAuth

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/circulohealth/sonar-backend/packages/common/authorizer"
	"github.com/golang-jwt/jwt"
	"github.com/sethvargo/go-password/password"
)

var GetEmailFromToken = func(headers map[string]string) (string, error) {
	tokenStr := authorizer.GetAuthorizationToken(headers)

	svc := authorizer.Jwt{}

	verifyKey, err := svc.ParseRSAPublicKeyFromPEM([]byte(authorizer.OlivePublicKey))
	if err != nil {
		log.Fatal(err)
	}

	token, err := svc.ParseWithClaims(tokenStr, &authorizer.OliveClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err != nil || !token.Valid {
		log.Printf("Error parsing token: %v\n", err)
		return "", errors.New("error: invalid token")
	}

	if claims, ok := token.Claims.(*authorizer.OliveClaims); ok {
		return claims.Email, nil
	}

	return "", errors.New("error: parsing claims")
}

var CreatePassword = func() (string, error) {
	// Generate a password that is 64 characters long with 10 digits, 10 symbols,
	// allowing upper and lower case letters, disallowing repeat characters.
	characters := 64
	numDigits := 10
	numSymbols := 10
	pass, err := password.Generate(characters, numDigits, numSymbols, false /*noUpper*/, false /*allowRepeat*/)
	if err != nil {
		return "", err
	}

	return pass, nil
}

func ValidateInitiateAuthOutput(initiateAuth *cognitoidentityprovider.InitiateAuthOutput) error {
	if initiateAuth == nil {
		return errors.New("expected output to not be null")
	}

	if initiateAuth.AuthenticationResult == nil {
		return errors.New("expected authentication result to not be null")
	}

	if initiateAuth.AuthenticationResult.IdToken == nil {
		return errors.New("expected id token to not be null")
	}

	if initiateAuth.AuthenticationResult.AccessToken == nil {
		return errors.New("expected access token to not be null")
	}

	return nil
}
