package dba

import (
	"testing"
)

//测试连接	go test -v -run TestLinkToMysqlServer
func TestLinkToMysqlServer(t *testing.T) {
	db := new(MySQLIns)
	err := db.Open("mysql5856")
	if err != nil {
		_, err = db.QueryAllThemes()
	}
	if err != nil {
		t.Fatalf("x失败：连接mysql服务器：" + err.Error())
	} else {
		t.Logf("成功：连接至mysql数据库")
	}
	db.Close()
}
