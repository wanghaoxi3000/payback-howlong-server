package controllers

import (
	"howlong/models"
	"strconv"

	"github.com/astaxie/beego"
)

type CreditController struct {
	authController
}

// Create : Create a credit card
// @router / [post]
func (o *CreditController) Create() {
	var creditInfo creditSerializer
	o.UnserializeStruct(&creditInfo)

	creditModel := creditInfo.unserializer()
	creditModel.User = o.user

	id, err := models.AddCredit(creditModel)
	if err != nil {
		o.ServerError(err, httpBadRequest)
	}

	o.Data["json"] = map[string]int64{"id": id}
	o.Ctx.Output.SetStatus(httpCreated)
	o.ServeJSON()
}

// Retrieve : Retrieve a credit card info
// @router /:creditID [get]
func (o *CreditController) Retrieve() {
	creditID := o.Ctx.Input.Param(":creditID")
	intid, _ := strconv.ParseInt(creditID, 10, 64)

	c, err := models.GetUserCreditByID(o.user, intid)
	if err != nil {
		beego.Warning("Get credit ID", creditID, "error:", err)
		o.ServerError("Not found", httpNotFound)
		return
	}

	creditSerial := new(creditSerializer)
	creditSerial.serializer(c)

	o.Data["json"] = creditSerial
	o.ServeJSON()
}

// List : List all credit card info
// @router / [get]
func (o *CreditController) List() {
	var credits []*models.Credit
	num, err := models.GetUserAllCredit(o.user, &credits)
	if err != nil {
		beego.Warning("Get credits error, user ID", o.user.Id, "error:", err)
		o.ServerError("Get fail", notAvailable)
		return
	}

	length := int(num)
	creditsStructs := make([]*creditSerializer, length)
	for i := range credits {
		creditsStructs[i] = new(creditSerializer)
		creditsStructs[i].serializer(credits[i])
	}

	var sortFlag int
	for k := range creditsStructs {
		sortFlag = k + 1
		for i := sortFlag; i < length; i++ {
			if creditsStructs[k].DateDetail.IntervalPay < creditsStructs[i].DateDetail.IntervalPay {
				creditsStructs[k], creditsStructs[i] = creditsStructs[i], creditsStructs[k]
			}
		}

	}

	o.Data["json"] = creditsStructs
	o.ServeJSON()
}

// Update : Update credit card info
// @router /:creditID [put]
func (o *CreditController) Update() {
	creditID := o.Ctx.Input.Param(":creditID")
	intid, _ := strconv.ParseInt(creditID, 10, 64)

	_, err := models.GetUserCreditByID(o.user, intid)
	if err != nil {
		o.ServerError("Not found", httpNotFound)
		return
	}

	var creditInfo creditSerializer
	o.UnserializeStruct(&creditInfo)
	creditModel := creditInfo.unserializer()
	creditModel.Id = intid
	creditModel.User = o.user

	err = models.UpdateUserCredit(creditModel)
	if err != nil {
		o.ServerError(err, httpBadRequest)
		return
	}

	creditInfo.serializer(creditModel)
	o.Data["json"] = creditInfo
	o.ServeJSON()
}

// Destroy : Remove a credit card info from database
// @router /:creditID [delete]
func (o *CreditController) Destroy() {
	creditID := o.Ctx.Input.Param(":creditID")
	intid, _ := strconv.ParseInt(creditID, 10, 64)

	c, err := models.GetUserCreditByID(o.user, intid)
	if err != nil {
		beego.Warning("Get credits error, user ID", o.user.Id, "error:", err)
		o.ServerError("Not found", httpNotFound)
		return
	}

	err = models.DeleteCredit(c)
	if err != nil {
		beego.Error("Delete credit ID:", creditID, "error:", err)
		o.ServerError("Delete fail", notAvailable)
	}

	o.Data["json"] = "delete success!"
	o.ServeJSON()
}
