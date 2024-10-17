package authorizer

import "github.com/golang-jwt/jwt"

// type for sonar user access token claims.
type AccessTokenClaims struct {
	Sub           string   `json:"sub"`
	CognitoGroups []string `json:"cognito:groups"`
	Iss           string   `json:"iss"`
	Version       int      `json:"version"`
	ClientID      string   `json:"client_id"`
	OriginJti     string   `json:"origin_jti"`
	TokenUse      string   `json:"token_use"`
	Scope         string   `json:"scope"`
	AuthTime      int      `json:"auth_time"`
	Exp           int      `json:"exp"`
	Iat           int      `json:"iat"`
	Jti           string   `json:"jti"`
	Username      string   `json:"username"`
	jwt.StandardClaims
}
