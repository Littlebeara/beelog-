package controllers

import (
	"log"
	"beelog/models"
	"github.com/astaxie/beego"
)

type CategoryController struct {
	beego.Controller
}

func (this *CategoryController) Get() {
	// 检查是否有操作
	op := this.Input().Get("op")
	switch op{
		case "add":
			name := this.Input().Get("name")
			if len(name) == 0{
				break
			}
			err := models.AddCategory(name)
			if err != nil{
				log.Fatal(err)
			}
			this.Redirect("/category",301)
			return
		case "del":
			id := this.Input().Get("id")
			if len(id) == 0{
				break
			}
				
			err := models.DelCategory(id)
			if err != nil{
				log.Fatal(err)
			}
				
			this.Redirect("/category", 301)
			return	
			
	}

	this.Data["IsCategory"] = true
	this.TplName = "category.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	
	var err error
	this.Data["categories"], err = models.GetAllCategory()
	if err != nil{
		log.Fatal(err)
	}
}