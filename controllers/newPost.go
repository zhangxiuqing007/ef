package controllers

import (
	"ef/usecase"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type postInputVM struct {
	*loginInfo
	ThemeName string
	ThemeID   int
}

var postInputModel = template.Must(template.ParseFiles("views/postInput.html", "views/comp/login_info_head.html"))

//NewPostInput 新增帖子，请求编辑表单,GET
func NewPostInput(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s := getExsitOrCreateNewSession(w, r, true)
	if s.User == nil {
		sendErrorPage(w, "请先登录")
		return
	}
	themeID, err := strconv.Atoi(ps.ByName("themeID"))
	if err != nil {
		sendErrorPage(w, "尝试在错误的主题发帖")
		return
	}
	tm, err := usecase.QueryTheme(themeID)
	if err != nil {
		sendErrorPage(w, "主题不存在")
		return
	}
	vm := new(postInputVM)
	vm.loginInfo = buildLoginInfo(s)
	vm.ThemeID = tm.ID
	vm.ThemeName = tm.Name
	postInputModel.Execute(w, vm)
}

func readFormDataOfNewPost(r *http.Request) (themeID int, title string, content string, err error) {
	strs := r.Form["themeID"]
	if strs != nil && len(strs) != 0 {
		themeID, err = strconv.Atoi(strs[0])
		if err != nil {
			err = errors.New("无法解析主题ID")
		}
	} else {
		err = errors.New("无法解析主题ID")
	}

	strs = r.Form["title"]
	if strs != nil && len(strs) != 0 {
		title = strs[0]
		if len(title) == 0 {
			err = errors.New("标题为空")
		}
	} else {
		err = errors.New("无标题")
	}

	strs = r.Form["content"]
	if strs != nil && len(strs) != 0 {
		content = strs[0]
		if len(content) == 0 {
			err = errors.New("内容为空")
		}
	} else {
		err = errors.New("无内容")
	}
	return
}

//NewPostCommit 新增帖子，提交,POST
func NewPostCommit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s := getExsitOrCreateNewSession(w, r, true)
	if s.User == nil {
		sendErrorPage(w, "请先登录")
		return
	}
	data := new(usecase.PostAddData)
	data.UserID = s.User.ID
	var err error
	r.ParseForm()
	data.ThemeID, data.Title, data.Content, err = readFormDataOfNewPost(r)
	if err != nil {
		sendErrorPage(w, err.Error())
		return
	}
	/*检查用户权限*/
	err = usecase.AddPost(data)
	if err != nil {
		sendErrorPage(w, err.Error())
		return
	}
	//成功的话，直接发主题页，无论任何方式排序，都是在主题第一位
	sendThemePage(w, data.ThemeID, 0, s)
}
