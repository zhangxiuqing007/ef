package usecase

import (
	"ef/models"
	"errors"
	"time"
	"unicode/utf8"
)

var defaultHeadPhotoPath string

func InitDefaultHeadPhotoPath(path string) {
	defaultHeadPhotoPath = path
}

//UserSignUpData 新用户注册传输用数据结构，由Controller创建。
type UserSignUpData struct {
	Account  string
	Password string
	Name     string
}

func (data UserSignUpData) buildUserIns() *models.UserInDB {
	return &models.UserInDB{
		ID:            0,
		Account:       data.Account,
		PassWord:      data.Password,
		Name:          data.Name,
		HeadPhotoPath: defaultHeadPhotoPath,
		Type:          models.UserTypeNormalUser,
		State:         models.UserStateNormal,
		SignUpTime:    time.Now().UnixNano(),
		PostCount:     0,
		CommentCount:  0,
		ImageCount:    0,
		PraiseTimes:   0,
		BelittleTimes: 0,
		LastEditTime:  0,
	}
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
