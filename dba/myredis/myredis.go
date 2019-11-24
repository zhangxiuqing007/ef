package myredis

import (
	"ef/dba"
	"ef/models"
	"errors"
	"strconv"

	"github.com/go-redis/redis"
)

var CacheInitPostCount int = 1_0000 //1W
//var CacheMaxPostCount int = 10_0000 //10W
const redisURL = "127.0.0.1:6379"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func check1Err(_ interface{}, err error) {
	checkErr(err)
}

type MySQLRedisIns struct {
	*dba.MySQLIns
	rc *redis.Client
}

func (s *MySQLRedisIns) Open(dbFilePath string) {
	s.rc = redis.NewClient(&redis.Options{
		Addr: redisURL,
		DB:   1,
	})
	s.MySQLIns = new(dba.MySQLIns)
	s.MySQLIns.Open(dbFilePath)
}

func (s *MySQLRedisIns) Close() {
	s.MySQLIns.Close()
	checkErr(s.rc.Close())
}

//redis数据的初始化
func (s *MySQLRedisIns) CopyAllTypeDataFromMysqlToRedis() {
	s.copyThemesAllFromMysqlToRedis()
	s.copyUsersAllFromMysqlToRedis()

	s.copyPostDataFromMysqlToRedis()
	s.copyCmtDataFromMysqlToRedis()
	s.copyPbDataFromMysqlToRedis()
}

//最外层key前缀
const redisTopPreHashKeyTheme = "tm:"
const redisTopPreHashKeyUser = "ur:"
const redisTopPreHashKeyPost = "po:"
const redisTopPreHashKeyCmt = "cmt:"
const redisTopPreHashKeyPb = "pb:"

//附加查询最外层前缀
const redisTopPreHashKeyUserAccount = "uac:"

//hash字段 theme
const redisFieldThemeID = "ID"
const redisFieldThemeName = "NA"
const redisFieldThemePostCount = "PC"

//theme section
func (s *MySQLRedisIns) copyThemesAllFromMysqlToRedis() {
	pipe := s.rc.Pipeline()
	defer func() {
		checkErr(pipe.Close())
	}()
	themeIDsInRedis := s.getThemeIDsNowInRedis()
	tms := s.MySQLIns.QueryAllThemes()
	for _, v := range tms {
		if _, ok := themeIDsInRedis[v.ID]; ok {
			continue
		}
		key := redisTopPreHashKeyTheme + strconv.Itoa(v.ID)
		pipe.HSet(key, redisFieldThemeID, v.ID)
		pipe.HSet(key, redisFieldThemeName, v.Name)
		pipe.HSet(key, redisFieldThemePostCount, v.PostCount)
	}
	check1Err(pipe.Exec())
}

const redisFieldUserID = "ID"
const redisFieldUserAccount = "AC"
const redisFieldUserPassWord = "PW"
const redisFieldUserName = "NA"
const redisFieldUserHeadPhotoPath = "HE"
const redisFieldUserType = "TY"
const redisFieldUserState = "ST"
const redisFieldUserSignUpTime = "SI"
const redisFieldUserPostCount = "PC"
const redisFieldUserCommentCount = "CC"
const redisFieldUserImageCount = "IC"
const redisFieldUserPraiseTimes = "PT"
const redisFieldUserBelittleTimes = "BT"
const redisFieldUserLastEditTime = "LE"

const redisFieldUserAccountPwd = "pwd"

//user section
func (s *MySQLRedisIns) copyUsersAllFromMysqlToRedis() {
	pipe := s.rc.Pipeline()
	defer func() {
		checkErr(pipe.Close())
	}()
	userIDsInRedis := s.getUserIDsNowInRedis()
	users := s.queryAllUsers()
	for _, v := range users {
		if _, ok := userIDsInRedis[v.ID]; ok {
			continue
		}
		key := redisTopPreHashKeyUser + strconv.Itoa(v.ID)
		pipe.HSet(key, redisFieldUserID, v.ID)
		pipe.HSet(key, redisFieldUserAccount, v.Account)
		pipe.HSet(key, redisFieldUserPassWord, v.PassWord)
		pipe.HSet(key, redisFieldUserName, v.Name)
		pipe.HSet(key, redisFieldUserHeadPhotoPath, v.HeadPhotoPath)
		pipe.HSet(key, redisFieldUserType, v.Type)
		pipe.HSet(key, redisFieldUserState, v.State)
		pipe.HSet(key, redisFieldUserSignUpTime, v.SignUpTime)
		pipe.HSet(key, redisFieldUserPostCount, v.PostCount)
		pipe.HSet(key, redisFieldUserCommentCount, v.CommentCount)
		pipe.HSet(key, redisFieldUserImageCount, v.ImageCount)
		pipe.HSet(key, redisFieldUserPraiseTimes, v.PraiseTimes)
		pipe.HSet(key, redisFieldUserBelittleTimes, v.BelittleTimes)
		pipe.HSet(key, redisFieldUserLastEditTime, v.LastEditTime)
		//账户密码hash同步更新
		pipe.HSet(redisTopPreHashKeyUserAccount+v.Account, redisFieldUserAccountPwd, v.PassWord)
	}
	check1Err(pipe.Exec())
}

func (s *MySQLRedisIns) queryAllUsers() []*models.UserInDB {
	const sqlQueryUsers = "select * from tb_user"
	rows, err := s.DB.Query(sqlQueryUsers)
	checkErr(err)
	defer func() {
		checkErr(rows.Close())
	}()
	users := make([]*models.UserInDB, 0)
	for rows.Next() {
		user := s.ScanFromUserTbAllFields(rows)
		if user == nil {
			panic(errors.New("存在无效用户"))
		}
		users = append(users, user)
	}
	return users
}

const redisFieldPostID = "ID"
const redisFieldPostThemeID = "TID"
const redisFieldPostUserID = "UID"
const redisFieldPostTitle = "TIT"
const redisFieldPostState = "STA"
const redisFieldPostCreateTime = "CT"
const redisFieldPostCmtCount = "CC"
const redisFieldPostLastCmterID = "LCI"
const redisFieldPostLastCmtTime = "LCT"

//post section
func (s *MySQLRedisIns) copyPostDataFromMysqlToRedis() {
	pipe := s.rc.Pipeline()
	defer func() {
		checkErr(pipe.Close())
	}()
	postIDsInRedis := s.getPostIDsNowInRedis()
	posts := s.queryMostActivePosts(CacheInitPostCount)
	for _, v := range posts {
		if _, ok := postIDsInRedis[v.ID]; ok {
			continue
		}
		key := redisTopPreHashKeyPost + strconv.Itoa(v.ID)
		pipe.HSet(key, redisFieldPostID, v.ID)
		pipe.HSet(key, redisFieldPostThemeID, v.ThemeID)
		pipe.HSet(key, redisFieldPostUserID, v.UserID)
		pipe.HSet(key, redisFieldPostTitle, v.Title)
		pipe.HSet(key, redisFieldPostState, v.State)
		pipe.HSet(key, redisFieldPostCreateTime, v.CreatedTime)
		pipe.HSet(key, redisFieldPostCmtCount, v.CmtCount)
		pipe.HSet(key, redisFieldPostLastCmterID, v.LastCmterID)
		pipe.HSet(key, redisFieldPostLastCmtTime, v.LastCmtTime)
	}
	check1Err(pipe.Exec())
}

func (s *MySQLRedisIns) queryMostActivePosts(count int) []*models.PostInDB {
	const sqlQueryPosts = "select * from tb_post order by po_id desc limit ?"
	rows, err := s.DB.Query(sqlQueryPosts, count)
	checkErr(err)
	defer func() {
		checkErr(rows.Close())
	}()
	posts := make([]*models.PostInDB, 0)
	for rows.Next() {
		post := s.ScanFromPostTbAllFields(rows)
		if post == nil {
			panic(errors.New("存在无效的帖子"))
		}
		posts = append(posts, post)
	}
	return posts
}

const redisFieldCmtID = "ID"
const redisFieldCmtPostID = "PID"
const redisFieldCmtUserID = "UID"
const redisFieldCmtCTX = "CTX"
const redisFieldCmtState = "STA"
const redisFieldCmtCreateTime = "CT"
const redisFieldCmtLastEditTime = "LET"
const redisFieldCmtEditTimes = "ET"
const redisFieldCmtPriseTimes = "PT"
const redisFieldCmtBellowTimes = "BT"

//cmt section
func (s *MySQLRedisIns) copyCmtDataFromMysqlToRedis() {
	pipe := s.rc.Pipeline()
	defer func() {
		checkErr(pipe.Close())
	}()
	ids := s.getPostIDsNowInRedis()
	if len(ids) == 0 {
		return
	}
	cmtIDsInRedis := s.getCmtIDsNowInRedis()
	for id := range ids {
		cmts := s.MySQLIns.QueryComments(id)
		for _, cmt := range cmts {
			if _, ok := cmtIDsInRedis[cmt.ID]; ok {
				continue
			}
			key := redisTopPreHashKeyCmt + strconv.Itoa(cmt.ID)
			pipe.HSet(key, redisFieldCmtID, cmt.ID)
			pipe.HSet(key, redisFieldCmtPostID, cmt.PostID)
			pipe.HSet(key, redisFieldCmtUserID, cmt.UserID)
			pipe.HSet(key, redisFieldCmtCTX, cmt.Content)
			pipe.HSet(key, redisFieldCmtState, cmt.State)
			pipe.HSet(key, redisFieldCmtCreateTime, cmt.CreatedTime)
			pipe.HSet(key, redisFieldCmtLastEditTime, cmt.LastEditTime)
			pipe.HSet(key, redisFieldCmtEditTimes, cmt.EditTimes)
			pipe.HSet(key, redisFieldCmtPriseTimes, cmt.PraiseTimes)
			pipe.HSet(key, redisFieldCmtBellowTimes, cmt.BelittleTimes)
		}
		check1Err(pipe.Exec())
	}
}

const redisFieldPbID = "ID"
const redisFieldPbCmtID = "CID"
const redisFieldPbUserID = "UID"
const redisFieldPbPValue = "PV"
const redisFieldPbPTime = "PT"
const redisFieldPbPCTime = "PC"
const redisFieldPbBValue = "BV"
const redisFieldPbBTime = "BT"
const redisFieldPbBCTime = "BC"

func (s *MySQLRedisIns) copyPbDataFromMysqlToRedis() {
	pipe := s.rc.Pipeline()
	defer func() {
		checkErr(pipe.Close())
	}()
	postIDs := s.getPostIDsNowInRedis()
	if len(postIDs) == 0 {
		return
	}
	pbIDsInRedis := s.getPbIDsNowInRedis()
	for _, postID := range postIDs {
		for _, pbIns := range s.MySQLIns.QueryPbsOfPost(postID) {
			if _, ok := pbIDsInRedis[pbIns.ID]; ok {
				//不覆盖
				continue
			}
			key := redisTopPreHashKeyPb + strconv.Itoa(pbIns.ID)
			pipe.HSet(key, redisFieldPbID, pbIns.ID)
			pipe.HSet(key, redisFieldPbCmtID, pbIns.CmtID)
			pipe.HSet(key, redisFieldPbUserID, pbIns.UserID)
			pipe.HSet(key, redisFieldPbPValue, pbIns.PValue)
			pipe.HSet(key, redisFieldPbPTime, pbIns.PTime)
			pipe.HSet(key, redisFieldPbPCTime, pbIns.PCTime)
			pipe.HSet(key, redisFieldPbBValue, pbIns.BValue)
			pipe.HSet(key, redisFieldPbBTime, pbIns.BTime)
			pipe.HSet(key, redisFieldPbBCTime, pbIns.BCTime)
		}
		check1Err(pipe.Exec())
	}
}

func (s *MySQLRedisIns) getThemeIDsNowInRedis() map[int]int {
	return s.getItemIDsNowInRedis(redisTopPreHashKeyPost)
}

func (s *MySQLRedisIns) getUserIDsNowInRedis() map[int]int {
	return s.getItemIDsNowInRedis(redisTopPreHashKeyPost)
}

func (s *MySQLRedisIns) getPostIDsNowInRedis() map[int]int {
	return s.getItemIDsNowInRedis(redisTopPreHashKeyPost)
}

func (s *MySQLRedisIns) getCmtIDsNowInRedis() map[int]int {
	return s.getItemIDsNowInRedis(redisTopPreHashKeyPost)
}

func (s *MySQLRedisIns) getPbIDsNowInRedis() map[int]int {
	return s.getItemIDsNowInRedis(redisTopPreHashKeyPb)
}

func (s *MySQLRedisIns) getItemIDsNowInRedis(prefix string) map[int]int {
	ids := make(map[int]int, 64)
	preIndex := len(prefix)
	for _, str := range s.rc.Keys(prefix + "*").Val() {
		id, err := strconv.Atoi(str[preIndex:])
		checkErr(err)
		ids[id] = id
	}
	return ids
}

//以下为redis接口的实现

func (s *MySQLRedisIns) QueryTheme(themeID int) *models.ThemeInDB {
	data, err := s.rc.HGetAll(redisTopPreHashKeyTheme + strconv.Itoa(themeID)).Result()
	checkErr(err)
	return s.exchangeRedisStrMapToThemeIns(data)
}

func (s *MySQLRedisIns) QueryAllThemes() []*models.ThemeInDB {
	var err error
	keys, err := s.rc.Keys(redisTopPreHashKeyTheme + "*").Result()
	checkErr(err)
	pipe := s.rc.Pipeline()
	defer func() {
		checkErr(pipe.Close())
	}()
	for _, key := range keys {
		pipe.HGetAll(key)
	}
	results, err := pipe.Exec()
	checkErr(err)
	tms := make([]*models.ThemeInDB, 0, 16)
	for _, v := range results {
		strStrMap, _ := v.(*redis.StringStringMapCmd)
		dataMap, err := strStrMap.Result()
		checkErr(err)
		tms = append(tms, s.exchangeRedisStrMapToThemeIns(dataMap))
	}
	return tms
}

func (s *MySQLRedisIns) exchangeRedisStrMapToThemeIns(data map[string]string) *models.ThemeInDB {
	var err error
	tm := new(models.ThemeInDB)
	tm.ID, err = strconv.Atoi(data[redisFieldThemeID])
	checkErr(err)
	tm.Name = data[redisFieldThemeName]
	tm.PostCount, err = strconv.Atoi(data[redisFieldThemePostCount])
	checkErr(err)
	return tm
}

func (s *MySQLRedisIns) QueryPostCountOfTheme(themeID int) int {
	count, err := s.rc.HGet(redisTopPreHashKeyTheme+strconv.Itoa(themeID), redisFieldThemePostCount).Int()
	checkErr(err)
	return count
}

func (s *MySQLRedisIns) QueryPost(postID int) *models.PostInDB {
	dataMap, err := s.rc.HGetAll(redisTopPreHashKeyPost + strconv.Itoa(postID)).Result()
	checkErr(err)
	return s.exchangeRedisStrMapToPostIns(dataMap)
}

func (s *MySQLRedisIns) exchangeRedisStrMapToPostIns(data map[string]string) *models.PostInDB {
	var err error
	post := new(models.PostInDB)
	post.ID, err = strconv.Atoi(data[redisFieldPostID])
	checkErr(err)
	post.ThemeID, err = strconv.Atoi(data[redisFieldPostThemeID])
	checkErr(err)
	post.UserID, err = strconv.Atoi(data[redisFieldPostUserID])
	checkErr(err)
	post.Title = data[redisFieldPostTitle]
	post.State, err = strconv.Atoi(data[redisFieldPostState])
	checkErr(err)
	post.CreatedTime, err = strconv.ParseInt(data[redisFieldPostCreateTime], 10, 64)
	checkErr(err)
	post.CmtCount, err = strconv.Atoi(data[redisFieldPostCmtCount])
	checkErr(err)
	post.LastCmterID, err = strconv.Atoi(data[redisFieldPostLastCmterID])
	checkErr(err)
	post.LastCmtTime, err = strconv.ParseInt(data[redisFieldPostLastCmtTime], 10, 64)
	checkErr(err)
	return post
}

func (s *MySQLRedisIns) QueryPostTitle(postID int) string {
	title, err := s.rc.HGet(redisTopPreHashKeyPost+strconv.Itoa(postID), redisFieldPostTitle).Result()
	checkErr(err)
	return title
}

func (s *MySQLRedisIns) QueryPostsOfTheme(themeID int, count, offset, sortType int) []*models.PostOnThemePage {
	return nil
}

func (s *MySQLRedisIns) QueryPostsOfUser(userID int, count, offset int) []*models.PostOnThemePage {
	return nil
}

func (s *MySQLRedisIns) QueryPostOfPostPage(postID int) *models.PostOnPostPage {
	return nil
	//
	//ID        int
	//CreaterID int
	//Title     string
	//
	//CmtCount int
	//
	//ThemeID   int
	//ThemeName string
}
