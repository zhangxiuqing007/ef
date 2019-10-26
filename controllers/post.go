package controllers

type postFormData struct {
	PostID    int
	PageIndex int
}

type PostController struct {
	baseController
}

func (c *PostController) Get() {
	data := new(postFormData)
	if err := c.ParseForm(data); err != nil || data.PostID <= 0 || data.PageIndex < 0 {
		c.send400("请求信息错误")
		return
	}
	c.sendPostPage(data)
}
