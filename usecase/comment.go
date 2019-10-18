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
	cmt := new(models.CommentInDB)
	//cmt.ID = 0
	cmt.PostID = data.PostID
	cmt.UserID = data.UserID
	cmt.Content = data.Content
	cmt.State = models.CmtStateNormal
	cmt.CreatedTime = time.Now().UnixNano()
	cmt.LastEditTime = cmt.CreatedTime
	cmt.EditTimes = 1
	//cmt.PraiseTimes = 0
	//cmt.BelittleTimes =0
	return cmt
}

//QueryCommentsOfPostPage 查询评论内容，用户帖子页内展示
func QueryCommentsOfPostPage(postID int, count, offset int, userID int) ([]*models.CmtOnPostPage, error) {
	return db.QueryCommentsOfPostPage(postID, count, offset, userID)
}

//SetPB 尝试在满足条件时设置pb
func SetPB(cmtID int, userID int, isP bool, isD bool) error {
	return db.SetPB(cmtID, userID, isP, isD)
}

//AddComment 新增评论
func AddComment(cmt *CmtAddData) error {
	return db.AddComment(cmt.buildCmtDbIns())
}

//QueryComment 查询评论
func QueryComment(cmtID int) (*models.CommentInDB, error) {
	return db.QueryComment(cmtID)
}

//UpdateComment 更新帖子
func UpdateComment(cmt *models.CommentInDB) error {
	return db.UpdateComment(cmt)
}
