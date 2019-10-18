package models

//PostOnPostPage 帖子，在帖子页展示的内容
type PostOnPostPage struct {
	ID        int
	CreaterID int
	Title     string

	CmtCount int

	ThemeID   int
	ThemeName string

	//需要生成的信息
	AllowEditTitle bool
}

//FormatShowInfo 生成展示用的信息
func (p *PostOnPostPage) FormatShowInfo(viewerID int) {
	p.AllowEditTitle = viewerID == p.CreaterID
}

//PostInDB 帖子，数据库形态
type PostInDB struct {
	ID      int
	ThemeID int
	UserID  int

	Title string

	State int

	CreatedTime int64

	CmtCount    int
	LastCmterID int
	LastCmtTime int64
}

const (
	//PostStateNormal 帖子状态：正常
	PostStateNormal = iota
	//PostStateLock 帖子状态：锁定
	PostStateLock
	//PostStateHide 帖子状态：隐藏
	PostStateHide
)
