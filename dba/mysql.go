package dba

import (
	"database/sql"
	"errors"
	"fmt"

	//mysql driver
	_ "github.com/go-sql-driver/mysql"
)

//MySQLIns mysql数据库实现
type MySQLIns struct {
	sqlBase
}

//Open 打开
func (s *MySQLIns) Open(dbFilePath string) error {
	var err error
	s.DB, err = sql.Open("mysql", fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/efdb_bu?multiStatements=true", dbFilePath))
	return err
}

//Clear 清空
func (s *MySQLIns) Clear() error {
	return errors.New("mysql不支持清空全部数据，请手动操作")
}
