package game

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"

	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"time"
)

type San_Account struct {
	Uid       int64
	Account   string
	UserId    string
	Password  string
	ServerId  int
	Creator   string
	Channelid string
	Time      int64
}

type Platform_Info struct {
	Platform      string //! 平台 android/ios/wp/windows
	DeviceId      string //! 设置标识 安卓取imei ios取idfa
	Brand         string //! 手机品牌
	Model         string //! 手机型号
	UUID          string
	Fr            string
	Res           string
	Net           string
	Mac           string
	Operator      string
	Ip            string
	Ch            string
	SubCh         string
	AccountId     string
	Account_AppId string
}

type Session struct {
	ID        int64           //! 自增长id
	Ws        *websocket.Conn //! websocket
	SendChan  chan []byte     //! 发送消息管道
	RecvChan  chan []byte     //! 发送消息管道
	ShutDown  bool            //! 是否关闭
	ShutTime  int64           //! 关闭时间
	MsgTime   []int64         //! 时间切片
	MsgMaxId  int             //! 当前最大消息Id
	IP        string          //! 公网IP
	PlayerObj *Player         //! 角色对象
	TryNum    int             //! 重试次数
	Account   string          //! 帐号
	Password  string          //! 密码
	Third     string          //! 第三方登录

	Ctrl     string //! 操作
	Os       string //! 操作os
	ServerId int    //! 服务器ID
	Token    string //! ios专有token
	DevInfo  SendRZ_EnvInfo_DevInfo_Ios
	ChInfo   SendRZ_EnvInfo_ChInfo_Ios
	//lockchan *sync.RWMutex
}

//! 发送消息
func (self *Session) SendMsg(head string, body []byte) {
	if head == "" || string(body) == "" {
		return
	}

	//self.lockchan.Lock()
	//defer self.lockchan.Unlock()

	LogDebug("s2c:", head, "....", string(body))

	if self.ShutDown == true || self.SendChan == nil {
		return
	}

	if head == "shutdown" {
		//self.ShutDown = true
		self.ShutTime = TimeServer().Unix()
		self.SendChan <- []byte("")
		return
	}

	if len(self.SendChan) >= sendChanSize-100 {
		LogError("Send Chan Overlap...", self.ID)
		self.SendChan <- []byte("")
		self.ShutDown = true
		self.ShutTime = TimeServer().Unix()
		return
	}

	self.SendChan <- HF_DecodeMsg(head, body)
}

func (self *Session) SendMsgBatch(msg []byte) {
	if self.ShutDown == true {
		return
	}

	if self.SendChan == nil {
		return
	}

	self.SendChan <- msg
}

func (self *Session) CloseChan() {
	//self.lockchan.Lock()
	//defer self.lockchan.Unlock()

	if self.SendChan == nil {
		return
	}

	close(self.SendChan)
	self.SendChan = nil

	close(self.RecvChan)
	self.RecvChan = nil
}

//! 消息run
func (self *Session) Run() {
	go self.sendMsgRun()
	go self.logicRun()
	self.receiveMsgRun()
}

func (self *Session) logicRun() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	ticker := time.NewTicker(time.Second * 1)
	for {
		if GetServer().ShutDown { //! 关服
			break
		}

		if self.ShutDown == true {
			break
		}

		select {
		case packet := <-self.RecvChan:
			self.onReceiveNew(packet)
		case <-ticker.C:
			self.onTimer()
		}
	}

	ticker.Stop()

	self.onClose()
	GetSessionMgr().RemoveSession(self)

}

//! 发送消息循环
func (self *Session) sendMsgRun() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	for msg := range self.SendChan {
		if GetServer().ShutDown { //! 关服
			break
		}

		if string(msg) == "" || self.ShutDown == true {
			break
		}

		exit := false
		for {
			self.Ws.SetWriteDeadline(time.Now().Add(socketTimeOut * time.Second))
			err := websocket.Message.Send(self.Ws, msg)
			if err != nil {
				neterr, ok := err.(net.Error)
				if ok && neterr.Timeout() {
					self.TryNum++
					if self.TryNum >= socketTryNum {
						//LogError("send break...", self.ID, len(self.SendChan), len(msg))
						exit = true
						break
					}
					//LogError("send timeout...", self.ID, len(self.SendChan), len(msg))
					continue
				}
				//LogError("send err", string(msg))
				exit = true
				break
			} else {
				break
			}
		}
		if exit {
			break
		}
	}

	LogInfo("client close send", self.ID)
	self.ShutDown = true
	self.ShutTime = TimeServer().Unix()
	self.Ws.Close()
}

//! 接收消息循环
func (self *Session) receiveMsgRun() {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	for {
		if GetServer().ShutDown { //! 关服
			break
		}

		if self.ShutDown == true {
			break
		}

		var msg []byte
		//var tmp []byte

		self.Ws.SetReadDeadline(time.Now().Add(socketTimeOut * time.Second))
		err := websocket.Message.Receive(self.Ws, &msg)
		if err != nil {
			neterr, ok := err.(net.Error)
			if ok && neterr.Timeout() {
				//LogError("receive timeout")
				continue
			}
			if err == io.EOF {
				LogInfo("client disconnet")
			} else {
				LogInfo("receive err:", err)
			}
			break
		}

		if len(self.MsgTime) >= 30 { //! 取20个消息的间隔如果小于1秒
			if self.MsgTime[29]-self.MsgTime[0] <= 1 {
				break
			}
			self.MsgTime = make([]int64, 0)
		} else {
			self.MsgTime = append(self.MsgTime, TimeServer().Unix())
		}

		//for len(tmp) >= 2048 {
		//	msg = append(msg, tmp...)
		//	self.Ws.SetReadDeadline(TimeServer().Add(1 * time.Second))
		//	err = websocket.Message.Receive(self.Ws, &tmp)
		//	if err != nil {
		//		neterr, ok := err.(net.Error)
		//		if ok && neterr.Timeout() {
		//			continue
		//		}
		//		if err == io.EOF {
		//			LogInfo("client disconnet")
		//		} else {
		//			LogError("receive err:", err)
		//		}
		//		break
		//	}
		//}
		//if len(tmp) > 1 || tmp[0] != '\n' {
		//	msg = append(msg, tmp...)
		//}

		//self.onReceive(msg)
		self.onReceiveNew(msg) //! 新登录模块，兼容渠道SDK
		//self.RecvChan <- msg
	}
	LogInfo("client close recv", self.ID)
	//self.CloseChan()
	self.SendMsg("shutdown", []byte("shutdown"))
	self.ShutDown = true
	self.ShutTime = TimeServer().Unix()
	self.Ws.Close()
	LogInfo("关闭玩家websokcet连接")
	//self.onClose()
	//GetSessionMgr().RemoveSession(self)
}

func (self *Session) onReceive(msg []byte) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	if GetServer().ShutDown {
		return
	}

	head, body, _ := HF_EncodeMsg(msg)
	LogDebug("onReceive c2s:", head, "....", string(body))

	if head == "chat" { //! 聊天单独处理
		if self.PlayerObj == nil {
			return
		}
		//self.PlayerObj.SetSessionID(self, false, self.IP, "")
		self.PlayerObj.onReceive(head, "", body)
		return
	} else if head == "center" {
		self.onCenterMsg(body)
		return
	}

	var ctrlhead C2S_CtrlHead
	err := json.Unmarshal(body, &ctrlhead)
	if err != nil {
		//LogError("CtrlType err:", err)
		return
	}

	//! release 检测消息MsgId
	//if GetServer().Con.NetworkCon.MsgFilter == true {
	//	if ctrlhead.MsgId < 3 {
	//		self.MsgMaxId = ctrlhead.MsgId
	//	} else {
	//		if ctrlhead.MsgId == self.MsgMaxId+1 {
	//			self.MsgMaxId = ctrlhead.MsgId
	//		} else {
	//			return
	//		}
	//	}
	//}
	if GetServer().Con.NetworkCon.MsgFilter == true {
		if ctrlhead.MsgId > self.MsgMaxId {
			self.MsgMaxId = ctrlhead.MsgId
		} else {
			return
		}
	}

	if ctrlhead.Ver != 0 && ctrlhead.Ver < GetServer().Con.ServerVer {
		//!版本太低，需要更新
		if ctrlhead.Ver/1000000 != GetServer().Con.ServerVer/1000000 {
			var msg S2C_Result2Msg
			msg.Cid = "needupdate"
			smsg, _ := json.Marshal(&msg)
			self.SendMsg(msg.Cid, smsg)
		} else if ctrlhead.Ver%1000000 != GetServer().Con.ServerVer%1000000 {
			var msg S2C_Result2Msg
			msg.Cid = "needdownload"
			smsg, _ := json.Marshal(&msg)
			self.SendMsg(msg.Cid, smsg)
		}

		LogError("版本过低：", ctrlhead.Ver)
		return
	}

	if head == "passport.php" {
		switch ctrlhead.Ctrl {
		case "checkver": //! 客户端验证
			return
		case "cancellogin": //!取消排队
			GetLineUpMgr().CancelLogin(self)
			return
		case "login_guest": //! 游客登陆	//
			var c2s_msg C2S_Reg
			json.Unmarshal(body, &c2s_msg)

			LogInfo("正常登录排队：", GetPlayerMgr().GetPlayerOnline(), GetServer().Con.NetworkCon.MaxPlayer,
				c2s_msg.Account, c2s_msg.Password, self.ID)
			if GetPlayerMgr().GetPlayerOnline() >= GetServer().Con.NetworkCon.MaxPlayer {
				self.Ctrl = ctrlhead.Ctrl
				self.Os = ctrlhead.Os
				self.Account = c2s_msg.Account
				self.Password = c2s_msg.Password
				self.ServerId = c2s_msg.ServerId
				GetLineUpMgr().AddWaitClient(self)
				return
			}

			account := self.Reg(c2s_msg.Account, c2s_msg.Password, c2s_msg.ServerId)
			if account == nil {
				return
			}
			if !GetServer().IsWhiteID(account.Uid) && GetServer().IsBlackIP(self.IP) {
				//! 不是白名单id，是黑名单ip
				return
			}

			//GetLineUpMgr().AddLoginClient(self)
			ctrlhead.Uid = account.Uid
			player := GetPlayerMgr().GetPlayer(ctrlhead.Uid, true)
			reLogin := false
			is_creat := false
			if player == nil { //! 找不到新建
				is_creat = true
				player = NewPlayer(ctrlhead.Uid)
				player.InitUserBase(ctrlhead.Uid)
				player.Sql_UserBase.Init("san_userbase", &player.Sql_UserBase, false)
				player.SetAccount(account)
				player.InitPlayerData()
				player.OtherPlayerData()
				player.New()
				player.SaveTime = TimeServer().Unix()
				GetPlayerMgr().AddPlayer(ctrlhead.Uid, player)
				InsertTable("san_userbase", &player.Sql_UserBase, 0, false)
			} else {
				player.OtherPlayerData()

				//! 检测是否有用户连接，踢掉前一个连接并保存数据
				if player.SessionObj != nil && player.SessionObj != self {
					LogDebug("重复登录，踢掉在线用户", "kickout")
					//player.Save(false)
					player.Save(false, true)
					player.SendRet2("kickout")
					player.SendMsg("shutdown", []byte("shutdown"))
					reLogin = true
					//GetPlayerMgr().SetPlayerOffline()
				}
			}

			if is_creat {
				GetServer().sendLog_UserCreate(player)
			}
			GetServer().sendLog_Activation(player)

			if player.Sql_UserBase.IsBlock == 1 {
				LogDebug("冻结帐号，无法登录", "block")
				//player.SendRet2("block")
				var msg S2C_Result2Msg
				msg.Cid = "block"
				self.SendMsg("block", HF_JtoB(&msg))

				self.SendMsg("shutdown", []byte("shutdown"))
				return
			}

			player.Platform.Brand = c2s_msg.Platform_Brand
			player.Platform.DeviceId = c2s_msg.Platform_DeviceId
			player.Platform.Model = c2s_msg.Platform_Model
			player.Platform.Platform = c2s_msg.Platform_os

			LogDebug("检测版本信息：", player.Platform.Platform, c2s_msg.Platform_os)
			player.SetSessionID(self, true, self.IP, c2s_msg.Platform_os)
			if !reLogin {
				GetPlayerMgr().SetPlayerOnline()
			}
			player.SendInfo("userbaseinfo")
			//player.GetPlatfromInfo("getplatforminfo")
			self.PlayerObj = player
			//!刷新好友
			player.GetModule("friend").(*ModFriend).Rename(false)
			GetServer().sendLog_userLoginOK(player)
			return
		case "login_sdk": // 安卓登录
			var msg C2S_SDKLogin
			json.Unmarshal(body, &msg)

			LogInfo("SDK登录排队：", GetPlayerMgr().GetPlayerOnline(), GetServer().Con.NetworkCon.MaxPlayer, "sdk", msg.Password)
			if GetPlayerMgr().GetPlayerOnline() >= GetServer().Con.NetworkCon.MaxPlayer {
				self.Ctrl = ctrlhead.Ctrl
				self.Os = ctrlhead.Os
				self.Account = "sdk"
				self.Password = msg.Password
				self.ServerId = msg.ServerId
				GetLineUpMgr().AddWaitClient(self)
				return
			}
			reLogin := false
			account := self.SDKReg(msg.Password, msg.ServerId)
			if account == nil {
				return
			}
			if !GetServer().IsWhiteID(account.Uid) && GetServer().IsBlackIP(self.IP) {
				//! 不是白名单id，是黑名单ip
				return
			}
			is_creat := false
			ctrlhead.Uid = account.Uid
			player := GetPlayerMgr().GetPlayer(ctrlhead.Uid, true)
			if player == nil { //! 找不到新建
				is_creat = true
				player = NewPlayer(ctrlhead.Uid)
				player.InitUserBase(ctrlhead.Uid)
				player.Sql_UserBase.Init("san_userbase", &player.Sql_UserBase, false)
				player.SaveTime = TimeServer().Unix()
				player.SetAccount(account)
				player.InitPlayerData()
				player.OtherPlayerData()
				player.New()
				GetPlayerMgr().AddPlayer(ctrlhead.Uid, player)
				InsertTable("san_userbase", &player.Sql_UserBase, 0, false)
			} else {
				player.OtherPlayerData()
				if player.Account == nil {
					player.SetAccount(account)
				}
			}

			if is_creat {
				GetServer().sendLog_UserCreate(player)
			}
			GetServer().sendLog_Activation(player)
			//! 检测是否有用户连接，踢掉前一个连接并保存数据
			if player.SessionObj != nil && player.SessionObj != self {
				LogDebug("重复登录，踢掉在线用户", "kickout")
				//player.Save(false)
				player.SendRet2("kickout")
				player.SendMsg("shutdown", []byte("shutdown"))
				//GetPlayerMgr().SetPlayerOffline()
				reLogin = true
				//self.Ws.Close()

			}

			if player.Sql_UserBase.IsBlock == 1 {
				LogDebug("冻结帐号，无法登录", "kickout")
				player.SendRet2("block")
				player.SendMsg("shutdown", []byte("shutdown"))
				return
			}

			//			player.Platform.Brand = c2s_msg.Platform_Brand
			//			player.Platform.DeviceId = c2s_msg.Platform_DeviceId
			//			player.Platform.Model = c2s_msg.Platform_Model
			//			player.Platform.Platform = c2s_msg.Platform_os

			player.SetSessionID(self, true, self.IP, ctrlhead.Os)
			if reLogin == false {
				GetPlayerMgr().SetPlayerOnline()
			}
			player.SendInfo("userbaseinfo")
			self.PlayerObj = player
			//!刷新好友
			player.GetModule("friend").(*ModFriend).Rename(false)
			GetServer().sendLog_userLoginOK(player)
			return
		case "login_ios": // Ios登录
			var msg C2S_SDKIOSLogin
			json.Unmarshal(body, &msg)

			LogInfo("SDK登录排队：", GetPlayerMgr().GetPlayerOnline(), GetServer().Con.NetworkCon.MaxPlayer,
				msg.MemId, msg.AppId, msg.UserToken)
			if GetPlayerMgr().GetPlayerOnline() >= GetServer().Con.NetworkCon.MaxPlayer {
				self.Ctrl = ctrlhead.Ctrl
				self.Os = ctrlhead.Os
				self.Account = msg.AppId
				self.Password = msg.MemId
				self.Token = msg.UserToken
				GetLineUpMgr().AddWaitClient(self)
				return
			}
			reLogin := false
			account := self.SDKRegByIOS(msg.UserToken, msg.AppId, msg.MemId, msg.ServerId)
			if account == nil {
				return
			}
			if !GetServer().IsWhiteID(account.Uid) && GetServer().IsBlackIP(self.IP) {
				//! 不是白名单id，是黑名单ip
				return
			}
			ctrlhead.Uid = account.Uid
			player := GetPlayerMgr().GetPlayer(ctrlhead.Uid, true)
			is_creat := false
			if player == nil { //! 找不到新建
				player = NewPlayer(ctrlhead.Uid)
				player.InitUserBase(ctrlhead.Uid)
				player.Sql_UserBase.Init("san_userbase", &player.Sql_UserBase, false)
				player.SaveTime = TimeServer().Unix()
				player.SetAccount(account)
				player.setPlayrInfo(&msg.DevInfo, &msg.ChInfo, msg.MemId, msg.AppId)
				player.InitPlayerData()
				player.OtherPlayerData()
				player.New()
				is_creat = true
				GetPlayerMgr().AddPlayer(ctrlhead.Uid, player)
				InsertTable("san_userbase", &player.Sql_UserBase, 0, false)
			} else {
				player.setPlayrInfo(&msg.DevInfo, &msg.ChInfo, msg.MemId, msg.AppId)
				player.OtherPlayerData()
				if player.Account == nil {
					player.SetAccount(account)
				}
			}

			if is_creat {
				GetServer().sendLog_UserCreateIOS(player)
			}
			GetServer().sendLog_ActivationIOS(player)

			//! 检测是否有用户连接，踢掉前一个连接并保存数据
			if player.SessionObj != nil && player.SessionObj != self {
				LogDebug("重复登录，踢掉在线用户", "kickout")
				//player.Save(false)
				player.SendRet2("kickout")
				player.SendMsg("shutdown", []byte("shutdown"))
				//GetPlayerMgr().SetPlayerOffline()
				reLogin = true
				//self.Ws.Close()

			}

			if player.Sql_UserBase.IsBlock == 1 {
				LogDebug("冻结帐号，无法登录", "kickout")
				player.SendRet2("block")
				player.SendMsg("shutdown", []byte("shutdown"))
				return
			}

			//			player.Platform.Brand = c2s_msg.Platform_Brand
			//			player.Platform.DeviceId = c2s_msg.Platform_DeviceId
			//			player.Platform.Model = c2s_msg.Platform_Model
			//			player.Platform.Platform = c2s_msg.Platform_os

			player.SetSessionID(self, true, self.IP, ctrlhead.Os)
			if reLogin == false {
				GetPlayerMgr().SetPlayerOnline()
			}

			player.SendInfo("userbaseinfo")
			self.PlayerObj = player
			//!刷新好友
			player.GetModule("friend").(*ModFriend).Rename(false)
			GetServer().sendLog_userLoginOKIOS(player)
			return
		}
	}

	if ctrlhead.Uid <= 0 {
		return
	}

	if !GetServer().IsWhiteID(ctrlhead.Uid) && GetServer().IsBlackIP(self.IP) {
		//! 不是白名单id，是黑名单ip
		return
	}

	//断线重连，计算人数并且同步状态
	if self.PlayerObj == nil {
		player := GetPlayerMgr().GetPlayer(ctrlhead.Uid, false)
		if player == nil {
			var msg S2C_Result2Msg
			msg.Cid = "kickout1"
			smsg, _ := json.Marshal(&msg)
			self.SendMsg("kickout1", smsg)
			return
		} else {
			var msg C2S_AutoLogin
			json.Unmarshal(body, &msg)
			if player.CheckCode == msg.CheckCode {
				if player.SessionObj != nil {
					player.SessionObj.onClose()
				}

				self.PlayerObj = player
				self.PlayerObj.SessionObj = self
				GetPlayerMgr().SetPlayerOnline()
				player.Module.GetOtherData()

				player.CheckRefresh()

				player.SendRet("autologin", LOGIC_TRUE)
			} else {
				var msg S2C_Result2Msg
				msg.Cid = "kickout1"
				smsg, _ := json.Marshal(&msg)
				self.SendMsg("kickout1", smsg)

				return
			}
		}
	} else {
		self.PlayerObj.IsSave = true
	}

	if self.PlayerObj.SessionObj != nil && self.PlayerObj.SessionObj != self { //! 这个人已经被另一个链接上了
		//GetPlayerMgr().SetPlayerOffline(self.PlayerObj)
		LogError("连接错误，这个人已经被另外的连接登录。", self.PlayerObj.Sql_UserBase.Uid, self.ID, self.Account, self.PlayerObj.Sql_UserBase.Uid)
		//self.PlayerObj.SessionObj.onClose()
		//self.onClose()
		return
	}
	self.PlayerObj.SetSessionID(self, false, self.IP, ctrlhead.Os)
	cur := TimeServer().UnixNano()
	self.PlayerObj.onReceive(head, ctrlhead.Ctrl, body)

	delay := TimeServer().UnixNano() - cur

	if (delay / 1000000) > 50 {
		LogDebug(fmt.Sprintf("消息处理(%s)毫秒(%d)纳秒(%d)", ctrlhead.Ctrl, (delay / 1000000), delay))
	}

}

func (self *Session) AutoLogin() {
	//if head == "passport.php" {
	LogInfo("AutoLogin:", self.Account, self.Password)
	var account *San_Account
	switch self.Ctrl {
	case "login_guest": //! 登陆
		player := self.LoginGuest(self.Account, self.Password, self.ServerId)
		if player != nil {
			GetLineUpMgr().RemoveClientSelf(self)
		}

		return
	case "login_sdk":
		//var msg C2S_SDKLogin
		//json.Unmarshal(body, &msg)
		account = self.SDKReg(self.Password, self.ServerId)
		if account == nil {
			return
		}
		if !GetServer().IsWhiteID(account.Uid) && GetServer().IsBlackIP(self.IP) { //! 不是白名单id，是黑名单ip
			return
		}
		//ctrlhead.Uid = account.Uid
		is_creat := false
		player := GetPlayerMgr().GetPlayer(account.Uid, true)
		if player == nil { //! 找不到新建
			is_creat = true
			player = NewPlayer(account.Uid)
			player.InitUserBase(account.Uid)
			player.Sql_UserBase.Init("san_userbase", &player.Sql_UserBase, false)
			player.SaveTime = TimeServer().Unix()
			player.SetAccount(account)
			player.InitPlayerData()
			player.OtherPlayerData()
			player.New()
			GetPlayerMgr().AddPlayer(account.Uid, player)
			InsertTable("san_userbase", &player.Sql_UserBase, 0, false)
		} else {
			player.OtherPlayerData()
			if player.Account == nil {
				player.SetAccount(account)
			}
		}

		if is_creat {
			GetServer().sendLog_UserCreate(player)
		}
		GetServer().sendLog_Activation(player)

		//! 检测是否有用户连接，踢掉前一个连接并保存数据
		if player.SessionObj != nil && player.SessionObj != self {
			LogDebug("重复登录，踢掉在线用户", "kickout")
			//player.Save(false)
			player.SendRet2("kickout")
			player.SendMsg("shutdown", []byte("shutdown"))

		}

		if player.Sql_UserBase.IsBlock == 1 {
			LogDebug("冻结帐号，无法登录", "kickout")
			player.SendRet2("block")
			player.SendMsg("shutdown", []byte("shutdown"))
			return
		}
		player.SetSessionID(self, true, self.IP, self.Os)
		GetPlayerMgr().SetPlayerOnline()
		player.SendInfo("userbaseinfo")
		self.PlayerObj = player
		//!刷新好友
		player.GetModule("friend").(*ModFriend).Rename(false)
		GetServer().sendLog_userLoginOK(player)
		GetLineUpMgr().RemoveClientSelf(self)
		return
	case "login_ios":
		//var msg C2S_SDKLogin
		//json.Unmarshal(body, &msg)
		//account = self.SDKRegByIOS(self.Token, self.Account, self.Password, self.ServerId)

		if player := self.LoginIos(self.Token, self.Account, self.Password, self.ServerId, self.Os); player != nil {
			player.setPlayrInfo(&self.DevInfo, &self.ChInfo, self.Account, self.Password)
		}
		GetLineUpMgr().RemoveClientSelf(self)
		return
	}
	//}

	if account.Uid <= 0 {
		return
	}

	if !GetServer().IsWhiteID(account.Uid) && GetServer().IsBlackIP(self.IP) { //! 不是白名单id，是黑名单ip
		return
	}

	if self.PlayerObj.SessionObj != nil && self.PlayerObj.SessionObj != self { //! 这个人已经被另一个链接上了
		LogError("连接错误，这个人已经被另外的连接登录。", self.PlayerObj.Sql_UserBase.Uid, self.ID, self.Account, self.PlayerObj.Sql_UserBase.Uid)
		//self.onClose()
		//self.PlayerObj.SessionObj.onClose()
		//self.onClose()
		return
	}
	self.PlayerObj.SetSessionID(self, false, self.IP, self.Os)
	//cur := TimeServer().UnixNano()
	//self.PlayerObj.onReceive(head, self.Ctrl, body)
}

func (self *Session) onReceiveNew(msg []byte) {
	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	if GetServer().ShutDown {
		return
	}

	head, body, _ := HF_EncodeMsg(msg)
	LogDebug("onReceiveNew c2s:", head, "....", string(body))

	if head == "chat" { //! 聊天单独处理
		if self.PlayerObj == nil {
			return
		}
		//self.PlayerObj.SetSessionID(self, false, self.IP, "")
		self.PlayerObj.onReceive(head, "", body)
		return
	} else if head == "center" {
		self.onCenterMsg(body)
		return
	}

	var ctrlhead C2S_CtrlHead
	err := json.Unmarshal(body, &ctrlhead)
	if err != nil {
		LogError("CtrlType err:", err, string(body))
		return
	}
	if ctrlhead.Ctrl == "encryption" {
		var msg C2S_Encryption
		json.Unmarshal(body, &msg)
		decode := decryptData(msg.EnInfo, MsgKey)
		var ctrlheadNew C2S_CtrlHead
		err := json.Unmarshal([]byte(decode), &ctrlheadNew)
		if err != nil {
			LogError("CtrlType err:", err, string(body))
			return
		}
		ctrlhead.Ctrl = ctrlheadNew.Ctrl
		body = []byte(decode)
	}

	if GetServer().Con.NetworkCon.MsgFilter == true {
		if ctrlhead.MsgId > self.MsgMaxId {
			self.MsgMaxId = ctrlhead.MsgId
		} else {
			return
		}
	}

	if ctrlhead.Ver != 0 && ctrlhead.Ver < GetServer().Con.ServerVer {
		//!版本太低，需要更新
		if ctrlhead.Ver/1000000 != GetServer().Con.ServerVer/1000000 {
			var msg S2C_Result2Msg
			msg.Cid = "needupdate"
			smsg, _ := json.Marshal(&msg)
			self.SendMsg(msg.Cid, smsg)
		} else if ctrlhead.Ver%1000000 != GetServer().Con.ServerVer%1000000 {
			var msg S2C_Result2Msg
			msg.Cid = "needdownload"
			smsg, _ := json.Marshal(&msg)
			self.SendMsg(msg.Cid, smsg)
		}

		LogError("版本过低：", ctrlhead.Ver)
		return
	}
	if head == "passport.php" {
		switch ctrlhead.Ctrl {
		case "checkver": //! 客户端验证
			return
		case "cancellogin": //!取消排队
			GetLineUpMgr().CancelLogin(self)
			return
		case "login_guest": //! 游客登陆	//
			var c2s_msg C2S_Reg
			json.Unmarshal(body, &c2s_msg)

			LogInfo("正常登录排队：", GetPlayerMgr().GetPlayerOnline(), GetServer().Con.NetworkCon.MaxPlayer,
				c2s_msg.Account, c2s_msg.Password, self.ID)
			if GetPlayerMgr().GetPlayerOnline() >= GetServer().Con.NetworkCon.MaxPlayer {
				self.Ctrl = ctrlhead.Ctrl
				self.Os = ctrlhead.Os
				self.Account = c2s_msg.Account
				self.Password = c2s_msg.Password
				self.ServerId = c2s_msg.ServerId
				GetLineUpMgr().AddWaitClient(self)
				return
			}

			player := self.LoginGuest(c2s_msg.Account, c2s_msg.Password, c2s_msg.ServerId)
			if player != nil {
				player.Platform.Brand = c2s_msg.Platform_Brand
				player.Platform.DeviceId = c2s_msg.Platform_DeviceId
				player.Platform.Model = c2s_msg.Platform_Model
				player.Platform.Platform = c2s_msg.Platform_os

				GetServer().SqlLog(player.Sql_UserBase.Uid, LOG_USER_LOGIN, 0, 0, 0, "点击登录", 0, 0, player)
			}

			return
		case "login_sdk": // 安卓登录
			var msg C2S_SDKLogin
			json.Unmarshal(body, &msg)

			LogInfo("SDK登录排队：", GetPlayerMgr().GetPlayerOnline(), GetServer().Con.NetworkCon.MaxPlayer, "sdk", msg.Password)
			if GetPlayerMgr().GetPlayerOnline() >= GetServer().Con.NetworkCon.MaxPlayer {
				self.Ctrl = ctrlhead.Ctrl
				self.Os = ctrlhead.Os
				self.Account = msg.Username
				self.Password = msg.Password
				self.ServerId = msg.ServerId
				GetLineUpMgr().AddWaitClient(self)
				return
			}

			player := self.LoginSdk(ctrlhead.Ctrl, msg.Password, msg.ServerId, msg.Username)
			if player != nil {
				GetServer().SqlLog(player.Sql_UserBase.Uid, LOG_USER_LOGIN, 1, 0, 0, "点击登录", 0, 0, player)
			}
			return
		case "login_sdk_myx":
			var msg C2S_SDKLogin
			json.Unmarshal(body, &msg)

			LogInfo("锦游登录排队：", GetPlayerMgr().GetPlayerOnline(), GetServer().Con.NetworkCon.MaxPlayer, "sdk", msg.Password)
			if GetPlayerMgr().GetPlayerOnline() >= GetServer().Con.NetworkCon.MaxPlayer {
				self.Ctrl = ctrlhead.Ctrl
				self.Os = ctrlhead.Os
				self.Account = msg.Username
				self.Password = msg.Password
				self.ServerId = msg.ServerId
				GetLineUpMgr().AddWaitClient(self)
				return
			}

			self.LoginSdk(ctrlhead.Ctrl, msg.Password, msg.ServerId, msg.Username)

			return
		case "login_sdk_third":
			var msg C2S_SDKLogin
			json.Unmarshal(body, &msg)

			LogInfo("第三方登录排队：", GetPlayerMgr().GetPlayerOnline(), GetServer().Con.NetworkCon.MaxPlayer, "third", msg.Password)
			if GetPlayerMgr().GetPlayerOnline() >= GetServer().Con.NetworkCon.MaxPlayer {
				self.Ctrl = ctrlhead.Ctrl
				self.Os = ctrlhead.Os
				self.Account = msg.Username
				self.Password = msg.Password
				self.ServerId = msg.ServerId
				self.Third = msg.Third
				GetLineUpMgr().AddWaitClient(self)
				return
			}

			player := self.LoginThirdSdk(msg.Third, msg.Password, msg.ServerId, msg.Username)
			if nil != player {
				GetServer().SqlLog(player.Sql_UserBase.Uid, LOG_USER_LOGIN, 1, 0, 0, "点击登录", 0, 0, player)
			}

			return
		case "login_ios": // Ios登录
			var msg C2S_SDKIOSLogin
			json.Unmarshal(body, &msg)

			LogInfo("SDK登录排队：", GetPlayerMgr().GetPlayerOnline(), GetServer().Con.NetworkCon.MaxPlayer,
				msg.MemId, msg.AppId, msg.UserToken)
			if GetPlayerMgr().GetPlayerOnline() >= GetServer().Con.NetworkCon.MaxPlayer {
				self.Ctrl = ctrlhead.Ctrl
				self.Os = ctrlhead.Os
				self.Account = msg.AppId
				self.Password = msg.MemId
				self.Token = msg.UserToken
				GetLineUpMgr().AddWaitClient(self)
				return
			}
			player := self.LoginIos(msg.UserToken, msg.AppId, msg.MemId, msg.ServerId, ctrlhead.Os)
			if player != nil {
				player.setPlayrInfo(&msg.DevInfo, &msg.ChInfo, msg.MemId, msg.AppId)
			}

			return
		}
	}

	if ctrlhead.Uid <= 0 {
		return
	}

	if !GetServer().IsWhiteID(ctrlhead.Uid) && GetServer().IsBlackIP(self.IP) {
		//! 不是白名单id，是黑名单ip
		return
	}

	//断线重连，计算人数并且同步状态
	if self.PlayerObj == nil {
		player := GetPlayerMgr().GetPlayer(ctrlhead.Uid, false)
		if player == nil {
			var msg S2C_Result2Msg
			msg.Cid = "kickout1"
			smsg, _ := json.Marshal(&msg)
			self.SendMsg("kickout1", smsg)
			return
		} else {
			var msg C2S_AutoLogin
			json.Unmarshal(body, &msg)
			if player.CheckCode == msg.CheckCode {
				if player.SessionObj != nil {
					player.SessionObj.onClose()
				}

				self.PlayerObj = player
				self.PlayerObj.SessionObj = self
				GetPlayerMgr().SetPlayerOnline()
				player.Module.GetOtherData()

				player.CheckRefresh()
			} else {
				var msg S2C_Result2Msg
				msg.Cid = "kickout1"
				smsg, _ := json.Marshal(&msg)
				self.SendMsg("kickout1", smsg)

				return
			}
		}
	} else {
		self.PlayerObj.IsSave = true
	}

	if self.PlayerObj.SessionObj != nil && self.PlayerObj.SessionObj != self { //! 这个人已经被另一个链接上了
		//GetPlayerMgr().SetPlayerOffline(self.PlayerObj)
		LogError("连接错误，这个人已经被另外的连接登录。", self.PlayerObj.Sql_UserBase.Uid, self.ID, self.Account, self.PlayerObj.Sql_UserBase.Uid)
		//self.PlayerObj.SessionObj.onClose()
		//self.onClose()
		return
	}
	self.PlayerObj.SetSessionID(self, false, self.IP, ctrlhead.Os)
	cur := TimeServer().UnixNano()
	self.PlayerObj.onReceive(head, ctrlhead.Ctrl, body)
	delay := TimeServer().UnixNano() - cur

	if (delay / 1000000) > 10 {
		LogDebug(fmt.Sprintf("消息处理(%s)毫秒(%d)纳秒(%d)", ctrlhead.Ctrl, (delay / 1000000), delay))
	}

}

func (self *Session) LoginThirdSdk(ctrl string, password string, serverid int, username string) *Player {
	var account *San_Account
	switch ctrl {
	case "login_sdk":
		account = self.SDKReg(password, serverid)
	case "login_sdk_myx":
		account = self.SDKReg_MYX(password, serverid, username, "sdk_myx_android")
	case "login_sdk_myx_ios":
		account = self.SDKReg_MYX(password, serverid, username, "sdk_myx_ios")
	case "sdk_mzy":
		account = self.SDKReg_MZY(password, serverid, username, ctrl)
	case "sdk_mzy_ios":
		account = self.SDKReg_MZY(password, serverid, username, ctrl)
	case "sdk_jinmu":
		account = self.SDKReg_Jinmu(username, serverid, password, ctrl)
	case "sdk_jinmu_ios":
		account = self.SDKReg_Jinmu(username, serverid, password, ctrl)
	case "sdk_shuguo":
		account = self.SDKReg_Shuguo(username, serverid, password, ctrl)
	case "sdk_shuguo_ios":
		account = self.SDKReg_Shuguo(username, serverid, password, ctrl)
	case "sdk_zhish":
		account = self.SDKReg_ZhiSh(username, serverid, password, ctrl)
	case "sdk_zhish_ios":
		account = self.SDKReg_ZhiSh(username, serverid, password, ctrl)
	case "sdk_9377":
		account = self.SDKReg_9377(username, serverid, password, ctrl)
	case "sdk_yunke":
		account = self.SDKReg_YunKe(username, serverid, password, ctrl)
	case "sdk_ingcle":
		account = self.SDKReg_Ingcle(username, serverid, password, ctrl)
	case "sdk_yunke_ios":
		account = self.SDKReg_YunKe_IOS(username, serverid, password, ctrl)
	case "sdk_koyou":
		account = self.SDKReg_KoYou(username, serverid, password, ctrl)
	case "sdk_koyou_ios":
		account = self.SDKReg_KoYou(username, serverid, password, ctrl)
	case "sdk_huixuan_ios", "sdk_huixuan_new_ios":
		account = self.SDKReg_HuiXuan(password, serverid, username, ctrl)
	case "sdk_huixuan":
		account = self.SDKReg_HuiXuan(password, serverid, username, ctrl)
	case "sdk_common", "sdk_common_ios":
		account = self.SDKReg_Common(password, serverid, username, ctrl)
	case "sdk_youxifan":
		account = self.SDKReg_YouXiFan(password, serverid, username, ctrl)
	}
	//account := self.SDKReg(password, serverid)
	if account == nil {
		return nil
	}
	reLogin := false
	//account := self.SDKReg(msg.Password, msg.ServerId)
	if account == nil {
		return nil
	}
	if !GetServer().IsWhiteID(account.Uid) && GetServer().IsBlackIP(self.IP) {
		//! 不是白名单id，是黑名单ip
		return nil
	}
	is_creat := false
	player := GetPlayerMgr().GetPlayer(account.Uid, true)
	if player == nil { //! 找不到新建
		is_creat = true
		player = NewPlayer(account.Uid)
		player.InitUserBase(account.Uid)
		player.Sql_UserBase.Init("san_userbase", &player.Sql_UserBase, false)
		player.SaveTime = TimeServer().Unix()
		player.SetAccount(account)
		player.InitPlayerData()
		player.OtherPlayerData()
		player.New()
		GetPlayerMgr().AddPlayer(account.Uid, player)
		InsertTable("san_userbase", &player.Sql_UserBase, 0, false)
	} else {
		player.OtherPlayerData()
		if player.Account == nil {
			player.SetAccount(account)
		}
	}

	if is_creat {
		GetServer().sendLog_UserCreate(player)
	}
	GetServer().sendLog_Activation(player)
	//! 检测是否有用户连接，踢掉前一个连接并保存数据
	if player.SessionObj != nil && player.SessionObj != self {
		LogInfo("重复登录，踢掉在线用户", "kickout")
		//player.Save(false)
		player.CheckCode = GetRandomString(8)
		player.SendRet2("kickout")
		player.SendMsg("shutdown", []byte("shutdown"))
		//GetPlayerMgr().SetPlayerOffline()
		player.Save(false, true)
		reLogin = true
		//self.Ws.Close()

		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_THE_ACCOUNT_IS_ALREADY_ONLINE"))
		//self.SendReturn("kickout")
		var msgRel S2C_OnlineTip
		msgRel.Cid = "alreadyonline"
		self.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))
		self.SendReturn("shutdown")
		player.ReloginTimes++
		if player.ReloginTimes > 2 {
			player.SessionObj = nil
		}
		return nil
	}

	blockflag, blocktime := self.CheckBlock(player)
	if blockflag == true {
		LogDebug("冻结帐号，无法登录", "block", blocktime)
		//player.SendRet("block", blocktime)
		//player.SendMsg("shutdown", []byte("shutdown"))
		var msg S2C_ResultMsg
		msg.Cid = "block"
		msg.Ret = blocktime
		self.SendMsg("block", HF_JtoB(&msg))

		self.SendMsg("shutdown", []byte("shutdown"))
		return nil
	}

	player.SetSessionID(self, true, self.IP, self.Os)
	if reLogin == false {
		GetPlayerMgr().SetPlayerOnline()
	}
	player.SendInfo("userbaseinfo")
	self.PlayerObj = player
	//!刷新好友
	player.GetModule("friend").(*ModFriend).Rename(false)
	GetServer().sendLog_userLoginOK(player)
	return player
}

func (self *Session) LoginSdk(ctrl string, password string, serverid int, username string) *Player {
	var account *San_Account
	switch ctrl {
	case "login_sdk":
		account = self.SDKReg(password, serverid)
	case "login_sdk_myx":
		account = self.SDKReg_MYX(password, serverid, username, "sdk_myx")
	case "login_sdk_mzy":
		account = self.SDKReg_MZY(password, serverid, username, "sdk_mzy")
	}
	//account := self.SDKReg(password, serverid)
	if account == nil {
		return nil
	}
	reLogin := false
	//account := self.SDKReg(msg.Password, msg.ServerId)
	if account == nil {
		return nil
	}
	if !GetServer().IsWhiteID(account.Uid) && GetServer().IsBlackIP(self.IP) {
		//! 不是白名单id，是黑名单ip
		return nil
	}
	is_creat := false
	player := GetPlayerMgr().GetPlayer(account.Uid, true)
	if player == nil { //! 找不到新建
		is_creat = true
		player = NewPlayer(account.Uid)
		player.InitUserBase(account.Uid)
		player.Sql_UserBase.Init("san_userbase", &player.Sql_UserBase, false)
		player.SaveTime = TimeServer().Unix()
		player.SetAccount(account)
		player.InitPlayerData()
		player.OtherPlayerData()
		player.New()
		GetPlayerMgr().AddPlayer(account.Uid, player)
		InsertTable("san_userbase", &player.Sql_UserBase, 0, false)
	} else {
		player.OtherPlayerData()
		if player.Account == nil {
			player.SetAccount(account)
		}
	}

	if is_creat {
		GetServer().sendLog_UserCreate(player)
	}
	GetServer().sendLog_Activation(player)
	//! 检测是否有用户连接，踢掉前一个连接并保存数据
	if player.SessionObj != nil && player.SessionObj != self {
		LogInfo("重复登录，踢掉在线用户", "kickout")
		//player.Save(false)
		player.CheckCode = GetRandomString(8)
		player.SendRet2("kickout")
		player.SendMsg("shutdown", []byte("shutdown"))
		//GetPlayerMgr().SetPlayerOffline()
		player.Save(false, true)
		reLogin = true
		//self.Ws.Close()

		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_THE_ACCOUNT_IS_ALREADY_ONLINE"))
		//self.SendReturn("kickout")
		self.SendReturn("shutdown")
		player.ReloginTimes++
		if player.ReloginTimes > 2 {
			player.SessionObj = nil
		}
		return nil
	}

	//if player.Sql_UserBase.IsBlock == 1 {
	//	LogDebug("冻结帐号，无法登录", "kickout")
	//	player.SendRet2("block")
	//	player.SendMsg("shutdown", []byte(""))
	//	return nil
	//}

	blockflag, blocktime := self.CheckBlock(player)
	if blockflag == true {
		LogDebug("冻结帐号，无法登录", "block", blocktime)
		//player.SendRet("block", blocktime)
		//player.SendMsg("shutdown", []byte(""))
		var msg S2C_ResultMsg
		msg.Cid = "block"
		msg.Ret = blocktime
		self.SendMsg("block", HF_JtoB(&msg))

		self.SendMsg("shutdown", []byte("shutdown"))
		return nil
	}

	player.SetSessionID(self, true, self.IP, self.Os)
	if reLogin == false {
		GetPlayerMgr().SetPlayerOnline()
	}
	player.SendInfo("userbaseinfo")
	self.PlayerObj = player
	//!刷新好友
	player.GetModule("friend").(*ModFriend).Rename(false)
	GetServer().sendLog_userLoginOK(player)
	return player
}

func (self *Session) LoginGuest(user string, password string, serverid int) *Player {
	//var c2s_msg C2S_Reg
	//json.Unmarshal(body, &c2s_msg)

	LogInfo("正常登录排队：", GetPlayerMgr().GetPlayerOnline(), GetServer().Con.NetworkCon.MaxPlayer,
		user, password, self.ID)
	if GetPlayerMgr().GetPlayerOnline() >= GetServer().Con.NetworkCon.MaxPlayer {
		//self.Ctrl = ctrlhead.Ctrl
		//self.Os = ctrlhead.Os
		//self.Account = c2s_msg.Account
		//self.Password = c2s_msg.Password
		//self.ServerId = c2s_msg.ServerId
		//GetLineUpMgr().AddWaitClient(self)
		return nil
	}

	account := self.Reg(user, password, serverid)
	if account == nil {
		return nil
	}
	if !GetServer().IsWhiteID(account.Uid) && GetServer().IsBlackIP(self.IP) {
		//! 不是白名单id，是黑名单ip
		return nil
	}

	//GetLineUpMgr().AddLoginClient(self)
	Uid := account.Uid
	player := GetPlayerMgr().GetPlayer(Uid, true)
	reLogin := false
	is_creat := false
	if player == nil { //! 找不到新建
		is_creat = true
		player = NewPlayer(account.Uid)
		player.InitUserBase(account.Uid)
		player.Sql_UserBase.Init("san_userbase", &player.Sql_UserBase, false)
		player.SetAccount(account)
		player.InitPlayerData()
		player.OtherPlayerData()
		player.New()
		player.SaveTime = TimeServer().Unix()
		GetPlayerMgr().AddPlayer(account.Uid, player)
		InsertTable("san_userbase", &player.Sql_UserBase, 0, false)
	} else {
		player.OtherPlayerData()

		//! 检测是否有用户连接，踢掉前一个连接并保存数据
		if player.SessionObj != nil && player.SessionObj != self {
			LogInfo("重复登录，踢掉在线用户", "kickout")
			//player.Save(false)
			player.Save(false, true)
			player.CheckCode = GetRandomString(8)
			player.SendRet2("kickout")
			player.SendMsg("shutdown", []byte("shutdown"))

			//var msg S2C_OnlineTip
			//msg.Cid = "alreadyonline"
			//msg.Param = LOGIC_FALSE
			//player.SendMsg(msg.Cid, HF_JtoB(&msg))

			reLogin = true
			//GetPlayerMgr().SetPlayerOffline()

			self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_YOUR_ACCOUNT_IS_ALREADY_ONLINE"))

			var msgRel S2C_OnlineTip
			msgRel.Cid = "alreadyonline"
			self.SendMsg(msgRel.Cid, HF_JtoB(&msgRel))

			//self.SendReturn("kickout")
			self.SendReturn("shutdown")
			player.ReloginTimes++
			if player.ReloginTimes > 2 {
				player.SessionObj = nil
			}
			return nil
		}
	}

	if is_creat {
		GetServer().sendLog_UserCreate(player)
	}
	GetServer().sendLog_Activation(player)

	blockflag, blocktime := self.CheckBlock(player)
	if blockflag == true {
		LogDebug("冻结帐号，无法登录", "block", blocktime)
		//player.SendRet("block", blocktime)
		//player.SendMsg("shutdown", []byte(""))
		var msg S2C_ResultMsg
		msg.Cid = "block"
		msg.Ret = blocktime
		self.SendMsg("block", HF_JtoB(&msg))

		self.SendMsg("shutdown", []byte("shutdown"))
		return nil
	}

	player.SetSessionID(self, true, self.IP, "windows")
	if !reLogin {
		GetPlayerMgr().SetPlayerOnline()
	}

	player.SendInfo("userbaseinfo")
	self.PlayerObj = player
	//!刷新好友
	player.GetModule("friend").(*ModFriend).Rename(false)
	GetServer().sendLog_userLoginOK(player)

	return player
}

func (self *Session) LoginIos(usertoken string, appid string, memid string, serverid int, os string) *Player {
	reLogin := false
	account := self.SDKRegByIOS(usertoken, appid, memid, serverid)
	if account == nil {
		return nil
	}
	if !GetServer().IsWhiteID(account.Uid) && GetServer().IsBlackIP(self.IP) {
		//! 不是白名单id，是黑名单ip
		return nil
	}
	//ctrlhead.Uid = account.Uid
	player := GetPlayerMgr().GetPlayer(account.Uid, true)
	is_creat := false
	if player == nil { //! 找不到新建
		player = NewPlayer(account.Uid)
		player.InitUserBase(account.Uid)
		player.Sql_UserBase.Init("san_userbase", &player.Sql_UserBase, false)
		player.SaveTime = TimeServer().Unix()
		player.SetAccount(account)
		player.InitPlayerData()
		player.OtherPlayerData()
		player.New()
		is_creat = true
		GetPlayerMgr().AddPlayer(account.Uid, player)
		InsertTable("san_userbase", &player.Sql_UserBase, 0, false)
	} else {
		player.OtherPlayerData()
		if player.Account == nil {
			player.SetAccount(account)
		}
	}

	if is_creat {
		GetServer().sendLog_UserCreateIOS(player)
	}
	GetServer().sendLog_ActivationIOS(player)

	//! 检测是否有用户连接，踢掉前一个连接并保存数据
	if player.SessionObj != nil && player.SessionObj != self {
		LogInfo("重复登录，踢掉在线用户", "kickout")
		player.CheckCode = GetRandomString(8)
		player.SendRet2("kickout")
		player.SendMsg("shutdown", []byte("shutdown"))
		player.Save(false, true)
		reLogin = true

		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_YOUR_ACCOUNT_IS_ALREADY_ONLINE"))
		self.SendReturn("shutdown")
		player.ReloginTimes++
		if player.ReloginTimes > 2 {
			player.SessionObj = nil
		}

		return nil
	}

	blockflag, blocktime := self.CheckBlock(player)
	if blockflag == true {
		LogDebug("冻结帐号，无法登录", "block", blocktime)
		var msg S2C_ResultMsg
		msg.Cid = "block"
		msg.Ret = blocktime
		self.SendMsg("block", HF_JtoB(&msg))

		self.SendMsg("shutdown", []byte("shutdown"))
		return nil
	}

	player.SetSessionID(self, true, self.IP, os)
	if reLogin == false {
		GetPlayerMgr().SetPlayerOnline()
	}

	player.SendInfo("userbaseinfo")
	self.PlayerObj = player
	//!刷新好友
	player.GetModule("friend").(*ModFriend).Rename(false)
	GetServer().sendLog_userLoginOKIOS(player)

	return player
}

func (self *Session) onClose() {
	if self.PlayerObj != nil && self.PlayerObj.SessionObj == self {
		LogDebug("玩家离线，同步离线数据", self.ID, self.PlayerObj.Sql_UserBase.Uid)
		GetPlayerMgr().SetPlayerOffline()
		LogDebug("减少玩家离线数量")
		self.PlayerObj.onClose()
		self.PlayerObj = nil
	}
	GetLineUpMgr().RemoveClient(self)
	LogDebug("删除玩家连接信息")
}

// true: 表示冻结
func (self *Session) CheckBlock(player *Player) (bool, int) {
	// 合服的时候不能让玩家登陆
	if GetBackStageMgr().getMerge() == 1 {
		player.SafeClose()
		return true, 0
	}

	if player.Sql_UserBase.IsBlock == 1 {
		if player.Sql_UserBase.BlockTime == 0 {
			return true, 0
		}

		if TimeServer().Unix() > player.Sql_UserBase.BlockTime+int64(player.Sql_UserBase.BlockDay) {
			player.Sql_UserBase.IsBlock = 0

			return false, 0
		}

		return true, int(int64(player.Sql_UserBase.BlockDay*DAY_SECS) - player.Sql_UserBase.BlockTime)
	}

	return false, 0
}

func (self *Session) SendErrInfo(cid string, info string) {
	var msg S2C_ErrInfo
	msg.Cid = cid
	msg.Info = info
	smsg, _ := json.Marshal(&msg)
	self.SendMsg("1", smsg)
}

func (self *Session) SendReturn(cid string) {
	var msg S2C_Result2Msg
	msg.Cid = cid
	smsg, _ := json.Marshal(&msg)
	self.SendMsg("1", smsg)
}

func (self *Session) Reg(account string, password string, serverid int) *San_Account {
	if account == "" { //! 游客登录
		b := make([]byte, 48)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return nil
		}
		h := md5.New()
		h.Write([]byte(base64.URLEncoding.EncodeToString(b)))
		account = hex.EncodeToString(h.Sum(nil))
	}

	var _account San_Account
	sql := fmt.Sprintf("select * from `san_account` where `account` = '%s' and serverid=%d", account, serverid)
	GetServer().DBUser.GetOneData(sql, &_account, "", 0)
	//c := GetServer().GetRedisConn()
	//defer c.Close()

	//value, _ := redis.Bytes(c.Do("GET", fmt.Sprintf("%s_%s", "san_account", account)))
	//json.Unmarshal(value, &_account)

	if _account.Uid > 0 {
		if password != _account.Password { //! 密码错误
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_PASSWORD"))
			return nil
		}
	} else { //! 插入
		_account.Account = account
		_account.Password = password
		_account.ServerId = serverid

		_account.Creator = "guest"
		_account.Time = TimeServer().Unix()
		_account.Uid = InsertTable("san_account", &_account, 1, false)
		//_account.Uid = GetServer().GetRedisInc("san_account")
		if _account.Uid <= 0 {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return nil
		}

		//! 重新获取帐号-回避insertid出错
		sql := fmt.Sprintf("select * from `san_account` where `account` = '%s' and serverid=%d", account, serverid)
		GetServer().DBUser.GetOneData(sql, &_account, "", 0)
		//value := HF_JtoB(&_account)
		//c := GetServer().GetRedisConn()
		//defer c.Close()
		//_, err := c.Do("MSET", fmt.Sprintf("%s_%d", "san_account", _account.Uid), value,
		//	fmt.Sprintf("%s_%s", "san_account", _account.Account), value)
		//if err != nil {
		//	LogError("redis set err:", "san_account", ",", string(value))
		//	self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		//	return nil
		//}

		var msg S2C_Reg
		msg.Cid = "reg"
		msg.Uid = _account.Uid
		msg.Account = _account.Account
		msg.Password = _account.Password
		msg.Creator = _account.Creator
		self.SendMsg("1", HF_JtoB(&msg))
	}

	return &_account
}

type JS_SDKData struct {
	Token string `json:"token"`
}

type JS_SDKGame struct {
	Id string `json:"id"`
}

type JS_SDKState struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type JS_SDKHData struct {
	UserId    string `json:"userId"`
	Creator   string `json:"creator"`
	ChannelId string `json:"channelid"`
}

type JS_SDKBody struct {
	Id    string      `json:"id"`
	State JS_SDKState `json:"state"`
	Data  JS_SDKHData `json:"data"`
}

type JS_SDKLogin struct {
	Id   int64      `json:"id"`
	Game JS_SDKGame `json:"game"`
	Data JS_SDKData `json:"data"`
	Sign string     `json:"sign"`
}

type JS_SDKIOSLogin struct {
	AppId     int    `json:"app_id"`
	MemId     int    `json:"mem_id"`
	UserToken string `json:"user_token"`
	Sign      string `json:"sign"`
}

type JS_SDKIOSBody struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

// 安卓登录-火速登录
func (self *Session) SDKReg(token string, serverid int) *San_Account {
	log.Println(token)

	str := "token=" + token + GetServer().Con.AppKey
	h := md5.New()
	h.Write([]byte(str))
	md5str := fmt.Sprintf("%x", h.Sum(nil))
	log.Println(md5str)

	var m JS_SDKLogin
	m.Id = TimeServer().Unix()
	m.Game.Id = GetServer().Con.GameID
	m.Data.Token = token
	m.Sign = md5str
	body := bytes.NewBuffer(HF_JtoB(&m))
	url := "http://account.flysdk.cn/gs/account.verifyToken?ver=1.0&df=json"
	res, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_APPLICATION_ERROR"))
		return nil
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_APPLICATION_ERROR"))
		return nil
	}
	log.Println("result=", string(result))

	var ret JS_SDKBody
	err = json.Unmarshal(result, &ret)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_APPLICATION_ERROR"))
		return nil
	}

	switch ret.State.Code {
	case 4000000:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_REQUEST_PARAMETER_ERROR"))
		return nil
	case 4000001:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_BUSINESS_PARAMETER_ERROR"))
		return nil
	case 5000000:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_THE_NETWORK_IS_BUSY_PLEASE"))
		return nil
	case 5000003:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_THE_SYSTEM_IS_BUSY_PLEASE"))
		return nil
	case 4001001:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_INVALID_SACRED_MARK"))
		return nil
	case 4001003:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_THE_GREAT_SAINT_LOGO_HAS"))
		return nil
	}

	account := ret.Data.UserId
	password := "CYCYCY"

	if account == "" { //! 游客登录
		b := make([]byte, 48)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return nil
		}
		h := md5.New()
		h.Write([]byte(base64.URLEncoding.EncodeToString(b)))
		account = hex.EncodeToString(h.Sum(nil))
	}

	var _account San_Account
	sql := fmt.Sprintf("select * from `san_account` where `account` = '%s' and serverid =%d", account, serverid)
	GetServer().DBUser.GetOneData(sql, &_account, "", 0)
	//c := GetServer().GetRedisConn()
	//defer c.Close()

	//value, _ := redis.Bytes(c.Do("GET", fmt.Sprintf("%s_%s", "san_account", account)))
	//json.Unmarshal(value, &_account)

	if _account.Uid > 0 {
		if password != _account.Password { //! 密码错误
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_PASSWORD"))
			return nil
		}
	} else { //! 插入
		_account.Account = account
		_account.Password = password
		_account.Creator = ret.Data.Creator
		_account.Channelid = ret.Data.ChannelId
		_account.ServerId = serverid
		_account.Time = TimeServer().Unix()
		_account.Uid = InsertTable("san_account", &_account, 1, false)
		//_account.Uid = GetServer().GetRedisInc("san_account")
		if _account.Uid <= 0 {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return nil
		}

		//! 重新获取帐号-回避insertid出错
		sql := fmt.Sprintf("select * from `san_account` where `account` = '%s' and serverid=%d", account, serverid)
		GetServer().DBUser.GetOneData(sql, &_account, "", 0)
		//value := HF_JtoB(&_account)
		//c := GetServer().GetRedisConn()
		//defer c.Close()
		//_, err := c.Do("MSET", fmt.Sprintf("%s_%d", "san_account", _account.Uid), value,
		//	fmt.Sprintf("%s_%s", "san_account", _account.Account), value)
		//if err != nil {
		//	LogError("redis set err:", "san_account", ",", string(value))
		//	self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		//	return nil
		//}

		var msg S2C_Reg
		msg.Cid = "reg"
		msg.Uid = _account.Uid
		msg.Account = _account.Account
		msg.Password = _account.Password
		msg.Creator = _account.Creator
		self.SendMsg("1", HF_JtoB(&msg))
	}

	return &_account
}

//! ios登录
func (self *Session) SDKRegByIOS(token string, appid string, memid string, serverid int) *San_Account {
	log.Println(token)

	str := "app_id=" + appid + "&mem_id=" + memid + "&user_token=" + token + "&app_key=" + GetServer().Con.GetAppKeyByAppId(appid)
	h := md5.New()
	h.Write([]byte(str))
	md5str := fmt.Sprintf("%x", h.Sum(nil))
	log.Println("sign = ", md5str)

	str1 := "app_id=" + appid + "&mem_id=" + memid + "&app_key=" + GetServer().Con.GetAppKeyByAppId(appid)
	h1 := md5.New()
	h1.Write([]byte(str1))
	md5str1 := fmt.Sprintf("%x", h1.Sum(nil))
	log.Println("sign1 = ", md5str1)

	//var m JS_SDKIOSLogin
	//m.AppId = HF_Atoi(appid)
	//m.MemId = memid
	//m.UserToken = token
	//m.Sign = md5str
	//body := bytes.NewBuffer(HF_JtoB(&m))
	check := "app_id=" + appid + "&mem_id=" + memid + "&user_token=" + token + "&sign=" + md5str

	//body := bytes.NewBuffer([]byte(""))
	log.Println("check = ", check)
	//url := "https://aliapi.1tsdk.com/api/v7/cp/user/check?" + check
	//res, err := http.Get(url)
	//if err != nil {
	//	log.Println(err)
	//	self.SendErrInfo("err", "应用错误")
	//	return nil
	//}

	url := GetServer().Con.GetCheckUrlByAppId(appid) + check
	body := bytes.NewBuffer([]byte(""))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := TimeServer().Add(6 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*1)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_APPLICATION_ERROR"))
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}

	if resp.StatusCode != 200 {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_APPLICATION_ERROR"))
		return nil
	}

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_APPLICATION_ERROR"))
		return nil
	}
	log.Println("result=", string(result))

	var ret JS_SDKIOSBody
	err = json.Unmarshal(result, &ret)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_APPLICATION_ERROR"))
		return nil
	}

	switch ret.Status {
	case "0":
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_REQUEST_PARAMETER_ERROR"))
		return nil
	case "10":
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_SERVER_INTERNAL_ERROR"))
		return nil
	case "11":
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_SERIAL_NUMBER_ERROR"))
		return nil
	case "12":
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_SIGNATURE_ERROR"))
		return nil
	case "13":
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_PLAYER_IDENTIFICATION_ERROR"))
		return nil
	case "14":
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_PLAYER_ID_TIMEOUT_INDICATES_THAT"))
		return nil
	}

	account := md5str
	account1 := md5str1
	password := "CYCYCY"

	if account == "" { //! 游客登录
		b := make([]byte, 48)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return nil
		}
		h := md5.New()
		h.Write([]byte(base64.URLEncoding.EncodeToString(b)))
		account = hex.EncodeToString(h.Sum(nil))
	}

	var _account San_Account
	sql := fmt.Sprintf("select * from `san_account` where (`account` = '%s' or `account` = '%s') and serverid = %d ", account, account1, serverid)
	GetServer().DBUser.GetOneData(sql, &_account, "", 0)
	//c := GetServer().GetRedisConn()
	//defer c.Close()

	//value, _ := redis.Bytes(c.Do("GET", fmt.Sprintf("%s_%s", "san_account", account)))
	//json.Unmarshal(value, &_account)

	if _account.Uid > 0 {
		if password != _account.Password { //! 密码错误
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_PASSWORD"))
			return nil
		}
		_account.Account = account1
		updateQuery := fmt.Sprintf("update san_account set Account = '%s' where uid=%d limit 1", account1, _account.Uid)
		GetServer().SqlSet(updateQuery)
	} else { //! 插入
		_account.Account = account1
		_account.Password = password
		_account.Creator = "ios"
		_account.ServerId = serverid
		_account.Channelid = appid
		_account.Time = TimeServer().Unix()
		_account.Uid = InsertTable("san_account", &_account, 1, false)
		//_account.Uid = GetServer().GetRedisInc("san_account")
		if _account.Uid <= 0 {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return nil
		}

		//! 重新获取帐号-回避insertid出错
		sql := fmt.Sprintf("select * from `san_account` where (`account` = '%s' or `account` = '%s') and serverid = %d ", account, account1, serverid)
		GetServer().DBUser.GetOneData(sql, &_account, "", 0)
		//value := HF_JtoB(&_account)
		//c := GetServer().GetRedisConn()
		//defer c.Close()
		//_, err := c.Do("MSET", fmt.Sprintf("%s_%d", "san_account", _account.Uid), value,
		//	fmt.Sprintf("%s_%s", "san_account", _account.Account), value)
		//if err != nil {
		//	LogError("redis set err:", "san_account", ",", string(value))
		//	self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
		//	return nil
		//}

		var msg S2C_Reg
		msg.Cid = "reg"
		msg.Uid = _account.Uid
		msg.Account = _account.Account
		msg.Password = _account.Password
		msg.Creator = _account.Creator
		self.SendMsg("1", HF_JtoB(&msg))
	}

	return &_account
}

////////////////////////////////////////////////////////////////////
//! session 管理者

//! 发送消息管道缓冲
const sendChanSize = 2000
const recvChanSize = 2000
const socketTimeOut = 5
const socketTryNum = 10
const broadcast_max = 10

type mapSession map[int64]*Session ///定义客户列表类型

var sessionindex int64 = 0

type SessionMgr struct {
	MapSession        mapSession    //
	Lock              *sync.RWMutex // 锁定
	RemoveLock        *sync.RWMutex
	RemoveSessionList []*Session // 保存删除的session做后的删除
}

var sessionmgrsingleton *SessionMgr = nil

func GetSessionMgr() *SessionMgr {
	if sessionmgrsingleton == nil {
		sessionmgrsingleton = new(SessionMgr)
		sessionmgrsingleton.MapSession = make(mapSession)
		sessionmgrsingleton.Lock = new(sync.RWMutex)
		sessionmgrsingleton.RemoveLock = new(sync.RWMutex)
		sessionmgrsingleton.RemoveSessionList = make([]*Session, 0)
	}

	return sessionmgrsingleton
}

//消息广播
func (self *SessionMgr) Run() {
	for msg := range GetServer().BroadCastMsg {
		if GetServer().ShutDown { //! 关服
			break
		}

		LogDebug("broadcast message:", string(msg))
		var buffer bytes.Buffer
		buffer.Write(HF_DecodeMsg("1", msg))
		self.Lock.RLock()
		for _, value := range self.MapSession {
			if value.PlayerObj != nil {
				value.SendMsgBatch(buffer.Bytes())
			}
		}
		self.Lock.RUnlock()
	}
	/*NEXT:
	for {
		if GetServer().ShutDown {
			break
		}

		msgCount := 0
		if len(GetServer().BroadCastMsg) > broadcast_max {
			msgCount = broadcast_max
		} else {
			msgCount = len(GetServer().BroadCastMsg)
		}

		if msgCount == 0{
			time.Sleep(time.Millisecond * 10)
			continue NEXT
		}
		var buffer bytes.Buffer
		head := "1"
		msgItr := 0
		for i := 0; i < msgCount; i++ {
			if len(GetServer().BroadCastMsg) > 0 {
				msgItr++
				msg := <-GetServer().BroadCastMsg
				buffer.Write(HF_DecodeMsg(head, msg))
			} else {
				self.Lock.Lock()
				for _, value := range self.MapSession {
					value.SendMsgBatch(buffer.Bytes())
				}
				self.Lock.Unlock()

				continue NEXT
			}
		}

		if msgItr > 0 {
			self.Lock.Lock()
			for _, value := range self.MapSession {
				value.SendMsgBatch(buffer.Bytes())
			}
			self.Lock.Unlock()
		}

	}*/
}

func (self *SessionMgr) GetNewSession(ws *websocket.Conn) *Session {
	//if len(GetServer().AddSession) >= 10000 {
	//	return nil
	//}
	self.Lock.Lock()
	sessionindex += 1

	session := new(Session)
	session.ID = sessionindex
	session.Ws = ws
	session.SendChan = make(chan []byte, sendChanSize)
	session.RecvChan = make(chan []byte, recvChanSize)
	session.ShutDown = false
	//session.lockchan = new(sync.RWMutex)
	session.IP = HF_GetHttpIP(session.Ws.Request())
	session.PlayerObj = nil
	session.MsgMaxId = -1

	self.MapSession[sessionindex] = session
	self.Lock.Unlock()
	//GetServer().AddSession <- session

	return session
}

//! 删除session
func (self *SessionMgr) RemoveSession(session *Session) {
	self.Lock.Lock()
	delete(self.MapSession, session.ID)
	self.Lock.Unlock()

	self.RemoveLock.Lock()
	self.RemoveSessionList = append(self.RemoveSessionList, session)
	self.RemoveLock.Unlock()
	//GetServer().DelSession <- session
}

//集中清理session
func (self *SessionMgr) ClearRemoveSession() {
	self.RemoveLock.Lock()
	defer self.RemoveLock.Unlock()
	waitCloseSession := 0
	tNow := TimeServer().Unix()
	for i := 0; i < len(self.RemoveSessionList); i++ {
		if self.RemoveSessionList[i].ShutDown &&
			self.RemoveSessionList[i].SendChan != nil &&
			tNow > self.RemoveSessionList[i].ShutTime {
			close(self.RemoveSessionList[i].SendChan)
			self.RemoveSessionList[i].SendChan = nil
			close(self.RemoveSessionList[i].RecvChan)
			self.RemoveSessionList[i].RecvChan = nil
		} else {
			waitCloseSession++
		}
	}
	if waitCloseSession == 0 {
		self.RemoveSessionList = make([]*Session, 0)
	}
}

func (self *SessionMgr) Barrage(size int, text string, red int, green int, blue int, uid int64) {
	var msg S2C_Barrage
	msg.Cid = "barrage"
	msg.Size = size
	msg.Text = text
	msg.Red = red
	msg.Green = green
	msg.Blue = blue
	msg.Uid = uid
	self.BroadCastMsg("1", HF_JtoB(&msg))
}

//! 广播消息
func (self *SessionMgr) BroadCastMsg(head string, body []byte) {
	//return
	if len(GetServer().BroadCastMsg) >= 5000 {
		return
	}
	GetServer().BroadCastMsg <- body
}

//! 广播消息处理
func (self *SessionMgr) OrderBroadCastMsg(head string, body []byte) {
	//return
	//for _, value := range self.MapSession {
	//	value.SendMsg(head, body)
	//}
}

func (self *SessionMgr) GetSessionNum() int {
	self.Lock.RLock()
	defer self.Lock.RUnlock()

	return len(self.MapSession)
}

// 收到区服消息
func (self *Session) onCenterMsg(body []byte) {
	var centerCid S2S_CenterCid
	err := json.Unmarshal(body, &centerCid)
	if err != nil {
		LogError("CenterCid err:", err)
		return
	}

	cid := centerCid.Cid
	switch cid {
	case "mherotopinfo":
		//info := GetConsumerTop().DoHeroTopinfo(body)
		//self.SendMsg("mherotopinfo", info)
		break
	case "mherodamage":
		//var msg S2Center_UploadHeroDamage
		//err = json.Unmarshal(body, &msg)
		//if err == nil {
		//	GetConsumerTopSvr().UploadDamage(msg.Top)
		//}
		break
	case "mgeneral":
		//var msg S2Center_UploadGeneral
		//err = json.Unmarshal(body, &msg)
		//if err == nil {
		//	GetGeneralMgr().addRank(msg.Top)
		//}
		break
	case "servgeneralrank":
		//var msg S2Center_GeneralRank
		//err = json.Unmarshal(body, &msg)
		//if err == nil {
		//	info := GetGeneralMgr().DoGetRank(msg.ServerId)
		//	self.SendMsg("servgeneralrank", info)
		//}
		break
	default:
		LogError("unsported ctrl, ctrl:", cid)
	}
}

func (self *Session) onTimer() {

	if self.PlayerObj == nil {
		return
	}

	self.PlayerObj.OnTimer()
}
