package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID  int    `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID int, email string, isAdmin bool, secret string, expiry time.Duration) (string, error) {
	claims := Claims{
		UserID:  userID,
		Email:   email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateTokenWithClaims creates a token from pre-built claims
func GenerateTokenWithClaims(claims *Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}