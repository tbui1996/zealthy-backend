package interfaces

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt"
)

type Jwt interface {
	ParseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error)
	ParseWithClaims(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error)
}
