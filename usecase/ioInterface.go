package usecase

import "ef/models"

var db IDataIO

//SetDbInstance 设置当前的db实现
func SetDbInstance(dbIns IDataIO) {
	db = dbIns
}

//IDataIO IO接口
type IDataIO interface {
	Open(linkInfo string)
	Clear()
	Close()

	AddTheme(theme *models.ThemeInDB)         //新增主题
	DeleteTheme(themeID int)                  //删除主题
	UpdateTheme(theme *models.ThemeInDB)      //更新主题（名称）
	QueryTheme(themeID int) *models.ThemeInDB //查询主题
	QueryAllThemes() []*models.ThemeInDB      //查询所有主题
	QueryPostCountOfTheme(themeID int) int    //查询主题的帖子量

	AddPost(post *models.PostInDB, cmt *models.CommentInDB)                               //新增帖子
	AddPosts(post []*models.PostInDB, cmt []*models.CommentInDB)                          //批量新增帖子
	DeletePost(postID int)                                                                //删除帖子
	QueryPost(postID int) *models.PostInDB                                                //查询DB帖子
	QueryPostTitle(postID int) string                                                     //查询帖子标题
	UpdatePostTitle(*models.PostInDB)                                                     //修改帖子标题
	QueryPostsOfTheme(themeID int, count, offset, sortType int) []*models.PostOnThemePage //查询主题下的帖子列表
	QueryPostsOfUser(userID int, count, offset int) []*models.PostOnThemePage             //查询用户发的帖子列表
	QueryPostOfPostPage(postID int) *models.PostOnPostPage                                //查询帖子页内，帖子的展示内容

	AddComment(comment *models.CommentInDB)                                                    //新增评论
	AddComments(comments []*models.CommentInDB)                                                //批量增加评论
	DeleteComment(cmtID int)                                                                   //删除评论
	QueryComment(cmtID int) *models.CommentInDB                                                //查询DB评论
	UpdateComment(cmt *models.CommentInDB)                                                     //修改评论
	QueryComments(postID int) []*models.CommentInDB                                            //查询DB评论
	QueryCommentsOfPostPage(postID int, count, offset int, userID int) []*models.CmtOnPostPage //查询帖子页内，评论的展示内容

	AddPbItem(pb *models.PBInDB)                      //新增赞踩行
	QueryPbItem(cmtID int, userID int) *models.PBInDB //查询赞踩行
	QueryPbsOfPost(postID int) []*models.PBInDB       //查询post相关的pbs
	Praise(pb *models.PBInDB)                         //赞
	Belittle(pb *models.PBInDB)                       //贬
	PraiseCancel(pb *models.PBInDB)                   //取消赞
	BelittleCancel(pb *models.PBInDB)                 //取消贬
	SetPB(cmtID int, userID int, isP bool, isD bool)  //尝试设置PB

	AddUser(user *models.UserInDB)                                             //新增用户
	DeleteUser(userID int)                                                     //删除用户
	QueryUserByID(userID int) *models.UserInDB                                 //通过id查询用户
	QueryUserByAccountAndPwd(account string, password string) *models.UserInDB //通过账户密码查询用户
	QueryPostCountOfUser(userID int) int                                       //查询用户的发帖量
	QueryImageCountOfUser(userID int) int                                      //查询用户上传图片数量
	UpdateUserHeadPhoto(userID int, path string)                               //更新头像文件路径
	IsUserNameExist(name string) bool                                          //用户名是否存在
	IsUserAccountExist(account string) bool                                    //用户账号是否存在

	AddImages(images []*models.ImageInDB)                   //批量新增图片
	QueryImages(userID int, count int, offset int) []string //图片查询
}

const (
	//PostSortTypeCreatedTime 排序类型：发帖时间
	PostSortTypeCreatedTime = iota
	//PostSortTypeLastCmtTime 排序类型：最终评论时间
	PostSortTypeLastCmtTime
)
