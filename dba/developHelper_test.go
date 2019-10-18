package dba

import (
	"ef/models"
	"math/rand"
	"testing"
	"time"
)

//目前mysql不能通过sql清空
// begin;
// SET foreign_key_checks = 0;
// truncate tb_cmt;
// truncate tb_post;
// truncate tb_user;
// truncate tb_theme;
// truncate tb_pb;
// SET foreign_key_checks = 1;
// commit;

//清空数据库	go test -v -run Test_ClearCurrentDb
func Test_ClearCurrentDb(t *testing.T) {
	iotool := new(testResourceBuilder).buildCurrentTestSQLIns()
	defer iotool.Close()
	err := iotool.Clear()
	if err != nil {
		t.Fatalf(err.Error())
	}
}

//增加标准主题		go test -v -run Test_HelpAddStandardThemes
func Test_HelpAddStandardThemes(t *testing.T) {
	iotool := new(testResourceBuilder).buildCurrentTestSQLIns()
	defer iotool.Close()
	iotool.AddTheme(&models.ThemeInDB{ID: 0, Name: "要闻", PostCount: 0})
	iotool.AddTheme(&models.ThemeInDB{ID: 0, Name: "国内", PostCount: 0})
	iotool.AddTheme(&models.ThemeInDB{ID: 0, Name: "国际", PostCount: 0})
	iotool.AddTheme(&models.ThemeInDB{ID: 0, Name: "社会", PostCount: 0})
	iotool.AddTheme(&models.ThemeInDB{ID: 0, Name: "军事", PostCount: 0})
	iotool.AddTheme(&models.ThemeInDB{ID: 0, Name: "娱乐", PostCount: 0})
	iotool.AddTheme(&models.ThemeInDB{ID: 0, Name: "体育", PostCount: 0})
	iotool.AddTheme(&models.ThemeInDB{ID: 0, Name: "汽车", PostCount: 0})
	iotool.AddTheme(&models.ThemeInDB{ID: 0, Name: "科技", PostCount: 0})
}

//增加一些用户，其中包括二把刀	go test -v -run Test_HelpAddSomeUsers
func Test_HelpAddSomeUsers(t *testing.T) {
	const addCount = 11
	rander := new(testResourceBuilder)
	rander.initRandomSeed()
	iotool := rander.buildCurrentTestSQLIns()
	defer iotool.Close()
	users := rander.buildRandomUsers(addCount)
	users[0].Name = "二把刀"
	users[0].Account = "erbadao"
	users[0].PassWord = "erbadao"
	if err := iotool.AddUser(users[0]); err != nil {
		t.Fatalf("x失败：添加测试用户：" + err.Error())
	} else {
		t.Log("成功：添加测试用户")
	}
	for i := 1; i < addCount; i++ {
		user := users[i]
		for {
			user.Name = rander.buildRandomChineseName()
			if iotool.IsUserNameExist(user.Name) {
				continue
			}
			break
		}
		if iotool.AddUser(user) != nil {
			t.Fatalf("x失败：添加随机用户")
		} else {
			t.Log("成功：添加测试用户")
		}
	}
}

//增加一些帖子和评论和赞踩	go test -v -run Test_HelpAddSomePostAndCmts
func Test_HelpAddSomePostAndCmts(t *testing.T) {
	const userCount = 11
	const themeCount = 9
	//帖子总数 1w
	const postMaxCount = 10000
	//评论总数 50w
	const cmtMaxCount = 500000
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
	cmts := make([]*models.CommentInDB, 0, cmtMaxCount)
	for i := 0; i < postMaxCount; i++ {
		userID := userIDs[rand.Intn(userCount)]
		posts = append(posts, rander.buildRandomPost(themeIDs[rand.Intn(themeCount)], userID))
		cmts = append(cmts, rander.buildRandomCmt(0, userID))
	}
	if err := iotool.AddPosts(posts, cmts); err != nil {
		t.Fatalf("x失败：插入巨量测试帖子，" + err.Error())
	} else {
		t.Log("成功：插入巨量测试帖子")
	}
	allCmts := make([]*models.CommentInDB, 0, cmtMaxCount+postMaxCount*2)
	cmts = cmts[0:0]
	//多轮评论
	for cmti := 0; cmti < cmtMaxCount; cmti++ {
		cmts = append(cmts, rander.buildRandomCmt(posts[rand.Intn(len(posts))].ID, userIDs[rand.Intn(userCount)]))
		if cmti == cmtMaxCount-1 || len(cmts) >= 50000 {
			allCmts = append(allCmts, cmts...)
			if iotool.AddComments(cmts) != nil {
				t.Fatalf("x失败：插入数万条随机评论")
			} else {
				t.Log("成功：插入数万条随机评论")
			}
			cmts = cmts[0:0]
		}
	}
	//随机赞踩
	for _, c := range allCmts {
		for _, u := range userIDs {
			//99%的概率不进行赞踩
			if rand.Intn(100) != 0 {
				continue
			}
			pb := new(models.PBInDB)
			pb.UserID = u
			pb.CmtID = c.ID
			if err := iotool.AddPbItem(pb); err != nil {
				t.Fatalf("x失败：新增赞踩")
			}
			//一半的概率赞，一半的概率踩
			if rand.Intn(10) >= 5 {
				pb.PTime = time.Now().UnixNano()
				if err := iotool.Praise(pb); err != nil {
					t.Fatal("x失败：赞 " + err.Error())
				}
			} else {
				pb.BTime = time.Now().UnixNano()
				if err := iotool.Belittle(pb); err != nil {
					t.Fatal("x失败：踩 " + err.Error())
				}
			}
		}
	}
}
