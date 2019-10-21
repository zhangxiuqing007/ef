package controllers

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var cmtEditPageTemplate = template.Must(template.ParseFiles("views/cmt_get_edit.html", "views/comp/login_info_head.html"))

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
	//s := getExsitOrCreateNewSession(w, r, true)
	//if s.User == nil {
	//	w.Write([]byte("请先登录"))
	//	return
	//}
	////解析表单
	//cmtID, isP, isD, err := readFormDataFromCmtPb(r)
	//if err != nil {
	//	w.Write([]byte(err.Error()))
	//	return
	//}
	////处理，数据库输入内容
	//err = usecase.SetPB(cmtID, s.User.ID, isP, isD)
	//if err != nil {
	//	w.Write([]byte("赞踩请求失败"))
	//	return
	//}
	////返回
	//w.Write([]byte("赞踩成功"))
}
