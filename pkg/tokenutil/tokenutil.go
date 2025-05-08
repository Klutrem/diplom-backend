package tokenutil

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	jwt.RegisteredClaims
	ID string `json:"id"`
}

func ParseToken(requestToken string, secret string) (*jwt.Token, *claims, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(secret))
	if err != nil {
		return nil, nil, err
	}

	claims := &claims{}

	token, err := jwt.ParseWithClaims(requestToken, claims, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return key, nil
	})
	if err != nil {
		return nil, nil, err
	}

	return token, claims, nil
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	token, _, err := ParseToken(requestToken, secret)
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}

func ExtractIDFromClaims(claims *claims) (string, error) {
	if claims.ID != "" {
		return claims.ID, nil
	}

	return "", errors.New("invalid token")
}
