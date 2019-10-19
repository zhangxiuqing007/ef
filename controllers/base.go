package controllers

import (
	"ef/models"
	"fmt"
	"net/http"
	"time"

	"github.com/astaxie/beego"
)

const sessionCookieId string = "sid"

var cookieLifeTime = beego.BConfig.WebConfig.Session.SessionCookieLifeTime

//登录信息vm
type loginInfo struct {
	IsLogin  bool
	UserID   int
	UserName string
}

//Session 会话对象
type Session struct {
	CreatedTime     int64
	LastRequestTime int64
	PostSortType    int //帖子排序方式
	User            *models.UserInDB
}

//协助会话操作的基类控制器
type SessionBaseController struct {
	beego.Controller
}

func (c *SessionBaseController) getSession() *Session {
	i := c.GetSession(sessionCookieId)
	if i == nil {
		i = c.createNewSession()
		c.SetSession(sessionCookieId, i)
		//获取其sid
		sidStr := c.Ctx.GetCookie("sid")
		c.Ctx.SetCookie("sid", sidStr, cookieLifeTime, "/session", "", false, true)
		c.Ctx.SetCookie("sid", sidStr, cookieLifeTime, "/theme", "", false, true)
		c.Ctx.SetCookie("sid", sidStr, cookieLifeTime, "/user", "", false, true)
		c.Ctx.SetCookie("sid", sidStr, cookieLifeTime, "/userPosts", "", false, true)
		c.Ctx.SetCookie("sid", sidStr, cookieLifeTime, "/post", "", false, true)
		c.Ctx.SetCookie("sid", sidStr, cookieLifeTime, "/newPost", "", false, true)
		c.Ctx.SetCookie("sid", sidStr, cookieLifeTime, "/cmt", "", false, true)
		c.Ctx.SetCookie("sid", sidStr, cookieLifeTime, "/attitude", "", false, true)
	}
	//空接口 转 会话对象指针
	s := i.(*Session)
	s.LastRequestTime = time.Now().UnixNano()
	return s
}

func (c *SessionBaseController) createNewSession() *Session {
	session := new(Session)
	session.CreatedTime = time.Now().UnixNano()
	session.LastRequestTime = session.CreatedTime
	session.User = nil
	return session
}

//创建登录信息
func buildLoginInfo(s *Session) *loginInfo {
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

type pageNavi struct {
	Path   string
	Number int
}

type pageNavis struct {
	HeadPageNavis []*pageNavi //前导航页
	CurrentPage   *pageNavi   //当前页
	TailPageNavis []*pageNavi //后导航页
}

//创建导航页vm
func buildPageNavis(pathBuilder func(index int) string, beginIndex, currentIndex, endIndex int) *pageNavis {
	//制作导航页切片
	buildNavis := func(bi, ei int) []*pageNavi {
		slice := make([]*pageNavi, 0, 16)
		for i := bi; i <= ei; i++ {
			slice = append(slice, &pageNavi{pathBuilder(i), i + 1})
		}
		return slice
	}
	//确定导航页限制
	p := new(pageNavis)
	p.HeadPageNavis = buildNavis(beginIndex, currentIndex-1)
	p.CurrentPage = &pageNavi{"", currentIndex + 1}
	p.TailPageNavis = buildNavis(currentIndex+1, endIndex)
	return p
}

//获取提供的导航页
func getNaviPageIndexs(
	currentPageIndex int, /*当前页索引*/
	countOnePage int, /*一页元素数量*/
	maxHalfNaviPageCount int, /*最大的导航页数量的一半*/
	elementTotalCount int) /*元素的总数量*/ (beginIndex, endIndex int) {
	//先计算beginIndex
	beginIndex = currentPageIndex - maxHalfNaviPageCount
	if beginIndex < 0 {
		beginIndex = 0
	}
	endIndex = limitPageIndex(currentPageIndex+maxHalfNaviPageCount, countOnePage, elementTotalCount)
	return
}

//限制页索引
func limitPageIndex(currentIndex int, countOnePage int, totalCount int) int {
	maxIndex := totalCount / countOnePage
	if currentIndex > maxIndex {
		currentIndex = maxIndex
	}
	if currentIndex < 0 {
		currentIndex = 0
	}
	return currentIndex
}

//请求的统一捕获异常处理函数
func recoverErrAndSendErrorPage(w http.ResponseWriter) {
	if err := recover(); err != nil {
		sendErrorPage(w, fmt.Sprintf("操作失败 %v", err))
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getExsitOrCreateNewSession(w http.ResponseWriter, r *http.Request, recordTime bool) *Session {
	return nil
}
