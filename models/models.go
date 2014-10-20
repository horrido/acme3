package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type AuthUser struct {
	Id        int    `form:"id"`
	First     string `form:"first" valid:"Required"`
	Last      string `form:"last"`
	Email     string `form:"email" valid:"Email" orm:"unique"`
	Password  string `form:"password" valid:"MinSize(6);MaxSize(30)"`
	Reg_key   string
	Reg_date  time.Time `orm:"auto_now_add;type(datetime)"`
	Reset_key string
}

func init() {
	orm.RegisterModel(new(AuthUser))
}
