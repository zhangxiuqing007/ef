package tool

import (
	"fmt"
	"strconv"
	"testing"
)

func TestNavigationOperation(t *testing.T) {
	index := 1
	const imgCountOnePage = 8
	const totalImgCount = 10
	const halfPageCountToNavigationOfImage = 8
	oper := new(PageNavigationOperator)
	pageIndex := oper.LimitPageIndex(index, imgCountOnePage, totalImgCount)
	beginIndex, endIndex := oper.GetNavigationPageLimitIndex(pageIndex, imgCountOnePage, halfPageCountToNavigationOfImage, totalImgCount)
	pathBuilder := func(i int) string {
		return "/img?ImagePageIndex=" + strconv.Itoa(i)
	}
	result := oper.BuildPageNavigations(pathBuilder, beginIndex, pageIndex, endIndex)
	fmt.Print(result)
}
