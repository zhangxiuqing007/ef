package main

import (
	"ef/controllers"
	"ef/dba"
	_ "ef/routers"
	"ef/usecase"
	"encoding/gob"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	_ "github.com/astaxie/beego/session/redis"
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
	var err error
	//允许redis的session，注册结构体
	gob.Register(&controllers.Session{})
	//初始化数据库配置
	dba.MysqlUser = beego.AppConfig.String("mysqluser")
	dba.MysqlPwd = beego.AppConfig.String("mysqlpwd")
	dba.MysqlDb = beego.AppConfig.String("mysqldb")
	//初始化session配置：名称
	controllers.SessionCookieKey = beego.BConfig.WebConfig.Session.SessionName
	//初始化头像初始路径
	usecase.DefaultHeadPhotoPath = beego.AppConfig.String("defaultHeadPhotoPath")
	usecase.HeadPhotoMinWidth, err = beego.AppConfig.Int("headPhotoMinWidth")
	checkErr(err)
	usecase.HeadPhotoMaxWidth, err = beego.AppConfig.Int("headPhotoMaxWidth")
	checkErr(err)
	usecase.HeadPhotoMinHeight, err = beego.AppConfig.Int("headPhotoMinHeight")
	checkErr(err)
	usecase.HeadPhotoMaxHeight, err = beego.AppConfig.Int("headPhotoMaxHeight")
	checkErr(err)
	//初始化数据库对象
	sqlIns := new(dba.MySQLIns)
	sqlIns.Open("")
	beego.BConfig.RecoverPanic = true
	defer sqlIns.Close()
	usecase.SetDbInstance(sqlIns)
	//添加路由过滤器
	beego.InsertFilter("/*", beego.BeforeRouter, httpMethodRouterFilter)
	//设置静态文件路径
	beego.SetStaticPath("/static", "static")
	beego.Run()
}
