package controllers

import (
	"encoding/json"
	"howlong/models"
	"howlong/utils"

	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

type loginSerializer struct {
	OpenID string
}

// Create : Create a credit card
// @router / [post]
func (o *UserController) Create() {
	var (
		loginInfo *loginSerializer
		user      *models.User
		err       error
	)

	loginInfo = new(loginSerializer)
	if err = json.Unmarshal(o.Ctx.Input.RequestBody, loginInfo); err != nil {
		o.Data["json"] = err.Error()
		o.ServeJSON()
		return
	}

	if user, err = models.UpdateUserByOpenID(loginInfo.OpenID, utils.RandString(12)); err != nil {
		beego.Error("Retrieve user error, openid:", loginInfo.OpenID, "error:", err)
		return
	}

	var token string
	if token, err = GenToken(user.SessionKey); err != nil {
		beego.Error("Gen JWT token error: ", err)
		return
	}

	beego.Debug("user ID: ", user.Id, "openid: ", user.OpenId, "JWT token: ", token)
	CheckToken(token, user.SessionKey)
	o.Data["json"] = user
	o.ServeJSON()
}
