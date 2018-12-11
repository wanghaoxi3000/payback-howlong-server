package controllers

import (
	"crypto/tls"
	"errors"
	"fmt"
	"howlong/models"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	jwt "github.com/dgrijalva/jwt-go"
)

type sessionReponse struct {
	Openid  string
	Session string `json:"session_key"`
	Unionid string
	errcode int
	errmsg  string
}

// code2Session: use weixin api to get user's openID and session key
func code2Session(jsCode string) *sessionReponse {
	reqAddr := fmt.Sprintf("%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		beego.AppConfig.String("weiLoginApi"),
		beego.AppConfig.String("weiAppid"),
		beego.AppConfig.String("weiSecret"),
		jsCode)
	beego.Debug("User login get openID and session from: ", reqAddr)

	reqData := &sessionReponse{}
	req := httplib.Get(reqAddr).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	req.ToJSON(reqData)

	return reqData
}

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

	rep := code2Session(loginInfo.LoginCode)
	if rep.errcode != 0 {
		beego.Error("get openID error, code: ", rep.errcode)
		o.ServerError(errors.New("Get openID error"), notAvailable)
	}

	user, err := models.UpdateUserByOpenID(rep.Openid, rep.Session)
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
		beego.Info(httpBadRequest, httpUnauthorized, httpPaymentRequried)
		o.ServerError("no auth user", httpUnauthorized)
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
