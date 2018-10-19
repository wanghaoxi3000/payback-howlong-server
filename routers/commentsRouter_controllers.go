package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["howlong/controllers:AuthController"] = append(beego.GlobalControllerRouter["howlong/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:AuthController"] = append(beego.GlobalControllerRouter["howlong/controllers:AuthController"],
		beego.ControllerComments{
			Method: "RefreshToken",
			Router: `/refresh-token`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:CreditController"] = append(beego.GlobalControllerRouter["howlong/controllers:CreditController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["howlong/controllers:CreditController"] = append(beego.GlobalControllerRouter["howlong/controllers:CreditController"],
		beego.ControllerComments{
			Method: "List",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
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

}
