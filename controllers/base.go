package controllers

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GenToken : Generate JWT token, expire time 1 day
func GenToken(key string) (jwtToken string, err error) {
	claims := &jwt.StandardClaims{
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(24)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToken, err = token.SignedString([]byte(key))
	if err != nil {
		jwtToken = ""
		return
	}

	return
}

// CheckToken : Check JWT token is valid
func CheckToken(tokenString string, key string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(key), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["exp"], claims["nbf"])
	} else {
		fmt.Println(err)
	}

	return true
}
