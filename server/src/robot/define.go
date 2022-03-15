package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type AccountInfo struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

//! 配置
type RobotConfig struct {
	RoboTtype     int           `json:"robottype"`     //! 机器人类型
	IsLoop        bool          `json:"isloop"`        //! 是否循环
	MaxNum        int           `json:"maxnum"`        //! 最大机器人数量
	Serurl        string        `json:"serurl"`        //! 服务器地址
	GuoZhanCityId int           `json:"guozhancityid"` //! 国战指定ID
	Camp          int           `json:"camp"`          //! 机器人国家属性
	ServerId      int           `json:"serverid"`      //! 服务器ID
	Players       []AccountInfo `json:"players"`       //! 帐号列表
	PowerCheckNum int           `json:"PowerCheckNum"` //! 机器人购买军令检查数量
	Origin        string        `json:"origin"`        //! 服务器地址
	ActionDis     int           `json:"actiondis"`     //! 在线机器人操作间隔
}

type RobotFight struct {
	Owner *Robot

	Baglst        []JS_PassItem
	Power         int
	CampFightCity int
	fightinfo     [3]*Son_CampFightInfo
	CityInfo      []JS_City

	State      int
	Stage      int
	Timer      int
	TotalTimer int
	Wins       int
}

type UserBase struct {
	Gold  int //! 金币
	Gem   int //! 钻石
	Power int //!体力
	Level int //!等级
	Camp  int
}

type Robot struct {
	ID            int64
	Ws            *websocket.Conn
	Dead          bool
	IsLoginOut    bool  //是否在登出中
	OperTime      int64 //! 执行下一次操作的时间
	Uid           int64
	Uname         string
	Account       string //! 账号
	Password      string //! 密码
	ServerId      int
	userbase      UserBase
	msgId         int
	csvId         int
	guildId       int
	is_fightinfo  bool
	fightinfo     [3]*Son_CampFightInfo
	SoloId        int
	CampFightCity int //! 国战开启城市

	MsgSN     int64
	RobotType int
	ServerUrl string
	origin    string

	logic IRobotLogic
}

const IS_ROBOT = true
const robot_type_1 = 1 //常规固定流程执行
const robot_type_2 = 2 //根据CSV表流程执行
const robot_type_3 = 3 //国战机器人
const robot_type_5 = 5 //在线机器人

const (
	Sendlogin         int = 1  //登录
	Passmission           = 2  //通过任务
	Passlevelbeggin       = 3  //开始通关
	Passlevelend          = 4  //通关完成
	Getfriendcommend      = 5  //获取好友推荐
	Sendfriendapply       = 6  //一键申请好友
	Sendfriendorder       = 7  //一键通过好友
	Chatsystem            = 8  //世界聊天
	DanMusystem           = 9  //弹幕
	Createitem            = 10 //创建获取道具
	GmPassid              = 11 //GM通关一章
	ShenqiLvUp            = 12 //武将神器升级
	HeroColorUp           = 13 //武将进阶
	HeroSatrUp            = 14 //武将升星
	SoldierLvUp           = 15 //武将统兵升级
	LuckDraw              = 16 //抽奖系统
	CampFightMove         = 17 //参加国战
	CampFightPlayList     = 18 //选择单挑玩法
	CampFightSolo2        = 19 //单挑选择对手
	CampFightSoloEnd      = 20 //发送单挑结果
	CampFight56Req        = 21 //发送报名国战决战
	HeroLvUp              = 22 //发送武将升级
	SendLoginOut          = 23 //登出
	ConsumerTopAttack     = 24 //无双神将攻击
)

//! 解消息
func HF_EncodeMsg(msg []byte) (string, []byte, bool) {
	data := &MsgBase{}
	err := proto.Unmarshal(msg, data)
	if err != nil {
		log.Println(err)
		return "", []byte(""), false
	}

	return data.GetMsghead(), data.GetMsgdata(), true
}

//! 加密消息
func HF_DecodeMsg(msghead string, msgdata []byte) []byte {
	data := &MsgBase{
		Msghead: proto.String(msghead),
		Msgtime: proto.Int64(time.Now().Unix()),
		Msgsign: proto.String("111"),
		Msgdata: []byte(msgdata),
	}

	msg, err := proto.Marshal(data)
	if err != nil {
		log.Println(err)
		return []byte("")
	}

	return msg
}

//! dst为指针
func HF_DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

//! int取最小
func HF_MinInt(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

//! int取最大
func HF_MaxInt(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

//! 得到一个随机数
func HF_GetRandom(num int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
}

func HF_Atoi(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}

func HF_Atoi64(s string) int64 {
	num, _ := strconv.ParseInt(s, 10, 64)
	return num
}

func HF_JtoA(v interface{}) string {
	s, _ := json.Marshal(v)
	return string(s)
}

func HF_JtoB(v interface{}) []byte {
	s, _ := json.Marshal(v)
	return s
}

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}
