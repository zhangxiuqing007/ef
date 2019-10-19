package dba

import (
	"testing"
)

//测试连接	go test -v -run TestLinkToMysqlServer
func TestLinkToMysqlServer(t *testing.T) {
	db := new(MySQLIns)
	err := db.Open("root123")
	if err == nil {
		tms, _ := db.QueryAllThemes()
		t.Log(len(tms))
	}
	if err != nil {
		t.Fatal("x失败：连接mysql服务器：" + err.Error())
	} else {
		t.Log("成功：连接至mysql数据库")
	}
	checkErr(db.Close())
}
