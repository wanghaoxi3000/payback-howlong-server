package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["howlong/controllers:CreditController"] = append(beego.GlobalControllerRouter["howlong/controllers:CreditController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:CreditController"] = append(beego.GlobalControllerRouter["howlong/controllers:CreditController"],
		beego.ControllerComments{
			Method: "Retrieve",
			Router: `/:creditID`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:CreditController"] = append(beego.GlobalControllerRouter["howlong/controllers:CreditController"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/:creditID`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:CreditController"] = append(beego.GlobalControllerRouter["howlong/controllers:CreditController"],
		beego.ControllerComments{
			Method: "Destroy",
			Router: `/:creditID`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:ObjectController"] = append(beego.GlobalControllerRouter["howlong/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:ObjectController"] = append(beego.GlobalControllerRouter["howlong/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:ObjectController"] = append(beego.GlobalControllerRouter["howlong/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:ObjectController"] = append(beego.GlobalControllerRouter["howlong/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:ObjectController"] = append(beego.GlobalControllerRouter["howlong/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:UserController"] = append(beego.GlobalControllerRouter["howlong/controllers:UserController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:UserController"] = append(beego.GlobalControllerRouter["howlong/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:UserController"] = append(beego.GlobalControllerRouter["howlong/controllers:UserController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:UserController"] = append(beego.GlobalControllerRouter["howlong/controllers:UserController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:UserController"] = append(beego.GlobalControllerRouter["howlong/controllers:UserController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:UserController"] = append(beego.GlobalControllerRouter["howlong/controllers:UserController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:UserController"] = append(beego.GlobalControllerRouter["howlong/controllers:UserController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
