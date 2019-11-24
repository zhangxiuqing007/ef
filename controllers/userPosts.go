package controllers

import (
	"ef/tool"
	"ef/usecase"
	"fmt"
)

type userPostsFromData struct {
	UserID    int
	PageIndex int
}

type UserPostsController struct {
	baseController
}

func (c *UserPostsController) Get() {
	data := new(userPostsFromData)
	if err := c.ParseForm(data); err != nil || data.UserID <= 0 || data.PageIndex < 0 {
		c.send400("请求信息错误")
		return
	}
	vm := new(themeVm)
	vm.ThemeID = data.UserID
	vm.WebTitle = "边缘社区-用户发帖列表"
	totalPostCount := usecase.QueryPostCountOfUser(data.UserID)
	oper := new(tool.PageNavigationOperator)
	pageIndex := oper.LimitPageIndex(data.PageIndex, postCountOnePage, totalPostCount)
	vm.PostHeaders = usecase.QueryPostsOfUser(data.UserID, postCountOnePage, pageIndex*postCountOnePage)
	for _, v := range vm.PostHeaders {
		v.FormatShowInfo()
	}
	pathBuilder := func(i int) string {
		return fmt.Sprintf("/userPosts?UserID=%d&PageIndex=%d", data.UserID, i)
	}
	beginIndex, endIndex := oper.GetNavigationPageLimitIndex(pageIndex, postCountOnePage, halfPageCountToNavigationOfTheme, totalPostCount)
	nevis := oper.BuildPageNavigations(pathBuilder, beginIndex, pageIndex, endIndex)
	c.setNavigationVm(nevis)
	c.setLoginVmSelf()
	c.Data["vm"] = vm
	c.TplName = "theme_get.html"
}
