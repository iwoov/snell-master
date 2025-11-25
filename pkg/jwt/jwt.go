package jwtutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 封装 JWT 负载，包含用户角色信息。
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 创建带有过期时间的 JWT。
func GenerateToken(userID uint, username, role, secret string, expireHours int) (string, error) {
	if userID == 0 {
		return "", errors.New("user id is required")
	}
	if secret == "" {
		return "", errors.New("jwt secret is required")
	}
	if expireHours <= 0 {
		return "", errors.New("expire hours must be greater than zero")
	}

	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expireHours) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析并返回自定义 Claims。
func ParseToken(tokenString, secret string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token is required")
	}
	if secret == "" {
		return nil, errors.New("jwt secret is required")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// ValidateToken 仅返回 token 是否有效。
func ValidateToken(tokenString, secret string) bool {
	_, err := ParseToken(tokenString, secret)
	return err == nil
}
