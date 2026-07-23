package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTCustomClaims struct {
	UserID string   `json:"uid"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

// ParseToken parses and validates a JWT token string, returning the claims.
func ParseToken(tokenStr string, secret string) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

// GenerateToken generates a signed JWT string for a user.
func GenerateToken(userID, email string, roles []string, secret string, duration time.Duration) (string, error) {
	claims := &JWTCustomClaims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
