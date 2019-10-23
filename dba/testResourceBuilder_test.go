package dba

import (
	"ef/models"
	"ef/tool"
	"ef/usecase"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//测试资源创建者，放置包内有太多的全局函数
type testResourceBuilder struct {
	firstNameWords []string
	lastNameWords  []rune
}

//随机种子值
func (t *testResourceBuilder) initRandomSeed() {
	rand.Seed(time.Now().UnixNano())
}

//制造当前测试类型的sql对象
func (t *testResourceBuilder) buildCurrentTestSQLIns() usecase.IDataIO {
	db := new(MySQLIns)
	checkErr(db.Open(""))
	return db
	//db := SqliteIns{}
	//db.Open("../ef.db")
	//return &db
}

//生成随机主题
func (t *testResourceBuilder) buildRandomTheme(count int) []*models.ThemeInDB {
	tms := make([]*models.ThemeInDB, 0, count)
	for i := 0; i < count; i++ {
		newTheme := new(models.ThemeInDB)
		newTheme.Name = fmt.Sprintf("随机主题：%d", i+1)
		newTheme.PostCount = 0
		tms = append(tms, newTheme)
	}
	return tms
}

//生成随机用户
func (t *testResourceBuilder) buildRandomUsers(count int) []*models.UserInDB {
	users := make([]*models.UserInDB, 0, count)
	for i := 0; i < count; i++ {
		newUser := new(models.UserInDB)
		newUser.Account = tool.NewUUID()
		newUser.PassWord = tool.NewUUID()
		newUser.Name = "测试用户" + strconv.Itoa(i)
		if rand.Intn(2) == 1 {
			newUser.Type = models.UserTypeAdministrator
		} else {
			newUser.Type = models.UserTypeNormalUser
		}
		newUser.State = models.UserStateNormal
		newUser.SignUpTime = time.Now().UnixNano()
		newUser.PostCount = 0
		newUser.CommentCount = 0
		newUser.PraiseTimes = 0
		newUser.BelittleTimes = 0
		newUser.LastEditTime = 0 //最后一次编辑时间
		users = append(users, newUser)
	}
	return users
}

//生成随机帖子
func (t *testResourceBuilder) buildRandomPost(themeID, userID int) *models.PostInDB {
	post := new(models.PostInDB)
	post.ThemeID = themeID
	post.UserID = userID
	post.Title = t.buildRandomPostTitle()
	post.State = models.PostStateNormal
	post.CreatedTime = time.Now().UnixNano()
	post.CmtCount = 0
	post.LastCmterID = userID
	post.LastCmtTime = 0
	return post
}

//生成随机评论
func (t *testResourceBuilder) buildRandomCmt(postID, userID int) *models.CommentInDB {
	cmt := new(models.CommentInDB)
	cmt.PostID = postID
	cmt.UserID = userID
	cmt.Content = t.buildRandomPostContent()
	cmt.State = models.CmtStateNormal
	cmt.CreatedTime = time.Now().UnixNano()
	cmt.LastEditTime = cmt.CreatedTime
	cmt.EditTimes = 1
	cmt.PraiseTimes = 0
	cmt.BelittleTimes = 0
	return cmt
}

//生成随机标题
func (t *testResourceBuilder) buildRandomPostTitle() string {
	return t.combineUuids(rand.Int()%6 + 1)
}

//生成随机内容
func (t *testResourceBuilder) buildRandomPostContent() string {
	return t.combineUuids(rand.Int()%20 + 1)
}

//合并uuid
func (t *testResourceBuilder) combineUuids(count int) string {
	var uids = make([]string, 0, count)
	for i := 0; i < count; i++ {
		uids = append(uids, tool.NewUUID())
	}
	return strings.Join(uids, "#")
}

//判断两个帖子内容是否相同
func (t *testResourceBuilder) isTwoPostSame(post1, post2 *models.PostInDB) bool {
	return post1.ID == post2.ID &&
		post1.ThemeID == post2.ThemeID &&
		post1.UserID == post2.UserID &&
		post1.Title == post2.Title &&
		post1.State == post2.State &&
		post1.CreatedTime == post2.CreatedTime
	//post1.CmtCount == post2.CmtCount &&
	//post1.LastCmterID == post2.LastCmterID &&
	//post1.LastCmtTime == post2.LastCmtTime
}

//判断两个用户是否相同
func (t *testResourceBuilder) isTwoUserSame(user1, user2 *models.UserInDB) bool {
	return user1.ID == user2.ID &&
		user1.Account == user2.Account &&
		user1.PassWord != user2.PassWord &&
		user1.Name == user2.Name &&
		user1.Type == user2.Type &&
		user1.State == user2.State &&
		user1.SignUpTime == user2.SignUpTime
	//user1.PostCount == user2.PostCount &&
	//user1.CommentCount == user2.CommentCount &&
	//user1.PraiseTimes == user2.PraiseTimes &&
	//user1.BelittleTimes == user2.BelittleTimes &&
	//user1.LastEditTime == user2.LastEditTime
}

func (t *testResourceBuilder) isTwoPbSame(pb1, pb2 *models.PBInDB) bool {
	return pb1.ID == pb2.ID &&
		pb1.CmtID == pb2.CmtID &&
		pb1.UserID == pb2.UserID &&
		pb1.BValue == pb2.BValue &&
		pb1.PValue == pb2.PValue
}

//检查是否已经加载了名字资源
func (t *testResourceBuilder) checkInitNameResources() {
	if t.firstNameWords == nil {
		spe := []rune{' ', '\r', '\n'}
		t.firstNameWords = tool.SplitText(tool.MustStr(tool.ReadAllTextUtf8("../conf/中文姓氏.txt")), spe)
		t.lastNameWords = make([]rune, 0, 1800)
		for _, v := range tool.SplitText(tool.MustStr(tool.ReadAllTextUtf8("../conf/中文名字.txt")), spe) {
			rs := []rune(v)
			if len(rs) > 0 {
				t.lastNameWords = append(t.lastNameWords, rs[0])
			}
		}
	}
}

//随机生成中文名字
func (t *testResourceBuilder) buildRandomChineseName() string {
	t.checkInitNameResources()
	name := make([]rune, 0, 5)
	name = append(name, t.getRandomFirstNameWord()...)
	name = append(name, t.getRandomLastNameWord())
	//70%的人名字是两个字
	if rand.Intn(100) > 30 {
		name = append(name, t.getRandomLastNameWord())
	}
	return string(name)
}

func (t *testResourceBuilder) getRandomFirstNameWord() []rune {
	return []rune(t.firstNameWords[rand.Intn(len(t.firstNameWords))])
}

func (t *testResourceBuilder) getRandomLastNameWord() rune {
	return t.lastNameWords[rand.Intn(len(t.lastNameWords))]
}
