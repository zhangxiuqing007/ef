package controllers

import (
	"ef/usecase"
	"math"
	"time"
	"unicode/utf8"
)

type CmtController struct {
	SessionBaseController
}

func (c *CmtController) Get() {
	c.GetCmtFixPage()
}

//请求评论新建页面
func (c *CmtController) GetCmtNewInputPage() {

}

type cmtEditFormData struct {
	CmtID        int
	CmtPageIndex int
}

type cmtEditVm struct {
	OriCmtContent string
	CmtID         int
	CmtPageIndex  int
}

//请求评论修改页面
func (c *CmtController) GetCmtFixPage() {
	data := new(cmtEditFormData)
	if err := c.ParseForm(data); err != nil {
		c.send400()
		return
	}
	cmt, err := usecase.QueryComment(data.CmtID)
	if err != nil {
		c.send404()
		return
	}
	//无需检测权限，在编辑好之后，提交的时候检测即可
	//发送编辑页
	vm := new(cmtEditVm)
	vm.CmtID = data.CmtID
	vm.OriCmtContent = cmt.Content
	vm.CmtPageIndex = data.CmtPageIndex
	c.Data["vm"] = vm
	c.setLoginVmSelf()
	c.TplName = "cmt_get_edit.html"
}

type newCmtFormData struct {
	PostID     int
	CmtContent string
}

func (c *CmtController) Post() {
	s := c.getSession()
	if s.UserID == 0 {
		c.send401()
		return
	}
	data := new(newCmtFormData)
	if err := c.ParseForm(data); err != nil {
		c.send400()
		return
	}
	//检查输入内容合法性，用户权限
	if utf8.RuneCountInString(data.CmtContent) < 2 {
		c.send403()
		c.Ctx.WriteString("评论字符最少需要2个字")
		return
	}
	if err := usecase.AddComment(&usecase.CmtAddData{
		PostID:  data.PostID,
		UserID:  s.UserID,
		Content: data.CmtContent}); err != nil {
		c.send404()
		c.Ctx.WriteString(err.Error())
		return
	}
	c.sendPostPage(&postFormData{
		PostID:    data.PostID,
		PageIndex: math.MaxInt32,
	})
}

type cmtEditCommitFormData struct {
	CmtID        int
	CmtPageIndex int
	CmtContent   string
}

func (c *CmtController) Put() {
	s := c.getSession()
	if s.UserID == 0 {
		c.send401()
		return
	}
	data := new(cmtEditCommitFormData)
	if err := c.ParseForm(data); err != nil {
		c.send400()
		return
	}
	//先从当前db获取评论内容
	cmt, err := usecase.QueryComment(data.CmtID)
	if err != nil {
		c.send404()
		return
	}
	//查看是否有编辑权限
	if cmt.UserID != s.UserID {
		c.send401()
		c.Ctx.WriteString("必须要发表者身份才能编辑评论")
		return
	}
	/*其他编辑权限，包括帖子状态和用户状态*/
	//修改评论内容
	cmt.Content = data.CmtContent
	cmt.LastEditTime = time.Now().UnixNano()
	cmt.EditTimes++
	if err := usecase.UpdateComment(cmt); err != nil {
		c.send404()
		c.Ctx.WriteString(err.Error())
		return
	}
	//发送返回页
	c.sendPostPage(&postFormData{
		PostID:    cmt.PostID,
		PageIndex: data.CmtPageIndex,
	})
}
