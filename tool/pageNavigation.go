package tool

//导航对象
type PageNavigation struct {
	Path   string
	Number int
}

//导航信息
type PageNavigations struct {
	HeadPages   []*PageNavigation //前导航页
	CurrentPage *PageNavigation   //当前页
	TailPages   []*PageNavigation //后导航页
}

//分页导航操作器
type PageNavigationOperator struct {
}

//创建导航页vm
func (operator *PageNavigationOperator) BuildPageNavigations(pathBuilder func(index int) string, beginIndex, currentIndex, endIndex int) *PageNavigations {
	//制作导航页切片
	builder := func(bi, ei int) []*PageNavigation {
		slice := make([]*PageNavigation, 0, ei-bi+1)
		for i := bi; i <= ei; i++ {
			slice = append(slice, &PageNavigation{pathBuilder(i), i + 1})
		}
		return slice
	}
	//确定导航页限制
	return &PageNavigations{
		HeadPages:   builder(beginIndex, currentIndex-1),
		CurrentPage: &PageNavigation{"", currentIndex + 1},
		TailPages:   builder(currentIndex+1, endIndex),
	}
}

//获取提供的导航页码限制
func (operator *PageNavigationOperator) GetNavigationPageLimitIndex(
	currentPageIndex int, /*当前页索引*/
	countOnePage int, /*一页元素数量*/
	maxHalfPageCount int, /*最大的导航页数量的一半*/
	elementTotalCount int) /*元素的总数量*/ (beginIndex, endIndex int) {
	//先计算beginIndex
	beginIndex = currentPageIndex - maxHalfPageCount
	if beginIndex < 0 {
		beginIndex = 0
	}
	endIndex = operator.LimitPageIndex(currentPageIndex+maxHalfPageCount, countOnePage, elementTotalCount)
	return
}

//限制页索引
func (operator *PageNavigationOperator) LimitPageIndex(currentIndex int, countOnePage int, totalCount int) int {
	maxIndex := totalCount / countOnePage
	if currentIndex > maxIndex {
		currentIndex = maxIndex
	}
	if currentIndex < 0 {
		currentIndex = 0
	}
	return currentIndex
}
