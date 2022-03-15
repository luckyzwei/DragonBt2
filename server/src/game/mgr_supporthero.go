package game

import (
	"encoding/json"
)

const (
	HERO_SUPPORT_TYPE_ENTANGLEMENT = 1 // 缘分
	HERO_SUPPORT_TYPE_REWARD       = 2 // 悬赏
)

const (
	SUPPORT_HERO_END_TIME = 7 // 过期时间
	SUPPORT_HERO_CD_TIME  = 1 // cd时间
	SUPPORT_HERO_MAX      = 6 // 上阵支援英雄最多个数
	SUPPORT_HERO_USE_MAX  = 3 // 使用支援英雄最多个数
	SUPPORT_HERO_AUTO_MAX = 5 // 自动上阵个数
)

////!支援英雄
//type San_SupportHero struct {
//	Uid         int64  // 角色Id
//	SupportHero string // 支援英雄
//
//	supportHero []*SupportHero // 支援英雄
//
//	DataUpdate
//}

//// 检查过期英雄
//func (self *San_SupportHero) CheckEndTime() {
//	// 当前时间
//	timeNow := TimeServer().Unix()
//	nLen := len(self.supportHero)
//	for i := nLen - 1; i >= 0; i-- {
//		// 已过期
//		if self.supportHero[i].EndTime != 0 && timeNow >= self.supportHero[i].EndTime {
//			//self.supportHero = append(self.supportHero[:i], self.supportHero[i+1:]...)
//			self.supportHero[i].Type = 0
//			self.supportHero[i].UserUid = 0
//			self.supportHero[i].UserName = ""
//			self.supportHero[i].EndTime = 0
//		}
//	}
//}
//
//func (self *San_SupportHero) Decode() { //! 将数据库数据写入data
//	json.Unmarshal([]byte(self.SupportHero), &self.supportHero)
//}
//
//func (self *San_SupportHero) Encode() { //! 将data数据写入数据库
//	self.SupportHero = HF_JtoA(&self.supportHero)
//}

//type SupportHero struct {
//	HeroKey      int    `json:"herokey"`      // 英雄key值
//	HeroID       int    `json:"heroid"`       // 英雄id
//	HeroIcon     int    `json:"heroicon"`     // 英雄头像
//	HeroPortrait int    `json:"heroportrait"` // 英雄头像框
//	HeroStar     int    `json:"herostar"`     // 英雄星级
//	HeroCamp     int    `json:"herocamp"`     // 英雄阵营
//	HeroLevel    int    `json:"herolevel"`    // 英雄等级
//	MasterUid    int64  `json:"masteruid"`    // 主人uid
//	MasterName   string `json:"mastername"`   // 主人名字
//	Type         int    `json:"type"`         // 使用类型
//	UserUid      int64  `json:"useruid"`      // 使用者uid
//	UserName     string `json:"username"`     // 使用者名字
//	EndTime      int64  `json:"endtime"`      // 结束时间
//	CDTime       int64  `json:"cdtime"`       // cd结束时间
//}

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

type SupportHeroMgr struct {
	//Locker          *sync.RWMutex              //! 数据锁
	//Sql_SupportHero map[int64]*San_SupportHero //! 玩家支援英雄
}

var supportheromgrsingleton *SupportHeroMgr = nil

func GetSupportHeroMgr() *SupportHeroMgr {
	if supportheromgrsingleton == nil {
		supportheromgrsingleton = new(SupportHeroMgr)
		//supportheromgrsingleton.Locker = new(sync.RWMutex)
		//supportheromgrsingleton.Sql_SupportHero = make(map[int64]*San_SupportHero)
	}

	return supportheromgrsingleton
}
func (self *SupportHeroMgr) GetData() {
	//var supportHero San_SupportHero
	//sql := fmt.Sprintf("select * from `san_supporthero`")
	//res := GetServer().DBUser.GetAllData(sql, &supportHero)
	//for i := 0; i < len(res); i++ {
	//	data, ok1 := res[i].(*San_SupportHero)
	//	if !ok1 {
	//		continue
	//	}
	//
	//	data.Init("san_supporthero", data, false)
	//	data.Decode()
	//	_, ok2 := self.Sql_SupportHero[data.Uid]
	//	if !ok2 {
	//		self.Sql_SupportHero[data.Uid] = data
	//	}
	//}
}

// 存储数据库
func (self *SupportHeroMgr) Save() {
	//self.Locker.Lock()
	//defer self.Locker.Unlock()
	//
	//for _, v := range self.Sql_SupportHero {
	//	v.Encode()
	//	v.Update(true)
	//}
}

//// 添加支援英雄
//func (self *SupportHeroMgr) AddHero1(heros []int, herokey []int, player *Player) bool {
//	self.Locker.Lock()
//	defer self.Locker.Unlock()
//
//	// 次数错误
//	if len(herokey) != len(heros) {
//		return false
//	}
//
//	// 获得添加玩家数据
//	uid := player.GetUid()
//	data, ok := self.Sql_SupportHero[uid]
//	if !ok {
//		self.Sql_SupportHero[uid] = &San_SupportHero{}
//		data, ok = self.Sql_SupportHero[uid]
//		if ok {
//			data.Uid = uid
//			data.supportHero = make([]*SupportHero, 0)
//		} else {
//			return false
//		}
//	}
//
//	// 计算结束时间
//	timeNow := TimeServer().Unix()
//	endtime := timeNow + SUPPORT_HERO_END_TIME*DAY_SECS
//	cdtime := timeNow + SUPPORT_HERO_CD_TIME*HOUR_SECS
//
//	// 设置英雄
//	for i, v := range heros {
//		config := GetCsvMgr().GetHeroConfig(v)
//		if config == nil {
//			continue
//		}
//		find := false
//		for _, t := range data.supportHero {
//			if t.HeroKey == herokey[i] {
//				find = true
//				break
//			}
//		}
//		if find{
//			continue
//		}
//
//		temp := SupportHero{
//			HeroKey:      herokey[i],
//			HeroID:       config.HeroId,    // 英雄id
//			HeroIcon:     0,                // 英雄头像
//			HeroPortrait: 0,                // 英雄头像框
//			HeroStar:     config.HeroStar,  // 英雄星级
//			HeroCamp:     config.HeroCamp,  // 英雄阵营
//			HeroLevel:    0,                // 英雄等级
//			MasterUid:    player.GetUid(),  // 主人名字
//			MasterName:   player.GetName(), // 主人uid
//			Type:         0,                // 使用类型
//			UserUid:      0,                // 使用者uid
//			UserName:     "",               // 使用者名字
//			EndTime:      endtime,          // 结束时间
//			CDTime:       cdtime,           // cd时间
//		}
//
//		data.supportHero = append(data.supportHero, &temp)
//	}
//
//	return true
//}
//
//// 添加支援英雄
//func (self *SupportHeroMgr) AddSupportHero(index int, hero *Hero, player *Player) bool {
//	self.Locker.Lock()
//	defer self.Locker.Unlock()
//
//	// 获得玩家数据
//	uid := player.GetUid()
//	data, ok := self.Sql_SupportHero[uid]
//	if !ok {
//		return false
//	}
//
//	// 判断index
//	if index <= 0 || index > SUPPORT_HERO_MAX {
//		return false
//	}
//
//	// 判断设置个数
//	if len(data.supportHero) > SUPPORT_HERO_MAX {
//		return false
//	}
//
//	// 计算结束时间
//	timeNow := TimeServer().Unix()
//	cdtime := timeNow + SUPPORT_HERO_CD_TIME*HOUR_SECS
//	config := GetCsvMgr().GetHeroMapConfig(hero.getHeroId(), hero.GetStar())
//	if config == nil {
//		return false
//	}
//
//	// 判断是否已经设置
//	for _, v := range data.supportHero {
//		if v.HeroKey == hero.HeroKeyId {
//			return false
//		}
//
//		if v.Index == index {
//			return false
//		}
//	}
//
//	// 设置英雄
//	temp := SupportHero{
//		Index:      index,            // index
//		HeroKey:    hero.HeroKeyId,   // herokey
//		HeroID:     config.HeroId,    // 英雄id
//		HeroStar:   hero.GetStar(),   // 英雄星级
//		HeroLv:     hero.HeroLv,      // 英雄等级
//		HeroSkin:   hero.Skin,        // 皮肤
//		MasterUid:  player.GetUid(),  // 主人名字
//		MasterName: player.GetName(), // 主人uid
//		Type:       0,
//		UserUid:    0,
//		UserName:   "",
//		EndTime:    0,
//		CDTime:     cdtime, // cd时间
//	}
//	data.supportHero = append(data.supportHero, &temp)
//	return true
//}
//
//// 移除支援英雄
//func (self *SupportHeroMgr) RemoveSupportHero(herokey int, uid int64) bool {
//	self.Locker.Lock()
//	defer self.Locker.Unlock()
//
//	// 获得添加玩家数据
//	data, ok := self.Sql_SupportHero[uid]
//	if !ok {
//		return false
//	}
//
//	// 判断是否设置了
//	find := false
//	array := -1
//	for i, t := range data.supportHero {
//		if t.HeroKey == herokey {
//			find = true
//			array = i
//			break
//		}
//	}
//	if !find {
//		return false
//	}
//
//	// 是否在cd中
//	//if data.supportHero[array].CDTime > TimeServer().Unix() {
//	//	return false
//	//}
//
//	// 删除
//	data.supportHero = append(data.supportHero[:array], data.supportHero[array+1:]...)
//	return true
//}
//
//// 使用英雄
//func (self *SupportHeroMgr) UseHero(uid int64, nHeroKey int, player *Player, nType int, endtime int64) bool {
//	self.Locker.Lock()
//	defer self.Locker.Unlock()
//
//	// 不能使用自己的支援英雄
//	if uid == player.GetUid() {
//		return false
//	}
//
//	// 获得数据
//	data, ok := self.Sql_SupportHero[uid]
//	if !ok {
//		return false
//	}
//
//	for _, v := range data.supportHero {
//		// 英雄id判断
//		if v.HeroKey != nHeroKey {
//			continue
//		}
//
//		//// 已经被使用了
//		//if v.UserUid != 0 {
//		//	continue
//		//}
//
//		// 设置使用
//		v.UserUid = player.GetUid()
//		v.UserName = player.GetName()
//		v.Type = nType
//		v.EndTime = endtime
//		return true
//	}
//
//	return false
//}
//
//// 取消使用英雄
//func (self *SupportHeroMgr) CancelUseHero(uid int64, nHeroKey int, player *Player) bool {
//	self.Locker.Lock()
//	defer self.Locker.Unlock()
//
//	// 不能取消自己的支援英雄
//	if uid == player.GetUid() {
//		return false
//	}
//
//	// 获得数据
//	data, ok := self.Sql_SupportHero[uid]
//	if !ok {
//		return false
//	}
//
//	for _, v := range data.supportHero {
//		// 英雄id判断
//		if v.HeroKey != nHeroKey {
//			continue
//		}
//
//		// 没有被使用
//		if v.UserUid == 0 {
//			continue
//		}
//
//		// 不是自己使用的
//		if v.UserUid != player.GetUid() {
//			continue
//		}
//
//		v.UserUid = 0
//		v.UserName = ""
//		v.Type = 0
//		v.EndTime = 0
//		return true
//	}
//	return false
//}

// 获得玩家数据
func (self *SupportHeroMgr) GetPlayerData(uid int64, create bool) []*SupportHero {
	//self.Locker.Lock()
	//defer self.Locker.Unlock()
	//
	//// 获得配置
	//data, ok := self.Sql_SupportHero[uid]
	//if ok {
	//	return data
	//}
	//
	//// 是否创建
	//if create {
	//	self.Sql_SupportHero[uid] = &San_SupportHero{}
	//	data, ok = self.Sql_SupportHero[uid]
	//	if ok {
	//		data.Uid = uid
	//		data.supportHero = make([]*SupportHero, 0)
	//		data.Encode()
	//		InsertTable("san_supporthero", data, 0, false)
	//		data.Init("san_supporthero", data, false)
	//
	//		return data
	//	} else {
	//		return nil
	//	}
	//}

	var mastermsg S2M_SupportHeroGetPlayerData
	mastermsg.Uid = uid

	ret := GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_GET_DATA, &mastermsg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return nil
	}

	var backmsg M2S_SupportHeroGetPlayerData
	json.Unmarshal([]byte(ret.Data), &backmsg)
	return backmsg.Data
}

// 清理玩家数据
func (self *SupportHeroMgr) CleanPlayerData(uid int64) {
	var mastermsg S2M_SupportHeroCleanPlayerData
	mastermsg.Uid = uid

	ret := GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_CLEAN_DATA, &mastermsg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return
	}
}

// 获得可使用的英雄
func (self *SupportHeroMgr) GetCanUseHero1(uids map[int64]int64, heroID int) []*SupportHero {
	var mastermsg S2M_SupportHeroGetCanUseHero
	mastermsg.Uids = uids
	mastermsg.HeroID = heroID

	ret := GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_GET_CAN_USE_HERO, &mastermsg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return nil
	}

	var backmsg M2S_SupportHeroGetCanUseHero
	json.Unmarshal([]byte(ret.Data), &backmsg)
	return backmsg.Data
}

//
//// 获得可使用的英雄
//func (self *SupportHeroMgr) GetCanUseHero2(uid int64) []*SupportHero {
//	self.Locker.Lock()
//	defer self.Locker.Unlock()
//
//	ret := []*SupportHero{}
//	data, ok := self.Sql_SupportHero[uid]
//	if !ok {
//		return ret
//	}
//	// 检测过期
//	data.CheckEndTime()
//	for _, v := range data.supportHero {
//		ret = append(ret, v)
//	}
//	return ret
//}

// 获得我的的英雄
func (self *SupportHeroMgr) GetMyHero(uid int64) []*SupportHero {
	//self.Locker.Lock()
	//defer self.Locker.Unlock()
	//
	//ret := []*SupportHero{}
	//data, ok := self.Sql_SupportHero[uid]
	//if !ok {
	//	return ret
	//}
	//// 检测过期
	//data.CheckEndTime()
	//for _, v := range data.supportHero {
	//	ret = append(ret, v)
	//}
	uids := make(map[int64]int64, 0)
	uids[uid] = uid
	return self.GetCanUseHero1(uids, 0)
}

// 更新英雄
func (self *SupportHeroMgr) UpdateHero(uid int64, hero *Hero) bool {
	var mastermsg S2M_SupportHeroUpdate
	mastermsg.Uid = uid
	mastermsg.HeroKeyId = hero.HeroKeyId
	mastermsg.HeroLv = hero.HeroLv
	mastermsg.HeroStar = hero.StarItem.UpStar
	mastermsg.Skin = hero.Skin

	ret := GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_UPDATE_HERO, &mastermsg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	return true
}

// 修改名字
func (self *SupportHeroMgr) Rename(player *Player) bool {
	//self.Locker.Lock()
	//defer self.Locker.Unlock()
	//
	//// 获得玩家数据
	//uid := player.GetUid()
	//data, ok := self.Sql_SupportHero[uid]
	//if !ok {
	//	return false
	//}
	//
	//for _, v := range data.supportHero {
	//	v.MasterName = player.GetName()
	//}
	//return true

	var mastermsg S2M_SupportHeroRename
	mastermsg.Uid = player.GetUid()
	mastermsg.Name = player.GetUname()

	ret := GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_RENAME, &mastermsg)
	if ret == nil || ret.RetCode != UNION_SUCCESS {
		return false
	}

	return true
}
