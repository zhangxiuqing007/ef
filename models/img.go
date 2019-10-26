package models

//ImageInDB 图片
type ImageInDB struct {
	ID         int
	UserID     int
	UploadTime int64
	FilePath   string
}
