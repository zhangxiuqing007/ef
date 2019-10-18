package tool

import (
	"io/ioutil"
)

//MustStr 一旦有错误就引起崩溃
func MustStr(str string, err error) string {
	if err != nil {
		panic(err)
	}
	return str
}

//ReadAllTextUtf8 读取utf8编码文本文件的全部内容并转换成string
func ReadAllTextUtf8(filePath string) (string, error) {
	buffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}
