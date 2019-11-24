package dba

import (
	"crypto/md5"
	"database/sql"
	"ef/models"
	"encoding/hex"
	"errors"
	"time"
)

type sqlBase struct {
	*sql.DB
}

//Close 关闭
func (s *sqlBase) Close() {
	checkErr(s.DB.Close())
}

const sqlStrToAddTheme = "insert into tb_theme (tm_name,tm_post_count) values (?,?)"

//AddTheme 增加主题
func (s *sqlBase) AddTheme(theme *models.ThemeInDB) {
	result, err := s.Exec(sqlStrToAddTheme, theme.Name, theme.PostCount)
	checkErr(err)
	check1Err(result.RowsAffected())
	tmID, err := result.LastInsertId()
	checkErr(err)
	theme.ID = int(tmID)
}

const sqlStrToDeleteTheme = "delete from tb_theme where tm_id = ?"

//DeleteTheme 删除主题(不处理下面的帖子和评论，不影响计数器)
func (s *sqlBase) DeleteTheme(themeID int) {
	checkSqlResultErr(s.Exec(sqlStrToDeleteTheme, themeID))
}

const sqlStrToUpdateTheme = "update tb_theme set tm_name = ? where tm_id = ?"

//UpdateTheme 更新主题名称
func (s *sqlBase) UpdateTheme(theme *models.ThemeInDB) {
	checkSqlResultErr(s.Exec(sqlStrToUpdateTheme, theme.Name, theme.ID))
}

const sqlStrToQueryTheme = "select tm_id,tm_name,tm_post_count from tb_theme where tm_id = ?"

//QueryTheme 查询某个主题
func (s *sqlBase) QueryTheme(themeID int) *models.ThemeInDB {
	return s.scanFromThemeTbAllFields(s.QueryRow(sqlStrToQueryTheme, themeID))
}

const sqlStrToQueryThemes = "select * from tb_theme order by tm_id"

//QueryAllThemes 查询主题列表
func (s *sqlBase) QueryAllThemes() []*models.ThemeInDB {
	rows, err := s.Query(sqlStrToQueryThemes)
	checkErr(err)
	defer rows.Close()
	tms := make([]*models.ThemeInDB, 0, 16)
	for rows.Next() {
		if tm := s.scanFromThemeTbAllFields(rows); tm == nil {
			panic(errors.New("存在无效主题"))
		} else {
			tms = append(tms, tm)
		}
	}
	return tms
}

func (s *sqlBase) scanFromThemeTbAllFields(scanner Scanner) *models.ThemeInDB {
	tm := new(models.ThemeInDB)
	if scanner.Scan(&tm.ID, &tm.Name, &tm.PostCount) != nil {
		return nil
	}
	return tm
}

const sqlStrToQueryPostCountOfTheme = "select tm_post_count from tb_theme where tm_id = ?"

//QueryPostCountOfTheme 查询目标主题的帖子数量
func (s *sqlBase) QueryPostCountOfTheme(themeID int) int {
	var count int
	checkErr(s.QueryRow(sqlStrToQueryPostCountOfTheme, themeID).Scan(&count))
	return count
}

const sqlStrToAddPostPart1 = "update tb_user set ur_post_count = ur_post_count+1 where ur_id = ?;"
const sqlStrToAddPostPart2 = "update tb_theme set tm_post_count = tm_post_count+1 where tm_id = ?;"
const sqlStrToAddPostPart3 = `
insert into tb_post (
	po_theme_id,
	po_user_id,
	po_title,
	po_state,
	po_post_time,
	po_cmt_count,
	po_lcmter_id,
	po_lcmt_time)
	values 
	(?,?,?,?,?,?,?,?);`

//AddPost 新增帖子
func (s *sqlBase) AddPost(post *models.PostInDB, cmt *models.CommentInDB) {
	var err error
	insertPostResult, err := s.Exec(sqlStrToAddPostPart3, post.ThemeID, post.UserID, post.Title, post.State, post.CreatedTime, post.CmtCount, post.LastCmterID, post.LastCmtTime)
	checkErr(err)
	postID, err := insertPostResult.LastInsertId()
	checkErr(err)
	post.ID = int(postID)
	check1Err(insertPostResult.RowsAffected())
	cmt.PostID = post.ID
	tx, _ := s.Begin()
	checkSqlResultErr(tx.Exec(sqlStrToAddPostPart1, post.UserID))
	checkSqlResultErr(tx.Exec(sqlStrToAddPostPart2, post.ThemeID))
	insertCmtResult, err := tx.Exec(sqlStrToAddCommentPart1, cmt.PostID, cmt.UserID, cmt.Content, cmt.State, cmt.CreatedTime, cmt.LastEditTime, cmt.EditTimes, cmt.PraiseTimes, cmt.BelittleTimes)
	checkErr(err)
	cmtID, err := insertCmtResult.LastInsertId()
	checkErr(err)
	cmt.ID = int(cmtID)
	check1Err(insertCmtResult.RowsAffected())
	checkSqlResultErr(tx.Exec(sqlStrToAddCommentPart3, cmt.LastEditTime, cmt.UserID))
	checkErr(tx.Commit())
}

//AddPosts 批量新增帖子
func (s *sqlBase) AddPosts(posts []*models.PostInDB, cmts []*models.CommentInDB) {
	tx, _ := s.Begin()
	stmt1, _ := tx.Prepare(sqlStrToAddPostPart1)
	stmt2, _ := tx.Prepare(sqlStrToAddPostPart2)
	stmt3, _ := tx.Prepare(sqlStrToAddPostPart3)
	stmt4, _ := tx.Prepare(sqlStrToAddCommentPart1)
	stmt5, _ := tx.Prepare(sqlStrToAddCommentPart3)
	for i, v := range posts {
		cmt := cmts[i]
		insertPostResult, err := stmt3.Exec(v.ThemeID, v.UserID, v.Title, v.State, v.CreatedTime, v.CmtCount, v.LastCmterID, v.LastCmtTime)
		checkErr(err)
		pid, err := insertPostResult.LastInsertId()
		checkErr(err)
		v.ID = int(pid)
		cmt.PostID = v.ID
		check1Err(insertPostResult.RowsAffected())
		checkSqlResultErr(stmt1.Exec(v.UserID))
		checkSqlResultErr(stmt2.Exec(v.ThemeID))
		insertCmtResult, err := stmt4.Exec(cmt.PostID, cmt.UserID, cmt.Content, cmt.State, cmt.CreatedTime, cmt.LastEditTime, cmt.EditTimes, cmt.PraiseTimes, cmt.BelittleTimes)
		checkErr(err)
		cmtID, err := insertCmtResult.LastInsertId()
		checkErr(err)
		cmt.ID = int(cmtID)
		check1Err(insertCmtResult.RowsAffected())
		checkSqlResultErr(stmt5.Exec(cmt.LastEditTime, cmt.UserID))
	}
	checkErr(tx.Commit())
}

const sqlStrToDeletePost = "delete from tb_post where po_id = ?;"

//DeletePost 删除帖子（但不影响任何计数器）
func (s *sqlBase) DeletePost(postID int) {
	checkSqlResultErr(s.Exec(sqlStrToDeletePost, postID))
}

const sqlStrToQueryPost = "select * from tb_post where po_id = ?"

//QueryPost 查询帖子
func (s *sqlBase) QueryPost(postID int) *models.PostInDB {
	return s.ScanFromPostTbAllFields(s.QueryRow(sqlStrToQueryPost, postID))
}

func (s *sqlBase) ScanFromPostTbAllFields(scanner Scanner) *models.PostInDB {
	post := new(models.PostInDB)
	if scanner.Scan(
		&post.ID,
		&post.ThemeID,
		&post.UserID,
		&post.Title,
		&post.State,
		&post.CreatedTime,
		&post.CmtCount,
		&post.LastCmterID,
		&post.LastCmtTime) != nil {
		return nil
	}
	return post
}

const sqlStrToQueryPostTitle = "select po_title from tb_post where po_id = ?;"

//QueryPostTitle 查询标题
func (s *sqlBase) QueryPostTitle(postID int) string {
	var title string
	checkErr(s.QueryRow(sqlStrToQueryPostTitle, postID).Scan(&title))
	return title
}

const sqlStrToUpdatePostTitle = "update tb_post set po_title = ? where po_id = ?"

//UpdatePostTitle 修改帖子标题
func (s *sqlBase) UpdatePostTitle(post *models.PostInDB) {
	tx, _ := s.Begin()
	_, err := tx.Exec(sqlStrToUpdatePostTitle, post.Title, post.ID)
	checkErr(err)
	checkSqlResultErr(tx.Exec(sqlStrToAddCommentPart3, post.LastCmtTime, post.UserID))
	checkErr(tx.Commit())
}

const sqlStrToQueryPostsSortType0 = `
select
    p.po_id,
    p.po_title,
    p.po_cmt_count,
    u1.ur_id,
    u1.ur_name,
    p.po_post_time,
    u2.ur_id,
    u2.ur_name,
    p.po_lcmt_time
from
    (select po_id from tb_post where po_theme_id = ? order by po_id desc limit ? offset ?) as p1,
    tb_post as p,
    tb_user as u1,
    tb_user as u2
where
    p1.po_id = p.po_id
    and p.po_user_id = u1.ur_id
    and p.po_lcmter_id = u2.ur_id`

const sqlStrToQueryPostsSortType1 = `
select
    p.po_id,
    p.po_title,
    p.po_cmt_count,
    u1.ur_id,
    u1.ur_name,
    p.po_post_time,
    u2.ur_id,
    u2.ur_name,
    p.po_lcmt_time
from
    (select po_id from tb_post where po_theme_id = ? order by po_lcmt_time desc limit ? offset ?) as p1,
    tb_post as p,
    tb_user as u1,
    tb_user as u2
where
    p1.po_id = p.po_id
    and p.po_user_id = u1.ur_id
    and p.po_lcmter_id = u2.ur_id`

//QueryPostsOfTheme 查询某主题下的帖子列表
func (s *sqlBase) QueryPostsOfTheme(themeID int, count, offset, sortType int) []*models.PostOnThemePage {
	if sortType == 0 {
		rows, err := s.Query(sqlStrToQueryPostsSortType0, themeID, count, offset)
		checkErr(err)
		return turnToPostsOnThemePage(rows, count)
	}
	rows, err := s.Query(sqlStrToQueryPostsSortType1, themeID, count, offset)
	checkErr(err)
	return turnToPostsOnThemePage(rows, count)
}

const sqlStrToQueryPostsOfUser = `
select
    p.po_id,
    p.po_title,
    p.po_cmt_count,
    u1.ur_id,
    u1.ur_name,
    p.po_post_time,
    u2.ur_id,
    u2.ur_name,
    p.po_lcmt_time
from
    (select po_id from tb_post where po_user_id = ? order by po_id desc limit ? offset ?) as p1,
    tb_post as p,
    tb_user as u1,
    tb_user as u2
where
    p1.po_id = p.po_id
    and p.po_user_id = u1.ur_id
    and p.po_lcmter_id = u2.ur_id`

//QueryPostsOfUser 查询某用户的帖子列表
func (s *sqlBase) QueryPostsOfUser(userID int, count, offset int) []*models.PostOnThemePage {
	rows, err := s.Query(sqlStrToQueryPostsOfUser, userID, count, offset)
	checkErr(err)
	return turnToPostsOnThemePage(rows, count)
}

func turnToPostsOnThemePage(rows *sql.Rows, cap int) []*models.PostOnThemePage {
	defer rows.Close()
	posts := make([]*models.PostOnThemePage, 0, cap)
	for rows.Next() {
		post := new(models.PostOnThemePage)
		checkErr(rows.Scan(
			&post.ID,
			&post.Title,
			&post.CmtCount,
			&post.CreaterID,
			&post.CreaterName,
			&post.CreatedTime,
			&post.LastCmterID,
			&post.LastCmterName,
			&post.LastCmtTime))
		posts = append(posts, post)
	}
	return posts
}

const sqlStrToQueryPostOfPostPage = `
select 
   p.po_title,
   p.po_user_id,
   p.po_cmt_count,
   tm.tm_id,
   tm.tm_name
from
    tb_post as p,
	tb_theme as tm
where
	p.po_id = ? and tm.tm_id = p.po_theme_id`

//QueryPostOfPostPage 查询帖子页的帖子内容
func (s *sqlBase) QueryPostOfPostPage(postID int) *models.PostOnPostPage {
	post := new(models.PostOnPostPage)
	post.ID = postID
	checkErr(s.QueryRow(sqlStrToQueryPostOfPostPage, postID).Scan(
		&post.Title,
		&post.CreaterID,
		&post.CmtCount,
		&post.ThemeID,
		&post.ThemeName))
	return post
}

const sqlStrToAddCommentPart1 = `
insert into tb_cmt 
(
	cmt_post_id,
	cmt_user_id,
	cmt_content,
	cmt_state,
	cmt_c_time,
	cmt_le_time,
	cmt_e_times,
	cmt_p_times,
	cmt_b_times
)
values (?,?,?,?,?,?,?,?,?);`
const sqlStrToAddCommentPart2 = "update tb_user set ur_cmt_count = ur_cmt_count+1 where ur_id = ?;"
const sqlStrToAddCommentPart3 = "update tb_user set ur_le_time = ? where ur_id = ?;"
const sqlStrToAddCommentPart4 = "update tb_post set po_cmt_count = po_cmt_count+1 where po_id = ?;"
const sqlStrToAddCommentPart5 = "update tb_post set po_lcmter_id = ? where po_id = ?;"
const sqlStrToAddCommentPart6 = "update tb_post set po_lcmt_time = ? where po_id = ?;"

//AddComment 增加评论
func (s *sqlBase) AddComment(cmt *models.CommentInDB) {
	tx, _ := s.Begin()
	cmtInsertResult, err := s.Exec(sqlStrToAddCommentPart1, cmt.PostID, cmt.UserID, cmt.Content, cmt.State, cmt.CreatedTime, cmt.LastEditTime, cmt.EditTimes, cmt.PraiseTimes, cmt.BelittleTimes)
	checkErr(err)
	cmtID, err := cmtInsertResult.LastInsertId()
	checkErr(err)
	cmt.ID = int(cmtID)
	check1Err(cmtInsertResult.RowsAffected())
	checkSqlResultErr(tx.Exec(sqlStrToAddCommentPart2, cmt.UserID))
	checkSqlResultErr(tx.Exec(sqlStrToAddCommentPart3, cmt.LastEditTime, cmt.UserID))
	checkSqlResultErr(tx.Exec(sqlStrToAddCommentPart4, cmt.PostID))
	_, err = tx.Exec(sqlStrToAddCommentPart5, cmt.UserID, cmt.PostID)
	checkErr(err)
	checkSqlResultErr(tx.Exec(sqlStrToAddCommentPart6, cmt.LastEditTime, cmt.PostID))
	checkErr(tx.Commit())
}

//AddComments 批量增加评论
func (s *sqlBase) AddComments(comments []*models.CommentInDB) {
	tx, _ := s.Begin()
	stmt1, _ := tx.Prepare(sqlStrToAddCommentPart1)
	stmt2, _ := tx.Prepare(sqlStrToAddCommentPart2)
	stmt3, _ := tx.Prepare(sqlStrToAddCommentPart3)
	stmt4, _ := tx.Prepare(sqlStrToAddCommentPart4)
	stmt5, _ := tx.Prepare(sqlStrToAddCommentPart5)
	stmt6, _ := tx.Prepare(sqlStrToAddCommentPart6)
	for _, v := range comments {
		result, err := stmt1.Exec(v.PostID, v.UserID, v.Content, v.State, v.CreatedTime, v.LastEditTime, v.EditTimes, v.PraiseTimes, v.BelittleTimes)
		checkErr(err)
		cid, err := result.LastInsertId()
		checkErr(err)
		v.ID = int(cid)
		check1Err(result.RowsAffected())
		checkSqlResultErr(stmt2.Exec(v.UserID))
		checkSqlResultErr(stmt3.Exec(v.LastEditTime, v.UserID))
		checkSqlResultErr(stmt4.Exec(v.PostID))
		_, err = stmt5.Exec(v.UserID, v.PostID)
		checkErr(err)
		checkSqlResultErr(stmt6.Exec(v.LastEditTime, v.PostID))
	}
	checkErr(tx.Commit())
}

const sqlStrToDeleteComment = "delete from tb_cmt where cmt_id =?;"

//DeleteComment 删除单个评论，不处理任何计数器
func (s *sqlBase) DeleteComment(cmtID int) {
	checkSqlResultErr(s.Exec(sqlStrToDeleteComment, cmtID))
}

const sqlStrToQueryComment = "select * from tb_cmt where cmt_id = ?"

//QueryComment 查询单个评论
func (s *sqlBase) QueryComment(cmtID int) *models.CommentInDB {
	return s.scanFromCmtTbAllFields(s.QueryRow(sqlStrToQueryComment, cmtID))
}

const sqlStrToUpdateCommentPart1 = "update tb_cmt set cmt_content = ?, cmt_le_time = ?, cmt_e_times= ? where cmt_id = ? "
const sqlStrToUpdateCommentPart2 = "update tb_user set ur_le_time = ? where ur_id = ? "

//UpdateComment 修改评论内容
func (s *sqlBase) UpdateComment(cmt *models.CommentInDB) {
	tx, _ := s.Begin()
	checkSqlResultErr(tx.Exec(sqlStrToUpdateCommentPart1, cmt.Content, cmt.LastEditTime, cmt.EditTimes, cmt.ID))
	_, err := tx.Exec(sqlStrToUpdateCommentPart2, cmt.LastEditTime, cmt.PostID)
	checkErr(err)
	checkErr(tx.Commit())
}

const sqlStrToQueryComments = "select * from tb_cmt where cmt_post_id = ? order by cmt_id"

//QueryComments 查询评论，按照创建时间排序
func (s *sqlBase) QueryComments(postID int) []*models.CommentInDB {
	rows, err := s.Query(sqlStrToQueryComments, postID)
	checkErr(err)
	defer rows.Close()
	cmts := make([]*models.CommentInDB, 0, 32)
	for rows.Next() {
		cmt := s.scanFromCmtTbAllFields(rows)
		if cmt == nil {
			panic(errors.New("无效的评论"))
		}
		cmts = append(cmts, cmt)
	}
	return cmts
}

func (s *sqlBase) scanFromCmtTbAllFields(scanner Scanner) *models.CommentInDB {
	cmt := new(models.CommentInDB)
	if scanner.Scan(
		&cmt.ID,
		&cmt.PostID,
		&cmt.UserID,
		&cmt.Content,
		&cmt.State,
		&cmt.CreatedTime,
		&cmt.LastEditTime,
		&cmt.EditTimes,
		&cmt.PraiseTimes,
		&cmt.BelittleTimes) != nil {
		return nil
	}
	return cmt
}

const sqlStrToQueryCommentsOfPostPagePart1 = `
select 
       cmt.cmt_id,
	   cmt.cmt_content,
       u.ur_id,
	   u.ur_name,
	   u.ur_hp_path,
	   cmt.cmt_c_time,
	   cmt.cmt_e_times,
	   cmt.cmt_le_time,
       cmt.cmt_p_times,
       cmt.cmt_b_times
from
     (select cmt_id from tb_cmt where cmt_post_id = ? order by cmt_id limit ? offset ?) as c1,
	 tb_cmt as cmt,
	 tb_user as u
where 
	 c1.cmt_id = cmt.cmt_id
	 and cmt.cmt_user_id = u.ur_id`

const sqlStrToQueryCommentsOfPostPagePart2 = `
select 
       c1.cmt_id,
       pb.pb_p_value,
       pb.pb_b_value
from
     (select cmt_id from tb_cmt where cmt_post_id = ? order by cmt_id limit ? offset ?) as c1
     join tb_pb as pb on c1.cmt_id = pb.pb_cmt_id
where 
	 pb.pb_user_id = ?;`

//QueryCommentsOfPostPage 查询评论内容，用于显示在帖子页中
func (s *sqlBase) QueryCommentsOfPostPage(postID int, count int, offset int, userID int) []*models.CmtOnPostPage {
	cmts := make([]*models.CmtOnPostPage, 0, count)
	rows, err := s.Query(sqlStrToQueryCommentsOfPostPagePart1, postID, count, offset)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		cmt := new(models.CmtOnPostPage)
		checkErr(rows.Scan(
			&cmt.ID,
			&cmt.Content,
			&cmt.CmterID,
			&cmt.CmterName,
			&cmt.CmterHeadPhotoPath,
			&cmt.CmtTime,
			&cmt.CmtEditTimes,
			&cmt.LastEditTime,
			&cmt.PraiseTimes,
			&cmt.BelittleTimes))
		cmts = append(cmts, cmt)
	}
	//查看指定用户的赞踩情况
	rowspb, err := s.Query(sqlStrToQueryCommentsOfPostPagePart2, postID, count, offset, userID)
	checkErr(err)
	defer rowspb.Close()
	f := func(id int) *models.CmtOnPostPage {
		for _, v := range cmts {
			if v.ID == id {
				return v
			}
		}
		return nil
	}
	var tempID int
	var tempPvalue int
	var tempBvalue int
	for rowspb.Next() {
		checkErr(rowspb.Scan(&tempID, &tempPvalue, &tempBvalue))
		targetCmt := f(tempID)
		if targetCmt != nil {
			targetCmt.IsPraised = tempPvalue == 1
			targetCmt.IsBelittled = tempBvalue == 1
		}
	}
	return cmts
}

const sqlStrToAddPbItem = `
insert into tb_pb
(
	pb_cmt_id,
	pb_user_id,
	pb_p_value,
	pb_p_time,
	pb_p_ctime,
	pb_b_value,
	pb_b_time,
	pb_b_ctime
)values(?,?,?,?,?,?,?,?);`

//AddPbItem 新增赞踩行
func (s *sqlBase) AddPbItem(pb *models.PBInDB) {
	result, err := s.Exec(sqlStrToAddPbItem,
		pb.CmtID, pb.UserID,
		pb.PValue, pb.PTime, pb.PCTime,
		pb.BValue, pb.BTime, pb.BCTime)
	checkErr(err)
	pbID, err := result.LastInsertId()
	checkErr(err)
	pb.ID = int(pbID)
	check1Err(result.RowsAffected())
}

const sqlStrToQueryPbItem = "select * from tb_pb where pb_cmt_id = ? and pb_user_id = ?"

//QueryPbItem 查询赞踩行
func (s *sqlBase) QueryPbItem(cmtID int, userID int) *models.PBInDB {
	return s.scanFromPbTbAllFields(s.QueryRow(sqlStrToQueryPbItem, cmtID, userID))
}

const sqlStrToQueryPbsOfPost = "select * from tb_pb where pb_cmt_id in (select cmt_id from tb_cmt where cmt_post_id = ?)"

//QueryPbsOfPost 查询帖子的赞踩信息
func (s *sqlBase) QueryPbsOfPost(postID int) []*models.PBInDB {
	rows, err := s.DB.Query(sqlStrToQueryPbsOfPost, postID)
	checkErr(err)
	defer rows.Close()
	pbs := make([]*models.PBInDB, 0)
	for rows.Next() {
		pb := s.scanFromPbTbAllFields(rows)
		if pb == nil {
			panic(errors.New("无效的赞踩记录"))
		}
		pbs = append(pbs, pb)
	}
	return pbs
}

func (s *sqlBase) scanFromPbTbAllFields(scanner Scanner) *models.PBInDB {
	pbIns := new(models.PBInDB)
	if scanner.Scan(
		&pbIns.ID,
		&pbIns.CmtID,
		&pbIns.UserID,
		&pbIns.PValue,
		&pbIns.PTime,
		&pbIns.PCTime,
		&pbIns.BValue,
		&pbIns.BTime,
		&pbIns.BCTime) != nil {
		return nil
	}
	return pbIns
}

const sqlStrToPraisePart1 = "update tb_user set ur_p_times = ur_p_times + 1 where ur_id = ?"
const sqlStrToPraisePart2 = "update tb_cmt set cmt_p_times = cmt_p_times + 1 where cmt_id = ?"
const sqlStrToPraisePart3 = "update tb_pb set pb_p_value = 1, pb_p_time = ? where pb_id = ?"

//Praise 赞
func (s *sqlBase) Praise(pb *models.PBInDB) {
	if pb.PValue == 1 {
		return
	}
	tx, _ := s.Begin()
	checkSqlResultErr(tx.Exec(sqlStrToPraisePart1, pb.UserID))
	checkSqlResultErr(tx.Exec(sqlStrToPraisePart2, pb.CmtID))
	checkSqlResultErr(tx.Exec(sqlStrToPraisePart3, pb.PTime, pb.ID))
	checkErr(tx.Commit())
}

const sqlStrToBelittlePart1 = "update tb_user set ur_b_times = ur_b_times + 1 where ur_id = ?"
const sqlStrToBelittlePart2 = "update tb_cmt set cmt_b_times = cmt_b_times + 1 where cmt_id = ?"
const sqlStrToBelittlePart3 = "update tb_pb set pb_b_value = 1, pb_b_time = ? where pb_id = ?"

//Belittle 贬
func (s *sqlBase) Belittle(pb *models.PBInDB) {
	if pb.BValue == 1 {
		return
	}
	tx, _ := s.Begin()
	checkSqlResultErr(tx.Exec(sqlStrToBelittlePart1, pb.UserID))
	checkSqlResultErr(tx.Exec(sqlStrToBelittlePart2, pb.CmtID))
	checkSqlResultErr(tx.Exec(sqlStrToBelittlePart3, pb.BTime, pb.ID))
	checkErr(tx.Commit())
}

const sqlStrToPraiseCancelPart1 = "update tb_user set ur_p_times = ur_p_times - 1 where ur_id = ?"
const sqlStrToPraiseCancelPart2 = "update tb_cmt set cmt_p_times = cmt_p_times - 1 where cmt_id = ?"
const sqlStrToPraiseCancelPart3 = "update tb_pb set pb_p_value = 0, pb_p_ctime = ? where pb_id = ?"

//PraiseCancel 取消赞
func (s *sqlBase) PraiseCancel(pb *models.PBInDB) {
	if pb.PValue == 0 {
		return
	}
	tx, _ := s.Begin()
	checkSqlResultErr(tx.Exec(sqlStrToPraiseCancelPart1, pb.UserID))
	checkSqlResultErr(tx.Exec(sqlStrToPraiseCancelPart2, pb.CmtID))
	checkSqlResultErr(tx.Exec(sqlStrToPraiseCancelPart3, pb.PCTime, pb.ID))
	checkErr(tx.Commit())
}

const sqlStrToBelittleCancelPart1 = "update tb_user set ur_b_times = ur_b_times - 1 where ur_id = ?"
const sqlStrToBelittleCancelPart2 = "update tb_cmt set cmt_b_times = cmt_b_times - 1 where cmt_id = ?"
const sqlStrToBelittleCancelPart3 = "update tb_pb set pb_b_value = 0, pb_b_ctime = ? where pb_id = ?"

//BelittleCancel 取消贬
func (s *sqlBase) BelittleCancel(pb *models.PBInDB) {
	if pb.BValue == 0 {
		return
	}
	tx, _ := s.Begin()
	checkSqlResultErr(tx.Exec(sqlStrToBelittleCancelPart1, pb.UserID))
	checkSqlResultErr(tx.Exec(sqlStrToBelittleCancelPart2, pb.CmtID))
	checkSqlResultErr(tx.Exec(sqlStrToBelittleCancelPart3, pb.BCTime, pb.ID))
	checkErr(tx.Commit())
}

//SetPB 设置赞踩
func (s *sqlBase) SetPB(cmtID int, userID int, isP bool, isD bool) {
	//先尝试查询赞踩
	pbIns := s.QueryPbItem(cmtID, userID)
	if pbIns == nil {
		pbIns = new(models.PBInDB)
		pbIns.CmtID = cmtID
		pbIns.UserID = userID
		s.AddPbItem(pbIns)
	}
	unixTimeCount := time.Now().UnixNano()
	//现在需要分情况考虑
	if isP {
		if isD {
			//赞
			pbIns.PTime = unixTimeCount
			s.Praise(pbIns)
			return
		}
		//取消赞
		pbIns.PCTime = unixTimeCount
		s.PraiseCancel(pbIns)
		return
	}
	if isD {
		//贬
		pbIns.BTime = unixTimeCount
		s.Belittle(pbIns)
		return
	}
	//取消贬
	pbIns.BCTime = unixTimeCount
	s.BelittleCancel(pbIns)
}

const sqlStrToAddUser = `
insert into tb_user 
(
	ur_account,
	ur_pwd,
	ur_name,
	ur_hp_path,
	ur_type,
	ur_state,
	ur_su_time,
	ur_post_count,
	ur_cmt_count,
	ur_img_count,
	ur_p_times,
	ur_b_times,
	ur_le_time
)
values (?,?,?,?,?,?,?,?,?,?,?,?,?)`

//AddUser 新增用户
func (s *sqlBase) AddUser(user *models.UserInDB) {
	back, err := s.Exec(sqlStrToAddUser,
		user.Account,
		s.passwordMd5ToHexStr(user.PassWord),
		user.Name,
		user.HeadPhotoPath,
		user.Type,
		user.State,
		user.SignUpTime,
		user.PostCount,
		user.CommentCount,
		user.ImageCount,
		user.PraiseTimes,
		user.BelittleTimes,
		user.LastEditTime)
	checkErr(err)
	uid, err := back.LastInsertId()
	checkErr(err)
	user.ID = int(uid)
	check1Err(back.RowsAffected())
}

const sqlStrToDeleteUser = "delete from tb_user where ur_id =?"

//DeleteUser 删除用户，不更新任何计数器
func (s *sqlBase) DeleteUser(userID int) {
	checkSqlResultErr(s.Exec(sqlStrToDeleteUser, userID))
}

const sqlStrToQueryUserByID = "select * from tb_user where ur_id = ?"

//QueryUserByID 查询用户
func (s *sqlBase) QueryUserByID(userID int) *models.UserInDB {
	return s.ScanFromUserTbAllFields(s.QueryRow(sqlStrToQueryUserByID, userID))
}

const sqlStrToQueryUserByAccountAndPwd = "select * from tb_user where ur_account = ? and ur_pwd = ?"

//QueryUserByAccountAndPwd 查询用户
func (s *sqlBase) QueryUserByAccountAndPwd(account string, password string) *models.UserInDB {
	row := s.QueryRow(sqlStrToQueryUserByAccountAndPwd, account, s.passwordMd5ToHexStr(password))
	return s.ScanFromUserTbAllFields(row)
}

func (s *sqlBase) ScanFromUserTbAllFields(scanner Scanner) *models.UserInDB {
	user := new(models.UserInDB)
	if scanner.Scan(
		&user.ID,
		&user.Account,
		&user.PassWord,
		&user.Name,
		&user.HeadPhotoPath,
		&user.Type,
		&user.State,
		&user.SignUpTime,
		&user.PostCount,
		&user.CommentCount,
		&user.ImageCount,
		&user.PraiseTimes,
		&user.BelittleTimes,
		&user.LastEditTime) != nil {
		return nil
	}
	return user
}

const sqlStrToQueryPostCountOfUser = "select ur_post_count from tb_user where ur_id = ?;"

//QueryPostCountOfUser 统计用户发帖总量
func (s *sqlBase) QueryPostCountOfUser(userID int) int {
	var count int
	checkErr(s.QueryRow(sqlStrToQueryPostCountOfUser, userID).Scan(&count))
	return count
}

const sqlStrToQueryImageCountOfUser = "select ur_img_count from tb_user where ur_id = ?;"

//QueryImageCountOfUser 统计用户上传图片的量
func (s *sqlBase) QueryImageCountOfUser(userID int) int {
	var count int
	checkErr(s.QueryRow(sqlStrToQueryImageCountOfUser, userID).Scan(&count))
	return count
}

const sqlStrToUpdateUserHeadPhoto = "update tb_user set ur_hp_path = ? where ur_id = ?"

//UpdateUserHeadPhoto 更新头像文件路径
func (s *sqlBase) UpdateUserHeadPhoto(userID int, path string) {
	checkSqlResultErr(s.Exec(sqlStrToUpdateUserHeadPhoto, path, userID))
}

const sqlStrToIsUserNameExist = "select ur_id from tb_user where ur_name = ?"

//IsUserNameExist 是否昵称已存在
func (s *sqlBase) IsUserNameExist(name string) bool {
	row := s.QueryRow(sqlStrToIsUserNameExist, name)
	return row.Scan(new(int64)) == nil
}

const sqlStrToIsUserAccountExist = "select ur_id from tb_user where ur_account = ?"

//IsUserAccountExist 是否账号已存在
func (s *sqlBase) IsUserAccountExist(account string) bool {
	row := s.QueryRow(sqlStrToIsUserAccountExist, account)
	return row.Scan(new(int)) == nil
}

//把密码md计算成16进制字符串 长度36
func (s *sqlBase) passwordMd5ToHexStr(password string) string {
	buffer := md5.New().Sum([]byte(password))
	if len(buffer) > 18 {
		buffer = buffer[0:18]
	}
	md5Str := hex.EncodeToString(buffer)
	return md5Str
}

const sqlStrToAddImagesPart1 = "insert into tb_img (img_user_id,img_upload_time,img_path) values (?,?,?);"
const sqlStrToAddImagesPart2 = "update tb_user set ur_img_count = ur_img_count + 1 where ur_id = ?;"

//批量新增图片
func (s *sqlBase) AddImages(images []*models.ImageInDB) {
	tx, _ := s.DB.Begin()
	stmt1, _ := tx.Prepare(sqlStrToAddImagesPart1)
	stmt2, _ := tx.Prepare(sqlStrToAddImagesPart2)
	for _, v := range images {
		result, err := stmt1.Exec(v.UserID, v.UploadTime, v.FilePath)
		checkErr(err)
		newID, err := result.LastInsertId()
		checkErr(err)
		v.ID = int(newID)
		check1Err(result.RowsAffected())
		//更新用户上传图片计数器
		checkSqlResultErr(stmt2.Exec(v.UserID))
	}
	checkErr(tx.Commit())
}

const sqlStrToQueryImages = `select tb_img.img_path
from tb_img,
(select img_id from tb_img where img_user_id = ? order by img_id desc limit ? offset ?) as imgId
where tb_img.img_id = imgId.img_id`

//图片查询
func (s *sqlBase) QueryImages(userID int, count int, offset int) []string {
	rows, err := s.Query(sqlStrToQueryImages, userID, count, offset)
	checkErr(err)
	paths := make([]string, 0, count)
	for rows.Next() {
		var path string
		checkErr(rows.Scan(&path))
		paths = append(paths, path)
	}
	return paths
}
