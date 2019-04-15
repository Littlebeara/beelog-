package controllers

import (
	"path"
	"strings"
	"beelog/models"
	"log"

	"github.com/astaxie/beego"
)

type TopicController struct {
	beego.Controller
}

func (this *TopicController)Get(){
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	this.Data["IsTopic"] = true
	this.TplName = "topic.html"
	topics, err := models.GetTopics("", "", false)
	if err != nil{
		log.Fatal(err)
	}else{
		this.Data["Topics"] = topics
	}
}

func(this *TopicController)Post(){
	if !checkAccount(this.Ctx){
		this.Redirect("/login",302)
		return
	}
	tid := this.Input().Get("tid")
	title := this.Input().Get("title")
	content := this.Input().Get("content")
	category := this.Input().Get("category")
	lable := this.Input().Get("label")
	
	_, fh, err := this.GetFile("attachment")
	if err != nil{
		log.Fatal(err)
	}
	var attachment string
	if fh != nil{
		//保存附件
		attachment = fh.Filename
		//beego.Info(attachment)
		err = this.SaveToFile("attachment", path.Join("attachment", attachment))
		if err != nil{
			log.Fatal(err)
		}
	}              
	if len(tid) == 0{
		err = models.AddTopic(title, category, lable, content, attachment)
	}else{
		err = models.ModifyTopic(tid, title, category, lable, content, attachment)
	}
	
	if err!= nil{
		log.Fatal(err)
	}
}

func(this *TopicController)Add(){
	this.TplName = "topic_add.html"
	
}

func (this *TopicController)View(){
	this.TplName = "topic_view.html"
	reqUrl := this.Ctx.Request.RequestURI
	i := strings.LastIndex("reqUrl", "/")
	tid := reqUrl[i+1:]
//	a := this.Ctx.Input.Params()
	topic, err := models.GetonrTopics(tid)
	if err != nil{
		log.Fatal(err)
		this.Redirect("/",302)
		return
	}
	this.Data["topic"] = topic
	this.Data["Lables"] = strings.Split(topic.Lables, " ")
	
	replies, err := models.GetAllReplies(tid)
	if err != nil{
		log.Fatal(err)
		return
	}
	this.Data["replies"] = replies
	this.Data["IsLogin"] = checkAccount(this.Ctx)
}

func (this *TopicController)Modify(){
	this.TplName = "topic_modify.html"
	tid := this.Input().Get("tid")
	topic, err := models.GetonrTopics(tid)
	if err != nil{
		log.Fatal(err)
		this.Redirect("/",302)
	}
	this.Data["topic"] = topic
	this.Data["tid"] = tid
	
}

func (this *TopicController)Delete(){
	if !checkAccount(this.Ctx){
		this.Redirect("/login",302)
		return
	}
	err := models.DelTopic(this.Input().Get("tid"))
	if err != nil{
		log.Fatal(err)
		
	}
	this.Redirect("/",302)
}















