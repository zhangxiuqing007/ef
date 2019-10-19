package controllers

import (
	"ef/usecase"
)

type loginFormData struct {
	Account  string `form:"account"`
	Password string `form:"password"`
}

func (data *loginFormData) IsInputError() bool {
	return len(data.Account) == 0 || len(data.Password) == 0
}

type SessionController struct {
	SessionBaseController
}

//返回登录输入界面
func (c *SessionController) Get() {
	c.TplName = "session_get.html"
	c.Data["Tip"] = "请输入账号密码"
}

//登录
func (c *SessionController) Post() {
	reInput := func(tip string) {
		c.TplName = "session_get.html"
		c.Data["Tip"] = tip
	}
	//读取form信息
	data := new(loginFormData)
	if err := c.ParseForm(data); err != nil {
		reInput(err.Error())
	}
	if data.IsInputError() {
		reInput("请输入正确的账号密码")
	}
	//查询用户
	if user, err := usecase.QueryUserByAccountAndPwd(data.Account, data.Password); err != nil {
		reInput(err.Error())
	} else {
		s := c.getSession()
		s.User = user
		c.TplName = "session_post.html"
		c.Data["loginInfo"] = buildLoginInfo(s)
		c.Data["Name"] = user.Name
	}
}
