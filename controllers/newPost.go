package controllers

import (
	"ef/usecase"
	"time"
)

type NewPostController struct {
	SessionBaseController
}

type newPostGetFormData struct {
	ThemeID int
	PostID  int
}

func (c *NewPostController) Get() {
	data := new(newPostGetFormData)
	if err := c.ParseForm(data); err != nil || data.ThemeID < 0 || data.PostID < 0 {
		c.send400()
		return
	}
	//检查登录情况
	if c.getSession().User == nil {
		c.send404()
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
		c.send404()
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
		c.send404()
		return
	}
	c.Data["vm"] = &postTitleEditVm{data.PostID, title}
	c.setLoginVmSelf()
	c.TplName = "newPost_get_title.html"
}

//新增帖子
func (c *NewPostController) Post() {
	data := new(usecase.PostAddData)
	if err := c.ParseForm(data); err != nil {
		c.send400()
		return
	}
	s := c.getSession()
	if s.User == nil {
		c.send400()
		return
	}
	data.UserID = s.User.ID
	/*检查用户权限*/
	if err := usecase.AddPost(data); err != nil {
		c.send404()
		return
	}
	//成功的话，直接发主题页，无论任何方式排序，都是在主题第一位
	c.sendThemePage(&themeFormData{
		ThemeID:   data.ThemeID,
		PageIndex: 0,
	})
}

type FixTitleFormData struct {
	PostID int
	Title  string
}

//修改帖子标题
func (c *NewPostController) Put() {
	s := c.getSession()
	if s.User == nil {
		c.send400()
		return
	}
	data := new(FixTitleFormData)
	if err := c.ParseForm(data); err != nil {
		c.send400()
		return
	}
	//查询旧的Post
	post, err := usecase.QueryPost(data.PostID)
	if err != nil {
		c.send404()
		return
	}
	//验证合法性
	if s.User.ID != post.UserID {
		c.send404()
		return
		//panic(errors.New("必须要发表者身份才能编辑标题"))
	}
	/*验证其他状态，包括帖子状态，标题是否可以修改，用户是否还有编辑权限*/
	//更新内容
	post.Title = data.Title
	post.LastCmtTime = time.Now().UnixNano()
	//保存到DB
	if err := usecase.UpdatePostTitle(post); err != nil {
		c.send404()
		return
	}
	//发送帖子页
	c.sendPostPage(&postFormData{
		PostID:    data.PostID,
		PageIndex: 0,
	})
}
