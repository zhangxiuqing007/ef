package controllers

import (
	"ef/models"
	"ef/tool"
	"ef/usecase"
	"fmt"

	"github.com/astaxie/beego" //整个controller层都是依赖于beego的
)

const sessionCookieId string = "sid"

var cookieLifeTime = beego.BConfig.WebConfig.Session.SessionCookieLifeTime

//登录信息vm
type loginVm struct {
	IsLogin  bool
	UserID   int
	UserName string
}

//Session 会话对象
type session struct {
	Sid          string
	PostSortType int //帖子排序方式
	UserID       int
	UserName     string
}

//创建登录信息
func (s *session) buildLoginInfo() *loginVm {
	return &loginVm{
		IsLogin:  s.UserID > 0,
		UserID:   s.UserID,
		UserName: s.UserName,
	}
}

//基类控制器
type baseController struct {
	beego.Controller
}

//获取session对象
func (c *baseController) getSession() *session {
	inter := c.GetSession(sessionCookieId)
	if inter == nil {
		inter = new(session)
		c.SetSession(sessionCookieId, inter)
	}
	//空接口 转 会话对象指针
	s := inter.(*session)
	//更新其他cookie，这里"beeGo"会莫名其妙更新sid的cookie
	if sidStr := c.Ctx.GetCookie(sessionCookieId); s.Sid != sidStr {
		s.Sid = sidStr
		c.Ctx.SetCookie(sessionCookieId, sidStr, cookieLifeTime, "/account", "", false, true)
		c.Ctx.SetCookie(sessionCookieId, sidStr, cookieLifeTime, "/session", "", false, true)
		c.Ctx.SetCookie(sessionCookieId, sidStr, cookieLifeTime, "/theme", "", false, true)
		c.Ctx.SetCookie(sessionCookieId, sidStr, cookieLifeTime, "/user", "", false, true)
		c.Ctx.SetCookie(sessionCookieId, sidStr, cookieLifeTime, "/userPosts", "", false, true)
		c.Ctx.SetCookie(sessionCookieId, sidStr, cookieLifeTime, "/post", "", false, true)
		c.Ctx.SetCookie(sessionCookieId, sidStr, cookieLifeTime, "/newPost", "", false, true)
		c.Ctx.SetCookie(sessionCookieId, sidStr, cookieLifeTime, "/cmt", "", false, true)
		c.Ctx.SetCookie(sessionCookieId, sidStr, cookieLifeTime, "/attitude", "", false, true)
	}
	return s
}

func (c *baseController) setNavigationVm(n *tool.PageNavigations) {
	c.Data["navigationInfo"] = n
}

//自我获取并设置登录信息
func (c *baseController) setLoginVmSelf() {
	c.setLoginVm(c.getSession())
}

//设置登录Vm信息
func (c *baseController) setLoginVm(s *session) {
	c.Data["loginInfo"] = s.buildLoginInfo()
}

//发送首页
func (c *baseController) sendIndexPage() {
	if tms, err := usecase.QueryAllThemes(); err != nil || len(tms) == 0 {
		c.send404("无主题")
	} else {
		c.setLoginVmSelf()
		c.Data["Themes"] = tms
		c.TplName = "index_get.html"
	}
}

type themeVm struct {
	ThemeID     int
	WebTitle    string
	PostHeaders []*models.PostOnThemePage
}

//发送主题页
func (c *baseController) sendThemePage(data *themeFormData) {
	tm, err := usecase.QueryTheme(data.ThemeID)
	if err != nil {
		c.send404("无法找到主题")
		return
	}
	vm := new(themeVm)
	vm.ThemeID = tm.ID
	vm.WebTitle = "边缘社区-" + tm.Name
	session := c.getSession()
	oper := new(tool.PageNavigationOperator)
	pageIndex := oper.LimitPageIndex(data.PageIndex, postCountOnePage, tm.PostCount)
	vm.PostHeaders, err = usecase.QueryPostsOfTheme(tm.ID, postCountOnePage, pageIndex*postCountOnePage, session.PostSortType)
	if err != nil {
		c.send404("找不到主题内的帖子列表")
		return
	}
	for _, v := range vm.PostHeaders {
		v.FormatShowInfo()
	}
	pathBuilder := func(i int) string {
		return fmt.Sprintf("/theme?ThemeID=%d&PageIndex=%d", tm.ID, i)
	}
	beginIndex, endIndex := oper.GetNavigationPageLimitIndex(pageIndex, postCountOnePage, halfPageCountToNavigationOfTheme, tm.PostCount)
	nevis := oper.BuildPageNavigations(pathBuilder, beginIndex, pageIndex, endIndex)
	c.setNavigationVm(nevis)
	c.setLoginVm(session)
	c.Data["vm"] = vm
	c.TplName = "theme_get.html"
}

//帖子页vm
type postVm struct {
	*models.PostOnPostPage
	Comments []*models.CmtOnPostPage
}

//发送帖子页
func (c *baseController) sendPostPage(data *postFormData) {
	vm := new(postVm)
	var err error
	vm.PostOnPostPage, err = usecase.QueryPostOfPostPage(data.PostID)
	if err != nil {
		c.send404("找不到标题")
		return
	}
	oper := new(tool.PageNavigationOperator)
	s := c.getSession()
	//查询评论内容
	vm.PostOnPostPage.FormatShowInfo(s.UserID)
	//限制页Index
	pageIndex := oper.LimitPageIndex(data.PageIndex, cmtCountOnePage, vm.CmtCount)
	vm.Comments, err = usecase.QueryCommentsOfPostPage(data.PostID, cmtCountOnePage, pageIndex*cmtCountOnePage, s.UserID)
	if err != nil {
		c.send404("找不到评论")
		return
	}
	//生成文字的日期、赞和踩的checkBox属性和评论所在的楼层
	baseLayerCount := pageIndex * cmtCountOnePage
	for i, v := range vm.Comments {
		v.FormatCheckedStrOfPB()
		v.FormatStringTime()
		v.FormatIndex(baseLayerCount + i)
		v.FormatAllowEdit(s.UserID)
		v.FormatCmtPageIndex(pageIndex)
	}
	//制作导航链接
	beginIndex, endIndex := oper.GetNavigationPageLimitIndex(pageIndex, cmtCountOnePage, halfPageCountToNavigationOfPost, vm.CmtCount)
	pathBuilder := func(index int) string {
		return fmt.Sprintf("/post?PostID=%d&PageIndex=%d", data.PostID, index)
	}
	c.setNavigationVm(oper.BuildPageNavigations(pathBuilder, beginIndex, pageIndex, endIndex))
	c.Data["vm"] = vm
	c.setLoginVm(s)
	c.TplName = "post_get.html"
}

//200 	OK
//请求成功。一般用于GET与POST请求
func (c *baseController) send200(extraStr string) {
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.Ctx.WriteString("200 	OK	" + extraStr)
}

//400 Bad Request
//该状态码表示请求报文中存在语法错误。当错误发生时，需要修改请求的内容再次发送请求。另外，浏览器会像200OK一样处理该状态码。
func (c *baseController) send400(extraStr string) {
	c.Ctx.ResponseWriter.WriteHeader(400)
	c.Ctx.WriteString("400	Bad Request	" + extraStr)
}

//401 	Unauthorized
//请求要求用户的身份认证
func (c *baseController) send401(extraStr string) {
	c.Ctx.ResponseWriter.WriteHeader(401)
	c.Ctx.WriteString("401	Unauthorized	" + extraStr)
}

//403 	Forbidden
//服务器理解请求客户端的请求，但是拒绝执行此请求
func (c *baseController) send403(extraStr string) {
	c.Ctx.ResponseWriter.WriteHeader(403)
	c.Ctx.WriteString("403	Forbidden	" + extraStr)
}

//404 Not Found
//表明服务器上无法找到请求的资源。除此之外，也可以在服务端拒绝请求且不想说明理由时使用
func (c *baseController) send404(extraStr string) {
	c.Ctx.ResponseWriter.WriteHeader(404)
	c.Ctx.WriteString("404	Not Found	" + extraStr)
}

//406 	Not Acceptable
//服务器无法根据客户端请求的内容特性完成请求
func (c *baseController) send406(extraStr string) {
	c.Ctx.ResponseWriter.WriteHeader(406)
	c.Ctx.WriteString("406	Not Acceptable	" + extraStr)
}
