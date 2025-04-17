package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Info struct {
	UserGUID       string
	UserIP         string
	AccessTokenUID string
}

type Claims struct {
	jwt.StandardClaims
	UserGUID       string
	UserIP         string
	AccessTokenUID string
}

// GenerateToken generates jwt-access-token
func GenerateToken(info *Info, secretKey []byte, duration time.Duration) (string, error) {
	if info == nil {
		return "", fmt.Errorf("info is nil")
	}

	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		UserGUID:       info.UserGUID,
		UserIP:         info.UserIP,
		AccessTokenUID: info.AccessTokenUID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString(secretKey)
}

// VerifyToken validates jwt-access-token and returns its claims
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
