package game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

const (
	UPDATE_RANK_TIME_CONSUMERTOP = 300
)

//! 神将结构-本服
type JS_MagicialHero struct {
	Heroid    int    `json:"heroid"`    //! 英雄Id
	Name      string `json:"name"`      //! 名字
	Icon      string `json:"icon"`      //! 图标
	Level     int    `json:"level"`     //! 等级
	HP        int    `json:"hp"`        //! 当前HP
	MaxHP     int    `json:"maxhp"`     //! 最大HP
	Endtime   int64  `json:"endtime"`   //! 本轮结束时间
	MapId     int    `json:"mapid"`     //! 战斗关卡
	Step      int    `json:"step"`      //! 关卡数据
	ShopGroup int    `json:"shopgroup"` //! 商店组
	RankGroup int    `json:"rankgroup"` //! 积分组
}

//! 消费排行榜-跨服-个人
type JS_ConsumerTopUser struct {
	Uid      int64  `json:"uid"`
	SvrId    int    `json:"svrid"`
	SvrName  string `json:"svrname"`
	UName    string `json:"uname"`
	Level    int    `json:"level"`
	Vip      int    `json:"vip"`
	Icon     int    `json:"icon"`
	Point    int    `json:"point"`
	Portrait int    `json:"portrait"` // 边框  20190412 by zy
	Rank     int    `json:"rank"`
	Kill     int    `json:"kill"`
	Step     int    `json:"step"`
	KillAll  int    `json:"killall"` //新增，用来计算服务器数据(更改以前的算法) 20210415 by zy
}

//! 消费排行榜-跨服-服务器
type JS_ConsumerTopServer struct {
	SvrId   int    `json:"svrid"`
	SvrName string `json:"svrname"`
	Rank    int    `json:"rank"`
	Point   int    `json:"point"`
	Kill    int    `json:"kill"`
	Step    int    `json:"step"`
}

//! 消费排行榜-击杀记录-稀有掉落-个人榜第一
type JS_ConsumerMsg struct {
	Uid      int64  `json:"uid"`
	Uname    string `json:"uname"`
	BossName string `json:"bossname"`
	Level    int    `json:"level"`
}

type San_ConsumerTop struct {
	Id      int    `json:"id"`      //! 自增Id
	Start   string `json:"start"`   //! 开始时间
	Show    int    `json:"show"`    //! 领取-显示时间
	Hero    string `json:"hero"`    //! 神将信息
	TopUser string `json:"topuser"` //! 排行榜-个人
	TopSvr  string `json:"topsvr"`  //! 排行榜-服务器
	Record  string `json:"record"`  //! 伤害记录
	Msg     string `json:"msg"`     //! 同步消息-跨服数据
	Attack  string `json:"attack"`  //! 攻击用户
	Expire  int    `json:"expire"`  //! 过期-过期之后，保留数据，不再更新

	hero    JS_MagicialHero         //! 神将结构
	topuser []*JS_ConsumerTopUser   //! 跨服个人排行数据
	topsvr  []*JS_ConsumerTopServer //! 跨服服务器排行数据
	record  []*JS_DamageRecord      //! 伤害记录
	msg     []*JS_ConsumerMsg       //! 保存消息
	attack  map[int64]int           //! 当前攻击用户

	DataUpdate
}

//! 消费排行榜-跨服-个人
type San_ConsumerTopUser struct {
	Uid      int64  `json:"uid"`
	SvrId    int    `json:"svrid"`
	SvrName  string `json:"svrname"`
	UName    string `json:"uname"`
	Level    int    `json:"level"`
	Vip      int    `json:"vip"`
	Icon     int    `json:"icon"`
	Portrait int    `json:"portrait"`
	Point    int    `json:"point"`
	Rank     int    `json:"rank"`
	Step     int    `json:"step"`

	DataUpdate
}

//! 将数据库数据写入data
func (self *San_ConsumerTop) Decode() {
	json.Unmarshal([]byte(self.Hero), &self.hero)
	json.Unmarshal([]byte(self.TopUser), &self.topuser)
	json.Unmarshal([]byte(self.TopSvr), &self.topsvr)
	json.Unmarshal([]byte(self.Record), &self.record)
	json.Unmarshal([]byte(self.Msg), &self.msg)
	json.Unmarshal([]byte(self.Attack), &self.attack)
}

//! 将data数据写入数据库
func (self *San_ConsumerTop) Encode() {
	self.Hero = HF_JtoA(&self.hero)
	self.TopUser = HF_JtoA(&self.topuser)
	self.TopSvr = HF_JtoA(&self.topsvr)
	self.Record = HF_JtoA(&self.record)
	self.Msg = HF_JtoA(&self.msg)
	self.Attack = HF_JtoA(&self.attack)
}

//! 消耗者排行榜
type ConsumerTopMgr struct {
	Sql_Data       *San_ConsumerTop              //! 数据
	Sql_GlobalUser map[int64]*JS_ConsumerTopUser //! 本地服数据
	LastUpdateTime int64                         //! 上次刷新
	Locker         *sync.RWMutex                 //! 数据锁
	PlayerLocker   *sync.RWMutex                 //! 人数锁
	AttackLocker   *sync.RWMutex                 //! 攻击锁
	Ver            int                           //! 版本
	MapPlayer      map[int64]int                 //! 打开界面，需要同步人数
	UpdateTime     int64                         //! 初始化请求时间
	BossFightInfo  *JS_FightInfo
}

var consumertopsingleton *ConsumerTopMgr = nil

func GetConsumerTop() *ConsumerTopMgr {
	if consumertopsingleton == nil {
		consumertopsingleton = new(ConsumerTopMgr)
		consumertopsingleton.Sql_Data = nil
		consumertopsingleton.Sql_GlobalUser = make(map[int64]*JS_ConsumerTopUser)
		consumertopsingleton.MapPlayer = make(map[int64]int)
		consumertopsingleton.Locker = new(sync.RWMutex)
		consumertopsingleton.PlayerLocker = new(sync.RWMutex)
		consumertopsingleton.AttackLocker = new(sync.RWMutex)
		consumertopsingleton.LastUpdateTime = 0
		consumertopsingleton.UpdateTime = TimeServer().Unix() + 30
	}

	return consumertopsingleton
}

//! 定时执行
func (self *ConsumerTopMgr) OnTimer() {
	self.GetData()

	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		if self.UpdateTime < TimeServer().Unix() {
			self.ReqConsumerTopList()
			self.UpdateTime = TimeServer().Unix() + UPDATE_RANK_TIME_CONSUMERTOP
		}
	}

	ticker.Stop()
}

func (self *ConsumerTopMgr) GetData() {
	var consumertop San_ConsumerTop
	sql := fmt.Sprintf("select * from `san_consumertop` where expire = 0  order by id desc limit 1")
	res := GetServer().DBUser.GetAllData(sql, &consumertop)
	for i := 0; i < len(res); i++ {
		data := res[i].(*San_ConsumerTop)
		data.Init("san_consumertop", data, false)
		data.Decode()
		self.Sql_Data = data

		if data.msg == nil {
			data.msg = make([]*JS_ConsumerMsg, 0)
		}
	}

	//!插入测试数据
	if self.Sql_Data == nil {
		//! 检查活动是否过期,活动切换时，初始化
		activity := GetActivityMgr().GetActivity(MAGICHERO_ACTIVITY_ID)
		if activity != nil {
			step := self.GetStep()
			//if len(activity.items) > 0 {
			//	if len(activity.items[0].N) > 0 {
			//		step = activity.items[0].N[0]
			//	}
			//}
			LogInfo("无双神将活动存在，初始化活动：", step)
			self.InitMagicalHero(step)
			self.Sql_Data.hero.Endtime = activity.status.EndTime
		} else {
			self.InitMagicalHero(0)
			if activity != nil {
				self.Sql_Data.hero.Endtime = activity.status.EndTime
			}
		}
	} else {
		//! 检查活动是否过期,活动切换时，初始化
		activity := GetActivityMgr().GetActivity(MAGICHERO_ACTIVITY_ID)
		if activity == nil {
			self.Sql_Data.Expire = 1
			if self.Sql_Data.attack == nil {
				self.Sql_Data.attack = make(map[int64]int)
			}
		} else {
			if activity.status.Status == ACTIVITY_STATUS_OPEN {
				if self.Sql_Data.Expire == 1 {
					step := self.GetStep()
					LogInfo("无双神将活动过期：", step)
					//if len(activity.items) > 0 {
					//	if len(activity.items[0].N) > 0 {
					//		step = activity.items[0].N[0]
					//	}
					//}
					self.InitMagicalHero(step)
					self.Sql_Data.hero.Endtime = activity.status.EndTime
				} else {
					step := self.GetStep()
					LogInfo("无双神将检查期数：", step)
					//if len(activity.items) > 0 {
					//	if len(activity.items[0].N) > 0 {
					//		step = activity.items[0].N[0]
					//	}
					//}

					if step != self.Sql_Data.hero.Step {
						LogInfo("无双神将期数错误，重启活动：", step)
						//! 过期当前
						self.Sql_Data.Expire = 1
						self.Save()

						self.InitMagicalHero(step)
						self.Sql_Data.hero.Endtime = activity.status.EndTime
					}
				}
			}
			//self.Sql_Data.Expire = 0
			self.Sql_Data.hero.Endtime = activity.status.EndTime
		}
	}

	//self.ReqConsumerTopList()
}

func (self *ConsumerTopMgr) GetStep() int {
	step := 0
	activity := GetActivityMgr().GetActivity(MAGICHERO_ACTIVITY_ID)
	if activity != nil {

		if len(activity.items) > 0 {
			if len(activity.items[0].N) == 4 {
				step = activity.items[0].N[3]*1000 + activity.items[0].N[2]
			}
		}
	}

	return step
}

//! 每天早上5点刷新，重置等级和血量
func (self *ConsumerTopMgr) OnFresh() {
	//判断当前的期数是否刷新了
	activity := GetActivityMgr().GetActivity(MAGICHERO_ACTIVITY_ID)
	if activity != nil {
		step := self.GetStep()
		if step != self.Sql_Data.hero.Step {
			LogInfo("无双神将过期，自动刷新：", step)
			//! 过期当前
			self.Sql_Data.Expire = 1
			self.Save()

			self.InitMagicalHero(step)
			self.Sql_Data.hero.Endtime = activity.status.EndTime
		}
	}

	csvhero, ok := GetCsvMgr().Data["Consumetop_Boss"][self.Sql_Data.hero.Step/1000]
	if !ok {
		LogError("找不到神将配置，请检查配置！", self.Sql_Data.hero.Step)
		return
	}
	csvhp, ok1 := GetCsvMgr().Data["Consumetop_Hp"][1]
	if !ok1 {
		LogError("找不到血量配置，请检查配置！")
		return
	}

	if activity := GetActivityMgr().GetActivity(MAGICHERO_ACTIVITY_ID); activity != nil {
		if activity.status.Status == ACTIVITY_STATUS_OPEN {
			step := self.GetStep()
			self.Sql_Data.hero.Heroid = HF_Atoi(csvhero["boss"])
			self.Sql_Data.hero.Name = csvhero["name"]
			self.Sql_Data.hero.Level = 1
			self.Sql_Data.hero.Icon = csvhero["picture"]
			self.Sql_Data.hero.MapId = HF_Atoi(csvhero["mapid"])

			//! 同步服务器列表
			if self.Sql_Data.hero.Step != step {
				self.Sql_Data.hero.Step = step
			}

			self.Sql_Data.hero.MaxHP = HF_Atoi(csvhp["hp"])
			self.Sql_Data.hero.HP = self.Sql_Data.hero.MaxHP
			self.Sql_Data.hero.Endtime = activity.status.EndTime

			self.ReqConsumerTopList()

			self.Sql_Data.Expire = 0
			self.Sql_Data.msg = make([]*JS_ConsumerMsg, 0)

			self.AttackLocker.Lock()
			self.Sql_Data.attack = make(map[int64]int)
			self.AttackLocker.Unlock()
		} else {
			self.Sql_Data.Expire = 1
		}
	}
}

//! 初始化神将系统
func (self *ConsumerTopMgr) InitMagicalHero(step int) {
	self.Sql_Data = new(San_ConsumerTop)
	csvhero, ok := GetCsvMgr().Data["Consumetop_Boss"][step/1000]
	if !ok {
		LogError("找不到神将配置，请检查配置！", step)
		return
	}
	self.Sql_Data.hero.Heroid = HF_Atoi(csvhero["boss"])
	self.Sql_Data.hero.Name = csvhero["name"]
	self.Sql_Data.hero.Level = 1
	self.Sql_Data.hero.Icon = csvhero["picture"]
	self.Sql_Data.hero.MapId = HF_Atoi(csvhero["mapid"])
	self.Sql_Data.hero.ShopGroup = HF_Atoi(csvhero["shop"])
	self.Sql_Data.hero.RankGroup = HF_Atoi(csvhero["list"])

	csvhp, ok1 := GetCsvMgr().Data["Consumetop_Hp"][self.Sql_Data.hero.Level]
	if !ok1 {
		LogError("找不到血量配置，请检查配置！", self.Sql_Data.hero.Level)
		return
	}

	self.Sql_GlobalUser = make(map[int64]*JS_ConsumerTopUser)

	self.Sql_Data.hero.MaxHP = HF_Atoi(csvhp["hp"])
	self.Sql_Data.hero.HP = self.Sql_Data.hero.MaxHP
	self.Sql_Data.hero.Endtime = 0
	self.Sql_Data.hero.Step = step

	self.Sql_Data.topsvr = make([]*JS_ConsumerTopServer, 0)
	self.Sql_Data.topuser = make([]*JS_ConsumerTopUser, 0)
	self.Sql_Data.record = make([]*JS_DamageRecord, 0)
	self.Sql_Data.msg = make([]*JS_ConsumerMsg, 0)
	self.Sql_Data.attack = make(map[int64]int)

	self.Sql_Data.Expire = 0

	self.Sql_Data.Encode()
	InsertTable("san_consumertop", self.Sql_Data, 0, false)
	self.Sql_Data.Init("san_consumertop", self.Sql_Data, false)
}

func (self *ConsumerTopMgr) Save() {
	self.AttackLocker.RLock()
	defer self.AttackLocker.RUnlock()

	self.Sql_Data.Encode()
	self.Sql_Data.Update(true)
}

//! 活动是否有效
func (self *ConsumerTopMgr) IsValid() bool {
	//if self.Sql_Data.Expire == 1 {
	//	return false
	//}
	//
	//if self.Sql_Data.hero.HP <= 0 {
	//	return false
	//}
	actinfo := GetActivityMgr().GetActivity(MAGICHERO_ACTIVITY_ID)
	if actinfo != nil {
		if actinfo.status.Status == ACTIVITY_STATUS_CLOSED {
			return false
		} else {
			return true
		}
	}

	return false
}

//! 是否可以攻击
func (self *ConsumerTopMgr) IsAttackable() bool {
	tNow := TimeServer()
	if tNow.Hour() >= 12 && tNow.Hour() <= 23 {
		return true
	}

	return false
}

func (self *ConsumerTopMgr) GetServerRank(serverid int) int {
	for i := 0; i < len(self.Sql_Data.topsvr); i++ {
		if self.Sql_Data.topsvr[i].SvrId == serverid {
			return self.Sql_Data.topsvr[i].Rank
		}
	}

	return 0
}

func (self *ConsumerTopMgr) AddPlayer(uid int64, remove bool) {
	self.PlayerLocker.Lock()
	defer self.PlayerLocker.Unlock()

	if remove == true {
		delete(self.MapPlayer, uid)
	} else {
		self.MapPlayer[uid] = 1
	}
}

//! 接受伤害
func (self *ConsumerTopMgr) AddDamage(player *Player, damage int, point int, addpoint int) (int, int) {
	//self.Locker.Lock()
	//defer self.Locker.Unlock()
	kill := 0
	rank := 0
	if self.Sql_Data.Expire == 0 {
		if self.Sql_Data.hero.HP > 0 {
			self.Sql_Data.hero.HP -= damage
			if self.Sql_Data.hero.HP <= 0 {
				self.Sql_Data.hero.HP = 0
				kill = 1

				//! 神将升级
				self.Sql_Data.hero.Level += 1
				csvhero, ok := GetCsvMgr().Data["Consumetop_Hp"][self.Sql_Data.hero.Level]
				if !ok {
					self.Sql_Data.Expire = 0
				}

				self.Sql_Data.hero.MaxHP = HF_Atoi(csvhero["hp"])
				self.Sql_Data.hero.HP = HF_Atoi(csvhero["hp"])

				var updatemsg JS_ConsumerMsg
				updatemsg.Uid = player.Sql_UserBase.Uid
				updatemsg.Uname = player.Sql_UserBase.UName
				updatemsg.Level = self.Sql_Data.hero.Level - 1
				updatemsg.BossName = self.Sql_Data.hero.Name

				self.Sql_Data.msg = append(self.Sql_Data.msg, &updatemsg)

				if len(self.Sql_Data.msg) > 10 {
					self.Sql_Data.msg = self.Sql_Data.msg[len(self.Sql_Data.msg)-10:]
				}

				//! 发送击杀辅助奖励
				self.AttackLocker.Lock()
				for uid := range self.Sql_Data.attack {
					if uid == player.Sql_UserBase.Uid {
						continue
					}

					atkplayer := GetPlayerMgr().GetPlayer(uid, true)
					if atkplayer != nil {
						if herocsv, ok := GetCsvMgr().Data["Consumetop_Hp"][self.Sql_Data.hero.Level-1]; ok {
							lstItem := make([]PassItem, 0)
							for j := 3; j < 6; j++ {
								itemid := HF_Atoi(herocsv[fmt.Sprintf("item%d", j+1)])
								if itemid == 0 {
									continue
								}
								lstItem = append(lstItem, PassItem{itemid, HF_Atoi(herocsv[fmt.Sprintf("num%d", j+1)])})
							}
							atkplayer.GetModule("mail").(*ModMail).AddMail(1, 1, 0, GetCsvMgr().GetText("STR_CONSUME_KILL"),
								fmt.Sprintf(GetCsvMgr().GetText("STR_CONSUME_MAIL_CONTENT"), GetConsumerTop().GetMHero().Name), GetCsvMgr().GetText("STR_SYS"), lstItem, true, 0)
						}
					}
				}

				self.Sql_Data.attack = make(map[int64]int)
				self.AttackLocker.Unlock()

				GetServer().sendSysChat(fmt.Sprintf(GetCsvMgr().GetText("STR_CONSUME_CHAT_KILL"),
					HF_GetColorByCamp(player.Sql_UserBase.Camp), player.Sql_UserBase.UName,
					HF_GetColorByCamp(2), self.Sql_Data.hero.Name, self.Sql_Data.hero.Level-1))

			} else {
				self.AttackLocker.Lock()
				self.Sql_Data.attack[player.Sql_UserBase.Uid] = 1
				self.AttackLocker.Unlock()
			}

			//! 同步消息
			var killmsg S2C_ConsumerTopKillBoss
			killmsg.Cid = "consumertopkill"
			killmsg.Boss = self.Sql_Data.hero.Name
			killmsg.Level = self.Sql_Data.hero.Level - 1
			killmsg.Uname = player.Sql_UserBase.UName
			killmsg.HeroId = self.Sql_Data.hero.Heroid
			killmsg.Kill = kill
			killmsg.Damage = damage
			killmsg.BossLevel = self.Sql_Data.hero.Level
			killmsg.HP = self.Sql_Data.hero.HP
			killmsg.MaxHP = self.Sql_Data.hero.MaxHP

			self.BroadCastMsg("consumertopkill", HF_JtoB(&killmsg))

			//GetSessionMgr().BroadCastMsg("consumertopkill", HF_JtoB(&killmsg))

			self.Locker.Lock()
			topuser, ok := self.Sql_GlobalUser[player.Sql_UserBase.Uid]
			if !ok {
				topuser = new(JS_ConsumerTopUser)
				topuser.Uid = player.Sql_UserBase.Uid
				topuser.Rank = 0
				topuser.UName = player.Sql_UserBase.UName
				topuser.Level = player.Sql_UserBase.Level
				topuser.Icon = player.Sql_UserBase.IconId
				topuser.Portrait = player.Sql_UserBase.Portrait
				topuser.Point = point
				topuser.Step = self.Sql_Data.hero.Step
				topuser.SvrId = GetServer().Con.ServerId
				topuser.SvrName = GetServer().Con.ServerName
				topuser.Step = self.Sql_Data.hero.Step

				self.Sql_GlobalUser[topuser.Uid] = topuser
			} else {
				if len(self.Sql_Data.topsvr) > 0 {
					for i := 0; i < len(self.Sql_Data.topsvr); i++ {
						if self.Sql_Data.topsvr[i].SvrId == GetServer().Con.ServerId {
							self.Sql_Data.topsvr[i].Point += addpoint
						}
					}

				}

				rank = topuser.Rank
				topuser.Point = point
				topuser.Step = self.Sql_Data.hero.Step
			}
			self.Locker.Unlock()
			topuser.Icon = player.Sql_UserBase.IconId
			topuser.UName = player.Sql_UserBase.UName
			topuser.Portrait = player.Sql_UserBase.Portrait
			topuser.Kill = kill
			self.UploadDamage(topuser)
		}
	}

	return kill, rank
}

func (self *ConsumerTopMgr) GetTopUserList() []*JS_ConsumerTopUser {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	return self.Sql_Data.topuser
}

func (self *ConsumerTopMgr) GetTopServerList() []*JS_ConsumerTopServer {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	return self.Sql_Data.topsvr
}

func (self *ConsumerTopMgr) GetTopUser(uid int64) *JS_ConsumerTopUser {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for i := 0; i < len(self.Sql_Data.topuser); i++ {
		if self.Sql_Data.topuser[i].Uid == uid {
			return self.Sql_Data.topuser[i]
		}
	}

	return nil
}

func (self *ConsumerTopMgr) GetUserRank(uid int64) int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	topuser, ok := self.Sql_GlobalUser[uid]
	if !ok {
		return 0
	} else {
		return topuser.Rank
	}

}

func (self *ConsumerTopMgr) GetMHero() *JS_MagicialHero {
	return &self.Sql_Data.hero
}

//! 上传战斗数据，每次更新后上传
func (self *ConsumerTopMgr) UploadDamage(top *JS_ConsumerTopUser) {
	res := GetMasterMgr().MatchConsumerTopUpdate(top)
	if res != nil {
		self.Locker.Lock()
		defer self.Locker.Unlock()

		self.Sql_Data.topuser = res.TopUser
		self.Sql_Data.topsvr = res.TopSvr
		if res.SelfInfo != nil {
			self.Sql_GlobalUser[res.SelfInfo.Uid] = res.SelfInfo
		}
	}
}

//! 上传战斗数据，每次更新后上传
func (self *ConsumerTopMgr) ReqConsumerTopList() {
	//!插入测试数据
	if self.Sql_Data == nil {
		//! 检查活动是否过期,活动切换时，初始化
		activity := GetActivityMgr().GetActivity(MAGICHERO_ACTIVITY_ID)
		if activity != nil {
			step := self.GetStep()
			LogInfo("无双神将活动存在，初始化活动：", step)
			self.InitMagicalHero(step)
			self.Sql_Data.hero.Endtime = activity.status.EndTime
		} else {
			self.InitMagicalHero(0)
		}
	} else {
		//! 检查活动是否过期,活动切换时，初始化
		activity := GetActivityMgr().GetActivity(MAGICHERO_ACTIVITY_ID)
		if activity == nil {
			self.Sql_Data.Expire = 1
			if self.Sql_Data.attack == nil {
				self.Sql_Data.attack = make(map[int64]int)
			}
		} else {
			if activity.status.Status == ACTIVITY_STATUS_OPEN {
				if self.Sql_Data.Expire == 1 {
					step := self.GetStep()
					LogInfo("无双神将活动过期：", step)
					self.InitMagicalHero(step)
					self.Sql_Data.hero.Endtime = activity.status.EndTime
				} else {
					step := self.GetStep()
					LogInfo("无双神将检查期数：", step)

					if step != self.Sql_Data.hero.Step {
						LogInfo("无双神将期数错误，重启活动：", step)
						//! 过期当前
						self.Sql_Data.Expire = 1
						self.Save()

						self.InitMagicalHero(step)
						self.Sql_Data.hero.Endtime = activity.status.EndTime
					}
				}
			}
			self.Sql_Data.hero.Endtime = activity.status.EndTime
		}
	}

	res := GetMasterMgr().MatchConsumerTopGetAllRank(self.GetStep(), GetServer().Con.ServerId)
	if res != nil {
		self.Locker.Lock()
		defer self.Locker.Unlock()
		self.Sql_Data.topuser = res.TopUser
		self.Sql_Data.topsvr = res.TopSvr
		for _, v := range self.Sql_Data.topuser {
			if v.SvrId == GetServer().Con.ServerId {
				self.Sql_GlobalUser[v.Uid] = v
			}
		}
	}
}

//! 广播消息
func (self *ConsumerTopMgr) BroadCastMsg(head string, body []byte) {
	self.PlayerLocker.RLock()
	defer self.PlayerLocker.RUnlock()

	var buffer bytes.Buffer
	buffer.Write(HF_DecodeMsg(head, body))
	for uid := range self.MapPlayer {
		player := GetPlayerMgr().GetPlayer(uid, false)
		if player != nil && player.SessionObj != nil {
			player.SessionObj.SendMsgBatch(buffer.Bytes())
		}
	}
}

func (self *ConsumerTopMgr) GetBossFightInfo() *JS_FightInfo {
	if self.BossFightInfo == nil {
		nLen := len(GetCsvMgr().JJCRobotConfig)
		for i := 0; i < nLen; i++ {
			cfg := GetCsvMgr().JJCRobotConfig[i]
			if cfg.Type != MAGICHERO_ACTIVITY_ID {
				continue
			}
			self.BossFightInfo = self.GetRobot(cfg)
			break
		}
	}
	if self.BossFightInfo == nil {
		return nil
	} else {
		data := new(JS_FightInfo)
		HF_DeepCopy(data, self.BossFightInfo)
		return data
	}
}

func (self *ConsumerTopMgr) GetRobot(cfg *JJCRobotConfig) *JS_FightInfo {
	bossConfig := GetCsvMgr().GetActivityBossConfig(MAGICHERO_ACTIVITY_ID, 0)
	if bossConfig == nil {
		return nil
	}

	data := new(JS_FightInfo)
	data.Rankid = 0
	data.Uid = 0
	data.Uname = bossConfig.Name
	data.Iconid = 10000000 + bossConfig.BossId

	data.Portrait = 1000 //机器人边框  20190412 by zy
	data.Level = cfg.Level
	data.Defhero = make([]int, 0)
	data.Heroinfo = make([]JS_HeroInfo, 0)
	data.HeroParam = make([]JS_HeroParam, 0)

	if cfg.Category == 0 {
		data.Camp = HF_GetRandom(3) + 1
	} else {
		data.Camp = cfg.Category
	}

	var heroes []int
	for i := 0; i < len(cfg.Hero); i++ {
		if cfg.Hero[i] == 0 {
			continue
		}

		heroes = append(heroes, cfg.Hero[i])
	}

	heroes = HF_GetRandomArr(heroes, 5)

	num := int(cfg.Fight[0] - cfg.Fight[1])
	if num <= 0 {
		num = 2
	}
	data.Deffight = cfg.Fight[1] + int64(HF_RandInt(1, num))

	for i := 0; i < len(heroes); i++ {
		pos := i + 1 //相当于KEY
		fightPos := 0
		for j := 0; j < len(bossConfig.Position); j++ {
			if bossConfig.Position[j] == heroes[i] {
				fightPos = j
			}
		}
		data.FightTeamPos.FightPos[fightPos] = pos
		data.Defhero = append(data.Defhero, pos)
		var hero JS_HeroInfo
		hero.Heroid = heroes[i]
		hero.Color = cfg.NpcQuality
		hero.HeroKeyId = pos
		hero.Stars = cfg.NpcStar[i]
		hero.HeroQuality = cfg.NpcStar[i]
		hero.Levels = cfg.NpcLv[i]
		hero.Skin = 0
		hero.Soldiercolor = 6
		hero.Skilllevel1 = 0
		hero.Skilllevel2 = 0
		hero.Skilllevel3 = 0
		hero.Skilllevel4 = SKILL_LEVEL
		hero.Skilllevel5 = 0
		hero.Skilllevel6 = 0
		hero.Fervor1 = 0
		hero.Fervor2 = 0
		hero.Fervor3 = 0
		hero.Fervor4 = 0
		hero.Fight = data.Deffight / 5

		hero.ArmsSkill = make([]JS_ArmsSkill, 0)
		hero.TalentSkill = []Js_TalentSkill{}
		hero.MainTalent = 0

		config := GetCsvMgr().HeroBreakConfigMap[hero.Heroid]
		if config == nil {
			continue
		}
		HeroBreakId := 0
		//计算突破等级
		for _, v := range config {
			if hero.Levels >= v.Break {
				HeroBreakId = v.Id
			}
		}

		skillBreakConfig := GetCsvMgr().HeroBreakConfigMap[hero.Heroid][HeroBreakId]
		if skillBreakConfig == nil {
			continue
		}

		for i := 0; i < len(skillBreakConfig.Skill); i++ {
			if skillBreakConfig.Skill[i] > 0 {
				hero.ArmsSkill = append(hero.ArmsSkill, JS_ArmsSkill{Id: skillBreakConfig.Skill[i] / 100, Level: skillBreakConfig.Skill[i] % 100})
			}
		}
		att, att_ext, energy := hero.CountFight(cfg.BaseTypes, cfg.BaseValues)

		data.Heroinfo = append(data.Heroinfo, hero)
		var param JS_HeroParam
		param.Heroid = hero.Heroid
		param.Param = att
		param.ExtAttr = att_ext
		param.Hp = param.Param[AttrHp]
		param.Energy = energy
		param.Energy = 0

		data.HeroParam = append(data.HeroParam, param)

	}
	return data
}
