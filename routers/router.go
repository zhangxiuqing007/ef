package routers

import (
	"ef/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/account", &controllers.AccountController{})
	beego.Router("/session", &controllers.SessionController{})
	beego.Router("/theme", &controllers.ThemeController{})
	beego.Router("/user", &controllers.UserController{})
	beego.Router("/userPosts", &controllers.UserPostsController{})
	beego.Router("/post", &controllers.PostController{})
	beego.Router("/newPost", &controllers.NewPostController{})
	beego.Router("/cmt", &controllers.CmtController{})
	beego.Router("/attitude", &controllers.AttitudeController{})
}
