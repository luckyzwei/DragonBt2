package game

import (
	"encoding/json"
	"strings"

	//"encoding/json"
	"fmt"
	"sync"
	"time"
	//"zephyr"
	//"log"
	"bytes"
)

type mapPlayer map[int64]*Player

type PlayerMgr struct {
	MapPlayer      mapPlayer
	LastUpdateTime int64
	Online         []*Player
	//Offline        []*Player
	OnlineNum    int
	Locker       *sync.RWMutex
	OnlineLocker *sync.RWMutex
}

var playermgrsingleton *PlayerMgr = nil

//! public
func GetPlayerMgr() *PlayerMgr {
	if playermgrsingleton == nil {
		playermgrsingleton = new(PlayerMgr)
		playermgrsingleton.Locker = new(sync.RWMutex)
		playermgrsingleton.OnlineLocker = new(sync.RWMutex)
		playermgrsingleton.MapPlayer = make(mapPlayer)
		playermgrsingleton.LastUpdateTime = 0
		playermgrsingleton.OnlineNum = 0
		playermgrsingleton.Online = make([]*Player, 0)
		//playermgrsingleton.Offline = make([]*Player, 0)
	}

	return playermgrsingleton
}

//! 若create为true，则在找不到的时候重新读取数据，断线重练使用，一般为true
//! 优化locker，减少互斥
func (self *PlayerMgr) GetPlayer(id int64, create bool) *Player {
	if GetServer().ShutDown {
		return nil
	}

	if id <= 0 {
		return nil
	}

	self.Locker.RLock()
	//defer self.Locker.Unlock()
	p, ok := self.MapPlayer[id]
	if !ok {
		self.Locker.RUnlock()
		if create {
			var userbase San_UserBase
			GetServer().DBUser.GetOneData(fmt.Sprintf("select * from `san_userbase` where `uid` = %d", id), &userbase, "", 0)
			if userbase.Uid <= 0 { //! 数据库内找不到
				return nil
			} else {
				p = NewPlayer(id)
				p.Sql_UserBase = userbase
				p.Sql_UserBase.Init("san_userbase", &p.Sql_UserBase, false)
				p.SetAccount(nil)
				p.InitPlayerData()
				//p.OtherPlayerData()
				//self.Locker.Lock()
				//self.MapPlayer[p.ID] = p
				//self.Locker.Unlock()
				//GetServer().Wait.Add(1)
				//GetServer().Event++
				self.AddPlayer(p.ID, p)
			}
		} else {
			return nil
		}
	} else {
		self.Locker.RUnlock()
	}

	p.MsgTime = TimeServer().Unix()
	p.IsSave = true //! 有人获取了该用户的实例
	return p
}

func (self *PlayerMgr) GetPlayerFromName(name string, create bool) *Player {
	if GetServer().ShutDown {
		return nil
	}

	//self.Locker.Lock()
	//defer self.Locker.Unlock()

	self.Locker.RLock()
	for _, value := range self.MapPlayer {
		if value.Sql_UserBase.UName == name {
			self.Locker.RUnlock()
			return value
		}
	}
	self.Locker.RUnlock()
	if create {
		var userbase San_UserBase
		GetServer().DBUser.GetOneData(fmt.Sprintf("select * from `san_userbase` where `uname` = '%s'", name), &userbase, "", 0)
		if userbase.Uid <= 0 {
			//! 数据库内找不到
			return nil
		} else {
			p := NewPlayer(userbase.Uid)
			p.Sql_UserBase = userbase
			p.Sql_UserBase.Init("san_userbase", &p.Sql_UserBase, false)
			p.InitPlayerData()
			p.SetAccount(nil)
			//self.Locker.Lock()
			//self.MapPlayer[p.ID] = p
			//self.Locker.Unlock()
			//GetServer().Wait.Add(1)
			//GetServer().Event++
			self.AddPlayer(p.ID, p)
			p.MsgTime = TimeServer().Unix()
			return p
		}
	}

	return nil
}

func (self *PlayerMgr) GetPlayerFromNameList(name string, create bool) []*Player {
	if GetServer().ShutDown {
		return nil
	}

	self.Locker.RLock()
	playerList := []*Player{}
	for _, value := range self.MapPlayer {
		if strings.Contains(value.Sql_UserBase.UName, name) {
			playerList = append(playerList, value)
		}
	}
	self.Locker.RUnlock()

	if create {
		var userbase San_UserBase
		sql := fmt.Sprintf("select * from `san_userbase` where `uname` like '%%%s%%'", name)
		res := GetServer().DBUser.GetAllData(sql, &userbase)
		if len(res) <= 0 {
			//! 数据库内找不到
			return playerList
		} else {
			for i := 0; i < len(res); i++ {
				data := res[i].(*San_UserBase)

				p := NewPlayer(data.Uid)
				p.Sql_UserBase = *data
				playerList = append(playerList, p)
			}
			return playerList
		}
	}

	return nil
}

func (self *PlayerMgr) AddPlayer(id int64, player *Player) {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	// 停服状态不允许进入服务器
	if GetServer().ShutDown {
		LogInfo("GetServer().ShutDown begin!")
		return
	}
	_, ok := self.MapPlayer[id]
	if ok {
		return
	}

	self.MapPlayer[id] = player
	GetServer().Wait.Add(1)
	GetServer().Event++
}

func (self *PlayerMgr) SetPlayerOnline() {
	self.OnlineLocker.Lock()
	defer self.OnlineLocker.Unlock()

	self.OnlineNum++
}

func (self *PlayerMgr) SetPlayerOffline() {
	self.OnlineLocker.Lock()
	defer self.OnlineLocker.Unlock()

	if self.OnlineNum > 0 {
		self.OnlineNum--
	}

}

func (self *PlayerMgr) GetPlayerOnline() int {
	self.OnlineLocker.RLock()
	defer self.OnlineLocker.RUnlock()

	return self.OnlineNum
}

func (self *PlayerMgr) RemovePlayer(id int64) bool {
	//self.Locker.Lock()
	_, ok := self.MapPlayer[id]
	if ok {
		delete(self.MapPlayer, id)
		GetServer().Wait.Done()
		GetServer().Event--
	}

	return ok

	//self.Locker.Unlock()
}

func (self *PlayerMgr) SaveAll(shutdown bool) {
	//self.Locker.RLock()
	//defer self.Locker.RUnlock()

	for _, value := range self.MapPlayer {
		if shutdown {
			if value.GetSession() != nil {
				value.SafeClose()

				value.Sql_UserBase.LastUpdTime = TimeServer().Format(DATEFORMAT)
				GetOfflineInfoMgr().SetPlayerOffTime(value.Sql_UserBase.Uid, TimeServer().Unix())
				tll, _ := time.ParseInLocation(DATEFORMAT, value.Sql_UserBase.LastLoginTime, time.Local)
				value.Sql_UserBase.LineTime += (TimeServer().Unix() - tll.Unix())
				//GetServer().SqlLineLog(value.Sql_UserBase.Uid, value.Sql_UserBase.IP, int(TimeServer().Unix()-tll.Unix()), value.Account.Creator)
			}
		}
		//log.Println("保存玩家:", value.Sql_UserBase.Uid)
		value.Save(shutdown, true)
		//LogInfo("保存玩家ok:", value.Sql_UserBase.Uid)
		if shutdown {
			GetServer().Wait.Done()
			GetServer().Event--
			//self.RemovePlayer(key)
		}
	}
}

//! 广播消息
func (self *PlayerMgr) BroadCastMsgToCamp(camp int, head string, body []byte) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	var buffer bytes.Buffer
	buffer.Write(HF_DecodeMsg(head, body))
	for _, value := range self.MapPlayer {
		if value.Sql_UserBase.Camp != camp {
			continue
		}
		if value.SessionObj != nil {
			value.SessionObj.SendMsgBatch(buffer.Bytes())
		}
		//value.SendMsg(head, body)
	}
}

//! 得到在线
func (self *PlayerMgr) UpdateOnline() {
	self.OnlineLocker.Lock()
	defer self.OnlineLocker.Unlock()

	self.LastUpdateTime = TimeServer().Unix()
	self.Online = make([]*Player, 0)
	self.Locker.RLock()
	for _, value := range self.MapPlayer {
		if value.GetSession() != nil {
			self.Online = append(self.Online, value)
		}
	}
	self.Locker.RUnlock()

	self.OnlineNum = len(self.Online)
}

//!随机得到在线人数
func (self *PlayerMgr) GetOnlineRandom(num int) []*Player {
	if TimeServer().Unix()-self.LastUpdateTime > 60 {
		self.UpdateOnline()
	}

	self.OnlineLocker.RLock()
	defer self.OnlineLocker.RUnlock()
	lst := make([]*Player, 0)
	for i := 0; i < len(self.Online); i++ {
		lst = append(lst, self.Online[i])
	}

	return lst

	//优化以后在说
	/*
		//! 每分钟更新一次
		if TimeServer().Unix()-self.LastUpdateTime > 60 {
			self.UpdateOnline()
		}

		self.OnlineLocker.RLock()
		defer self.OnlineLocker.RUnlock()
		lst := make([]*Player, 0)
		if len(self.Online) < num {
			for i := 0; i < len(self.Online); i++ {
				lst = append(lst, self.Online[i])
			}

			return lst
		}

		for len(self.Online) > 0 && len(lst) < num {
			index := HF_GetRandom(len(self.Online))
			player := self.Online[index]
			lst = append(lst, player)
		}

		return lst
	*/
}

//! 加邮件
func (self *PlayerMgr) AddMail(mail *JS_Mail, id int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for _, value := range self.MapPlayer {
		if value.Sql_UserBase.Level < mail.MinLevel {
			continue
		}
		if mail.MaxLevel != 0 && value.Sql_UserBase.Level > mail.MaxLevel {
			continue
		}
		value.GetModule("mail").(*ModMail).AddGlobalMail(id, MAIL_CAN_ALL_GET, 2, 0, mail.Title,
			mail.Content, mail.Sender, mail.Item, true, 0)
	}
}

//! 得到在线区分渠道
func (self *PlayerMgr) GetOnlineByGameId() map[string]int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	res := make(map[string]int)
	for _, value := range self.MapPlayer {
		if value.GetSession() != nil {
			gameId := GetServer().Con.GetGameIdByAppId(value.GetAppleId())
			res[gameId]++
		}
	}

	return res

}

type Js_LuckShop struct {
	Uid  int64
	Info string
	info []JS_LuckShopItem
}

// 检查活动
func (self *PlayerMgr) CheckAct() {
	var data Js_LuckShop
	sql := "select `uid`, `info` from `san_luckshop` order by `uid`"
	res := GetServer().DBUser.GetAllData(sql, &data)
	for i := 0; i < len(res); i++ {
		data := res[i].(*Js_LuckShop)
		if data == nil {
			continue
		}

		json.Unmarshal([]byte(data.Info), &data.info)

		found := false
		for itemIndex := range data.info {
			item := data.info[itemIndex]
			if item.Done == 1 {
				found = true
				break
			}
		}

		//log.Println(data.UID)
		if found {
			player := self.GetPlayer(data.Uid, true)
			if player != nil {
				player.GetModule("luckshop").(*ModLuckShop).FitServer()
			}
		}

	}
}

func (self *PlayerMgr) KickoutPlayers() {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	for _, player := range self.MapPlayer {
		if player != nil {
			player.SafeClose()
		}
	}
}

//返回内存中的实际人数
func (self *PlayerMgr) GetPlayerByDebug() int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()
	return len(self.MapPlayer)
}
