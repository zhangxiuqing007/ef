package usecase

import (
	"ef/models"
	"errors"
	"time"
	"unicode/utf8"
)

//UserSignUpData 新用户注册传输用数据结构，由Controller创建。
type UserSignUpData struct {
	Account  string
	Password string
	Name     string
}

func (data UserSignUpData) buildUserIns() *models.UserInDB {
	user := new(models.UserInDB)
	user.ID = 0
	user.Account = data.Account
	user.PassWord = data.Password
	user.Name = data.Name
	user.Type = models.UserTypeNormalUser
	user.State = models.UserStateNormal
	user.SignUpTime = time.Now().UnixNano()
	user.PostCount = 0
	user.CommentCount = 0
	user.PraiseTimes = 0
	user.BelittleTimes = 0
	user.LastEditTime = 0
	return user
}

//AddUser signUp
func AddUser(data *UserSignUpData) error {
	//检查昵称合法性
	if utf8.RuneCountInString(data.Name) == 0 {
		return errors.New("昵称不合法（至少一个字）")
	}
	//检查账户合法性
	if utf8.RuneCountInString(data.Account) < 3 {
		return errors.New("账号不合法（至少三个字符）")
	}
	//检查密码合法性
	if utf8.RuneCountInString(data.Password) < 3 {
		return errors.New("密码不合法（至少三个字符）")
	}
	//检查昵称占用
	if db.IsUserNameExist(data.Name) {
		return errors.New("昵称被占用")
	}
	//检查账户占用
	if db.IsUserAccountExist(data.Account) {
		return errors.New("账号被占用")
	}
	//保存
	user := data.buildUserIns()
	return db.AddUser(user)
}

//QueryUserByAccountAndPwd 用户查询
func QueryUserByAccountAndPwd(account string, password string) (*models.UserInDB, error) {
	return db.QueryUserByAccountAndPwd(account, password)
}

//QueryUserByID 查询用户信息
func QueryUserByID(userID int) (*models.UserInDB, error) {
	return db.QueryUserByID(userID)
}

//QueryPostCountOfUser 查询用户的发帖量
func QueryPostCountOfUser(userID int) (int, error) {
	return db.QueryPostCountOfUser(userID)
}
