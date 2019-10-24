package controllers

type UserController struct {
	baseController
}

type userFromData struct {
	UserID int
}

func (c *UserController) Get() {
	data := new(userFromData)
	if err := c.ParseForm(data); err != nil || data.UserID <= 0 {
		c.send400("请求信息错误")
		return
	}
	c.sendUserPage(data)
}
