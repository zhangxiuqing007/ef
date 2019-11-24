package controllers

import (
	"bytes"
	"ef/usecase"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"path"
)

type HeadPhotoController struct {
	baseController
}

//获取头像上传页
func (c *HeadPhotoController) Get() {
	if c.getSession().UserID <= 0 {
		c.send401("请先登录")
		return
	}
	c.TplName = "headPhoto_get.html"
}

//修改头像图片
func (c *HeadPhotoController) Put() {
	s := c.getSession()
	if s.UserID <= 0 {
		c.send401("请先登录")
		return
	}
	file, header, err := c.GetFile("headPhotoFile")
	if err != nil {
		c.send400("无法读取文件")
		return
	}
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		c.send400("无法读取文件")
		return
	}
	img, _, err := image.Decode(bytes.NewBuffer(buf))
	if err != nil {
		c.send400("不是有效的图片文件")
		return
	}
	usecase.ChangeHeadPhotoGeneral(s.UserID, buf, img, path.Ext(header.Filename))
	//重新发送用户页
	c.sendUserPage(&userFromData{UserID: s.UserID})
}
