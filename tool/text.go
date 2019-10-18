package tool

import (
	"errors"
)

//SplitText 拆分文本
func SplitText(s string, spe []rune) []string {
	if len(s) == 0 || spe == nil || len(spe) == 0 {
		panic(errors.New("参数错误"))
	}
	runes := append([]rune(s), spe[0])
	result := make([]string, 0, len(runes)/4)
	//是否包含方法
	isContain := func(r rune) bool {
		for _, v := range spe {
			if v == r {
				return true
			}
		}
		return false
	}
	index := -1
	for i, v := range runes {
		if isContain(v) {
			if index >= 0 {
				result = append(result, string(runes[index:i]))
				index = -1
			}
		} else {
			if index == -1 {
				index = i
			}
		}
	}
	return result
}
