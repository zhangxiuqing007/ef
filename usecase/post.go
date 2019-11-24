package usecase

import (
	"ef/models"
	"time"
)

//PostAddData 申请发布帖子的数据
type PostAddData struct {
	ThemeID int
	UserID  int

	Title   string
	Content string
}

func (data *PostAddData) buildPostDb() *models.PostInDB {
	return &models.PostInDB{
		ID:          0,
		ThemeID:     data.ThemeID,
		UserID:      data.UserID,
		Title:       data.Title,
		State:       models.PostStateNormal,
		CreatedTime: time.Now().UnixNano(),
		CmtCount:    0,
		LastCmterID: data.UserID,
		LastCmtTime: 0,
	}
}

func (data *PostAddData) buildCmtDb() *models.CommentInDB {
	unixTime := time.Now().UnixNano()
	return &models.CommentInDB{
		ID:            0,
		PostID:        0,
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

//QueryPost 帖子查询
func QueryPost(postID int) *models.PostInDB {
	return db.QueryPost(postID)
}

//QueryPostTitle 查询帖子标题
func QueryPostTitle(postID int) string {
	return db.QueryPostTitle(postID)
}

//UpdatePostTitle 更新帖子标题
func UpdatePostTitle(post *models.PostInDB) {
	db.UpdatePostTitle(post)
}

//QueryPostsOfTheme 查询帖子列表
func QueryPostsOfTheme(themeID int, count, offset, sortType int) []*models.PostOnThemePage {
	return db.QueryPostsOfTheme(themeID, count, offset, sortType)
}

//QueryPostsOfUser 查询某个用户发的帖子的列表
func QueryPostsOfUser(userID int, count, offset int) []*models.PostOnThemePage {
	return db.QueryPostsOfUser(userID, count, offset)
}

//QueryPostOfPostPage 帖子页内容查询
func QueryPostOfPostPage(postID int) *models.PostOnPostPage {
	return db.QueryPostOfPostPage(postID)
}

//AddPost 新增帖子
func AddPost(data *PostAddData) {
	//先成PostDB
	post := data.buildPostDb()
	//生成CmtDB
	cmt := data.buildCmtDb()
	//添加帖子
	db.AddPost(post, cmt)
}
