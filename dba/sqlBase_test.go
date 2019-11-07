package dba

import (
	"ef/models"
	"ef/tool"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

//测试主题表相关操作	go test -v -run TestThemeTableOperations$
func TestThemeTableOperations(t *testing.T) {
	randTool := new(testResourceBuilder)
	randTool.initRandomSeed()
	sqlIns := randTool.buildCurrentTestSQLIns()
	defer sqlIns.Close()
	const testCount = 5
	//逐个新增主题
	tms := randTool.buildRandomTheme(testCount)
	for i := 0; i < testCount; i++ {
		if err := sqlIns.AddTheme(tms[i]); err != nil {
			t.Fatalf("x失败：新增主题 " + err.Error())
		}
	}
	t.Logf("成功：新增主题")
	//逐个更新主题名称
	for i, v := range tms {
		v.Name = fmt.Sprintf("主题改名%d", i)
		if sqlIns.UpdateTheme(v) != nil {
			t.Fatalf("x失败：修改主题名称")
		}
	}
	t.Log("成功：修改主题名")
	//逐个查询主题并对比信息
	for _, v := range tms {
		qtm, err := sqlIns.QueryTheme(v.ID)
		if err != nil || qtm.ID != v.ID || qtm.Name != v.Name {
			t.Fatalf("x失败：查询主题")
		}
	}
	t.Log("成功：查询主题，一致")
	//逐个查询主题的帖子数量
	for _, v := range tms {
		if _, err := sqlIns.QueryPostCountOfTheme(v.ID); err != nil {
			t.Fatal("x失败：查询主题的帖子数量")
		}
	}
	t.Log("成功：查询主题帖子总量")
	//查询所有主题
	if _, err := sqlIns.QueryAllThemes(); err != nil {
		t.Fatalf("x失败：查询所有主题失败")
	}
	t.Log("成功：查询所有主题")
	//删除刚才新增的主题
	for _, v := range tms {
		if err := sqlIns.DeleteTheme(v.ID); err != nil {
			t.Fatalf("x失败：删除主题")
		}
	}
	t.Logf("成功：删除主题")
}

//测试用户表相关操作	go test -v -run TestUserTableOperations
func TestUserTableOperations(t *testing.T) {
	randTool := new(testResourceBuilder)
	randTool.initRandomSeed()
	sqlIns := randTool.buildCurrentTestSQLIns()
	defer sqlIns.Close()
	//创建随机个用户
	const testCount = 5
	users := randTool.buildRandomUsers(testCount)
	//新增用户
	for _, v := range users {
		if sqlIns.AddUser(v) != nil {
			t.Fatalf("x失败：新增用户" + v.Name)
		}
	}
	t.Log("成功：新增用户")
	//通过id查询用户
	for _, v := range users {
		quser, err := sqlIns.QueryUserByID(v.ID)
		if err != nil || !randTool.isTwoUserSame(v, quser) {
			t.Fatalf("x失败：通过id查询用户")
		}
	}
	t.Log("成功：通过id查询用户")
	//通过账户密码查询用户
	for _, v := range users {
		quser, err := sqlIns.QueryUserByAccountAndPwd(v.Account, v.PassWord)
		if err != nil || !randTool.isTwoUserSame(v, quser) {
			t.Fatalf("x失败：通过账户密码查询用户")
		}
	}
	t.Log("成功：通过账户密码查询用户")
	//查询用户名是否存在
	//查询用户账号是否存在
	for _, v := range users {
		if !sqlIns.IsUserNameExist(v.Name) || !sqlIns.IsUserAccountExist(v.Account) {
			t.Fatalf("x失败：查询是否用户名或账号已存在")
		}
	}
	t.Log("成功：查询是否用户名或账号已存在")
	//查询用户发帖量
	for _, v := range users {
		if _, err := sqlIns.QueryPostCountOfUser(v.ID); err != nil {
			t.Fatal("x失败：查询用户发帖量")
		}
		t.Log("成功：查询用户发帖量")
	}
	//删除用户
	for _, v := range users {
		if sqlIns.DeleteUser(v.ID) != nil {
			t.Error("x失败：删除用户" + v.Name)
		}
	}
	t.Log("成功：删除用户")
}

//帖子和评论增删改查，赞踩操作	go test -v -run Test_PostAndCmt
func Test_PostAndCmt(t *testing.T) {
	randTool := new(testResourceBuilder)
	randTool.initRandomSeed()
	sqlIns := randTool.buildCurrentTestSQLIns()
	defer sqlIns.Close()
	const testUserCount = 5
	const cmtCount = 50
	//创建临时主题
	tmIns := randTool.buildRandomTheme(1)[0]
	if err := sqlIns.AddTheme(tmIns); err != nil {
		t.Fatal("x失败：新增测试主题 " + err.Error())
	}
	//创建临时用户
	users := randTool.buildRandomUsers(testUserCount)
	for _, v := range users {
		if err := sqlIns.AddUser(v); err != nil {
			t.Fatal("x失败，新增测试用户：" + err.Error())
		}
	}
	//每个用户创建2个帖子以及主内容
	posts := make([]*models.PostInDB, 0, testUserCount*2)
	cmts := make([]*models.CommentInDB, 0, testUserCount*2)
	for i := 0; i < testUserCount; i++ {
		userID := users[i].ID
		posts = append(posts, randTool.buildRandomPost(tmIns.ID, userID))
		posts = append(posts, randTool.buildRandomPost(tmIns.ID, userID))
		cmts = append(cmts, randTool.buildRandomCmt(0, userID))
		cmts = append(cmts, randTool.buildRandomCmt(0, userID))
	}
	//逐个新增帖子，前半组
	for i := 0; i < testUserCount; i++ {
		if err := sqlIns.AddPost(posts[i], cmts[i]); err != nil {
			t.Fatalf("x失败：新增帖子：" + err.Error())
		}
	}
	t.Log("成功：新增帖子")
	//批量新增帖子，后半组
	if err := sqlIns.AddPosts(posts[testUserCount:], cmts[testUserCount:]); err != nil {
		t.Fatalf("x失败：批量新增帖子 " + err.Error())
	}
	t.Log("成功：批量新增帖子")

	//逐个新增评论（其实是帖子的主内容）
	for _, v := range cmts {
		if err := sqlIns.AddComment(v); err != nil {
			t.Fatalf("x失败：新增帖子主内容 " + err.Error())
		}
	}
	t.Log("成功：逐个新增评论")
	//追加一定数量的评论
	cmts = make([]*models.CommentInDB, 0, cmtCount)
	for i := 0; i < cmtCount; i++ {
		postID := posts[rand.Intn(len(posts))].ID
		userID := users[rand.Intn(len(users))].ID
		cmts = append(cmts, randTool.buildRandomCmt(postID, userID))
	}
	//批量增加评论
	if err := sqlIns.AddComments(cmts); err != nil {
		t.Fatalf("x失败：批量增加评论 " + err.Error())
	}
	t.Log("成功：批量新增评论")
	//查询帖子
	for _, v := range posts {
		p, err := sqlIns.QueryPost(v.ID)
		if err != nil || !randTool.isTwoPostSame(p, v) {
			t.Fatalf("x失败：查询帖子失败或内容不一致 ")
		}
		title, err := sqlIns.QueryPostTitle(v.ID)
		if err != nil || title != v.Title {
			t.Fatalf("x失败：标题不一致 ")
		}
		v.Title = title
		v.LastCmtTime = time.Now().UnixNano()
		err = sqlIns.UpdatePostTitle(v)
		if err != nil {
			t.Fatal("x失败：无法更新帖子标题 " + err.Error())
		}
	}
	t.Log("成功：查询单个帖子，查询帖子标题，更新帖子标题")
	//查询主题帖子总数量
	if count, err := sqlIns.QueryPostCountOfTheme(tmIns.ID); err != nil || count != testUserCount*2 {
		t.Fatalf("x失败：查询主题帖子总量")
	}
	t.Log("成功：查询主题帖子总量")
	//查询用户发帖总数量
	for _, v := range users {
		if count, err := sqlIns.QueryPostCountOfUser(v.ID); err != nil || count != 2 {
			t.Fatalf("x失败：查询用户发帖总量")
		}
	}
	t.Log("成功：查询用户发帖总量")
	//查询主题下的帖子列表
	if ps, err := sqlIns.QueryPostsOfTheme(tmIns.ID, testUserCount, testUserCount, 0); err != nil || len(ps) != testUserCount {
		t.Fatalf("x失败：查询主题下的帖子列表，按发帖顺序排序")
	}
	if ps, err := sqlIns.QueryPostsOfTheme(tmIns.ID, testUserCount, testUserCount, 1); err != nil || len(ps) != testUserCount {
		t.Fatalf("x失败：查询主题下的帖子列表，按最后评论顺序排序")
	}
	t.Log("成功：查询主题下的帖子列表，按照两种排序")
	//查询用户发的帖子列表
	for _, v := range users {
		if ps, err := sqlIns.QueryPostsOfUser(v.ID, 1, 1); err != nil || len(ps) != 1 {
			t.Fatalf("x失败：查询用户发的帖子列表")
		}
	}
	t.Log("成功：查询用户发的的帖子列表")
	for _, v := range posts {
		//查询帖子页内，帖子的展示内容
		if p, err := sqlIns.QueryPostOfPostPage(v.ID); err != nil || p.Title != v.Title || p.ThemeName != tmIns.Name {
			t.Fatalf("x失败：查询帖子页内，帖子的展示内容")
		}
		//查询DB评论
		cs, err := sqlIns.QueryComments(v.ID)
		if err != nil || cs[0].PostID != v.ID || cs[0].UserID != v.UserID {
			t.Fatalf("x失败：查询DB评论")
		}
		//查询帖子页内，评论的展示内容
		scs, err := sqlIns.QueryCommentsOfPostPage(v.ID, 20, 0, v.UserID)
		if err != nil || scs[0].CmterID != v.UserID {
			t.Fatalf("x失败：查询帖子页内，评论的展示内容 ")
		}
	}
	//查询单个评论Item
	for _, v := range cmts {
		cmt, err := sqlIns.QueryComment(v.ID)
		if err != nil {
			t.Fatal("x失败：查询单个评论失败")
		}
		cmt.Content = randTool.buildRandomPostContent()
		cmt.LastEditTime = time.Now().UnixNano()
		cmt.EditTimes++
		if err := sqlIns.UpdateComment(cmt); err != nil {
			t.Fatal("x失败：修改单个评论 " + err.Error())
		}
	}
	t.Log("成功：查询帖子内容、查询帖子所有评论，帖子页内评论内容数组，查询单个评论，修改单个评论")

	//赞踩测试
	for _, u := range users {
		for _, c := range cmts {
			pb := new(models.PBInDB)
			pb.CmtID = c.ID
			pb.UserID = u.ID
			if err := sqlIns.AddPbItem(pb); err != nil {
				t.Fatal("x失败：" + err.Error())
			}
			if err := sqlIns.Praise(pb); err != nil {
				t.Fatal("x失败：" + err.Error())
			}
			pb.PValue = 1
			if err := sqlIns.PraiseCancel(pb); err != nil {
				t.Fatal("x失败：" + err.Error())
			}
			pb.PValue = 0
			if err := sqlIns.Belittle(pb); err != nil {
				t.Fatal("x失败：" + err.Error())
			}
			pb.BValue = 1
			if err := sqlIns.BelittleCancel(pb); err != nil {
				t.Fatal("x失败：" + err.Error())
			}
			pb.BValue = 0
			pbq, err := sqlIns.QueryPbItem(pb.CmtID, pb.UserID)
			if err != nil {
				t.Fatal("x失败：查询pbItem：" + err.Error())
			}
			if !randTool.isTwoPbSame(pbq, pb) {
				t.Fatal("x失败：两PB不一致")
			}
			checkErr(sqlIns.SetPB(c.ID, u.ID, rand.Float64() > 0.5, rand.Float64() > 0.5))
		}
	}
	//赞踩验证
	for _, v := range users {
		if ur, err := sqlIns.QueryUserByID(v.ID); err != nil || ur.PraiseTimes != 0 || ur.BelittleTimes != 0 {
			t.Fatal("x失败：赞踩数量不对")
		}
	}
	//赞踩验证
	for _, v := range posts {
		tempCmts, err := sqlIns.QueryComments(v.ID)
		if err != nil {
			t.Fatal("x失败：查询帖子的所有评论 " + err.Error())
		}
		for _, tc := range tempCmts {
			if tc.PraiseTimes != 0 || tc.BelittleTimes != 0 {
				t.Fatal("x失败：评论赞踩数量不对")
			}
		}
	}
	t.Log("成功：所有赞踩执行成功")

	//逐个删除评论
	for i := 0; i < len(cmts)/2; i++ {
		if sqlIns.DeleteComment(cmts[i].ID) != nil {
			t.Fatalf("x失败：删除评论")
		}
	}
	t.Log("成功：删除评论")
	//删除帖子（连同其剩余的评论）
	for _, v := range posts {
		if err := sqlIns.DeletePost(v.ID); err != nil {
			t.Fatalf("x失败：删除帖子（连同其剩余的评论）：" + err.Error())
		}
	}
	t.Log("成功：删除帖子")
	//批量新增图片
	imgs := make([]*models.ImageInDB, 0, 100)
	for i := 0; i < 100; i++ {
		imgs = append(imgs, &models.ImageInDB{
			ID:         0,
			UserID:     users[rand.Intn(len(users))].ID,
			UploadTime: time.Now().UnixNano(),
			FilePath:   tool.NewUUID() + ".png",
		})
	}
	checkErr(sqlIns.AddImages(imgs))
	//图片查询
	for _, v := range users {
		_, err := sqlIns.QueryImages(v.ID, rand.Intn(10), rand.Intn(100))
		checkErr(err)
	}
	//查询用户上传图片数量
	for _, v := range users {
		_, err := sqlIns.QueryImageCountOfUser(v.ID)
		checkErr(err)
	}
	//更新用户头像
	for _, v := range users {
		checkErr(sqlIns.UpdateUserHeadPhoto(v.ID, tool.NewUUID()+".png"))
	}

	//删除用户
	for _, v := range users {
		checkErr(sqlIns.DeleteUser(v.ID))
	}
	//删除主题
	checkErr(sqlIns.DeleteTheme(tmIns.ID))
	t.Log("成功：帖子和评论增删改查操作")
}
