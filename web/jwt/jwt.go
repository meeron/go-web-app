package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
	"web-app/shared"
)

var signingKey = []byte("THIS_IS_SECRET")

func Create(customClaims map[string]string) string {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().UTC().Add(time.Hour).Unix()

	for key, value := range customClaims {
		claims[key] = value
	}

	return shared.Unwrap(token.SignedString(signingKey))
}

func Validate(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}
