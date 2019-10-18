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
	post := new(models.PostInDB)
	//post.ID = 0
	post.ThemeID = data.ThemeID
	post.UserID = data.UserID
	post.Title = data.Title
	post.State = models.PostStateNormal
	post.CreatedTime = time.Now().UnixNano()
	//post.CmtCount = 0
	post.LastCmterID = data.UserID
	//post.LastCmtTime = 0
	return post
}

func (data *PostAddData) buildCmtDb() *models.CommentInDB {
	cmt := new(models.CommentInDB)
	//cmt.ID =0
	//cmt.PostID =0
	cmt.UserID = data.UserID
	cmt.Content = data.Content
	cmt.State = models.CmtStateNormal
	cmt.CreatedTime = time.Now().UnixNano()
	cmt.LastEditTime = cmt.CreatedTime
	cmt.EditTimes = 1
	//cmt.PraiseTimes =0
	//cmt.BelittleTimes =0
	return cmt
}

//QueryPost 帖子查询
func QueryPost(postID int) (*models.PostInDB, error) {
	return db.QueryPost(postID)
}

//QueryPostTitle 查询帖子标题
func QueryPostTitle(postID int) (string, error) {
	return db.QueryPostTitle(postID)
}

//UpdatePostTitle 更新帖子标题
func UpdatePostTitle(post *models.PostInDB) error {
	return db.UpdatePostTitle(post)
}

//QueryPostsOfTheme 查询帖子列表
func QueryPostsOfTheme(themeID int, count, offset, sortType int) ([]*models.PostOnThemePage, error) {
	return db.QueryPostsOfTheme(themeID, count, offset, sortType)
}

//QueryPostsOfUser 查询某个用户发的帖子的列表
func QueryPostsOfUser(userID int, count, offset int) ([]*models.PostOnThemePage, error) {
	return db.QueryPostsOfUser(userID, count, offset)
}

//QueryPostOfPostPage 帖子页内容查询
func QueryPostOfPostPage(postID int) (*models.PostOnPostPage, error) {
	return db.QueryPostOfPostPage(postID)
}

//AddPost 新增帖子
func AddPost(data *PostAddData) error {
	//先成PostDB
	post := data.buildPostDb()
	//生成CmtDB
	cmt := data.buildCmtDb()
	//添加帖子
	//返回结果
	return db.AddPost(post, cmt)
}
