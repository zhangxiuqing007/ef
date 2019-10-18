package dba

import (
	"ef/models"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

//测试查询所有主题的速度
func BenchmarkReadAllThemes(b *testing.B) {
	rander := new(testResourceBuilder)
	sqlIns := rander.buildCurrentTestSQLIns()
	for i := 0; i < b.N; i++ {
		_, _ = sqlIns.QueryAllThemes()
	}
}

//测试逐个新增评论速度
func BenchmarkInsertComment(b *testing.B) {
	const userCount = 11
	const themeCount = 9
	const postMaxCount = 10

	rander := new(testResourceBuilder)
	rander.initRandomSeed()
	iotool := rander.buildCurrentTestSQLIns()
	defer iotool.Close()
	userIDs := [userCount]int{}
	for i := 0; i < userCount; i++ {
		userIDs[i] = 1 + i
	}
	themeIDs := [themeCount]int{}
	for i := 0; i < themeCount; i++ {
		themeIDs[i] = 1 + i
	}
	posts := make([]*models.PostInDB, 0, postMaxCount)
	cmts := make([]*models.CommentInDB, 0, 200)
	for i := 0; i < postMaxCount; i++ {
		userID := userIDs[rand.Intn(userCount)]
		posts = append(posts, rander.buildRandomPost(themeIDs[rand.Intn(themeCount)], userID))
		cmts = append(cmts, rander.buildRandomCmt(0, userID))
	}
	if err := iotool.AddPosts(posts, cmts); err != nil {
		b.Fatalf("x失败：插入一些测试帖子，" + err.Error())
	}
	//多轮评论
	for i := 0; i < b.N; i++ {
		cmtIns := rander.buildRandomCmt(posts[rand.Intn(len(posts))].ID, userIDs[rand.Intn(userCount)])
		if err := iotool.AddComment(cmtIns); err != nil {
			b.Fatal("x失败：插入评论 " + err.Error())
		}
	}
	b.Log("成功：插入大量评论：" + strconv.Itoa(b.N))
}

//测试批量新增评论速度
func BenchmarkBatchInsertComment(b *testing.B) {
	const userCount = 11
	const themeCount = 9
	const postMaxCount = 10

	rander := new(testResourceBuilder)
	rander.initRandomSeed()
	iotool := rander.buildCurrentTestSQLIns()
	defer iotool.Close()
	userIDs := [userCount]int{}
	for i := 0; i < userCount; i++ {
		userIDs[i] = 1 + i
	}
	themeIDs := [themeCount]int{}
	for i := 0; i < themeCount; i++ {
		themeIDs[i] = 1 + i
	}
	posts := make([]*models.PostInDB, 0, postMaxCount)
	cmts := make([]*models.CommentInDB, 0, 200)
	for i := 0; i < postMaxCount; i++ {
		userID := userIDs[rand.Intn(userCount)]
		posts = append(posts, rander.buildRandomPost(themeIDs[rand.Intn(themeCount)], userID))
		cmts = append(cmts, rander.buildRandomCmt(0, userID))
	}
	if err := iotool.AddPosts(posts, cmts); err != nil {
		b.Fatalf("x失败：插入一些测试帖子，" + err.Error())
	}
	//多轮评论
	for i := 0; i < b.N; i++ {
		cmts := make([]*models.CommentInDB, 0, 500)
		for i := 0; i < 500; i++ {
			cmts = append(cmts, rander.buildRandomCmt(posts[rand.Intn(len(posts))].ID, userIDs[rand.Intn(userCount)]))
		}
		if err := iotool.AddComments(cmts); err != nil {
			b.Fatal("x失败：批量插入评论 " + err.Error())
		}
	}
	b.Log("成功：插入大量评论：" + strconv.Itoa(b.N*500))
}

//测试查询主题页速度
func BenchmarkQueryThemePage(b *testing.B) {
	rander := new(testResourceBuilder)
	rander.initRandomSeed()
	iotool := rander.buildCurrentTestSQLIns()
	defer iotool.Close()
	tms, err := iotool.QueryAllThemes()
	if err != nil {
		b.Fatal("x错误：" + err.Error())
	}
	postCountOnce := 20
	//开始随机查询主题页帖子头列表
	for i := 0; i < b.N; i++ {
		//随机选取一个主题
		tm := tms[rand.Intn(len(tms))]
		//随意产生一个页Index
		pageIndex := tm.PostCount / postCountOnce
		if pageIndex > 0 {
			pageIndex = rand.Intn(pageIndex + 1)
		}
		//定义一个排序类型
		sortType := 1
		if _, err := iotool.QueryPostsOfTheme(tm.ID, postCountOnce, postCountOnce*pageIndex, sortType); err != nil {
			b.Fatal("x错误：" + err.Error())
		}
	}
}

//测试查询帖子页速度
func BenchmarkQueryPostPage(b *testing.B) {
	rander := new(testResourceBuilder)
	rander.initRandomSeed()
	iotool := rander.buildCurrentTestSQLIns()
	defer iotool.Close()
	tms, err := iotool.QueryAllThemes()
	if err != nil {
		b.Fatal("x错误：" + err.Error())
	}
	maxPostCount := 1000
	tmPostMap := make(map[int][]*models.PostOnThemePage, 16)
	for _, v := range tms {
		tmPostMap[v.ID] = make([]*models.PostOnThemePage, 0, maxPostCount)
		posts, err := iotool.QueryPostsOfTheme(v.ID, maxPostCount, 0, 0)
		if err != nil {
			b.Fatal("x错误：" + err.Error())
		}
		for _, p := range posts {
			tmPostMap[v.ID] = append(tmPostMap[v.ID], p)
		}
	}
	cmtCountOnce := 20
	//开始随机查询主题页帖子头列表
	for i := 0; i < b.N; i++ {
		//随机选取一个主题
		rnadTm := tms[rand.Intn(len(tms))]
		postHeaders := tmPostMap[rnadTm.ID]
		//随机选取一个帖子的ID
		post := postHeaders[rand.Intn(len(postHeaders))]
		//随意产生一个页Index
		pageIndex := post.CmtCount / cmtCountOnce
		if pageIndex > 0 {
			pageIndex = rand.Intn(pageIndex + 1)
		}
		//随机产生一个用户ID
		userID := rand.Intn(11) + 1
		if _, err := iotool.QueryCommentsOfPostPage(post.ID, cmtCountOnce, cmtCountOnce*pageIndex, userID); err != nil {
			b.Fatal("x错误：" + err.Error())
		}
	}
}

//测试赞踩评论速度
func BenchmarkPB(b *testing.B) {
	sqlIns := new(testResourceBuilder).buildCurrentTestSQLIns()
	rand.Seed(time.Now().UnixNano())
	const userCount = 11
	const cmtCount = 100000
	for i := 0; i < b.N; i++ {
		randUserID := rand.Intn(userCount) + 1
		randCmtID := rand.Intn(cmtCount) + 1
		isP := rand.Intn(2) == 1
		isD := rand.Intn(2) == 1
		if err := sqlIns.SetPB(randCmtID, randUserID, isP, isD); err != nil {
			b.Fatal("x失败：赞踩请求失败 " + err.Error())
		}
	}
}
