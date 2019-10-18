package controllers

import (
	"ef/models"
	"ef/tool"
	"ef/usecase"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

var userTemplate = template.Must(template.ParseFiles("views/user.html"))

type userVM struct {
	ID              int
	Name            string
	SignUpTime      string
	Type            string
	State           string
	LastOperateTime string

	PostTotalCount      int
	CmtTotalCount       int
	TotalPraisedTimes   int
	TotalBelittledTimes int
}

//UserInfo 查看用户个人资料
func UserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userIDStr := ps.ByName("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		sendErrorPage(w, err.Error())
		return
	}
	sendUserPage(w, userID, getExsitOrCreateNewSession(w, r, true))
}

//统计并发送用户资料页面
func sendUserPage(w http.ResponseWriter, userID int, s *Session) {
	//db统计用户信息
	saInfo, err := usecase.QueryUserByID(userID)
	if err != nil {
		sendErrorPage(w, err.Error())
		return
	}
	vm := new(userVM)
	vm.ID = saInfo.ID
	vm.Name = saInfo.Name
	vm.SignUpTime = tool.FormatTimeDetail(time.Unix(0, saInfo.SignUpTime))
	vm.Type = models.GetUserTypeShowName(saInfo.Type)
	vm.State = models.GetUserStateShowName(saInfo.State)
	if saInfo.LastEditTime == 0 {
		vm.LastOperateTime = "无"
	} else {
		vm.LastOperateTime = tool.FormatTimeDetail(time.Unix(0, saInfo.LastEditTime))
	}

	vm.PostTotalCount = saInfo.PostCount
	vm.CmtTotalCount = saInfo.CommentCount
	vm.TotalPraisedTimes = saInfo.PraiseTimes
	vm.TotalBelittledTimes = saInfo.BelittleTimes
	userTemplate.Execute(w, vm)
}
