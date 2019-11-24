package myredis

import "testing"

func buildReadySqlIns() *MySQLRedisIns {
	ins := new(MySQLRedisIns)
	ins.Open("")
	return ins
}

func Test_QueryAllThemes(t *testing.T) {
	tms := buildReadySqlIns().QueryAllThemes()
	for _, tm := range tms {
		t.Log(tm)
	}
}
