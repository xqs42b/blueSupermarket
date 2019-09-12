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

func (c *BillListController) BillList() {
	var billSlice []*Bill
	goodsName := c.GetString("goodsName")
	tiGong := c.GetString("tigong")
	fuKuan := c.GetString("fukuan")

	o := orm.NewOrm()
	bills := o.QueryTable("bill").Filter("GoodsName__icontains",goodsName)
	if tiGong != "0" {
		bills = bills.Filter("Provider__contains", tiGong)
	}
	if fuKuan != "" {
		isPay, _ := strconv.ParseBool(fuKuan)
		bills = bills.Filter("IsPay", isPay)
	}
	_, err := bills.RelatedSel("provider").OrderBy("-id").All(&billSlice)
	if err != nil {
		fmt.Println(err)
	}

	var dataSlice []map[string]interface{}
	for _, bVal := range billSlice {
		oneMap := make(map[string]interface{})
		oneMap["Id"] = bVal.Id
		oneMap["GoodsNumber"] = bVal.GoodsNumber
		oneMap["GoodsName"] = bVal.GoodsName
		oneMap["ProviderName"] = bVal.Provider.ProviderName
		oneMap["Amount"] = bVal.Amount
		oneMap["IsPay"] = bVal.IsPay
		createTime := bVal.CreateTime
		createTimeFormat := time.Unix(createTime.Unix(),0).Format("2006-01-02")
		oneMap["CreateTime"] = createTimeFormat
		dataSlice = append(dataSlice, oneMap)
	}

	var providerNameMaps []orm.Params
	num, _ := o.QueryTable("provider").Values(&providerNameMaps, "Id", "ProviderName")
	if num > 0 {
		c.Data["providerNames"] = providerNameMaps
	} else {
		c.Data["providerNames"] = nil
	}

	c.Data["bill"] = dataSlice
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
	_, err := o.QueryTable("bill").Filter("id", int64Id).RelatedSel("provider").Limit(1).All(&billSlice)
	if err != nil {
		fmt.Println(err)
		c.Data["oneMap"] = nil
	} else {
		oneMap := make(map[string]interface{})
		for _, bVal := range billSlice {
			oneMap["GoodsNumber"] = bVal.GoodsNumber
			oneMap["GoodsName"] = bVal.GoodsName
			oneMap["GoodsUnit"] = bVal.GoodsUnit
			oneMap["GoodsCount"] = bVal.GoodsCount
			oneMap["Amount"] = bVal.Amount
			oneMap["ProviderName"] = bVal.Provider.ProviderName
			oneMap["IsPay"] = bVal.IsPay
		}
		c.Data["oneMap"] = oneMap
	}
	c.TplName = "blueTpl/billView.html"
}
