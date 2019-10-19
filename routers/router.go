package routers

import (
	"ef/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/session", &controllers.SessionController{})
	beego.Router("/account", &controllers.AccountController{})
}
