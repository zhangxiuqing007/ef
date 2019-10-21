package controllers

//IndexController 首页控制器
type IndexController struct {
	SessionBaseController
}

//Get Get方法路由
func (c *IndexController) Get() {
	c.sendIndexPage()
}
