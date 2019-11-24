package controllers

import (
	"ef/tool"
	"ef/usecase"
	"io/ioutil"
	"path"
	"strconv"
)

type ImageController struct {
	baseController
}

//获取上传图片页
func (c *ImageController) Put() {
	c.TplName = "img_get.html"
	c.setLoginVmSelf()
}

type imageGetFormData struct {
	ImagePageIndex int
}

type ImageSelectImageVm struct {
	ImagePaths []string
}

//获取图片内容
func (c *ImageController) Get() {
	data := new(imageGetFormData)
	if err := c.ParseForm(data); err != nil {
		c.send400("请求数据错误")
		return
	}
	s := c.getSession()
	if s.UserID <= 0 {
		c.send401("请先登录")
		return
	}
	//查看图片总数量
	totalImgCount := usecase.GetUserImageCount(s.UserID)
	if totalImgCount == 0 {
		c.send404("查询不到已上传图片信息")
		return
	}
	oper := new(tool.PageNavigationOperator)
	pageIndex := oper.LimitPageIndex(data.ImagePageIndex, imgCountOnePage, totalImgCount)
	beginIndex, endIndex := oper.GetNavigationPageLimitIndex(pageIndex, imgCountOnePage, halfPageCountToNavigationOfImage, totalImgCount)
	pathBuilder := func(i int) string {
		return "/img?ImagePageIndex=" + strconv.Itoa(i)
	}
	c.setNavigationVm(oper.BuildPageNavigations(pathBuilder, beginIndex, pageIndex, endIndex))
	vm := new(ImageSelectImageVm)
	vm.ImagePaths = usecase.QueryImages(s.UserID, imgCountOnePage, imgCountOnePage*pageIndex)
	if len(vm.ImagePaths) == 0 {
		c.send404("无法找到图片")
		return
	}
	//发送compHTML出去
	c.TplName = "comp/comp_img_view_page.html"
	c.Data["vm"] = vm
}

//上传(新建)图片
func (c *ImageController) Post() {
	//检测权限
	s := c.getSession()
	if s.UserID <= 0 {
		c.send401("请先登录")
		return
	}
	//读取图片内容
	fileHeaders, err := c.GetFiles("images")
	if err != nil {
		c.send400("没有提交内容")
		return
	}
	/*验证图片合法性*/

	//开始解析图片
	imageUploads := make([]*usecase.ImageUpload, len(fileHeaders))
	for i, v := range fileHeaders {
		file, err := v.Open()
		if err != nil {
			c.send400("提交内容无法解读")
			return
		}
		buffer, err := ioutil.ReadAll(file)
		if err != nil {
			c.send400("提交内容无法解读")
			return
		}
		imageUploads[i] = new(usecase.ImageUpload).SetContent(buffer, s.UserID, path.Ext(v.Filename))
	}
	usecase.SaveImages(imageUploads)
	c.send200("上传成功")
}
