package controllers

import (
	"ef/models"
	"ef/tool"
	"ef/usecase"
	"time"
)

type userVm struct {
	*models.UserInDB
	//下面需要格式化
	SignUpTime   string
	Type         string
	State        string
	LastEditTime string
}

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
	saInfo, err := usecase.QueryUserByID(data.UserID)
	if err != nil {
		c.send404("用户不存在")
		return
	}
	vm := new(userVm)
	vm.UserInDB = saInfo
	vm.SignUpTime = tool.FormatTimeDetail(time.Unix(0, saInfo.SignUpTime))
	vm.Type = models.GetUserTypeShowName(saInfo.Type)
	vm.State = models.GetUserStateShowName(saInfo.State)
	if saInfo.LastEditTime == 0 {
		vm.LastEditTime = "无"
	} else {
		vm.LastEditTime = tool.FormatTimeDetail(time.Unix(0, saInfo.LastEditTime))
	}
	c.Data["vm"] = vm
	c.TplName = "user_get.html"
}
