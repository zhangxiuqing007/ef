package controllers

import (
	"fmt"
	"net/http"
)

type loginInfo struct {
	IsLogin  bool
	UserID   int
	UserName string
}

//创建登录信息
func buildLoginInfo(s *Session) *loginInfo {
	result := new(loginInfo)
	result.IsLogin = s.User != nil
	if result.IsLogin {
		result.UserID = s.User.ID
		result.UserName = s.User.Name
	}
	return result
}

type pageNavi struct {
	Path   string
	Number int
}

type pageNavis struct {
	HeadPageNavis []*pageNavi //前导航页
	CurrentPage   *pageNavi   //当前页
	TailPageNavis []*pageNavi //后导航页
}

//创建导航页vm
func buildPageNavis(pathBuilder func(index int) string, beginIndex, currentIndex, endIndex int) *pageNavis {
	//制作导航页切片
	buildNavis := func(bi, ei int) []*pageNavi {
		slice := make([]*pageNavi, 0, 16)
		for i := bi; i <= ei; i++ {
			slice = append(slice, &pageNavi{pathBuilder(i), i + 1})
		}
		return slice
	}
	//确定导航页限制
	p := new(pageNavis)
	p.HeadPageNavis = buildNavis(beginIndex, currentIndex-1)
	p.CurrentPage = &pageNavi{"", currentIndex + 1}
	p.TailPageNavis = buildNavis(currentIndex+1, endIndex)
	return p
}

//获取提供的导航页
func getNaviPageIndexs(
	currentPageIndex int, /*当前页索引*/
	countOnePage int, /*一页元素数量*/
	maxHalfNaviPageCount int, /*最大的导航页数量的一半*/
	elementTotalCount int) /*元素的总数量*/ (beginIndex, endIndex int) {
	//先计算beginIndex
	beginIndex = currentPageIndex - maxHalfNaviPageCount
	if beginIndex < 0 {
		beginIndex = 0
	}
	endIndex = limitPageIndex(currentPageIndex+maxHalfNaviPageCount, countOnePage, elementTotalCount)
	return
}

//限制页索引
func limitPageIndex(currentIndex int, countOnePage int, totalCount int) int {
	maxIndex := totalCount / countOnePage
	if currentIndex > maxIndex {
		currentIndex = maxIndex
	}
	if currentIndex < 0 {
		currentIndex = 0
	}
	return currentIndex
}

//请求的统一捕获异常处理函数
func recoverErrAndSendErrorPage(w http.ResponseWriter) {
	if err := recover(); err != nil {
		sendErrorPage(w, fmt.Sprintf("操作失败 %v", err))
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
