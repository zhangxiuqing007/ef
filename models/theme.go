package models

import (
	"ef/tool"
	"time"
)

//ThemeInDB 主题
type ThemeInDB struct {
	ID        int
	Name      string
	PostCount int
}

//PostOnThemePage 帖子在主题页中展示的信息
type PostOnThemePage struct {
	ID    int
	Title string

	CmtCount int

	CreaterID   int
	CreaterName string
	CreatedTime int64

	LastCmterID   int
	LastCmterName string
	LastCmtTime   int64

	//生成的内容
	CreatedTimeF string
	LastCmtTimeF string
	IsHasCmt     bool
}

//FormatShowInfo 生成展示类信息
func (p *PostOnThemePage) FormatShowInfo() {
	p.CreatedTimeF = tool.FormatTimeDetail(time.Unix(0, p.CreatedTime))
	p.IsHasCmt = p.CmtCount > 0 || p.LastCmtTime > 0
	if p.IsHasCmt {
		p.LastCmtTimeF = tool.FormatTimeDetail(time.Unix(0, p.LastCmtTime))
	}
}
