package test

import (
	"fmt"
	"game"
	"reflect"
	"strings"
	"testing"
)

type JS_RechargeRecord2 struct {
	Type      int   `json:"type"`
	Money     int   `json:"money"`
	Addgem    int   `json:"addgem"`
	ExtraGem  int   `json:"extragem"`
	BeforeGem int   `json:"befroegem"`
	AfterGem  int   `json:"aftergem"`
	Time      int64 `json:"time"`
	OrderId   int   `json:"orderid"`
	Isok      int   `json:"isok"`
}

//! 充值数据库
type UserRecharge struct {
	Uid         int64  `json:"uid"`
	Money       int    `json:"money"`
	Getallgem   int    `json:"getallgem"`
	Type1       int    `json:"type1"`
	Type2       int    `json:"type2"`
	Type3       int    `json:"type3"`
	Type4       int    `json:"type4"`
	Type5       int    `json:"type5"`
	Type6       int    `json:"type6"`
	Record      string `json:"record"`
	Firstaward  int    `json:"firstaward"`
	MoneyDay    int    `json:"moneyday"`
	MoneyWeek   int    `json:"moneyweek"`
	MonthCount1 int    `json:"monthcount1"`
	MonthCount2 int    `json:"monthcount2"`
	MonthCount3 int    `json:"monthcount3"`
	VipBox      int64  `json:"vipbox"`
	FundType    int64  `json:"fundtype"`
	FundGet     string `json:"fundget"`

	record  []JS_RechargeRecord2
	fundget []int64

	DataTest
}

type DataTest struct {
	baseData reflect.Value //! 原始数据
	newData  reflect.Value //! 新数据
}

//! 初始化
func (self *DataTest) Init(data interface{}) {
	self.newData = reflect.ValueOf(data).Elem()
	self.baseData = reflect.New(self.newData.Type()).Elem()
	self.baseData.Set(self.newData)
}

func (self *DataTest) Update(sql bool) {
	valueList := ""
	res := string(game.HF_JtoB(self.newData.Interface()))
	if !strings.Contains(res, "63") {
		panic("may be error")
	}
	//! 跳过tableKey
	for i := 1; i < self.baseData.NumField(); i++ {
		baseInt, newInt := int64(0), int64(0)
		baseStr, newStr := "", ""
		baseFloat, newFloat := float64(0.0), float64(0.0)

		//! 类型不同
		if self.baseData.Field(i).Type() != self.newData.Field(i).Type() {
			continue
		}

		switch self.baseData.Field(i).Kind() {
		case reflect.Int64:
			baseInt = self.baseData.Field(i).Int()
			newInt = self.newData.Field(i).Int()
		case reflect.Int:
			baseInt = self.baseData.Field(i).Int()
			newInt = self.newData.Field(i).Int()
		case reflect.Int8:
			baseInt = self.baseData.Field(i).Int()
			newInt = self.newData.Field(i).Int()
		case reflect.String:
			baseStr = self.baseData.Field(i).String()
			newStr = self.newData.Field(i).String()
		case reflect.Float32:
			baseFloat = self.baseData.Field(i).Float()
			newFloat = self.newData.Field(i).Float()
		case reflect.Float64:
			baseFloat = self.baseData.Field(i).Float()
			newFloat = self.newData.Field(i).Float()
		default:
			continue
		}

		rowName := strings.ToLower(self.baseData.Type().Field(i).Name)

		if baseInt != newInt {
			valueList += fmt.Sprintf("`%s`=%d,", rowName, int(newInt))
			self.baseData.Field(i).SetInt(newInt)
		} else if baseStr != newStr {
			valueList += fmt.Sprintf("`%s`='%s',", rowName, newStr)
			self.baseData.Field(i).SetString(newStr)
		} else if baseFloat != newFloat {
			valueList += fmt.Sprintf("`%s`=%f,", rowName, newFloat)
			self.baseData.Field(i).SetFloat(newFloat)
		}
	}

	if valueList != "" {
		valueKey := self.baseData.Field(0).Int()
		tableKey := strings.ToLower(self.baseData.Type().Field(0).Name)

		updateQuery := fmt.Sprintf("update `%s` set %s where `%s`=%d limit 1", "table", valueList, tableKey, valueKey)
		//! 去掉多余逗号
		updateQuery = strings.Replace(updateQuery, ", where", " where", 1)
	}
	return
}

func TestRecharge(t *testing.T) {
	pInfo := UserRecharge{}
	pInfo.VipBox = 31
	pInfo.Init(&pInfo)
	pInfo.VipBox = 63
	for i := 0; i < 200000; i++ {
		pInfo.Update(true)
	}

}
