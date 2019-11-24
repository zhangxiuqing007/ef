package dba

import (
	"database/sql"
	"strconv"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func check1Err(num int64, err error) {
	checkErr(err)
	if num != 1 {
		panic("expect 1, get " + strconv.FormatInt(num, 10))
	}
}

func checkSqlResultErr(result sql.Result, err error) {
	checkErr(err)
	check1Err(result.RowsAffected())
}

type Scanner interface {
	Scan(dest ...interface{}) error
}
