package controllers

import (
	"ef/usecase"
	"math"
	"time"
	"unicode/utf8"
)

type CmtController struct {
	baseController
}

type cmtEditFormData struct {
	CmtID        int
	CmtPageIndex int
}

func (c *CmtController) Get() {
	data := new(cmtEditFormData)
	if err := c.ParseForm(data); err != nil {
		c.send400("请求信息错误")
		return
	}
	//id==0，是要新增评论
	if data.CmtID <= 0 {
		c.GetCmtNewInputPage(data)
	} else {
		c.GetCmtFixPage(data)
	}
}

type cmtNewInputVm struct {
	PostID    int
	PostTitle string
}

//请求评论新建页面
func (c *CmtController) GetCmtNewInputPage(data *cmtEditFormData) {
	s := c.getSession()
	if s.UserID <= 0 {
		c.send401("请先登录")
		return
	}
	vm := new(cmtNewInputVm)
	vm.PostID = data.CmtPageIndex
	vm.PostTitle = usecase.QueryPostTitle(data.CmtPageIndex)
	c.setLoginVm(s)
	c.Data["vm"] = vm
	c.TplName = "cmt_get_new.html"
}

type cmtEditVm struct {
	OriCmtContent string
	CmtID         int
	CmtPageIndex  int
}

//请求评论修改页面
func (c *CmtController) GetCmtFixPage(data *cmtEditFormData) {
	cmt := usecase.QueryComment(data.CmtID)
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
		c.send401("请先登录")
		return
	}
	data := new(newCmtFormData)
	if err := c.ParseForm(data); err != nil {
		c.send400("请求信息错误")
		return
	}
	//检查输入内容合法性，用户权限
	if utf8.RuneCountInString(data.CmtContent) < 2 {
		c.send403("评论字符最少需要2个字")
		return
	}
	//添加评论
	usecase.AddComment(&usecase.CmtAddData{
		PostID:  data.PostID,
		UserID:  s.UserID,
		Content: data.CmtContent})
	//发送帖子的最后一个评论页
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
		c.send401("请先登录")
		return
	}
	data := new(cmtEditCommitFormData)
	if err := c.ParseForm(data); err != nil {
		c.send400("请求信息错误")
		return
	}
	//先从当前db获取评论内容
	cmt := usecase.QueryComment(data.CmtID)
	//查看是否有编辑权限
	if cmt.UserID != s.UserID {
		c.send401("您无修改权限，必须要发表者身份才能编辑评论")
		return
	}
	/*其他编辑权限，包括帖子状态和用户状态*/
	//修改评论内容
	cmt.Content = data.CmtContent
	cmt.LastEditTime = time.Now().UnixNano()
	cmt.EditTimes++
	usecase.UpdateComment(cmt)
	//发送返回页
	c.sendPostPage(&postFormData{
		PostID:    cmt.PostID,
		PageIndex: data.CmtPageIndex,
	})
}
