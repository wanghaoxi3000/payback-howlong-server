package controllers

import (
	"howlong/models"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type CreditController struct {
	beego.Controller
}

// Create : Create a credit card
// @router / [post]
func (o *CreditController) Create() {
	credit := new(models.Credit)

	if err := o.ParseForm(credit); err != nil {
		logs.Error("Parse credit struct error: %v", err.Error())
		o.Ctx.Output.SetStatus(400)
		o.Ctx.Output.Body([]byte("Request data error"))
		return
	}

	if validateRet, err := credit.Validate(); err != nil {
		logs.Error("validate credit struct error: %v", err.Error())
		o.Abort("500")
	} else if len(validateRet) > 0 {
		o.Data["json"] = validateRet
		o.Ctx.Output.SetStatus(400)
		o.ServeJSON()
		return
	}

	if id, err := models.AddCredit(credit); err == nil {
		o.Data["json"] = map[string]int64{"id": id}
		o.ServeJSON()
	} else {
		logs.Error("Save credit struct error: %v", err.Error())
		o.Abort("500")
	}
}

// Retrieve : Retrieve a credit card info
// @router /:creditID [get]
func (o *CreditController) Retrieve() {
	creditID := o.Ctx.Input.Param(":creditID")
	intid, _ := strconv.ParseInt(creditID, 10, 64)

	t, err := models.GetCreditById(intid)
	if err != nil {
		logs.Warning("GetCredit Id %v error: %v", creditID, err.Error())
		o.Ctx.Output.SetStatus(404)
		o.Ctx.Output.Body([]byte("Not found"))
		return
	}

	o.Data["json"] = t
	o.ServeJSON()
}

// Destroy : Remove a credit card info from database
// @router /:creditID [delete]
func (o *CreditController) Destroy() {
	creditID := o.Ctx.Input.Param(":creditID")
	intid, _ := strconv.ParseInt(creditID, 10, 64)

	if err := models.DeleteCredit(intid); err != nil {
		logs.Warning("Delete credit Id %v error: %v", creditID, err.Error())
		o.Ctx.Output.SetStatus(404)
		o.Ctx.Output.Body([]byte("Not found"))
		return
	}

	o.Data["json"] = "delete success!"
	o.ServeJSON()
}
