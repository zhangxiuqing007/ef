package controllers

type ThemeController struct {
	baseController
}

type themeFormData struct {
	ThemeID   int
	PageIndex int
}

func (c *ThemeController) Get() {
	data := new(themeFormData)
	if err := c.ParseForm(data); err != nil || data.ThemeID <= 0 || data.PageIndex < 0 {
		c.send400("请求信息错误") //表单信息错误
		return
	}
	c.sendThemePage(data)
}
