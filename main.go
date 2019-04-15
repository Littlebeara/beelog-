package main

import (
	"os"
	"beelog/controllers"
	"beelog/models"
	
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {
	// 注册数据库
	models.RegisterDB()
}

func main() {
	// 开启 ORM 调试模式
	orm.Debug = true
	// 自动建表
	orm.RunSyncdb("default", false, true)

	// 注册 beego 路由
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/category", &controllers.CategoryController{})
	beego.Router("/topic", &controllers.TopicController{})
	beego.Router("/reply", &controllers.ReplyController{})
	beego.Router("/reply/add", &controllers.ReplyController{}, "post:Add")
	beego.Router("/reply/delete", &controllers.ReplyController{}, "Get:Delete")
	beego.AutoRouter(&controllers.TopicController{})
	beego.Router("/login", &controllers.LoginController{})
	
	//创建附件目录
	os.Mkdir("attachment", os.ModePerm)
	
	//作为静态文件
	//beego.SetStaticPath("/attachment","attachment")
	
	//作为一个单独的控制器来处理
	beego.Router("/attachment/:all", &controllers.AttachController{})
	
	// 启动 beego
	beego.Run()
}