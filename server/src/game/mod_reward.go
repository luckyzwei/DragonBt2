package game

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	REWARD_SET       = "reward_set"       // 设置
	REWARD_GET       = "reward_get"       // 领取
	REWARD_GET_ALL   = "reward_get_all"   // 一键领取
	REWARD_REFRESH   = "reward_refresh"   // 刷新
	REWARD_SEND_INFO = "reward_send_info" //发送信息
	REWARD_RED_POINT = "reward_red_point" //发送红点信息
)

const (
	REWARD_GET_STATE_NONE        = 0 // 未开启
	REWARD_GET_STATE_RUNING_TIME = 1 // 开启未完成
	REWARD_GET_STATE_CAN_GET     = 2 // 完成可领取
)

const (
	REWARD_TASK_TYPE_PERSON = 0 // 个人悬赏类型
	REWARD_TASK_TYPE_TEAM   = 1 // 团队悬赏类型
	REWARD_TASK_TYPE_MAX    = 2 // 类型最大值

	REWARD_PERSON_TASK_MAX   = 5  // 个人任务刷新最大个数
	REWARD_TEAM_TASK_MAX     = 1  // 团队任务刷新最大个数
	REWARD_TEAM_REFRESH_COST = 50 // 个人任务刷新花费
)

type RewardHero struct {
	Index    int   `json:"index"`
	Uid      int64 `json:"uid"`      // 谁的
	HeroKey  int   `json:"herokey"`  // 英雄key值
	HeroID   int   `json:"heroid"`   // 英雄id
	HeroStar int   `json:"herostar"` // 英雄star
	HeroSkin int   `json:"skin"`     // 英雄皮肤
}

type JS_Reward struct {
	ID         int           `json:"id"`         // id
	IsTeam     int           `json:"isteam"`     // 是否是团队
	FinishTime int64         `json:"finishtime"` // 完成时间
	DeleteTime int64         `json:"deletetime"` // 清除时间
	RewardHero []*RewardHero `json:"rewardhero"` // 上阵英雄
	State      int           `json:"state"`
	Items      []PassItem    `json:"items"`
}

// 失败清理所有英雄
func (self *JS_Reward) ClearHero() {
	self.RewardHero = []*RewardHero{}
}

// 检查是否完成
func (self *JS_Reward) CheckFinish() bool {
	timeNow := TimeServer().Unix()
	// 获得配置
	config := GetCsvMgr().GetRewardConfig(self.IsTeam, self.ID)
	if nil == config {
		return false
	}

	attMap := make(map[int]int)
	for _, v := range self.RewardHero {
		// 获得配置
		heroConfig := GetCsvMgr().GetHeroMapConfig(v.HeroID, v.HeroStar)
		if heroConfig == nil {
			continue
		}

		_, ok := attMap[heroConfig.Attribute]
		if !ok {
			attMap[heroConfig.Attribute] = 1
		} else {
			attMap[heroConfig.Attribute]++
		}
	}

	for i, v := range config.NeedCamp {
		if v == 0 {
			continue
		}
		_, ok := attMap[i+1]
		if !ok {
			return false
		}
		if attMap[i+1] < v {
			return false
		}
	}

	// 星级是否达到
	nStar := 0
	for _, v := range self.RewardHero {
		if v.HeroStar >= config.NeedStar {
			nStar++
		}
	}

	// 星级是否匹配
	if nStar < config.StarNum {
		return false
	}

	// 设置完成时间和状态
	self.FinishTime = timeNow + int64(config.SepndTime)
	self.State = REWARD_GET_STATE_RUNING_TIME
	return true
}

// 获得上阵的英雄
func (self *JS_Reward) GetHero(index int) *RewardHero {
	for _, v := range self.RewardHero {
		if v.Index == index {
			return v
		}
	}
	return nil
}

// 悬赏任务
type San_Reward struct {
	Uid       int64  // 角色ID
	Info      string // 任务
	Level     int    // 等级
	TaskCount string // 任务进度

	info      []*JS_Reward              // 任务
	taskCount [REWARD_TASK_TYPE_MAX]int // 等级进度
	redPoint  bool
	DataUpdate
}

// 悬赏
type ModReward struct {
	player *Player // 玩家

	Sql_Reward San_Reward // 信息数据
}

// 获得数据
func (self *ModReward) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_reward` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Reward, "san_reward", self.player.ID)

	if self.Sql_Reward.Uid <= 0 {
		self.Sql_Reward.Uid = self.player.ID
		self.Sql_Reward.info = make([]*JS_Reward, 0)
		self.Sql_Reward.Level = 1
		self.Encode()
		InsertTable("san_reward", &self.Sql_Reward, 0, true)
		self.Sql_Reward.Init("san_reward", &self.Sql_Reward, true)
	} else {
		self.Decode()
		self.Sql_Reward.Init("san_reward", &self.Sql_Reward, true)
	}
}

// 获得数据
func (self *ModReward) OnGetOtherData() {
	// 新手引导特殊处理
	if len(self.Sql_Reward.info) <= 0 && self.Sql_Reward.Level <= 1 {
		self.player.GetModule("reward").(*ModReward).OnRefresh(true)
	}
}

// save
func (self *ModReward) Decode() {
	json.Unmarshal([]byte(self.Sql_Reward.Info), &self.Sql_Reward.info)
	json.Unmarshal([]byte(self.Sql_Reward.TaskCount), &self.Sql_Reward.taskCount)
}

// read
func (self *ModReward) Encode() {
	self.Sql_Reward.Info = HF_JtoA(self.Sql_Reward.info)
	self.Sql_Reward.TaskCount = HF_JtoA(self.Sql_Reward.taskCount)
}

// 存储
func (self *ModReward) OnSave(sql bool) {
	self.Encode()
	self.Sql_Reward.Update(sql)
}

// 消息
func (self *ModReward) OnMsg(ctrl string, body []byte) bool {
	return false
}

// 注册消息
func (self *ModReward) onReg(handlers map[string]func(body []byte)) {
	handlers[REWARD_SET] = self.RewardSet            // 设置
	handlers[REWARD_GET] = self.RewardGet            // 领取
	handlers[REWARD_GET_ALL] = self.RewardGetAll     // 一键领取
	handlers[REWARD_REFRESH] = self.RewardRefresh    // 刷新
	handlers[REWARD_SEND_INFO] = self.SendInfo       // 发送消息
	handlers[REWARD_RED_POINT] = self.RewardRedPoint // 发送消息
}

// 固定刷新
func (self *ModReward) OnRefresh(special bool) {
	self.CheckRewardState()
	self.RefreshPerson(true)
	self.RefreshTeam(special)

	self.Sql_Reward.redPoint = true
	self.SendInfo([]byte{})
}

// 红点信息
func (self *ModReward) RewardRedPoint(body []byte) {
	var backmsg S2C_RewardRedPoint
	backmsg.Cid = REWARD_RED_POINT
	backmsg.RedPoint = self.Sql_Reward.redPoint
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))

	self.Sql_Reward.redPoint = false
}

// 手动刷新
func (self *ModReward) RewardRefresh(body []byte) {
	// 检查悬赏任务状态
	self.CheckRewardState()

	// 获得等级配置
	config := GetCsvMgr().GetRewardForbarLvUpConfig(self.Sql_Reward.Level)
	if config == nil {
		self.player.SendErrInfo("err", "手动刷新获得等级配置"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	csv_vip := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if csv_vip == nil || len(csv_vip.GuildHunting) < 2 {
		return
	}

	// 没有未完成的任务
	if self.GetUnfinishCount(false) <= 0 {
		self.player.SendErrInfo("err", "没有未完成的任务"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	// 货币不足
	if self.player.GetObjectNum(DEFAULT_GEM) < REWARD_TEAM_REFRESH_COST {
		self.player.SendErrInfo("err", "货币不足"+GetCsvMgr().GetText("STR_ERROR"))
		return
	}

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_REWARD_SET_REFRESH, 1, self.Sql_Reward.Level, self.GetUnfinishCount(false), "刷新赏金任务", 0, csv_vip.RewardItask+REWARD_PERSON_TASK_MAX, self.player)

	// 删除货币
	self.player.RemoveObjectEasy(DEFAULT_GEM, REWARD_TEAM_REFRESH_COST, "悬赏刷新", 0, 0, 0)
	// 刷新个人悬赏任务
	self.RefreshPerson(false)

	var backmsg S2C_RewardRefresh
	backmsg.Cid = REWARD_REFRESH
	backmsg.Info = self.Sql_Reward.info
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

//领取奖励
func (self *ModReward) RewardGet(body []byte) {
	var msg C2S_RewardGet
	json.Unmarshal(body, &msg)
	// 检查悬赏任务状态
	self.CheckRewardState()

	// 获得数据
	info := self.GetRewardInfo(msg.ID)
	if nil == info {
		return
	}

	// 获得配置
	config := GetCsvMgr().GetRewardConfig(info.IsTeam, msg.ID)
	if nil == config {
		return
	}

	// 状态不对
	if info.State != REWARD_GET_STATE_CAN_GET {
		return
	}

	heros := []*Hero{}

	// 获得玩家uid
	myUid := self.player.GetUid()
	for _, v := range info.RewardHero {
		// 是自己的英雄
		if myUid == v.Uid {
			hero := self.player.getHero(v.HeroKey)
			if nil != hero {
				// 恢复使用类型
				if config.IsTeam == 1 {
					hero.UseType[HERO_USE_TYPE_REWARD_TEAM] = 0
				} else {
					hero.UseType[HERO_USE_TYPE_REWARD] = 0
				}

				heros = append(heros, hero)
			}
		} else {
			var mastermsg S2M_SupportHeroCancelUse
			mastermsg.Uid = v.Uid
			mastermsg.HeroKeyId = v.HeroKey
			mastermsg.Useruid = self.player.Sql_UserBase.Uid
			GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_CANCEL_USE, &mastermsg)
			//GetSupportHeroMgr().CancelUseHero(v.Uid, v.HeroKey, self.player)
		}
	}

	// 获得等级配置
	oldLevel := self.Sql_Reward.Level
	levelConfig := GetCsvMgr().GetRewardForbarLvUpConfig(self.Sql_Reward.Level)
	if levelConfig != nil {
		// 如果是个人任务
		if info.IsTeam == REWARD_TASK_TYPE_PERSON {
			if config.Color >= levelConfig.Persontaskstar {
				self.Sql_Reward.taskCount[REWARD_TASK_TYPE_PERSON] += 1
				self.CheckRewardLevel()
			}
		} else { // 如果是团队任务
			if config.Color >= levelConfig.Teamtaskstar {
				self.Sql_Reward.taskCount[REWARD_TASK_TYPE_TEAM] += 1
				self.CheckRewardLevel()
			}
		}
	}
	logStr := ""
	if info.IsTeam == REWARD_TASK_TYPE_PERSON {
		self.player.HandleTask(TASK_TYPE_REWARD_TASK_GET, 1, 1, config.Color)
		self.player.HandleTask(TASK_TYPE_REWARD_TASK_GET_EQUAL, 1, 1, config.Color)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_REWARD_GET_NO_TEAM, config.ID, oldLevel, self.Sql_Reward.Level, "领取个人赏金任务", 0, 0, self.player)
		logStr = "领取个人赏金任务奖励"
	} else {
		self.player.HandleTask(TASK_TYPE_REWARD_TASK_GET, 1, 2, config.Color)
		self.player.HandleTask(TASK_TYPE_REWARD_TASK_GET_EQUAL, 1, 2, config.Color)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_REWARD_GET_TEAM, config.ID, oldLevel, self.Sql_Reward.Level, "领取团队赏金任务", 0, 0, self.player)
		logStr = "领取团队赏金任务奖励"
	}
	ret := make(map[int]*Item, 0)
	AddItemMapHelper2(ret, info.Items)
	var backmsg S2C_RewardGet
	backmsg.Item = self.player.AddObjectItemMap(ret, logStr, config.ID, self.Sql_Reward.Level, 0)
	backmsg.Cid = REWARD_GET
	backmsg.ID = msg.ID
	backmsg.Level = self.Sql_Reward.Level
	backmsg.TaskCount = self.Sql_Reward.taskCount
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))

	//删除该条悬赏
	self.DeleteInfo(msg.ID)

	var backmsg2 S2C_ResonanceCrystaUpdateResonance
	backmsg2.Cid = RESONANCE_CRYSTAL_UPDATE_RESONANCE
	backmsg2.Heros = heros
	self.player.SendMsg(backmsg2.Cid, HF_JtoB(&backmsg2))
}

//领取奖励
func (self *ModReward) RewardGetAll(body []byte) {
	var msg C2S_RewardGetAll
	json.Unmarshal(body, &msg)

	vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if vipcsv == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PVP_VIP_CONFIGURATION_TABLE_ERROR"))
		return
	}

	if vipcsv.RewardOnekey != 1 {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PVP_VIP_CONFIGURATION_TABLE_ERROR"))
		return
	}

	// 检查悬赏任务状态
	self.CheckRewardState()
	heros := []*Hero{}
	backitem := make(map[int]*Item, 0)
	var backmsg S2C_RewardGetAll
	len := len(self.Sql_Reward.info)
	for i := len - 1; i >= 0; i-- {
		info := self.Sql_Reward.info[i]
		// 状态不对
		if info.State != REWARD_GET_STATE_CAN_GET {
			continue
		}

		if msg.Type != info.IsTeam {
			continue
		}
		// 获得配置
		config := GetCsvMgr().GetRewardConfig(info.IsTeam, info.ID)
		if nil == config {
			continue
		}

		// 获得玩家uid
		myUid := self.player.GetUid()
		for _, v := range info.RewardHero {
			// 是自己的英雄
			if myUid == v.Uid {
				hero := self.player.getHero(v.HeroKey)
				if nil != hero {
					// 恢复使用类型
					if config.IsTeam == 1 {
						hero.UseType[HERO_USE_TYPE_REWARD_TEAM] = 0
					} else {
						hero.UseType[HERO_USE_TYPE_REWARD] = 0
					}

					heros = append(heros, hero)
				}
			} else {
				var mastermsg S2M_SupportHeroCancelUse
				mastermsg.Uid = v.Uid
				mastermsg.HeroKeyId = v.HeroKey
				mastermsg.Useruid = self.player.Sql_UserBase.Uid
				GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_CANCEL_USE, &mastermsg)
				//GetSupportHeroMgr().CancelUseHero(v.Uid, v.HeroKey, self.player)
			}
		}

		// 获得等级配置
		oldLevel := self.Sql_Reward.Level
		levelConfig := GetCsvMgr().GetRewardForbarLvUpConfig(self.Sql_Reward.Level)
		if levelConfig != nil {
			// 如果是个人任务
			if info.IsTeam == REWARD_TASK_TYPE_PERSON {
				if config.Color >= levelConfig.Persontaskstar {
					self.Sql_Reward.taskCount[REWARD_TASK_TYPE_PERSON] += 1
					self.CheckRewardLevel()
				}
			} else { // 如果是团队任务
				if config.Color >= levelConfig.Teamtaskstar {
					self.Sql_Reward.taskCount[REWARD_TASK_TYPE_TEAM] += 1
					self.CheckRewardLevel()
				}
			}
		}
		if info.IsTeam == REWARD_TASK_TYPE_PERSON {
			self.player.HandleTask(TASK_TYPE_REWARD_TASK_GET, 1, 1, config.Color)
			self.player.HandleTask(TASK_TYPE_REWARD_TASK_GET_EQUAL, 1, 1, config.Color)
			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_REWARD_GET_NO_TEAM_ALL, config.ID, oldLevel, self.Sql_Reward.Level, "一键领取个人赏金任务", 0, 0, self.player)
		} else {
			self.player.HandleTask(TASK_TYPE_REWARD_TASK_GET, 1, 2, config.Color)
			self.player.HandleTask(TASK_TYPE_REWARD_TASK_GET_EQUAL, 1, 2, config.Color)
			GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_REWARD_GET_TEAM_ALL, config.ID, oldLevel, self.Sql_Reward.Level, "一键领取团队赏金任务", 0, 0, self.player)
		}

		ret := make(map[int]*Item, 0)
		AddItemMapHelper2(ret, info.Items)
		AddItemMapHelper2(backitem, info.Items)
		self.player.AddObjectItemMap(ret, "悬赏领取", 0, 0, 0)

		backmsg.ID = append(backmsg.ID, info.ID)
		//删除该条悬赏
		self.DeleteInfo(info.ID)
	}

	for _, v := range backitem {
		backmsg.Item = append(backmsg.Item, PassItem{v.ItemId, v.ItemNum})
	}

	backmsg.Cid = REWARD_GET_ALL
	backmsg.Level = self.Sql_Reward.Level
	backmsg.TaskCount = self.Sql_Reward.taskCount
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))

	var backmsg2 S2C_ResonanceCrystaUpdateResonance
	backmsg2.Cid = RESONANCE_CRYSTAL_UPDATE_RESONANCE
	backmsg2.Heros = heros
	self.player.SendMsg(backmsg2.Cid, HF_JtoB(&backmsg2))
}

// 发送信息
func (self *ModReward) SendInfo(body []byte) {
	self.CheckRewardState()
	var backmsg S2C_RewardInfo
	backmsg.Cid = REWARD_SEND_INFO
	backmsg.Info = self.Sql_Reward.info
	backmsg.Level = self.Sql_Reward.Level
	backmsg.TaskCount = self.Sql_Reward.taskCount
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))
}

func (self *ModReward) IsUseSupportHero(uid int64, keyid int) bool {
	for _, v := range self.Sql_Reward.info {
		for _, t := range v.RewardHero {
			if t.Uid == uid && t.HeroKey == keyid {
				return true
			}
		}
	}

	return false
}

//消息使用英雄
func (self *ModReward) RewardSet(body []byte) {
	var msg C2S_RewardSet
	json.Unmarshal(body, &msg)
	// 检查悬赏任务状态
	self.CheckRewardState()

	// 获得数据
	info := self.GetRewardInfo(msg.ID)
	if nil == info {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("没有次任务"))
		return
	}
	// 获得配置
	config := GetCsvMgr().GetRewardConfig(info.IsTeam, msg.ID)
	if nil == config {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("没有次任务配置"))
		return
	}

	// 已经有英雄上阵
	if len(info.RewardHero) > 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("已经有英雄上阵"))
		return
	}

	// 状态错误
	if info.State != REWARD_GET_STATE_NONE {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("状态错误"))
		return
	}

	// 参数错误
	nLen := len(msg.Uids)
	nCount := 0
	for _, v := range config.NeedCamp {
		if v != 0 {
			nCount += v
		}
	}
	if nLen != len(msg.HeroKeys) {
		if self.Sql_Reward.Level != 1 {
			if nLen != nCount {
				return
			}
		}
	}

	heros := []*Hero{}

	useHeros := []int{}
	addHero := []*RewardHero{}
	attMap := make(map[int]int)
	//我的uid
	myUid := self.player.GetUid()
	for index := 0; index < nLen; index++ {
		heroKey := msg.HeroKeys[index]
		uid := msg.Uids[index]
		heroID := 0
		heroStar := 0
		heroSkin := 0

		// 借用别人的英雄
		if myUid != uid {
			if config.IsTeam == 0 {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("不是团队任务无法使用佣兵"))
				return
			}

			// 不是最后一个则返回
			if index != nLen-1 {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("佣兵顺序错误"))
				return
			}

			// 判断是不是好友或者一个军团
			if !self.player.GetModule("support").(*ModSupportHero).IsFriendOrUnion(uid) {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("不是好友或者一个军团"))
				return
			}

			// 寻找被借人的data
			data := GetSupportHeroMgr().GetPlayerData(uid, false)
			if nil == data {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("未找到玩家"))
				return
			}

			// 判断有没有这个英雄提供援助
			find := false
			heroIndex := -1
			for i, v := range data {
				if v.HeroKey == heroKey {
					find = true
					heroIndex = i
					break
				}
			}
			if !find {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("没有这个英雄提供援助"))
				return
			}

			if self.IsUseSupportHero(uid, heroKey) {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("英雄已派遣"))
				return
			}

			//// 超出借出数量
			//if self.player.GetModule("support").(*ModSupportHero).GetUseCount(uid) >= SUPPORT_HERO_USE_MAX {
			//	self.player.SendErrInfo("err", GetCsvMgr().GetText("超出借出数量"))
			//	return
			//}

			// 用以判断该英雄是不是属于该悬赏的的
			heroID = data[heroIndex].HeroID
			heroStar = data[heroIndex].HeroStar
			heroSkin = data[heroIndex].HeroSkin
		} else {
			// 是最后一个则返回
			if config.IsTeam == 1 {
				if index == nLen-1 {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("没有上阵佣兵"))
					return
				}
			}

			// 获得自己的英雄
			hero := self.player.getHero(heroKey)
			if nil == hero {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("没找到英雄"))
				return
			}

			if config.IsTeam == 1 {
				if hero.UseType[HERO_USE_TYPE_REWARD_TEAM] == 1 {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("英雄已派遣"))
					return
				}
			} else {
				if hero.UseType[HERO_USE_TYPE_REWARD] == 1 {
					self.player.SendErrInfo("err", GetCsvMgr().GetText("英雄已派遣"))
					return
				}
			}

			heroID = hero.HeroId
			heroStar = hero.GetStar()
			heroSkin = hero.Skin
		}

		for _, usehero := range useHeros {
			if usehero == heroID {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("英雄已派遣"))
				return
			}
		}

		// 英雄配置
		heroConfig := GetCsvMgr().GetHeroMapConfig(heroID, heroStar)
		if heroConfig == nil {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("英雄配置错误"))
			return
		}

		_, ok := attMap[heroConfig.Attribute]
		if !ok {
			attMap[heroConfig.Attribute] = 1
		} else {
			attMap[heroConfig.Attribute]++
		}
		//// 属性不匹配
		//if heroConfig.Attribute != config.NeedCamp[index] {
		//	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		//	return
		//}

		useHeros = append(useHeros, heroID)

		addHero = append(addHero, &RewardHero{index + 1, uid, heroKey, heroID, heroStar, heroSkin})
	}

	for i, v := range config.NeedCamp {
		if v == 0 {
			continue
		}
		_, ok := attMap[i+1]
		if !ok {
			self.player.SendErrInfo("err", "属性不匹配")
			return
		}
		if attMap[i+1] < v {
			self.player.SendErrInfo("err", "属性不匹配")
			return
		}
	}

	// 设置英雄级使用状态
	info.RewardHero = addHero
	for _, v := range addHero {
		if myUid == v.Uid {
			hero := self.player.getHero(v.HeroKey)
			if nil != hero {
				if config.IsTeam == 1 {
					hero.UseType[HERO_USE_TYPE_REWARD_TEAM] = 1
				} else {
					hero.UseType[HERO_USE_TYPE_REWARD] = 1
				}

				heros = append(heros, hero)
			}
		}
	}

	// 检测是否完成
	if !info.CheckFinish() {
		//取消使用状态
		myUid := self.player.GetUid()
		for _, v := range info.RewardHero {
			if myUid == v.Uid {
				hero := self.player.getHero(v.HeroKey)
				if nil != hero {
					if config.IsTeam == 1 {
						hero.UseType[HERO_USE_TYPE_REWARD_TEAM] = 0
					} else {
						hero.UseType[HERO_USE_TYPE_REWARD] = 0
					}
				}
			} else {
				var mastermsg S2M_SupportHeroCancelUse
				mastermsg.Uid = v.Uid
				mastermsg.HeroKeyId = v.HeroKey
				mastermsg.Useruid = self.player.Sql_UserBase.Uid
				GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_CANCEL_USE, &mastermsg)
				//GetSupportHeroMgr().CancelUseHero(v.Uid, v.HeroKey, self.player)
			}
		}
		// 清理英雄
		info.ClearHero()
		self.player.SendErrInfo("err", GetCsvMgr().GetText("检测是失败"))
		return
	}

	if config.IsTeam != 0 {
		var mastermsg S2M_SupportHeroUse
		mastermsg.Uid = msg.Uids[len(msg.Uids)-1]
		mastermsg.HeroKeyId = msg.HeroKeys[len(msg.HeroKeys)-1]
		mastermsg.Useruid = self.player.Sql_UserBase.Uid
		mastermsg.Username = self.player.Sql_UserBase.UName
		mastermsg.Type = HERO_SUPPORT_TYPE_REWARD
		mastermsg.Endtime = info.FinishTime

		ret := GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_USE, &mastermsg)
		if ret == nil || ret.RetCode != UNION_SUCCESS {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("设置失败"))
			return
		}
		//// 简单设置提供者的数据
		//if !GetSupportHeroMgr().UseHero(msg.Uids[len(msg.Uids)-1], msg.HeroKeys[len(msg.HeroKeys)-1], self.player, HERO_SUPPORT_TYPE_REWARD, info.FinishTime) {
		//	self.player.SendErrInfo("err", GetCsvMgr().GetText("设置失败"))
		//	return
		//}

		self.player.HandleTask(TASK_TYPE_REWARD_TASK_SET, 2, 0, 0)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_REWARD_SET_TEAM, config.ID, self.Sql_Reward.Level, 0, "团队赏金任务派遣", 0, 0, self.player)

	} else {
		self.player.HandleTask(TASK_TYPE_REWARD_TASK_SET, 1, 0, 0)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_REWARD_SET_NO_TEAM, config.ID, self.Sql_Reward.Level, 0, "个人赏金任务派遣", 0, 0, self.player)
	}

	var backmsg S2C_RewardSet
	backmsg.Cid = REWARD_SET
	backmsg.Uids = msg.Uids
	backmsg.ID = config.ID
	backmsg.HeroKeys = msg.HeroKeys
	backmsg.Info = append(backmsg.Info, info)
	self.player.SendMsg(backmsg.Cid, HF_JtoB(&backmsg))

	var backmsg2 S2C_ResonanceCrystaUpdateResonance
	backmsg2.Cid = RESONANCE_CRYSTAL_UPDATE_RESONANCE
	backmsg2.Heros = heros
	self.player.SendMsg(backmsg2.Cid, HF_JtoB(&backmsg2))
}

//获得数据
func (self *ModReward) GetRewardInfo(id int) *JS_Reward {
	for _, v := range self.Sql_Reward.info {
		if v.ID == id {
			return v
		}
	}
	return nil
}

//删除数据
func (self *ModReward) DeleteInfo(id int) bool {
	for i, v := range self.Sql_Reward.info {
		if v.ID == id {
			self.Sql_Reward.info = append(self.Sql_Reward.info[:i], self.Sql_Reward.info[i+1:]...)
			return true
		}
	}
	return false
}

//创建
func (self *ModReward) NewInfo(startTime int64, isteam, id int) {
	// 获得配置
	config := GetCsvMgr().GetRewardConfig(isteam, id)
	if config != nil {
		data := JS_Reward{}
		data.ID = id
		data.IsTeam = isteam
		data.DeleteTime = startTime + int64(config.ElapsedTime)
		data.Items = self.RandomAward(isteam, config.Group)

		self.Sql_Reward.info = append(self.Sql_Reward.info, &data)
	}
}

func (self *ModReward) RandomAward(isteam int, group int) []PassItem {
	ret := []PassItem{}
	configs := GetCsvMgr().GetRewardAwardConfig(isteam, group)
	if len(configs) > 0 {
		nTotalChance := 0
		for _, v := range configs {
			nTotalChance += v.P
		}

		nRandNum := HF_GetRandom(nTotalChance)
		total := 0
		for _, v := range configs {
			total += v.P
			if nRandNum < total {
				ret = append(ret, PassItem{v.Prize, v.Prizenum})
				return ret
			}
		}
	}
	return ret
}

// 检查完成状态
func (self *ModReward) CheckRewardState() {
	timeNow := TimeServer().Unix()
	// 获得玩家uid
	myUid := self.player.GetUid()
	nLen := len(self.Sql_Reward.info)
	for i := nLen - 1; i >= 0; i-- {
		// 过期删除
		if self.Sql_Reward.info[i].DeleteTime != 0 &&
			self.Sql_Reward.info[i].DeleteTime <= timeNow &&
			self.Sql_Reward.info[i].State == REWARD_GET_STATE_NONE {
			// 获得配置
			config := GetCsvMgr().GetRewardConfig(self.Sql_Reward.info[i].IsTeam, self.Sql_Reward.info[i].ID)
			if nil != config {
				for _, v := range self.Sql_Reward.info[i].RewardHero {
					// 是自己的英雄
					if myUid == v.Uid {
						hero := self.player.getHero(v.HeroKey)
						if nil != hero {
							// 恢复使用类型
							if config.IsTeam == 1 {
								hero.UseType[HERO_USE_TYPE_REWARD_TEAM] = 0
							} else {
								hero.UseType[HERO_USE_TYPE_REWARD] = 0
							}
						}
					} else {
						var mastermsg S2M_SupportHeroCancelUse
						mastermsg.Uid = v.Uid
						mastermsg.HeroKeyId = v.HeroKey
						mastermsg.Useruid = self.player.Sql_UserBase.Uid
						GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_CANCEL_USE, &mastermsg)
						//GetSupportHeroMgr().CancelUseHero(v.Uid, v.HeroKey, self.player)
					}
				}
			}

			self.Sql_Reward.info = append(self.Sql_Reward.info[:i], self.Sql_Reward.info[i+1:]...)
			continue
		}

		// 完成设置状态 并设置支援英雄
		if self.Sql_Reward.info[i].FinishTime != 0 && // 有完成时间
			self.Sql_Reward.info[i].State == REWARD_GET_STATE_RUNING_TIME && // 检查状态
			self.Sql_Reward.info[i].FinishTime <= timeNow { // 时间未过期
			self.Sql_Reward.info[i].State = REWARD_GET_STATE_CAN_GET // 设置完成状态
			if self.Sql_Reward.info[i].IsTeam != 0 { // 如果是团队任务
				if self.Sql_Reward.info[i].RewardHero[len(self.Sql_Reward.info[i].RewardHero)-1].Uid != self.player.GetUid() { //最后一个英雄不是自己英雄
					var mastermsg S2M_SupportHeroCancelUse
					mastermsg.Uid = self.Sql_Reward.info[i].RewardHero[len(self.Sql_Reward.info[i].RewardHero)-1].Uid
					mastermsg.HeroKeyId = self.Sql_Reward.info[i].RewardHero[len(self.Sql_Reward.info[i].RewardHero)-1].HeroKey
					mastermsg.Useruid = self.player.Sql_UserBase.Uid
					GetMasterMgr().FriendPRC.SupportHeroAction(RPC_SUPPORT_HERO_CANCEL_USE, &mastermsg)

					//GetSupportHeroMgr().CancelUseHero(self.Sql_Reward.info[i].RewardHero[len(self.Sql_Reward.info[i].RewardHero)-1].Uid,
					//	self.Sql_Reward.info[i].RewardHero[len(self.Sql_Reward.info[i].RewardHero)-1].HeroKey, self.player)
					continue
				}
			}
		}
	}
}

// 检查升级
func (self *ModReward) CheckRewardLevel() {
	// 获得等级配置
	config := GetCsvMgr().GetRewardForbarLvUpConfig(self.Sql_Reward.Level)
	if config == nil {
		return
	}

	// 升到顶了
	if config.Persontasknum <= 0 && config.Teamtasknum <= 0 {
		return
	}

	// 检查数量是否达到
	team := self.Sql_Reward.taskCount[REWARD_TASK_TYPE_TEAM]
	teamDone := false
	if team >= config.Teamtasknum {
		teamDone = true
	}
	person := self.Sql_Reward.taskCount[REWARD_TASK_TYPE_PERSON]
	personDone := false
	if person >= config.Persontasknum {
		personDone = true
	}
	if personDone && teamDone {
		self.Sql_Reward.Level += 1
		self.Sql_Reward.taskCount[REWARD_TASK_TYPE_TEAM] = 0
		self.Sql_Reward.taskCount[REWARD_TASK_TYPE_PERSON] = 0

		// 新手引导特殊处理
		if self.Sql_Reward.Level == 2 {
			self.RefreshPerson(true)
			self.SendInfo([]byte{})
		}
		self.player.HandleTask(TASK_TYPE_REWARD_LEVLE, self.Sql_Reward.Level, 0, 0)
	}
}

// 刷新个人悬赏
func (self *ModReward) RefreshPerson(isall bool) {
	// 刷新全部则清理全部
	min := 0
	if isall {
		nLen := len(self.Sql_Reward.info)
		for i := nLen - 1; i >= 0; i-- {
			if self.Sql_Reward.info[i].IsTeam == 0 {
				if self.Sql_Reward.info[i].State == REWARD_GET_STATE_NONE {
					self.Sql_Reward.info = append(self.Sql_Reward.info[:i], self.Sql_Reward.info[i+1:]...)
				} else if self.Sql_Reward.info[i].State == REWARD_GET_STATE_RUNING_TIME {
					min++
				}
			}
		}
	}

	// 计算刷新出的个数
	nCount := 0
	if isall {
		// 获得等级配置
		config := GetCsvMgr().GetRewardForbarLvUpConfig(self.Sql_Reward.Level)
		if config == nil {
			return
		}

		csv_vip := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
		if csv_vip == nil || len(csv_vip.GuildHunting) < 2 {
			return
		}
		nCount = csv_vip.RewardItask + REWARD_PERSON_TASK_MAX - min
	} else {
		nCount = self.GetUnfinishCount(true)
	}

	// 随机任务
	self.RandomTask(self.GetTime(), REWARD_TASK_TYPE_PERSON, self.Sql_Reward.Level, nCount)

}

//刷新队伍悬赏
func (self *ModReward) RefreshTeam(special bool) {
	csv_vip := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if csv_vip == nil || len(csv_vip.GuildHunting) < 2 {
		return
	}
	// 先刷一轮
	nCount := csv_vip.RewardTeamtask + REWARD_TEAM_TASK_MAX
	self.RandomTask(self.GetTime(), REWARD_TASK_TYPE_TEAM, self.Sql_Reward.Level, nCount)

	// 是否特殊处理 非特殊处理则需要补刷
	if special {
		return
	}

	// 检测补刷
	// 当前时间
	now := TimeServer()
	// 最后一次更新的时间 玩家可能一直在线更新 也可能很久没上线
	tll, _ := time.ParseInLocation(DATEFORMAT, self.player.Sql_UserBase.LastUpdTime, time.Local)
	// 最后一次刷新时间
	lastupdate := time.Date(tll.Year(), tll.Month(), tll.Day(), 5, 0, 0, 0, now.Location()).Unix()
	// 如果时间小于五点 说明上次刷的时间是上一轮 这一轮还没刷
	if tll.Hour() < 5 {
		lastupdate -= DAY_SECS
	}

	// 取当天的时间 五点 如果时间在五点前 算昨天 今天的还没开始刷
	today := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
	if now.Hour() < 5 {
		today -= DAY_SECS
	}

	// 时间相同 直接返回 说明不用刷
	if lastupdate == today {
		return
	}
	// 两种情况 在线时经过五点 以及五点以后上线补刷
	//
	// 当一直在线时 同一天 则说明一直在线 或者是四点以前上过线 不需要补刷

	// 补刷天数要扣掉一开始补刷的一次
	day := (today-lastupdate)/DAY_SECS - 1

	for i := int64(1); i <= day; i++ {
		nCount := csv_vip.RewardTeamtask + REWARD_TEAM_TASK_MAX
		self.RandomTask(lastupdate+i*DAY_SECS, REWARD_TASK_TYPE_TEAM, self.Sql_Reward.Level, nCount)
	}
}

// 随机任务通用
func (self *ModReward) RandomTask(startTime int64, isTeam, lv, count int) []int {
	// 获得等级配置
	config := GetCsvMgr().GetRewardForbarLvUpConfig(lv)
	if config == nil {
		return nil
	}

	// 随机出颜色
	taskGroup := []int{}
	for i := 0; i < count; i++ {
		group := self.RandomGroup(config)
		if group != 0 {
			taskGroup = append(taskGroup, group)
		}
	}

	//随机出任务
	taskID := []int{}
	for _, v := range taskGroup {
		id := self.RandomTaskByGroup(isTeam, v)
		if id != 0 {
			taskID = append(taskID, id)
			self.NewInfo(startTime, isTeam, id)
		}
	}

	return taskID
}

// 通过颜色随机出一个任务列表中没有的任务
func (self *ModReward) RandomTaskByGroup(isTeam, group int) int {
	config := GetCsvMgr().GetRewardForbarGroupConfig(isTeam, group)
	data := []int{}
	for _, v := range config {
		info := self.GetRewardInfo(v.ID)
		if info == nil {
			data = append(data, v.ID)
		}
	}

	// 计算列表
	if len(data) <= 0 {
		return 0
	}

	// 随机
	index := HF_GetRandom(len(data))
	return data[index]
}

// 根据概率随机出颜色
func (self *ModReward) RandomGroup(config *RewardForbarLvUpConfig) int {
	nTotalChance := 0
	for _, v := range config.Renovatepro {
		nTotalChance += v
	}

	nRandNum := HF_GetRandom(nTotalChance)
	total := 0
	for i, v := range config.Renovatestar {
		total += config.Renovatepro[i]
		if nRandNum < total {
			return v
		}
	}
	return 0
}

// 获得未完成任务数量
func (self *ModReward) GetUnfinishCount(delete bool) int {
	nCount := 0
	nLen := len(self.Sql_Reward.info)
	for i := nLen - 1; i >= 0; i-- {
		if self.Sql_Reward.info[i].IsTeam == 0 &&
			self.Sql_Reward.info[i].State == REWARD_GET_STATE_NONE {
			// 是否删除
			if delete {
				self.Sql_Reward.info = append(self.Sql_Reward.info[:i], self.Sql_Reward.info[i+1:]...)
			}
			nCount++
		}
	}
	return nCount
}

// 获得数量
func (self *ModReward) GetUseCount(uid int64) int {
	nCount := 0
	for _, v := range self.Sql_Reward.info {
		if v.IsTeam == 0 {
			continue
		}

		for _, g := range v.RewardHero {
			if g.Uid == uid {
				nCount++
			}
		}
	}
	return nCount
}

// gm完成
func (self *ModReward) GmFinishTask() {
	now := TimeServer().Unix()
	for _, v := range self.Sql_Reward.info {
		if v.FinishTime > 0 {
			v.FinishTime = now
		}
	}
}

func (self *ModReward) GmRefreshTask() {
	self.OnRefresh(true)
}

func (self *ModReward) VipLevelChange(oldlevel int) {
	vipcsv := GetCsvMgr().GetVipConfig(self.player.Sql_UserBase.Vip)
	if vipcsv == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PVP_VIP_CONFIGURATION_TABLE_ERROR"))
		return
	}

	oldvipcsv := GetCsvMgr().GetVipConfig(oldlevel)
	if oldvipcsv == nil {
		self.player.SendErr(GetCsvMgr().GetText("STR_MOD_PVP_VIP_CONFIGURATION_TABLE_ERROR"))
		return
	}

	timestamp := self.GetTime()

	personCount := vipcsv.RewardItask - oldvipcsv.RewardItask
	// 随机任务
	self.RandomTask(timestamp, REWARD_TASK_TYPE_PERSON, self.Sql_Reward.Level, personCount)

	teamCount := vipcsv.RewardTeamtask - oldvipcsv.RewardTeamtask

	self.RandomTask(timestamp, REWARD_TASK_TYPE_TEAM, self.Sql_Reward.Level, teamCount)

	self.SendInfo([]byte{})
}

func (self *ModReward) GetTime() int64 {
	now := TimeServer()
	timestamp := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
	if now.Hour() < 5 {
		timestamp -= DAY_SECS
	}
	return timestamp
}
