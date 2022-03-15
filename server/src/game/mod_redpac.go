package game

import (
	"encoding/json"
	"fmt"
	"sort"
	//"time"
)

const (
	MaxRedPac      = 200
	RedPacTaken    = 1
	RedPacGone     = 2
	RedPacThank    = 3
	RedCampType    = 1
	RedUnionType   = 2
	RedPoolType    = 3
	RedActType     = 4
	MaxCampRedNum  = 200
	MaxUnionRedNum = 200
)

type ModRedPac struct {
	player *Player
	Data   San_RedPac //! 数据库结构
}

type RedRecord struct {
	ItemId  int `json:"itemid"`  //! 道具Id
	SendNum int `json:"sendnum"` //! 发送的道具数量
	GotNum  int `json:"gotnum"`  //! 抢到的红包数量
}

type RedWait struct {
	KeyId int64 `json:"keyid"` //! 待发送的红包唯一Id
	RedId int   `json:"redid"` //! 配置id
	Time  int64 `json:"time"`  //! 产生的时间
}

type RedRefresh struct {
	Id        int   `json:"id"`        //! 物品Id
	Num       int   `json:"num"`       //! 数量
	FreshTime int64 `json:"freshtime"` //! 刷新时间
}

type RedPacInfo struct {
	GemNum     int                   `json:"gemnum"`     //! 钻石上限
	WaitKey    int64                 `json:"waitkey"`    //! 待发送红包key
	SendLimit  map[int]*RedRefresh   `json:"limitnum"`   //! 发送上限
	GotLimit   map[int]*RedRefresh   `json:"gotnum"`     //! 接收上限
	TodayItems map[int]*RedRecord    `json:"todayitems"` //! 今天的红包历史
	TotalItems map[int]*RedRecord    `json:"totalitems"` //! 总的红包历史
	RedWait    []*RedWait            `json:"redwait"`    //! 待发送的系统
	UserRedPac map[int64]*UserRedPac `json:"userredpac"` //! 玩家的红包信息[已经领取或者领取完的红包]
}

type UserRedPac struct {
	KeyId     int64 `json:"keyid"`     //! 红包id
	Status    int   `json:"status"`    //! 红包状态, 0 未抢 1 已领取 2 被抢完 3 已答谢[只有领取的玩家才能答谢] 存放到全局
	Num       int   `json:"num"`       //! 领取的金额
	TimeStamp int64 `json:"timestamp"` //! 抢到的时间戳
}

type San_RedPac struct {
	Uid  int64
	Info string

	info RedPacInfo //! 每一期的次数
	DataUpdate
}

type RedPacShow struct {
	KeyId     int64  `json:"keyid"`     //! 红包唯一Id
	Uname     string `json:"uname"`     //! 姓名
	ItemId    int    `json:"itemid"`    //! 道具id
	ItemNum   int    `json:"itemnum"`   //! 道具数量
	TimeStamp int64  `json:"timestamp"` //! 时间戳
	Status    int    `json:"status"`    //! 状态  0 未抢 1 已领取 2 被抢完 3 已答谢
	Duration  int64  `json:"duration"`  //! 持续时间
	Uid       int64  `json:"uid"`       //! 玩家Id
}

type lstRedStatus []*RedStatus

func (s lstRedStatus) Len() int      { return len(s) }
func (s lstRedStatus) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstRedStatus) Less(i, j int) bool {
	if s[i].Num > s[j].Num {
		return true
	}

	if s[i].Num < s[j].Num {
		return false
	}

	if s[i].TimeStamp < s[j].TimeStamp {
		return true
	}

	if s[i].TimeStamp > s[j].TimeStamp {
		return false
	}

	return false
}

type lstRedPacShow []*RedPacShow

func (s lstRedPacShow) Len() int      { return len(s) }
func (s lstRedPacShow) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// 红包状态>红包持续时间>红包总钻石>玩家uid
func (s lstRedPacShow) Less(i, j int) bool {
	if s[i].Status < s[j].Status {
		return true
	}

	if s[i].Status > s[j].Status {
		return false
	}

	if s[i].Duration < s[j].Duration {
		return true
	}

	if s[i].Duration > s[j].Duration {
		return false
	}

	if s[i].ItemNum < s[j].ItemNum {
		return true
	}

	if s[i].ItemNum > s[j].ItemNum {
		return false
	}

	if s[i].Uid < s[j].Uid {
		return true
	}

	if s[i].Uid > s[j].Uid {
		return false
	}

	return false
}

func (self *RedPacInfo) init() {
	if self.SendLimit == nil {
		self.SendLimit = make(map[int]*RedRefresh)
	}

	if self.GotLimit == nil {
		self.GotLimit = make(map[int]*RedRefresh)
	}

	if self.TodayItems == nil {
		self.TodayItems = make(map[int]*RedRecord)
	}

	if self.TotalItems == nil {
		self.TotalItems = make(map[int]*RedRecord)
	}

	if self.RedWait == nil {
		self.RedWait = make([]*RedWait, 0)
	}

	if self.UserRedPac == nil {
		self.UserRedPac = make(map[int64]*UserRedPac)
	}
}

func (self *ModRedPac) Decode() {
	json.Unmarshal([]byte(self.Data.Info), &self.Data.info)
}

func (self *ModRedPac) Encode() {
	self.Data.Info = HF_JtoA(self.Data.info)
}

func (self *ModRedPac) OnGetData(player *Player) {
	self.player = player
	self.Data.info.init()
}

func (*ModRedPac) TableName() string {
	return "san_userredpac"
}

func (self *ModRedPac) OnGetOtherData() {
	redTableName := self.TableName()
	sql := fmt.Sprintf("select * from `%s` where uid = %d", redTableName, self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Data, redTableName, self.player.ID)
	if self.Data.Uid <= 0 {
		self.Data.Uid = self.player.ID
		self.Data.info.init()
		self.Encode()
		InsertTable(redTableName, &self.Data, 0, true)
	} else {
		self.Decode()
		self.Data.info.init()
		self.checkUserPac()
	}

	self.Data.Init(redTableName, &self.Data, true)
}

func (self *ModRedPac) OnSave(sql bool) {
	self.Encode()
	self.Data.Update(sql)
}

// 每天清零
func (self *ModRedPac) OnRefresh() {
	pInfo := &self.Data.info

	for key := range pInfo.TodayItems {
		pInfo.TodayItems[key] = &RedRecord{ItemId: key, SendNum: 0, GotNum: 0}
	}
}

// 消息处理
func (self *ModRedPac) OnMsg(ctrl string, body []byte) bool {
	switch ctrl {
	case "sendPool":
		self.sendPool(ctrl, body)
		return true
	case "getredpool":
		self.getRedPool(ctrl)
		return true
	case "gotredpac":
		self.gotRedPac(ctrl, body)
		return true
	case "lookredpac":
		self.lookRedpac(ctrl, body)
		return true
	case "thankredpac":
		self.thankRedPac(ctrl, body)
		return true
	case "sendglobalred":
		self.sendGlobalRed(ctrl, body)
		return true
	case "getredhis":
		self.getRedHis(ctrl, body)
		return true
	case "getredpacs":
		self.getRedPacs(ctrl, body)
		return true
	}

	return false
}

// 获取红包界面信息
func (self *ModRedPac) getRedPacs(cid string, body []byte) {
	if !self.isPassOk() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_PASS_THROUGH_THE_PLOT_LEVEL"))
		return
	}

	var msg C2S_GetRedPacs
	json.Unmarshal(body, &msg)
	res := make(map[int64]*RedPacShow)
	if msg.RedType == RedCampType {
		res = GetRedPacMgr().getCampRedPacShow(self.player.Sql_UserBase.Camp)
	} else if msg.RedType == RedUnionType {
		modUnion := self.player.GetModule("union").(*ModUnion)
		unionId := modUnion.Sql_UserUnionInfo.Unionid
		res = GetRedPacMgr().getUnionRedPacShow(unionId)
	}

	// 删除多余的红包状态并且纠正状态
	pInfo := &self.Data.info
	for _, v := range pInfo.UserRedPac {
		redShow, ok := res[v.KeyId]
		if ok {
			redShow.TimeStamp = v.TimeStamp
			redShow.Status = v.Status
		}
	}

	var handleRes []*RedPacShow
	for key := range res {
		handleRes = append(handleRes, res[key])
	}
	sort.Sort(lstRedPacShow(handleRes))

	var sendLimit []*RedRefresh
	var gotLimit []*RedRefresh
	for index := range GetCsvMgr().RedpacketConfig {
		pConfig := GetCsvMgr().RedpacketConfig[index]
		if pConfig == nil {
			continue
		}

		info, ok := pInfo.SendLimit[pConfig.Item]
		if !ok {
			sendLimit = append(sendLimit, &RedRefresh{
				Id:        pConfig.Item,
				Num:       pConfig.Ceiling,
				FreshTime: 0,
			})
		} else {
			// 检查发送
			info.checkSendNum()
			leftNum := pConfig.Ceiling - info.Num
			if leftNum <= 0 {
				leftNum = 0
			}
			sendLimit = append(sendLimit, &RedRefresh{
				Id:        pConfig.Item,
				Num:       leftNum,
				FreshTime: info.FreshTime,
			})
		}

		info2, ok := pInfo.GotLimit[pConfig.Item]
		if !ok {
			gotLimit = append(gotLimit, &RedRefresh{
				Id:        pConfig.Item,
				Num:       pConfig.Upperlimit,
				FreshTime: 0,
			})
		} else {
			// 检查接收
			info2.checkGotNum()
			leftNum := pConfig.Upperlimit - info2.Num
			if leftNum <= 0 {
				leftNum = 0
			}
			gotLimit = append(gotLimit, &RedRefresh{
				Id:        pConfig.Item,
				Num:       leftNum,
				FreshTime: info2.FreshTime,
			})
		}

	}

	msgSend := &S2C_GetRedPacs{
		Cid:        cid,
		RedPacShow: handleRes,
		GemNum:     self.Data.info.GemNum,
		SendLimit:  sendLimit,
		GotLimit:   gotLimit,
	}
	self.player.SendMsg(msgSend.Cid, HF_JtoB(msgSend))

}

func (self *ModRedPac) getItems() []int {
	config := GetCsvMgr().RedpacketConfig
	var res []int
	for index := range config {
		res = append(res, config[index].Item)
	}
	return res
}

func (self *ModRedPac) NewRedRecord(itemId int, sendNum int, GotNum int) *RedRecord {
	return &RedRecord{ItemId: itemId, SendNum: sendNum, GotNum: GotNum}
}

func (self *ModRedPac) getRedHis(cid string, body []byte) {
	if !self.isPassOk() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_PASS_THROUGH_THE_PLOT_LEVEL"))
		return
	}

	items := self.getItems()
	var today []*RedRecord
	var total []*RedRecord
	pInfo := &self.Data.info
	for index := range items {
		itemId := items[index]
		r1, ok := pInfo.TodayItems[itemId]
		if ok {
			today = append(today, r1)
		} else {
			today = append(today, self.NewRedRecord(itemId, 0, 0))
		}

		r2, ok := pInfo.TotalItems[itemId]
		if ok {
			total = append(total, r2)
		} else {
			total = append(total, self.NewRedRecord(itemId, 0, 0))
		}
	}

	msg := &S2C_GetRedHis{
		Cid:   cid,
		Today: today,
		Total: total,
	}
	self.player.SendMsg(msg.Cid, HF_JtoB(msg))
}

func (self *ModRedPac) isPeopleOk(num int, people int) bool {
	return false
}

// 发送军团红包, 区分钻石和其他物品
func (self *ModRedPac) sendGlobalRed(cid string, body []byte) {
	if !self.isPassOk() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_PASS_THROUGH_THE_PLOT_LEVEL"))
		return
	}

	var msg C2S_SendGlobalRed
	json.Unmarshal(body, &msg)

	redId := msg.RedId
	pConfig := GetCsvMgr().getRedPacConfig(redId)
	if pConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_RED_PACKET_CONFIGURATION_DOES_NOT"))
		return
	}

	if msg.RedType != RedCampType && msg.RedType != RedUnionType {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_RED_PACKET_TYPE_ERROR"))
		return
	}

	selectIndex := msg.SelectIndex
	if selectIndex < 1 || selectIndex > 3 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_ERROR_IN_GEAR_SELECTION"))
		return
	}

	index := selectIndex - 1

	if index < 0 || index >= len(pConfig.Vips) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_MISALLOCATION_OF_RED_ENVELOPE_ARISTOCRACY"))
		return
	}

	if index < 0 || index >= len(pConfig.Prices) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_MISALLOCATION_OF_RED_ENVELOPE_PRICE"))
		return
	}

	if index < 0 || index >= len(pConfig.Costs) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_RED_PACKET_CONSUMPTION_CONFIGURATION_ERROR"))
		return
	}

	if index < 0 || index >= len(pConfig.Peoples) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_INCORRECT_ALLOCATION_OF_RED_ENVELOPE"))
		return
	}

	vipNeed := pConfig.Vips[index]
	if self.player.Sql_UserBase.Vip < vipNeed {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_PLAYERS_NOBILITY_LEVEL_IS_INSUFFICIENT"))
		return
	}

	itemId := pConfig.Item
	price := pConfig.Prices[index]
	cost := pConfig.Costs[index]
	pInfo := &self.Data.info

	people := 0
	resUnionId := 0
	if msg.RedType == RedUnionType {
		modUnion := self.player.GetModule("union").(*ModUnion)
		unionId := modUnion.Sql_UserUnionInfo.Unionid
		unionNum := GetUnionMgr().GetUnionNum(unionId)

		if unionNum == 0 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_AT_PRESENT_THERE_IS_NO"))
			return
		}

		// 检查惩罚时间
		if TimeServer().Unix()-modUnion.Sql_UserUnionInfo.LastUpdTime < 86400 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_IN_LEGION_PUNISHMENT_TIME_UNABLE"))
			return
		}

		// 判断是否是军团长或者副军团长
		pos := modUnion.Sql_UserUnionInfo.Position
		if pos != 1 && pos != 2 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_LEGION_RED_PACK_CAN_ONLY"))
			return
		}

		data := GetUnionMgr().GetUnion(unionId)
		if data != nil {
			csv_community := GetCsvMgr().CommunityConfig[data.Level]
			people = csv_community.Membernum
		} else {
			people = unionNum
		}
		resUnionId = unionId
	}

	leftItemNum := 0
	if itemId == DEFAULT_GEM { // 钻石
		if pInfo.GemNum < price {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_INSUFFICIENT_AVAILABLE_BALANCE"))
			return
		}

		if self.player.Sql_UserBase.Gem < cost {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_UNION_DIAMOND_SHORTAGE"))
			return
		}

		leftItemNum = pInfo.GemNum - cost

	} else { // 从配置读取发送上限
		pLimit, ok := pInfo.SendLimit[itemId]
		if !ok {
			pLimit = self.NewSendRedFresh(itemId)
			pInfo.SendLimit[itemId] = pLimit
		}

		if pConfig.Ceiling > 0 {
			if pLimit.Num+price > pConfig.Ceiling {
				self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_INSUFFICIENT_AVAILABLE_BALANCE"))
				return
			}
		}
		leftItemNum = pConfig.Ceiling - pLimit.Num - price
		// 检查背包物品数量
		if self.player.GetObjectNum(DEFAULT_GEM) < cost {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_UNION_DIAMOND_SHORTAGE"))
			return
		}
	}

	camp := self.player.Sql_UserBase.Camp
	iconId := self.player.Sql_UserBase.IconId

	self.addSendNum(pConfig.Item, price)
	redPacParam := &RedPacParam{
		Uid:      self.player.Sql_UserBase.Uid,
		UnionId:  resUnionId,
		Camp:     camp,
		RedType:  msg.RedType,
		UName:    self.player.Sql_UserBase.UName,
		People:   people,
		TotalNum: price,
		ItemId:   pConfig.Item,
		Msg:      msg.Msg,
		IconId:   iconId,
	}

	redPac := GetRedPacMgr().CreateRedPac(redPacParam)

	// 扣除对应的信息
	if itemId == DEFAULT_GEM { // 钻石
		pInfo.GemNum -= cost
		self.player.AddObject(DEFAULT_GEM, -cost, 0, 0, 0, "发送红包")
	} else {
		self.player.AddObject(DEFAULT_GEM, -cost, 0, 0, 0, "发送红包")

		// 更新Limit
		pLimit, ok := pInfo.SendLimit[itemId]
		if !ok {
			pLimit = self.NewSendRedFresh(itemId)
			pInfo.GotLimit[itemId] = pLimit
		}
		pLimit.checkSendNum()
		pLimit.updateSendNum(price)
	}

	// 通知玩家
	msgSend := &S2C_SendUnionRed{
		Cid:        cid,
		GemNum:     self.player.Sql_UserBase.Gem,
		LeftNum:    leftItemNum,
		LeftItemId: pConfig.Item,
	}
	self.player.SendMsg(msgSend.Cid, HF_JtoB(msgSend))
	// 根据类型播报, union or camp
	if msg.RedType == RedUnionType {
		GetRedPacMgr().checkUnionRedPac(resUnionId)
		GetRedPacMgr().addUnionRed(resUnionId, redPac.RedId)
		party := GetUnionMgr().GetUnion(resUnionId)
		if party != nil {
			GetRedPacMgr().Send2Union(party, redPac.RedId)
			//GetRedPacMgr().sendUnionChat(party, camp, self.player.Sql_UserBase.UName, pConfig.Item)
		} else {
			LogError("party is nil!")
		}
	} else if msg.RedType == RedCampType {
		GetRedPacMgr().checkCampRedPac(camp)
		GetRedPacMgr().addCampRed(camp, redPac.RedId)
		GetRedPacMgr().Send2Camp(camp, redPac.RedId)
		//GetRedPacMgr().sendSystemChat(self.player.Sql_UserBase.UName, camp, pConfig.Item)
	}
}

// 答谢红包
func (self *ModRedPac) thankRedPac(cid string, body []byte) {
	if !self.isPassOk() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_PASS_THROUGH_THE_PLOT_LEVEL"))
		return
	}

	var msg C2S_ThankRedPac
	json.Unmarshal(body, &msg)

	uid := self.player.Sql_UserBase.Uid
	// 检查玩家是否领取过这个红包
	keyId := msg.KeyId
	userRedPac := self.getUserRedPac(keyId)
	if userRedPac == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_RED_ENVELOPE_DOES_NOT_EXIST"))
		return
	}

	if userRedPac.Status == RedPacThank {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_IVE_ALREADY_THANKED_THE_RED"))
		return
	}

	if userRedPac.Num <= 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_AT_PRESENT_RED_ENVELOPES_HAVE"))
		return
	}

	if userRedPac.Status != RedPacTaken {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_THE_PLAYER_DID_NOT_RECEIVE"))
		return
	}

	redPac, _, _ := GetRedPacMgr().GetRedPac(uid, msg.KeyId)
	if redPac == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_RED_ENVELOPE_DOES_NOT_EXIST"))
		return
	}

	userRedPac.Status = RedPacThank
	msgSend := &S2C_ThankRedPac{
		Cid:        cid,
		UserRedPac: userRedPac,
	}
	self.player.SendMsg(msgSend.Cid, HF_JtoB(msgSend))
}

// 查看红包状态
func (self *ModRedPac) lookRedpac(cid string, body []byte) {
	if !self.isPassOk() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_PASS_THROUGH_THE_PLOT_LEVEL"))
		return
	}

	var msg C2S_LookRedPac
	json.Unmarshal(body, &msg)

	uid := self.player.Sql_UserBase.Uid
	redPac, _, _ := GetRedPacMgr().GetRedPac(uid, msg.KeyId)
	if redPac == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_RED_ENVELOPE_DOES_NOT_EXIST"))
		return
	}

	userRedPac := self.getUserRedPac(redPac.RedId)
	total := 0
	if userRedPac == nil {
		total = 0
	} else {
		total = userRedPac.Num
	}

	takenPacNum := len(redPac.Status)

	takenNum := 0
	for _, v := range redPac.Status {
		takenNum += v.Num
	}

	redStatusLst := make([]*RedStatus, 0, len(redPac.Status))
	for key := range redPac.Status {
		redStatusLst = append(redStatusLst, redPac.Status[key])
	}

	sort.Sort(lstRedStatus(redStatusLst))

	redInfo := make([]*RedInfo, 0, len(redStatusLst))
	for index := range redStatusLst {
		if redStatusLst[index] == nil {
			continue
		}
		redInfo = append(redInfo, &RedInfo{
			Name: redStatusLst[index].Uname,
			Num:  redStatusLst[index].Num,
		})
	}

	thank := 0
	if userRedPac != nil && userRedPac.Status == 3 {
		thank = 1
	}
	msgSend := &S2C_LookRedPac{
		Cid:         cid,
		IconId:      redPac.Iconid,
		Name:        redPac.UName,
		Msg:         redPac.Msg,
		Num:         total,
		ItemId:      redPac.ItemId,
		TakenPacNum: takenPacNum,
		AllPacNum:   redPac.People,
		TakenNum:    takenNum,
		AllNum:      redPac.TotalNum,
		RedInfo:     redInfo,
		Uid:         redPac.Uid,
		Thank:       thank,
	}
	self.player.SendMsg(msgSend.Cid, HF_JtoB(msgSend))
}

func (self *ModRedPac) getUserRedPac(keyId int64) *UserRedPac {
	pInfo := &self.Data.info
	v, ok := pInfo.UserRedPac[keyId]
	if ok {
		return v
	}

	return nil
}

// 抢红包, 抢到, 没抢到
func (self *ModRedPac) gotRedPac(cid string, body []byte) {
	if !self.isPassOk() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_PASS_THROUGH_THE_PLOT_LEVEL"))
		return
	}

	var msg C2S_GotRedPac
	json.Unmarshal(body, &msg)

	// 检查是否已经有这个红包
	pInfo := &self.Data.info
	keyId := msg.KeyId
	if v, ok := pInfo.UserRedPac[keyId]; ok {
		if v.Status == 1 || v.Status == 3 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_DONT_RECEIVE_IT_AGAIN"))
			return
		} else if v.Status == 2 {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_ALREADY_ROBBED"))
			return
		}
	}

	uid := self.player.Sql_UserBase.Uid
	uname := self.player.Sql_UserBase.UName

	redPac, isTaken, allTaken := GetRedPacMgr().GetRedPac(uid, msg.KeyId)
	if redPac == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_RED_ENVELOPE_DOES_NOT_EXIST"))
		return
	}

	if redPac.EndTime <= TimeServer().Unix() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_THE_RED_ENVELOPE_IS_OVER"))
		return
	}

	if isTaken {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_DONT_RECEIVE_IT_AGAIN"))
		return
	}

	pConfig := GetCsvMgr().getRedPacConfByItemId(redPac.ItemId)
	if pConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_RED_PACKET_CONFIGURATION_DOES_NOT"))
		return
	}

	// 检查是否达到抢的上限

	pLimit, ok := pInfo.GotLimit[redPac.ItemId]
	if !ok {
		pLimit = self.NewGotRedFresh(redPac.ItemId)
		pInfo.GotLimit[redPac.ItemId] = pLimit
	}

	// 检查是否需要刷新
	pLimit.checkGotNum()

	if pConfig.Upperlimit > 0 {
		if pLimit.Num >= pConfig.Upperlimit {
			itemName := GetCsvMgr().GetItemName(pConfig.Item)
			self.player.SendErrInfo("err", fmt.Sprintf(GetCsvMgr().GetText("STR_REDPAC_LIMIT_RAPE"), itemName))
			self.player.SendMsg("redpacode", HF_JtoB(&S2C_RedPacCode{
				Cid:  "redpacode",
				Code: 1,
			}))
			return
		}
	}

	// 通过红包类型判断是否可以领取
	camp := self.player.Sql_UserBase.Camp
	if redPac.RedType == RedCampType {
		if camp != redPac.Camp {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_DIFFERENT_COUNTRIES_CANT_GET_RED"))
			return
		}
	} else if redPac.RedType == RedUnionType || redPac.RedType == RedPoolType {
		modUnion := self.player.GetModule("union").(*ModUnion)
		unionId := modUnion.Sql_UserUnionInfo.Unionid
		if unionId != redPac.UnionId {
			self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_LEGIONS_CANT_GET_RED_ENVELOPES"))
			return
		}
	}

	if allTaken {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_ITS_TOO_LATE_THE_RED"))
		self.handleEmptyRed(cid, redPac.RedId)
		return
	}

	pRedStatus := GetRedPacMgr().AllocateRedPac(keyId, uid, uname)
	if pRedStatus == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_ITS_TOO_LATE_THE_RED"))
		self.handleEmptyRed(cid, redPac.RedId)
		return
	}
	pLimit.updateGotNum(pRedStatus.Num)
	userRedPac := self.createRedPac(keyId, pRedStatus)

	pInfo.UserRedPac[keyId] = userRedPac
	self.player.AddObject(pRedStatus.ItemId, pRedStatus.Num, 0, 0, 0, "抢红包")
	self.addGotNum(pRedStatus.ItemId, pRedStatus.Num)

	// 更新玩家红包信息
	self.makeGotSyn(cid, userRedPac, pConfig.Item)
}

func (self *ModRedPac) handleEmptyRed(cid string, keyId int64) {
	userRedPac := self.createEmptyRedPac(keyId)
	pInfo := &self.Data.info
	if _, ok := pInfo.UserRedPac[keyId]; ok {
		return
	}
	pInfo.UserRedPac[keyId] = userRedPac
	//self.makeGotSyn(cid, userRedPac)
}

func (self *ModRedPac) makeGotSyn(cid string, userRedPac *UserRedPac, itemId int) {
	msgSend := &S2C_GotRedPac{
		Cid:        cid,
		UserRedPac: userRedPac,
		ItemId:     itemId,
	}
	self.player.SendMsg(msgSend.Cid, HF_JtoB(msgSend))
}

func (self *ModRedPac) createRedPac(keyId int64, pRedStatus *RedStatus) *UserRedPac {
	return &UserRedPac{
		KeyId:     keyId,
		Status:    pRedStatus.Status,
		Num:       pRedStatus.Num,
		TimeStamp: pRedStatus.TimeStamp,
	}
}

func (self *ModRedPac) createEmptyRedPac(keyId int64) *UserRedPac {
	now := TimeServer().Unix()
	return &UserRedPac{
		KeyId:     keyId,
		Status:    RedPacGone,
		Num:       0,
		TimeStamp: now,
	}
}

// 发送红包池的红包
func (self *ModRedPac) sendPool(cid string, body []byte) {
	if !self.isPassOk() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_PASS_THROUGH_THE_PLOT_LEVEL"))
		return
	}

	var msg C2S_SendPool
	json.Unmarshal(body, &msg)

	keyId := msg.KeyId
	redWait, index := self.findRedWait(keyId)
	if redWait == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_CURRENT_RED_ENVELOPE_DOES_NOT"))
		return
	}

	modUnion := self.player.GetModule("union").(*ModUnion)
	unionId := modUnion.Sql_UserUnionInfo.Unionid
	camp := self.player.Sql_UserBase.Camp
	iconId := self.player.Sql_UserBase.IconId
	unionNum := GetUnionMgr().GetUnionNum(unionId)
	if unionNum == 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_AT_PRESENT_THERE_IS_NO"))
		return
	}

	pConfig := self.getConfigById(redWait.RedId)
	if pConfig == nil {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_RED_PACKET_CONFIGURATION_DOES_NOT"))
		return
	}

	if TimeServer().Unix()-modUnion.Sql_UserUnionInfo.LastUpdTime < 86400 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_IN_LEGION_PUNISHMENT_TIME_UNABLE"))
		return
	}

	redPacParam := &RedPacParam{
		Uid:      self.player.Sql_UserBase.Uid,
		UnionId:  unionId,
		Camp:     camp,
		RedType:  RedPoolType,
		UName:    self.player.Sql_UserBase.UName,
		People:   pConfig.People,
		TotalNum: pConfig.Num,
		ItemId:   pConfig.Item,
		Msg:      msg.Msg,
		IconId:   iconId,
	}
	redPac := GetRedPacMgr().CreateRedPac(redPacParam)
	self.removeRedWait(index)

	self.addSendNum(pConfig.Item, pConfig.Num)
	// 通知玩家
	msgSend := &S2C_SendPool{
		Cid:     cid,
		RedWait: self.Data.info.RedWait,
	}
	self.player.SendMsg(msgSend.Cid, HF_JtoB(msgSend))
	GetRedPacMgr().checkUnionRedPac(unionId)
	GetRedPacMgr().addUnionRed(unionId, redPac.RedId)

	// 全服播报
	party := GetUnionMgr().GetUnion(self.player.GetModule("union").(*ModUnion).Sql_UserUnionInfo.Unionid)
	if party != nil {
		GetRedPacMgr().Send2Union(party, redPac.RedId)
		//GetRedPacMgr().sendUnionChat(party, camp, self.player.Sql_UserBase.UName, pConfig.Item)
	} else {
		LogError("party is nil!")
	}

}

func (self *ModRedPac) getConfigByNum(addGem int) *RedpacketmoneyConfig {
	config := GetCsvMgr().RedpacketmoneyConfig
	var pConfig *RedpacketmoneyConfig
	for index := range config {
		if config[index].Money == addGem {
			pConfig = config[index]
			break
		}
	}

	return pConfig
}

func (self *ModRedPac) getConfigById(redId int) *RedpacketmoneyConfig {
	return GetCsvMgr().getConfigById(redId)
}

// 创建红包池
func (self *ModRedPac) CreateRedWait(addGem int) {
	self.Data.info.GemNum += addGem
	pConfig := self.getConfigByNum(addGem)
	if pConfig != nil {
		self.addRedWait(pConfig)
	}
	self.synGemNum()
}

// 创建红包
func (self *ModRedPac) addRedWait(config *RedpacketmoneyConfig) {
	pInfo := &self.Data.info
	pInfo.WaitKey += 1
	redWait := &RedWait{
		KeyId: pInfo.WaitKey,
		RedId: config.Id,
		Time:  TimeServer().Unix(),
	}

	if len(pInfo.RedWait) >= MaxRedPac {
		pInfo.RedWait = append(pInfo.RedWait[:0], pInfo.RedWait[1:]...)
	}

	pInfo.RedWait = append([]*RedWait{redWait}, pInfo.RedWait...)
}

// 获取红包信息
func (self *ModRedPac) getRedPool(cid string) {
	if !self.isPassOk() {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_REDPAC_PASS_THROUGH_THE_PLOT_LEVEL"))
		return
	}

	msg := &S2C_GetRedPac{
		Cid:     cid,
		RedWait: self.Data.info.RedWait,
	}
	self.player.SendMsg(msg.Cid, HF_JtoB(msg))
}

func (self *ModRedPac) findRedWait(keyId int64) (*RedWait, int) {
	redWait := self.Data.info.RedWait
	for index := range redWait {
		if redWait[index].KeyId == keyId {
			return redWait[index], index
		}
	}
	return nil, 0
}

// 删除红包
func (self *ModRedPac) removeRedWait(index int) {
	pInfo := &self.Data.info
	if index < 0 || index >= len(pInfo.RedWait) {
		return
	}

	pInfo.RedWait = append(pInfo.RedWait[:index], pInfo.RedWait[index+1:]...)
}

func (self *ModRedPac) isPassOk() bool {
	return GetCsvMgr().IsLevelOpen2(self.player, 54)
}

func (self *ModRedPac) addSendNum(itemid int, num int) {
	pInfo := &self.Data.info

	if _, ok := pInfo.TodayItems[itemid]; ok {
		pInfo.TodayItems[itemid].SendNum += num
	} else {
		pInfo.TodayItems[itemid] = self.NewRedRecord(itemid, num, 0)
	}

	if _, ok := pInfo.TotalItems[itemid]; ok {
		pInfo.TotalItems[itemid].SendNum += num
	} else {
		pInfo.TotalItems[itemid] = self.NewRedRecord(itemid, num, 0)
	}
}

func (self *ModRedPac) addGotNum(itemid int, num int) {
	pInfo := &self.Data.info

	if _, ok := pInfo.TodayItems[itemid]; ok {
		pInfo.TodayItems[itemid].GotNum += num
	} else {
		pInfo.TodayItems[itemid] = self.NewRedRecord(itemid, 0, num)
	}

	if _, ok := pInfo.TotalItems[itemid]; ok {
		pInfo.TotalItems[itemid].GotNum += num
	} else {
		pInfo.TotalItems[itemid] = self.NewRedRecord(itemid, 0, num)
	}
}

// 创建一个刷新纪录
func (self *ModRedPac) NewSendRedFresh(itemId int) *RedRefresh {
	return &RedRefresh{
		Id:        itemId,
		Num:       0,
		FreshTime: GetCsvMgr().getSendFreshTime(itemId),
	}
}

func (self *ModRedPac) NewGotRedFresh(itemId int) *RedRefresh {
	return &RedRefresh{
		Id:        itemId,
		Num:       0,
		FreshTime: GetCsvMgr().getGotFreshTime(itemId),
	}
}

// 更新发送和抢红包时自动刷新
func (self *RedRefresh) updateSendNum(add int) {
	self.Num += add
}

func (self *RedRefresh) checkSendNum() {
	now := TimeServer()
	if self.FreshTime > 0 && self.FreshTime <= now.Unix() {
		self.Num = 0
		self.FreshTime = GetCsvMgr().getSendFreshTime(self.Id)
	}

	if self.FreshTime >= now.Unix()+86400 {
		self.FreshTime = GetCsvMgr().getSendFreshTime(self.Id)
	}

	if self.FreshTime == 0 {
		self.FreshTime = GetCsvMgr().getSendFreshTime(self.Id)
	}
}

func (self *RedRefresh) updateGotNum(add int) {
	self.Num += add
}

func (self *RedRefresh) checkGotNum() {
	now := TimeServer()
	if self.FreshTime > 0 && self.FreshTime <= now.Unix() {
		self.Num = 0
		self.FreshTime = GetCsvMgr().getGotFreshTime(self.Id)
	}

	if self.FreshTime >= now.Unix()+86400 {
		self.FreshTime = GetCsvMgr().getGotFreshTime(self.Id)
	}

	if self.FreshTime == 0 {
		self.FreshTime = GetCsvMgr().getGotFreshTime(self.Id)
	}
}

func (self *ModRedPac) checkUserPac() {
	var keyIds []int64
	pInfo := &self.Data.info
	for _, info := range pInfo.UserRedPac {
		keyIds = append(keyIds, info.KeyId)
	}

	if len(keyIds) <= 0 {
		return
	}
	removeIds := GetRedPacMgr().CheckRedPacOk(keyIds)
	for key := range removeIds {
		delete(pInfo.UserRedPac, key)
	}
}

func (self *ModRedPac) synGemNum() {
	self.player.SendMsg("synredgemnum", HF_JtoB(&S2C_SynRedGemNum{
		Cid:    "synredgemnum",
		GemNum: self.Data.info.GemNum,
	}))
}
