package game

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// 特殊活动Id
const DOUBLE_PASS = 900101000
const DOUBLE_CITY_SEARCH = 900201
const DOUBLE_POWER = 900301

//! 活动定义结构:activitynew.csv
type JS_ActivityType struct {
	Id         int    `json:"id"`          //! 活动Id
	Type       int    `json:"btn_type"`    //! 活动子类型
	Status     int    `json:"status"`      //! 活动状态 0 关闭 1 开启 2 活动展示阶段
	Name       string `json:"name"`        //! 活动名字
	TaskType   int    `json:"tasktype"`    //! 活动任务类型
	Start      string `json:"start"`       //! 活动开始时间
	Continued  int    `json:"continued"`   //! 活动持续时间
	CD         int    `json:"cd"`          //! 活动周期时间
	Show       int    `json:"show"`        //! 活动展示时间,为0表示直接结束
	Renovate   int    `json:"renovate"`    //! 刷新单个活动
	Reset      int    `json:"reset"`       //! 刷新整个活动
	Sort       int    `json:"sort"`        //! 分类,暂时没有什么用
	Mode       int    `json:"mode"`        //! 分类,主要给客户端分类显示
	ModeType   int    `json:"modetype"`    //! 分类,主要给客户端分类显示
	Backicon   string `json:"backicon"`    //! 图片路径
	Nameicon   string `json:"nameicon"`    //! 名字路径
	Dec        string `json:"dec"`         //! 描述路径
	Button     int    `json:"button"`      //! 按钮
	ButtonType int    `json:"buttontype"`  //! 按钮
	Refresh    int    `json:"refresh"`     //! 是否立即刷新
	RefreshMax int    `json:"refresh_max"` //! 刷新次数
	Label      int    `json:"label"`       //! 标签
	EffectShow string `json:"effect_show"` //! 特效
}

// 获取活动总的周期时间
func (self *JS_ActivityType) getTotalTime() int {
	return self.Show + self.Continued + self.CD
}

//! 活动子类型
type JS_ActivityItem struct {
	Id       int    `json:"id"`                 //! 子活动Id
	Step     int    `json:"step"`               //! 活动期数
	Txt      string `json:"txt"`                //! 子活动名字
	N        []int  `json:"n,omitempty"`        //! 子活动任务条件
	Item     []int  `json:"item,omitempty"`     //! 子活动物品Id
	Num      []int  `json:"num,omitempty"`      //! 子活动物品数量
	CostItem []int  `json:"costitem,omitempty"` //! 子活动消耗物品Id
	CostNum  []int  `json:"costnum,omitempty"`  //! 子活动消耗物品数量
}

//! 幸运礼包, activitybox.csv
type JS_ActivityBox struct {
	BoxId        int    `json:"boxid"`        //! 礼包唯一Id
	GearId       int    `json:"gearid"`       //! 分类Id
	ActivityType int    `json:"activitytype"` //! 活动类型
	BoxName      string `json:"boxname"`      //! 宝箱名字
	Sort         int    `json:"sort"`         //! 子礼包类型
	Look         int    `json:"look"`         //! 原价
	Sale         int    `json:"sale"`         //! 实际价格
	Type         int    `json:"type"`         //! 随机类型
	TaskTypes    int    `json:"tasktypes"`    //! 任务类型
	N            [4]int `json:"n"`            //! 任务条件
	Item         [6]int `json:"item"`         //! 物品Id
	Num          [6]int `json:"num"`          //! 物品数量
	Start        string `json:"start"`        //! 开始时间:开服或者实际时间
	Continue     int    `json:"continue"`     //! 持续时间
	NeedLv       int    `json:"needlv"`       //! 需要等级
	NeedVip      int    `json:"needvip"`      //! 需要vip等级,暂时不用
	Print        string `json:"printicon"`    //! 图标类型给客户端用
	Backicon     string `json:"backicon"`     //! 图标给客户端用
	NeedVip1     int    `json:"needvip1"`     //! vip区间
	NeedVip2     int    `json:"needvip2"`     //! vip区间
	MonthCard    int    `json:"monthcard"`    //! 月卡类型
	Times        int    `json:"time"`         //! 次数
	PicType      int    `json:"pictype"`      //! 类型
	StarHero     int    `json:"starhero"`     //! 需要判断的英雄
}

//! 活动结构, san_activitymask
type Sql_ActivityMask struct {
	Id       int    `json:"id"`       //! 活动类型
	Info     string `json:"info"`     //! 活动信息
	Items    string `json:"items"`    //! 奖励道具
	Topfight string `json:"topfight"` //! 战力
	Toplevel string `json:"toplevel"` //! 等级

	info     JS_ActivityType   //! 活动开启时间、 活动状态[第一次从配置表读取]、 活动类型信息
	items    []JS_ActivityItem //! 活动获得道具、消耗道具、任务条件
	topfight []JS_Top          //! 战力榜信息作废
	toplevel []JS_Top          //! 等级榜信息作废
	status   JS_ActivityInfo   //! 活动结束时间、活动状态[动态改变]

	DataUpdate
}

// 子活动状态
type JS_ActivityInfo struct {
	EndTime int64 //! 结束时间
	Status  int   //! 活动状态, 0 活动关闭, 1 活动开启 2 活动结束，可领奖
}

//! 反序列化
func (self *Sql_ActivityMask) Decode() {
	json.Unmarshal([]byte(self.Info), &self.info)
	json.Unmarshal([]byte(self.Items), &self.items)
	json.Unmarshal([]byte(self.Topfight), &self.topfight)
	json.Unmarshal([]byte(self.Toplevel), &self.toplevel)
}

//! 序列化
func (self *Sql_ActivityMask) Encode() {
	self.Info = HF_JtoA(&self.info)
	self.Items = HF_JtoA(&self.items)
	self.Topfight = HF_JtoA(&self.topfight)
	self.Toplevel = HF_JtoA(&self.toplevel)
}

//! 获得某一期活动获得的道具信息等
func (self *Sql_ActivityMask) GetStep(step int) *JS_ActivityItem {
	for i := 0; i < len(self.items); i++ {
		if self.items[i].Step == step {
			return &self.items[i]
		}
	}
	return nil
}

//! 活动管理类
type ActivityMgr struct {
	Sql_Activity map[int]*Sql_ActivityMask //! 配置以及活动状态数据[发生变化的位置:UpdateActivityStatus, ReloadActivity]
	MapActivity  map[int]*JS_ActivityItem  //! 子活动道具、任务条件等信息
	ActivityType []int                     //! 活动类型
	MaskVer      int                       //! 活动版本启动是1, 每次更新之后++
	Locker       *sync.RWMutex             //! 活动锁
	UpdateTime   int64                     //! 更新状态时间, 每小时更新一次活动状态

	//! 幸运礼包配置
	Sql_ActivityBox map[int]*JS_ActivityBox //! 幸运礼包
	BoxGroup        map[int][]int           //! 幸运礼包4个切页的boxId数据
}

// 活动管理器
var actMgrSingleton *ActivityMgr = nil

//! public
func GetActivityMgr() *ActivityMgr {
	if actMgrSingleton == nil {
		actMgrSingleton = new(ActivityMgr)
		actMgrSingleton.MaskVer = 1
		actMgrSingleton.Sql_Activity = make(map[int]*Sql_ActivityMask, 0)
		actMgrSingleton.MapActivity = make(map[int]*JS_ActivityItem, 0)
		actMgrSingleton.UpdateTime = 0
		actMgrSingleton.Locker = new(sync.RWMutex)
		actMgrSingleton.ActivityType = make([]int, 0)
		actMgrSingleton.Sql_ActivityBox = make(map[int]*JS_ActivityBox)
		actMgrSingleton.BoxGroup = make(map[int][]int)
	}

	return actMgrSingleton
}

// 从san_activitymask读取活动相关配置信息
func (self *ActivityMgr) GetData() {
	var sql_activity Sql_ActivityMask
	sql := fmt.Sprintf("select * from san_activitymask")
	res := GetServer().DBUser.GetAllData(sql, &sql_activity)
	// 活动信息初始化
	self.Sql_Activity = make(map[int]*Sql_ActivityMask, 0)
	for i := 0; i < len(res); i++ {
		data := res[i].(*Sql_ActivityMask)
		data.Decode()
		data.Init("san_activitymask", data, false)
		self.Sql_Activity[data.Id] = data
		// 读取道具信息
		for j := 0; j < len(data.items); j++ {
			self.MapActivity[data.items[j].Id] = &data.items[j]
		}
	}

	// 如果san_activitymask为空, 则从配置表读取信息
	if len(self.Sql_Activity) == 0 {
		self.Init()
	}
	// 初始化幸运礼包
	self.InitActivityBox()
	self.MaskVer++
}

// 后来刷新活动信息
func (self *ActivityMgr) UpdateActivity(actMask *Sql_ActivityMask) {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	if actMask != nil {
		actOld, ok := self.Sql_Activity[actMask.Id]
		if !ok {
			InsertTable("san_activitymask", actMask, 0, false)
			actMask.Init("san_activitymask", actMask, false)
			self.Sql_Activity[actMask.Id] = actMask
		} else {
			HF_DeepCopy(&actOld.info, &actMask.info)
			HF_DeepCopy(&actOld.items, &actMask.items)
		}
	}
}

// 读取活动配置表信息
func (self *ActivityMgr) Init() {
	var lst []CsvNode

	//! 主表
	actNew := GetCsvMgr().Data2["Activitynew"]
	for i := 0; i < len(actNew); i++ {
		lst = append(lst, actNew[i])
	}

	actNew1 := GetCsvMgr().Data2["Activitynew1"]
	for i := 0; i < len(actNew1); i++ {
		lst = append(lst, actNew1[i])
	}

	// 读取活动开启时间以及客户端需要的数据
	for i := 0; i < len(lst); i++ {
		activityid := HF_Atoi(lst[i]["id"])
		activitytype := activityid / 100
		sql_activity, ok := self.Sql_Activity[activitytype]
		if !ok {
			sql_activity := new(Sql_ActivityMask)
			sql_activity.Id = activitytype
			sql_activity.info.Id = activitytype
			sql_activity.info.Name = lst[i]["name"]
			sql_activity.info.Backicon = lst[i]["backicon"]
			sql_activity.info.Nameicon = lst[i]["nameicon"]
			sql_activity.info.Start = lst[i]["start"]
			sql_activity.info.CD = HF_Atoi(lst[i]["cd"])
			sql_activity.info.Continued = HF_Atoi(lst[i]["continued"])
			sql_activity.info.Show = HF_Atoi(lst[i]["show"])
			sql_activity.info.Reset = HF_Atoi(lst[i]["reset"])
			sql_activity.info.Sort = HF_Atoi(lst[i]["sort"])
			sql_activity.info.Type = HF_Atoi(lst[i]["type"])
			sql_activity.info.Mode = HF_Atoi(lst[i]["mode"])
			sql_activity.info.ModeType = HF_Atoi(lst[i]["modetype"])
			sql_activity.info.Renovate = HF_Atoi(lst[i]["renovate"])
			sql_activity.info.TaskType = HF_Atoi(lst[i]["tasktypes"])
			sql_activity.info.Dec = lst[i]["dec"]
			sql_activity.info.Status = HF_Atoi(lst[i]["status"])
			sql_activity.info.Refresh = HF_Atoi(lst[i]["refresh"])
			sql_activity.info.RefreshMax = HF_Atoi(lst[i]["refresh_max"])
			sql_activity.info.Label = HF_Atoi(lst[i]["label"])
			sql_activity.info.EffectShow = lst[i]["effect_show"]
			buttonValue := HF_Atoi(lst[i]["button"])
			buttonType := HF_Atoi(lst[i]["buttontype"])
			if buttonValue == 0 {
				buttonValue = 1
			}
			sql_activity.info.Button = buttonValue
			sql_activity.info.ButtonType = buttonType

			sql_activity.items = make([]JS_ActivityItem, 0)
			self.Sql_Activity[activitytype] = sql_activity
			self.ActivityType = append(self.ActivityType, activitytype)
		} else {
			sql_activity.items = make([]JS_ActivityItem, 0)
		}
	}

	// 读取道具、任务条件等信息
	for i := 0; i < len(lst); i++ {
		activityid := HF_Atoi(lst[i]["id"])
		activitytype := activityid / 100
		sql_activity, ok1 := self.Sql_Activity[activitytype]
		if ok1 {
			var activityitem JS_ActivityItem
			activityitem.Id = HF_Atoi(lst[i]["id"])
			activityitem.Step = HF_Atoi(lst[i]["step"])
			activityitem.Txt = lst[i]["txt"]
			activityitem.Item = make([]int, 0)
			activityitem.Num = make([]int, 0)
			activityitem.CostItem = make([]int, 0)
			activityitem.CostNum = make([]int, 0)
			for j := 1; j < 5; j++ {
				activityitem.N = append(activityitem.N, HF_Atoi(lst[i][fmt.Sprintf("n%d", j)]))
				activityitem.CostItem = append(activityitem.CostItem, HF_Atoi(lst[i][fmt.Sprintf("costitem%d", j)]))
				activityitem.CostNum = append(activityitem.CostNum, HF_Atoi(lst[i][fmt.Sprintf("costnum%d", j)]))
			}

			for j := 1; j < 5; j++ {
				activityitem.Item = append(activityitem.Item, HF_Atoi(lst[i][fmt.Sprintf("item%d", j)]))
				activityitem.Num = append(activityitem.Num, HF_Atoi(lst[i][fmt.Sprintf("num%d", j)]))
			}

			sql_activity.items = append(sql_activity.items, activityitem)
			sql_activity.topfight = make([]JS_Top, 0)
			sql_activity.toplevel = make([]JS_Top, 0)
			self.MapActivity[activityid] = &sql_activity.items[len(sql_activity.items)-1]
		}
	}

	LogDebug("初始化活动表", len(self.Sql_Activity))

	// 将活动配置数据、活动状态插入到san_activitymask表
	for _, value := range self.Sql_Activity {
		value.Encode()
		InsertTable("san_activitymask", value, 0, false)
		value.Init("san_activitymask", value, false)
	}
}

// 读取幸运礼包配置
func (self *ActivityMgr) InitActivityBox() {
	lst := GetCsvMgr().ActivityGiftConfig
	for i := 0; i < len(lst); i++ {
		if lst[i].Group == 0 {
			continue
		}
		activitybox := new(JS_ActivityBox)
		activitybox.BoxId = lst[i].ActivityType*100000 + lst[i].Group*100 + lst[i].Index
		activitybox.ActivityType = lst[i].ActivityType
		activitybox.BoxName = lst[i].Name
		activitybox.Look = lst[i].RechargeAmount
		activitybox.Sale = lst[i].Sale
		activitybox.GearId = lst[i].Group
		activitybox.Type = lst[i].Type
		activitybox.Sort = lst[i].Index
		activitybox.TaskTypes = lst[i].TaskTypes
		activitybox.Start = lst[i].Start
		activitybox.Continue = lst[i].Continued

		activity := self.GetActivityUnsafe(lst[i].ActivityType)
		if activity != nil {
			activitybox.Start = activity.info.Start
			activitybox.Continue = activity.info.Continued
		}

		activitybox.NeedLv = 0
		activitybox.NeedVip = 0
		activitybox.Print = lst[i].Pic[0]
		activitybox.Backicon = lst[i].Pic[1]
		activitybox.NeedVip1 = 0
		activitybox.NeedVip2 = 0
		activitybox.MonthCard = 0
		activitybox.Times = lst[i].Times
		activitybox.PicType = lst[i].Pic2Type
		activitybox.StarHero = lst[i].StarHero

		for j := 0; j < 4; j++ {
			activitybox.N[j] = lst[i].N[j]
		}
		for j := 0; j < 6; j++ {
			if j >= len(lst[i].Items) {
				continue
			}

			activitybox.Item[j] = lst[i].Items[j]
			activitybox.Num[j] = lst[i].Nums[j]
		}
		self.Sql_ActivityBox[activitybox.BoxId] = activitybox

		_, ok := self.BoxGroup[activitybox.GearId]
		if !ok {
			self.BoxGroup[activitybox.GearId] = make([]int, 0)
		}

		self.BoxGroup[activitybox.GearId] = append(self.BoxGroup[activitybox.GearId], activitybox.BoxId)
	}
}

// 获取活动所有类型
func (self *ActivityMgr) GetActivityType() []int {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	return self.ActivityType
}

// 根据活动类型获取活动信息
func (self *ActivityMgr) GetActivity(actType int) *Sql_ActivityMask {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	activityType := actType
	if activityType > 100000 {
		activityType = actType / 100
	}

	if act, ok := self.Sql_Activity[activityType]; !ok {
		return nil
	} else {
		return act
	}
}

func (self *ActivityMgr) GetActivityUnsafe(actType int) *Sql_ActivityMask {
	activityType := actType
	if activityType > 100000 {
		activityType = actType / 100
	}

	if act, ok := self.Sql_Activity[activityType]; !ok {
		return nil
	} else {
		return act
	}
}

func (self *ActivityMgr) GetActivityOverflow() *Sql_ActivityMask {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for _, v := range self.Sql_Activity {
		if v.Id >= ACT_OVERFLOW_GIFT_MIN && v.Id <= ACT_OVERFLOW_GIFT_MAX && v.status.Status == 1 {
			return v
		}
	}
	return nil
}

//! 返回道具双倍，经验双倍
func (self *ActivityMgr) GetDoubleStatus(actid int) (int, int) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	value, ok := self.Sql_Activity[actid/100]
	if !ok || value == nil || value.status.Status == ACTIVITY_STATUS_CLOSED {
		//LogDebug("双倍LOG，不存在：", actid)
		return 1, 1
	} else {
		activityitem := self.MapActivity[actid]
		if activityitem == nil {
			//LogDebug("双倍LOG，不存在1：", actid)
			return 1, 1
		}
		if activityitem.N[0] <= 1 {
			//LogDebug("双倍LOG，默认1：", actid)
			return 1, 1
		} else {
			item, exp := 1, 1
			//LogDebug("双倍LOG，双倍道具：", activityitem.N)
			if activityitem.N[1] > 0 {
				//LogDebug("双倍LOG，双倍道具：", actid)
				item = activityitem.N[0]
			}
			if activityitem.N[2] > 0 {
				//LogDebug("双倍LOG，双倍经验：", actid)
				exp = activityitem.N[0]
			}

			return item, exp
		}
	}
}

// 根据活动唯一Id获得活动相关道具
func (self *ActivityMgr) GetActivityItem(actid int) *JS_ActivityItem {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	return self.MapActivity[actid]
}

//! 更新活动状态, force=true表示强制刷新活动状态信息
func (self *ActivityMgr) UpdateActivityStatus(force bool) {
	//每小时0分刷新一次
	tNow := TimeServer()
	// 每天五点定时刷新
	if force == false {
		if self.UpdateTime == 0 || (tNow.Minute() < 5 && tNow.Unix() > self.UpdateTime+360) {
			self.UpdateTime = TimeServer().Unix()
		} else {
			return
		}
	}

	awardActList := make([]*Sql_ActivityMask, 0)

	self.Locker.Lock()

	// 每刷新一次活动版本信息+1
	self.MaskVer++

	// 重新设置活动类型
	self.ActivityType = make([]int, 0)

	for _, value := range self.Sql_Activity {
		self.ActivityType = append(self.ActivityType, value.info.Id)
		// 如果策划填错了,把开始时间填成0或者空
		if value.info.Start == "0" || value.info.Start == "" {
			// 独立充值特殊处理, 不影响其他逻辑,全局变量永远开启,个人逻辑通过注册时间判断
			if value.info.Type == ActSingleRecharge || value.info.Type == ActTimeGift {
				value.info.Status = ACTIVITY_STATUS_OPEN
				value.status.Status = ACTIVITY_STATUS_OPEN
				value.status.EndTime = tNow.Unix() + 864000000
				continue
			}

			if value.info.Type == ActivityOpenFundType {
				value.info.Status = ACTIVITY_STATUS_OPEN
				value.status.Status = ACTIVITY_STATUS_OPEN
				value.status.EndTime = int64(value.info.Continued) + int64(value.info.Show)
				continue
			}

			// 如果持续时间为0, 默认把活动打开且开100天
			if value.info.Continued == 0 || value.info.Start == "0" {
				value.status.Status = ACTIVITY_STATUS_OPEN
				value.status.EndTime = tNow.Unix() + 86400*730
			} else {
				// 否则开服时间加上持续时间比当前时间小就关闭活动
				if GetServer().GetOpenServer()+int64(value.info.Continued) < TimeServer().Unix() {
					value.status.Status = ACTIVITY_STATUS_CLOSED
					value.status.EndTime = 0
				} else {
					// 否则开启活动并且设置结束时间
					value.status.Status = ACTIVITY_STATUS_OPEN
					value.status.EndTime = GetServer().GetOpenServer() + int64(value.info.Continued) + int64(value.info.Show)
				}
			}

			continue
		}

		var starttime int64 = 0
		//! 开放时间为整数时，为开服天数，一次性有效
		startday := HF_Atoi(value.info.Start)
		// 如果是按照开服时间算
		if startday > 0 {
			// 计算开服时间
			starttime = GetServer().GetOpenServer() + int64((startday-1)*DAY_SECS)
			if value.info.Continued == 0 {
				//为0的时候需要判断下开始时间
				if TimeServer().Unix()-starttime >= 0 {
					value.status.Status = ACTIVITY_STATUS_OPEN
				} else {
					value.status.Status = ACTIVITY_STATUS_CLOSED
				}
				continue
			}
		} else if startday < 0 {
			value.info.Status = ACTIVITY_STATUS_OPEN
			value.status.Status = ACTIVITY_STATUS_OPEN
			value.status.EndTime = int64(value.info.Continued) + int64(value.info.Show)
			continue
		} else {
			// 否则按照实际时间算
			t, err := time.ParseInLocation(DATEFORMAT, value.info.Start, time.Local)

			if err != nil {
				starttime = 0
				LogError("开服时间填写错误:", err.Error(), ", 活动类型:", value.info.Type)
				t, err = NewTimeUtil(TimeServer()).Parse(value.info.Start)
				if err != nil {
					LogError("开服时间填写错误:", err.Error(), ", 活动类型:", value.info.Type)
				}
			} else {
				starttime = t.Unix()
			}
		}

		tmp := TimeServer().Unix() - starttime
		// tmp判断活动是否开启
		if tmp <= 0 {
			// 活动没有开启
			value.status.Status = ACTIVITY_STATUS_CLOSED
			value.status.EndTime = 0
		} else {
			// 判断CD逻辑
			if value.info.CD > 0 {
				// 判断活动周期数
				tmpTime := tmp % int64(value.info.getTotalTime())
				// 减去中间间隔+开启间隔 = 活动结束时间
				endtime := tNow.Unix() - tmpTime + int64(value.info.getTotalTime())
				// 计算活动持续时间
				lastTime := tNow.Unix() - tmpTime + int64(value.info.Continued-1)
				// 检查活动是否结束
				if tNow.Unix() <= endtime-1 {
					// 判断活动的持续时间有没有结束
					if tNow.Unix() <= lastTime {
						// 活动持续时间没有结束, 表示活动还开启着
						value.status.Status = ACTIVITY_STATUS_OPEN
						value.status.EndTime = endtime
					} else {
						// 活动持续阶段结束, 判断活动的展示时间,如果展示时间为0,表示活动可以结束了
						if value.info.Show == 0 {
							value.status.Status = ACTIVITY_STATUS_CLOSED
							value.status.EndTime = 0
						} else {
							// show > 0, 表示进入活动的展示阶段
							value.status.Status = ACTIVITY_STATUS_SHOW
							// 设置活动结束时间
							value.status.EndTime = endtime
							// 展示阶段活动回调逻辑
							//self.onActShow(value)
							awardActList = append(awardActList, value)
						}
					}
				} else {
					// 否则活动结束
					value.status.Status = ACTIVITY_STATUS_CLOSED
					value.status.EndTime = 0
				}
			} else {
				// 不可重复刷新
				// 检查活动是否过了持续时间, 没有表示活动开启
				if tNow.Unix() <= starttime+int64(value.info.Continued-1) {
					value.status.Status = ACTIVITY_STATUS_OPEN
					value.info.Status = ACTIVITY_STATUS_OPEN
					value.status.EndTime = starttime + int64(value.info.Continued+value.info.Show)
				} else if tNow.Unix() <= starttime+int64(value.info.Continued+value.info.Show-1) {
					// 判断活动是否在展示阶段, 如果在, 表示活动进入展示阶段
					value.status.Status = ACTIVITY_STATUS_SHOW
					value.info.Status = ACTIVITY_STATUS_SHOW
					value.status.EndTime = starttime + int64(value.info.Continued+value.info.Show)
					//self.onActShow(value)
					awardActList = append(awardActList, value)
				} else {
					// 过了展示阶段, 活动默认结束
					value.status.Status = ACTIVITY_STATUS_CLOSED
					value.info.Status = ACTIVITY_STATUS_CLOSED
					value.status.EndTime = 0
				}
			}
		}
	}

	self.Locker.Unlock() //!解锁
}

// 活动逻辑存盘
func (self *ActivityMgr) Save() {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	for _, value := range self.Sql_Activity {
		// 序列化
		value.Encode()
		// 插入数据库
		value.Update(true)
	}
}

//! 同步全局活动状态给客户端
func (self *ActivityMgr) SendInfo(ver int, player *Player) {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	var msg S2C_ActivityMask
	msg.Cid = "activitymask"
	msg.Ver = self.MaskVer
	// 活动版本发生了变化
	if ver != self.MaskVer {
		// 刷新活动的任务信息
		player.GetModule("activity").(*ModActivity).HandleTask(0, 0, 0, 0)

		for _, value := range self.Sql_Activity {
			activity := GetActivityMgr().GetActivity(value.info.Type)
			if activity == nil {
				continue
			}

			if activity.Id == ACT_WARORDERLIMIT_1 || activity.Id == ACT_WARORDERLIMIT_2 {
				msg.Info = append(msg.Info, value.info)
				// 道具信息
				for j := 0; j < len(value.items); j++ {
					msg.Items = append(msg.Items, value.items[j])
				}
				continue
			}

			startday := HF_Atoi(activity.info.Start)
			if startday < 0 {
				//if value.info.Type == ACT_GROWTH_GIFT {
				//if !player.GetModule("activity").(*ModActivity).isActivityGrowthGiftOpen() {
				//	continue
				//}
				//} else
				if value.info.Type == ACT_REG_TOTAL_DAY {
					if !player.GetModule("activity").(*ModActivity).isActivityRegTotalDayOpen() {
						continue
					}
				} else {
					if !player.GetModule("activity").(*ModActivity).isActivityGiftOpen(value.info.Type) {
						continue
					}
				}
			}
			// 只发送持续、展示阶段的活动
			if value.status.Status > ACTIVITY_STATUS_CLOSED {
				msg.Info = append(msg.Info, value.info)
				// 道具信息
				for j := 0; j < len(value.items); j++ {
					msg.Items = append(msg.Items, value.items[j])
				}
			}
		}
	} else {
		// 否则发空数据给客户端
		msg.Info = make([]JS_ActivityType, 0)
		msg.Items = make([]JS_ActivityItem, 0)
	}
	msg.PassIdRecord = player.GetModule("onhook").(*ModOnHook).GetDailyPassId()

	player.GetModule("activity").(*ModActivity).CheckTaskDone()

	msg.MsgVer = "20180702"
	player.SendMsg("1", HF_JtoB(&msg))
}

//!修改活动开关
func (self *ActivityMgr) ReloadActivity() {
	self.Locker.Lock()
	self.UpdateTime = 0
	self.GetData()
	self.Locker.Unlock()

	//! 更新状态
	self.UpdateActivityStatus(true)
}

// 通过活动获取限时神将keyId
func (self *Sql_ActivityMask) getTaskN4() int {
	if len(self.items) <= 0 {
		return 0
	}

	item := self.items[0]
	if len(item.N) < 4 {
		return 0
	}

	return item.N[3]
}

func (self *Sql_ActivityMask) getTaskN3() int {
	if len(self.items) <= 0 {
		return 0
	}

	item := self.items[0]
	if len(item.N) < 4 {
		return 0
	}

	return item.N[2]
}

// 为0表示没有开启
func (self *ActivityMgr) getActN4(actType int) int {
	activity := self.GetActivity(actType)
	if activity == nil {
		return 0
	}

	n4 := activity.getTaskN4()
	return n4
}

func (self *ActivityMgr) getActN3(actType int) int {
	activity := self.GetActivity(actType)
	if activity == nil {
		return 0
	}

	n3 := activity.getTaskN3()
	return n3
}

// 通过task类型找到活动类型
func (self *ActivityMgr) GetTaskAct(taskType int) *Sql_ActivityMask {
	self.Locker.RLock()
	defer self.Locker.RUnlock()

	for _, value := range self.Sql_Activity {
		if value.status.Status == 0 {
			continue
		}

		if value.info.TaskType != taskType {
			continue
		}

		return value
	}
	return nil
}

func (self *Sql_ActivityMask) getTaskType() int {
	return self.info.TaskType
}

func (self *Sql_ActivityMask) getN4() []int {
	if len(self.items) <= 0 {
		return []int{}
	}

	item := self.items[0]
	if len(item.N) < 4 {
		return []int{}
	}
	return item.N
}

func (self *Sql_ActivityMask) getTxt() string {
	if len(self.items) <= 0 {
		return ""
	}

	return self.items[0].Txt
}

func (self *Sql_ActivityMask) getActTime() int64 {
	var starttime int64 = 0
	//! 开放时间为整数时，为开服天数，一次性有效
	startday := HF_Atoi(self.info.Start)
	//! 如果是按照开服时间算
	if startday > 0 {
		//! 计算开服时间
		starttime = GetServer().GetOpenServer() + int64((startday-1)*DAY_SECS)
	} else {
		//! 否则按照实际时间算
		t, err := time.ParseInLocation(DATEFORMAT, self.info.Start, time.Local)
		if err != nil {

			t, err = NewTimeUtil(TimeServer()).Parse(self.info.Start)
			if err != nil {
				LogError("开服时间填写错误:", err.Error(), ", 活动类型:", self.info.Type)
			}
			starttime = t.Unix()
		} else {
			starttime = t.Unix()
		}
	}
	return starttime
}

func (self *ActivityMgr) JudgeOpenAllIndex(actTypeMin int, actTypeMax int) (bool, int) {

	for i := actTypeMin; i <= actTypeMax; i++ {
		activity := GetActivityMgr().GetActivity(i)
		if activity == nil {
			continue
		}

		startday := HF_Atoi(activity.info.Start)
		if startday > 0 {
			StartTime := GetServer().GetOpenServer() + int64(startday-1)*86400
			EndTime := StartTime + int64(activity.info.Continued) + int64(activity.info.Show)
			now := TimeServer().Unix()
			if now >= StartTime && now <= EndTime {
				return true, i
			}
		} else if startday < 0 {

		} else {
			//! 否则按照实际时间算
			starttime := int64(0)
			t, err := time.ParseInLocation(DATEFORMAT, activity.info.Start, time.Local)
			if err != nil {
				t, err = NewTimeUtil(TimeServer()).Parse(activity.info.Start)
				if err != nil {
					LogError("开服时间填写错误:", err.Error(), ", 活动类型:", activity.info.Type)
				}
			} else {
				starttime = t.Unix()
			}
			endtime := starttime + int64(activity.info.Continued) + int64(activity.info.Show)

			if TimeServer().Unix() > starttime && TimeServer().Unix() < endtime {
				return true, i
			}
		}
	}
	return false, 0
}

func (self *ActivityMgr) JudgeOpenIndex(actTypeMin int, actTypeMax int) (bool, int) {

	for i := actTypeMin; i <= actTypeMax; i++ {
		activity := GetActivityMgr().GetActivity(i)
		if activity == nil {
			continue
		}

		startday := HF_Atoi(activity.info.Start)
		if startday > 0 {
			StartTime := GetServer().GetOpenServer() + int64(startday-1)*86400
			EndTime := StartTime + int64(activity.info.Continued)
			if TimeServer().Unix() >= StartTime && TimeServer().Unix() <= EndTime {
				return true, i
			}
		} else if startday < 0 {

		} else {
			//! 否则按照实际时间算
			starttime := int64(0)
			t, err := time.ParseInLocation(DATEFORMAT, activity.info.Start, time.Local)
			if err != nil {
				t, err = NewTimeUtil(TimeServer()).Parse(activity.info.Start)
				if err != nil {
					LogError("开服时间填写错误:", err.Error(), ", 活动类型:", activity.info.Type)
				}
			} else {
				starttime = t.Unix()
			}
			endtime := starttime + int64(activity.info.Continued)

			if TimeServer().Unix() > starttime && TimeServer().Unix() < endtime {
				return true, i
			}
		}
	}
	return false, 0
}
