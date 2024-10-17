package authorizer

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt"
)

type Jwt struct{}

func (j *Jwt) ParseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error) {
	return jwt.ParseRSAPublicKeyFromPEM(key)
}

func (j *Jwt) ParseWithClaims(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, keyFunc)
}
