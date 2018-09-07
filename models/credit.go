package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Credit struct {
	Id        int64  `orm:"auto" form:"-"`
	Name      string `orm:"size(30)" form:"name" valid:"MinSize(1); MaxSize(20)"`
	BillDay   int    `form:"billDay" valid:"Range(1, 28)"`
	Payday    int    `form:"payDay" valid:"Range(1, 28)"`
	PayFix    bool   `form:"payFix"`
	SameMonth bool   `form:"sameMonth"`
}

func init() {
	orm.RegisterModel(new(Credit))
}

// Validate : Validate credit struct data
func (c *Credit) Validate() (map[string]string, error) {
	validInfo := make(map[string]string)
	valid := validation.Validation{}

	if info, err := valid.Valid(c); err != nil {
		fmt.Println("Go to here 1111")
		return nil, err
	} else if !info {
		fmt.Println("Go to here 2222")
		for _, err := range valid.Errors {
			validInfo[err.Key] = err.Message
		}
	}

	return validInfo, nil
}

// ListCreditInfo : list all credit card info
func ListCreditInfo() []*Credit {
	o := orm.NewOrm()
	var credits []*Credit
	o.QueryTable("Credit").All(&credits)

	return credits
}

// AddCredit insert a new Credit into database and returns
// last inserted Id on success.
func AddCredit(m *Credit) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCreditById retrieves Credit by Id. Returns error if
// Id doesn't exist
func GetCreditById(id int64) (v *Credit, err error) {
	o := orm.NewOrm()
	v = &Credit{Id: id}
	if err = o.QueryTable(new(Credit)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
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

// UpdateCredit updates Credit by Id and returns error if
// the record to be updated doesn't exist
func UpdateCreditById(m *Credit) (err error) {
	o := orm.NewOrm()
	v := Credit{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCredit deletes Credit by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCredit(id int64) (err error) {
	o := orm.NewOrm()
	v := Credit{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Credit{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
