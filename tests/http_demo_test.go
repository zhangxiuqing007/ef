package tests

import (
	"net/http"
	"testing"
)

const testURL = "http://127.0.0.1:8080"

//尝试用http的Client对象发起http请求
func Test_Get_Index(t *testing.T) {
	c := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	r, err := c.Get(testURL)
	if err != nil {
		panic(err)
	}
	t.Log(r.Cookies())
}
