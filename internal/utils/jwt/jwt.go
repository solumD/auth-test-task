package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Info struct {
	AccessTokenUID string
	UserIP         string
}

type Claims struct {
	jwt.StandardClaims
	AccessTokenUID string
	UserIP         string
}

// GenerateToken генерирует jwt-access-токен
func GenerateToken(info *Info, secretKey []byte, duration time.Duration) (string, error) {
	if info == nil {
		return "", fmt.Errorf("info is nil")
	}
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		AccessTokenUID: info.AccessTokenUID,
		UserIP:         info.UserIP,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString(secretKey)
}

// VerifyToken валидирует jwt-токен и возвращает его claims
func VerifyToken(tokenStr string, secretKey []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return secretKey, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
