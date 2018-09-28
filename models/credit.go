package models

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type dateInfo struct {
	NextBill string // 下一账单日
	NextPay  string // 下一还款日
	CurPay   string // 当日消费还款日

	WaitDay      int // 最长等待时间
	IntervalBill int // 账单日间隔
	IntervalPay  int // 还款日间隔
}

// Credit : Credit card info orm struct
type Credit struct {
	Id      int64  `orm:"auto" form:"-"`
	Name    string `orm:"size(30)" form:"name" valid:"MinSize(1); MaxSize(20)"`
	BillDay int    `form:"billDay" valid:"Range(1, 29)"` // 29 代表月底最后一天
	PayDay  int    `form:"payDay" valid:"Range(1, 28)"`
	PayFix  bool   `form:"payFix"` // true: 固定账单日 false: 固定天数

	User *User `orm:"rel(fk);on_delete(cascade)"`

	DateInfo *dateInfo `orm:"-"`
}

func init() {
	orm.RegisterModel(new(Credit))
}

// Validate : Validate credit struct data
func (c *Credit) Validate() (map[string]string, error) {
	validInfo := make(map[string]string)
	valid := validation.Validation{}

	if info, err := valid.Valid(c); err != nil {
		return nil, err
	} else if !info {
		for _, err := range valid.Errors {
			validInfo[err.Key] = err.Message
		}
	}

	// 还款日为间隔, 还款间隔和账单日应当不大于28
	if !c.PayFix && c.PayDay+c.BillDay > 28 {
		validInfo["Payday"] = "The sum of pay and bill day could less than 29"
	}

	return validInfo, nil
}

// CreditDetail : Calculate credit card detail info
func (c *Credit) CreditDetail(nowTime time.Time) {
	var billDate, payDate, curPayDate time.Time
	var waitDay time.Duration

	if c.BillDay > 28 {
		billDate = nowTime.AddDate(0, 1, -nowTime.Day())
	} else {
		billDate = nowTime.AddDate(0, 0, c.BillDay-nowTime.Day())
	}
	if nowTime.After(billDate) {
		billDate = billDate.AddDate(0, 1, 0)
	}

	if c.PayFix {
		payDate = nowTime.AddDate(0, 0, c.PayDay-nowTime.Day())
		waitDay = billDate.AddDate(0, 2, c.PayDay-billDate.Day()).Sub(billDate)
		if nowTime.After(payDate) {
			payDate = payDate.AddDate(0, 1, 0)
		}
		curPayDate = billDate.AddDate(0, 1, c.PayDay-billDate.Day())
	} else {
		payDate = billDate.AddDate(0, -1, c.PayDay)
		waitDay = billDate.AddDate(0, 1, c.PayDay).Sub(billDate)
		if nowTime.After(payDate) {
			payDate = billDate.AddDate(0, 0, c.PayDay)
		}
		curPayDate = billDate.AddDate(0, 0, c.PayDay)
	}

	c.DateInfo = &dateInfo{
		billDate.Format("2006-01-02"),
		payDate.Format("2006-01-02"),
		curPayDate.Format("2006-01-02"),
		int(math.Ceil(waitDay.Hours() / 24)),
		int(math.Ceil(billDate.Sub(nowTime).Hours() / 24)),
		int(math.Ceil(curPayDate.Sub(nowTime).Hours() / 24)),
	}
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

// GetSortedCredit : Get all sorted credit card info
func GetSortedCredit() (l []*Credit, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Credit))

	if _, err = qs.All(&l); err != nil {
		return nil, err
	}

	nowTime := time.Now()
	for _, v := range l {
		if v.DateInfo == nil {
			v.CreditDetail(nowTime)
		}
	}

	length := len(l)
	var sortFlag int
	for k := range l {
		sortFlag = k + 1
		for i := sortFlag; i < length; i++ {
			if l[k].DateInfo.IntervalPay < l[i].DateInfo.IntervalPay {
				l[k], l[i] = l[i], l[k]
			}
		}

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

// UpdateCreditById updates Credit by Id and returns error if
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
