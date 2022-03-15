package game

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	//"time"
)

const (
	CHAT_WORLD       = 1  // 世界
	CHAT_PARTY       = 2  // 公会
	CHAT_CAMP        = 3  // 阵营
	CHAT_PRIVATE     = 4  // 私聊
	CHAT_SYSTEM      = 5  // 公告
	CHAT_CAMP_SYSTEM = 6  // 阵营系统聊天
	CHAT_TEAM        = 7  // 组队系统聊天
	CHAT_MINE        = 8  // 矿点战玩家
	CHAT_MINE_SYSTEM = 9  // 矿点战系统
	CHAT_GVE         = 10 // 孤山夺宝
	CHAT_PRIVATE_SVR = 11 //! 跨服聊天

	MAX_BARRAGE_INTERVAL    = 10 // 弹幕间隔
	MAX_WORLD_CHAT_INTERVAL = 10 // 世界聊天间隔
)

var channelName = []string{"", "世界", "帮派", "阵营", "私聊", "公告", "阵营系统", "组队系统", "矿点战玩家", "矿点战系统", "孤山夺宝"}

const (
	SENDSAVECHAT_ADDRESS = "http://chat-api.kokoyou.com/api/chat/save"
	//SENDSAVECHAT_ADDRESS = "http://usdk.api.kokoyou.com/cp/cp/check"
	GameID  = 2004
	GAMEKEY = "dbf77ad9205f89a174d9ba3d79468a3f"
)

type ModChat struct {
	player         *Player
	Barrage        int64 // 下一次可发弹幕时间
	WolrdChat      int64 // 下一次可发时间聊天时间
	WorldChatCount int   //世界聊天计数
	WorldChannel   int   //世界频道ID
}

func (self *ModChat) OnGetData(player *Player) {
	self.player = player
}

func (self *ModChat) OnGetOtherData() {
}

func (self *ModChat) OnMsg(ctrl string, body []byte) bool {
	if ctrl != "chat" {
		return false
	}
	if self.player.Sql_UserBase.IsGag != 0 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_FIGHT_YOU_HAVE_BEEN_BANNED_PLEASE"))
		return true
	}
	data := &C2S_NewChat{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return true
	}

	//敏感词直接return
	if GetServer().IsSensitiveWord(data.Content) {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_UIchatui_Ban_Text"))
		return true
	}

	//! 保存聊天记录
	GetServer().SqlChatLog(self.player.GetUid(), self.player.Sql_UserBase.UName+":"+data.Content, int(data.Channel))
	self.SendChatMessage(data, "")

	switch data.Channel {
	case CHAT_WORLD:
		self.player.UpdateNoticeChatTime()
		GetMasterMgr().ChatRPC.SendWorldMessage(self.player, data.Content, self.WorldChannel)
		return true
	case CHAT_PARTY:
		self.player.UpdateNoticeChatTime()
		unionId := self.player.GetUnionId()
		if unionId != 0 {
			GetMasterMgr().ChatRPC.SendUnionMessage(self.player, data.Content, unionId)
		}
		return true
	case CHAT_PRIVATE:
		self.player.UpdateNoticeChatTime()
		GetMasterMgr().ChatRPC.SendPrivateMessage(self.player, data.Content, int(data.PrivateUid))
		//保存私聊记录
		GetMasterMgr().FriendPRC.SavePrivateMessage(self.player, data.Content, int(data.PrivateUid))
		return true
	}
	return true
}

// 组队频道全服推送
func (self *ModChat) handleChatTeam(msg *S2C_Chat, content string) {

}

func (self *ModChat) OnSave(sql bool) {

}

// 检查弹幕时间间隔是否Ok
func (self *ModChat) CheckBarrage() bool {
	now := TimeServer().Unix()
	if self.Barrage > now {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CHAT_YOUR_SPEECH_IS_TOO_FAST"))
		return false
	}

	if self.player.Sql_UserBase.Vip < 4 && self.player.Sql_UserBase.Level < 50 {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CHAT_LEVEL_50_OR_NOBLE_4"))
		return false
	}

	self.Barrage = now + MAX_BARRAGE_INTERVAL
	return true
}

// 检查世界聊天时间间隔是否Ok.
func (self *ModChat) CheckWorldChat() bool {
	now := TimeServer().Unix()
	if self.WolrdChat > now {
		self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CHAT_YOUR_SPEECH_IS_TOO_FAST"))
		return false
	}

	self.WolrdChat = now + MAX_WORLD_CHAT_INTERVAL
	self.WorldChatCount++
	return true
}

// 矿点战消息推送
func (self *ModChat) handleChatMine(msg *S2C_Chat, content string) {

}

//! 趣炫-聊天后台接入
//! 格式：游戏缩写|平台别名|服id|时间|发送者平台账号|发送者角色ID|发送者角色名|vip等级|聊天类型|接受者平台账号|接收者角色ID|接收者角色名|接收者vip等级|聊天内容

func (self *ModChat) SendChatMessage(msg *C2S_NewChat, tar string) {
	gameTag := "S15"
	gameFlag := "s15qdazzle"

	switch msg.Channel {
	case CHAT_WORLD:
		{
			//! 世界
			//! g2|mixed|1|123456879|oppo_36505593|7974|陶震|8|1|||||我老婆要买xxxx。。
			//[时间] User[平台账号:角色ID:角色名] (avatar:数字vip:数字) ToChannel[-1] [聊天内容]
			//[2014-11-26 19:36:51] User[uc_302839541:190:狂] (avatar:51 vip:1) ToChannel[1000002] [情丝通通的交出来]
			sendMessage := fmt.Sprintf("%s|%s|%d|%d|%s|%d|%s|%d|%d|||||%s\r\n", gameTag, gameFlag,
				GetServer().Con.ServerId, TimeServer().Unix(), self.player.Account.UserId, self.player.GetUid(), self.player.GetUname(), self.player.Sql_UserBase.Vip,
				1, msg.Content)
			LogDebug("world msg: ", sendMessage)
			GetServer().SendMsg2Chat([]byte(sendMessage))
		}
	case CHAT_PRIVATE:
		{
			toPlayer := GetMasterMgr().GetPlayer(msg.PrivateUid)
			if toPlayer != nil {
				//! 私聊
				//! g2|mixed|1|123456879|oppo_36505593|7974|陶震|8|0|baidu_165412427|1637|枪炮玫瑰|10|我老婆要买xxxx。。
				//[时间] FromUser[平台账号:角色ID:角色名] chat_flag:数字 avatar:数字 vip:数字 ToUser[平台账号:角色ID:角色名] target_vip:数字 [聊天内容]
				// [2014-11-05 10:05:36] FromUser[xy_760260:80:冷晓烁] chat_flag:0 avatar:82 vip:7 ToUser[xy_523622:1428:车行天下] target_vip:10 [等一下]
				sendMessage := fmt.Sprintf("%s|%s|%d|%d|%s|%d|%s|%d|%d|%s|%d|%s|%d|%s\r\n", gameTag, gameFlag,
					GetServer().Con.ServerId, TimeServer().Unix(), self.player.Account.UserId, self.player.GetUid(), self.player.GetUname(), self.player.Sql_UserBase.Vip,
					0, toPlayer.Data.UserID, toPlayer.Data.UId, toPlayer.Data.UName, toPlayer.Data.Vip, msg.Content)
				LogDebug("private msg: ", sendMessage)
				GetServer().SendMsg2Chat([]byte(sendMessage))
			}
		}
	case CHAT_PARTY:
		{
			//! 公会
			//! g2|mixed|1|123456879|oppo_36505593|7974|陶震|8|2|10001||||我老婆要买xxxx。。
			// [2017-04-27 11:40:01] ChannelChat (Info): User[xd_102981428601:3190:雅绿念梦] (avatar:2 vip:0) ToChannel[12] [1|对香香好点，小心萧揍你]
			// [时间] User[平台账号:角色ID:角色名] (avatar:数字vip:数字) ToChannel[帮派id] [聊天内容]
			sendMessage := fmt.Sprintf("%s|%s|%d|%d|%s|%d|%s|%d|%d|%d||||%s\r\n", gameTag, gameFlag,
				GetServer().Con.ServerId, TimeServer().Unix(), self.player.Account.UserId, self.player.GetUid(), self.player.GetUname(), self.player.Sql_UserBase.Vip,
				2, self.player.GetUnionId(), msg.Content)
			LogDebug("party msg: ", sendMessage)
			GetServer().SendMsg2Chat([]byte(sendMessage))
		}

	}
}

//!  增加  接受聊天信息接口  20190430 by zy
func (self *ModChat) SendSaveChat(msgInfo *S2C_Chat, toName string) {

	toUserName := ""
	var toRoleId int64 = 0
	toRoleName := ""
	toVipInfo := 0

	toPlayer := GetPlayerMgr().GetPlayerFromName(toName, false)
	if toPlayer != nil && toPlayer.SessionObj != nil {
		toUserName = toPlayer.Account.Account
		toRoleId = toPlayer.Sql_UserBase.Uid
		toRoleName = toPlayer.Sql_UserBase.UName
		toVipInfo = toPlayer.Sql_UserBase.Vip
	}

	str := fmt.Sprintf("%d%d%d%d%d%s",
		GameID, msgInfo.Channel, self.player.Sql_UserBase.Uid, GetServer().Con.ServerId, msgInfo.Time, GAMEKEY)
	h := md5.New()
	h.Write([]byte(str))
	md5str := fmt.Sprintf("%x", h.Sum(nil))

	postvalue := url.Values{
		"gameId":          {fmt.Sprintf("%d", GameID)},
		"gamePlayerId":    {fmt.Sprintf("%s", self.player.Account.UserId)},
		"serverId":        {fmt.Sprintf("%d", GetServer().Con.ServerId)},
		"serverName":      {fmt.Sprintf("%s", GetServer().Con.ServerName)},
		"userName":        {fmt.Sprintf("%s", self.player.Account.Account)},
		"roleId":          {fmt.Sprintf("%d", self.player.Sql_UserBase.Uid)},
		"roleName":        {fmt.Sprintf("%s", self.player.Sql_UserBase.UName)},
		"vipInfo":         {fmt.Sprintf("%d", self.player.Sql_UserBase.Vip)},
		"toUserName":      {fmt.Sprintf("%s", toUserName)},
		"toRoleId":        {fmt.Sprintf("%d", toRoleId)},
		"toRoleName":      {fmt.Sprintf("%s", toRoleName)},
		"toVipInfo":       {fmt.Sprintf("%d", toVipInfo)},
		"chatChannelId":   {fmt.Sprintf("%d", msgInfo.Channel)},
		"chatChannelName": {fmt.Sprintf("%s", channelName[msgInfo.Channel])},
		"ip":              {fmt.Sprintf("%s", self.player.Sql_UserBase.IP)},
		"imei":            {fmt.Sprintf("%s", self.player.Platform.DeviceId)},
		"chatContent":     {fmt.Sprintf("%s", msgInfo.Content)},
		"chatTime":        {fmt.Sprintf("%d", msgInfo.Time)},
		"sign":            {fmt.Sprintf("%s", md5str)},
	}

	res, err := http.PostForm(SENDSAVECHAT_ADDRESS, postvalue)
	if err != nil {
		//log.Println(err)
		//self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CAMPTASK_CONFIGURATION_DATA_ERROR"))
		return
	}

	result, err := ioutil.ReadAll(res.Body)
	if res.Body != nil {
		defer res.Body.Close()
	}
	if err != nil {
		//log.Println(err)
		//self.player.SendErrInfo("err", GetCsvMgr().GetText("STR_MOD_CAMPTASK_CONFIGURATION_DATA_ERROR"))
		return
	}
	stringInfo := string(result)
	log.Println("result=", stringInfo)

}

func (self *ModChat) EnterChannel(channelId int) {
	self.WorldChannel = channelId
}

func (self *ModChat) GetWorldChannel() int {
	return self.WorldChannel
}

func (self *ModChat) WorldMessageRecord(msgList []*ChatMessage) {
	var msg S2C_NewChat
	msg.Cid = "chatrecord"
	msg.Channel = CHAT_WORLD
	msg.MsgList = msgList
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModChat) UnionMessageRecord(msgList []*ChatMessage) {
	var msg S2C_NewChat
	msg.Cid = "chatrecord"
	msg.Channel = CHAT_PARTY
	msg.MsgList = msgList
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}

func (self *ModChat) PrivateMessageRecord(msgList []*ChatMessage) {
	var msg S2C_NewChat
	msg.Cid = "chatrecord"
	msg.Channel = CHAT_PRIVATE
	msg.MsgList = msgList
	self.player.SendMsg(msg.Cid, HF_JtoB(&msg))
}
