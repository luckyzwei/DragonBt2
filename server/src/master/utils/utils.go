package utils

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	MIN_SECS   = 60
	HOUR_SECS  = 3600
	DAY_SECS   = 86400                 //! 每天的秒
	DATEFORMAT = "2006-01-02 15:04:05" // 时间格式化

	LOG_COLOR_NORMAL = 0
	LOG_COLOR_RED    = 1
	LOG_COLOR_YELLOW = 2
	LOG_COLOR_BULE   = 3
	LOG_COLOR_GREEN  = 3
)

//! 阵营定义-三国-蜀魏吴
const (
	CAMP_SHU = 1 // 帝国
	CAMP_WEI = 2 // 联邦
	CAMP_WU  = 3 // 圣堂
	CAMP_QUN = 4 //! 群众
	CAMP_NUM = 3 //! 阵营数量
)

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

//! 得到ip
func HF_GetHttpIP(req *http.Request) string {
	ip := req.Header.Get("Remote_addr")
	if ip == "" {
		ip = req.RemoteAddr
	}
	return strings.Split(ip, ":")[0]
}

//生成随机字符串
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63n(1000)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func HF_Atoi(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}

func HF_AtoI64(s string) int64 {
	num, _ := strconv.ParseInt(s, 10, 64)
	return num
}

func HF_JtoA(v interface{}) string {
	s, err := json.Marshal(v)
	if err != nil {
		LogError("HF_JtoA err:", string(debug.Stack()))
	}
	return string(s)
}

func HF_JtoB(v interface{}) []byte {
	s, err := json.Marshal(v)
	if err != nil {
		LogError("HF_JtoB err:", string(debug.Stack()))
	}
	return s
}

func HF_MaxInt64(a int64, b int64) int64 {
	if a > b {
		return a
	}

	return b
}

func HF_Atof(s string) float32 {
	num, _ := strconv.ParseFloat(s, 32)
	return float32(num)
}

func HF_Atof64(s string) float64 {
	num, _ := strconv.ParseFloat(s, 64)
	return num
}

func HF_Itof64(v int) float64 {

	stringfight := strconv.Itoa(v)
	int64fight, _ := strconv.ParseFloat(stringfight, 64)
	return int64fight
}
func HF_Itoi64(v int) int64 {

	stringfight := strconv.Itoa(v)
	int64fight, _ := strconv.ParseInt(stringfight, 10, 64)
	return int64fight
}

//! 得到一个随机数[0, num-1]
func HF_GetRandom(num int) int {
	if num == 0 {
		defer func() {
			LogError("出现了一个0随机", string(debug.Stack()))
		}()
		return 0
	}
	return rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63n(1000))).Intn(num)
}

//进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

//进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//! 压缩并转码base64
func HF_CompressAndBase64(data []byte) string {
	var buf bytes.Buffer
	compressor := zlib.NewWriter(&buf)
	compressor.Write(data)
	compressor.Close()

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

//! 解码并解压缩
func HF_Base64AndDecompress(data string) []byte {
	defer func() {
		x := recover()
		if x != nil {
			log.Println("HF_Base64AndDecompress:", x, string(debug.Stack()))
			LogDebug("HF_Base64AndDecompress:", string(debug.Stack()))
		}
	}()

	// 对上面的编码结果进行base64解码
	decodeBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		LogDebug("base64 decode err:", err)
	}
	//fmt.Println(string(decodeBytes))

	b := bytes.NewReader(decodeBytes)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)

	return out.Bytes()
}

//! 得到主城id
func HF_GetMainCityID(camp int) int {
	//if camp == constant.CAMP_SHU {
	//	return constant.CITY_SHU
	//} else if camp == constant.CAMP_WEI {
	//	return constant.CITY_WEI
	//}
	//return constant.CITY_WU

	return 0
}

// 随机生成随机数
func HF_RandInt(min, max int) int {
	if max < min {
		LogError("rand error:max < min:", min, max)
	}

	if min >= max || min == 0 || max == 0 {
		return max
	}

	check := max - min
	if check <= 0 {
		check = 1
	}
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	return source.Intn(check) + min
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

//! 检查数据库名字是否重复
//func HF_IsHasName(name string) bool {
//var sql San_UserBase
//db.GetDBMgr().DBUser.GetOneData(fmt.Sprintf("select * from `san_userbase` where `uname` = '%s'", name), &sql, "", 0)
//return sql.Uid > 0
//}

func HF_AbsInt(a int) int {
	if a > 0 {
		return a
	}
	return a * -1
}

//! 克隆对象 dst为指针
func HF_DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func HF_GetWeekTime() int64 {
	now := time.Now()
	if now.Weekday() > time.Monday {
		return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix() + int64((8-int(now.Weekday()))*DAY_SECS)
	} else if now.Weekday() == time.Sunday {
		return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix() + DAY_SECS
	} else {
		if now.Hour() < 5 {
			return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
		} else {
			return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix() + 7*DAY_SECS
		}
	}
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

//! 是否合法
func HF_IsLicitName(name []byte) bool {
	for i := 0; i < len(name); i++ {
		switch name[i] {
		case '\r', '\'', '\n', ' ', '	', '"', '\\':
			return false
		default:
		}
	}

	return true
}

//! 判断是否新的一天
func HF_IsNewDate(checkDate string) bool {
	lct, err := time.ParseInLocation(DATEFORMAT, checkDate, time.Local)
	if err != nil {
		LogError("HF_IsNewDate:", err.Error(), checkDate)
		return false
	}

	checkToday := time.Date(lct.Year(), lct.Month(), lct.Day(), 5, 0, 0, 0, time.Now().Location()).Unix()
	if lct.Hour() >= 5 {
		checkToday += DAY_SECS
	}

	if time.Now().Unix() >= checkToday {
		return true
	} else {
		return false
	}
}

// 获得日志颜色
func HF_GetColorByLog(colorType int) string {
	strColor := ""
	if colorType == LOG_COLOR_RED {
		strColor = "20|FF0000|"
	} else if colorType == LOG_COLOR_YELLOW {
		strColor = "20|FFFF00|"
	} else if colorType == LOG_COLOR_BULE {
		strColor = "20|0080FF|"
	} else if colorType == LOG_COLOR_GREEN {
		strColor = "20|00FF00|"
	} else {
		strColor = "20|acd7ff|"
	}
	return strColor
}

// 获得城池颜色
func HF_GetColorByCamp(camp int) string {
	strColor := ""
	if camp == CAMP_SHU {
		strColor = "40#190#8#"
	} else if camp == CAMP_WEI {
		strColor = "83#140#237#"
	} else if camp == CAMP_WU {
		strColor = "255#31#0#"
	} else {
		strColor = "236#151#30#"
	}
	return strColor
}

//! 克隆对象 dst为指针
func CopyByJson(dst interface{}, src interface{}) error {
	if dst == nil {
		return fmt.Errorf("dst cannot be nil")
	}
	if src == nil {
		return fmt.Errorf("src cannot be nil")
	}
	bytes, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("Unable to marshal src: %s", err)
	}
	err = json.Unmarshal(bytes, dst)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal into dst: %s", err)
	}
	return nil
}

// 得到一个随机数组
func HF_GetRandomArr(arr []int, num int) []int {
	if len(arr) <= num {
		return arr
	}

	lst := make([]int, 0)
	for len(arr) > 0 && len(lst) < num {
		index := HF_GetRandom(len(arr))
		lst = append(lst, arr[index])
		copy(arr[index:], arr[index+1:])
		arr = arr[:len(arr)-1]
	}

	return lst
}

func DeepCopyByJson(dst interface{}, src interface{}) error {
	if dst == nil {
		return fmt.Errorf("dst cannot be nil")
	}
	if src == nil {
		return fmt.Errorf("src cannot be nil")
	}
	bytes, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("unable to marshal src: %s", err)
	}
	err = json.Unmarshal(bytes, dst)
	if err != nil {
		return fmt.Errorf("unable to unmarshal into dst: %s", err)
	}
	return nil
}

//! Team => Arms Type
func HF_Team2ArmsType(team int) int {
	//switch team {
	//case 0:
	//	return constant.TYPE_CASERN_BU
	//case 1:
	//	return constant.TYPE_CASERN_GONG
	//case 2:
	//	return constant.TYPE_CASERN_QI
	//}
	//
	//return constant.TYPE_CASERN_BU

	return 0
}
