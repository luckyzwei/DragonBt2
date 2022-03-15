package gate

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"master/core"
	"master/network"
	"master/utils"
	"runtime/debug"
	"time"
)

//! 消息处理句柄类型
//! username
//! serverid
//! password
//! ctrl
//type SDK_Handler map[string]func(string, int, string, string) (*core.San_Account, string)

//! 登录模块
type GateApp struct {
	//SdkHandler SDK_Handler
}

var s_gateapp *GateApp

func GetGateApp() *GateApp {
	if s_gateapp == nil {
		s_gateapp = new(GateApp)

		//s_centerapp.SdkHandler = make(SDK_Handler)
	}

	return s_gateapp
}

//! 得到一个websocket处理句柄
func (self *GateApp) GetConnectHandler() websocket.Handler {
	connectHandler := func(ws *websocket.Conn) {
		if core.GetMasterApp().IsClosed() {
			return
		}
		session := network.GetSessionMgr().GetNewSession(ws)
		if session == nil {
			return
		}
		utils.LogDebug("add session:", session.ID)

		session.SetOnMessage(self.OnRecv)
		session.SetOnClose(self.OnCloseSession)

		session.Run()
	}
	return websocket.Handler(connectHandler)
}

//! 初始化处理接口
func (self *GateApp) InitHandler() {
	//self.SdkHandler[SDK_DEFAULT] = sdk.OnReg
	//self.SdkHandler[SDK_KOYOU] = sdk.SDKReg_KoYou
	//self.SdkHandler[SDK_KOYOU_IOS] = sdk.SDKReg_KoYou
}

//! 连接断开消息处理
func (self *GateApp) OnCloseSession(session *network.Session) {
	if session.PlayerObj != nil && session.PlayerObj.GetSession().(*network.Session) == session {
		utils.LogDebug("玩家离线，同步离线数据", session.ID, session.PlayerObj.GetUid())
		//core.GetPlayerMgr().SetPlayerOffline()
		utils.LogDebug("减少玩家离线数量")
		session.PlayerObj.OnClose()
		session.PlayerObj = nil
	}

	//GetLineUpMgr().RemoveClient(session)
	utils.LogDebug("删除玩家连接信息")
}

//! 处理数据
func (self *GateApp) OnRecv(session *network.Session, msg []byte) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			utils.LogError(x, string(debug.Stack()))
		}
	}()

	if session.ShutDown {
		return
	}

	head, body, _ := network.HF_EncodeMsg(msg)
	utils.LogDebug("onReceiveNew c2s:", head, "....", string(body))

	var ctrlhead C2S_CtrlHead
	err := json.Unmarshal(body, &ctrlhead)
	if err != nil {
		utils.LogError("CtrlType err:", err, string(body))
		return
	}

	//! 检查版本
	verCheck := self.CheckVer(ctrlhead.Ver)
	if verCheck == 1 {
		session.SendReturn("need_update")
	} else if verCheck == 2 {
		session.SendReturn("need_update")
	}

	//! 特殊处理
	//switch head {
	//case HEAD_PASSPORT:
	//	if serverConf.NetworkCon.MaxPlayer > core.GetPlayerMgr().GetOnline() {
	//		GetLineUpMgr().AddWaitClient(session)
	//		session.LoginBuf = body
	//		return
	//	}
	//
	//	var msg C2S_Reg
	//	err := json.Unmarshal(body, &msg)
	//	if err != nil {
	//		utils.LogError("Reg Message Parse err:", err, string(body))
	//		return
	//	}
	//
	//	self.LoginThirdSdk(session, ctrlhead.Ctrl, msg.Account, msg.Password, msg.ServerId)
	//	return
	//case HEAD_CHAT:
	//	if session.PlayerObj == nil {
	//		return
	//	}
	//	session.PlayerObj.OnReceive(head, "", body)
	//	return
	//case HEAD_CENTER:
	//	self.OnCenterMessage(session, body)
	//	return
	//case HEAD_Reconnect:
	//	self.OnRecconect(session, ctrlhead.Uid, body)
	//	return
	//}

	if ctrlhead.Uid <= 0 {
		return
	}

	//! 如果没有关联的角色，消息不予处理
	if utils.IsNil(session.PlayerObj) {
		return
	}

	//session.PlayerObj.SetSessionID(session, false, session.IP, session.Os)
	cur := time.Now().UnixNano()
	//! 消息处理 暂时屏蔽内容
	//session.PlayerObj.OnReceive(head, ctrlhead.Ctrl, body)

	delay := time.Now().UnixNano() - cur

	if (delay / 1000000) > 10 {
		utils.LogDebug(fmt.Sprintf("消息处理(%s)毫秒(%d)纳秒(%d)", ctrlhead.Ctrl, (delay / 1000000), delay))
	}
}

//! 0-版本OK，1-版本过低，2-app版本过低
func (self *GateApp) CheckVer(ver int) int {
	conf := core.GetMasterApp().GetConfig()

	if ver != 0 && ver < conf.ServerVer {
		//!版本太低，需要更新
		if ver/1000000 != conf.ServerVer/1000000 {
			return 1
			//var msg S2C_Result2Msg
			//msg.Cid = "needupdate"
			//smsg, _ := json.Marshal(&msg)
			//session.SendMsg(msg.Cid, smsg)
		} else if ver%1000000 != conf.ServerVer%1000000 {
			return 2
			//var msg S2C_Result2Msg
			//msg.Cid = "needdownload"
			//smsg, _ := json.Marshal(&msg)
			//self.SendMsg(msg.Cid, smsg)
		}

		utils.LogError("版本过低：", ver)
		return 0
	}

	return 0
}

//! 处理中心服务器消息
func (self GateApp) OnCenterMessage(session *network.Session, body []byte) {
	var centerCid S2S_CenterCid
	err := json.Unmarshal(body, &centerCid)
	if err != nil {
		utils.LogError("CenterCid err:", err)
		return
	}

	//cid := centerCid.Cid
	//switch cid {
	//case "mherotopinfo":
	//	info := GetConsumerTop().DoHeroTopinfo(body)
	//	self.SendMsg("mherotopinfo", info)
	//	break
	//case "mherodamage":
	//	var msg S2Center_UploadHeroDamage
	//	err = json.Unmarshal(body, &msg)
	//	if err == nil {
	//		GetConsumerTopSvr().UploadDamage(msg.Top)
	//	}
	//	break
	//case "mgeneral":
	//	var msg S2Center_UploadGeneral
	//	err = json.Unmarshal(body, &msg)
	//	if err == nil {
	//		GetGeneralMgr().addRank(msg.Top)
	//	}
	//	break
	//case "servgeneralrank":
	//	var msg S2Center_GeneralRank
	//	err = json.Unmarshal(body, &msg)
	//	if err == nil {
	//		info := GetGeneralMgr().DoGetRank(msg.ServerId)
	//		self.SendMsg("servgeneralrank", info)
	//	}
	//	break
	//default:
	//	utils.LogError("unsported ctrl, ctrl:", cid)
	//}
}
