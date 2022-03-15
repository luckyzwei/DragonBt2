package game

import (
	"encoding/json"
	"fmt"
	//"time"
)

const (
	MSG_ACTIVATE_SKIN  = "msg_activate_skin"  // 激活皮肤
	MSG_SET_SKIN       = "msg_set_skin"       // 设置皮肤
	MSG_SEND_SKIN_INFO = "msg_send_skin_info" // 发送信息
)

type JS_Skin struct {
	ID      int   `json:"id"`      // id
	HeroID  int   `json:"heroid"`  // heroid
	EndTime int64 `json:"endtime"` // 结束时间
}

// 悬赏任务
type San_Skin struct {
	Uid  int64  // 角色ID
	Info string // 玩家皮肤信息

	info []*JS_Skin // 玩家皮肤信息
	DataUpdate
}

// 悬赏
type ModSkin struct {
	player *Player // 玩家

	Sql_Skin San_Skin // 信息数据
}

// 获得数据
func (self *ModSkin) OnGetData(player *Player) {
	self.player = player

	sql := fmt.Sprintf("select * from `san_skin` where uid = %d", self.player.ID)
	GetServer().DBUser.GetOneData(sql, &self.Sql_Skin, "san_skin", self.player.ID)

	if self.Sql_Skin.Uid <= 0 {
		self.Sql_Skin.Uid = self.player.ID
		self.CheckSkin()

		self.Encode()
		InsertTable("san_skin", &self.Sql_Skin, 0, true)
		self.Sql_Skin.Init("san_skin", &self.Sql_Skin, true)
	} else {
		self.Decode()
		self.Sql_Skin.Init("san_skin", &self.Sql_Skin, true)
	}
}

// 获得数据
func (self *ModSkin) OnGetOtherData() {

}

// save
func (self *ModSkin) Decode() {
	json.Unmarshal([]byte(self.Sql_Skin.Info), &self.Sql_Skin.info)
}

// read
func (self *ModSkin) Encode() {
	self.Sql_Skin.Info = HF_JtoA(self.Sql_Skin.info)
}

// 存储
func (self *ModSkin) OnSave(sql bool) {
	self.Encode()
	self.Sql_Skin.Update(sql)
}

// 消息
func (self *ModSkin) OnMsg(ctrl string, body []byte) bool {
	return false
}

// 注册消息
func (self *ModSkin) onReg(handlers map[string]func(body []byte)) {
	handlers[MSG_ACTIVATE_SKIN] = self.ActivateSkin
	handlers[MSG_SET_SKIN] = self.SetSkin
	handlers[MSG_SEND_SKIN_INFO] = self.SendInfo
}
func (self *ModSkin) CheckSkin() {
	now := TimeServer().Unix()
	nLen := len(self.Sql_Skin.info)
	// 检查过期
	for i := nLen - 1; i >= 0; i-- {
		if self.Sql_Skin.info[i].EndTime == 0 {
			continue
		}

		if now >= self.Sql_Skin.info[i].EndTime {
			id := self.Sql_Skin.info[i].ID
			self.Sql_Skin.info = append(self.Sql_Skin.info[:i], self.Sql_Skin.info[i+1:]...)
			// 获得配置
			config := GetCsvMgr().GetHeroSkinConfig(id)
			if config == nil {
				continue
			}

			// 恢复默认皮肤
			heros := self.player.getHeroes()
			for _, v := range heros {
				if v.HeroId == config.Hero && v.Skin == config.ModId {
					v.Skin = 0
				}
			}
		}
	}

	//// 检查默认皮肤激活
	//configs := GetCsvMgr().HeroSkinConfig
	//for _, config := range configs {
	//	if config.ModId == 0 {
	//		find := false
	//		for _, v := range self.Sql_Skin.info {
	//			if v.ID == config.ID {
	//				find = true
	//				break
	//			}
	//		}
	//		if !find {
	//			self.Sql_Skin.info = append(self.Sql_Skin.info, &JS_Skin{config.ID, config.Hero, 0})
	//		}
	//	}
	//}
}

// 激活皮肤
func (self *ModSkin) ActivateSkin(body []byte) {
	var msg C2S_ActivateSkin
	json.Unmarshal(body, &msg)

	// 检查皮肤
	self.CheckSkin()
	for _, v := range self.Sql_Skin.info {
		if v.ID == msg.ID {
			self.player.SendErrInfo("err", "已解锁")
			return
		}
	}

	// 获得配置
	config := GetCsvMgr().GetHeroSkinConfig(msg.ID)
	if config == nil {
		self.player.SendErrInfo("err", "没找到配置")
		return
	}

	if config.ModId == 0 {
		self.player.SendErrInfo("err", "默认皮肤不用激活")
		return
	}

	// 扣除物品
	if self.player.HasObjectOkEasy(config.Item, config.Num) != nil {
		self.player.SendErrInfo("err", "解锁物品不足")
		return
	}
	var backmsg S2C_ActivateSkin
	backmsg.Items = self.player.RemoveObjectEasy(config.Item, config.Num, "激活皮肤", msg.ID, 0, 0)

	data := &JS_Skin{msg.ID, config.Hero, 0}
	// 激活
	self.Sql_Skin.info = append(self.Sql_Skin.info, data)
	backmsg.Data = data

	backmsg.Cid = MSG_ACTIVATE_SKIN
	backmsg.ID = msg.ID
	self.player.SendMsg(MSG_ACTIVATE_SKIN, HF_JtoB(&backmsg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SKIN_ACTIVITE, msg.ID, config.Hero, 0, "激活皮肤", 0, 0, self.player)
}

// 设置皮肤
func (self *ModSkin) SetSkin(body []byte) {
	var msg C2S_SetSkin
	json.Unmarshal(body, &msg)

	// 检查皮肤
	self.CheckSkin()

	// 获得配置
	config := GetCsvMgr().GetHeroSkinConfig(msg.ID)
	if config == nil {
		self.player.SendErrInfo("err", "没找到配置")
		return
	}

	if config.ModId != 0 {
		find := false
		for _, v := range self.Sql_Skin.info {
			if v.ID == msg.ID {
				find = true
				break
			}
		}
		if !find {
			self.player.SendErrInfo("err", "未解锁")
			return
		}
	}

	// 检查英雄
	hero := self.player.getHero(msg.HeroIndex)
	if hero == nil {
		self.player.SendErrInfo("err", "英雄不存在")
		return
	}
	if hero.HeroId != config.Hero {
		self.player.SendErrInfo("err", "不是该英雄")
		return
	}

	oldskin := hero.Skin
	// 设置
	hero.Skin = config.ModId
	GetOfflineInfoMgr().UpdateHeroSkin(self.player)
	var backmsg S2C_SetSkin
	backmsg.Cid = MSG_SET_SKIN
	backmsg.ID = msg.ID
	backmsg.HeroIndex = msg.HeroIndex
	self.player.SendMsg(MSG_SET_SKIN, HF_JtoB(&backmsg))

	GetServer().SqlLog(self.player.Sql_UserBase.Uid, LOG_SKIN_SET, msg.ID, oldskin, hero.HeroId, "更换皮肤", 0, hero.HeroKeyId, self.player)
}

// 设置皮肤
func (self *ModSkin) SendInfo(body []byte) {
	self.CheckSkin()

	var backmsg S2C_SendInfo
	backmsg.Cid = MSG_SEND_SKIN_INFO
	backmsg.Info = self.Sql_Skin.info
	self.player.SendMsg(MSG_SEND_SKIN_INFO, HF_JtoB(&backmsg))
}
