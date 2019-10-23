package controllers

import "ef/usecase"

type AttitudeController struct {
	SessionBaseController
}

type PBFormData struct {
	CmtID int
	IsP   bool
	IsD   bool
}

func (c *AttitudeController) Post() {
	//解析表单数据
	data := new(PBFormData)
	if err := c.ParseForm(data); err != nil {
		c.send400()
		return
	}
	//先查看登录状态
	s := c.getSession()
	if s.UserID == 0 {
		c.send401()
		return
	}
	//无法完成请求的内容
	if err := usecase.SetPB(data.CmtID, s.UserID, data.IsP, data.IsD); err != nil {
		c.send406()
		return
	}
	//直接返回语句，代表成功
	c.send200()
}
