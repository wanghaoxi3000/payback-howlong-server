package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id         int64  `orm:"auto"`
	OpenId     string `orm:"unique;index"`
	SessionKey string

	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`

	Credits []*Credit `orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(new(User))
}
