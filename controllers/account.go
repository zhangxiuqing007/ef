package controllers

import (
	"ef/usecase"

	"github.com/astaxie/beego"
)

//AccountController 账号密码控制器
type AccountController struct {
	beego.Controller
}

func (c *AccountController) sendInputPage(tip string) {
	c.TplName = "account_get.html"
	c.Data["Tip"] = tip
}

//Get Get方法路由
func (c *AccountController) Get() {
	c.sendInputPage("请输入注册资料")
}

type newAccountFormData struct {
	Name      string
	Account   string
	Password1 string
	Password2 string
}

func (data *newAccountFormData) isTwoPwdSame() bool {
	return data.Password1 == data.Password2
}

func (data *newAccountFormData) isContentHasError() bool {
	return len(data.Name) == 0 || len(data.Account) == 0 || len(data.Password1) == 0 || len(data.Password2) == 0
}

//Get Get方法路由
func (c *AccountController) Post() {
	data := new(newAccountFormData)
	if err := c.ParseForm(data); err != nil {
		c.sendInputPage("请输入完整的注册资料")
		return
	}
	if data.isContentHasError() {
		c.sendInputPage("请输入完整的注册资料")
		return
	}
	if !data.isTwoPwdSame() {
		c.sendInputPage("两次密码输入不一致")
		return
	}
	//组织申请数据
	addUser := &usecase.UserSignUpData{
		Name:     data.Name,
		Account:  data.Account,
		Password: data.Password1,
	}
	//调用用例层代码，尝试添加账户，如果失败，则直接崩溃
	usecase.AddUser(addUser)
	//注册成功，发送结果
	c.TplName = "account_post.html"
	c.Data["Name"] = data.Name
}
