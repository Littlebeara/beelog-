package models


import (
	"strings"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// 设置数据库路径
	_DB_NAME = "data/beeblog.db"
	// 设置数据库名称
	_SQLITE3_DRIVER = "sqlite3"
)

// 分类
type Category struct {
	Id              int64
	Title           string
	Created         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}

// 文章
type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Category		string
	Lables			string
	Content         string `orm:"size(5000)"`
	Attachment      string
	Created         time.Time `orm:"index"`
	Updated         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time `orm:"index"`
	ReplyCount      int64
	ReplyLastUserId int64
}

type Commnet struct{
	Id		int64
	Tid  	int64
	Name    string
	Content string    `orm:"size(1000)"`
	Created time.Time `orm:index`
	
}
func RegisterDB() {
	// 检查数据库文件
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}

	// 注册模型
	orm.RegisterModel(new(Category), new(Topic), new(Commnet))
	// 注册驱动（“sqlite3” 属于默认注册，此处代码可省略）
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)
	// 注册默认数据库
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)
}

func AddCategory(name string)error{
	o := orm.NewOrm()
	
	cate := &Category{Title: name}
	
	qs := o.QueryTable("category")
	err := qs.Filter("title", name).One(cate)
	if err != nil{
		return err
	}
	
	//插入数据
	_, err = o.Insert(cate)
	if err != nil{
		return err
	}
	return nil
}

func DelCategory(id string)error{
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil{
		return err
	}
	o := orm.NewOrm()
	cate := &Category{Id: cid}
	_, err = o.Delete(cate)
	return err
}
//获取所有的category
func GetAllCategory()([]*Category, error){
	o := orm.NewOrm()
	
	cates := make([]*Category, 0)
	
	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

func AddTopic(title, category, label, content, attachment string)error{
	o := orm.NewOrm()
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"
	
	topic := &Topic{
		Title      : title,
		Category   : category,
		Lables	   : label,
		Attachment : attachment,
		Content    : content,
		Created    : time.Now(),
		Updated    : time.Now(),
	}
	
	_, err := o.Insert(topic)
	return err
	
	//更新分类统计
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("titile", category).One(cate)
	if err == nil{
		cate.TopicCount ++
		_, err = o.Update(cate)
	}
	return err
}

func GetTopics(cate string, label string, isDesc bool)(topics []*Topic, err  error){
	o := orm.NewOrm()

	topics = make([]*Topic, 0)

	qs := o.QueryTable("topic")
	if isDesc{
		if len(cate)> 0{
			qs = qs.Filter("category", cate)
		}
		if len(label) > 0{
			qs = qs.Filter("label_contains", "$"+label+"#")
		}
		_, err = qs.OrderBy("-created").All(&topics)
	}else{
		_, err = qs.All(&topics)
	}
	return topics, err
	
}

func GetonrTopics(tid string)(*Topic, error){
	tidnum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil{
		return nil, err
	}
	
	o := orm.NewOrm()
	topic := new(Topic)
	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidnum).One(topic)
	if err != nil{
		return nil, err
	}
	topic.Views ++
	_, err = o.Update(topic)
	
	topic.Lables = strings.Replace(strings.Replace(topic.Lables, "#"," ",-1), "$", "", -1)
	return topic, err
}

func ModifyTopic(tid, title, category, label, content, attachment string)error{
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"
	
	var Oldcate, OldAttach string
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		Oldcate = topic.Category
		OldAttach = topic.Attachment
		topic.Title = title
		topic.Attachment = attachment
		topic.Content = content
		topic.Lables = label
		topic.Category = category
		
		topic.Updated = time.Now()
		_, err := o.Update(topic)
		if err != nil{
			return err
		}
	}
	//更新分类统计
	if len(Oldcate) > 0{
		cate := new(Category)
		qs := o.QueryTable("category")
		err := qs.Filter("title", Oldcate).One(cate)
		if err == nil{
			cate.TopicCount --
			_, err = o.Update(cate)
			return err
		}	
		
	}
	
	//删除旧的附件
	
	if len(OldAttach) > 0{
		os.Remove(path.Join("attachment", OldAttach))
	}
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil{
		cate.TopicCount ++
		_, err = o.Update(cate)
			return err
	}
	return nil

}

func DelTopic(tid string)error{
	tidnum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil{
		return err
	}
	var oldCate string
	o := orm.NewOrm()
	topic := &Topic{Id: tidnum}
	if o.Read(topic) == nil{
		oldCate = topic.Category
		_, err = o.Delete(topic)
		if err != nil{
			return err
		}
	}
	
	if len(oldCate) >0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", oldCate).One(cate)
		if err == nil{
			cate.TopicCount --
			_, err = o.Update(cate)
		}
	}
	return err
}

func AddReply(tid, nickname, content string)error{
	tidnum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil{
		return err
	}
	reply := &Commnet{
		Tid 	: tidnum,
		Name	: nickname,
		Content : content,
		Created : time.Now(),
	}
	o := orm.NewOrm()
	 _, err = o.Insert(reply)
	if err != nil{
		return err
	}
	topic := &Topic{Id: tidnum}
	if o.Read(topic) == nil{
		topic.ReplyTime = time.Now()
		topic.ReplyCount ++
		_, err = o.Update(topic)
	}
	return err
	
}

func GetAllReplies(tid string)(replies []*Commnet, err error){
	tidnum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil{
		return nil, err
	}
	replies = make([]*Commnet, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidnum).All(&replies)
	return replies, err
}

func DeleteReply(rid string)error{
	ridum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil{
		return err
	}
	
	var tidnum int64
	o := orm.NewOrm()
	reply := &Commnet{Id: ridum}
	if o.Read(reply) ==nil {
		tidnum = reply.Id
		_, err = o.Delete(reply)
		if err != nil{
			return err
		}
		
	}
	replies := make([]*Commnet, 0)
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidnum).OrderBy("-created").All(&replies)//降序排序
	if err != nil{
		return err
	}
	
	topic := &Topic{Id: tidnum}
	if o.Read(topic) == nil{
		topic.ReplyTime = replies[0].Created
		topic.ReplyCount = int64(len(replies))
		_, err = o.Update(topic)
	}
	
	return err
	
}




