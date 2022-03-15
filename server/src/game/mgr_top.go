package game

import (
	"time"
)

type TopMgr struct{}

var topmgr *TopMgr = nil

//! public
func GetTopMgr() *TopMgr {
	if topmgr == nil {
		topmgr = new(TopMgr)
	}
	return topmgr
}

//! run
func (self *TopMgr) Run() {
	self.GetData()
	GetTopFightMgr().InitTopActual()
	GetTopUnionMgr().InitTopActual()
	GetTopLevelMgr().InitTopActual()
	GetTopCityMgr().getTopRedis()
	//新写的排行榜就开服调用一次
	GetTopInterstellarMgr().GetData()
	ticker := time.NewTicker(time.Minute * 1)
	for {
		<-ticker.C
		GetTopCityMgr().setTopRedis()
		minute := TimeServer().Minute()
		if minute%5 == 0 {
			self.GetData()
			self.updateRank()
		}
	}

	ticker.Stop()
}

// 更新排行榜信息
func (self *TopMgr) updateRankCache() {
	self.GetData()
	self.updateRank()
}

//! 每1个小时重新获取数据
func (self *TopMgr) GetData() {
	unionName := GetUnionMgr().GetUnionNameMap()
	GetTopFightMgr().GetData(unionName)
	GetTopUnionMgr().GetData()
	GetTopLevelMgr().GetData(unionName)
	GetTopPassMgr().GetData(unionName)
	GetTopTowerMgr().GetData(unionName)
	GetTopTalentMgr().GetData(unionName)
	GetTopPvpMgr().GetData(unionName)
	GetTopHeroStarsMgr().GetData(unionName)
	GetTopEquipGemMgr().GetData(unionName)
	GetTopHorseFightMgr().GetData(unionName)
	GetTopTigerFightMgr().GetData(unionName)
	GetTopArenaMgr().GetData(unionName)
	//GetTopInterstellarMgr().GetData() 开服调用一次  以后不用调用了
}

func (self *TopMgr) GetTopCurNum(topType int, id int64) int {
	switch topType {
	case TOP_RANK_FIGHT: // 战力
		return GetTopFightMgr().GetTopCurNum(topType, id)
	case TOP_RANK_UNION: // 军团战力
		return GetTopUnionMgr().GetTopCurNum(topType, id)
	case TOP_RANK_HERO_TALENT, TOP_RANK_HERO_TALENT_CAMP1, TOP_RANK_HERO_TALENT_CAMP2, TOP_RANK_HERO_TALENT_CAMP3, TOP_RANK_HERO_TALENT_CAMP4: // 天赋总星级
		return GetTopTalentMgr().GetTopCurNum(topType, id)
	case TOP_RANK_PASS_STAR: // 关卡
		return GetTopPassMgr().GetTopCurNum(topType, id)
	case TOP_RANK_TOWER, TOP_RANK_TOWER1, TOP_RANK_TOWER2, TOP_RANK_TOWER3, TOP_RANK_TOWER4: // 爬塔
		return GetTopTowerMgr().GetTopCurNum(topType, id)
	}

	return 0
}

func (self *TopMgr) GetTopOldNum(topType int, id int64) int {
	switch topType {
	case TOP_RANK_FIGHT: // 战力
		return GetTopFightMgr().GetTopOldNum(topType, id)
	case TOP_RANK_UNION: // 军团战力
		return GetTopUnionMgr().GetTopOldNum(topType, id)
	case TOP_RANK_HERO_TALENT, TOP_RANK_HERO_TALENT_CAMP1, TOP_RANK_HERO_TALENT_CAMP2, TOP_RANK_HERO_TALENT_CAMP3, TOP_RANK_HERO_TALENT_CAMP4: // 天赋总星级
		return GetTopTalentMgr().GetTopOldNum(topType, id)
	case TOP_RANK_PASS_STAR: // 关卡
		return GetTopPassMgr().GetTopOldNum(topType, id)
	case TOP_RANK_TOWER, TOP_RANK_TOWER1, TOP_RANK_TOWER2, TOP_RANK_TOWER3, TOP_RANK_TOWER4: // 爬塔
		return GetTopTowerMgr().GetTopOldNum(topType, id)
	}

	return 0
}

//! 更名更新
func (self *TopMgr) SyncPlayerName(player *Player) {
	//GetTopFightMgr().SyncPlayerName(player)
	GetTopPassMgr().SyncPlayerName(player)
	GetTopTowerMgr().Rename(player)
	GetTopTalentMgr().Rename(player)
	GetTopArenaMgr().Rename(player)
}

func (self *TopMgr) updateRank() {
	unionName := GetUnionMgr().GetUnionNameMap()
	//GetTopEquipGemMgr().GetData(unionName)
	//GetTopTigerFightMgr().GetData(unionName)
	//GetTopHorseFightMgr().GetData(unionName)
	//GetTopPvpMgr().GetData(unionName)
	//GetTopHeroStarsMgr().GetData(unionName)
	GetTopPassMgr().GetData(unionName)
	GetTopTowerMgr().GetData(unionName)
	GetTopTalentMgr().GetData(unionName)
	GetTopArenaMgr().GetData(unionName)
}
