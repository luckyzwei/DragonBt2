package main

import (
	"bytes"
	//"compress/zlib"
	"encoding/gob"
	"encoding/json"
	//"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
)

const TIMEFORMAT = "2006-01-02 15:04:05" //! 时间格式化
const POWERMAX = 1500                    //! 最大体力限制
const SKILLPOINTMAX = 100                //! 最大技能点限制
const LEVELMAX = 80                      //! 最大等级限制
const ADDPOWERTIME = 300                 //! 体力回复时间(秒)
const ADDSPTIME = 20                     //! 技能点回复时间(秒)
//const DEFAULT_OPEN_ARTIFACE_LEVEL = 26   //! 开放神器等级
const DEFAULT_WEEKPLAN_MAXCOUNT = 90 //!
const DEFAULT_JJC_FIGHTTIME = 0
const DEFAULT_JJC_WORSHIP_MONEY = 2000
const DEFAULT_JJC_CHG_GEM = 1
const DEFAULT_GOLD = 91000001
const DEFAULT_GEM = 91000002
const DEFAULT_JJC_FIGHT_MAX = 5

const CAMP_SHU = 1
const CAMP_WEI = 2
const CAMP_WU = 3
const CAMP_QUN = 4

const CAMP_CW = 1
const CAMP_BH = 2

//! 表里面的type要+1
const TYPE_CASERN_GONG = 0 //! 弓
const TYPE_CASERN_QI = 1   //! 骑
const TYPE_CASERN_BU = 2   //! 步

const LEVEL_OPEN_TITLE = 6       //! 称号开启等级
const LEVEL_OPEN_CASERN = 32     //! 兵种开启等级
const LEVEL_OPEN_ARTAFACT = 28   //! 神器开启等级
const LEVEL_OPEN_FATE = 21       //! 缘分开启等级
const LEVEL_OPEN_TONGYU = 20     //! 统御开启等级
const PASS_OPEN_TONGYU = 1201100 //! 统御开启关卡

//! 主城
var CITY_SHU int = 1104
var CITY_WEI int = 1059
var CITY_WU int = 1018

// 国王争夺战出生点
var CITY_CW int = 1001
var CITY_BH int = 1010

var CAMP_NAME = []string{"蜀国", "魏国", "吴国"}

const FIGHTTYPE_DEF = -200       //! 城防军
const FIGHTTYPE_YUANZHENG = -300 //! 远征军
const FIGHTTYPE_ROBOT = -500     //! 机器人

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

//! 拷贝指针结构
func HF_CloneType(obj interface{}) interface{} {
	newObj := reflect.New(reflect.TypeOf(obj).Elem()).Elem()
	return newObj.Addr().Interface()
}

//! 克隆对象 dst为指针
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

func HF_MaxInt64(a int64, b int64) int64 {
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

func HF_Atof(s string) float32 {
	num, _ := strconv.ParseFloat(s, 32)
	return float32(num)
}

func HF_Atof64(s string) float64 {
	num, _ := strconv.ParseFloat(s, 64)
	return num
}

func HF_Data2Jsonstr(data interface{}, str *string) {
	s, _ := json.Marshal(data)
	*str = string(s)
}

//! 过滤 emoji 表情
func HF_FilterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			new_content += string(value)
		}
	}
	return new_content
}

//! 得到ip
func HF_GetHttpIP(req *http.Request) string {
	ip := req.Header.Get("Remote_addr")
	if ip == "" {
		ip = req.RemoteAddr
	}
	return strings.Split(ip, ":")[0]
}
