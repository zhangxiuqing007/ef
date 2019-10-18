package dba

import (
	"ef/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

type sqlBase struct {
	*sql.DB
}

//Close 关闭
func (s *sqlBase) Close() error {
	return s.DB.Close()
}

const sqlStrToAddTheme = "insert into tb_theme (tm_name,tm_post_count) values (?,?)"

//AddTheme 增加主题
func (s *sqlBase) AddTheme(theme *models.ThemeInDB) error {
	result, err := s.Exec(sqlStrToAddTheme, theme.Name, theme.PostCount)
	if err != nil {
		return err
	}
	tmID, err := result.LastInsertId()
	theme.ID = int(tmID)
	return err
}

const sqlStrToDeleteTheme = "delete from tb_theme where tm_id = ?"

//DeleteTheme 删除主题(不处理下面的帖子和评论，不影响计数器)
func (s *sqlBase) DeleteTheme(themeID int) error {
	_, err := s.Exec(sqlStrToDeleteTheme, themeID)
	return err
}

const sqlStrToUpdateTheme = "update tb_theme set tm_name = ? where tm_id = ?"

//UpdateTheme 更新主题名称
func (s *sqlBase) UpdateTheme(theme *models.ThemeInDB) error {
	_, err := s.Exec(sqlStrToUpdateTheme, theme.Name, theme.ID)
	return err
}

const sqlStrToQueryTheme = "select tm_name,tm_post_count from tb_theme where tm_id = ?"

//QueryTheme 查询某个主题
func (s *sqlBase) QueryTheme(themeID int) (*models.ThemeInDB, error) {
	tm := &models.ThemeInDB{ID: themeID}
	err := s.QueryRow(sqlStrToQueryTheme, themeID).Scan(&tm.Name, &tm.PostCount)
	return tm, err
}

const sqlStrToQueryThemes = "select * from tb_theme order by tm_id"

//QueryAllThemes 查询主题列表
func (s *sqlBase) QueryAllThemes() ([]*models.ThemeInDB, error) {
	rows, err := s.Query(sqlStrToQueryThemes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tms := make([]*models.ThemeInDB, 0, 20)
	for rows.Next() {
		tm := new(models.ThemeInDB)
		err = rows.Scan(&tm.ID, &tm.Name, &tm.PostCount)
		if err != nil {
			return tms, err
		}
		tms = append(tms, tm)
	}
	return tms, err
}

const sqlStrToQueryPostCountOfTheme = "select tm_post_count from tb_theme where tm_id = ?"

//QueryPostCountOfTheme 查询目标主题的帖子数量
func (s *sqlBase) QueryPostCountOfTheme(themeID int) (int, error) {
	var count int
	err := s.QueryRow(sqlStrToQueryPostCountOfTheme, themeID).Scan(&count)
	return count, err
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
func (s *sqlBase) AddPost(post *models.PostInDB, cmt *models.CommentInDB) error {
	tx, _ := s.Begin()
	if result, err := s.Exec(sqlStrToAddPostPart3, post.ThemeID, post.UserID, post.Title, post.State,
		post.CreatedTime, post.CmtCount, post.LastCmterID, post.LastCmtTime); err == nil {
		if postID, err := result.LastInsertId(); err == nil {
			post.ID = int(postID)
			cmt.PostID = post.ID
		} else {
			tx.Rollback()
			return err
		}
	} else {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToAddPostPart1, post.UserID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToAddPostPart2, post.ThemeID); err != nil {
		tx.Rollback()
		return err
	}
	if result, err := tx.Exec(sqlStrToAddCommentPart1, cmt.PostID, cmt.UserID, cmt.Content, cmt.State,
		cmt.CreatedTime, cmt.LastEditTime, cmt.EditTimes, cmt.PraiseTimes, cmt.BelittleTimes); err == nil {
		if cmtID, err := result.LastInsertId(); err == nil {
			cmt.ID = int(cmtID)
		} else {
			tx.Rollback()
			return err
		}
	} else {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToAddCommentPart3, cmt.LastEditTime, cmt.UserID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

//AddPosts 批量新增帖子
func (s *sqlBase) AddPosts(posts []*models.PostInDB, cmts []*models.CommentInDB) error {
	tx, _ := s.Begin()
	stmt1, _ := tx.Prepare(sqlStrToAddPostPart1)
	stmt2, _ := tx.Prepare(sqlStrToAddPostPart2)
	stmt3, _ := tx.Prepare(sqlStrToAddPostPart3)
	stmt4, _ := tx.Prepare(sqlStrToAddCommentPart1)
	stmt5, _ := tx.Prepare(sqlStrToAddCommentPart3)
	for i, v := range posts {
		cmt := cmts[i]
		if _, err := stmt1.Exec(v.UserID); err != nil {
			tx.Rollback()
			return err
		}
		if _, err := stmt2.Exec(v.ThemeID); err != nil {
			tx.Rollback()
			return err
		}
		if result, err := stmt3.Exec(v.ThemeID, v.UserID, v.Title, v.State, v.CreatedTime, v.CmtCount, v.LastCmterID, v.LastCmtTime); err == nil {
			if pid, err := result.LastInsertId(); err == nil {
				v.ID = int(pid)
				cmt.PostID = v.ID
			} else {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return err
		}
		if result, err := stmt4.Exec(cmt.PostID, cmt.UserID, cmt.Content, cmt.State, cmt.CreatedTime, cmt.LastEditTime, cmt.EditTimes, cmt.PraiseTimes, cmt.BelittleTimes); err == nil {
			if cmtID, err := result.LastInsertId(); err == nil {
				cmt.ID = int(cmtID)
			} else {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return err
		}
		if _, err := stmt5.Exec(cmt.LastEditTime, cmt.UserID); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

const sqlStrToDeletePost = "delete from tb_post where po_id = ?;"

//DeletePost 删除帖子（但不影响任何计数器）
func (s *sqlBase) DeletePost(postID int) error {
	_, err := s.Exec(sqlStrToDeletePost, postID)
	return err
}

const sqlStrToQueryPost = `
select
    po_theme_id,
    po_user_id,
    po_title,
    po_state,
    po_post_time,
    po_cmt_count,
    po_lcmter_id,
    po_lcmt_time
from tb_post where po_id = ?`

//QueryPost 查询帖子
func (s *sqlBase) QueryPost(postID int) (*models.PostInDB, error) {
	post := new(models.PostInDB)
	post.ID = postID
	err := s.QueryRow(sqlStrToQueryPost, postID).Scan(
		&post.ThemeID,
		&post.UserID,
		&post.Title,
		&post.State,
		&post.CreatedTime,
		&post.CmtCount,
		&post.LastCmterID,
		&post.LastCmtTime)
	return post, err
}

const sqlStrToQueryPostTitle = "select po_title from tb_post where po_id = ?;"

//QueryPostTitle 查询主题
func (s *sqlBase) QueryPostTitle(postID int) (string, error) {
	var title string
	err := s.QueryRow(sqlStrToQueryPostTitle, postID).Scan(&title)
	return title, err
}

const sqlStrToUpdatePostTitle = "update tb_post set po_title = ? where po_id = ?"

//UpdatePostTitle 修改帖子标题
func (s *sqlBase) UpdatePostTitle(post *models.PostInDB) error {
	tx, _ := s.Begin()
	if _, err := tx.Exec(sqlStrToUpdatePostTitle, post.Title, post.ID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToAddCommentPart3, post.LastCmtTime, post.UserID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
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
func (s *sqlBase) QueryPostsOfTheme(themeID int, count, offset, sortType int) ([]*models.PostOnThemePage, error) {
	var rows *sql.Rows
	var err error
	if sortType == 0 {
		rows, err = s.Query(sqlStrToQueryPostsSortType0, themeID, count, offset)
	} else {
		rows, err = s.Query(sqlStrToQueryPostsSortType1, themeID, count, offset)
	}
	if err != nil {
		return nil, err
	}
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
func (s *sqlBase) QueryPostsOfUser(userID int, count, offset int) ([]*models.PostOnThemePage, error) {
	rows, err := s.Query(sqlStrToQueryPostsOfUser, userID, count, offset)
	if err != nil {
		return nil, err
	}
	return turnToPostsOnThemePage(rows, count)
}

func turnToPostsOnThemePage(rows *sql.Rows, cap int) ([]*models.PostOnThemePage, error) {
	defer rows.Close()
	posts := make([]*models.PostOnThemePage, 0, cap)
	for rows.Next() {
		post := new(models.PostOnThemePage)
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.CmtCount,
			&post.CreaterID,
			&post.CreaterName,
			&post.CreatedTime,
			&post.LastCmterID,
			&post.LastCmterName,
			&post.LastCmtTime)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, nil
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
func (s *sqlBase) QueryPostOfPostPage(postID int) (*models.PostOnPostPage, error) {
	post := new(models.PostOnPostPage)
	post.ID = postID
	err := s.QueryRow(sqlStrToQueryPostOfPostPage, postID).Scan(
		&post.Title,
		&post.CreaterID,
		&post.CmtCount,
		&post.ThemeID,
		&post.ThemeName)
	return post, err
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
func (s *sqlBase) AddComment(cmt *models.CommentInDB) error {
	tx, _ := s.Begin()
	result, err := s.Exec(sqlStrToAddCommentPart1, cmt.PostID, cmt.UserID, cmt.Content, cmt.State,
		cmt.CreatedTime, cmt.LastEditTime, cmt.EditTimes, cmt.PraiseTimes, cmt.BelittleTimes)
	if err != nil {
		tx.Rollback()
		return err
	}
	cmtID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	cmt.ID = int(cmtID)
	if _, err := tx.Exec(sqlStrToAddCommentPart2, cmt.UserID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToAddCommentPart3, cmt.LastEditTime, cmt.UserID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToAddCommentPart4, cmt.PostID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToAddCommentPart5, cmt.UserID, cmt.PostID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToAddCommentPart6, cmt.LastEditTime, cmt.PostID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

//AddComments 批量增加评论
func (s *sqlBase) AddComments(comments []*models.CommentInDB) error {
	tx, _ := s.Begin()
	stmt1, _ := tx.Prepare(sqlStrToAddCommentPart1)
	stmt2, _ := tx.Prepare(sqlStrToAddCommentPart2)
	stmt3, _ := tx.Prepare(sqlStrToAddCommentPart3)
	stmt4, _ := tx.Prepare(sqlStrToAddCommentPart4)
	stmt5, _ := tx.Prepare(sqlStrToAddCommentPart5)
	stmt6, _ := tx.Prepare(sqlStrToAddCommentPart6)
	for _, v := range comments {
		result, err := stmt1.Exec(v.PostID, v.UserID, v.Content, v.State, v.CreatedTime,
			v.LastEditTime, v.EditTimes, v.PraiseTimes, v.BelittleTimes)
		if err != nil {
			tx.Rollback()
			return err
		}
		cid, err := result.LastInsertId()
		if err != nil {
			tx.Rollback()
			return err
		}
		v.ID = int(cid)
		if _, err := stmt2.Exec(v.UserID); err != nil {
			tx.Rollback()
			return err
		}
		if _, err := stmt3.Exec(v.LastEditTime, v.UserID); err != nil {
			tx.Rollback()
			return err
		}
		if _, err := stmt4.Exec(v.PostID); err != nil {
			tx.Rollback()
			return err
		}
		if _, err := stmt5.Exec(v.UserID, v.PostID); err != nil {
			tx.Rollback()
			return err
		}
		if _, err := stmt6.Exec(v.LastEditTime, v.PostID); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

const sqlStrToDeleteComment = "delete from tb_cmt where cmt_id =?;"

//DeleteComment 删除单个评论，不处理任何计数器
func (s *sqlBase) DeleteComment(cmtID int) error {
	_, err := s.Exec(sqlStrToDeleteComment, cmtID)
	return err
}

const sqlStrToQueryComment = "select * from tb_cmt where cmt_id = ?"

//QueryComment 查询单个评论
func (s *sqlBase) QueryComment(cmtID int) (*models.CommentInDB, error) {
	cmt := new(models.CommentInDB)
	err := s.QueryRow(sqlStrToQueryComment, cmtID).Scan(
		&cmt.ID,
		&cmt.PostID,
		&cmt.UserID,
		&cmt.Content,
		&cmt.State,
		&cmt.CreatedTime,
		&cmt.LastEditTime,
		&cmt.EditTimes,
		&cmt.PraiseTimes,
		&cmt.BelittleTimes)
	return cmt, err
}

const sqlStrToUpdateCommentPart1 = "update tb_cmt set cmt_content = ?, cmt_le_time = ?, cmt_e_times= ? where cmt_id = ? "
const sqlStrToUpdateCommentPart2 = "update tb_user set ur_le_time = ? where ur_id = ? "

//UpdateComment 修改评论内容
func (s *sqlBase) UpdateComment(cmt *models.CommentInDB) error {
	tx, _ := s.Begin()
	if _, err := tx.Exec(sqlStrToUpdateCommentPart1, cmt.Content, cmt.LastEditTime, cmt.EditTimes, cmt.ID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToUpdateCommentPart2, cmt.LastEditTime, cmt.PostID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

const sqlStrToQueryComments = "select * from tb_cmt where cmt_post_id = ? order by cmt_id"

//QueryComments 查询评论，按照创建时间排序
func (s *sqlBase) QueryComments(postID int) ([]*models.CommentInDB, error) {
	rows, err := s.Query(sqlStrToQueryComments, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cmts := make([]*models.CommentInDB, 0, 32)
	for rows.Next() {
		cmt := new(models.CommentInDB)
		err = rows.Scan(
			&cmt.ID,
			&cmt.PostID,
			&cmt.UserID,
			&cmt.Content,
			&cmt.State,
			&cmt.CreatedTime,
			&cmt.LastEditTime,
			&cmt.EditTimes,
			&cmt.PraiseTimes,
			&cmt.BelittleTimes)
		if err != nil {
			return nil, err
		}
		cmts = append(cmts, cmt)
	}
	return cmts, nil
}

const sqlStrToQueryCommentsOfPostPagePart1 = `
select 
       cmt.cmt_id,
	   cmt.cmt_content,
       u.ur_id,
	   u.ur_name,
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
func (s *sqlBase) QueryCommentsOfPostPage(postID int, count int, offset int, userID int) ([]*models.CmtOnPostPage, error) {
	cmts := make([]*models.CmtOnPostPage, 0, count)
	rows, err := s.Query(sqlStrToQueryCommentsOfPostPagePart1, postID, count, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		cmt := new(models.CmtOnPostPage)
		err = rows.Scan(
			&cmt.ID,
			&cmt.Content,
			&cmt.CmterID,
			&cmt.CmterName,
			&cmt.CmtTime,
			&cmt.CmtEditTimes,
			&cmt.LastEditTime,
			&cmt.PraiseTimes,
			&cmt.BelittleTimes)
		if err != nil {
			return nil, err
		}
		cmts = append(cmts, cmt)
	}
	//查看指定用户的赞踩情况
	rowspb, err := s.Query(sqlStrToQueryCommentsOfPostPagePart2, postID, count, offset, userID)
	if err != nil {
		return nil, err
	}
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
		err = rowspb.Scan(&tempID, &tempPvalue, &tempBvalue)
		if err != nil {
			return nil, err
		}
		targetCmt := f(tempID)
		if targetCmt != nil {
			targetCmt.IsPraised = tempPvalue == 1
			targetCmt.IsBelittled = tempBvalue == 1
		}
	}
	return cmts, nil
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
func (s *sqlBase) AddPbItem(pb *models.PBInDB) error {
	result, err := s.Exec(sqlStrToAddPbItem,
		pb.CmtID, pb.UserID,
		pb.PValue, pb.PTime, pb.PCTime,
		pb.BValue, pb.BTime, pb.BCTime)
	if err != nil {
		return err
	}
	pbID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	pb.ID = int(pbID)
	return nil
}

const sqlStrToQueryPbItem = `
select
	pb_id,
	pb_p_value,
	pb_p_time,
	pb_p_ctime,
	pb_b_value,
	pb_b_time,
	pb_b_ctime 
from tb_pb where pb_cmt_id = ? and pb_user_id = ?;`

//QueryPbItem 查询赞踩行
func (s *sqlBase) QueryPbItem(cmtID int, userID int) (*models.PBInDB, error) {
	var pbIns = new(models.PBInDB)
	pbIns.CmtID = cmtID
	pbIns.UserID = userID
	err := s.QueryRow(sqlStrToQueryPbItem, cmtID, userID).Scan(
		&pbIns.ID,
		&pbIns.PValue,
		&pbIns.PTime,
		&pbIns.PCTime,
		&pbIns.BValue,
		&pbIns.BTime,
		&pbIns.BCTime)
	return pbIns, err
}

const sqlStrToPraisePart1 = "update tb_user set ur_p_times = ur_p_times + 1 where ur_id = ?"
const sqlStrToPraisePart2 = "update tb_cmt set cmt_p_times = cmt_p_times + 1 where cmt_id = ?"
const sqlStrToPraisePart3 = "update tb_pb set pb_p_value = 1, pb_p_time = ? where pb_id = ?"

//Praise 赞
func (s *sqlBase) Praise(pb *models.PBInDB) error {
	if pb.PValue == 1 {
		return nil
	}
	tx, _ := s.Begin()
	if _, err := tx.Exec(sqlStrToPraisePart1, pb.UserID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToPraisePart2, pb.CmtID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToPraisePart3, pb.PTime, pb.ID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

const sqlStrToBelittlePart1 = "update tb_user set ur_b_times = ur_b_times + 1 where ur_id = ?"
const sqlStrToBelittlePart2 = "update tb_cmt set cmt_b_times = cmt_b_times + 1 where cmt_id = ?"
const sqlStrToBelittlePart3 = "update tb_pb set pb_b_value = 1, pb_b_time = ? where pb_id = ?"

//Belittle 贬
func (s *sqlBase) Belittle(pb *models.PBInDB) error {
	if pb.BValue == 1 {
		return nil
	}
	tx, _ := s.Begin()
	if _, err := tx.Exec(sqlStrToBelittlePart1, pb.UserID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToBelittlePart2, pb.CmtID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToBelittlePart3, pb.BTime, pb.ID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

const sqlStrToPraiseCancelPart1 = "update tb_user set ur_p_times = ur_p_times - 1 where ur_id = ?"
const sqlStrToPraiseCancelPart2 = "update tb_cmt set cmt_p_times = cmt_p_times - 1 where cmt_id = ?"
const sqlStrToPraiseCancelPart3 = "update tb_pb set pb_p_value = 0, pb_p_ctime = ? where pb_id = ?"

//PraiseCancel 取消赞
func (s *sqlBase) PraiseCancel(pb *models.PBInDB) error {
	if pb.PValue == 0 {
		return nil
	}
	tx, _ := s.Begin()
	if _, err := tx.Exec(sqlStrToPraiseCancelPart1, pb.UserID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToPraiseCancelPart2, pb.CmtID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToPraiseCancelPart3, pb.PCTime, pb.ID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

const sqlStrToBelittleCancelPart1 = "update tb_user set ur_b_times = ur_b_times - 1 where ur_id = ?"
const sqlStrToBelittleCancelPart2 = "update tb_cmt set cmt_b_times = cmt_b_times - 1 where cmt_id = ?"
const sqlStrToBelittleCancelPart3 = "update tb_pb set pb_b_value = 0, pb_b_ctime = ? where pb_id = ?"

//BelittleCancel 取消贬
func (s *sqlBase) BelittleCancel(pb *models.PBInDB) error {
	if pb.BValue == 0 {
		return nil
	}
	tx, _ := s.Begin()
	if _, err := tx.Exec(sqlStrToBelittleCancelPart1, pb.UserID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToBelittleCancelPart2, pb.CmtID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(sqlStrToBelittleCancelPart3, pb.BCTime, pb.ID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

//SetPB 设置赞踩
func (s *sqlBase) SetPB(cmtID int, userID int, isP bool, isD bool) error {
	//先尝试查询赞踩
	pbIns, err := s.QueryPbItem(cmtID, userID)
	//如果还没有记录
	if err != nil || pbIns.ID == 0 {
		err = s.AddPbItem(pbIns)
		//如果新增失败的话
		if err != nil {
			return err
		}
	}
	unixTimeCount := time.Now().UnixNano()
	//现在需要分情况考虑
	if isP {
		//赞
		if isD {
			pbIns.PTime = unixTimeCount
			return s.Praise(pbIns)
			//取消赞
		}
		pbIns.PCTime = unixTimeCount
		return s.PraiseCancel(pbIns)
	}
	//贬
	if isD {
		pbIns.BTime = unixTimeCount
		return s.Belittle(pbIns)
		//取消贬
	}
	pbIns.BCTime = unixTimeCount
	return s.BelittleCancel(pbIns)
}

const sqlStrToAddUser = `
insert into tb_user 
(
	ur_account,
	ur_pwd,
	ur_name,
	ur_type,
	ur_state,
	ur_su_time,
	ur_post_count,
	ur_cmt_count,
	ur_p_times,
	ur_b_times,
	ur_le_time
)
values (?,?,?,?,?,?,?,?,?,?,?)`

//AddUser 新增用户
func (s *sqlBase) AddUser(user *models.UserInDB) error {
	back, err := s.Exec(sqlStrToAddUser,
		user.Account,
		s.passwordMd5ToHexStr(user.PassWord),
		user.Name,
		user.Type,
		user.State,
		user.SignUpTime,
		user.PostCount,
		user.CommentCount,
		user.PraiseTimes,
		user.BelittleTimes,
		user.LastEditTime)
	if err != nil {
		return err
	}
	urid, err := back.LastInsertId()
	user.ID = int(urid)
	return err
}

const sqlStrToDeleteUser = "delete from tb_user where ur_id =?"

//DeleteUser 删除用户，不更新任何计数器
func (s *sqlBase) DeleteUser(userID int) error {
	_, err := s.Exec(sqlStrToDeleteUser, userID)
	return err
}

const sqlStrToQueryUserByID = "select * from tb_user where ur_id = ?"

//QueryUserByID 查询用户
func (s *sqlBase) QueryUserByID(userID int) (*models.UserInDB, error) {
	user := new(models.UserInDB)
	err := s.QueryRow(sqlStrToQueryUserByID, userID).Scan(
		&user.ID,
		&user.Account,
		&user.PassWord,
		&user.Name,
		&user.Type,
		&user.State,
		&user.SignUpTime,
		&user.PostCount,
		&user.CommentCount,
		&user.PraiseTimes,
		&user.BelittleTimes,
		&user.LastEditTime)
	if err != nil {
		err = errors.New("无此用户")
		return nil, err
	}
	return user, nil
}

const sqlStrToQueryUserByAccountAndPwd = "select * from tb_user where ur_account = ? and ur_pwd = ?"

//QueryUserByAccountAndPwd 查询用户
func (s *sqlBase) QueryUserByAccountAndPwd(account string, password string) (*models.UserInDB, error) {
	user := new(models.UserInDB)
	err := s.QueryRow(sqlStrToQueryUserByAccountAndPwd, account, s.passwordMd5ToHexStr(password)).Scan(
		&user.ID,
		&user.Account,
		&user.PassWord,
		&user.Name,
		&user.Type,
		&user.State,
		&user.SignUpTime,
		&user.PostCount,
		&user.CommentCount,
		&user.PraiseTimes,
		&user.BelittleTimes,
		&user.LastEditTime)
	if err != nil {
		err = errors.New("无此用户")
		return nil, err
	}
	return user, nil
}

const sqlStrToQueryPostCountOfUser = "select ur_post_count from tb_user where ur_id = ?;"

//QueryPostCountOfUser 统计用户发帖总量
func (s *sqlBase) QueryPostCountOfUser(userID int) (int, error) {
	var count int
	err := s.QueryRow(sqlStrToQueryPostCountOfUser, userID).Scan(&count)
	return count, err
}

const sqlStrToIsUserNameExist = "select ur_id from tb_user where ur_name = ?"

//IsUserNameExist 是否昵称已存在
func (s *sqlBase) IsUserNameExist(name string) bool {
	row := s.QueryRow(sqlStrToIsUserNameExist, name)
	err := row.Scan(new(int64))
	return err == nil
}

const sqlStrToIsUserAccountExist = "select ur_id from tb_user where ur_account = ?"

//IsUserAccountExist 是否账号已存在
func (s *sqlBase) IsUserAccountExist(account string) bool {
	row := s.QueryRow(sqlStrToIsUserAccountExist, account)
	err := row.Scan(new(int))
	return err == nil
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
