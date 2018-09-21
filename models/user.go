package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// User : User database model
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

// UpdateUserByOpenID : retrieves User by OpenID, Update session_key.
func UpdateUserByOpenID(openID string, seesion string) (obj *User, err error) {
	o := orm.NewOrm()
	obj = &User{OpenId: openID}

	if _, _, err = o.ReadOrCreate(obj, "OpenId"); err != nil {
		return nil, err
	}

	obj.SessionKey = seesion
	if _, err = o.Update(obj); err != nil {
		return nil, err
	}

	return
}
