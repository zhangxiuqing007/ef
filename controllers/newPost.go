package controllers

import (
	"ef/usecase"
	"time"
	"unicode/utf8"
)

type NewPostController struct {
	baseController
}

type newPostGetFormData struct {
	ThemeID int
	PostID  int
}

func (c *NewPostController) Get() {
	//检查登录情况
	if c.getSession().UserID <= 0 {
		c.send401("请先登录")
		return
	}
	data := new(newPostGetFormData)
	if err := c.ParseForm(data); err != nil || data.ThemeID < 0 || data.PostID < 0 {
		c.send400("请求信息错误")
		return
	}
	if data.PostID == 0 {
		c.GetToBuildNewPost(data)
	} else {
		c.GetToEditPostTitle(data)
	}
}

type postInputVm struct {
	ThemeName string
	ThemeID   int
}

//请求新帖输入页
func (c *NewPostController) GetToBuildNewPost(data *newPostGetFormData) {
	tm, err := usecase.QueryTheme(data.ThemeID)
	if err != nil {
		c.send404("主题不存在")
		return
	}
	c.Data["vm"] = &postInputVm{tm.Name, tm.ID}
	c.setLoginVmSelf()
	c.TplName = "newPost_get_input.html"
}

type postTitleEditVm struct {
	PostID   int
	OriTitle string
}

//获取标题修改页
func (c *NewPostController) GetToEditPostTitle(data *newPostGetFormData) {
	title, err := usecase.QueryPostTitle(data.PostID)
	if err != nil {
		c.send404("帖子不存在")
		return
	}
	c.Data["vm"] = &postTitleEditVm{data.PostID, title}
	c.setLoginVmSelf()
	c.TplName = "newPost_get_title.html"
}

//新增帖子
func (c *NewPostController) Post() {
	s := c.getSession()
	if s.UserID <= 0 {
		c.send401("请先登录")
		return
	}
	data := new(usecase.PostAddData)
	if err := c.ParseForm(data); err != nil {
		c.send400("请求信息错误")
		return
	}
	//检查发布的内容
	if utf8.RuneCountInString(data.Title) < 2 || utf8.RuneCountInString(data.Content) < 2 {
		c.send403("标题或内容至少需要2个字")
		return
	}
	data.UserID = s.UserID
	/*检查用户权限*/
	if err := usecase.AddPost(data); err != nil {
		c.send406("操作失败：" + err.Error())
		return
	}
	//成功的话，直接发主题页，无论任何方式排序，都是在主题第一位
	c.sendThemePage(&themeFormData{
		ThemeID:   data.ThemeID,
		PageIndex: 0,
	})
}

type fixTitleFormData struct {
	PostID int
	Title  string
}

//修改帖子标题
func (c *NewPostController) Put() {
	s := c.getSession()
	if s.UserID <= 0 {
		c.send401("请先登录")
		return
	}
	data := new(fixTitleFormData)
	if err := c.ParseForm(data); err != nil {
		c.send400("请求信息错误")
		return
	}
	//查询旧的Post
	post, err := usecase.QueryPost(data.PostID)
	if err != nil {
		c.send404("帖子不存在")
		return
	}
	//验证合法性
	if s.UserID != post.UserID {
		c.send403("无修改权限")
		return
	}
	/*验证其他状态，包括帖子状态，标题是否可以修改，用户是否还有编辑权限*/

	//更新内容
	post.Title = data.Title
	post.LastCmtTime = time.Now().UnixNano()
	//保存到DB
	if err := usecase.UpdatePostTitle(post); err != nil {
		c.send406("操作失败：" + err.Error())
		return
	}
	//发送帖子页
	c.sendPostPage(&postFormData{
		PostID:    data.PostID,
		PageIndex: 0,
	})
}
