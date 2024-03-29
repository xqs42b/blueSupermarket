package controllers

import (
	. "blueSupermarket/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type UserListController struct {
	beego.Controller
}

type UserAddController struct {
	beego.Controller
}

type UserAddDataController struct {
	beego.Controller
}

type UserDelController struct {
	beego.Controller
}

type UserUpdateController struct {
	beego.Controller
}

type UserUpdateDataController struct {
	beego.Controller
}

type UserViewController struct {
	beego.Controller
}

func (c *UserListController) UserList() {
	var userMaps []orm.Params
	o := orm.NewOrm()
	username := c.GetString("username")
	_, err := o.QueryTable("user").Filter("username__icontains", username).OrderBy("-id").Values(&userMaps)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, userMap := range userMaps {
			birthday := userMap["Birthday"]
			if birthday != nil {
				bornYear := birthday.(time.Time).Year()
				currentYear := time.Now().Year()
				userMap["Age"] = currentYear - bornYear
			} else {
				userMap["Age"] = 0
			}
		}
	}
	c.Data["user"] = &userMaps
	c.Data["username"] = c.GetSession("user")
	c.TplName = "blueTpl/userList.html"
}

func (c *UserAddController) UserAdd() {
	c.Data["username"] = c.GetSession("user")
	c.TplName = "blueTpl/userAdd.html"
}

func (c *UserAddDataController) UserAddData() {
	isSuccess := true
	username := c.GetString("userName")
	password := c.GetString("password")
	gender := c.GetString("gender")
	phone := c.GetString("phone")
	birthday := c.GetString("birthday")
	timeDate,err := time.Parse("2006-01-02", birthday)
	if err != nil {
		fmt.Println(err)
		timeDate = time.Now()
	}

	resGender := false
	if gender == "man" {
		resGender = true
	}
	address := c.GetString("address")
	userLei := c.GetString("userlei")
	role,err := strconv.Atoi(userLei)
	if err != nil {
		role = 3
	}

	o := orm.NewOrm()
	var user User
	user.Username = username
	user.Password = password
	user.Gender = resGender
	user.Phone_number = phone
	user.Birthday = timeDate
	user.Site = address
	user.Role = role

	id, err:= o.Insert(&user)
	if id == 0{
		fmt.Println(err)
		isSuccess = false
	}
	c.Data["json"] = isSuccess
	c.ServeJSON()
}

func (c *UserDelController) UserDel() {
	id := c.GetString("id")
	intId,err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Data["json"] = false
		c.ServeJSON()
		return
	}
	o := orm.NewOrm()
	if _,err := o.Delete(&User{Id:intId});err == nil {
		c.Data["json"] = true
		c.ServeJSON()
		return
	}
	c.Data["json"] = false
	c.ServeJSON()
}

func (c *UserUpdateController) UserUpdate() {
	var userMaps []orm.Params
	o := orm.NewOrm()
	id := c.GetString("id")
	_, err := o.QueryTable("user").Filter("Id", id).Limit(1).Values(&userMaps)
	if err != nil {
		fmt.Println(err)
		c.Data["user"] = nil
	} else {
		for _, userMap := range userMaps {
			if userMap["Birthday"] != nil {
				birthday := userMap["Birthday"].(time.Time)
				datetime := time.Unix(birthday.Unix(), 0).Format("2006-01-02")
				userMap["Birthday"] = datetime
			} else {
				userMap["Birthday"] = nil
			}
		}
		c.Data["user"] = userMaps
	}
	c.Data["username"] = c.GetSession("user")
	c.TplName = "blueTpl/userUpdate.html"
}

func (c *UserUpdateDataController) UserUpdateData() {
	id := c.GetString("id")
	int64Id, err := strconv.ParseInt(id,10,64)
	if err != nil {
		c.Data["json"] = false
		c.ServeJSON()
		return
	}
	username := c.GetString("userName")
	gender := c.GetString("gender")
	phone := c.GetString("phone")
	birthday := c.GetString("birthday")
	timeDate, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		timeDate = time.Now()
	}
	resGender := false
	if gender == "man" {
		resGender = true
	}
	address := c.GetString("address")
	userLei := c.GetString("userlei")
	role, err := strconv.Atoi(userLei)
	if err != nil {
		role = 3
	}

	o := orm.NewOrm()
	user := User{Id:int64Id}
	if o.Read(&user) == nil {
		user.Username = username
		user.Gender = resGender
		user.Phone_number = phone
		user.Birthday = timeDate
		user.Site = address
		user.Role = role
		_,err := o.Update(&user)
		if err == nil {
			c.Data["json"] = true
			c.ServeJSON()
			return
		}
	}
	c.Data["json"] = false
	c.ServeJSON()
}

func (c *UserViewController) UserView() {
	var userMaps []orm.Params
	o := orm.NewOrm()
	id := c.GetString("id")
	_, err := o.QueryTable("user").Filter("Id", id).Limit(1).Values(&userMaps)
	if err != nil {
		fmt.Println(err)
		c.Data["user"] = nil
	} else {
		for _, userMap := range userMaps {
			if userMap["Birthday"] != nil {
				birthday := userMap["Birthday"].(time.Time)
				datetime := time.Unix(birthday.Unix(), 0).Format("2006-01-02")
				userMap["Birthday"] = datetime
			} else {
				userMap["Birthday"] = nil
			}
		}
		c.Data["user"] = userMaps
	}

	c.Data["username"] = c.GetSession("user")
	c.TplName = "blueTpl/userView.html"
}
