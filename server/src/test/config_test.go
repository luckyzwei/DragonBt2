package test

import (
	"fmt"
	. "game"
	"github.com/kr/pretty"
	"github.com/sanity-io/litter"
	"testing"
)
func TestConfig(t *testing.T)  {
	GetCsvMgr().InitData()
	//tConfig := GetCsvMgr().GetTariffConfig(8, 5)
	//if tConfig == nil {
	//	t.Error("config must be wrong!")
	//}
	//fmt.Printf("%# v\n", pretty.Formatter(tConfig))
	//p := GetCsvMgr().GetTechAttrConfig(18, 1)
	p := GetArmyMgr().GetInstanceConfig(100006)
	fmt.Printf("%# v\n", pretty.Formatter(p))
}

func TestEnum(t *testing.T)  {
	//t.Log("hello world!")
	//var a  = [4]int{1,2,3,4}
	//var b = a
	//b[2] = 10
	//t.Log(b)
}

func TestDeclare(t *testing.T)   {
	GetCsvMgr().InitData()
	//declareTime,declareTim2 := GetUnionFightMgr().GetCallTime()
	//t.Log(time.Unix(declareTime, 0).Format(DATEFORMAT))
	//t.Log(time.Unix(declareTim2, 0).Format(DATEFORMAT))
	//readyTime := GetUnionFightMgr().GetDelay()
	//t.Log("readyTime:", readyTime)
	items := GetLootMgr().LootItems([]int{110001,100001})
	litter.Dump(items)
}