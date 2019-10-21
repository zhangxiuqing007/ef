package controllers

import (
	"ef/models"
	"ef/usecase"
	"fmt"

	"github.com/astaxie/beego"
)

const sessionCookieId string = "sid"

var cookieLifeTime = beego.BConfig.WebConfig.Session.SessionCookieLifeTime

//Session 会话对象
type Session struct {
	sid          string
	PostSortType int //帖子排序方式
	User         *models.UserInDB
}

//创建登录信息
func (s *Session) buildLoginInfo() *loginInfo {
	result := new(loginInfo)
	if s == nil {
		return result
	}
	result.IsLogin = s.User != nil
	if result.IsLogin {
		result.UserID = s.User.ID
		result.UserName = s.User.Name
	}
	return result
}

//协助会话操作的基类控制器
type SessionBaseController struct {
	beego.Controller
}

//获取session对象
func (c *SessionBaseController) getSession() *Session {
	inter := c.GetSession(sessionCookieId)
	if inter == nil {
		inter = new(Session)
		c.SetSession(sessionCookieId, inter)
	}
	//空接口 转 会话对象指针
	s := inter.(*Session)
	//更新其他cookie，这里"beeGo"会莫名其妙更新sid的cookie
	if sidStr := c.Ctx.GetCookie(sessionCookieId); s.sid != sidStr {
		s.sid = sidStr
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

func (c *SessionBaseController) setNavigationVm(n *pageNavigations) {
	c.Data["navigationInfo"] = n
}

//设置登录Vm信息
func (c *SessionBaseController) setLoginVm(s *Session) {
	c.Data["loginInfo"] = s.buildLoginInfo()
}
func (c *SessionBaseController) setLoginVmSelf() {
	c.setLoginVm(c.getSession())
}

//发送首页
func (c *SessionBaseController) sendIndexPage() {
	if tms, err := usecase.QueryAllThemes(); err != nil || len(tms) == 0 {
		c.send404()
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
func (c *SessionBaseController) sendThemePage(data *themeFormData) {
	tm, err := usecase.QueryTheme(data.ThemeID)
	if err != nil {
		c.send404() //无法找到主题对象
		return
	}
	vm := new(themeVm)
	vm.ThemeID = tm.ID
	vm.WebTitle = "边缘社区-" + tm.Name
	s := c.getSession()
	pageIndex := limitPageIndex(data.PageIndex, postCountOnePage, tm.PostCount)
	vm.PostHeaders, err = usecase.QueryPostsOfTheme(tm.ID, postCountOnePage, pageIndex*postCountOnePage, s.PostSortType)
	if err != nil {
		c.send404() //无法找到帖子标题
		return
	}
	for _, v := range vm.PostHeaders {
		v.FormatShowInfo()
	}
	pathBuilder := func(i int) string {
		return fmt.Sprintf("/theme?ThemeID=%d&PageIndex=%d", tm.ID, i)
	}
	beginIndex, endIndex := getNavigationPageLimitIndex(pageIndex, postCountOnePage, halfPageCountToNavigationOfTheme, tm.PostCount)
	nevis := buildPageNavigations(pathBuilder, beginIndex, pageIndex, endIndex)
	c.setNavigationVm(nevis)
	c.setLoginVm(s)
	c.Data["vm"] = vm
	c.TplName = "theme_get.html"
}

type postVm struct {
	*models.PostOnPostPage
	Comments []*models.CmtOnPostPage
}

//发送帖子页
func (c *SessionBaseController) sendPostPage(data *postFormData) {
	vm := new(postVm)
	var err error
	vm.PostOnPostPage, err = usecase.QueryPostOfPostPage(data.PostID)
	if err != nil {
		c.send404()
		return
	}
	//发起请求的用户ID
	userID := 0
	s := c.getSession()
	if s.User != nil {
		userID = s.User.ID
	}
	//查询评论内容
	vm.PostOnPostPage.FormatShowInfo(userID)
	//限制页Index
	pageIndex := limitPageIndex(data.PageIndex, cmtCountOnePage, vm.CmtCount)
	vm.Comments, err = usecase.QueryCommentsOfPostPage(data.PostID, cmtCountOnePage, pageIndex*cmtCountOnePage, userID)
	if err != nil {
		c.send404()
		return
	}
	//生成文字的日期、赞和踩的checkBox属性和评论所在的楼层
	baseLayerCount := pageIndex * cmtCountOnePage
	for i, v := range vm.Comments {
		v.FormatCheckedStrOfPB()
		v.FormatStringTime()
		v.FormatIndex(baseLayerCount + i)
		v.FormatAllowEdit(userID)
		v.FormatCmtPageIndex(pageIndex)
	}
	//制作导航链接
	beginIndex, endIndex := getNavigationPageLimitIndex(pageIndex, cmtCountOnePage, halfPageCountToNavigationOfPost, vm.CmtCount)
	pathBuilder := func(index int) string {
		return fmt.Sprintf("/post?PostID=%d&PageIndex=%d", data.PostID, index)
	}
	c.setNavigationVm(buildPageNavigations(pathBuilder, beginIndex, pageIndex, endIndex))
	c.Data["vm"] = vm
	c.setLoginVm(s)
	c.TplName = "post_get.html"
}

//400 Bad Request
//该状态码表示请求报文中存在语法错误。当错误发生时，需要修改请求的内容再次发送请求。另外，浏览器会像200OK一样处理该状态码。
func (c *SessionBaseController) send400() {
	c.Ctx.ResponseWriter.WriteHeader(400)
	c.Ctx.WriteString("400：请求报文中存在语法错误")
}

//404 Not Found
//表明服务器上无法找到请求的资源。除此之外，也可以在服务端拒绝请求且不想说明理由时使用
func (c *SessionBaseController) send404() {
	c.Ctx.ResponseWriter.WriteHeader(404)
	c.Ctx.WriteString("404：服务器上无法找到请求的资源")
}