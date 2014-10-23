package controllers

import (
	"acme3/models"
	pk "acme3/utilities/pbkdf2"
	"encoding/hex"
	"fmt"
	"github.com/alexcesaro/mail/gomail"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/twinj/uuid"
	"strings"
	"time"
)

func (this *MainController) Login() {
	this.activeContent("user/login")

	sess := this.GetSession("acme")
	if sess != nil {
		this.Redirect("/home", 302)
	}

	back := strings.Replace(this.Ctx.Input.Param(":back"), ">", "/", -1) // allow for deeper URL such as l1/l2/l3 represented by l1>l2>l3
	fmt.Println("back is", back)
	if this.Ctx.Input.Method() == "POST" {
		flash := beego.NewFlash()
		email := this.GetString("email")
		password := this.GetString("password")
		valid := validation.Validation{}
		valid.Email(email, "email")
		valid.Required(password, "password")
		if valid.HasErrors() {
			errormap := make(map[string]string)
			for _, err := range valid.Errors {
				errormap[err.Key] = err.Message
			}
			this.Data["Errors"] = errormap
			return
		}
		fmt.Println("Authorization is", email, ":", password)

		//******** Read password hash from database
		var x pk.PasswordHash

		x.Hash = make([]byte, 32)
		x.Salt = make([]byte, 16)

		o := orm.NewOrm()
		o.Using("default")
		user := models.AuthUser{Email: email}
		err := o.Read(&user, "Email")
		if err == nil {
			if user.Reg_key != "" {
				flash.Error("Account not verified")
				flash.Store(&this.Controller)
				return
			}

			// scan in the password hash/salt
			fmt.Println("Password to scan:", user.Password)
			if x.Hash, err = hex.DecodeString(user.Password[:64]); err != nil {
				fmt.Println("ERROR:", err)
			}
			if x.Salt, err = hex.DecodeString(user.Password[64:]); err != nil {
				fmt.Println("ERROR:", err)
			}
			fmt.Println("decoded password is", x)
		} else {
			flash.Error("No such user/email")
			flash.Store(&this.Controller)
			return
		}

		//******** Compare submitted password with database
		if !pk.MatchPassword(password, &x) {
			flash.Error("Bad password")
			flash.Store(&this.Controller)
			return
		}

		//******** Create session and go back to previous page
		m := make(map[string]interface{})
		m["first"] = user.First
		m["username"] = email
		m["timestamp"] = time.Now()
		this.SetSession("acme", m)
		this.Redirect("/"+back, 302)
	}
}

func (this *MainController) Logout() {
	this.activeContent("logout")
	this.DelSession("acme")
	this.Redirect("/home", 302)
}

type user1 struct {
	First    string `form:"first" valid:"Required"`
	Last     string `form:"last"`
	Email    string `form:"email" valid:"Email"`
	Password string `form:"password" valid:"MinSize(6)"`
	Confirm  string `form:"password2" valid:"Required"`
}

func (this *MainController) Register() {
	this.activeContent("user/register")

	if this.Ctx.Input.Method() == "POST" {
		flash := beego.NewFlash()
		u := user1{}
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
		if u.Password != u.Confirm {
			flash.Error("Passwords don't match")
			flash.Store(&this.Controller)
			return
		}
		h := pk.HashPassword(u.Password)

		//******** Save user info to database
		o := orm.NewOrm()
		o.Using("default")

		user := models.AuthUser{First: u.First, Last: u.Last, Email: u.Email}

		// Convert password hash to string
		user.Password = hex.EncodeToString(h.Hash) + hex.EncodeToString(h.Salt)

		// Add user to database with new uuid and send verification email
		key := uuid.NewV4()
		user.Reg_key = key.String()
		_, err := o.Insert(&user)
		if err != nil {
			flash.Error(u.Email + " already registered")
			flash.Store(&this.Controller)
			return
		}

		domainname := this.Data["domainname"]
		if !sendVerification(u.Email, key.String(), domainname.(string)) {
			flash.Error("Unable to send verification email")
			flash.Store(&this.Controller)
			return
		}
		flash.Notice("Your account has been created. You must verify the account in your email.")
		flash.Store(&this.Controller)
		this.Redirect("/notice", 302)
	}
}

func sendVerification(email, u string, domainname string) bool {
	link := "http://" + domainname + "/user/verify/" + u
	host := "smtp.gmail.com"
	port := 587
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", "acmecorp@gmail.com", "ACME Corporation")
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Account Verification for ACME Corporation")
	msg.SetBody("text/html", "To verify your account, please click on the link: <a href=\""+link+
		"\">"+link+"</a><br><br>Best Regards,<br>ACME Corporation")
	m := gomail.NewMailer(host, "youraccount@gmail.com", "YourPassword", port)
	if err := m.Send(msg); err != nil {
		return false
	}
	return true
}

func (this *MainController) Verify() {
	this.activeContent("user/verify")

	u := this.Ctx.Input.Param(":uuid")
	o := orm.NewOrm()
	o.Using("default")
	user := models.AuthUser{Reg_key: u}
	err := o.Read(&user, "Reg_key")
	if err == nil {
		this.Data["Verified"] = 1
		user.Reg_key = ""
		if _, err := o.Update(&user); err != nil {
			delete(this.Data, "Verified")
		}
	}
}

type user2 struct {
	First   string `form:"first" valid:"Required"`
	Last    string `form:"last"`
	Email   string `form:"email" valid:"Email"`
	Current string `form:"current" valid:"Required"`
}

func (this *MainController) Profile() {
	this.activeContent("user/profile")

	//******** This page requires login
	sess := this.GetSession("acme")
	if sess == nil {
		this.Redirect("/user/login/home", 302)
		return
	}
	m := sess.(map[string]interface{})

	flash := beego.NewFlash()

	//******** Read password hash from database
	var x pk.PasswordHash

	x.Hash = make([]byte, 32)
	x.Salt = make([]byte, 16)

	o := orm.NewOrm()
	o.Using("default")
	user := models.AuthUser{Email: m["username"].(string)}
	err := o.Read(&user, "Email")
	if err == nil {
		// scan in the password hash/salt
		if x.Hash, err = hex.DecodeString(user.Password[:64]); err != nil {
			fmt.Println("ERROR:", err)
		}
		if x.Salt, err = hex.DecodeString(user.Password[64:]); err != nil {
			fmt.Println("ERROR:", err)
		}
	} else {
		flash.Error("Internal error")
		flash.Store(&this.Controller)
		return
	}

	if this.Ctx.Input.Method() == "POST" {
		u := user2{}
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

		password := this.GetString("password")
		password2 := this.GetString("password2")
		if password != "" {
			valid.MinSize(password, 6, "password")
			valid.Required(password2, "password2")
			if valid.HasErrors() {
				errormap := make(map[string]string)
				for _, err := range valid.Errors {
					errormap[err.Key] = err.Message
				}
				this.Data["Errors"] = errormap
				return
			}

			if password != password2 {
				flash.Error("Passwords don't match")
				flash.Store(&this.Controller)
				return
			}
			h := pk.HashPassword(password)

			// Convert password hash to string
			user.Password = hex.EncodeToString(h.Hash) + hex.EncodeToString(h.Salt)
		}

		//******** Compare submitted password with database
		if !pk.MatchPassword(u.Current, &x) {
			flash.Error("Bad current password")
			flash.Store(&this.Controller)
			return
		}

		//******** Save user info to database
		user.First = u.First
		user.Last = u.Last
		user.Email = u.Email

		_, err := o.Update(&user)
		if err == nil {
			flash.Notice("Profile updated")
			flash.Store(&this.Controller)
			m["username"] = u.Email
		} else {
			flash.Error("Internal error")
			flash.Store(&this.Controller)
			return
		}
	} else {
		this.Data["User"] = user
	}
}

func (this *MainController) Remove() {
	this.activeContent("user/remove")

	//******** This page requires login
	sess := this.GetSession("acme")
	if sess == nil {
		this.Redirect("/user/login/home", 302)
		return
	}
	m := sess.(map[string]interface{})

	if this.Ctx.Input.Method() == "POST" {
		current := this.GetString("current")
		valid := validation.Validation{}
		valid.Required(current, "current")
		if valid.HasErrors() {
			errormap := make(map[string]string)
			for _, err := range valid.Errors {
				errormap[err.Key] = err.Message
			}
			this.Data["Errors"] = errormap
			return
		}

		flash := beego.NewFlash()

		//******** Read password hash from database
		var x pk.PasswordHash

		x.Hash = make([]byte, 32)
		x.Salt = make([]byte, 16)

		o := orm.NewOrm()
		o.Using("default")
		user := models.AuthUser{Email: m["username"].(string)}
		err := o.Read(&user, "Email")
		if err == nil {
			// scan in the password hash/salt
			if x.Hash, err = hex.DecodeString(user.Password[:64]); err != nil {
				fmt.Println("ERROR:", err)
			}
			if x.Salt, err = hex.DecodeString(user.Password[64:]); err != nil {
				fmt.Println("ERROR:", err)
			}
		} else {
			flash.Error("Internal error")
			flash.Store(&this.Controller)
			return
		}

		//******** Compare submitted password with database
		if !pk.MatchPassword(current, &x) {
			flash.Error("Bad current password")
			flash.Store(&this.Controller)
			return
		}

		//******** Delete user record
		_, err = o.Delete(&user)
		if err == nil {
			flash.Notice("Your account is deleted.")
			flash.Store(&this.Controller)
			this.DelSession("acme")
			this.Redirect("/notice", 302)
		} else {
			flash.Error("Internal error")
			flash.Store(&this.Controller)
			return
		}
	}
}

func (this *MainController) Forgot() {
	this.activeContent("user/forgot")

	if this.Ctx.Input.Method() == "POST" {
		email := this.GetString("email")
		valid := validation.Validation{}
		valid.Email(email, "email")
		if valid.HasErrors() {
			errormap := make(map[string]string)
			for _, err := range valid.Errors {
				errormap[err.Key] = err.Message
			}
			this.Data["Errors"] = errormap
			return
		}

		flash := beego.NewFlash()

		o := orm.NewOrm()
		o.Using("default")
		user := models.AuthUser{Email: email}
		err := o.Read(&user, "Email")
		if err != nil {
			flash.Error("No such user/email in our records")
			flash.Store(&this.Controller)
			return
		}

		u := uuid.NewV4()
		user.Reset_key = u.String()
		_, err = o.Update(&user)
		if err != nil {
			flash.Error("Internal error")
			flash.Store(&this.Controller)
			return
		}
		domainname := this.Data["domainname"]
		sendRequestReset(email, u.String(), domainname.(string))
		flash.Notice("You've been sent a reset password link. You must check your email.")
		flash.Store(&this.Controller)
		this.Redirect("/notice", 302)
	}
}

func sendRequestReset(email, u string, domainname string) bool {
	link := "http://" + domainname + "/user/reset/" + u
	host := "smtp.gmail.com"
	port := 587
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", "acmecorp@gmail.com", "ACME Corporation")
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Request Password Reset for ACME Corporation")
	msg.SetBody("text/html", "To reset your password, please click on the link: <a href=\""+link+
		"\">"+link+"</a><br><br>Best Regards,<br>ACME Corporation")
	m := gomail.NewMailer(host, "youraccount@gmail.com", "YourPassword", port)
	if err := m.Send(msg); err != nil {
		return false
	}
	return true
}

func (this *MainController) Reset() {
	this.activeContent("user/reset")

	flash := beego.NewFlash()

	u := this.Ctx.Input.Param(":uuid")
	o := orm.NewOrm()
	o.Using("default")
	user := models.AuthUser{Reset_key: u}
	err := o.Read(&user, "Reset_key")
	if err == nil {
		if this.Ctx.Input.Method() == "POST" {
			password := this.GetString("password")
			password2 := this.GetString("password2")
			valid := validation.Validation{}
			valid.MinSize(password, 6, "password")
			valid.Required(password2, "password2")
			if valid.HasErrors() {
				errormap := make(map[string]string)
				for _, err := range valid.Errors {
					errormap[err.Key] = err.Message
				}
				this.Data["Errors"] = errormap
				return
			}

			if password != password2 {
				flash.Error("Passwords don't match")
				flash.Store(&this.Controller)
				return
			}
			h := pk.HashPassword(password)

			// Convert password hash to string
			user.Password = hex.EncodeToString(h.Hash) + hex.EncodeToString(h.Salt)

			user.Reset_key = ""
			if _, err := o.Update(&user); err != nil {
				flash.Error("Internal error")
				flash.Store(&this.Controller)
				return
			}
			flash.Notice("Password updated.")
			flash.Store(&this.Controller)
			this.Redirect("/notice", 302)
		}
	} else {
		flash.Notice("Invalid key.")
		flash.Store(&this.Controller)
		this.Redirect("/notice", 302)
	}
}
