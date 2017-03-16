package utils

import (
	jwt "github.com/dgrijalva/jwt-go"
)

// UserJWTClaims wraps jwt.StandardClaims and adds any additional claims
type UserJWTClaims struct {
	jwt.StandardClaims
}

// GenerateJWT creates a new JWT from the signingKey and userID
func GenerateJWT(signingKey []byte, userID string) (string, error) {
	claims := UserJWTClaims{
		jwt.StandardClaims{
			Issuer:  "discapi",
			Subject: userID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}
