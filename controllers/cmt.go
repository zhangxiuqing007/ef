package controllers

import (
	"ef/usecase"
	"errors"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
)

var cmtEditPageTemplate = template.Must(template.ParseFiles("views/cmtEdit.html", "views/login.html"))

type cmtEditVM struct {
	*loginInfo
	OriCmtContent string
	CmtID         int
	CmtPageIndex  int
}

func readFormDataFromCmtPb(r *http.Request) (cmtID int, isP bool, isD bool, err error) {
	if r.ParseForm() != nil {
		err = errors.New("分析提交内容失败")
		return
	}
	strs := r.Form["cmtID"]
	if strs != nil && len(strs) > 0 {
		cmtID, err = strconv.Atoi(strs[0])
		if err != nil {
			err = errors.New("获取评论目标失败")
		}
	} else {
		err = errors.New("缺失评论目标")
	}
	strs = r.Form["type"]
	if strs != nil && len(strs) > 0 {
		isP = strs[0] == "p"
	} else {
		err = errors.New("缺失赞踩类别")
	}
	strs = r.Form["dc"]
	if strs != nil && len(strs) > 0 {
		isD = strs[0] == "d"
	} else {
		err = errors.New("缺失赞踩操作类型")
	}
	return
}

//CmtPb 对评论进行赞或踩
func CmtPb(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s := getExsitOrCreateNewSession(w, r, true)
	if s.User == nil {
		w.Write([]byte("请先登录"))
		return
	}
	//解析表单
	cmtID, isP, isD, err := readFormDataFromCmtPb(r)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	//处理，数据库输入内容
	err = usecase.SetPB(cmtID, s.User.ID, isP, isD)
	if err != nil {
		w.Write([]byte("赞踩请求失败"))
		return
	}
	//返回
	w.Write([]byte("赞踩成功"))
}

func readFormDataFromCmt(r *http.Request) (postID int, cmtContent string) {
	checkErr(r.ParseForm())
	postID, err := strconv.Atoi(r.Form["postID"][0])
	checkErr(err)
	cmtContent = r.Form["cmtContent"][0]
	return
}

//Cmt 评论
func Cmt(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer recoverErrAndSendErrorPage(w)
	s := getExsitOrCreateNewSession(w, r, true)
	if s.User == nil {
		panic(errors.New("请先登录"))
	}
	postID, cmtContent := readFormDataFromCmt(r)
	if utf8.RuneCountInString(cmtContent) < 2 {
		panic(errors.New("评论字符最少需要2个字"))
	}
	err := usecase.AddComment(&usecase.CmtAddData{
		PostID:  postID,
		UserID:  s.User.ID,
		Content: cmtContent})
	checkErr(err)
	sendPostPage(w, postID, math.MaxInt32, s)
}

//CmtEdit 编辑评论
func CmtEdit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer recoverErrAndSendErrorPage(w)
	cmtID, err := strconv.Atoi(ps.ByName("cmtID"))
	checkErr(err)
	cmtPageIndex, err := strconv.Atoi(ps.ByName("cmtPageIndex"))
	checkErr(err)
	//先获取评论内容
	cmt, err := usecase.QueryComment(cmtID)
	checkErr(err)
	//无需检测权限，在编辑好之后，提交的时候检测即可
	//发送编辑页
	vm := new(cmtEditVM)
	vm.loginInfo = buildLoginInfo(getExsitOrCreateNewSession(w, r, true))
	vm.CmtID = cmt.ID
	vm.OriCmtContent = cmt.Content
	vm.CmtPageIndex = cmtPageIndex
	cmtEditPageTemplate.Execute(w, vm)
}

func readFormDataFromCmtEdit(r *http.Request) (cmtID int, CmtPageIndex int, cmtContent string) {
	checkErr(r.ParseForm())
	cmtID, err := strconv.Atoi(r.Form["cmtID"][0])
	checkErr(err)
	CmtPageIndex, err = strconv.Atoi(r.Form["CmtPageIndex"][0])
	checkErr(err)
	cmtContent = r.Form["cmtContent"][0]
	return
}

//CmtEditCommit 编辑评论提交
func CmtEditCommit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer recoverErrAndSendErrorPage(w)
	s := getExsitOrCreateNewSession(w, r, true)
	if s.User == nil {
		panic(errors.New("请先登录"))
	}
	cmtID, CmtPageIndex, cmtNewContent := readFormDataFromCmtEdit(r)
	//先获取评论内容
	cmt, err := usecase.QueryComment(cmtID)
	checkErr(err)
	//查看是否有编辑权限
	if cmt.UserID != s.User.ID {
		panic(errors.New("必须要发表者身份才能编辑评论"))
	}
	/*其他编辑权限，包括帖子状态和用户状态*/
	//修改评论内容
	cmt.Content = cmtNewContent
	cmt.LastEditTime = time.Now().UnixNano()
	cmt.EditTimes++
	checkErr(usecase.UpdateComment(cmt))
	//发送返回页
	sendPostPage(w, cmt.PostID, CmtPageIndex, s)
}
