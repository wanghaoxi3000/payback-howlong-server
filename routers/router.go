package routers

import (
	"howlong/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/users",
			beego.NSRouter("/login", &controllers.UserController{}, "post:Login"),

			beego.NSInclude(&controllers.UserController{}),
		),
		beego.NSNamespace("/credits",
			beego.NSInclude(&controllers.CreditController{}),
		),
	)
	beego.AddNamespace(ns)
}
