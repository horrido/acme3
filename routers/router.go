package routers

import (
	"acme3/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/home", &controllers.MainController{})
	beego.Router("/about", &controllers.MainController{}, "get:About")
	beego.Router("/photos", &controllers.MainController{}, "get:Photos")
	beego.Router("/user/login/:back", &controllers.MainController{}, "get,post:Login")
	beego.Router("/user/logout", &controllers.MainController{}, "get:Logout")
	beego.Router("/user/register", &controllers.MainController{}, "get,post:Register")
	beego.Router("/user/profile", &controllers.MainController{}, "get,post:Profile")
	beego.Router("/user/verify/:uuid({[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}})", &controllers.MainController{}, "get:Verify")
	beego.Router("/user/remove", &controllers.MainController{}, "get,post:Remove")
	beego.Router("/user/forgot", &controllers.MainController{}, "get,post:Forgot")
	beego.Router("/user/reset/:uuid({[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}})", &controllers.MainController{}, "get,post:Reset")
	beego.Router("/notice", &controllers.MainController{}, "get:Notice")
	beego.Router("/appadmin/index/:parms", &controllers.AdminController{}, "get,post:Index")
	beego.Router("/appadmin/add/:parms", &controllers.AdminController{}, "get,post:Add")
	beego.Router("/appadmin/update/:username", &controllers.AdminController{}, "get,post:Update")
}
