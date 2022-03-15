package game

import (
	"encoding/json"
	"fmt"
)

type ModGm struct {
	player *Player
}

func (self *ModGm) OnGetData(player *Player) {
	self.player = player
}

func (self *ModGm) OnGetOtherData() {
}

func (self *ModGm) Gmgogogo() {
	for _, value := range GetCsvMgr().Data["1GM"] {
		types := HF_Atoi(value["type"])
		if types == 1 {
			self.player.GetModule("hero").(*ModHero).AddHero(HF_Atoi(value["id"]), 4, 0, 0, "gm礼包")
		} else if types == 2 {
			self.player.AddObject(HF_Atoi(value["id"]), HF_Atoi(value["num"]), 4, 0, 0, "gm礼包")
		}
	}

	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GM_PLEASE_RESTART_AFTER_THE_GOODS"))
}

func (self *ModGm) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "kickout01":
		self.player.SendRet2("kickout01")
		//self.SendMsg("shutdown", []byte(""))
		return true
	case "autologin":
		var c2s_msg C2S_GMAutoLogin
		json.Unmarshal(body, &c2s_msg)
		self.player.SendRet("autologin", LOGIC_TRUE)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_AUTO_LOGIN, c2s_msg.LoginType, 0, 0, "断线重连", 0, 0, self.player)
		return true
	case "swtichserver":
		var c2s_msg C2S_SwitchServer
		json.Unmarshal(body, &c2s_msg)
		GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_USER_SWITCH_SERVER, c2s_msg.ServerID, 0, 0, "切换服务器", 0, 0, self.player)
		return true
	case "gmgogogo":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		self.Gmgogogo()
		return true
	case "gmaddexp":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		var c2s_msg C2S_Gmaddexp
		json.Unmarshal(body, &c2s_msg)

		self.GmAddExp(c2s_msg.Addexp)
		return true
	case "gm_trainwar_num":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		self.player.GetModule("pass").(*ModPass).OnRefresh()
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GM_SUCCESS"))
		return true
	case "gm_reset_dailytask":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		self.player.GetModule("task").(*ModTask).OnRefresh()
		self.player.GetModule("union").(*ModUnion).OnRefresh()
		self.player.GetModule("onhook").(*ModOnHook).OnRefresh()
		self.player.Refresh()

		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_GM_SUCCESS"))
		return true
	case "gmstr":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		var c2s_msg C2S_GMStr
		json.Unmarshal(body, &c2s_msg)
		self.GMStr(c2s_msg.Gmstr, c2s_msg.Herolst)
		return true
	case "gmchat":
		var c2s_msg C2S_GMStr
		json.Unmarshal(body, &c2s_msg)
		if len(c2s_msg.Gmstr) > 200 {
			self.player.SendErr(GetCsvMgr().GetText("STR_MOD_GM_DONT_MAKE_IT_TOO_LONG"))
			return true
		}
		GetServer().sendSysChat(c2s_msg.Gmstr)
		return true
	case "gm_pass_mission":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		// [{"type":"gm_pass_mission","chid":10}]
		var msg C2S_PassMission
		json.Unmarshal(body, &msg)
		self.player.GetModule("pass").(*ModPass).GMPassChapter(msg.ChpaterId)
		self.player.SendInfo("updateuserinfo")
		return true
	case "gm_clear_teampos":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		// [{"type":"gm_clear_teampos"}]
		self.GMTokeOffAllEuiup()
		return true
	case "gm_set_union_exp_limit":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		var msg C2S_AddUnionExpLimit
		// [{"type":"gm_set_union_exp_limit","add_num":10}]
		json.Unmarshal(body, &msg)
		self.player.GetModule("union").(*ModUnion).GMAddUnionExpLimit(msg.AddNum)
		return true
	case "gm_mail_test":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		pMail := self.player.GetModule("mail").(*ModMail)
		if pMail == nil {
			return true
		}
		var msg C2S_GMTestMail
		json.Unmarshal(body, &msg)

		if len(msg.ItemId) != len(msg.ItemNum) {
			return true
		}

		outItem := make([]PassItem, 0)
		for i := 0; i < len(msg.ItemNum); i++ {
			if msg.ItemNum[i] > 0 {
				outItem = append(outItem, PassItem{msg.ItemId[i], msg.ItemNum[i]})
			}
		}
		pMail.AddMail(1, 1, 0, msg.Title, msg.Content, "邮件测试", outItem, true, 0)
		return true
	case "gm_add_union_exp":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		var msg C2S_AddUnionExpLimit
		json.Unmarshal(body, &msg)
		self.player.GetModule("union").(*ModUnion).GMAddUnionExp(msg.AddNum)
		return true
	case "gm_add_union_activity":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		var msg C2S_AddUnionExpLimit
		json.Unmarshal(body, &msg)
		self.player.GetModule("union").(*ModUnion).GMAddUnionActivity(msg.AddNum)
		return true
	case "gm_set_reward_level":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		var msg C2S_AddUnionExpLimit
		json.Unmarshal(body, &msg)
		self.player.GetModule("reward").(*ModReward).Sql_Reward.Level = msg.AddNum
		return true
	case "gm_finish_reward":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}

		self.player.GetModule("reward").(*ModReward).GmFinishTask()
		return true
	case "gm_refresh_reward":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}

		self.player.GetModule("reward").(*ModReward).GmRefreshTask()
		return true
	case "gm_set_tower_level":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		var msg C2S_SetTowerLevel
		json.Unmarshal(body, &msg)
		self.player.GetModule("tower").(*ModTower).GmSetToweLevel(msg.Type, msg.Level)
		return true
	case "gmresetpit":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		//GetServer().Notice("测试公告123456789", 0, 0)

		self.player.GetModule("newpit").(*ModNewPit).GmResetPit(1)
		return true
	case "gmresetpit2":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}

		self.player.GetModule("newpit").(*ModNewPit).GmResetPit(3)
		return true
	case "gmfinishgrowgift":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		var msg C2S_AddUnionExpLimit
		json.Unmarshal(body, &msg)
		self.player.GetModule("growthgift").(*ModGrowthGift).GMFinishGrowthGift(msg.AddNum)
		return true
	case "gm_help":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		var msg C2S_GMSuperHelp
		json.Unmarshal(body, &msg)
		self.GmSuperHelp(msg.Index)
		return true
	case "gmshadowreset":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		self.GmShadowReset()
		return true
	case "gm_nobilityup":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		var msg C2S_GMNobilityUp
		json.Unmarshal(body, &msg)
		self.player.GetModule("nobilitytask").(*ModNobilityTask).GmLevelUp(msg.Level)
		return true
	case "gm_interstellar":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		var msg C2S_GMInterstellar
		json.Unmarshal(body, &msg)
		self.player.GetModule("interstellar").(*ModInterStellar).GMOrder(&msg)
		return true
	case "gm_interstellarall":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		self.player.GetModule("interstellar").(*ModInterStellar).GMOrderAll()
		return true
	case "gminstancepass":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		var msg C2S_GMInstancePass
		json.Unmarshal(body, &msg)
		self.player.GetModule("instance").(*ModInstance).GmGMInstancePass(&msg)
		return true
	case "gmclearhero":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		self.player.GetModule("hero").(*ModHero).GmClearHero()
		return true
	case "gmclearitem":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		self.player.GetModule("bag").(*ModBag).GmClearItem()
		self.player.GetModule("equip").(*ModEquip).GmClearItem()
		self.player.GetModule("team").(*ModTeam).GmClearItem()
		return true
	case "gmtemptest":
		if self.player.Sql_UserBase.Uid >= GetServer().Con.GM {
			return true
		}
		self.player.GetModule("activitybossfestival").(*ModActivityBossFestival).GMReset()
		return true
	}
	return false
}

func (self *ModGm) OnSave(sql bool) {
}

func (self *ModGm) GmAddExp(exp int) {
	self.player.AddExp(exp, 4, 0, "gm礼包")
	self.player.SendInfo("updateuserinfo")
}

func (self *ModGm) GMStr(gmstr string, herolst []int) {
	if gmstr == "getfightnum" {

		fnum := self.player.GetModule("hero").(*ModHero).GetFight(herolst)

		self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_GM_FIGHT"), fnum))

	}
}

func (self *ModGm) SendDebugStr(gmstr string) {
	var msg S2C_DebugString
	msg.Cid = "debugstr"
	msg.Debugstr = gmstr
	smsg, _ := json.Marshal(&msg)
	self.player.SendMsg("debugstr", smsg)

}

func (self *ModGm) GMTokeOffAllEuiup() {
	self.player.GetModule("team").(*ModTeam).clearTeamPos()
}

func (self *ModGm) GmSuperHelp(index int) {
	switch index {
	case 1:
		self.GmSuper1()
	case 2:
		self.GmSuper2()
	case 3:
		self.GmSuper3()
	case 4:
		self.GmSuper4()
	case 5:
		self.GmSuper5()
	case 6:
		self.GmSuper6()
	case 7:
		self.GmSuper7()
	case 8:
		self.GmSuper8()
	case 9:
		self.GmSuper9()
	case 10:
		self.GmSuper10()
	case 11:
		self.GmSuper11()
	default:
		self.GmSuper1()
		self.GmSuper2()
		self.GmSuper3()
		self.GmSuper4()
		self.GmSuper5()
		self.GmSuper6()
		self.GmSuper7()
		self.GmSuper8()
		self.GmSuper9()
		self.GmSuper10()
		self.GmSuper11()
	}
}

//1.角色300级，英雄220级
func (self *ModGm) GmSuper1() {
	self.player.LvToMax()
	self.player.GetModule("hero").(*ModHero).HeroLvToMax()
}

//2.英雄们均穿戴最高阶装备，装备强化至最高阶
func (self *ModGm) GmSuper2() {
	self.player.GetModule("hero").(*ModHero).HeroEquipToBest()
}

//3.英雄们均升阶至最高品阶
func (self *ModGm) GmSuper3() {
	self.player.GetModule("hero").(*ModHero).HeroQuToMax()
}

//4.激活全部羁绊   csj
func (self *ModGm) GmSuper4() {
	self.player.GetModule("entanglement").(*ModEntanglement).GmAddAllFate()
	self.player.GetModule("entanglement").(*ModEntanglement).SendInfo([]byte{})
	self.player.GetModule("hero").(*ModHero).SendInfo()
}

//5.创建最高级公会  csj
func (self *ModGm) GmSuper5() {
	self.player.GetModule("union").(*ModUnion).GMAddUnionLevel()
	self.player.GetModule("union").(*ModUnion).GetUserUnionInfo()
}

//6.冒险关卡全部通关
func (self *ModGm) GmSuper6() {
	levelId := 112060
	self.player.GetModule("pass").(*ModPass).GMPassChapter(levelId)
	self.player.SendInfo("updateuserinfo")
}

//7.试炼之塔、四族高塔至最高层  csj
func (self *ModGm) GmSuper7() {
	for i := TOWER_TYPE_0; i < TOWER_TYPE_MAX; i++ {
		self.player.GetModule("tower").(*ModTower).GmSetToweLevel(i, 200)
	}
	self.player.GetModule("tower").(*ModTower).sendInfo()
}

//8.地牢通关至第三层
func (self *ModGm) GmSuper8() {
	self.player.GetModule("newpit").(*ModNewPit).GmSuperPass()
}

//9.完成所有日常，周常，主线任务   csj
func (self *ModGm) GmSuper9() {
	self.player.GetModule("task").(*ModTask).GMFinishAllTask()
	self.player.GetModule("task").(*ModTask).SendInfo()
}

//10.激活最高级贵族并领取特权
func (self *ModGm) GmSuper10() {
	self.player.VipToMax()
}

//11.大量金币、钻石、工会币、遣退币、地牢币
func (self *ModGm) GmSuper11() {
	lstItem := make([]PassItem, 0)
	lstItem = append(lstItem, PassItem{ITEM_GOLD, 100000000})
	lstItem = append(lstItem, PassItem{ITEM_GEM, 1000000})
	lstItem = append(lstItem, PassItem{ITEM_UNION, 100000})
	lstItem = append(lstItem, PassItem{ITEM_BACK_COIN, 100000})
	lstItem = append(lstItem, PassItem{ITEM_NEW_PIT_COIN, 100000})

	items := self.player.AddObjectPassItem(lstItem, "GM_SUPER", 0, 0, 0)
	self.player.GetModule("bag").(*ModBag).SendOnItem(items)

}

func (self *ModGm) GmShadowReset() {
	self.player.GetModule("instance").(*ModInstance).GmShadowReset()
}
