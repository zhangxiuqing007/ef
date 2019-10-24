package usecase

import (
	"fmt"
	"os"
	"time"
)

var cacheCurrentHeadPhotoFolderName string = ""
var cacheCurrentGeneralImageFolderName string = ""

func GetCurrentHeadPhotoFolder() string {
	t := time.Now()
	folderName := fmt.Sprintf("static/img/headPhoto/%d-%d", t.Year(), t.Month())
	if cacheCurrentHeadPhotoFolderName != folderName {
		checkErr(os.MkdirAll(folderName, os.ModeType))
		cacheCurrentHeadPhotoFolderName = folderName
	}
	return folderName
}

//获取当前普通图片文件夹
func GetCurrentGeneralImageFolder() string {
	t := time.Now()
	folderName := fmt.Sprintf("static/img/personal/%d-%d-%d", t.Year(), t.Month(), t.Day())
	if cacheCurrentGeneralImageFolderName != folderName {
		checkErr(os.MkdirAll(folderName, os.ModeType))
		cacheCurrentGeneralImageFolderName = folderName
	}
	return folderName
}
