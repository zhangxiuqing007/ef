package controllers

import (
	"ef/models"
	"ef/usecase"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var themeTemplate = template.Must(template.ParseFiles("views/theme.html", "views/comp/pageNavi.html", "views/comp/login_info_head.html"))

const postCountOnePage = 20                //主题页，一页帖子数量
const halfPageCountToNavigationOfTheme = 8 //帖子导航页数量

type themeVM struct {
	ThemeID     int
	WebTitle    string                    //网页Header
	*loginInfo                            //登录信息
	PostHeaders []*models.PostOnThemePage //帖子简要内容
	*pageNavis
}

//Theme 请求主题内帖子列表
func Theme(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	themeID, err := strconv.Atoi(ps.ByName("themeID"))
	if err != nil {
		sendErrorPage(w, "尝试访问错误的主题ID")
		return
	}
	pageIndex, err := strconv.Atoi(ps.ByName("pageIndex"))
	if err != nil {
		sendErrorPage(w, "尝试访问错误的分页")
		return
	}
	sendThemePage(w, themeID, pageIndex, getExsitOrCreateNewSession(w, r, true))
}

//发送主题页，帖子列表
func sendThemePage(w http.ResponseWriter, themeID int, pageIndex int, s *Session) {
	tm, err := usecase.QueryTheme(themeID)
	if err != nil {
		sendErrorPage(w, "访问主题失败")
		return
	}
	//创建 viewModel对象
	vm := new(themeVM)
	vm.ThemeID = themeID
	//给vm赋基本值
	vm.WebTitle = "边缘社区-" + tm.Name
	vm.loginInfo = buildLoginInfo(s)
	//限制请求页Index
	totalPostCount, err := usecase.QueryPostCountOfTheme(tm.ID)
	if err != nil {
		sendErrorPage(w, "访问主题失败")
		return
	}
	pageIndex = limitPageIndex(pageIndex, postCountOnePage, totalPostCount)
	//获取帖子列表，根据请求的页码查询帖子列表
	vm.PostHeaders, err = usecase.QueryPostsOfTheme(tm.ID, postCountOnePage, pageIndex*postCountOnePage, s.PostSortType)
	if err != nil {
		sendErrorPage(w, "查询帖子列表失败")
		return
	}
	for _, v := range vm.PostHeaders {
		v.FormatShowInfo()
	}
	pathBuilder := func(i int) string {
		return fmt.Sprintf("/Theme/%d/%d", tm.ID, i)
	}
	beginIndex, endIndex := getNaviPageIndexs(pageIndex, postCountOnePage, halfPageCountToNavigationOfTheme, totalPostCount)
	vm.pageNavis = buildPageNavis(pathBuilder, beginIndex, pageIndex, endIndex)
	themeTemplate.Execute(w, vm)
}

//UserPosts 请求查看某个用户所发的所有的帖子
func UserPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		sendErrorPage(w, "尝试访问错误的用户")
		return
	}
	pageIndex, err := strconv.Atoi(ps.ByName("pageIndex"))
	if err != nil {
		sendErrorPage(w, "尝试访问错误的分页")
		return
	}
	sendUserPostPage(w, userID, pageIndex, getExsitOrCreateNewSession(w, r, true))
}

//发送用户发帖列表
func sendUserPostPage(w http.ResponseWriter, userID int, pageIndex int, s *Session) {
	vm := new(themeVM)
	vm.loginInfo = buildLoginInfo(s)
	vm.WebTitle = "边缘社区-用户发帖列表"
	//限制请求页Index
	totalPostCount, err := usecase.QueryPostCountOfUser(userID)
	if err != nil {
		sendErrorPage(w, "访问用户发帖列表失败")
		return
	}
	pageIndex = limitPageIndex(pageIndex, postCountOnePage, totalPostCount)
	//获取帖子列表，根据请求的页码查询帖子列表
	vm.PostHeaders, err = usecase.QueryPostsOfUser(userID, postCountOnePage, pageIndex*postCountOnePage)
	if err != nil {
		sendErrorPage(w, "查询帖子列表失败")
		return
	}
	for _, v := range vm.PostHeaders {
		v.FormatShowInfo()
	}
	pathBuilder := func(i int) string {
		return fmt.Sprintf("/User/%d/%d", userID, i)
	}
	beginIndex, endIndex := getNaviPageIndexs(pageIndex, postCountOnePage, halfPageCountToNavigationOfTheme, totalPostCount)
	vm.pageNavis = buildPageNavis(pathBuilder, beginIndex, pageIndex, endIndex)
	themeTemplate.Execute(w, vm)
}
