/*
@Time : 2020/4/22 10:10
@Author : 96121
@File : data_friend
@Software: GoLand
*/
package player

import (
	"master/utils"
	"time"
)

///////////////////////////////////////////////////////

const (
	SUPPORT_HERO_END_TIME = 7 // 过期时间
	SUPPORT_HERO_CD_TIME  = 1 // cd时间
	SUPPORT_HERO_MAX      = 6 // 上阵支援英雄最多个数
	SUPPORT_HERO_USE_MAX  = 3 // 使用支援英雄最多个数
)

type SupportHero struct {
	Index      int    `json:"index"`      // index
	HeroKey    int    `json:"herokey"`    // 英雄key值
	HeroID     int    `json:"heroid"`     // 英雄id
	HeroStar   int    `json:"herostar"`   // 英雄星级
	HeroLv     int    `json:"herolv"`     // 英雄等级
	HeroSkin   int    `json:"skin"`       // 皮肤
	MasterUid  int64  `json:"masteruid"`  // 主人uid
	MasterName string `json:"mastername"` // 主人名字
	Type       int    `json:"type"`       // 使用类型
	UserUid    int64  `json:"useruid"`    // 使用者uid
	UserName   string `json:"username"`   // 使用者名字
	EndTime    int64  `json:"endtime"`    // 结束时间
	CDTime     int64  `json:"cdtime"`     // cd结束时间
}

type MySupportHero struct {
	Index    int    `json:"index"`    // index
	HeroKey  int    `json:"herokey"`  // 英雄key值
	Type     int    `json:"type"`     // 使用类型
	UserUid  int64  `json:"useruid"`  // 使用者uid
	UserName string `json:"username"` // 使用者名字
	CDTime   int64  `json:"cdtime"`   // cd时间
}

type MsgSupportHero struct {
	HeroKeyId int `json:"herokeyid"` // 英雄key值
	HeroID    int `json:"heroid"`    // 英雄id
	HeroStar  int `json:"herostar"`  // 英雄星级
	HeroLv    int `json:"herolv"`    // 英雄等级
	Skin      int `json:"skin"`      // 皮肤
}

// 检查过期英雄
func (self *ModFriend) CheckEndTime() {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	// 当前时间
	timeNow := time.Now().Unix()
	nLen := len(self.Data.supportHero)
	for i := nLen - 1; i >= 0; i-- {
		// 已过期
		if self.Data.supportHero[i].EndTime != 0 && timeNow >= self.Data.supportHero[i].EndTime {
			//self.supportHero = append(self.supportHero[:i], self.supportHero[i+1:]...)
			self.Data.supportHero[i].Type = 0
			self.Data.supportHero[i].UserUid = 0
			self.Data.supportHero[i].UserName = ""
			self.Data.supportHero[i].EndTime = 0
		}
	}
}

// 添加支援英雄
func (self *ModFriend) AddSupportHero(index int, hero *MsgSupportHero) bool {
	// 判断index
	if index <= 0 || index > SUPPORT_HERO_MAX {
		return false
	}

	// 判断设置个数
	if len(self.Data.supportHero) > SUPPORT_HERO_MAX {
		return false
	}

	// 计算结束时间
	timeNow := time.Now().Unix()
	cdtime := timeNow + SUPPORT_HERO_CD_TIME*utils.HOUR_SECS

	// 判断是否已经设置
	for _, v := range self.Data.supportHero {
		if v.HeroKey == hero.HeroKeyId {
			return false
		}

		if v.Index == index {
			return false
		}
	}

	// 设置英雄
	temp := SupportHero{
		Index:      index,                  // index
		HeroKey:    hero.HeroKeyId,         // herokey
		HeroID:     hero.HeroID,            // 英雄id
		HeroStar:   hero.HeroStar,          // 英雄星级
		HeroLv:     hero.HeroLv,            // 英雄等级
		HeroSkin:   hero.Skin,              // 皮肤
		MasterUid:  self.player.GetUid(),   // 主人名字
		MasterName: self.player.GetUname(), // 主人uid
		Type:       0,
		UserUid:    0,
		UserName:   "",
		EndTime:    0,
		CDTime:     cdtime, // cd时间
	}
	self.Locker.Lock()
	self.Data.supportHero = append(self.Data.supportHero, &temp)
	self.Locker.Unlock()
	return true
}

// 移除支援英雄
func (self *ModFriend) RemoveSupportHero(herokey int) bool {
	// 判断是否设置了
	find := false
	array := -1
	for i, t := range self.Data.supportHero {
		if t.HeroKey == herokey {
			find = true
			array = i
			break
		}
	}
	if !find {
		return false
	}

	// 是否在cd中
	//if data.supportHero[array].CDTime > time.Now().Unix() {
	//	return false
	//}

	self.Locker.Lock()
	// 删除
	self.Data.supportHero = append(self.Data.supportHero[:array], self.Data.supportHero[array+1:]...)
	self.Locker.Unlock()
	return true
}

// 使用英雄
func (self *ModFriend) UseHero(nHeroKey int, useruid int64, username string, nType int, endtime int64) bool {
	// 不能使用自己的支援英雄
	if self.player.GetUid() == useruid {
		return false
	}

	self.Locker.Lock()
	defer self.Locker.Unlock()
	for _, v := range self.Data.supportHero {
		// 英雄id判断
		if v.HeroKey != nHeroKey {
			continue
		}

		//// 已经被使用了
		//if v.UserUid != 0 {
		//	continue
		//}

		// 设置使用
		v.UserUid = useruid
		v.UserName = username
		v.Type = nType
		v.EndTime = endtime
		return true
	}

	return false
}

// 取消使用英雄
func (self *ModFriend) CancelUseHero(nHeroKey int, useruid int64) bool {
	// 不能取消自己的支援英雄
	if self.player.GetUid() == useruid {
		return false
	}

	self.Locker.RLock()
	defer self.Locker.RUnlock()
	for _, v := range self.Data.supportHero {
		// 英雄id判断
		if v.HeroKey != nHeroKey {
			continue
		}

		// 没有被使用
		if v.UserUid == 0 {
			continue
		}

		// 不是自己使用的
		if v.UserUid != useruid {
			continue
		}

		v.UserUid = 0
		v.UserName = ""
		v.Type = 0
		v.EndTime = 0
		return true
	}
	return false
}

// 获得玩家数据
func (self *ModFriend) CleanPlayerData(uid int64) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	self.Data.supportHero = []*SupportHero{}
}

//// 获得我的的英雄
//func (self *ModFriend) GetMyHero() []*SupportHero {
//	ret := []*SupportHero{}
//	// 检测过期
//	self.CheckEndTime()
//	self.Locker.RLock()
//	defer self.Locker.RUnlock()
//	for _, v := range self.Data.supportHero {
//		ret = append(ret, v)
//	}
//	return ret
//}

// 更新英雄
func (self *ModFriend) UpdateHero(hero *MsgSupportHero) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	for _, v := range self.Data.supportHero {
		if v.HeroKey == hero.HeroKeyId {
			v.HeroStar = hero.HeroStar
			v.HeroLv = hero.HeroLv
			v.HeroSkin = hero.Skin
			return true
		}
	}

	return false
}

// 修改名字
func (self *ModFriend) Rename(name string) bool {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	for _, v := range self.Data.supportHero {
		v.MasterName = name
	}
	return true
}

