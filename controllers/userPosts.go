package controllers

import (
	"ef/usecase"
	"fmt"
)

type userPostsFromData struct {
	UserID    int
	PageIndex int
}

type UserPostsController struct {
	SessionBaseController
}

func (c *UserPostsController) Get() {
	data := new(userPostsFromData)
	if err := c.ParseForm(data); err != nil || data.UserID <= 0 || data.PageIndex < 0 {
		c.send400()
		return
	}
	vm := new(themeVm)
	vm.ThemeID = data.UserID
	vm.WebTitle = "边缘社区-用户发帖列表"
	totalPostCount, err := usecase.QueryPostCountOfUser(data.UserID)
	if err != nil {
		c.send404()
		return
	}
	pageIndex := limitPageIndex(data.PageIndex, postCountOnePage, totalPostCount)
	vm.PostHeaders, err = usecase.QueryPostsOfUser(data.UserID, postCountOnePage, pageIndex*postCountOnePage)
	if err != nil {
		c.send404()
		return
	}
	for _, v := range vm.PostHeaders {
		v.FormatShowInfo()
	}
	pathBuilder := func(i int) string {
		return fmt.Sprintf("/userPosts?UserID=%d&PageIndex=%d", data.UserID, i)
	}
	beginIndex, endIndex := getNavigationPageLimitIndex(pageIndex, postCountOnePage, halfPageCountToNavigationOfTheme, totalPostCount)
	nevis := buildPageNavigations(pathBuilder, beginIndex, pageIndex, endIndex)
	c.setNavigationVm(nevis)
	c.setLoginVmSelf()
	c.Data["vm"] = vm
	c.TplName = "theme_get.html"
}
