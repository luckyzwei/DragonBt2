package game
//
//import "encoding/json"
//
//const (
//	ELITEUPGRADE_END = 7	// 精英直升活动可领取的时间
//)
//
//// 精英直升
//type JS_EliteUpGrade struct {
//	BuyState     int               `json:"buystate"` 				//购买状态
//	EliteUpTask []JS_EliteUpTask   `json:"eliteupgradetask"`
//	//Plan 		int				   `json:"plan"`					//当前领取进度
//}
//
//type JS_EliteUpTask struct {
//	Id     int `json:"id"`     // id
//	State  int `json:"state"`  // 0未完成 1可领取  2已领取
//}
//
//// 检查精英直升数据
//func (self *ModRecharge) CheckEliteUpgrade()  {
//
//	if len(self.Sql_UserRecharge.eliteUpGrade.EliteUpTask) != ELITEUPGRADE_END {
//		var eliteUpGrade JS_EliteUpGrade
//		var eliteUpGradeTask JS_EliteUpTask
//		eliteUpGrade.BuyState = LOGIC_FALSE
//		for _,i := range GetCsvMgr().ActivitynewConfig {
//			if i.Type == ACT_ELITE_UPGRADE{
//				eliteUpGradeTask.Id = i.Id
//				eliteUpGradeTask.State = CANTFINISH
//				eliteUpGrade.EliteUpTask = append(eliteUpGrade.EliteUpTask,eliteUpGradeTask)
//			}
//		}
//		self.Sql_UserRecharge.eliteUpGrade = eliteUpGrade
//	}
//
//
//}
//
//// 领取精英直升奖励
//func (self *ModActivity) GetEliteUpgradeAward(body []byte) {
//
//	// 解析客户端消息
//	var msg C2S_GetEliteUpgradeAward
//	json.Unmarshal(body, &msg)
//
//	// 判断购买状态和活动时间
//	eliteupgradeaward := self.Sql_Activity.info[msg.Id]
//	//if eliteupgradeaward.Done == LOGIC_FALSE {
//	//	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_ELITEUPGRADE_NOT_OPEN"))		// 活动未开放
//	//	return
//	//}
//	//eliteupgradeaward := self.Sql_UserRecharge.eliteUpGrade
//	//if eliteupgradeaward.BuyState == LOGIC_FALSE {
//	//	self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_ELITEUPGRADE_NOT_OPEN"))		// 活动未开放
//	//	return
//	//}
//
//	var rep S2C_GetEliteUpgradeAward
//	rep.Cid = "geteliteupgradeaward"
//	actConfig := GetCsvMgr().ActivitynewConfig[msg.Id]
//	if actConfig == nil {
//		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_ELITEUPGRADE_CONFIG_ERROR"))
//		return
//	}
//	items,nums := actConfig.Items,actConfig.Nums
//	if len(items) != len(nums){
//		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_RECHARGE_ELITEUPGRADE_CONFIG_ERROR"))
//		return
//	}
//	var getItems []PassItem
//	for i := 0;i<len(items);i++{
//		if items[i] != 0{
//			getItems = append(getItems,PassItem{
//				items[i],nums[i],
//			})
//		}
//	}
//
//	if eliteupgradeaward.Done == CANTAKE && eliteupgradeaward.Progress >= actConfig.Step {
//		self.player.AddObjectLst(items,nums,"精英直升奖励领取",0,0,0)
//		// 奖励已领取
//		self.Sql_Activity.info[msg.Id].Done = TAKEN
//		//eliteupgradeaward.EliteUpTask[1].State = TAKEN
//		//eliteupgradeaward.Plan ++
//	}
//	rep.Items = getItems
//	self.player.SendMsg(rep.Cid, HF_JtoB(&rep))
//}
//
//
