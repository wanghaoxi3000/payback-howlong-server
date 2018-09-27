package controllers

type dateInfo struct {
	NextBill string // 下一账单日
	NextPay  string // 下一还款日
	CurPay   string // 当日消费还款日

	WaitDay      int // 最长等待时间
	IntervalBill int // 账单日间隔
	IntervalPay  int // 还款日间隔
}

type creditSerializer struct {
	Id      int64
	Name    string `valid:"MinSize(1); MaxSize(20)"`
	BillDay int    `valid:"Range(1, 29)"`
	PayDay  int    `valid:"Range(1, 28)"`
	PayFix  bool

	DateInfo *dateInfo
}

func (c *creditSerializer) Validate() (validInfo map[string]string, err error) {
	// validInfo := make(map[string]string)
	if validInfo, err = validateStruct(c); err != nil {
		return
	}

	// 还款日为间隔, 还款间隔和账单日应当不大于28
	if !c.PayFix && c.PayDay+c.BillDay > 28 {
		validInfo["Payday"] = "The sum of pay and bill day could less than 29"
	}

	return
}
