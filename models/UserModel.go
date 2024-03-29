package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// 用户表
type User struct {
	Id int64
	Username string `orm:"unique;size(32)"`
	Password string `orm:"size(32)"`
	Gender bool `orm:"default(true)"`
	Birthday time.Time `orm:"type(datetime);null;"`
	Phone_number string `orm:"size(11);null"`
	Site string `orm:"size(200);null"`
	Createtime time.Time `orm:"type(datetime);auto_now_add"`
	Role int `orm:"default(1)"`
}
//Lastlogintime time.Time `orm:"null;type(datetime)"`

func init() {
	orm.RegisterModel(new(User))
}
