package game

import (
	//	"log"
	//"fmt"
	//"sort"
	//"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

//!订单表
type San_Recharge struct {
	Id        int
	Order     string //! 订单号
	Uid       int64  //! 角色Id
	Account   string //! 帐号Id
	Type      int    //! 类型
	Money     int    //! 金额
	Timestamp int64  //! 时间
	Flag      int    //! 处理标志
	Dealtime  int64  //! 处理时间
	Level     int    //! 角色等级

	DataUpdate
}

func (self *San_Recharge) Decode() { //! 将数据库数据写入data

}

func (self *San_Recharge) Encode() { //! 将data数据写入数据库

}

type RechargeMgr struct {
	Locker       *sync.RWMutex         //! 数据锁
	Sql_Recharge map[int]*San_Recharge //! 订单记录
	MaxOrderId   int                   //! 最大处理订单Id
	WaitOrder    int                   //! 等待处理订单
	StartTime    int64
	NeedUpdate   bool
}

var rechargemgrsingleton *RechargeMgr = nil

func GetRechargeMgr() *RechargeMgr {
	if rechargemgrsingleton == nil {
		rechargemgrsingleton = new(RechargeMgr)
		rechargemgrsingleton.Locker = new(sync.RWMutex)
		rechargemgrsingleton.Sql_Recharge = make(map[int]*San_Recharge)
		rechargemgrsingleton.WaitOrder = 0
		rechargemgrsingleton.StartTime = TimeServer().Unix()
	}

	return rechargemgrsingleton
}

//! 获取数据
func (self *RechargeMgr) GetData(all bool) {
	//LogDebug("RechargeMgr:GetData()")
	self.WaitOrder = 0
	var order San_Recharge
	var sql string
	if all == false {
		sql = fmt.Sprintf("select * from `san_recharge` where flag = 0")
	} else {
		sql = fmt.Sprintf("select * from `san_recharge` where flag > 1")
	}

	res := GetServer().DBUser.GetAllData(sql, &order)
	if res == nil {
		return
	}
	for i := 0; i < len(res); i++ {
		data := res[i].(*San_Recharge)
		data.Init("san_recharge", data, false)
		data.Decode()
		//data.locker = new(sync.RWMutex)
		self.Sql_Recharge[data.Id] = data
		self.WaitOrder++
	}
}

//! 3秒从数据库返回一次
func (self *RechargeMgr) Run() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			self.Fresh()
		}
	}
}

//! 同步客户端
func (self *RechargeMgr) Fresh() {
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

	//!开服120秒内不计算
	passtime := TimeServer().Unix() - self.StartTime
	if passtime < 30 {
		return
	}
	if passtime%5 > 0 && !self.NeedUpdate {
		return
	}
	if passtime%60 == 0 {
		self.GetData(true)
	} else {
		self.GetData(false)
	}
	if self.NeedUpdate==true{
		//LogDebug("后台触发RechargeMgr.Fresh()")
		self.NeedUpdate = false
	}
	if self.WaitOrder == 0 {
		return
	}

	for _, value := range self.Sql_Recharge {
		if value.Flag == 0 || value.Flag == 2 || value.Flag == 10 {
			virtualrecharge := false
			if value.Flag == 10 {
				virtualrecharge = true
			}
			player := GetPlayerMgr().GetPlayer(value.Uid, false)
			if player != nil {
				if player.IsOnline() {
					if value.Type > 1000 {
						value.Flag = player.GetModule("recharge").(*ModRecharge).RechargeEx(value.Type, value.Id, 1, value.Money, value.Level)
					} else {
						value.Flag = player.GetModule("recharge").(*ModRecharge).Recharge(value.Type, value.Id, 1, value.Money)
					}

					if player.SessionObj == nil {
						//! 未连接直接保存
						LogDebug("充值成功，保存")

						player.Save(false, true)
					} else {
						//! 下次收发消息存档
						LogDebug("充值成功，准备保存")
						player.MsgWaitSave = 500
						player.SaveTimes = 3
					}
					if virtualrecharge {
						value.Flag = -10
					}
					//value.Flag = 1
					value.Dealtime = TimeServer().Unix()
				} else {
					if value.Flag == 0 {
						value.Flag = 2
						value.Dealtime = TimeServer().Unix()
					}
				}
				value.Update(true)
			} else {
				if value.Flag == 0 {
					value.Flag = 2
					value.Dealtime = TimeServer().Unix()
					value.Update(true)
				}
				//value.Flag = 2
				//value.Dealtime = TimeServer().Unix()
			}

		}
	}
}
