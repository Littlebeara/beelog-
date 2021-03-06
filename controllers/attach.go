package controllers

import (
	"io"
	"os"
	"net/url"
	"github.com/astaxie/beego"
)

type AttachController struct{
	beego.Controller
	
}

func (this *AttachController)Get(){
	filepath, err := url.QueryUnescape(this.Ctx.Request.RequestURI[1:])
	if err != nil{
		this.Ctx.WriteString(err.Error())
		return
	}
	
	f, err := os.Open(filepath)
	if err != nil{
		this.Ctx.WriteString(err.Error())
		return
	}
	defer f.Close()
	
	_, err = io.Copy(this.Ctx.ResponseWriter, f)
	if err != nil{
		this.Ctx.WriteString(err.Error())
		return
	}
}
















