package controllers

import (
	"beelog/models"
    "log"

	"github.com/astaxie/beego"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	this.Data["IsHome"] = true
	this.TplName = "home.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	

	topics, err := models.GetTopics(this.Input().Get("cate"), this.Input().Get("label"),true)
	if err != nil {
		log.Fatal(err)
	}
	this.Data["Topics"] = topics
	
	categories, err := models.GetAllCategory()
	if err != nil{
		log.Fatal(err)
	}
	this.Data["Categories"] = categories
}