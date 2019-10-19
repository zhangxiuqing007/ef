package controllers

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var loginInputTemplate = template.Must(template.ParseFiles("views/session_get.html"))
var loginSuccessTemplate = template.Must(template.ParseFiles("views/session_post.html"))

type loginVM struct {
	Tip string
}

func readFormDataOfLogin(r *http.Request) (account string, pwd string) {
	strs := r.Form["account"]
	if strs != nil && len(strs) != 0 {
		account = strs[0]
	}
	strs = r.Form["password"]
	if strs != nil && len(strs) != 0 {
		pwd = strs[0]
	}
	return
}

//Exit 登出
func Exit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s := getExsitOrCreateNewSession(w, r, true)
	s.User = nil
	sendIndexPage(w, s)
}

func sendIndexPage(w http.ResponseWriter, s *Session) {

}
