package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type baseSerializer struct {
}

// Validate : Validate a serializer
func (s baseSerializer) Validate() (map[string]string, error) {
	validInfo := make(map[string]string)
	valid := validation.Validation{}

	if info, err := valid.Valid(s); err != nil {
		beego.Error("vailidate error: %v", err.Error())
		return nil, err
	} else if !info {
		for _, err := range valid.Errors {
			validInfo[err.Key] = err.Message
		}
	}
	fmt.Println("Go to here", validInfo, s, len(""))
	return validInfo, nil
}

// BaseController : Add some function from beego.Controller
type BaseController struct {
	beego.Controller
}

type errorInfo struct {
	Error string
}

// ServerError : Deal server message
func (o *BaseController) ServerError(msg interface{}, code int) {
	errMsg := make(map[string]string)

	switch v := msg.(type) {
	case error:
		beego.Debug("server error: %v code: %d", v.Error(), code)
		errMsg["error"] = v.Error()

	case map[string]string:
		errMsg = v

	default:
		beego.Error("server error, unexpected type %T", v)
		o.Abort("500")
	}

	o.Data["json"] = errMsg
	o.Ctx.Output.SetStatus(code)
	o.ServeJSON()
}
