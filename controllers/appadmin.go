package controllers

import (
	"acme3/models"
	pk "acme3/utilities/pbkdf2"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/twinj/uuid"
	"html/template"
	"strconv"
	"strings"
	//"time"
)

type AdminController struct {
	beego.Controller
}

func (this *AdminController) activeAdminContent(view string) {
	this.Layout = "admin-layout.tpl"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "header.tpl"
	this.LayoutSections["Footer"] = "footer.tpl"
	this.TplName = view + ".tpl"
	this.Data["domainname"] = "localhost:8080"
	//this.Data["domainname"] = "yourdomainname"
}

type compareform struct {
	Comparefield string `form:"comparefield"`
	Compareop    string `form:"compareop"`
	Compareval   string `form:"compareval" valid:"Required"`
}

func (this *AdminController) setCompare(query string) (orm.QuerySeter, bool) {
	o := orm.NewOrm()
	qs := o.QueryTable("auth_user")
	if this.Ctx.Input.Method() == "POST" {
		f := compareform{}
		if err := this.ParseForm(&f); err != nil {
			fmt.Println("cannot parse form")
			return qs, false
		}
		valid := validation.Validation{}
		if b, _ := valid.Valid(&f); !b {
			this.Data["Errors"] = valid.ErrorsMap
			return qs, false
		}
		if len(f.Compareop) >= 5 && f.Compareop[:5] == "__not" {
			qs = qs.Exclude(f.Comparefield+f.Compareop[5:], f.Compareval)
		} else {
			qs = qs.Filter(f.Comparefield+f.Compareop, f.Compareval)
		}
		this.Data["query"] = f.Comparefield + f.Compareop + "," + f.Compareval
	} else {
		str := strings.Split(query, ",")
		i := strings.Index(str[0], "__")
		if len(str[0][i:]) >= 5 && str[0][i:i+5] == "__not" {
			qs = qs.Exclude(str[0][:i]+str[0][i+5:], str[1])
		} else {
			qs = qs.Filter(str[0], str[1])
		}
		this.Data["query"] = query
	}
	return qs, true
}

func max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

func (this *AdminController) Index() {
	this.activeAdminContent("appadmin/index")

	defer func(this *AdminController) {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Index", r)
			this.Redirect("/home", 302)
		}
	}(this)

	const pagesize = 10
	parms := this.Ctx.Input.Param(":parms")
	this.Data["parms"] = parms
	str := strings.Split(parms, "!")
	fmt.Println("parms is", str)
	order := str[0]
	off, _ := strconv.Atoi(str[1])
	offset := int64(off)
	if offset < 0 {
		offset = 0
	}
	query := str[2]

	var users []*models.AuthUser
	rows := ""

	qs, ok := this.setCompare(query)
	if !ok {
		fmt.Println("cannot set QuerySeter")
		o := orm.NewOrm()
		qs := o.QueryTable("auth_user")
		qs = qs.Filter("id__gte", 0)
		this.Data["query"] = "id__gte,0"
	}

	count, _ := qs.Count()
	this.Data["count"] = count
	if offset >= count {
		offset = 0
	}
	num, err := qs.Limit(pagesize, offset).OrderBy(order).All(&users)
	if err != nil {
		fmt.Println("Query table failed:", err)
	}
	domainname := this.Data["domainname"]
	for x := range users {
		i := strings.Index(users[x].Reg_date.String(), " ")
		rows += fmt.Sprintf("<tr><td><a href='http://%s/appadmin/update/%s!%s'>%d</a></td>"+
			"<td>%s</td><td>%s</td><td>%s</td><td>%s...</td><td>%s</td><td>%s</td><td>%s</td></tr>", domainname, users[x].Email, parms,
			users[x].Id, users[x].First, users[x].Last, users[x].Email, users[x].Password[:20],
			users[x].Reg_key, users[x].Reg_date.String()[:i], users[x].Reset_key)
	}
	this.Data["Rows"] = template.HTML(rows)

	this.Data["order"] = order
	this.Data["offset"] = offset
	this.Data["end"] = max(0, count-pagesize)
	if num+offset < count {
		this.Data["next"] = num + offset
	}
	if offset-pagesize >= 0 {
		this.Data["prev"] = offset - pagesize
		this.Data["showprev"] = true
	} else if offset > 0 && offset < pagesize {
		this.Data["prev"] = 0
		this.Data["showprev"] = true
	}

	if count > pagesize {
		this.Data["ShowNav"] = true
	}
	this.Data["progress"] = float64(offset*100) / float64(max(count, 1))
}

func (this *AdminController) Add() {
	this.activeAdminContent("appadmin/add")

	parms := this.Ctx.Input.Param(":parms")
	this.Data["parms"] = parms

	if this.Ctx.Input.Method() == "POST" {
		flash := beego.NewFlash()
		u := models.AuthUser{}
		if err := this.ParseForm(&u); err != nil {
			fmt.Println("cannot parse form")
			return
		}
		this.Data["User"] = u
		valid := validation.Validation{}
		if b, _ := valid.Valid(&u); !b {
			this.Data["Errors"] = valid.ErrorsMap
			return
		}
		h := pk.HashPassword(u.Password)

		//******** Save user info to database
		o := orm.NewOrm()
		o.Using("default")

		// Convert password hash to string
		u.Password = hex.EncodeToString(h.Hash) + hex.EncodeToString(h.Salt)

		// Add user to database with new uuid
		key := uuid.NewV4()
		u.Reg_key = key.String()
		_, err := o.Insert(&u)
		this.Data["User"] = u
		if err != nil {
			flash.Error(u.Email + " already registered")
			flash.Store(&this.Controller)
			return
		}

		flash.Notice("User added")
		flash.Store(&this.Controller)
	}
}

type authUser struct {
	Id        int    `form:"id"`
	First     string `form:"first" valid:"Required"`
	Last      string `form:"last"`
	Email     string `form:"email" valid:"Email"`
	Password  string `form:"password"`
	Reg_key   string `form:"reg_key"`
	Reg_date  string `form:"reg_date"` // ParseForm cannot deal with time.Time in the form definition
	Reset_key string `form:"reset_key"`
	Delete    string `form:"delete,checkbox"`
}

func (this *AdminController) Update() {
	this.activeAdminContent("appadmin/update")

	defer func(this *AdminController) {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Update", r)
			this.Redirect("/home", 302)
		}
	}(this)

	flash := beego.NewFlash()

	str := this.Ctx.Input.Param(":username")
	i := strings.Index(str, "!")
	username := str[:i]
	this.Data["parms"] = str[i+1:]
	o := orm.NewOrm()
	o.Using("default")
	user := models.AuthUser{Email: username}
	if err := o.Read(&user, "Email"); err != nil {
		flash.Error("Internal error")
		flash.Store(&this.Controller)
		return
	}

	if this.Ctx.Input.Method() == "POST" {
		u := authUser{}
		if err := this.ParseForm(&u); err != nil {
			fmt.Println("cannot parse form")
			return
		}
		this.Data["User"] = u
		valid := validation.Validation{}
		if b, _ := valid.Valid(&u); !b {
			this.Data["Errors"] = valid.ErrorsMap
			return
		}

		if u.Delete == "on" {
			fmt.Println("about to delete record...")
			_, err := o.Delete(&user)
			if err == nil {
				flash.Notice("Record deleted")
				flash.Store(&this.Controller)
				return
			} else {
				flash.Error("Internal error")
				flash.Store(&this.Controller)
				return
			}
		}

		//******** Save user info to database
		user.First = u.First
		user.Last = u.Last
		user.Email = u.Email
		user.Reg_key = u.Reg_key
		user.Reset_key = u.Reset_key

		o := orm.NewOrm()
		o.Using("default")

		// Update user record
		_, err := o.Update(&user)
		if err != nil {
			flash.Error("Update failed")
			flash.Store(&this.Controller)
			return
		}

		flash.Error("User updated")
		flash.Store(&this.Controller)
	} else {
		this.Data["User"] = user
	}
}
