package game

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
)

//! 解消息
func HF_EncodeMsg(msg []byte) (string, []byte, bool) {
	data := &MsgBase{}
	err := proto.Unmarshal(msg, data)
	if err != nil {
		log.Println("HF_EncodeMsg:", err)
		return "", []byte(""), false
	}

	return data.GetMsghead(), data.GetMsgdata(), true
}

//! 加密消息
func HF_DecodeMsg(msghead string, msgdata []byte) []byte {
	timeStart := TimeServer()
	data := &MsgBase{
		Msghead: proto.String(msghead),
		Msgtime: proto.Int64(TimeServer().Unix()),
		Msgsign: proto.String("111"),
		Msgdata: []byte(msgdata),
	}

	msg, err := proto.Marshal(data)
	if err != nil {
		log.Println("HF_DecodeMsg:", err)
		return []byte("")
	}
	var buf bytes.Buffer
	compressor := zlib.NewWriter(&buf)
	compressor.Write(msg)
	compressor.Close()

	size := []byte(fmt.Sprintf("%d ", len(msg)))
	newmsg := []byte("")
	newmsg = append(newmsg, size...)
	newmsg = append(newmsg, buf.Bytes()...)

	LogDebug("message zlib cost: ", len(msgdata), buf.Len(), time.Since(timeStart), ", cid:", msghead)
	return newmsg
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

func HF_MinInt64(a int64, b int) int64 {
	if a < int64(b) {
		return int64(a)
	}

	return int64(b)
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

//! 得到一个随机数[0, num-1]
var randSource *rand.Rand
var randCount int

func HF_GetRandom(num int) int {
	if num == 0 {
		defer func() {
			LogError("出现了一个0随机", string(debug.Stack()))
		}()
		return 0
	}

	randCount++
	nowTime := time.Now().UnixNano()
	if randSource == nil || randCount > 2000 {
		randCount = 0
		randSource = rand.New(rand.NewSource(nowTime))
	}

	return (randSource.Intn(1000000)) % num
}

//! Team => Arms Type
func HF_Team2ArmsType(team int) int {
	switch team {
	case 0:
		return TYPE_CASERN_BU
	case 1:
		return TYPE_CASERN_GONG
	case 2:
		return TYPE_CASERN_QI
	}

	return TYPE_CASERN_BU
}

//! 系统时间的偏移
var ServerTimeOffset int64 = 0

func TimeServer() time.Time {
	return time.Unix(time.Now().Unix()+ServerTimeOffset, 0)
}

// 获得下周刷新时间
func HF_GetWeekTime() int64 {
	now := TimeServer()
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

//生成随机字符串
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(TimeServer().UnixNano() + rand.Int63n(1000)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//扩展支持各种格式
func HF_GetStartTime(start string) int64 {
	startday := HF_Atoi(start)
	startTime := GetServer().GetOpenServer() + int64(startday-1)*86400
	return startTime
}

func HF_Atoi(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}

func HF_ItoA(s int) string {
	num := strconv.Itoa(s)
	return num
}

func HF_I64toA(s int64) string {
	num := strconv.FormatInt(s, 10)
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

//! 检查数据库名字是否重复
func HF_IsHasName(name string) bool {
	var sql San_UserBase
	GetServer().DBUser.GetOneData(fmt.Sprintf("select * from `san_userbase` where `uname` = '%s'", name), &sql, "", 0)
	return sql.Uid > 0
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
	info := strings.Split(ip, ":")
	if len(info) <= 0 {
		return "127.0.0.1"
	}
	return info[0]
}

//! 从掉落包里选出一个
func HF_DropForItemBag(bag int) (int, int) {
	csv_bag, ok := GetCsvMgr().ItemBagMap[bag]
	if !ok {
		LogError("bag =", bag, " config error!")
		return 0, 0
	}

	itemid := make([]int, 0)
	itemnum := make([]int, 0)
	itempro := make([]int, 0)
	totalpro := 0
	for i := 0; i < len(csv_bag.ItemIds); i++ {
		id := csv_bag.ItemIds[i]
		if id == 0 {
			continue
		}
		itemid = append(itemid, id)
		if csv_bag.Nums[i] == 0 {
			continue
		}
		itemnum = append(itemnum, csv_bag.Nums[i])

		pro := csv_bag.Weights[i]
		itempro = append(itempro, pro)
		totalpro += pro
	}

	if totalpro == 0 {
		LogError("totalpro == 0, bag = ", bag)
		return 0, 0
	}

	cur := HF_GetRandom(totalpro)
	curpro := 0
	for i := 0; i < len(itempro); i++ {
		curpro += itempro[i]
		if cur < curpro {
			return itemid[i], itemnum[i]
		}
	}

	return 0, 0
}

//! 从宝箱掉落组中选出掉落
func HF_DropForGemsweeperGroup(groupid int) []PassItem {
	var outitem []PassItem
	pro := HF_GetRandom(GetCsvMgr().Gemsweeper_itembag_SUM[groupid])
	cur := 0
	for i := 0; i < 100; i++ {
		config, ok := GetCsvMgr().GetGemsweeper_itembag_GG(groupid, i+1)
		if !ok {
			LogError("groupId:", groupid, ", i+1:", i+1, " not exist!")
			return outitem
		}

		if len(config) <= 0 {
			continue
		}

		cur += HF_Atoi(config[0]["gvalue"])
		if pro < cur {
			sum := 0
			for _, value := range config {
				sum += HF_Atoi(value["value"])
			}
			pro_in := HF_GetRandom(sum)
			cur_in := 0
			for _, value := range config {
				cur_in += HF_Atoi(value["value"])
				if pro_in < cur_in {
					var item PassItem
					item.ItemID, item.Num = HF_Atoi(value["itemid"]), HF_Atoi(value["count"])
					if item.ItemID != 0 {
						outitem = append(outitem, item)
					}
					return outitem
				}
			}
			break
		}
	}

	return outitem
}

//! 从远征BUFF组里面 选一个BUFFid
func HF_GetBuffIdFromBuffGroup(groupid int) int {
	csv, ok := GetCsvMgr().GetExpeditionBuffGroup(groupid)
	if !ok {
		LogDebug("csv没有数据", groupid)
		return 0
	}
	allweight := 0
	for i := 0; i < len(csv.Weights); i++ {
		allweight += csv.Weights[i]
	}

	pro := HF_GetRandom(allweight)
	cur := 0
	for i := 0; i < len(csv.Weights); i++ {
		cur += csv.Weights[i]
		if pro < cur {
			return csv.Ids[i]
		}
	}
	return 0
}

//! 从掉落组里选出若干
func HF_DropForItemBagGroup(groupid int) []PassItem {
	return GetCsvMgr().DropItem(groupid)
}

func HF_GetLevelBoxItem(id int) []PassItem {
	var outitem []PassItem
	csv, ok := GetCsvMgr().LevelboxConfig[id]
	if !ok {
		LogError("id:", id, "levelboxItem not found!")
		return outitem
	}

	for i := 0; i < len(csv.Items); i++ {
		itemid := csv.Items[i]
		if itemid == 0 {
			continue
		}
		num := csv.Nums[i]
		if num == 0 {
			continue
		}
		outitem = append(outitem, PassItem{itemid, num})
	}
	return outitem
}

//! 得到主城id
func HF_GetMainCityID(camp int) int {
	if camp == CAMP_SHU {
		return CITY_SHU
	} else if camp == CAMP_WEI {
		return CITY_WEI
	}
	return CITY_WU
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

// 重写生成连接池方法
func NewPool(ip string, db int, auth string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   12000, // max number of connections
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ip)
			if err != nil {
				panic(err.Error())
			}
			if auth != "" {
				c.Do("AUTH", auth)
			}
			c.Do("SELECT", db)
			return c, err
		},
	}
}

//获取当前队伍战斗力
func HF_GetCurTeamFightNum(playerUID int64) int64 {
	player := GetPlayerMgr().GetPlayer(playerUID, true)
	if player == nil {
		return 0
	}
	return player.Sql_UserBase.Fight
}

//获取总队伍战斗力
func HF_GetAllTeamFightNum(playerUID int64) int64 {
	return HF_GetCurTeamFightNum(playerUID)
}
func HF_GetHeroType(heroid int) string {
	heroConfig := GetCsvMgr().GetHeroConfig(heroid)
	if heroConfig == nil {
		return GetCsvMgr().GetText("STR_HELPFUN_NEAR")
	}

	if heroConfig.HeroArms == 1 {
		return GetCsvMgr().GetText("STR_HELPFUN_NEAR")
	} else if heroConfig.HeroArms == 2 {
		return GetCsvMgr().GetText("STR_HELPFUN_FAR")
	}

	return GetCsvMgr().GetText("STR_HELPFUN_NEAR")
}

//!获取武将武魂数量
func HF_GetHeroWuHun(player *Player, heroid int) int {
	if player == nil {
		LogError("player is nil.")
		return 0
	}
	itemId := 12000000 + heroid/10*100 + 1
	return player.GetObjectNum(itemId)
}

//! 加密消息
func HF_DecodeCenterMsg(msghead string, msgdata []byte) []byte {
	data := &MsgBase{
		Msghead: proto.String(msghead),
		Msgtime: proto.Int64(TimeServer().Unix()),
		Msgsign: proto.String("222"),
		Msgdata: []byte(msgdata),
	}

	msg, err := proto.Marshal(data)
	if err != nil {
		log.Println("HF_DecodeCenterMsg:", err)
		return []byte("")
	}

	return msg
}

// 修正玩家uid, king需要特殊处理
func HF_SetPUid(uid *int64, playerIndex int64) {
	if *uid > 0 && *uid < playerIndex {
		*uid = *uid + playerIndex
	}
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
	source := rand.New(rand.NewSource(TimeServer().UnixNano()))
	return source.Intn(check) + min
}

//! 判断是否新的一天
func HF_IsNewDate(checkDate string) bool {
	lct, err := time.ParseInLocation(DATEFORMAT, checkDate, time.Local)
	if err != nil {
		LogError("HF_IsNewDate:", err.Error(), checkDate)
		return false
	}

	checkToday := time.Date(lct.Year(), lct.Month(), lct.Day(), 5, 0, 0, 0, TimeServer().Location()).Unix()
	if lct.Hour() >= 5 {
		checkToday += DAY_SECS
	}

	if TimeServer().Unix() >= checkToday {
		return true
	} else {
		return false
	}
}

func HF_AbsInt(a int) int {
	if a > 0 {
		return a
	}
	return a * -1
}

// 将attLst的属性合并到attMap中去
func AddAttrHelper(attMap map[int]*Attribute, attLst []*Attribute) {
	if attMap == nil {
		LogError("attMap is nil.")
		return
	}

	for _, pAttr := range attLst {
		if pAttr.AttValue == 0 {
			continue
		}

		_, ok := attMap[pAttr.AttType]
		if !ok {
			attMap[pAttr.AttType] = &Attribute{pAttr.AttType, pAttr.AttValue}
		} else {
			if attMap[pAttr.AttType] != nil {
				attMap[pAttr.AttType].AttValue += pAttr.AttValue
			} else {
				LogError("nil ptr in AddAttrHelper")
			}
		}
	}
}

func AddAttrHelperForTimes(attMap map[int]*Attribute, attType []int, attValue []int64, times int64) {
	if attMap == nil {
		return
	}

	if len(attType) <= 0 || len(attType) != len(attValue) || times <= 0 {
		LogError("param is error.")
		return
	}

	for i := 0; i < len(attType); i++ {
		if attType[i] <= 0 || attValue[i] <= 0 {
			continue
		}

		_, ok := attMap[attType[i]]
		if !ok {
			attMap[attType[i]] = &Attribute{attType[i], attValue[i] * times}
		} else {
			if attMap[attType[i]] != nil {
				attMap[attType[i]].AttValue += attValue[i] * times
			} else {
				LogError("nil ptr in AddAttrHelperForTimes")
			}
		}
	}
}

func AddAttrHelperForGrowth(attMap map[int]*Attribute, attType []int, attValue []int64, growthType []int, growthValue []int64, quaType []int, quaValue []int64) {
	if attMap == nil {
		LogError("attMap is nil.")
		return
	}

	if len(attType) <= 0 || len(attType) != len(attValue) || len(attType) != len(growthType) || len(attType) != len(growthValue) || len(attType) != len(quaType) || len(attType) != len(quaValue) {
		LogError("param is error.")
		return
	}

	for i := 0; i < len(attType); i++ {
		if attType[i] <= 0 || attValue[i] <= 0 {
			continue
		}

		value := float32(attValue[i])
		if attType[i] == growthType[i] && attType[i] == quaType[i] {
			value = value * (float32(growthValue[i]) / float32(PER_BIT)) * (float32(quaValue[i]) / float32(PER_BIT))
		} else {
			continue
		}

		_, ok := attMap[attType[i]]
		if !ok {
			attMap[attType[i]] = &Attribute{attType[i], int64(value)}
		} else {
			if attMap[attType[i]] != nil {
				attMap[attType[i]].AttValue += int64(value)
			} else {
				LogError("nil ptr in AddAttrHelperForTimes")
			}
		}
	}
}

// 通过技能Id获得技能所增加的属性, 默认1级
func GetSkillAttr(skillIds []int) map[int]*Attribute {
	//skillLv := 1
	//slevel := float32(skillLv - 1)
	res := make(map[int]*Attribute)
	for _, pSkill := range skillIds {
		skillConfig, ok := GetCsvMgr().SkillConfigMap[pSkill]
		if !ok {
			continue
		}

		if skillConfig.SkillType != 1 {
			continue
		}
		if len(skillConfig.SkillValueType) != len(skillConfig.SkillCount) {
			continue
		}

		for i := 0; i < len(skillConfig.SkillValueType); i++ {
			if skillConfig.SkillCount[i] == 0 {
				continue
			}
			skillType := skillConfig.SkillValueType[i]
			_, ok := res[skillType]

			if !ok {
				res[skillType] = &Attribute{skillType, int64(skillConfig.SkillCount[i])}
			} else {
				if res[skillType] != nil {
					res[skillType].AttValue += int64(skillConfig.SkillCount[i])
				} else {
					LogError("nil ptr in GetSkillAttr")
				}
			}
		}
	}

	return res
}

// 将src的属性合并到dst中去
func AddAttrMapHelper(dst map[int]*Attribute, src map[int]*Attribute) {
	if dst == nil {
		LogError("dst is nil.")
		return
	}

	for _, pAttr := range src {
		if pAttr == nil {
			continue
		}

		if pAttr.AttValue == 0 {
			continue
		}

		_, ok := dst[pAttr.AttType]
		if !ok {
			dst[pAttr.AttType] = &Attribute{pAttr.AttType, pAttr.AttValue}
		} else {
			if dst[pAttr.AttType] != nil {
				dst[pAttr.AttType].AttValue += pAttr.AttValue
			} else {
				LogError("nil ptr in AddAttrMapHelper")
			}
		}
	}
}

func HF_MergePassitem(passItem []PassItem) []PassItem {
	rel := make([]PassItem, 0)
	if passItem == nil {
		return rel
	}
	tempRel := make(map[int]*Item, 0)
	AddItemMapHelper2(tempRel, passItem)

	for _, v := range tempRel {
		rel = append(rel, PassItem{v.ItemId, v.ItemNum})
	}
	return rel
}

// 更加简化的属性计算公式
func AddAttrDirect(dst map[int]*Attribute, attTypes []int, attValues []int64) {
	if dst == nil {
		LogError("dst is nil.")
		return
	}

	if len(attTypes) <= 0 || len(attValues) <= 0 {
		LogError("len(attTypes) <= 0 || len(attValues) <= 0")
		return
	}

	if len(attTypes) != len(attValues) {
		LogError("len(attTypes) != len(attValues)")
		return
	}

	for index, attType := range attTypes {
		attValue := attValues[index]
		if attValue == 0 {
			continue
		}

		_, ok := dst[attType]
		if !ok {
			//这里不好处理int64 暂时搁置 20190506 by zy
			dst[attType] = &Attribute{attType, int64(attValue)}
		} else {
			if dst[attType] != nil {
				dst[attType].AttValue += int64(attValue)
			} else {
				LogError("nil ptr in AddAttrDirect")
			}
		}
	}
}

// 道具增加
func AddItemMapHelper(dst map[int]*Item, itemIds []int, itemNums []int) {
	if dst == nil {
		return
	}

	if len(itemIds) != len(itemNums) {
		LogError("len(itemIds) != len(itemNums)")
		return
	}

	for index := range itemIds {
		itemId := itemIds[index]
		itemNum := itemNums[index]
		if itemId == 0 || itemNum == 0 {
			continue
		}
		_, ok := dst[itemId]
		if !ok {
			dst[itemId] = &Item{itemId, itemNum}
		} else {
			if dst[itemId] != nil {
				dst[itemId].ItemNum += itemNum
			} else {
				LogError("nil ptr in AddItemMapHelper")
			}
		}
	}
}

// 道具增加
func AddItemMapHelper2(dst map[int]*Item, items []PassItem) {
	if dst == nil {
		return
	}

	for index := range items {
		itemId := items[index].ItemID
		itemNum := items[index].Num
		if itemId == 0 || itemNum == 0 {
			continue
		}
		_, ok := dst[itemId]
		if !ok {
			dst[itemId] = &Item{itemId, itemNum}
		} else {
			if dst[itemId] != nil {
				dst[itemId].ItemNum += itemNum
			} else {
				LogError("nil ptr in AddItemMapHelper2")
			}
		}
	}
}

func AddItemMapHelper3(dst map[int]*Item, itemId int, itemNum int) {
	if dst == nil {
		return
	}

	if itemId == 0 || itemNum == 0 {
		return
	}
	_, ok := dst[itemId]
	if !ok {
		dst[itemId] = &Item{itemId, itemNum}
	} else {
		if dst[itemId] != nil {
			dst[itemId].ItemNum += itemNum
		} else {
			LogError("nil ptr in AddItemMapHelper2")
		}
	}
}

// 道具增加
func AddItemMapHelper4(dst map[int]*Item, items map[int]*Item) {
	if dst == nil {
		return
	}

	for index, _ := range items {
		itemId := items[index].ItemId
		itemNum := items[index].ItemNum
		if itemId == 0 || itemNum == 0 {
			continue
		}
		_, ok := dst[itemId]
		if !ok {
			dst[itemId] = &Item{itemId, itemNum}
		} else {
			if dst[itemId] != nil {
				dst[itemId].ItemNum += itemNum
			} else {
				LogError("nil ptr in AddItemMapHelper3")
			}
		}
	}
}

// 道具增加
func AddItemMap(dst map[int]*Item, srcMap map[int]*Item) {
	if dst == nil || srcMap == nil {
		return
	}

	for _, item := range srcMap {
		itemId := item.ItemId
		itemNum := item.ItemNum
		if itemId == 0 || itemNum == 0 {
			continue
		}
		_, ok := dst[itemId]
		if !ok {
			dst[itemId] = &Item{itemId, itemNum}
		} else {
			if dst[itemId] != nil {
				dst[itemId].ItemNum += itemNum
			} else {
				LogError("nil ptr in additemMap")
			}
		}
	}
}

// 增加战斗力属性
func AddFightAtt(dst map[int]*Attribute, attType int, attValue int64) {
	if attValue > 0 {
		_, ok := dst[attType]
		if !ok {
			dst[attType] = &Attribute{attType, attValue}
		} else {
			if dst[attType] != nil {
				dst[attType].AttValue += attValue
			} else {
				LogError("nil ptr in AddFightAtt")
			}
		}
	}
}

// 打印堆栈信息
func DumpStacks() {
	buf := make([]byte, 16384)
	buf = buf[:runtime.Stack(buf, true)]
	fmt.Printf("=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===", buf)
}

func HF_Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func HF_ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(HF_Round(num*output)) / output
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

func HF_SetBit(n int64, pos uint) int64 {
	n |= (1 << pos)
	return n
}

func HF_ClearBit(n int64, pos uint) int64 {
	var mask int64
	mask = ^(1 << pos)
	n &= int64(mask)
	return n
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

//计算下次月刷新时间
func HF_GetNextMonthStart() int64 {
	now := TimeServer()
	now = now.AddDate(0, 0, -now.Day()+1)
	return now.AddDate(0, 1, 0).Unix()
}

//计算下次周刷新时间
func HF_GetNextWeekStart() int64 {

	now := TimeServer()

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return weekStart.Unix() + 7*DAY_SECS
}

//计算下次日刷新时间
func HF_GetNextDayStart() int64 {
	now := TimeServer()
	if now.Hour() < 5 {
		return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix()
	} else {
		return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location()).Unix() + DAY_SECS
	}
}

func HF_GetNextTimeToLog() int64 {
	now := TimeServer()
	nextTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 0, 0, now.Location()).Unix()
	if nextTime < now.Unix() {
		nextTime += DAY_SECS
	}
	return nextTime
}

//计算地牢结束时间
func HF_GetNewPitEnd() int64 {
	for _, config := range GetCsvMgr().TimeResetConfig {
		if config.System == NEWPIT_TIME_SYSTEM {
			switch config.TimeType {
			case TIME_RESET_TYPE_TIME:
				stage := config.Continue + config.Cd
				calCount := (TimeServer().Unix() - config.Time[0]) / stage
				calTime := (TimeServer().Unix() - config.Time[0]) % stage
				continueTime := config.Continue

				isOpen, index := GetActivityMgr().JudgeOpenIndex(ACT_NEWPIT_HALF_MIN, ACT_NEWPIT_HALF_MAX)
				if isOpen {
					activity := GetActivityMgr().GetActivity(index)
					if activity != nil {
						startday := HF_Atoi(activity.info.Start)
						if startday >= 0 {
							if activity.status.Status > ACTIVITY_STATUS_CLOSED {
								continueTime = continueTime / 2
							}
						}
					}
				}

				//处于冷却期则计算期数加1
				if calTime > continueTime {
					calCount += 1
				}
				startTime := config.Time[0] + calCount*stage
				return startTime + continueTime
			}
		}
	}
	return 0
}

func HF_GetNewPitMapId() int {
	now := TimeServer().Unix()
	openTime := GetServer().GetOpenServer()
	times := int(now-openTime) / (DAY_SECS * 2)
	return times%2 + 1
}

//角色创建时间为基准，参数为天数偏移
func HF_CalPlayerCreateTime(calTime int64, disDay int) int64 {
	cal := time.Unix(calTime, 0)
	rel := int64(0)
	if cal.Hour() < 5 {
		rel = time.Date(cal.Year(), cal.Month(), cal.Day(), 5, 0, 0, 0, cal.Location()).Unix() + DAY_SECS*int64(disDay-1)
	} else {
		rel = time.Date(cal.Year(), cal.Month(), cal.Day(), 5, 0, 0, 0, cal.Location()).Unix() + DAY_SECS*int64(disDay)
	}

	return rel
}

// 增加多个道具
func HF_RewardCsvToPassItem(ids []int, nums []int) []PassItem {
	var res []PassItem
	if len(ids) != len(nums) {
		return res
	}

	for index := range ids {
		id := ids[index]
		if id == 0 {
			continue
		}
		num := nums[index]
		if num == 0 {
			continue
		}
		res = append(res, PassItem{ItemID: id, Num: num})
	}

	return res
}

//计算时间 动态支持3种模式
func HF_CalTimeForConfig(stime string, rtime string) int64 {
	relTime := int64(0)
	//! 开放时间为整数时，为开服天数，一次性有效
	startday := HF_Atoi(stime)
	//! 如果是按照开服时间算
	if startday > 0 {
		//! 计算开服时间
		relTime = GetServer().GetOpenServer() + int64((startday-1)*DAY_SECS)
	} else if startday < 0 {
		rTime, _ := time.ParseInLocation(DATEFORMAT, rtime, time.Local)
		relTime = rTime.Unix() + int64((-startday-1)*DAY_SECS)
	} else {
		//! 否则按照实际时间算
		t, err := time.ParseInLocation(DATEFORMAT, stime, time.Local)
		if err != nil {
			t, err = NewTimeUtil(TimeServer()).Parse(stime)
			if err != nil {
				LogError("时间填写错误:", err.Error())
				return relTime
			}
		}
		relTime = t.Unix()
	}
	return relTime
}

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func HF_GetGoroutineID() uint64 {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
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

// 获得装备颜色
func HF_GetColorByQuality(quality int) string {
	strColor := ""
	switch quality {
	case 1:
		strColor = "228#228#228#"
	case 2:
		strColor = "66#216#51#"
	case 3:
		strColor = "50#200#255#"
	case 4:
		strColor = "255#50#255#"
	case 5:
		strColor = "255#185#50#"
	case 6:
		strColor = "255#50#50#"
	case 7:
		strColor = "255#240#0#"
	case 8:
		strColor = "228#228#228#"
	}
	return strColor
}

