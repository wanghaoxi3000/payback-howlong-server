package controllers

import (
	"encoding/json"
	"errors"
	"howlong/models"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

const (
	badRequest = 400
	notAuthenticated

	innerError   = 500
	notAvailable = 503
)

type serializer interface {
	Validate() (map[string]string, error)
}

func validateStruct(obj interface{}) (map[string]string, error) {
	valid := validation.Validation{}
	validInfo := make(map[string]string)

	if info, err := valid.Valid(obj); err != nil {
		return nil, err
	} else if !info {
		for _, e := range valid.Errors {
			validInfo[e.Key] = e.Message
		}
	}

	return validInfo, nil
}

type baseController struct {
	beego.Controller

	user *models.User
}

type noAuthError interface {
	ServerNoAuth()
}

// Prepare : Implemented Prepare method for BaseController.
func (o *baseController) Prepare() {
	auth := strings.SplitN(o.Ctx.Input.Header("Authorization"), "Bearer ", 2)
	if len(auth) > 1 {
		token := strings.TrimSpace(auth[1])
		o.user = checkToken(token)
	}

	if o.user == nil {
		if app, ok := o.AppController.(noAuthError); ok {
			app.ServerNoAuth()
		}
	}
}

// ServerError : Deal server message
func (o *baseController) ServerError(msg interface{}, code int) {
	errMsg := make(map[string]string)

	switch v := msg.(type) {
	case string:
		errMsg["error"] = v

	case error:
		errMsg["error"] = v.Error()

	case map[string]string:
		errMsg = v

	default:
		beego.Error("server error, unexpected type %T", v)
		o.Abort("500")
	}

	beego.Warning("server err:", errMsg, "code:", code)
	o.Data["json"] = errMsg
	o.Ctx.Output.SetStatus(code)
	o.ServeJSON()
	o.StopRun()
}

// UnserializeStruct : Unserialize to struct
func (o *baseController) UnserializeStruct(model serializer) error {
	if err := json.Unmarshal(o.Ctx.Input.RequestBody, model); err != nil {
		o.ServerError(err, badRequest)
		return err
	}

	if validateRet, err := model.Validate(); err != nil {
		beego.Error("validate struct %T error: %v", model, err)
		o.ServerError(errors.New("unknown data"), badRequest)
		return err
	} else if len(validateRet) > 0 {
		o.ServerError(validateRet, badRequest)
		return errors.New("invalid data")
	}

	return nil
}

type authController struct {
	baseController
}

func (o *authController) ServerNoAuth() {
	o.ServerError("no auth user", notAuthenticated)
}
