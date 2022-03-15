package game

import (
	"crypto/md5"
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"time"

	//"time"
)

const (
	SDKUP_EVENT_ID_OFFLINE         = 114
	SDKUP_EVENT_ID_MONEY_CHANGE    = 146
	SDKUP_EVENT_ID_RESOURCE_CHANGE = 147
	SDKUP_EVENT_ID_ONLINE          = 1001

	SKDUP_ID = "0"

	SKDUP_ADDR_URL = "https://api.kylinmobi.net/v2/logs/common/" //正式
	//SKDUP_ADDR_URL = "https://api-java-t.kylinmobi.net/v2/logs/common/" //测试

	//NEW
	SKDUP_ADDR_URL_AIWAN_SDK                = "https://cpl-api.jxywl.cn/game/userdata/v2"
	SKDUP_ADDR_URL_AIWAN_SDK_EVENT_ENTERSVR = "entersvr"
	SKDUP_ADDR_URL_AIWAN_SDK_EVENT_LEVELUP  = "levelup"
	SKDUP_ADDR_URL_AIWAN_SDK_APPID          = "9t9n7byk3l8q0qyy"                 //爱玩提供
	SKDUP_ADDR_URL_AIWAN_SDK_CHANNELID1     = "14097"                            //
	SKDUP_ADDR_URL_AIWAN_SDK_CHANNELID2     = "14098"                            //
	SKDUP_ADDR_URL_AIWAN_SDK_APPKEY         = "LYNIscgVfoSZLARJPfEEJuLJkVqrCPNl" //
)

//! 经分-日志接口

//! 设置envInfo
func (self *Server) GetEnvInfo(player *Player) SendRZ_envInfo {

	var data2 SendRZ_envInfo

	if player.getPlatNo() == PLATFORM_IOS {
		data2.AccInfo.AccountId = fmt.Sprintf("%s_%s", player.Platform.AccountId, self.Con.GetGameIdByAppId(player.GetAppleId()))
	} else {
		if player.Account != nil {
			data2.AccInfo.AccountId = player.Account.Account
		}
	}

	//data2.AccInfo.Account_appId = player.Platform.Account_AppId
	if player.Account != nil {
		data2.AccInfo.Creator = player.Account.Creator
	}
	data2.AccInfo.UserType = ""

	// data2.DevInfo.Brand = player.Platform.Brand
	// data2.DevInfo.DeviceId = player.Platform.DeviceId
	// data2.DevInfo.Model = player.Platform.Model
	data2.DevInfo.Brand = ""
	data2.DevInfo.DeviceId = ""
	data2.DevInfo.Model = ""
	data2.DevInfo.Os = player.Platform.Platform

	data2.GmInfo.AppVer = "3.2.0"
	data2.GmInfo.BuildVer = "3.2.0"
	data2.GmInfo.PkgName = "com.congyu.xrzp"
	data2.GmInfo.ResVer = "3.2.0"

	data2.RunId = ""

	return data2
}

//! 设置envInfo
func (self *Server) GetEnvInfoIos(player *Player) SendRZ_envInfo_ios {

	var data2 SendRZ_envInfo_ios
	data2.AccInfo.AccountId = fmt.Sprintf("%s_%s", player.Platform.AccountId, self.Con.GetGameIdByAppId(player.GetAppleId()))
	//data2.AccInfo.Account_appId = player.Platform.Account_AppId
	data2.AccInfo.Creator = player.Account.Creator
	data2.AccInfo.UserType = ""

	data2.DevInfo.DeviceId = player.Platform.DeviceId
	data2.DevInfo.UUID = player.Platform.UUID
	data2.DevInfo.Brand = player.Platform.Brand
	data2.DevInfo.Model = player.Platform.Model
	data2.DevInfo.Os = player.Platform.Platform

	data2.DevInfo.Fr = player.Platform.Fr
	data2.DevInfo.Res = player.Platform.Res
	data2.DevInfo.Net = player.Platform.Net
	data2.DevInfo.Mac = player.Platform.Mac
	data2.DevInfo.Operator = player.Platform.Operator
	data2.DevInfo.Ip = player.Platform.Ip

	data2.ChInfo.Ch = player.Platform.Ch
	data2.ChInfo.SubCh = player.Platform.SubCh
	//data2.GmInfo.AppVer = "3.2.0"
	//data2.GmInfo.BuildVer = "3.2.0"
	//data2.GmInfo.PkgName = "com.congyu.xrzp"
	//data2.GmInfo.ResVer = "3.2.0"

	data2.RunId = ""

	return data2
}

// 用户在游戏服务器登录成功事件
func (self *Server) sendLog_userLoginOK(player *Player) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}
	var data SendRZ_Loginok
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "user.online"

	var data1 SendRZ_Loginok_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	//data1.ServerName = self.Con.ServerName
	//data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	//data1.RoleLevel = player.Sql_UserBase.Level
	//data1.RoleName = player.Sql_UserBase.UName
	//data1.VipLevel = player.Sql_UserBase.Vip
	//data1.Force = player.Sql_UserBase.Fight
	//data1.Diamond = player.Sql_UserBase.Gem
	//data1.Gold = player.Sql_UserBase.Gold
	//data1.Power = player.Sql_UserBase.TiLi
	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

// 用户在游戏服务器登录成功事件
func (self *Server) sendLog_userLoginOKIOS(player *Player) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}
	var data SendRZ_LoginokIOS
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "user.online"

	var data1 SendRZ_LoginokIOS_params
	data.Params = data1
	data.EnvInfo = self.GetEnvInfoIos(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

// 创建角色成功
func (self *Server) sendLog_createOK(player *Player) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	if player.getPlatNo() == PLATFORM_IOS {
		self.sendLog_createOKIOS(player)
		return
	}

	var data SendRZ_Create
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "role.create"

	var data1 SendRZ_Create_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

// 创建角色成功
func (self *Server) sendLog_createOKIOS(player *Player) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_CreateIOS
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "role.create"

	var data1 SendRZ_Create_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName

	data.Params = data1
	data.EnvInfo = self.GetEnvInfoIos(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

// 角色升级成功
func (self *Server) sendLog_LevelupOk(player *Player) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	if player.getPlatNo() == PLATFORM_IOS {
		self.sendLog_LevelupOkIOS(player)
		return
	}

	var data SendRZ_Levelup
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "role.levelup"

	var data1 SendRZ_Levelup_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.RoleName = player.Sql_UserBase.UName
	data1.VipLevel = player.Sql_UserBase.Vip

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

	//	str += HF_JtoA(data1)

	//	str += "`"
	//	str += "envInfo"
	//	str += HF_JtoA(self.GetEnvInfo(player))
	//	self.Log(data.AppId, []byte(str))

	//	str += "`"
	//	str += "deviceInfo="

	//	var data2 SendRZ_deviceInfo
	//	data2.Runid = 1
	//	data2.BindIds.UserType = ""
	//	data2.BindIds.Creator = player.Account.Creator
	//	data2.BindIds.AccountId = player.Account.Account
	//	data2.Metrics.Os = player.Platform.Platform
	//	str += HF_JtoA(data2)

	//	//log.Println("sendLog_LevelupOk  platform:", data2.Metrics.Os)
	//	self.Log(data.AppId, []byte(str))
	//log.Println("————————————————————登录成功：", str)
}

// 登录成功
func (self *Server) sendLog_LoginOK(player *Player) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	if player.getPlatNo() == PLATFORM_IOS {
		self.sendLog_LoginOKIOS(player)
		return
	}

	var data SendRZ_Loginok
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "role.online"

	var data1 SendRZ_Loginok_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.RoleName = player.Sql_UserBase.UName
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.VipExp = player.Sql_UserBase.VipExp
	data1.Force = player.Sql_UserBase.Fight
	data1.Diamond = player.Sql_UserBase.Gem
	data1.Gold = player.Sql_UserBase.Gold
	data1.Power = player.Sql_UserBase.TiLi
	data1.Feats = player.GetObjectNum(91000018)

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

}

// 登录成功
func (self *Server) sendLog_LoginOKIOS(player *Player) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}
	var data SendRZ_LoginokIOSEX
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "role.online"

	var data1 SendRZ_Loginok_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Force = player.Sql_UserBase.Fight
	data1.Diamond = player.Sql_UserBase.Gem
	data1.Gold = player.Sql_UserBase.Gold
	data1.Power = player.Sql_UserBase.TiLi
	data1.Feats = player.GetObjectNum(91000018)
	data.Params = data1
	data.EnvInfo = self.GetEnvInfoIos(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

}

// 登出
func (self *Server) sendLog_Offline(player *Player, logintime int64, logouttime int64, duration int64) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}
	var data SendRZ_Offline
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "role.offline"

	var data1 SendRZ_Offline_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName
	data1.RoleLevel = player.Sql_UserBase.Level

	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Vip_exp = player.Sql_UserBase.VipExp
	data1.LoginTime = logintime
	data1.LogoutTime = logouttime
	data1.Duration = duration

	data1.Force = player.Sql_UserBase.Fight
	data1.Diamond = player.Sql_UserBase.Gem
	data1.Gold = player.Sql_UserBase.Gold
	data1.Power = player.Sql_UserBase.TiLi
	data1.Feats = player.GetObjectNum(91000018)

	for _, vhero := range player.GetModule("hero").(*ModHero).Sql_Hero.info {
		var heroinfo SendRZ_HeroInfo_params
		heroinfo.HeroId = vhero.HeroId
		heroinfo.HeroLevel = player.Sql_UserBase.Level
		heroinfo.HeroStar = 1
		heroinfo.HeroRank = 1
		data1.HeroInfo = append(data1.HeroInfo, heroinfo)
	}
	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

}

// 获得虚拟币
func (self *Server) sendLog_GetMoney(player *Player, num int, source string, coinType string, coins int) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	defer func() {
		x := recover()
		if x != nil {
			log.Println(x, string(debug.Stack()))
			LogError(x, string(debug.Stack()))
		}
	}()

	var data SendRZ_GetMoney
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "coins.gain"

	//	str := ""
	//	str += fmt.Sprintf("ts=%d", TimeServer().Unix())
	//	str += "`"
	//	str += fmt.Sprintf("appId=%s", "10000008")
	//	str += "`"
	//	str += fmt.Sprintf("event=%s", "coins.gain")
	//	str += "`"
	//	str += "params="

	var data1 SendRZ_GetMoney_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.RoleName = player.Sql_UserBase.UName
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Number = num
	data1.Source = source
	data1.CoinType = coinType
	data1.Coins = coins

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

	//	str += HF_JtoA(data1)

	//	str += "`"
	//	str += "envInfo"
	//	str += HF_JtoA(self.GetEnvInfo(player))
	//	self.Log(data.AppId, []byte(str))

	//	str += "`"
	//	str += "deviceInfo="

	//	var data2 SendRZ_deviceInfo
	//	data2.Runid = 1
	//	data2.BindIds.UserType = ""
	//	data2.BindIds.Creator = player.Account.Creator
	//	data2.BindIds.AccountId = player.Account.Account
	//	data2.Metrics.Os = player.Platform.Platform
	//	str += HF_JtoA(data2)

	//	//log.Println("sendLog_GetMoney  platform:", data2.Metrics.Os)
	//	self.Log(data.AppId, []byte(str))
	//log.Println("————————————————————购买虚拟币：", str)
}

// 消耗虚拟币
func (self *Server) sendLog_UseMoney(player *Player, num int, destination string, coinType string, coins int) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_UseMoney
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "coins.consume"

	//	str := ""
	//	str += fmt.Sprintf("ts=%d", TimeServer().Unix())
	//	str += "`"
	//	str += fmt.Sprintf("appId=%s", "10000008")
	//	str += "`"
	//	str += fmt.Sprintf("event=%s", "coins.consume")
	//	str += "`"
	//	str += "params="

	var data1 SendRZ_UseMoney_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.RoleName = player.Sql_UserBase.UName
	data1.VipLevel = player.Sql_UserBase.Vip
	if num < 0 {
		num = num * -1
	}
	data1.Number = num
	data1.Destination = destination
	data1.CoinType = coinType
	data1.Coins = coins

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

	//	str += HF_JtoA(data1)

	//	str += "`"
	//	str += "envInfo"
	//	str += HF_JtoA(self.GetEnvInfo(player))
	//	self.Log(data.AppId, []byte(str))
	//	str += "`"
	//	str += "deviceInfo="

	//	var data2 SendRZ_deviceInfo
	//	data2.Runid = 1
	//	data2.BindIds.UserType = ""
	//	data2.BindIds.Creator = player.Account.Creator
	//	data2.BindIds.AccountId = player.Account.Account
	//	data2.Metrics.Os = player.Platform.Platform
	//	str += HF_JtoA(data2)

	//	//log.Println("sendLog_UseMoney  platform:", data2.Metrics.Os)

	//	self.Log(data.AppId, []byte(str))
	//log.Println("————————————————————消耗虚拟币：", str)
}

// 获得物品
func (self *Server) sendLog_GetItem(player *Player, num int, source string, itemId int, item string, curcount int) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_GetItem
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "items.gain"

	//	str := ""
	//	str += fmt.Sprintf("ts=%d", TimeServer().Unix())
	//	str += "`"
	//	str += fmt.Sprintf("appId=%s", "10000008")
	//	str += "`"
	//	str += fmt.Sprintf("event=%s", "items.gain")
	//	str += "`"
	//	str += "params="

	var data1 SendRZ_GetItem_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.RoleName = player.Sql_UserBase.UName
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Source = source
	data1.ItemId = fmt.Sprintf("%d", itemId)
	data1.ItemType = GetCsvMgr().GetItemType(itemId)
	data1.Item = item
	data1.Number = num
	data1.ItemCount = curcount

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

	//	str += HF_JtoA(data1)

	//	str += "`"
	//	str += "envInfo"
	//	str += HF_JtoA(self.GetEnvInfo(player))
	//	self.Log(data.AppId, []byte(str))

	//	str += "`"
	//	str += "deviceInfo="

	//	var data2 SendRZ_deviceInfo
	//	data2.Runid = 1
	//	data2.BindIds.UserType = ""
	//	data2.BindIds.Creator = player.Account.Creator
	//	data2.BindIds.AccountId = player.Account.Account
	//	data2.Metrics.Os = player.Platform.Platform
	//	str += HF_JtoA(data2)

	//	self.Log(data.AppId, []byte(str))
	//	//log.Println("————————————————————获得物品：", str)
}

// 消耗物品
func (self *Server) sendLog_UseItem(player *Player, num int, destination string, itemId int, item string, curcount int) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_UseItem
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "items.consume"

	//	str := ""
	//	str += fmt.Sprintf("ts=%d", TimeServer().Unix())
	//	str += "`"
	//	str += fmt.Sprintf("appId=%s", "10000008")
	//	str += "`"
	//	str += fmt.Sprintf("event=%s", "items.consume")
	//	str += "`"
	//	str += "params="

	var data1 SendRZ_UseItem_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.RoleName = player.Sql_UserBase.UName
	data1.Destination = destination
	data1.ItemType = GetCsvMgr().GetItemType(itemId)
	data1.ItemId = fmt.Sprintf("%d", itemId)
	data1.Item = item
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.ItemCount = curcount
	if num < 0 {
		num = num * -1
	}
	data1.Number = num

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

	//	str += HF_JtoA(data1)

	//	str += "`"
	//	str += "envInfo"
	//	str += HF_JtoA(self.GetEnvInfo(player))
	//	self.Log(data.AppId, []byte(str))

	//	str += "`"
	//	str += "deviceInfo="

	//	var data2 SendRZ_deviceInfo
	//	data2.Runid = 1
	//	data2.BindIds.UserType = ""
	//	data2.BindIds.Creator = player.Account.Creator
	//	data2.BindIds.AccountId = player.Account.Account
	//	data2.Metrics.Os = player.Platform.Platform
	//	str += HF_JtoA(data2)

	//	self.Log(data.AppId, []byte(str))
	//	//log.Println("————————————————————消耗物品：", str)
}

// 道具购买 玩家，消耗货币，消费去向，道具类型，道具名字，货币类型，购买数量，购买地址
func (self *Server) sendLog_BuyItem(player *Player, num int, destination string, itemId int, item string, coinType string, itemnum int, place string) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_BuyItem
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "items.pay"

	//	str := ""
	//	str += fmt.Sprintf("ts=%d", TimeServer().Unix())
	//	str += "`"
	//	str += fmt.Sprintf("appId=%s", "10000008")
	//	str += "`"
	//	str += fmt.Sprintf("event=%s", "items.pay")
	//	str += "`"
	//	str += "params="

	var data1 SendRZ_BuyItem_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.RoleName = player.Sql_UserBase.UName
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Place = place
	data1.Number = num
	//data1.Destination = destination
	data1.CoinType = coinType
	data1.ItemType = GetCsvMgr().GetItemType(itemId)
	data1.ItemId = fmt.Sprintf("%d", itemId)
	data1.Item = item
	data1.Amount = itemnum

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

	//	str += HF_JtoA(data1)

	//	str += "`"
	//	str += "envInfo"
	//	str += HF_JtoA(self.GetEnvInfo(player))
	//	self.Log(data.AppId, []byte(str))
	//	str += "`"
	//	str += "deviceInfo="

	//	var data2 SendRZ_deviceInfo
	//	data2.Runid = 1
	//	data2.BindIds.UserType = ""
	//	data2.BindIds.Creator = player.Account.Creator
	//	data2.BindIds.AccountId = player.Account.Account
	//	data2.Metrics.Os = player.Platform.Platform
	//	LogDebug("平台：", player.Platform.Platform)
	//	str += HF_JtoA(data2)

	//	self.Log(data.AppId, []byte(str))
}

// PVE战斗结束
func (self *Server) sendLog_pve(player *Player, fightnum int, levelname string, pass string, fighttime int, passedType string) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}
	var data SendRZ_PVE
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "activity.finish"

	//	str := ""
	//	str += fmt.Sprintf("ts=%d", TimeServer().Unix())
	//	str += "`"
	//	str += fmt.Sprintf("appId=%s", "10000008")
	//	str += "`"
	//	str += fmt.Sprintf("event=%s", "activity.finish")
	//	str += "`"
	//	str += "params="

	var data1 SendRZ_PVE_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.RoleName = player.Sql_UserBase.UName
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Begin_time = TimeServer().UnixNano() / 1e6
	data1.Force = player.Sql_UserBase.Fight
	data1.TeamForce = fightnum
	data1.MissionType = "剧情关卡"
	data1.Mission = levelname
	data1.Passed = pass
	data1.Duration = fighttime
	data1.PassedType = passedType // 战斗，劝降，免战
	data1.End_roleLevel = player.Sql_UserBase.Level

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

// PVE战斗结束
func (self *Server) sendLog_pve_ex(player *Player, fightnum int, levelname string, pass string, fighttime int, passedType string, missiontype string) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}
	var data SendRZ_PVE
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "activity.finish"

	var data1 SendRZ_PVE_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.RoleName = player.Sql_UserBase.UName
	data1.VipLevel = player.Sql_UserBase.Vip

	data1.Begin_time = TimeServer().UnixNano() / 1e6
	data1.Force = player.Sql_UserBase.Fight
	data1.TeamForce = fightnum
	data1.MissionType = missiontype
	data1.Mission = levelname
	data1.Passed = pass
	data1.Duration = fighttime
	data1.PassedType = passedType // 战斗，劝降，免战
	data1.End_roleLevel = player.Sql_UserBase.Level

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

// PVP战斗结束
func (self *Server) sendLog_pvp(player *Player, fightnum int64, levelname string, pass string, fighttime int, pvpid int64, teamid int, pvpname string) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}
	var data SendRZ_PVP
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "activity.pvp.finish"

	var data1 SendRZ_PVP_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.RoleName = player.Sql_UserBase.UName
	data1.VipLevel = player.Sql_UserBase.Vip

	data1.Begin_time = TimeServer().UnixNano() / 1e6
	data1.Force = player.Sql_UserBase.Fight
	data1.TeamForce = fightnum
	data1.MissionType = "PVP挑战"
	data1.Mission = levelname
	data1.Passed = pass
	data1.Duration = fighttime
	data1.End_roleLevel = player.Sql_UserBase.Level

	data1.Pvpid = fmt.Sprintf("%d", pvpid) // 没有
	data1.Pvpname = pvpname
	data1.Teamid = fmt.Sprintf("%d", teamid) // 没有

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

// 在线玩家统计, 需要区分不同的游戏类型
func (self *Server) sendLog_onTime() {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	res := GetPlayerMgr().GetOnlineByGameId()
	for gameId, onLineNum := range res {
		var data SendRZ_onTime
		data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
		data.AppId = gameId
		data.Event = "online.cu"

		//	str := ""
		//	str += fmt.Sprintf("ts=%d", TimeServer().Unix())
		//	str += "`"
		//	str += fmt.Sprintf("appId=%s", "10000008")
		//	str += "`"
		//	str += fmt.Sprintf("event=%s", "online.cu")
		//	str += "`"
		//	str += "params="

		// 当前时间转换
		str_time := ""
		str_time += fmt.Sprintf("%d", TimeServer().Year())
		if TimeServer().Month() < 10 {
			str_time += fmt.Sprintf("%d", 0)
			str_time += fmt.Sprintf("%d", TimeServer().Month())
		} else {
			str_time += fmt.Sprintf("%d", TimeServer().Month())
		}
		if TimeServer().Day() < 10 {
			str_time += fmt.Sprintf("%d", 0)
			str_time += fmt.Sprintf("%d", TimeServer().Day())
		} else {
			str_time += fmt.Sprintf("%d", TimeServer().Day())
		}
		if TimeServer().Hour() < 10 {
			str_time += fmt.Sprintf("%d", 0)
			str_time += fmt.Sprintf("%d", TimeServer().Hour())
		} else {
			str_time += fmt.Sprintf("%d", TimeServer().Hour())
		}
		if TimeServer().Minute() < 10 {
			str_time += fmt.Sprintf("%d", 0)
			str_time += fmt.Sprintf("%d", TimeServer().Minute())
		} else {
			str_time += fmt.Sprintf("%d", TimeServer().Minute())
		}

		var data1 SendRZ_onTime_params
		data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
		data1.ServerName = self.Con.ServerName
		data1.Timevalue = str_time
		data1.UserCnt = onLineNum

		data.Params = data1
		var data2 SendRZ_envInfo
		data.EnvInfo = data2
		//data.EnvInfo.DevInfo.Os = "android"
		//data.EnvInfo.DevInfo.Os = "ios"

		self.Log(data.AppId, []byte(HF_JtoA(data)))
	}

}

// 聊天日志
func (self *Server) sendLog_chat(player *Player, channel string, text string) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_Chat
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "chat"

	var data1 SendRZ_Chat_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.RoleName = player.Sql_UserBase.UName
	data1.Type = channel
	data1.Text = text

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.LogChat(data.AppId, []byte(HF_JtoA(data)))
}

// 获得武将
func (self *Server) sendLog_herogain(player *Player, hero *Hero, source string) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_HeroGain
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "hero.gain"

	var data1 SendRZ_HeroGain_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.VipLevel = player.Sql_UserBase.Vip

	config := GetCsvMgr().GetHeroConfig(hero.HeroId)
	data1.HeroType = HF_GetHeroType(hero.HeroId)
	data1.HeroId = strconv.Itoa(int(hero.HeroId))
	data1.HeroName = config.HeroName
	data1.HeroLevel = player.Sql_UserBase.Level
	data1.HeroStar = 1
	data1.HeroRank = 1

	data1.HeroForce = hero.Fight

	data1.Source = source
	data1.ChipsCount = HF_GetHeroWuHun(player, hero.HeroId)

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

}

// 武将升级
func (self *Server) sendLog_herolevelup(player *Player, hero *Hero, source string) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_HeroLevelUp
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "hero.levelup"

	var data1 SendRZ_HeroLevelUp_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Force = player.Sql_UserBase.Fight

	config := GetCsvMgr().GetHeroConfig(hero.HeroId)
	data1.HeroType = HF_GetHeroType(hero.HeroId)
	data1.HeroId = strconv.Itoa(int(hero.HeroId))
	data1.HeroName = config.HeroName
	data1.HeroLevel = player.Sql_UserBase.Level
	data1.HeroStar = 1
	data1.HeroRank = 1

	data1.HeroForce = hero.Fight

	data1.Source = source
	data1.ChipsCount = HF_GetHeroWuHun(player, hero.HeroId)

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

}

// 武将升星
func (self *Server) sendLog_herostarup(player *Player, hero *Hero) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_HeroStarup
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "hero.starup"

	var data1 SendRZ_HeroStarup_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Force = player.Sql_UserBase.Fight

	config := GetCsvMgr().GetHeroConfig(hero.HeroId)
	data1.HeroType = HF_GetHeroType(hero.HeroId)
	data1.HeroId = strconv.Itoa(int(hero.HeroId))
	data1.HeroName = config.HeroName
	data1.HeroLevel = player.Sql_UserBase.Level
	data1.HeroStar = 1
	data1.HeroRank = 1

	data1.HeroForce = hero.Fight

	data1.ChipsCount = HF_GetHeroWuHun(player, hero.HeroId)

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

}

// 武将升阶
func (self *Server) sendLog_herorankup(player *Player, hero *Hero) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_HeroStarup
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "hero.rankup"

	var data1 SendRZ_HeroStarup_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Force = player.Sql_UserBase.Fight

	config := GetCsvMgr().GetHeroConfig(hero.HeroId)
	data1.HeroType = HF_GetHeroType(hero.HeroId)
	data1.HeroId = strconv.Itoa(int(hero.HeroId))
	data1.HeroName = config.HeroName
	data1.HeroLevel = player.Sql_UserBase.Level
	data1.HeroStar = 1
	data1.HeroRank = 1

	data1.HeroForce = hero.Fight

	data1.ChipsCount = HF_GetHeroWuHun(player, hero.HeroId)

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

}

//！武将上阵
func (self *Server) sendLog_herojoin(player *Player, hero *Hero, teamidx int) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_HeroJoin
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "hero.join"

	var data1 SendRZ_HeroJoin_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Force = player.Sql_UserBase.Fight

	config := GetCsvMgr().GetHeroConfig(hero.HeroId)
	data1.HeroType = HF_GetHeroType(hero.HeroId)
	data1.HeroId = strconv.Itoa(int(hero.HeroId))
	data1.HeroName = config.HeroName
	data1.HeroLevel = player.Sql_UserBase.Level
	data1.HeroStar = 1
	data1.HeroRank = 1

	data1.HeroForce = hero.Fight

	data1.ChipsCount = HF_GetHeroWuHun(player, hero.HeroId)

	if teamidx == TYPE_CASERN_GONG {
		data1.CampName = "弓兵营"
	} else if teamidx == TYPE_CASERN_QI {
		data1.CampName = "骑兵营"
	} else if teamidx == TYPE_CASERN_BU {
		data1.CampName = "步兵营"
	}

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

}

//武将统兵升级
func (self *Server) sendLog_armylevelup(player *Player, hero *Hero, costitem int, costnum int) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_ArmyLevelUp
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "army.levelup"

	var data1 SendRZ_ArmyLevelUp_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Force = player.Sql_UserBase.Fight

	config := GetCsvMgr().GetHeroConfig(hero.HeroId)
	data1.HeroType = HF_GetHeroType(hero.HeroId)
	data1.HeroId = strconv.Itoa(int(hero.HeroId))
	data1.HeroName = config.HeroName
	data1.HeroLevel = player.Sql_UserBase.Level
	data1.HeroStar = 1
	data1.HeroRank = 1

	data1.HeroForce = hero.Fight

	data1.ArmyLevel = 0
	data1.ArmyForce = 0 //兵营战力？

	data1.Book = costnum
	data1.BookCount = player.GetObjectNum(costitem)

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

}

//1.3.27 	威名进阶
func (self *Server) sendLog_famerank(player *Player, rank int) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_FameRank
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "fame.rank"

	var data1 SendRZ_FameRank_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Force = player.Sql_UserBase.Fight

	data1.Rank = rank

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

}

//1.3.28 	武将置换
func (self *Server) sendLog_herodisplace(player *Player, heroBase *Hero, heroTarget *Hero) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_HeroDisplace
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "hero.displace"

	var data1 SendRZ_HeroDisplace_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Force = player.Sql_UserBase.Fight

	data1.HeroIdBase = heroBase.HeroId
	//data1.HeroLevelBase = heroBase.Levels
	data1.ArmyLevelBase = 0
	data1.WeaponLevelBase = 0
	data1.ClothLevelBase = 0
	data1.OrnamentLevelBase = 0

	data1.HeroIdTarget = heroTarget.HeroId
	//data1.HeroLevelTarget = heroTarget.Levels
	data1.ArmyLevelTarget = 0
	data1.WeaponLevelTarget = 0
	data1.ClothLevelTarget = 0
	data1.OrnamentLevelTarget = 0

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))

}

//1.3.29 	每日政务
func (self *Server) sendLog_TaskFinish(player *Player, taskid int, tasktype int, taskname string, taskstar int) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_TaskFinish
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "task.finish"

	var data1 SendRZ_TaskFinish_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Force = player.Sql_UserBase.Fight

	data1.Task_Type = taskid
	data1.Task_Id = tasktype
	data1.Task_Name = taskname
	data1.Task_Star = taskstar

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

//3.1.31 	客户端激活
func (self *Server) sendLog_Activation(player *Player) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_Activation
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "activation"

	var data1 SendRZ_Activation_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

//3.1.31 	客户端激活
func (self *Server) sendLog_ActivationIOS(player *Player) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_ActivationIOS
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "activation"

	var data1 SendRZ_Activation_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName

	data.Params = data1
	data.EnvInfo = self.GetEnvInfoIos(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

//3.1.32 	账号创建
func (self *Server) sendLog_UserCreate(player *Player) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_Activation
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "user.create"

	var data1 SendRZ_Activation_params
	//data1.ServerId = fmt.Sprintf("%d", self.Con.ServerId)
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

//3.1.32 	账号创建
func (self *Server) sendLog_UserCreateIOS(player *Player) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_ActivationIOS
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "user.create"

	var data1 SendRZ_Activation_params
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName

	data.Params = data1
	data.EnvInfo = self.GetEnvInfoIos(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

// 3.1.33 	角色升级 ios
func (self *Server) sendLog_LevelupOkIOS(player *Player) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}
	var data SendRZ_LevelupIOS
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "role.levelup"

	var data1 SendRZ_LevelupIOS_params
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.RoleName = player.Sql_UserBase.UName
	data1.VipLevel = player.Sql_UserBase.Vip

	data.Params = data1
	data.EnvInfo = self.GetEnvInfoIos(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

// 3.1.34 	账号充值成功
func (self *Server) sendLog_AccountChargeSuccess(player *Player, payment string, dsorderId string, currency string, amount float64, goodsInfo string) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}
	var data SendRZ_AccountChargeSuccess
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "account.charge.success"

	var data1 SendRZ_AccountChargeSuccess_params
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.RoleName = player.Sql_UserBase.UName

	data1.PayMent = payment
	data1.DsorderId = dsorderId
	data1.Currency = currency
	data1.Amount = amount
	data1.GoodsInfo = goodsInfo

	data.Params = data1
	data.EnvInfo = self.GetEnvInfoIos(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

//1.3.14 	加入商会
func (self *Server) sendLog_CommerceJoin(player *Player, type_id int, type_name string) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}

	var data SendRZ_CommerceJoin
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "commerce.join"

	var data1 SendRZ_CommerceJoin_params
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Force = player.Sql_UserBase.Fight

	data1.Type_Id = type_id
	data1.Type_Name = type_name

	data.Params = data1
	data.EnvInfo = self.GetEnvInfo(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

// 心跳上报
func (self *Server) sendLog_HeartBeatIOS(player *Player) {
	if self.Con.ServerExtCon.UpRecord == 0 {
		return
	}
	var data SendRZ_LoginokIOSEX
	data.Ts = fmt.Sprintf("%d", TimeServer().UnixNano()/1e6)
	data.AppId = self.Con.GetGameIdByAppId(player.GetAppleId())
	data.Event = "heartbeat"

	var data1 SendRZ_Loginok_params
	data1.ServerId = fmt.Sprintf("%d", player.GetServerId())
	data1.ServerName = self.Con.ServerName
	data1.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data1.RoleName = player.Sql_UserBase.UName
	data1.RoleLevel = player.Sql_UserBase.Level
	data1.VipLevel = player.Sql_UserBase.Vip
	data1.Force = player.Sql_UserBase.Fight
	data1.Diamond = player.Sql_UserBase.Gem
	data1.Gold = player.Sql_UserBase.Gold
	data1.Power = player.Sql_UserBase.TiLi
	data.Params = data1
	data.EnvInfo = self.GetEnvInfoIos(player)
	self.Log(data.AppId, []byte(HF_JtoA(data)))
}

func (self *Server) SendLog_SDKUP_Offline(player *Player) {
	now := time.Now().Unix()

	str := fmt.Sprintf("%s%d", SKDUP_ID, now)
	h := md5.New()
	h.Write([]byte(str))
	md5str := fmt.Sprintf("%x", h.Sum(nil))

	url := fmt.Sprintf("%s/%s/roleLogout?x-dbsv1-sign=%s&x-dbsv1-timestamp=%d&x-dbsv1-appId=%s",
		SKDUP_ADDR_URL,
		SKDUP_ID,
		md5str,
		now,
		SKDUP_ID)

	var data SendSDK_Offline
	data.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data.RoleName = player.Sql_UserBase.UName
	data.VipLevel = player.Sql_UserBase.Vip
	data.AppId = HF_Atoi(SKDUP_ID)
	data.EventId = SDKUP_EVENT_ID_OFFLINE
	self.SDKLog(url, []byte(HF_JtoA(data)))
}

func (self *Server) SendLog_SDKUP_MoneyChange(player *Player, currency_type string, change_type string, reason string, currency_count int,
	currency_balance int, itmes_id int, items_name string, items_num int) {

	now := time.Now().Unix()

	str := fmt.Sprintf("%s%d", SKDUP_ID, now)
	h := md5.New()
	h.Write([]byte(str))
	md5str := fmt.Sprintf("%x", h.Sum(nil))

	url := fmt.Sprintf("%s/%s/coinChange?x-dbsv1-sign=%s&x-dbsv1-timestamp=%d&x-dbsv1-appId=%s",
		SKDUP_ADDR_URL,
		SKDUP_ID,
		md5str,
		now,
		SKDUP_ID)

	var data SendSDK_MoneyChange
	data.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data.RoleName = player.Sql_UserBase.UName
	data.VipLevel = player.Sql_UserBase.Vip
	data.AppId = HF_Atoi(SKDUP_ID)
	data.EventId = SDKUP_EVENT_ID_MONEY_CHANGE
	data.CurrencyType = currency_type
	data.ChangeType = change_type
	data.Reason = reason
	data.CurrencyCount = currency_count
	data.CurrencyBalance = currency_balance
	data.ItemsId = itmes_id
	data.ItemsName = items_name
	data.ItemsNum = items_num
	self.SDKLog(url, []byte(HF_JtoA(data)))
}

func (self *Server) SendLog_SDKUP_ItemChange(player *Player, change_type int, reason string, itmes_id int, items_type int, items_num int, items_balance int) {

	now := time.Now().Unix()

	str := fmt.Sprintf("%s%d", SKDUP_ID, now)
	h := md5.New()
	h.Write([]byte(str))
	md5str := fmt.Sprintf("%x", h.Sum(nil))

	url := fmt.Sprintf("%s/%s/itemChange?x-dbsv1-sign=%s&x-dbsv1-timestamp=%d&x-dbsv1-appId=%s",
		SKDUP_ADDR_URL,
		SKDUP_ID,
		md5str,
		now,
		SKDUP_ID)

	var data SendSDK_ItemChange
	data.RoleId = fmt.Sprintf("%s", fmt.Sprintf("%d", int(player.Sql_UserBase.Uid)))
	data.RoleName = player.Sql_UserBase.UName
	data.VipLevel = player.Sql_UserBase.Vip
	data.AppId = HF_Atoi(SKDUP_ID)
	data.EventId = SDKUP_EVENT_ID_RESOURCE_CHANGE

	data.ChangeType = change_type
	data.Reason = reason
	data.ItemsId = itmes_id
	data.ItemsType = items_type
	data.ItemsNum = items_num
	data.ItemsBalance = items_balance
	self.SDKLog(url, []byte(HF_JtoA(data)))
}

func (self *Server) SendLog_SDKUP_Online() {

	now := time.Now().Unix()

	str := fmt.Sprintf("%s%d", SKDUP_ID, now)
	h := md5.New()
	h.Write([]byte(str))
	md5str := fmt.Sprintf("%x", h.Sum(nil))

	url := fmt.Sprintf("%s/%s/onlineNumber?x-dbsv1-sign=%s&x-dbsv1-timestamp=%d&x-dbsv1-appId=%s",
		SKDUP_ADDR_URL,
		SKDUP_ID,
		md5str,
		now,
		SKDUP_ID)

	res := GetPlayerMgr().GetOnlineByGameId()
	for gameId, onLineNum := range res {
		var data SendSDK_Online
		data.AppId = HF_Atoi(gameId)
		data.ServerId = self.Con.ServerId
		data.OnlineRoleCount = onLineNum
		data.EventId = SDKUP_EVENT_ID_ONLINE
		data.EventTime = time.Now().Unix()
		data.DistributorId = 0
		self.SDKLog(url, []byte(HF_JtoA(data)))
	}
}

func (self *Server) SendLog_SDKUP_AIWAN_LOGIN(player *Player, eventType string) {

	if player.Account.Channelid != SKDUP_ADDR_URL_AIWAN_SDK_CHANNELID1 &&
		player.Account.Channelid != SKDUP_ADDR_URL_AIWAN_SDK_CHANNELID2 {
		return
	}

	now := time.Now().Unix()

	url := SKDUP_ADDR_URL_AIWAN_SDK

	var data SendSDK_AIWAN
	data.EventType = eventType
	data.AppId = SKDUP_ADDR_URL_AIWAN_SDK_APPID
	data.ServerId = HF_ItoA(player.Account.ServerId)
	if len(player.Account.UserId) >= 5 {
		data.OpenId = player.Account.UserId[5:]
	} else {
		data.OpenId = player.Account.UserId
	}
	data.RoleId = HF_I64toA(player.Sql_UserBase.Uid)
	data.NickName = player.Sql_UserBase.UName
	data.RegTime = int(player.GetRegStampTime())
	data.PostTime = int(now)
	data.Level = player.Sql_UserBase.Level

	ext := make(map[string]int)
	ext["power"] = player.Sql_UserBase.PassMax
	data.Ext = HF_JtoA(ext)
	str := fmt.Sprintf("appid=%s&ext=%s&level=%d&openid=%s&posttime=%d&roleid=%s&serverid=%s%s",
		data.AppId,
		data.Ext,
		data.Level,
		data.OpenId,
		data.PostTime,
		data.RoleId,
		data.ServerId,
		SKDUP_ADDR_URL_AIWAN_SDK_APPKEY)
	h := md5.New()
	h.Write([]byte(str))
	data.Sign = fmt.Sprintf("%x", h.Sum(nil))

	self.SDKLog(url, []byte(HF_JtoA(data)))
}
