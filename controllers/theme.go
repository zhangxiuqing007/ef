package controllers

const postCountOnePage = 20                //主题页，一页帖子数量
const halfPageCountToNavigationOfTheme = 8 //帖子导航页数量

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
