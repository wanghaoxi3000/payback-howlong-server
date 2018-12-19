package controllers

import (
	"howlong/models"
	"math"
	"time"
)

type dateInfo struct {
	NextBill string // 下一账单日
	NextPay  string // 下一还款日
	CurPay   string // 当日消费还款日

	WaitDay        int // 最长等待时间
	IntervalBill   int // 下一账账单日间隔
	IntervalPay    int // 下一账还款日间隔
	IntervalCurPay int // 当日消费还款日间隔
}

type creditSerializer struct {
	Id      int64
	Name    string `valid:"MinSize(1); MaxSize(20)"`
	BillDay int    `valid:"Range(1, 29)"`
	PayDay  int    `valid:"Range(1, 28)"`
	PayFix  bool

	DateDetail *dateInfo
}

func (c *creditSerializer) Validate() (validInfo map[string]string, err error) {
	if validInfo, err = validateStruct(c); err != nil {
		return
	}

	// 还款日为间隔, 还款日和账单日间隔应当不大于28
	if !c.PayFix && c.PayDay > 28 {
		validInfo["Payday"] = "The sum of pay and bill day could less than 29"
	}

	return
}

// CreditDetail : Calculate credit card detail info
func (c *creditSerializer) CreditDetail(nowTime time.Time) {
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

	c.DateDetail = &dateInfo{
		billDate.Format("2006-01-02"),
		payDate.Format("2006-01-02"),
		curPayDate.Format("2006-01-02"),
		int(math.Ceil(waitDay.Hours() / 24)),
		int(math.Ceil(billDate.Sub(nowTime).Hours() / 24)),
		int(math.Ceil(payDate.Sub(nowTime).Hours() / 24)),
		int(math.Ceil(curPayDate.Sub(nowTime).Hours() / 24)),
	}
}

func (c *creditSerializer) serializer(credit *models.Credit) {
	c.Id = credit.Id
	c.Name = credit.Name
	c.BillDay = credit.BillDay
	c.PayDay = credit.PayDay
	c.PayFix = credit.PayFix

	c.CreditDetail(time.Now())
}

func (c *creditSerializer) unserializer() *models.Credit {
	creditModel := &models.Credit{
		Name:    c.Name,
		BillDay: c.BillDay,
		PayDay:  c.PayDay,
		PayFix:  c.PayFix,
	}

	return creditModel
}
