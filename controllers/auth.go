package controllers

import (
	"errors"
	"fmt"
	"howlong/models"
	"howlong/utils"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	jwt "github.com/dgrijalva/jwt-go"
)

// genToken : Generate JWT token, expire time 1 hour default
func genToken(userID int64, key string) (jwtToken string, err error) {
	aud := strconv.FormatInt(userID, 10)
	claims := &jwt.StandardClaims{
		Audience:  aud,
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

// checkToken : Check JWT token is valid and set user to controller
func checkToken(tokenString string) *models.User {
	var authUser *models.User

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("Unknown claims")
		}

		userID, err := strconv.ParseInt(claims["aud"].(string), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Parse claims aud error %v", err)
		}

		user, err := models.GetUserByID(userID)
		if err != nil {
			return nil, fmt.Errorf("Get seesion key error %v", err)
		}

		authUser = user
		return []byte(user.SessionKey), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return authUser
	}

	beego.Warning("token invalid:", err)
	return nil
}

type AuthController struct {
	baseController
}

// Login : User login to get token
// @router /login [post]
func (o *AuthController) Login() {
	var loginInfo loginSerializer
	o.UnserializeStruct(&loginInfo)

	user, err := models.UpdateUserByOpenID(loginInfo.LoginCode, utils.RandString(12))
	if err != nil {
		beego.Error("Retrieve user error, openid:", loginInfo.LoginCode, "error:", err)
		o.ServerError(fmt.Errorf("Find user %v error", loginInfo.LoginCode), httpBadRequest)
	}

	token, err := genToken(user.Id, user.SessionKey)
	if err != nil {
		beego.Error("Gen JWT token error: ", err)
		o.ServerError(errors.New("Gen token error"), notAvailable)
	}

	beego.Debug("user ID: ", user.Id, "openid: ", user.OpenId, "login")
	authInfo := make(map[string]string)
	authInfo["token"] = token
	o.Data["json"] = authInfo
	o.ServeJSON()
}

// RefreshToken : Refresh token expired time
// @router /refresh-token [get]
func (o *AuthController) RefreshToken() {
	if o.user == nil {
		o.Abort("401")
	}

	token, err := genToken(o.user.Id, o.user.SessionKey)
	if err != nil {
		beego.Error("Gen JWT token error: ", err)
		o.ServerError(errors.New("Gen token error"), notAvailable)
	}

	authInfo := make(map[string]string)
	authInfo["token"] = token
	o.Data["json"] = authInfo
	o.ServeJSON()
}
