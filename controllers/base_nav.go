package controllers

//导航对象
type pageNavigation struct {
	Path   string
	Number int
}

//导航信息
type pageNavigations struct {
	HeadPages   []*pageNavigation //前导航页
	CurrentPage *pageNavigation   //当前页
	TailPages   []*pageNavigation //后导航页
}

//分页导航操作器
type pageNavigationOperator struct {
}

//创建导航页vm
func (operator *pageNavigationOperator) buildPageNavigations(pathBuilder func(index int) string, beginIndex, currentIndex, endIndex int) *pageNavigations {
	//制作导航页切片
	builder := func(bi, ei int) []*pageNavigation {
		slice := make([]*pageNavigation, 0, 16)
		for i := bi; i <= ei; i++ {
			slice = append(slice, &pageNavigation{pathBuilder(i), i + 1})
		}
		return slice
	}
	//确定导航页限制
	p := new(pageNavigations)
	p.HeadPages = builder(beginIndex, currentIndex-1)
	p.CurrentPage = &pageNavigation{"", currentIndex + 1}
	p.TailPages = builder(currentIndex+1, endIndex)
	return p
}

//获取提供的导航页码限制
func (operator *pageNavigationOperator) getNavigationPageLimitIndex(
	currentPageIndex int, /*当前页索引*/
	countOnePage int, /*一页元素数量*/
	maxHalfPageCount int, /*最大的导航页数量的一半*/
	elementTotalCount int) /*元素的总数量*/ (beginIndex, endIndex int) {
	//先计算beginIndex
	beginIndex = currentPageIndex - maxHalfPageCount
	if beginIndex < 0 {
		beginIndex = 0
	}
	endIndex = operator.limitPageIndex(currentPageIndex+maxHalfPageCount, countOnePage, elementTotalCount)
	return
}

//限制页索引
func (operator *pageNavigationOperator) limitPageIndex(currentIndex int, countOnePage int, totalCount int) int {
	maxIndex := totalCount / countOnePage
	if currentIndex > maxIndex {
		currentIndex = maxIndex
	}
	if currentIndex < 0 {
		currentIndex = 0
	}
	return currentIndex
}
