package controllers

import (
	"howlong/models"
	"strconv"

	"github.com/astaxie/beego/logs"
)

type CreditController struct {
	authController
}

// Create : Create a credit card
// @router / [post]
func (o *CreditController) Create() {
	var creditInfo creditSerializer
	o.UnserializeStruct(&creditInfo)
	logs.Debug(creditInfo)

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
		logs.Warning("Get credit ID %v error: %v", creditID, err.Error())
		o.Ctx.Output.SetStatus(404)
		o.Ctx.Output.Body([]byte("Not found"))
		return
	}

	o.Data["json"] = t
	o.ServeJSON()
}

// List : List all credit card info
// @router / [get]
func (o *CreditController) List() {
	// if creditList, err := models.GetSortedCredit(); err != nil {
	// 	logs.Error("Get all credit from database error: %v", err.Error())
	// 	o.Abort("500")
	// } else {
	// 	o.Data["json"] = creditList
	// 	o.ServeJSON()
	// }
	o.Data["json"] = o.user
	o.ServeJSON()
}

// Update : Update credit card info
// @router /:creditID [put]
func (o *CreditController) Update() {
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

	creditID := o.Ctx.Input.Param(":creditID")
	intid, _ := strconv.ParseInt(creditID, 10, 64)
	credit.Id = intid

	err := models.UpdateCreditById(credit)
	if err != nil {
		logs.Warning("Update credit ID %v error: %v", creditID, err.Error())
		o.Ctx.Output.SetStatus(404)
		o.Ctx.Output.Body([]byte("Not found"))
		return
	}

	o.Data["json"] = credit
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
