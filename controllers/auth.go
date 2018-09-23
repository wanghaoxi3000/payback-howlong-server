package controllers

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	jwt "github.com/dgrijalva/jwt-go"
)

// GenToken : Generate JWT token, expire time 1 day
func GenToken(key string) (jwtToken string, err error) {
	claims := &jwt.StandardClaims{
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(beego.AppConfig.DefaultInt("jwtExpireHour", 1))).Unix(),
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
		nowTime := time.Now()
		if value, ok := claims["nbf"].(float64); ok {
			if !nowTime.After(time.Unix(int64(value), 0)) {
				logs.Error("token nbf err, token: %s now time: %s", time.Unix(int64(value), 0), nowTime)
				return false
			}
		} else {
			logs.Error("parse token nbf value error: %v", claims["nbf"])
			return false
		}

		if value, ok := claims["exp"].(float64); ok {
			if nowTime.After(time.Unix(int64(value), 0)) {
				logs.Error("token nbf err, token: %s now time: %s", time.Unix(int64(value), 0), nowTime)
				return false
			}
		} else {
			logs.Error("parse token exp value error: %v", claims["exp"])
			return false
		}
	} else {
		logs.Error(err)
		return false
	}

	return true
}
