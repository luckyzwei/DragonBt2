package game

import (
	"encoding/json"
	"fmt"
	//"time"
)

const (
	SetHead     = 1
	SetPortrait = 2
)

// 头像以及配件
type ModHead struct {
	player *Player
	Data   SanHead //! 数据库数据
}

// 拥有的头像ID
type SanHead struct {
	Uid  int64
	Info string
	info map[int]*HeadInfo
	DataUpdate
}

type HeadInfo struct {
	Id       int   `json:"id"`        // 头像/挂件Id, 对应配置表里面的id
	EndTime  int64 `json:"end_time"`  // 结束时间
	TimeType int   `json:"time_type"` // 时间类型
}

func (m *ModHead) Decode() {
	err := json.Unmarshal([]byte(m.Data.Info), &m.Data.info)
	if err != nil {
		LogError(err.Error())
	}
}

func (m *ModHead) Encode() {
	m.Data.Info = HF_JtoA(m.Data.info)
}

func (m *ModHead) getTableName() string {
	return "san_userhead"
}

func (m *ModHead) init(uid int64) {
	m.Data.Uid = uid
	m.CheckInfo()
}

// 初始化所有道具信息
func (m *ModHead) CheckInfo() {
	if m.Data.info == nil {
		m.Data.info = make(map[int]*HeadInfo)
	}
	m.InitHead(DEFAULT_HEAD_ICON)
}

func (m *ModHead) OnGetData(player *Player) {
	m.player = player
}

func (m *ModHead) OnGetOtherData() {
	tableName := m.getTableName()
	sql := fmt.Sprintf("select * from `%s` where uid = %d", tableName, m.player.ID)
	GetServer().DBUser.GetOneData(sql, &m.Data, tableName, m.player.ID)
	if m.Data.Uid <= 0 {
		m.init(m.player.ID)
		m.CheckInfo()
		m.Encode()
		InsertTable(tableName, &m.Data, 0, true)
	} else {
		m.Decode()
		m.CheckInfo()
	}

	m.Data.Init(tableName, &m.Data, true)
}

func (m *ModHead) OnSave(sql bool) {
	m.Encode()
	m.Data.Update(sql)
}

func (m *ModHead) OnRefresh() {

}

// 老的消息处理
func (m *ModHead) OnMsg(ctrl string, body []byte) bool {
	return false
}

// 注册消息
func (m *ModHead) onReg(handlers map[string]func(body []byte)) {
	handlers["head_action"] = m.onHeadAction

}

type C2S_HeadAction struct {
	Id     int `json:"id"`
	Action int `json:"action"`
}

type S2C_HeadAction struct {
	Cid      string `json:"cid"`
	Action   int    `json:"action"`
	IconId   int    `json:"iconid"`   //玩家头像
	Portrait int    `json:"portrait"` // 玩家头像框
}

// 设置头像信息
func (m *ModHead) onHeadAction(body []byte) {
	var msg C2S_HeadAction
	err := json.Unmarshal(body, &msg)
	if err != nil {
		LogError(err.Error())
	}

	if msg.Action == SetHead {
		m.setHead(&msg)
	} else if msg.Action == SetPortrait {
		m.SetPortrait(&msg)
	}
	m.player.NoticeCenterBaseInfo()
	GetArenaMgr().Rehead(m.player)
	GetArenaSpecialMgr().Rehead(m.player)
	GetTopArenaMgr().Rehead(m.player)
}

// 同步头像信息
func (m *ModHead) SynHeadAction(action int) {
	msg := &S2C_HeadAction{}
	msg.Cid = "head_action"
	msg.Action = action
	msg.IconId = m.player.Sql_UserBase.IconId
	msg.Portrait = m.player.Sql_UserBase.Portrait
	m.player.Send(msg.Cid, msg)
}

func (m *ModHead) setHead(msg *C2S_HeadAction) {
	id := msg.Id
	config, ok := GetCsvMgr().HeadConfigMap[id]
	if !ok {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_HEAD_NO_AVATAR_CONFIGURATION_EXISTS"))
		return
	}

	if config.Type != 1 {
		return
	}

	if m.Data.info == nil {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_HEAD_DATA_IS_NOT_INITIALIZED"))
		return
	}

	head, ok := m.Data.info[id]
	if !ok {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_HEAD_THE_AVATAR_CONFIGURATION_IS_NOT"))
		return
	}

	if head == nil {
		m.player.SendErr(fmt.Sprintf("head == nil, id:%d", id))
		return
	}

	if head.TimeType == 2 && head.EndTime < TimeServer().Unix() {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_HEAD_HEADFRAME_HAS_EXPIRED"))
		return
	}

	oldhead := m.player.Sql_UserBase.IconId
	m.player.Sql_UserBase.IconId = id
	if m.player.GetUnionId() != 0 {
		GetUnionMgr().FreshMember(m.player.GetUnionId(), &m.player.Sql_UserBase, &m.player.GetModule("union").(*ModUnion).Sql_UserUnionInfo)
	}

	// 还需要设置军团中的头像框
	m.SynHeadAction(msg.Action)
	GetOfflineInfoMgr().ReIconId(m.player)

	GetServer().SqlLog(m.player.Sql_UserBase.Uid, LOG_PLAYER_CHANGE_ICON, oldhead, id, 0, "用户改头像", 0, 0, m.player)
}

func (m *ModHead) SetPortrait(msg *C2S_HeadAction) {
	id := msg.Id
	config, ok := GetCsvMgr().HeadConfigMap[id]
	if !ok {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_HEAD_NO_AVATAR_CONFIGURATION_EXISTS"))
		return
	}

	if config.Type != 2 {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_HEAD_INCORRECT_HEADFRAME_TYPE"))
		return
	}

	if m.Data.info == nil {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_HEAD_DATA_IS_NOT_INITIALIZED"))
		return
	}

	head, ok := m.Data.info[id]
	if !ok {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_HEAD_HEADFRAME_CONFIGURATION_NOT_YET_AVAILABLE"))
		return
	}

	if head.TimeType == 2 && head.EndTime < TimeServer().Unix() {
		m.player.SendErr(GetCsvMgr().GetText("STR_MOD_HEAD_HEADFRAME_HAS_EXPIRED"))
		return
	}

	oldportrait := m.player.Sql_UserBase.Portrait
	m.player.Sql_UserBase.Portrait = id
	if m.player.GetUnionId() != 0 {
		GetUnionMgr().FreshMember(m.player.GetUnionId(), &m.player.Sql_UserBase, &m.player.GetModule("union").(*ModUnion).Sql_UserUnionInfo)
	}
	m.SynHeadAction(msg.Action)
	GetOfflineInfoMgr().RePortrait(m.player)

	GetServer().SqlLog(m.player.Sql_UserBase.Uid, LOG_PLAYER_CHANGE_PORTRAIT, oldportrait, id, 0, "用户改边框", 0, 0, m.player)
}

type S2C_HeadInfo struct {
	Cid      string            `json:"cid"`       // 消息cid
	HeadInfo map[int]*HeadInfo `json:"head_info"` // headportrait 配置ID
}

// 拥有的头像
func (m *ModHead) SendInfo() {
	m.check()
	const cid = "headinfo"
	msg := &S2C_HeadInfo{}
	msg.Cid = cid
	msg.HeadInfo = m.Data.info
	m.player.Send(cid, msg)
}

func (m *ModHead) check() {
	if m.Data.info == nil {
		return
	}
	for k, v := range m.Data.info {
		if v == nil {
			LogError("v is nil, check head, playerId:", m.player.Sql_UserBase.Uid)
			continue
		}
		if v.Id == 0 {
			v.Id = k
		}
	}
}

// 检查召唤英雄是否获得头像
func (m *ModHead) CheckHero(heroId int) {
	for _, v := range GetCsvMgr().HeadConfigMap {
		if v.Open == 1 && v.Condition == heroId {
			m.AddHead(v)
		}
	}
}

func (m *ModHead) AddHead(config *HeadConfig) bool {
	if m.Data.info == nil {
		LogError("m.Data.info is nil")
		return false
	}

	v, ok := m.Data.info[config.Id]
	if ok {
		// cd 增加逻辑
		if config.Timetype == 1 {
			v.Update(config.Timevalue)
			m.SynHeadInfo(v)

			return true
		}
	} else {
		pInfo := NewHeadInfo(config)
		m.Data.info[config.Id] = pInfo
		m.SynHeadInfo(pInfo)

		return true
	}

	return false
}

func NewHeadInfo(config *HeadConfig) *HeadInfo {
	// 时间类型
	headInfo := &HeadInfo{}
	headInfo.Id = config.Id
	headInfo.TimeType = config.Timetype
	if config.Timetype == 0 {
		headInfo.EndTime = 0
	} else if config.Timetype == 1 {
		if headInfo.EndTime <= TimeServer().Unix() {
			headInfo.EndTime = TimeServer().Unix()
		}
		headInfo.EndTime += int64(config.Timevalue)
	}

	return headInfo
}

func (m *HeadInfo) Update(timeValue int) {
	if m.EndTime < TimeServer().Unix() {
		m.EndTime = TimeServer().Unix()
	}
	m.EndTime += int64(timeValue)

}

// 检查使用道具
func (m *ModHead) CheckUseItem(itemId int) bool {
	for _, v := range GetCsvMgr().HeadConfigMap {
		if v.Item == itemId {
			return m.AddHead(v)
		}
	}

	return false
}

type S2C_SynHeadInfo struct {
	Cid      string    `json:"cid"`       // 消息cid
	HeadInfo *HeadInfo `json:"head_info"` // headportrait 配置ID
}

// 同步头像信息
func (m *ModHead) SynHeadInfo(headInfo *HeadInfo) {
	msg := &S2C_SynHeadInfo{}
	msg.Cid = "syn_head_info"
	msg.HeadInfo = headInfo
	if headInfo != nil {
		m.player.Send(msg.Cid, msg)
	}
}

func (m *ModHead) AddHead2(id int) {
	config, ok := GetCsvMgr().HeadConfigMap[id]
	if !ok {
		LogError("config is nil")
		return
	}

	if m.Data.info == nil {
		LogError("m.Data.info is nil")
		return
	}

	v, ok := m.Data.info[config.Id]
	if ok {
		// cd 增加逻辑
		if config.Timetype == 1 {
			v.Update(config.Timevalue)
			m.SynHeadInfo(v)
		}
	} else {
		pInfo := NewHeadInfo(config)
		if pInfo == nil {
			return
		}
		m.Data.info[config.Id] = pInfo
		m.SynHeadInfo(pInfo)
	}
}

func (m *ModHead) InitHead(id int) {
	config, ok := GetCsvMgr().HeadConfigMap[id]
	if !ok {
		LogError("config is nil")
		return
	}

	if m.Data.info == nil {
		LogError("m.Data.info is nil")
		return
	}

	v, ok := m.Data.info[config.Id]
	if ok {
		// cd 增加逻辑
		if config.Timetype == 1 {
			v.Update(config.Timevalue)
		}
	} else {
		pInfo := NewHeadInfo(config)
		if pInfo == nil {
			return
		}
		m.Data.info[config.Id] = pInfo
	}
}
