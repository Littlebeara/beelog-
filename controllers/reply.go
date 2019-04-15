package controllers
//136 1748 2207
import (
	"log"
	"beelog/models"
	"github.com/astaxie/beego"
)

type ReplyController struct{
	beego.Controller
}

func (this *ReplyController)Add(){
	tid := this.Input().Get("tid")
	err := models.AddReply(tid, this.Input().Get("nickname"), this.Input().Get("content"))
	if err != nil{
		log.Fatal(err)
	}
	
	this.Redirect("/topic/views/" + tid, 302)
	
}

func (this ReplyController)Delete(){
	if !checkAccount(this.Ctx){
		return 
	}
	
	tid := this.Input().Get("tid")
	err := models.DeleteReply(this.Input().Get("rid"))
	if err != nil{
		log.Fatal(err)
	}
	
	this.Redirect("/topic/view/" + tid, 302)
}