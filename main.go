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
	//读取数据库配置
	dba.MysqlUser = beego.AppConfig.String("mysqluser")
	dba.MysqlPwd = beego.AppConfig.String("mysqlpwd")
	dba.MysqlDb = beego.AppConfig.String("mysqldb")
	sqlIns := new(dba.MySQLIns)
	checkErr(sqlIns.Open(""))
	defer sqlIns.Close()
	usecase.SetDbInstance(sqlIns)
	//添加路由过滤器
	beego.InsertFilter("/*", beego.BeforeRouter, httpMethodRouterFilter)
	//设置静态文件路径
	beego.SetStaticPath("/static", "static")
	beego.Run()
}
