package models

//PBInDB 赞踩Item，数据库的结构
type PBInDB struct {
	ID     int
	CmtID  int
	UserID int

	PValue int
	PTime  int64
	PCTime int64

	BValue int
	BTime  int64
	BCTime int64
}
