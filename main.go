package main

import (
	"ef/dba"
	_ "ef/routers"
	"ef/usecase"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//过滤器，保证使用除GET和POST外的方法
func httpMethodRouterFilter(ctx *context.Context) {
	if !ctx.Input.IsPost() {
		return
	}
	if method := ctx.Input.Query("_method"); method != "" {
		ctx.Request.Method = method
	}
}

func main() {
	sqlIns := new(dba.MySQLIns)
	checkErr(sqlIns.Open(""))
	defer sqlIns.Close()
	usecase.SetDbInstance(sqlIns)
	//读取数据库配置
	dba.MysqlUser = beego.AppConfig.String("mysqluser")
	dba.MysqlPwd = beego.AppConfig.String("mysqlpwd")
	dba.MysqlDb = beego.AppConfig.String("mysqldb")
	//添加路由过滤器
	beego.InsertFilter("/*", beego.BeforeRouter, httpMethodRouterFilter)
	beego.Run()

	////URL路由
	//router := httprouter.New()
	//router.GET("/", controllers.Index)
	//
	//router.GET("/UserRegist", controllers.UserRegist)
	//router.POST("/UserRegistCommit", controllers.UserRegistCommit)
	//
	//router.GET("/Login", controllers.Login)
	//router.POST("/LoginCommit", controllers.LoginCommit)
	//
	//router.GET("/Exit", controllers.Exit)
	//
	//router.GET("/Theme/:themeID/:pageIndex", controllers.Theme)
	//
	//router.GET("/User/:userID", controllers.UserInfo)
	//router.GET("/User/:userID/:pageIndex", controllers.UserPosts)
	//
	//router.GET("/Post/Content/:postID/:pageIndex", controllers.PostInfo)
	//router.GET("/Post/TitleEdit/:postID", controllers.PostTitleEdit)
	//router.POST("/Post/TitleEditSubmit", controllers.PostTitleEditCommit)
	//
	//router.GET("/NewPostInput/:themeID", controllers.NewPostInput)
	//router.POST("/NewPostCommit", controllers.NewPostCommit)
	//
	//router.POST("/Cmt", controllers.Cmt)
	//router.GET("/Cmt/Edit/:cmtID/:cmtPageIndex", controllers.CmtEdit)
	//router.POST("/Cmt/EditSubmit", controllers.CmtEditCommit)
	//router.POST("/Cmt/PG", controllers.CmtPb)
	//
	//fmt.Println("开始监听HTTP请求...")
	//if err := http.ListenAndServe("localhost:8080", router); err != nil {
	//	fmt.Print("程序启动失败：" + err.Error())
	//}
}
