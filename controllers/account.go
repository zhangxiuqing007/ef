package controllers

import (
	"ef/usecase"

	"github.com/astaxie/beego"
)

//IndexController 首页控制器
type AccountController struct {
	beego.Controller
}

//Get Get方法路由
func (c *AccountController) Get() {
	c.TplName = "account_get.html"
	c.Data["Tip"] = "请输入注册资料"
}

type newAccountFormData struct {
	Name      string `form:"name"`
	Account   string `form:"account"`
	Password1 string `form:"password1"`
	Password2 string `form:"password2"`
}

func (data *newAccountFormData) isTwoPwdSame() bool {
	return data.Password1 == data.Password2
}

func (data *newAccountFormData) isContentHasError() bool {
	return len(data.Name) == 0 || len(data.Account) == 0 || len(data.Password1) == 0 || len(data.Password2) == 0
}

//Get Get方法路由
func (c *AccountController) Post() {
	reInput := func(tip string) {
		c.TplName = "account_get.html"
		c.Data["Tip"] = tip
	}
	u := new(newAccountFormData)
	if err := c.ParseForm(u); err != nil {
		reInput("请输入完整的注册资料")
		return
	}
	if u.isContentHasError() {
		reInput("请输入完整的注册资料")
		return
	}
	if !u.isTwoPwdSame() {
		reInput("两次密码输入不一致")
		return
	}
	//组织申请数据
	data := &usecase.UserSignUpData{
		Name:     u.Name,
		Account:  u.Account,
		Password: u.Password1,
	}
	//调用用例层代码，尝试添加账户，并返回错误
	if err := usecase.AddUser(data); err != nil {
		reInput("注册失败：" + err.Error())
		return
	}
	//注册成功
	c.TplName = "account_post.html"
	c.Data["Name"] = u.Name
}
