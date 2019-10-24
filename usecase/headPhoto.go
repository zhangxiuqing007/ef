package usecase

import (
	"ef/tool"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"os"
)

var HeadPhotoMinWidth int
var HeadPhotoMaxWidth int
var HeadPhotoMinHeight int
var HeadPhotoMaxHeight int

//修改头像，这里兼容gif文件
func ChangeHeadPhotoGeneral(userID int, file []byte, img image.Image, ext string) error {
	if len(file) == 0 {
		return errors.New("错误的文件内容")
	}
	if img == nil {
		return errors.New("必须是图片才能作为头像")
	}
	size := img.Bounds().Size()
	if size.X < HeadPhotoMinWidth || size.X > HeadPhotoMaxWidth || size.Y < HeadPhotoMinHeight || size.Y > HeadPhotoMaxHeight {
		return errors.New("头像图片尺寸不符合要求")
	}
	path := fmt.Sprintf("%s/%d-%s%s", GetCurrentHeadPhotoFolder(), userID, tool.NewUUID(), ext)
	//保存到数据库中
	err := db.UpdateUserHeadPhoto(userID, path)
	if err != nil {
		return err
	}
	//保存到本地硬盘中
	return ioutil.WriteFile(path, file, os.ModeType)
}
