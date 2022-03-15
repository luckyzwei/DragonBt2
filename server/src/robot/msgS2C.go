package main

type S2C_MsgId struct {
	Cid string `json:"cid"`
}

//! 注册
type S2C_Reg struct {
	Cid      string `json:"cid"`
	Uid      int64  `json:"uid"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Creator  string `json:"creator"`
}

//! 基本信息
type S2C_UserBaseInfo struct {
	Cid      string           `json:"cid"`
	Baseinfo Son_UserBaseInfo `json:"baseinfo"`
}

type Son_UserBaseInfo struct {
	Checkinaward       int    `json:"checkinaward"`
	Checkinnum         int    `json:"checkinnum"`
	Exp                int    `json:"exp"`
	Face               int    `json:"face"`
	Gem                int    `json:"gem"`
	Gold               int    `json:"gold"`
	Iconid             string `json:"iconid"`
	Ischeckin          bool   `json:"ischeckin"`
	Isrename           int    `json:"isrename"`
	Juqingid           int    `json:"juqingid"`
	Juqingid2          int    `json:"juqingid2"`
	Lastcheckintime    string `json:"lastcheckintime"`
	Lastlivetime       string `json:"lastlivetime"`
	Lastlogintime      string `json:"lastlogintime"`
	Level              int    `json:"level"`
	Levelaward         int    `json:"levelaward"`
	Loginaward         int    `json:"loginaward"`
	Logindays          int    `json:"logindays"`
	Morale             int    `json:"morale"`
	Partyid            int    `json:"partyid"`
	Position           int    `json:"position"`
	Regtime            string `json:"regtime"`
	Skillpoint         int    `json:"skillpoint"`
	Splastupdatatime   int64  `json:"splastupdatatime"`
	Tili               int    `json:"tili"`
	Tililastupdatatime int64  `json:"tililastupdatatime"`
	Uid                int64  `json:"uid"`
	Uname              string `json:"uname"`
	Vip                int    `json:"vip"`
	Vipexp             int    `json:"vipexp"`
	Worldaward         int    `json:"worldaward"`
	Zhiyinid           int    `json:"zhiyinid"`
	Zhiyinid1          int    `json:"zhiyinid1"`
	Citylevel          int    `json:"citylevel"`
	Camp               int    `json:"camp"`
	City               int    `json:"city"`
	Day                int    `json:"day"`
	Promotebox         int    `json:"promotebox"`
	OpenServer         int64  `json:"openserver"`
	Channelid          string `json:"channelid"`
	Account            string `json:"account"`
	FitServer          int    `json:"fitserver"`
	Soul               int    `json:"soul"`       //! 魂石
	TechPoint          int    `json:"techpoint"`  //! 科技点
	BossMoney          int    `json:"bossmoney"`  //! 水晶币
	TowerStone         int    `json:"towerstone"` //! 镇魂石
	Portrait           int    `json:"portrait"`   //! 头像框
	CampOk             int    `json:"camp_ok"`    // 阵营ok
	NameOk             int    `json:"name_ok"`    // 名字ok
}

type JS_PassItem struct {
	ItemID int `json:"itemid"`
	Num    int `json:"num"`
}

type S2C_OnItem struct {
	Cid     string     `json:"cid"`
	Itemlst []PassItem `json:"itemlst"`
}

type PassItem struct {
	ItemID int `json:"itemid"` // 道具ID
	Num    int `json:"num"`    // 道具数量
}

//! 发送结果3
type S2C_Result3Msg struct {
	Cid string `json:"cid"`
	Ret bool   `json:"ret"`
}

//! 背包
type S2C_BagInfo struct {
	Cid    string        `json:"cid"`
	Baglst []JS_PassItem `json:"baglst"`
}

//! 同步行动力
type S2C_TeamPower struct {
	Cid       string `json:"cid"`
	Power     int    `json:"power"`
	PowerTime int64  `json:"powertime"`
}

type JS_City struct {
	Id     int `json:"id"`
	Camp   int `json:"camp"`   //! 当前势力
	Attack int `json:"attack"` //! 正在被攻击
	DefNum int `json:"defnum"` //! 城防军数量

	line []int //! 连接城池
}

//! 城池状态
type S2C_CityInfo struct {
	Cid   string    `json:"cid"`
	Level int       `json:"level"`
	Info  []JS_City `json:"info"`
}

//! 登陆成功
type S2C_LoginRet struct {
	Cid        string `json:"cid"`
	Ret        int    `json:"ret"`
	CheckCode  string `json:"checkcode"`
	Servertime int64  `json:"servertime"`
}

//! 更新经验
type S2C_UpdateExp struct {
	Cid    string `json:"cid"`
	Old    int    `json:"old"`
	New    int    `json:"new"`
	Newexp int    `json:"newexp"`
}

//! 更新体力
type S2C_UpdateTiLi struct {
	Cid  string `json:"cid"`
	Tili int    `json:"tili"`
	Time int    `json:"time"`
	Uid  int64  `json:"uid"`
}

//! 发送抽奖消息
type S2C_FindOK struct {
	Cid         string     `json:"cid"`
	FindType    int        `json:"findtype"`
	FreeTime    int64      `json:"freetime"`    //! 免费倒计时
	HgFreeTime  int64      `json:"hgfreetime"`  //! 后宫免费倒计时
	BoxFreeTime int64      `json:"boxfreetime"` //! 后宫免费倒计时
	Num         int        `json:"num"`         //! 神将招募多少次
	MJNum       int        `json:"mjnum"`       //! 名将招募次数
	GoldNum     int        `json:"goldnum"`
	GemNum      int        `json:"gemnum"`
	Item        []PassItem `json:"item"`
	Cost        []PassItem `json:"cost"`
	Gem         int        `json:"gem"`
	BoxNum      int        `json:"boxnum"`
	BeautyNum   int        `json:"hgnum"`
	Gold        int        `json:"gold"`
	YuFu        int        `json:"yufu"`
	GoldEndTime int64      `json:"goldendtime"`  // 免费道具抽结束时间
	LootEnergy  int        `json:"loot_energy"`  // 高级召唤能量
	SummonTimes int        `json:"summon_times"` // 高级抽奖次数
}
