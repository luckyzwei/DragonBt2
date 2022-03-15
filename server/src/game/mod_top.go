package game

import (
	"encoding/json"
)

const (
	TOP_RANK_FIGHT               = 1  // 战力排行  排行榜 走全服及阵营
	TOP_RANK_TOWER               = 2  // 爬塔排行总榜
	TOP_RANK_UNION               = 3  // 军团战力排行 排行榜 走全服
	TOP_RANK_HERO_TALENT         = 4  // 英雄天赋排行 活动
	TOP_RANK_PASS_STAR           = 5  // 关卡星级排行 排行榜 走全服及阵营
	TOP_RANK_BEAUTY              = 6  // 圣物排行  走自己的模块
	TOP_RANK_CITY_NUM            = 8  // 城池排行 活动
	TOP_RANK_HERO_STAR           = 9  // 英雄星级排行 活动
	TOP_RANK_ARENA_NORMAL        = 10 // 竞技场排行 排行榜
	TOP_RANK_EQUIP_GEM           = 11 // 装备宝石等级排行 活动
	TOP_RANK_HORSE_FIGHT         = 12 // 魔宠战力排行 活动
	TOP_RANK_TIGER_FIGHT         = 13 // 纹章排行
	TOP_RANK_HERO_TALENT_CAMP1   = 14 // 英雄天赋排行 活动
	TOP_RANK_HERO_TALENT_CAMP2   = 15 // 英雄天赋排行 活动
	TOP_RANK_HERO_TALENT_CAMP3   = 16 // 英雄天赋排行 活动
	TOP_RANK_HERO_TALENT_CAMP4   = 17 // 英雄天赋排行 活动
	TOP_RANK_TOWER1              = 18 // 爬塔排行1
	TOP_RANK_TOWER2              = 19 // 爬塔排行2
	TOP_RANK_TOWER3              = 20 // 爬塔排行3
	TOP_RANK_TOWER4              = 21 // 爬塔排行4
	TOP_RANK_ARENA_SPECIAL_RANK  = 22 // 高级竞技场排名排行
	TOP_RANK_ARENA_SPECIAL_POINT = 23 // 高级竞技场积分排行
	TOP_RANK_ARENA_HIGHEST       = 24 // 巅峰竞技场排行
)

type ModTop struct {
	player *Player
}

func (self *ModTop) OnGetData(player *Player) {
	self.player = player
}

func (self *ModTop) OnGetOtherData() {
}

func (self *ModTop) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "gettop":
		var c2s_msg C2S_GetTop
		json.Unmarshal(body, &c2s_msg)
		self.GetTop(c2s_msg.Index, c2s_msg.Ver)
		return true
	}

	return false
}

func (self *ModTop) OnSave(sql bool) {

}

//获取排行榜
func (self *ModTop) GetTop(index int, ver int) {
	topType := index
	switch topType {
	case TOP_RANK_FIGHT:
		self.getFightTop(topType, ver)
	//case TOP_RANK_UNION:
	//	self.getUnionTop(topType, ver)
	case TOP_RANK_HERO_TALENT, TOP_RANK_HERO_TALENT_CAMP1, TOP_RANK_HERO_TALENT_CAMP2, TOP_RANK_HERO_TALENT_CAMP3, TOP_RANK_HERO_TALENT_CAMP4:
		self.getTalentTop(topType, ver)
	case TOP_RANK_PASS_STAR:
		self.getPassTop(topType, ver)
	case TOP_RANK_HERO_STAR:
		self.getHeroStarTop(topType, ver)
	case TOP_RANK_EQUIP_GEM:
		self.getEquipGemTop(topType, ver)
	case TOP_RANK_HORSE_FIGHT:
		self.getHorseFightTop(topType, ver)
	case TOP_RANK_TOWER, TOP_RANK_TOWER1, TOP_RANK_TOWER2, TOP_RANK_TOWER3, TOP_RANK_TOWER4:
		self.getTowerTop(topType, ver)
	case TOP_RANK_ARENA_NORMAL, TOP_RANK_ARENA_SPECIAL_RANK, TOP_RANK_ARENA_SPECIAL_POINT, TOP_RANK_ARENA_HIGHEST:
		self.getArenaTop(topType, ver)
	}
}

// 战斗力显示
func (self *ModTop) getFightTop(topType int, ver int) {
	var msg S2C_Top
	msg.Cid = "gettop"
	_, curver := GetTopFightMgr().GetTopShow()
	if curver != ver || ver == 0 {
		//msg.Top = lst
		msg.Type = topType
		msg.Ver = curver
		msg.Cur = GetTopMgr().GetTopCurNum(topType, self.player.Sql_UserBase.Uid)
		msg.Old = GetTopMgr().GetTopOldNum(topType, self.player.Sql_UserBase.Uid)
	}
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("gettop", smsg)
}

//// 军团显示
//func (self *ModTop) getUnionTop(topType int, ver int) {
//	var msg S2C_Top
//	msg.Cid = "gettop"
//	_, curver := GetTopUnionMgr().GetTopUnionShow()
//	if curver != ver || ver == 0 {
//		//msg.Top = lst
//		msg.Type = topType
//		msg.Ver = curver
//		unionid := self.player.GetModule("union").(*ModUnion).Sql_UserUnionInfo.Unionid
//		if unionid > 0 {
//			msg.Cur = GetTopMgr().GetTopCurNum(topType, int64(unionid))
//			msg.Old = GetTopMgr().GetTopOldNum(topType, int64(unionid))
//		}
//	}
//	smsg, _ := json.Marshal(&msg)
//	self.player.SendMsg("gettop", smsg)
//}

// 关卡星级
func (self *ModTop) getPassTop(topType int, ver int) {
	var msg S2C_Top
	msg.Cid = "gettop"
	lst, curver := GetTopPassMgr().GetTopPassShow()
	if curver != ver || ver == 0 {
		msg.Top = lst
		msg.Type = topType
		msg.Ver = curver
		msg.Cur = GetTopPassMgr().GetTopCurNum(topType, self.player.Sql_UserBase.Uid)
		msg.Old = GetTopPassMgr().GetTopOldNum(topType, self.player.Sql_UserBase.Uid)
		msg.Num = int64(self.player.GetModule("onhook").(*ModOnHook).Sql_OnHook.OnHookStage)
	}
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("gettop", smsg)
}

// 天赋星级
func (self *ModTop) getTalentTop(topType int, ver int) {
	var msg S2C_Top
	msg.Cid = "gettop"
	lst, curver := GetTopTalentMgr().GetTopShow(topType)
	if curver != ver || ver == 0 {
		msg.Top = lst
		msg.Type = topType
		msg.Ver = curver
		msg.Cur = GetTopTalentMgr().GetTopCurNum(topType, self.player.Sql_UserBase.Uid)
		msg.Old = GetTopTalentMgr().GetTopOldNum(topType, self.player.Sql_UserBase.Uid)
		nType := GetTopTalentMgr().GetTopType(topType)
		msg.Num = int64(self.player.GetModule("hero").(*ModHero).Sql_Hero.totalStars[nType].Stars)
	}
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("gettop", smsg)
}

// 塔等级
func (self *ModTop) getTowerTop(topType int, ver int) {
	var msg S2C_Top
	msg.Cid = "gettop"
	lst, curver := GetTopTowerMgr().GetTopShow(topType)
	if curver != ver || ver == 0 {
		msg.Top = lst
		msg.Type = topType
		msg.Ver = curver
		msg.Cur = GetTopTowerMgr().GetTopCurNum(topType, self.player.Sql_UserBase.Uid)
		msg.Old = GetTopTowerMgr().GetTopOldNum(topType, self.player.Sql_UserBase.Uid)
		nType := GetTopTowerMgr().GetTopType(topType)
		msg.Num = int64(self.player.GetModule("tower").(*ModTower).data.info[nType].MaxLevel)
	}
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("gettop", smsg)
}

// 英雄星级
func (self *ModTop) getHeroStarTop(topType int, ver int) {
	var msg S2C_Top
	msg.Cid = "gettop"
	_, curver := GetTopHeroStarsMgr().GetTopHeroStarsShow()
	if curver != ver || ver == 0 {
		//msg.Top = lst
		msg.Type = topType
		msg.Ver = curver
		msg.Cur = GetTopHeroStarsMgr().GetTopCurNum(topType, self.player.Sql_UserBase.Uid)
		msg.Old = 0
		msg.Num = int64(self.player.GetModule("hero").(*ModHero).Sql_Hero.totalStars[HERO_STAR_TOTAL].Stars)
	}
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("gettop", smsg)
}

// 装备宝石等级
func (self *ModTop) getEquipGemTop(topType int, ver int) {
	var msg S2C_Top
	msg.Cid = "gettop"
	_, curver := GetTopEquipGemMgr().GetTopEquipGemShow()
	if curver != ver || ver == 0 {
		//msg.Top = lst
		msg.Type = topType
		msg.Ver = curver
		msg.Cur = GetTopEquipGemMgr().GetTopCurNum(topType, self.player.Sql_UserBase.Uid)
		msg.Old = 0
		msg.Num = int64(self.player.GetModule("equip").(*ModEquip).Data.TotalGemLevel)
	}
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("gettop", smsg)
}

// 魔宠战力
func (self *ModTop) getHorseFightTop(topType int, ver int) {
	var msg S2C_Top
	msg.Cid = "gettop"
	_, curver := GetTopHorseFightMgr().GetTopHorseFightShow()
	if curver != ver || ver == 0 {
		//msg.Top = lst
		msg.Type = topType
		msg.Ver = curver
		msg.Cur = GetTopHorseFightMgr().GetTopCurNum(topType, self.player.Sql_UserBase.Uid)
	}
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("gettop", smsg)
}

// 排行版
func (self *ModTop) getArenaTop(topType int, ver int) {
	var msg S2C_Top
	msg.Cid = "gettop"
	lst, curver := GetTopArenaMgr().GetTopShow(topType)
	if curver != ver || ver == 0 {
		msg.Top = lst
		msg.Type = topType
		msg.Ver = curver
		msg.Cur = GetTopArenaMgr().GetTopCurNum(topType, self.player.Sql_UserBase.Uid)
		msg.Old = GetTopArenaMgr().GetTopOldNum(topType, self.player.Sql_UserBase.Uid)
		//nType := GetTopArenaMgr().GetTopType(topType)
		if TOP_RANK_ARENA_NORMAL == topType {
			data := GetArenaMgr().GetPlayerArenaData(self.player.GetUid())
			if data == nil {
				return
			}
			msg.Num = data.Point
		} else if TOP_RANK_ARENA_SPECIAL_RANK == topType {
			data := GetArenaSpecialMgr().GetPlayerData(self.player.GetUid())
			if data == nil {
				return
			}

			config := GetCsvMgr().GetArenaSpecialClassConfig(data.Class, data.Dan)
			msg.Num = int64(config.Ranking)
		} else if TOP_RANK_ARENA_SPECIAL_POINT == topType {
			data := GetArenaSpecialMgr().GetPlayerData(self.player.GetUid())
			if data == nil {
				return
			}

			msg.Num = int64(data.Point)
		}
	}
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("gettop", smsg)
}
