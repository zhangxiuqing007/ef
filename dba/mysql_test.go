package dba

import (
	"testing"
)

//测试连接	go test -v -run TestLinkToMysqlServer
func TestLinkToMysqlServer(t *testing.T) {
	db := new(MySQLIns)
	db.Open("")
	defer db.Close()
	t.Log(len(db.QueryAllThemes()))
	t.Log("成功：连接至mysql数据库")
}
