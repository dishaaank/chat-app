package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your-256-bit-secret")

func GenerateJWT(username string) (string, error) {
	claims := &CustomClaims{
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   username, // still include in 'sub'
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func VerifyJWT(tokenStr string) (*jwt.Token, *CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-256-bit-secret"), nil // Replace with your real key
	})
	if err != nil {
		return nil, nil, err
	}
	return token, claims, nil
}

// CustomClaims defines your JWT payload structure
type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
