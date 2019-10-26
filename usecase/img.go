package usecase

import (
	"ef/models"
	"ef/tool"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type ImageUpload struct {
	Buffer   []byte
	UserID   int
	FilePath string
}

func (i *ImageUpload) SetContent(buffer []byte, userID int, fileExt string) *ImageUpload {
	i.Buffer = buffer
	i.UserID = userID
	i.FilePath = fmt.Sprintf("%s/%d-%s%s", GetCurrentGeneralImageFolder(), userID, tool.NewUUID(), fileExt)
	return i
}

func (i *ImageUpload) buildDBIns() *models.ImageInDB {
	return &models.ImageInDB{
		ID:         0,
		UserID:     i.UserID,
		UploadTime: time.Now().UnixNano(),
		FilePath:   i.FilePath,
	}
}

//保存图片，到本地硬盘和数据库
func SaveImages(images []*ImageUpload) error {
	//逐张图片保存到本地硬盘中，并转换成数据库结构
	dbImages := make([]*models.ImageInDB, len(images))
	for i, v := range images {
		err := ioutil.WriteFile(v.FilePath, v.Buffer, os.ModeType)
		if err != nil {
			return err
		}
		dbImages[i] = v.buildDBIns()
	}
	//全部保存到数据库中
	return db.AddImages(dbImages)
}

func GetUserImageCount(userID int) (int, error) {
	return db.QueryImageCountOfUser(userID)
}

func QueryImages(userID int, count int, offset int) ([]string, error) {
	return db.QueryImages(userID, count, offset)
}
