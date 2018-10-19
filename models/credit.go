package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

// Credit : Credit card info orm struct
type Credit struct {
	Id      int64  `orm:"auto" form:"-"`
	Name    string `orm:"size(30)" form:"name" valid:"MinSize(1); MaxSize(20)"`
	BillDay int    `form:"billDay" valid:"Range(1, 29)"` // 29 代表月底最后一天
	PayDay  int    `form:"payDay" valid:"Range(1, 28)"`
	PayFix  bool   `form:"payFix"` // true: 固定账单日 false: 固定天数

	User *User `orm:"rel(fk);on_delete(cascade)"`
}

func init() {
	orm.RegisterModel(new(Credit))
}

// AddCredit insert a new Credit into database and returns
// last inserted Id on success.
func AddCredit(m *Credit) (id int64, err error) {
	o := orm.NewOrm()
	cnt, err := o.QueryTable(m).Filter("User", m.User.Id).Filter("Name", m.Name).Count()
	if err != nil {
		return
	}
	if cnt > 0 {
		err = fmt.Errorf("%s 已存在", m.Name)
		return 0, err
	}

	id, err = o.Insert(m)
	return
}

// GetUserCreditByID : Retrieves user credit by ID. Returns error if
// Id doesn't exist
func GetUserCreditByID(user *User, id int64) (v *Credit, err error) {
	o := orm.NewOrm()
	v = &Credit{Id: id}

	if err = o.QueryTable(v).Filter("User", user.Id).Filter("Id", id).One(v); err == nil {
		return v, nil
	}

	return nil, err
}

// GetUserAllCredit : Get user's all credit and return err if fail
func GetUserAllCredit(user *User, v *[]*Credit) (num int64, err error) {
	o := orm.NewOrm()
	num, err = o.QueryTable("credit").Filter("User", user.Id).RelatedSel().All(v)
	if err == orm.ErrNoRows {
		return 0, nil
	}
	return
}

// GetAllCredit retrieves all Credit matches certain condition. Returns empty list if
// no records exist
func GetAllCredit(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Credit))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Credit
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateUserCredit : Updates a user's credit and returns error if
// the record can not to be updated
func UpdateUserCredit(m *Credit) (err error) {
	o := orm.NewOrm()
	if exist := o.QueryTable(m).Exclude("Id", m.Id).Filter("User", m.User.Id).Filter("Name", m.Name).Exist(); exist == true {
		return fmt.Errorf("%s 已存在", m.Name)
	}
	_, err = o.Update(m)
	return
}

// DeleteCredit : Deletes Credit and returns error if
// the record can not to be deleted
func DeleteCredit(credit *Credit) (err error) {
	o := orm.NewOrm()
	_, err = o.Delete(credit)
	return
}
