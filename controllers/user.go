package controllers

import (
	"encoding/json"
	"howlong/models"
	"howlong/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type UserController struct {
	BaseController
}

type loginSerializer struct {
	baseSerializer

	OpenID string `valid:"MinSize(1); MaxSize(20)"`
}

// Login : Login a credit card
func (o *UserController) Login() {
	var (
		loginInfo loginSerializer
		user      *models.User
		err       error
	)

	if err = json.Unmarshal(o.Ctx.Input.RequestBody, &loginInfo); err != nil {
		o.ServerError(err, 400)
		return
	}

	valid := validation.Validation{}
	validInfo := make(map[string]string)
	if info, e := valid.Valid(loginInfo); e != nil {
		beego.Error("vailidate error: %v", e.Error())
		// return nil, err
	} else if !info {
		for _, e := range valid.Errors {
			validInfo[e.Key] = e.Message
		}
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
