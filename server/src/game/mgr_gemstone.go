package game

type GemStoneMgr struct {
}

var gemStoneMgrSingleton *GemStoneMgr = nil

func GetGemStoneMgr() *GemStoneMgr {

	if gemStoneMgrSingleton == nil {
		gemStoneMgrSingleton = new(GemStoneMgr)
	}

	return gemStoneMgrSingleton
}

func (self *GemStoneMgr) GetLevelConfig(id int) *GemstoneLevelConfig {

	ret, ok := GetCsvMgr().GemstoneLevelConfigMap[id]
	if ok == false {
		return nil
	}
	return ret
}

func (self *GemStoneMgr) GetChapterConfig(id int) *GemstoneChapterConfig {
	ret, ok := GetCsvMgr().GemstoneChapterConfigMap[id]
	if ok == false {
		return nil
	}
	return ret
}

func (self *GemStoneMgr) GetChapterMaxId(cid int) int {
	ret := 0
	for i := 0; i < len(GetCsvMgr().GemstoneLevelConfig); i++ {
		if GetCsvMgr().GemstoneLevelConfig[i].ChaptrGroup == cid {
			if GetCsvMgr().GemstoneLevelConfig[i].LevelIndex > ret {
				ret = GetCsvMgr().GemstoneLevelConfig[i].LevelIndex
			}
		}
	}

	return ret
}

func (self *GemStoneMgr) GetChapterMaxLevel(cid int) int {
	ret := 0
	for i := 0; i < len(GetCsvMgr().GemstoneLevelConfig); i++ {
		if GetCsvMgr().GemstoneLevelConfig[i].ChaptrGroup == cid {
			if GetCsvMgr().GemstoneLevelConfig[i].LevelIndex > ret {
				ret = GetCsvMgr().GemstoneLevelConfig[i].LevelIndex
			}
		}
	}

	return ret
}

func (self *GemStoneMgr) GetChapterLevelConfig(cid int, level int) *GemstoneLevelConfig {
	for i := 0; i < len(GetCsvMgr().GemstoneLevelConfig); i++ {
		lvlCfg := GetCsvMgr().GemstoneLevelConfig[i]
		if lvlCfg.ChaptrGroup == cid && lvlCfg.LevelIndex == level {
			return lvlCfg
		}
	}

	return nil
}

func (self *GemStoneMgr) GetSweepTariffConfig(times int) *TariffConfig {
	return GetCsvMgr().GetTariffConfig(TariffGemStone, times)
}
