package myredis

import "testing"

func Test_CopyAllTypeDataFromMysqlToRedis(t *testing.T) {
	buildReadySqlIns().CopyAllTypeDataFromMysqlToRedis()
}
