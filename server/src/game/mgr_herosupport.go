package game

import (
	//"encoding/json"
	"fmt"
	//"log"
	//"github.com/garyburd/redigo/redis"
	"sync"
)

type San_Num struct {
	Num int
}

type HeroSupportMgr struct {
	CampNum [3]int
	JJNum   int //开服基金
	Locker  *sync.RWMutex
}

var herosupportmgrsingleton *HeroSupportMgr = nil

//! public
func GetHeroSupportMgr() *HeroSupportMgr {
	if herosupportmgrsingleton == nil {
		herosupportmgrsingleton = new(HeroSupportMgr)
		herosupportmgrsingleton.Locker = new(sync.RWMutex)
	}

	return herosupportmgrsingleton
}

func (self *HeroSupportMgr) GetData() {
	for i := CAMP_SHU; i <= CAMP_WU; i++ {
		sql := fmt.Sprintf("select count(*) as num from san_userbase where camp = %d", i)
		var num San_Num
		GetServer().DBUser.GetOneData(sql, &num, "", 0)
		self.CampNum[i-1] = num.Num
	}

	//LogDebug("各国家人数:", self.CampNum)

	//! 获取基金人数
	sql1 := "select count(*) as num from san_userrecharge where fundtype != 0"
	var num San_Num
	GetServer().DBUser.GetOneData(sql1, &num, "", 0)
	self.JJNum = num.Num
	self.JJNum += 500
	//if self.JJNum < 100 {
	//	self.JJNum = 100
	//}
}

func (self *HeroSupportMgr) Save() {
}

func (self *HeroSupportMgr) GetCamp() int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	camp := CAMP_SHU
	if self.CampNum[CAMP_WEI-1] < self.CampNum[camp-1] {
		camp = CAMP_WEI
	}
	if self.CampNum[CAMP_WU-1] < self.CampNum[camp-1] {
		camp = CAMP_WU
	}

	LogDebug("随机国家:", camp)

	return camp
}

func (self *HeroSupportMgr) AddCamp(camp int) {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.CampNum[camp-1]++

	//LogDebug("各国家人数:", self.CampNum)
}

func (self *HeroSupportMgr) AddJJ() {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	self.JJNum++
	//LogDebug("开服基金购买人数增加1")
}

func (self *HeroSupportMgr) GetJJ() int {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	return self.JJNum
}
