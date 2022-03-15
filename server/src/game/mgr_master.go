/*
@Time : 2020/5/6 8:16
@Author : 96121
@File : mgr_master
@Software: GoLand
*/
package game

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

const (
	WARN_ERROR           = 1
	WARN_OK              = 2
	WARN_MASTER_TIME_OUT = 3
)

//! 中心服务器
type MasterMgr struct {
	Init bool   //! 是否初始化
	Host string //! 中心服务器地址

	WarnTime           int64 //!警告通知时间
	WarnTimeFroTimeOut int64 //!警告通知时间

	Client    *rpc.Client   //! RPC连接缓存
	PlayerPRC *RPC_Player   //! 角色信息PRC接口
	FriendPRC *RPC_Friend   //! 好友RPC接口
	ServerRPC *RPC_Server   //! 服务器RPC接口
	UnionRPC  *RPC_Union    //! 公会RPC接口
	TowerRPC  *RPC_Tower    //! 塔RPC接口
	ChatRPC   *RPC_Chat     //! 聊天
	MatchRPC  *RPC_Match    //! 跨服比赛
	Locker    *sync.RWMutex //! 数据锁
}

var s_mastermgr *MasterMgr

func GetMasterMgr() *MasterMgr {
	if s_mastermgr == nil {
		s_mastermgr = new(MasterMgr)

		s_mastermgr.Client = nil
		s_mastermgr.ServerRPC = new(RPC_Server)

		s_mastermgr.PlayerPRC = new(RPC_Player)
		s_mastermgr.FriendPRC = new(RPC_Friend)
		s_mastermgr.UnionRPC = new(RPC_Union)
		s_mastermgr.TowerRPC = new(RPC_Tower)
		s_mastermgr.ChatRPC = new(RPC_Chat)
		s_mastermgr.MatchRPC = new(RPC_Match)

		s_mastermgr.WarnTimeFroTimeOut = TimeServer().Unix() + MIN_SECS*5
		s_mastermgr.Init = false
	}

	return s_mastermgr
}

func (self *MasterMgr) InitService() bool {
	if self.Init {
		return true
	}
	self.Locker = new(sync.RWMutex)
	self.PlayerPRC.Init()
	self.FriendPRC.Init()
	self.ServerRPC.Init()
	self.UnionRPC.Init()
	self.TowerRPC.Init()
	self.ChatRPC.Init()
	self.MatchRPC.Init()

	timeout := time.Second * 30
	conn, err := net.DialTimeout("tcp", GetServer().Con.ServerExtCon.MasterSvr, timeout)
	if err != nil {
		self.Init = false
		//log.Println("dailing error: ", err)

		return false
	}

	self.Init = true
	self.Client = rpc.NewClient(conn)

	// timeout 2s
	readAndWriteTimeout := 2 * time.Second
	err = conn.SetDeadline(TimeServer().Add(readAndWriteTimeout))
	if err != nil {
		log.Println("SetDeadline failed:", err)
		return false
	}

	//! 初始化连接
	self.ServerRPC.Client = self.Client
	self.PlayerPRC.Client = self.Client
	self.FriendPRC.Client = self.Client
	self.UnionRPC.Client = self.Client
	self.TowerRPC.Client = self.Client
	self.ChatRPC.Client = self.Client
	self.MatchRPC.Client = self.Client

	return true
}

func (self *MasterMgr) ResumeService(force bool) bool {
	if self.Init == true || force == true {
		conn, err := rpc.DialHTTP("tcp", GetServer().Con.ServerExtCon.MasterSvr)
		if err != nil {
			//self.Init = false
			self.Client = nil
			self.ServerRPC.Client = nil
			self.PlayerPRC.Client = nil
			self.FriendPRC.Client = nil
			self.UnionRPC.Client = nil
			self.TowerRPC.Client = nil
			self.ChatRPC.Client = nil
			self.MatchRPC.Client = nil
			//log.Println("dailing error: ", err)

			if GetServer().Con.ServerExtCon.MasterWarnClose == LOGIC_FALSE && self.WarnTime < TimeServer().Unix() {
				self.WarnTime = TimeServer().Unix() + MIN_SECS*10
				warnStr := self.MakeWarnStr(WARN_ERROR)
				GetWechatWarningMgr().SendWarning(warnStr)
			}
		} else {
			//! 重新处理连接
			self.Client = conn

			//! 初始化连接
			self.ServerRPC.Client = conn
			self.PlayerPRC.Client = conn
			self.FriendPRC.Client = conn
			self.UnionRPC.Client = conn
			self.TowerRPC.Client = conn
			self.ChatRPC.Client = conn
			self.MatchRPC.Client = conn

			if GetServer().Con.ServerExtCon.MasterWarnClose == LOGIC_FALSE {
				warnStr := self.MakeWarnStr(WARN_OK)
				GetWechatWarningMgr().SendWarning(warnStr)
			}
		}

		return true
	}

	return false
}

// 增加红包活动检查, 一分钟检查一下,差不多12,18点的时候发送红包
func (self *MasterMgr) OnTimer() {
	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		//self.RegServer()
		self.ReqEvents()
		//self.ReqUpdatePlayer()   //改在PLAYER里刷新，避免多线程MAP同时读写
	}

	ticker.Stop()
}

//! 逻辑处理
func (self *MasterMgr) OnLogic() {
	for event := range self.ServerRPC.EventChan {
		LogDebug("Proc Event :", event.EventCode, event.Target, event.Param1, event.Param2)
		switch event.EventCode {
		case PLAYER_EVENT_ADD_FRIEND:
			player := GetPlayerMgr().GetPlayer(event.UId, false)
			if player != nil {
				info := &JS_Friend{}
				json.Unmarshal([]byte(event.Param2), &info)
				player.GetModule("friend").(*ModFriend).AddApply(info)
			}
		case PLAYER_EVENT_DEL_FRIEND, PLAYER_EVENT_BLACK_FRIEND:
			player := GetPlayerMgr().GetPlayer(event.UId, false)
			if player != nil {
				info := &JS_Friend{}
				json.Unmarshal([]byte(event.Param2), &info)
				player.GetModule("friend").(*ModFriend).delApply(info.Uid)
				player.GetModule("friend").(*ModFriend).DelFriend(info.Uid, true, LOGIC_FALSE)
				player.GetModule("friend").(*ModFriend).DelHasapply(info.Uid)
			}
		case PLAYER_EVENT_AGREEE_FRIEND:
			player := GetPlayerMgr().GetPlayer(event.UId, false)
			if player != nil {
				info := &JS_Friend{}
				json.Unmarshal([]byte(event.Param2), &info)
				if info == nil || info.Uid == 0 {
					return
				}
				player.GetModule("friend").(*ModFriend).delApply(info.Uid)
				player.GetModule("friend").(*ModFriend).AddFriendByCenter(info)
			}
		case PLAYER_EVENT_REFUSE_FRIEND:
			player := GetPlayerMgr().GetPlayer(event.UId, false)
			if player != nil {
				info := &JS_Friend{}
				json.Unmarshal([]byte(event.Param2), &info)
				player.GetModule("friend").(*ModFriend).delApply(info.Uid)
			}
		case PLAYER_EVENT_POWER_FRIEND:
			player := GetPlayerMgr().GetPlayer(event.UId, true)
			if player != nil {
				player.GetModule("friend").(*ModFriend).SetGift(event.Target, 1)
			}
		case PLAYER_EVENT_UPDATE_HIRE_HERO:
			player := GetPlayerMgr().GetPlayer(event.UId, true)
			if player != nil {
				hire := make([]*HireHero, 0)
				json.Unmarshal([]byte(event.Param2), &hire)
				for _, v := range hire {
					player.GetModule("friend").(*ModFriend).HireStateUpdate(v)
				}
			}
		case PLAYER_EVENT_UPDATE_HIRE_HERO_SINGLE:
			player := GetPlayerMgr().GetPlayer(event.UId, true)
			if player != nil {
				singleHire := new(HireHero)
				json.Unmarshal([]byte(event.Param2), &singleHire)
				player.GetModule("friend").(*ModFriend).HireStateUpdate(singleHire)
			}
		case PLAYER_EVENT_AGREE_HIRE_HERO:
			player := GetPlayerMgr().GetPlayer(event.UId, true)
			if player != nil {
				hire := &HireHero{}
				json.Unmarshal([]byte(event.Param2), &hire)
				if hire == nil || hire.OwnPlayer == nil {
					return
				}
				player.GetModule("friend").(*ModFriend).AddHireHero(hire)
				player.GetModule("friend").(*ModFriend).RemoveApplyHire(hire)
			}
		case PLAYER_EVENT_REFUSE_HIRE_HERO:
			player := GetPlayerMgr().GetPlayer(event.UId, true)
			if player != nil {
				hire := &HireHero{}
				json.Unmarshal([]byte(event.Param2), &hire)
				if hire == nil || hire.OwnPlayer == nil {
					return
				}
				player.GetModule("friend").(*ModFriend).RemoveApplyHire(hire)
			}
		case PLAYER_EVENT_DELETE_HIRE:
			player := GetPlayerMgr().GetPlayer(event.UId, false)
			if player != nil {
				hire := &HireHero{}
				json.Unmarshal([]byte(event.Param2), &hire)
				if hire == nil || hire.OwnPlayer == nil {
					return
				}
				player.GetModule("friend").(*ModFriend).DeleteHire(hire)
			}
		case UNION_EVENT_MASTER_FAIL:
			player := GetPlayerMgr().GetPlayer(event.UId, true)
			if nil != player {
				player.GetModule("union").(*ModUnion).CancelApplyUnion(event.Param1, event.UId, false)
			}
		case UNION_EVENT_MASTER_OK:
			player := GetPlayerMgr().GetPlayer(event.UId, true)
			if nil != player {
				player.GetModule("union").(*ModUnion).ClearApplyUnion()
				player.GetModule("union").(*ModUnion).AddPlayer(event.Param1)

				player.SendInfo("joinunionok")
				////! 同步数据
				//(player.GetModule("union").(*ModUnion)).GetUserUnionInfo()
			}

			break
		case UNION_EVENT_OUT_PLAYER:
			player := GetPlayerMgr().GetPlayer(event.UId, true)
			if nil != player {
				playerunionmod := player.GetModule("union").(*ModUnion)

				playerunionmod.OutUnionData(event.Param1)

				if event.Target == 1 {
					player.SendRet("outplayerbymaster", 0)
				} else {
					player.SendRet("outunion", 0)
				}
			}
		case UNION_EVENT_UNION_MODIFY:
			player := GetPlayerMgr().GetPlayer(event.UId, true)
			if nil != player {
				player.GetModule("union").(*ModUnion).ModifyPosition(event.Param1)
				var msg S2C_UnionModify
				msg.Cid = "unionmodify"
				msg.Uid = event.Target
				msg.Destuid = HF_AtoI64(event.Param2)
				msg.Op = event.Param1
				smsg, _ := json.Marshal(&msg)
				player.SendMsg("unionmodify", smsg)
			}
		case UNION_EVENT_UNION_SET_BRAVE_HAND:
			player := GetPlayerMgr().GetPlayer(event.UId, true)
			if nil != player {
				var msg S2C_SetBraveHand
				msg.Cid = "setbravehand"
				msg.Uid = event.Target
				msg.Destuid = HF_AtoI64(event.Param2)
				msg.Op = event.Param1
				smsg, _ := json.Marshal(&msg)
				player.SendMsg("setbravehand", smsg)

			}
		case UNION_EVENT_UNION_HUNTER_AWARD:
			player := GetPlayerMgr().GetPlayer(event.UId, true)
			if nil != player {
				config := GetCsvMgr().GetUnionHuntConfigByID(event.Param1)
				if config == nil {
					continue
				}
				dropConfig := GetCsvMgr().GetUnionHuntDropConfig(config.Group, event.Target)
				if dropConfig == nil {
					continue
				}
				if len(dropConfig.Guilditem) != len(dropConfig.Guildnum) {
					continue
				}

				pMail := player.GetModule("mail").(*ModMail)
				if pMail == nil {
					continue
				}

				var itemlst []PassItem

				for q, p := range dropConfig.Guilditem {
					if p != 0 {
						itemlst = append(itemlst, PassItem{p, dropConfig.Guildnum[q]})
					}
				}

				title := fmt.Sprintf(GetCsvMgr().GetText("公会狩猎奖励"))
				text := fmt.Sprintf(GetCsvMgr().GetText("公会狩猎奖励"))
				pMail.AddMail(1, 1, 0, title, text, GetCsvMgr().GetText("STR_SYS"), itemlst, false, 0)
			}
		case UNION_EVENT_UNION_UPDATE:
			GetUnionMgr().UpdateUnion(event.Param1)
		case UNION_EVENT_UNION_SEND_MAIL:
			player := GetPlayerMgr().GetPlayer(event.UId, true)
			if nil != player {
				pMail := player.GetModule("mail").(*ModMail)
				if pMail != nil {
					var mail UnionMail
					json.Unmarshal([]byte(event.Param2), &mail)
					pMail.AddMail(1, 1, 0, mail.Title, mail.Text, "公会管理员", []PassItem{}, true, 0)
				}
			}
		case CHAT_NEW_WORLD_MESSAGE:
			player := GetPlayerMgr().GetPlayer(event.UId, false)
			if nil != player {
				player.UpdateNoticeChatTime()
				chatMsg := make([]*ChatMessage, 0)
				json.Unmarshal([]byte(event.Param2), &chatMsg)
				if chatMsg != nil && len(chatMsg) > 0 {
					var msg S2C_NewChat
					msg.Cid = "chat"
					msg.Channel = CHAT_WORLD
					msg.MsgList = chatMsg
					player.SendMsg(msg.Cid, HF_JtoB(&msg))
				}
			}
		case CHAT_NEW_UNION_MESSAGE:
			player := GetPlayerMgr().GetPlayer(event.UId, false)
			if nil != player {
				player.UpdateNoticeChatTime()
				chatMsg := make([]*ChatMessage, 0)
				json.Unmarshal([]byte(event.Param2), &chatMsg)
				if chatMsg != nil && len(chatMsg) > 0 {
					var msg S2C_NewChat
					msg.Cid = "chat"
					msg.Channel = CHAT_PARTY
					msg.MsgList = chatMsg
					player.SendMsg(msg.Cid, HF_JtoB(&msg))
				}
			}
		case CHAT_NEW_PRIVATE_MESSAGE:
			player := GetPlayerMgr().GetPlayer(event.UId, false)
			if nil != player {
				player.UpdateNoticeChatTime()
				chatMsg := make([]*ChatMessage, 0)
				json.Unmarshal([]byte(event.Param2), &chatMsg)
				if chatMsg != nil && len(chatMsg) > 0 {
					var msg S2C_NewChat
					msg.Cid = "chat"
					msg.Channel = CHAT_PRIVATE
					msg.MsgList = chatMsg
					player.SendMsg(msg.Cid, HF_JtoB(&msg))
				}
			}
		case CHAT_GAP_PLAYER:
			player := GetPlayerMgr().GetPlayer(event.UId, false)
			if nil != player {
				var msg S2C_NewChatGap
				msg.Cid = "gapchat"
				msg.Channel = event.Param1
				msg.GapUid = event.Target
				player.SendMsg(msg.Cid, HF_JtoB(&msg))
			}
		case MATCH_CROSSARENA_UPDATE:
			data := &Js_CrossArenaUser{}
			json.Unmarshal([]byte(event.Param2), &data)
			if data != nil {
				GetCrossArenaMgr().UpdateInfo(data)
			}
		case MATCH_CROSSARENA_3V3_UPDATE:
			data := &Js_CrossArena3V3User{}
			json.Unmarshal([]byte(event.Param2), &data)
			if data != nil {
				GetCrossArena3V3Mgr().UpdateInfo(data)
			}
		}
	}
	self.ServerRPC.EventChan = make(chan *RPC_ServerEvent)
}

func (self *MasterMgr) GetClient() *rpc.Client {
	return self.Client
}

//! 服务器结构
type GameSvrNode struct {
	ID     int    //! 服务器Id
	Name   string //! 名字
	Online int    //! 在线人数
}

//! 注册服务，心跳操作
func (self *MasterMgr) RegServer() {
	if TimeServer().Unix()%10 == 0 {
		//LogDebug("update server ...", TimeServer())
		self.ServerRPC.RegServer(GetServer().Con.ServerId,
			GetPlayerMgr().GetPlayerOnline(),
			GetServer().Con.ServerName)
	}

	//! 请求事件
	//self.ServerRPC.ReqEvent(GetServer().Con.ServerId)
}

//! 请求服务器事件
func (self *MasterMgr) ReqEvents() {
	if self.ServerRPC != nil {
		ret := self.ServerRPC.ReqEvent(GetServer().Con.ServerId)
		if ret == true {
			//! 获取到事件
			ret = self.ServerRPC.ReqEventArr(GetServer().Con.ServerId)

			if ret == true {
				//! 递归获取事件-慎用
				self.ReqEvents()
			}
		}
	}
}

//! 同步角色信息  20200730已取消调用  改由player调用同步
func (self *MasterMgr) ReqUpdatePlayer() {
	playerList := GetPlayerMgr().GetOnlineRandom(1)

	for i := 0; i < len(playerList); i++ {
		self.UpdatePlayer(playerList[i])
	}
}

//! 注册服务，心跳操作
func (self *MasterMgr) UpdatePlayer(player *Player) {
	if player == nil {
		return
	}
	if !player.NoticeBaseInfo {
		return
	}
	LogDebug("update player =>", player.Sql_UserBase.Uid, player.Sql_UserBase.UName)
	//! 调用接口
	self.PlayerPRC.RegPlayer(player)
}

//! 获取远程角色数据
func (self *MasterMgr) GetPlayer(uid int64) *RPC_PlayerData_Req {
	//! 从接口获取，接口缓存
	return self.PlayerPRC.GetPlayer(uid)
}

//! 设置离线
func (self *MasterMgr) SetPlayerOffline(player *Player) {
	self.PlayerPRC.SetPlayerOffline(player)
}

//加入聊天频道
func (self *MasterMgr) EnterChat(player *Player) {
	if player == nil {
		return
	}
	//! 调用接口
	rel := self.ChatRPC.EnterChat(player.Sql_UserBase.Uid)
	if rel != nil {
		player.GetModule("chat").(*ModChat).EnterChannel(rel.Channel)
	}
}

func (self *MasterMgr) ExitChat(player *Player) {
	if player == nil {
		return
	}
	//! 调用接口
	self.ChatRPC.ExitChat(player.Sql_UserBase.Uid)
}

//获取世界聊天
func (self *MasterMgr) QueryWorldMessage(player *Player) {
	if player == nil {
		return
	}
	//! 调用接口
	rel := self.ChatRPC.QueryWorldMessage(player)
	if rel != nil {
		player.GetModule("chat").(*ModChat).WorldMessageRecord(rel.MsgList)
	}
}

func (self *MasterMgr) QueryUnionMessage(player *Player) {
	if player == nil {
		return
	}
	//! 调用接口
	rel := self.ChatRPC.QueryUnionMessage(player)
	if rel != nil {
		player.GetModule("chat").(*ModChat).UnionMessageRecord(rel.MsgList)
	}
}

func (self *MasterMgr) QueryPrivateMessage(player *Player) {
	if player == nil {
		return
	}
	//! 调用接口
	rel := self.FriendPRC.QueryPrivateMessage(player)
	if rel != nil {
		player.GetModule("chat").(*ModChat).PrivateMessageRecord(rel.MsgList)
	}
}

//生成警告内容  1断开连接    恢复连接
func (self *MasterMgr) MakeWarnStr(nType int) string {
	strTitle := ""

	strHost := fmt.Sprintf("serverid:%d\n", GetServer().Con.ServerId)
	strCenter := fmt.Sprintf("masterip:%s\n", GetServer().Con.ServerExtCon.MasterSvr)
	strTime := fmt.Sprintf("Time:%s\n", time.Now().Format("2006-01-02 15:04:05"))

	strMsg := ""
	switch nType {
	case WARN_ERROR:
		strTitle = fmt.Sprintf("error:%s连接中心服监控\n", GetServer().Con.GameName)
		strMsg = fmt.Sprintf("message:%s\n", "连接断开")
	case WARN_OK:
		strTitle = fmt.Sprintf("ok:%s连接中心服监控\n", GetServer().Con.GameName)
		strMsg = fmt.Sprintf("message:%s\n", "连接成功")
	case WARN_MASTER_TIME_OUT:
		strTitle = fmt.Sprintf("ok:%s中心服超时\n", GetServer().Con.GameName)
		strMsg = fmt.Sprintf("message:%s\n", "连接超时")
	default:
		return fmt.Sprintf("未知警报类型:%d\n", nType)
	}
	return strTitle + strHost + strCenter + strTime + strMsg
}

func (self *MasterMgr) MatchGeneralUpdate(top *Js_GeneralUser, records []*GeneralRecord) *RPC_GeneralActionRes {
	return self.MatchRPC.MatchGeneralUpdate(top, records)
}

func (self *MasterMgr) MatchGeneralGetAllRank(keyId int, serverId int) *RPC_GeneralActionRes {
	return self.MatchRPC.MatchGeneralGetAllRank(keyId, serverId)
}

func (self *MasterMgr) MatchCrossArenaGetAllRank(keyId int) *RPC_CrossArenaActionRes {
	return self.MatchRPC.MatchCrossArenaGetAllRank(keyId)
}

func (self *MasterMgr) MatchCrossArenaAdd(keyId int, data *Js_CrossArenaUser, fightInfo *JS_FightInfo) *RPC_CrossArenaActionRes {
	return self.MatchRPC.MatchCrossArenaAdd(keyId, data, fightInfo)
}

func (self *MasterMgr) MatchCrossArenaGetDefence(keyId int, player *Player) *RPC_CrossArenaGetDefenceRes {
	return self.MatchRPC.MatchCrossArenaGetDefence(keyId, player)
}

func (self *MasterMgr) MatchCrossArenaGetInfo(keyId int, uid int64) *RPC_CrossArenaGetInfoRes {
	return self.MatchRPC.MatchCrossArenaGetInfo(keyId, uid)
}

func (self *MasterMgr) MatchCrossArenaFightEnd(keyId int, attack *JS_FightInfo, defend *JS_FightInfo, battleInfo BattleInfo) *RPC_CrossArenaActionRes {
	return self.MatchRPC.MatchCrossArenaFightEnd(keyId, attack, defend, battleInfo)
}

func (self *MasterMgr) MatchCrossArenaGetBattleInfo(keyId int64) *RPC_CrossArenaBattleInfoRes {
	return self.MatchRPC.MatchCrossArenaGetBattleInfo(keyId)
}

func (self *MasterMgr) MatchCrossArenaGetBattleRecord(keyId int64) *RPC_CrossArenaBattleRecordRes {
	return self.MatchRPC.MatchCrossArenaGetBattleRecord(keyId)
}

func (self *MasterMgr) CallEx(client *rpc.Client, serviceMethod string, args interface{}, reply interface{}) error {
	timeout := time.Duration(3 * time.Second)
	done := make(chan error, 1)
	go func() {
		call := <-client.Go(serviceMethod, args, reply, make(chan *rpc.Call, 1)).Done
		done <- call.Error
	}()

	select {
	case <-time.After(timeout):
		if self.WarnTimeFroTimeOut < TimeServer().Unix() {
			self.WarnTimeFroTimeOut = TimeServer().Unix() + MIN_SECS*10
			warnStr := self.MakeWarnStr(WARN_MASTER_TIME_OUT)
			GetWechatWarningMgr().SendWarning(warnStr)
		}
		return fmt.Errorf("timeout")
	case err := <-done:
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *MasterMgr) MatchConsumerTopGetAllRank(keyId int, serverId int) *RPC_ConsumerTopActionRes {
	return self.MatchRPC.MatchConsumerTopGetAllRank(keyId, serverId)
}

func (self *MasterMgr) MatchConsumerTopUpdate(top *JS_ConsumerTopUser) *RPC_ConsumerTopActionRes {
	return self.MatchRPC.MatchConsumerTopUpdate(top)
}

func (self *MasterMgr) MatchCrossArena3V3GetAllRank(keyId int) *RPC_CrossArena3V3ActionRes {
	return self.MatchRPC.MatchCrossArena3V3GetAllRank(keyId)
}

func (self *MasterMgr) MatchCrossArena3V3Add(keyId int, data *Js_CrossArena3V3User, fightInfo [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo) *RPC_CrossArena3V3ActionRes {
	return self.MatchRPC.MatchCrossArena3V3Add(keyId, data, fightInfo)
}

func (self *MasterMgr) MatchCrossArena3V3GetDefence(keyId int, player *Player) *RPC_CrossArena3V3GetDefenceRes {
	return self.MatchRPC.MatchCrossArena3V3GetDefence(keyId, player)
}

func (self *MasterMgr) MatchCrossArena3V3GetInfo(keyId int, uid int64) *RPC_CrossArena3V3GetInfoRes {
	return self.MatchRPC.MatchCrossArena3V3GetInfo(keyId, uid)
}

func (self *MasterMgr) MatchCrossArena3V3FightEnd(keyId int, attack [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo, defend [CROSSARENA3V3_TEAM_MAX]*JS_FightInfo, battleInfo [CROSSARENA3V3_TEAM_MAX]BattleInfo) *RPC_CrossArena3V3ActionRes {
	return self.MatchRPC.MatchCrossArena3V3FightEnd(keyId, attack, defend, battleInfo)
}

func (self *MasterMgr) MatchCrossArena3V3GetBattleInfo(keyId int64) *RPC_CrossArena3V3BattleInfoRes {
	return self.MatchRPC.MatchCrossArena3V3GetBattleInfo(keyId)
}

func (self *MasterMgr) MatchCrossArena3V3GetBattleRecord(keyId int64) *RPC_CrossArena3V3BattleRecordRes {
	return self.MatchRPC.MatchCrossArena3V3GetBattleRecord(keyId)
}