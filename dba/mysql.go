package dba

import (
	"database/sql"
	"errors"
	"fmt"

	//mysql driver
	_ "github.com/go-sql-driver/mysql"
)

//默认值
var MysqlUser = "root"
var MysqlPwd = "mysql5856"
var MysqlDb = "efdb_bu"

//MySQLIns mysql数据库实现
type MySQLIns struct {
	sqlBase
}

//Open 打开 这里的参数没用
func (s *MySQLIns) Open(dbFilePath string) error {
	var err error
	s.DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?multiStatements=true", MysqlUser, MysqlPwd, MysqlDb))
	return err
}

//Clear 清空
func (s *MySQLIns) Clear() error {
	return errors.New("mysql不支持清空全部数据，请手动操作")
}
