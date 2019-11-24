package controllers

import (
	"ef/usecase"
)

type loginFormData struct {
	Account  string
	Password string
}

func (data *loginFormData) IsInputError() bool {
	return len(data.Account) == 0 || len(data.Password) == 0
}

type SessionController struct {
	baseController
}

//发送输入页
func (c *SessionController) sendInputPage(tip string) {
	c.TplName = "session_get.html"
	c.Data["Tip"] = tip
}

//返回登录输入界面
func (c *SessionController) Get() {
	c.sendInputPage("请输入账号密码")
}

//登录
func (c *SessionController) Post() {
	//读取form信息
	data := new(loginFormData)
	if err := c.ParseForm(data); err != nil {
		c.sendInputPage(err.Error())
		return
	}
	if data.IsInputError() {
		c.sendInputPage("请输入正确的账号密码")
		return
	}
	//查询用户
	if user := usecase.QueryUserByAccountAndPwd(data.Account, data.Password); user == nil {
		c.sendInputPage("账号密码错误")
	} else {
		s := c.getSession()
		s.UserID = user.ID
		s.UserName = user.Name
		c.TplName = "session_post.html"
		c.Data["Name"] = user.Name
	}
}

//登出
func (c *SessionController) Delete() {
	s := c.getSession()
	//清空登录记录
	s.UserID = 0
	s.UserName = ""
	//发送首页
	c.sendIndexPage()
}
