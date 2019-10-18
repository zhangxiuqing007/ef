package controllers

import (
	"ef/usecase"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var userRegistInputTemplate = template.Must(template.ParseFiles("views/registInput.html"))
var userRegistSuccessTemplate = template.Must(template.ParseFiles("views/registSuccess.html"))

type registInputVM struct {
	Tip string
}

//UserRegist 新用户注册
func UserRegist(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userRegistInputTemplate.Execute(w, registInputVM{Tip: "请输入注册资料"})
}

func readFormDataFromRegist(r *http.Request) (name string, account string, pwd1 string, pwd2 string) {
	strs := r.Form["name"]
	if strs != nil && len(strs) != 0 {
		name = strs[0]
	}
	strs = r.Form["account"]
	if strs != nil && len(strs) != 0 {
		account = strs[0]
	}
	strs = r.Form["password1"]
	if strs != nil && len(strs) != 0 {
		pwd1 = strs[0]
	}
	strs = r.Form["password2"]
	if strs != nil && len(strs) != 0 {
		pwd2 = strs[0]
	}
	return
}

//UserRegistCommit 提交注册信息
func UserRegistCommit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		userRegistInputTemplate.Execute(w, registInputVM{Tip: "请输入完整的注册资料"})
		return
	}
	name, account, pwd1, pwd2 := readFormDataFromRegist(r)
	//如果缺少任意一个注册资料
	if len(name) == 0 || len(account) == 0 || len(pwd1) == 0 || len(pwd2) == 0 {
		userRegistInputTemplate.Execute(w, registInputVM{Tip: "请输入完整的注册资料"})
		return
	}
	//后端再次检查一遍
	if pwd1 != pwd2 {
		userRegistInputTemplate.Execute(w, registInputVM{Tip: "两次密码输入不一致"})
		return
	}
	//组织申请数据
	data := &usecase.UserSignUpData{
		Name:     name,
		Account:  account,
		Password: pwd1,
	}
	//调用用例层代码，尝试添加账户，并返回错误
	err = usecase.AddUser(data)
	if err != nil {
		userRegistInputTemplate.Execute(w, registInputVM{Tip: err.Error()})
		return
	}
	//注册成功
	userRegistSuccessTemplate.Execute(w, name)
}
