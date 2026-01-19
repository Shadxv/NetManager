package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type TokenPayload struct {
	RolePermissions       int64    `json:"rolePermissions"`
	AdditionalPermissions []string `json:"additionalPermissions"`
	RequiresPasswordReset bool     `json:"requiresPasswordReset"`
	jwt.RegisteredClaims
}

func ValidateToken(tokenString string, secret string) (*TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenPayload); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
