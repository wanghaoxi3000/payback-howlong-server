package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type UserController struct {
	beego.Controller
}

type loginInfo struct {
	OpenID string
}

// Create : Create a credit card
// @router / [post]
func (o *UserController) Create() {
	var ob loginInfo

	if err := json.Unmarshal(o.Ctx.Input.RequestBody, &ob); err == nil {
		logs.Info("Json %v", ob)
		o.Data["json"] = ob
	} else {
		o.Data["json"] = err.Error()
	}
	o.ServeJSON()
}
