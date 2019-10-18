package controllers

import (
	"ef/models"
	"ef/usecase"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

var postTemplate = template.Must(template.ParseFiles("views/post.html", "views/pageNavi.html", "views/login.html"))
var postTitleEditTemplate = template.Must(template.ParseFiles("views/postTitleEdit.html", "views/login.html"))

const cmtCountOnePage = 20                //帖子页，一页评论的数量
const halfPageCountToNavigationOfPost = 8 //评论导航页数量

type postVM struct {
	*loginInfo
	*models.PostOnPostPage
	Comments []*models.CmtOnPostPage
	*pageNavis
}

//PostInfo 查看帖子
func PostInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	postStrID := ps.ByName("postID")
	pageStrIndex := ps.ByName("pageIndex")
	//转化成int型的postId
	postID, err := strconv.Atoi(postStrID)
	if err != nil {
		sendErrorPage(w, err.Error())
		return
	}
	pageIndex, err := strconv.Atoi(pageStrIndex)
	if err != nil || pageIndex < 0 {
		pageIndex = 0
	}
	sendPostPage(w, postID, pageIndex, getExsitOrCreateNewSession(w, r, true))
}

//发送帖子页
func sendPostPage(w http.ResponseWriter, postID int, pageIndex int, s *Session) {
	vm := new(postVM)
	vm.loginInfo = buildLoginInfo(s)
	//查询帖子主体信息
	var err error
	vm.PostOnPostPage, err = usecase.QueryPostOfPostPage(postID)
	if err != nil {
		sendErrorPage(w, "帖子查询失败")
		return
	}
	//限制页Index
	pageIndex = limitPageIndex(pageIndex, cmtCountOnePage, vm.CmtCount)
	//查询评论内容
	userID := 0
	if s.User != nil {
		userID = s.User.ID
	}
	vm.PostOnPostPage.FormatShowInfo(userID)
	vm.Comments, err = usecase.QueryCommentsOfPostPage(postID, cmtCountOnePage, pageIndex*cmtCountOnePage, userID)
	if err != nil {
		sendErrorPage(w, "评论查询失败")
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
	pathBuilder := func(index int) string {
		return fmt.Sprintf("/Post/Content/%d/%d", postID, index)
	}
	beginIndex, endIndex := getNaviPageIndexs(pageIndex, cmtCountOnePage, halfPageCountToNavigationOfPost, vm.CmtCount)
	vm.pageNavis = buildPageNavis(pathBuilder, beginIndex, pageIndex, endIndex)
	//发送帖子页
	postTemplate.Execute(w, vm)
}

type postTitleEditVM struct {
	*loginInfo
	PostID   int
	OriTitle string
}

//PostTitleEdit 编辑标题
func PostTitleEdit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	postID, err := strconv.Atoi(ps.ByName("postID"))
	if err != nil {
		sendErrorPage(w, "编辑标题请求失败")
		return
	}
	s := getExsitOrCreateNewSession(w, r, true)
	vm := new(postTitleEditVM)
	vm.loginInfo = buildLoginInfo(s)
	//查询帖子现有标题
	title, err := usecase.QueryPostTitle(postID)
	if err != nil {
		sendErrorPage(w, "编辑标题请求失败")
		return
	}
	vm.PostID = postID
	vm.OriTitle = title
	postTitleEditTemplate.Execute(w, vm)
}

func readFormDataFromTitleEditCommit(r *http.Request) (postID int, newTitle string) {
	checkErr(r.ParseForm())
	postID, err := strconv.Atoi(r.Form["postID"][0])
	checkErr(err)
	newTitle = r.Form["title"][0]
	return
}

//PostTitleEditCommit 提交新标题
func PostTitleEditCommit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer recoverErrAndSendErrorPage(w)
	s := getExsitOrCreateNewSession(w, r, true)
	if s.User == nil {
		panic(errors.New("请先登录"))
	}
	postID, newTitle := readFormDataFromTitleEditCommit(r)
	//查询旧的Post
	post, err := usecase.QueryPost(postID)
	checkErr(err)
	//验证合法性
	if s.User.ID != post.UserID {
		panic(errors.New("必须要发表者身份才能编辑标题"))
	}
	/*验证其他状态，包括帖子状态，标题是否可以修改，用户是否还有编辑权限*/
	//更新内容
	post.Title = newTitle
	post.LastCmtTime = time.Now().UnixNano()
	//保存到DB
	checkErr(usecase.UpdatePostTitle(post))
	//发送帖子页
	sendPostPage(w, postID, 0, s)
}
