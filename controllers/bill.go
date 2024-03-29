package controllers

import (
	. "blueSupermarket/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type BillListController struct {
	beego.Controller
}

type BillAddController struct {
	beego.Controller
}

type BillAddDataController struct {
	beego.Controller
}

type BillViewController struct {
	beego.Controller
}

type BillDelController struct {
	beego.Controller
}

type BillUpdateController struct {
	beego.Controller
}

type BillUpdateDataController struct {
	beego.Controller
}

func (c *BillListController) BillList() {
	goodsName := c.GetString("goodsName")
	tiGong := c.GetString("tigong", "null")
	fuKuan := c.GetString("fukuan")

	o := orm.NewOrm()
	var billSlice []*Bill
	bills := o.QueryTable("bill").Filter("GoodsName__icontains", goodsName)
	if fuKuan != "" {
		isPay, _ := strconv.ParseBool(fuKuan)
		bills = bills.Filter("IsPay", isPay)
	}
	_, err := bills.OrderBy("-id").All(&billSlice)
	if err != nil {
		fmt.Println(err)
	}

	var providerMaps []orm.Params
	_, providerErr := o.QueryTable("provider").Values(&providerMaps, "Id","ProviderName")
	if providerErr != nil {
		fmt.Println(providerErr)
	}

	var dataSlice []map[string]interface{}
	for _, bVal := range billSlice {
		isContinue := false
		oneMap := make(map[string]interface{})
		for _, pVal := range providerMaps {
			if tiGong == "null" {
				if bVal.Provider.Id == 0 {
					oneMap["ProviderName"] = "无"
					continue
				}
				if pVal["Id"] == bVal.Provider.Id {
					oneMap["ProviderName"] = pVal["ProviderName"]
				}
			} else {
				tiGongInt, _ := strconv.ParseInt(tiGong, 10, 64)
				if bVal.Provider.Id == tiGongInt{
					if tiGongInt == 0 {
						oneMap["ProviderName"] = "无"
						break
					}
					if pVal["Id"] == bVal.Provider.Id {
						oneMap["ProviderName"] = pVal["ProviderName"]
					}
				} else {
					isContinue = true
					continue
				}
			}
		}
		if isContinue{
			continue
		}
		oneMap["Id"] = bVal.Id
		oneMap["GoodsNumber"] = bVal.GoodsNumber
		oneMap["GoodsName"] = bVal.GoodsName
		oneMap["Amount"] = fmt.Sprintf("%.2f",bVal.Amount)
		oneMap["IsPay"] = bVal.IsPay
		createTime := bVal.CreateTime
		createTimeFormat := time.Unix(createTime.Unix(),0).Format("2006-01-02")
		oneMap["CreateTime"] = createTimeFormat
		dataSlice = append(dataSlice, oneMap)
	}
	c.Data["bill"] = dataSlice
	c.Data["provider"] = providerMaps
	c.Data["username"] = c.GetSession("user")
	c.TplName = "blueTpl/billList.html"
}

func (c *BillAddController) BillAdd() {
	var providerMaps []orm.Params
	o := orm.NewOrm()
	_, err := o.QueryTable("provider").Values(&providerMaps, "Id", "ProviderName")
	if err == nil {
		c.Data["provider"] = providerMaps
	} else {
		c.Data["provider"] = nil
	}
	c.Data["username"] = c.GetSession("user")
	c.TplName = "blueTpl/billAdd.html"
}

func (c *BillAddDataController) BillAddData() {
	billId := c.GetString("billId")
	billName := c.GetString("billName")
	billCom := c.GetString("billCom")
	billNum := c.GetString("billNum")
	goodsCount,err := strconv.ParseUint(billNum, 10, 32)
	if err != nil {
		c.Data["json"] = false
		c.ServeJSON()
		return
	}
	money := c.GetString("money")
	moneyFloat64, err := strconv.ParseFloat(money, 64)
	if err != nil {
		c.Data["json"] = false
		c.ServeJSON()
		return
	}
	account, err := strconv.ParseFloat(fmt.Sprintf("%.2f", moneyFloat64),32)
	if err != nil {
		c.Data["json"] = false
		c.ServeJSON()
		return
	}
	supplier := c.GetString("supplier")
	providerId, err := strconv.ParseInt(supplier,10,64)
	if err != nil {
		c.Data["json"] = false
		c.ServeJSON()
		return
	}
	zhiFu := c.GetString("zhifu")
	isPay,err := strconv.ParseBool(zhiFu)
	if err != nil {
		c.Data["json"] = false
		c.ServeJSON()
		return
	}

	var bill Bill
	bill.GoodsName = billName
	bill.GoodsNumber = billId
	bill.GoodsUnit = billCom
	bill.GoodsCount = uint32(goodsCount)
	bill.Amount = account
	bill.IsPay = isPay
	bill.Provider = &Provider{Id:providerId}

	o := orm.NewOrm()
	id, err := o.Insert(&bill)
	if id == 0 {
		fmt.Println(err)
		c.Data["json"] = false
	} else {
		c.Data["json"] = true
	}
	c.ServeJSON()
}

func (c *BillViewController) BillView() {
	var billSlice []*Bill
	id := c.GetString("id")
	int64Id, _ := strconv.ParseInt(id,10,64)
	o := orm.NewOrm()
	_, err := o.QueryTable("bill").Filter("id", int64Id).Limit(1).All(&billSlice)
	if err != nil {
		fmt.Println(err)
		c.Data["oneMap"] = nil
	} else {
		var providerMaps []orm.Params
		_, providerErr := o.QueryTable("provider").Values(&providerMaps,"Id", "ProviderName")
		if providerErr != nil {
			fmt.Println(providerErr)
		}

		oneMap := make(map[string]interface{})
		for _, bVal := range billSlice {
			oneMap["GoodsNumber"] = bVal.GoodsNumber
			oneMap["GoodsName"] = bVal.GoodsName
			oneMap["GoodsUnit"] = bVal.GoodsUnit
			oneMap["GoodsCount"] = bVal.GoodsCount
			oneMap["Amount"] = fmt.Sprintf("%.2f",bVal.Amount)
			if bVal.Provider.Id == 0 {
				oneMap["ProviderName"] = "无"
			} else {
				for _, pVal := range providerMaps {
					if pVal["Id"] == bVal.Provider.Id {
						oneMap["ProviderName"] = pVal["ProviderName"]
					}
				}
			}
			oneMap["IsPay"] = bVal.IsPay
		}
		c.Data["oneMap"] = oneMap
	}
	c.Data["username"] = c.GetSession("user")
	c.TplName = "blueTpl/billView.html"
}

func (c *BillDelController) BillDel() {
	id := c.GetString("id")
	int64Id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Data["json"] = false
		c.ServeJSON()
		return
	}
	o := orm.NewOrm()
	_, delErr := o.Delete(&Bill{Id:int64Id})
	if delErr != nil {
		c.Data["json"] = false
		c.ServeJSON()
		return
	}
	c.Data["json"] = true
	c.ServeJSON()
}

func (c *BillUpdateController) BillUpdate() {
	var billSlice []*Bill
	id := c.GetString("id")
	int64Id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	o := orm.NewOrm()
	var providerMaps []orm.Params
	_, providerErr := o.QueryTable("provider").Values(&providerMaps, "Id", "ProviderName")
	if providerErr == nil {
		c.Data["provider"] = providerMaps
	} else {
		c.Data["provider"] = nil
	}

	_, queryErr := o.QueryTable("bill").Filter("id",int64Id).Limit(1).All(&billSlice)
	if queryErr != nil {
		fmt.Println(queryErr)
		c.Data["oneMap"] = nil
	} else {
		oneMap := make(map[string]interface{})
		for _, bVal := range billSlice {
			oneMap["Id"] = bVal.Id
			oneMap["GoodsName"] = bVal.GoodsName
			oneMap["GoodsNumber"] = bVal.GoodsNumber
			oneMap["GoodsUnit"] = bVal.GoodsUnit
			oneMap["GoodsCount"] = bVal.GoodsCount
			oneMap["Amount"] = fmt.Sprintf("%.2f",bVal.Amount)
			oneMap["ProviderId"] = bVal.Provider.Id
			oneMap["IsPay"] = bVal.IsPay
		}
		c.Data["oneMap"] = oneMap
	}
	c.Data["username"] = c.GetSession("user")
	c.TplName = "blueTpl/billUpdate.html"
}

func (c *BillUpdateDataController) BillUpdateData() {
	id := c.GetString("id")
	intId, idErr := strconv.ParseInt(id, 10, 64)
	if idErr != nil {
		c.Data["json"] = false
		c.ServeJSON()
		return
	}
	goodsId := c.GetString("providerId")
	goodsName := c.GetString("providerName")
	goodsUnit := c.GetString("people")
	billNum := c.GetString("phone")
	goodsCount,err := strconv.ParseUint(billNum, 10, 32)
	if err != nil {
		c.Data["json"] = false
		c.ServeJSON()
		return
	}
	amount := c.GetString("Amount")
	amountFloat64, amountErr := strconv.ParseFloat(amount,64)
	if amountErr != nil {
		c.Data["json"] = false
		c.ServeJSON()
		return
	}
	supplier := c.GetString("supplier")
	providerId, providerErr := strconv.ParseInt(supplier,10, 64)
	if providerErr != nil {
		c.Data["json"] = false
		c.ServeJSON()
		return
	}
	zhiFu := c.GetString("zhifu")
	isPay,zhiFuErr := strconv.ParseBool(zhiFu)
	if zhiFuErr != nil {
		c.Data["json"] = false
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	bill := Bill{Id:intId}
	if o.Read(&bill) == nil {
		bill.GoodsNumber = goodsId
		bill.GoodsName = goodsName
		bill.GoodsUnit = goodsUnit
		bill.GoodsCount = uint32(goodsCount)
		bill.IsPay = isPay
		bill.Amount = amountFloat64
		bill.Provider = &Provider{Id:providerId}
		_,updateErr := o.Update(&bill)
		if updateErr == nil {
			c.Data["json"] = true
			c.ServeJSON()
			return
		}
	}
	c.Data["json"] = false
	c.ServeJSON()
}
