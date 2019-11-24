package usecase

import (
	"ef/models"
	"time"
)

//CmtAddData 新增评论需要提供的资料
type CmtAddData struct {
	PostID  int
	UserID  int
	Content string
}

func (data *CmtAddData) buildCmtDbIns() *models.CommentInDB {
	unixTime := time.Now().UnixNano()
	return &models.CommentInDB{
		ID:            0,
		PostID:        data.PostID,
		UserID:        data.UserID,
		Content:       data.Content,
		State:         models.CmtStateNormal,
		CreatedTime:   unixTime,
		LastEditTime:  unixTime,
		EditTimes:     1,
		PraiseTimes:   0,
		BelittleTimes: 0,
	}
}

//QueryCommentsOfPostPage 查询评论内容，用户帖子页内展示
func QueryCommentsOfPostPage(postID int, count, offset int, userID int) []*models.CmtOnPostPage {
	return db.QueryCommentsOfPostPage(postID, count, offset, userID)
}

//SetPB 尝试在满足条件时设置pb
func SetPB(cmtID int, userID int, isP bool, isD bool) {
	db.SetPB(cmtID, userID, isP, isD)
}

//AddComment 新增评论
func AddComment(cmt *CmtAddData) {
	db.AddComment(cmt.buildCmtDbIns())
}

//QueryComment 查询评论
func QueryComment(cmtID int) *models.CommentInDB {
	return db.QueryComment(cmtID)
}

//UpdateComment 更新帖子
func UpdateComment(cmt *models.CommentInDB) {
	db.UpdateComment(cmt)
}
