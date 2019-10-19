package controllers

import (
	"ef/usecase"
)

//IndexController 首页控制器
type IndexController struct {
	SessionBaseController
}

//Get Get方法路由
func (c *IndexController) Get() {
	//查询主题
	tms, err := usecase.QueryAllThemes()
	if err != nil {
		c.Data["Tip"] = "查询主题列表失败"
		c.TplName = "error.html"
	} else if len(tms) == 0 {
		c.Data["Tip"] = "无主题"
		c.TplName = "error.html"
	} else {
		c.Data["loginInfo"] = buildLoginInfo(c.getSession())
		c.Data["Themes"] = tms
		c.TplName = "index.html"
	}
}
