package dba

import (
	"database/sql"

	//sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

//SqliteIns sqlite实现
type SqliteIns struct {
	sqlBase
}

//Open 打开
func (s *SqliteIns) Open(dbFilePath string) error {
	var err error
	s.DB, err = sql.Open("sqlite3", dbFilePath)
	return err
}

//Clear 清空
func (s *SqliteIns) Clear() error {
	const sqlStrToClear = `
delete from tb_cmt;
delete from tb_post;
delete from tb_theme;
delete from tb_user;
vacuum;`
	_, err := s.Exec(sqlStrToClear)
	return err
}
