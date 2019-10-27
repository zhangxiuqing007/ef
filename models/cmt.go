package models

import (
	"ef/tool"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"
)

//CmtOnPostPage 评论在帖子页展示的信息
type CmtOnPostPage struct {
	//查询的值
	ID         int
	Content    string
	ContentESC template.HTML

	CmterID            int
	CmterName          string
	CmterHeadPhotoPath string

	CmtTime      int64
	CmtEditTimes int
	LastEditTime int64

	PraiseTimes   int
	BelittleTimes int

	IsPraised   bool
	IsBelittled bool

	//生成的值
	IndexStr     string
	CmtTimeF     string
	IsPChecked   string
	IsBChecked   string
	AllowEdit    bool
	CmtPageIndex int
}

//FormatCheckedStrOfPB 生成PB的值
func (cmt *CmtOnPostPage) FormatCheckedStrOfPB() {
	if cmt.IsPraised {
		cmt.IsPChecked = "checked"
	}
	if cmt.IsBelittled {
		cmt.IsBChecked = "checked"
	}
}

//FormatStringTime 生成文字类型的时间
func (cmt *CmtOnPostPage) FormatStringTime() {
	timeStr := tool.FormatTimeDetail(time.Unix(0, cmt.CmtTime))
	if cmt.CmtEditTimes >= 2 {
		cmt.CmtTimeF = fmt.Sprintf("初次评论后，修改过%d次，最后编辑时间：%s", cmt.CmtEditTimes-1, timeStr)
	} else {
		cmt.CmtTimeF = timeStr
	}
}

//FormatIndex 生成楼层字符
func (cmt *CmtOnPostPage) FormatIndex(index int) {
	if index == 0 {
		cmt.IndexStr = "楼主"
	} else {
		cmt.IndexStr = strconv.Itoa(index) + "楼"
	}
}

//FormatAllowEdit 生成AllowEdit的值
func (cmt *CmtOnPostPage) FormatAllowEdit(userID int) {
	if cmt.CmterID == userID {
		cmt.AllowEdit = true
	}
}

//FormatCmtPageIndex 生成当前 评论 页面的index
func (cmt *CmtOnPostPage) FormatCmtPageIndex(index int) {
	cmt.CmtPageIndex = index
}

//替换图片和对应的style
func (cmt *CmtOnPostPage) FormatImageWithStyle() {
	if strings.Contains(cmt.Content, "[img]") {
		content := cmt.Content
		//必须手动转义其他可能存在的尖括号
		if strings.Contains(content, "<") {
			content = strings.ReplaceAll(content, "<", "&lt")
			content = strings.ReplaceAll(content, ">", "&gt")
		}
		//再转义图片符号
		content = strings.ReplaceAll(content, "[img]", "<img")
		content = strings.ReplaceAll(content, "[/img]", "></img>")
		cmt.Content = content
	}
	cmt.ContentESC = template.HTML(cmt.Content)
}

//CommentInDB 评论，数据库形态
type CommentInDB struct {
	ID     int
	PostID int
	UserID int

	Content string

	State int

	CreatedTime  int64
	LastEditTime int64
	EditTimes    int

	PraiseTimes   int
	BelittleTimes int
}

const (
	//CmtStateNormal 评论状态：正常
	CmtStateNormal = iota
	//CmtStateLock 评论状态：锁定
	CmtStateLock
	//CmtStateHide 评论状态：隐藏
	CmtStateHide
)
