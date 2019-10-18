package controllers

import (
	"github.com/astaxie/beego"
)

//MainController 主控制器
type MainController struct {
	beego.Controller
}

//Get Get方法路由
func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
