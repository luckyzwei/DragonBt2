package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	RedPacTime = 72000
	RedPacLast = 3600
	ActTypeRed = 9010
)

func (self *RedPacMgr) GlobalRedPacSql() string {
	return `CREATE TABLE IF NOT EXISTS san_redpac (
		  keyid bigint(20) NOT NULL AUTO_INCREMENT,
		  redid bigint(20) NOT NULL DEFAULT '0' COMMENT '红包id',
		  unionid int(11) NOT NULL DEFAULT '0' COMMENT '军团Id',
		  camp int(11) NOT NULL DEFAULT '0' COMMENT '国家Id',
		  redtype int(11) NOT NULL DEFAULT '1' COMMENT '红包类型',
		  endtime bigint(20) NOT NULL DEFAULT '0' COMMENT '红包结束时间',
		  uid bigint(20) NOT NULL DEFAULT '0' COMMENT '发起人uid',
		  uname varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL COMMENT '用户名称',
		  people int(11) NOT NULL DEFAULT '0' COMMENT '人数',
		  totalnum int(11) NOT NULL DEFAULT '0' COMMENT '红包道具总数',
		  itemid int(11) NOT NULL DEFAULT '0' COMMENT '道具Id',
		  redstatus text NOT NULL COMMENT '红包状态',
		  msg text NOT NULL COMMENT '留言',
          iconid int(11) NOT NULL DEFAULT '0' COMMENT '头像Id',
          rednum text NOT NULL COMMENT '红包金额',
 		  starttime bigint(20) NOT NULL DEFAULT '0' COMMENT '红包开始时间',
		  PRIMARY KEY (keyid,redid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

func (self *RedPacMgr) UserRedPacSql() string {
	return `CREATE TABLE IF NOT EXISTS san_userredpac (
		    uid bigint(20) NOT NULL DEFAULT '0' COMMENT '玩家Id',
  			info text NOT NULL COMMENT '红包信息',
  			PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
}

type RedStatus struct {
	Uid       int64  `json:"uid"`       //! 抢到人的uid
	Uname     string `json:"uname"`     //! 抢到的人名字
	Num       int    `json:"num"`       //! 抢到的金额
	TimeStamp int64  `json:"timestamp"` //! 抢到的时间戳
	Status    int    `json:"status"`    //! 0 未抢 1 已领取 2 被抢完 3 已答谢[只有领取的玩家才能答谢] 存放到全局
	ItemId    int    `json:"itemid"`    //! 道具Id
}

type San_GloRedPac struct {
	KeyId     int64  //! 自增长id
	RedId     int64  //! 红包Id
	UnionId   int    //! 军团Id: 红包信息=0 说明是军团红包
	Camp      int    //! 国家Id: 红包信息=0, 说明是国家红包
	RedType   int    //! 红包类型: 1国家 2军团 3个人红包池
	EndTime   int64  //! 结束时间
	Uid       int64  //! 发起人uid
	UName     string //! 发起人姓名
	People    int    //! 人数
	TotalNum  int    //! 金额
	ItemId    int    //! 道具Id
	RedStatus string //! 红包状态
	Msg       string //! 留言
	Iconid    int    //! 头像
	RedNum    string //! 红包金额
	StartTime int64  //! 开始时间

	Status map[int64]*RedStatus //! 红包状态
	Num    []int                //! 红包金额

	DataUpdate
}

func (self *San_GloRedPac) Decode() {
	json.Unmarshal([]byte(self.RedStatus), &self.Status)
	json.Unmarshal([]byte(self.RedNum), &self.Num)
}

func (self *San_GloRedPac) Encode() {
	self.RedStatus = HF_JtoA(self.Status)
	self.RedNum = HF_JtoA(self.Num)
}

/**
 * @param count 红包数
 * @param money 总金额
 * @return []int 返回随机金额
 */
func RandRedPac(money int, count int) []int {
	result := make([]int, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, 0)
	}

	source := rand.New(rand.NewSource(TimeServer().UnixNano()))
	for i := 0; i < count-1; i++ {
		n1 := money / 100
		n2 := money / (count - i)
		n3 := n2 * 2
		check := n3 - n1
		if check <= 0 {
			check = 1
		}

		if n1 <= 0 {
			n1 = 1
		}

		result[i] = source.Intn(check) + n1
		money -= result[i]
	}
	result[count-1] = money
	return result
}

// 红包全局管理器
type RedPacMgr struct {
	RedPacMap       map[int64]*San_GloRedPac //! 红包管理器
	RedPacLock      *sync.RWMutex            //! 红包锁
	RedPacCamp      map[int][]int64          //! 国家红包, key:camp, value: redpac keyId
	RedPacUnion     map[int][]int64          //! 军团红包, key:unionId, value: redpac keyId
	MaxKeyId        int64                    //! 服务器最大Id, 初始值serverId * 1000000
	RedPacCampLock  *sync.RWMutex            //! 国家红包锁
	RedPacUnionLock *sync.RWMutex            //! 国家红包锁
}

var redPacMgr *RedPacMgr = nil

// 需要做camp, union映射的keyId
func GetRedPacMgr() *RedPacMgr {
	if redPacMgr == nil {
		redPacMgr = new(RedPacMgr)
		redPacMgr.RedPacMap = make(map[int64]*San_GloRedPac)
		redPacMgr.RedPacLock = new(sync.RWMutex)
		redPacMgr.RedPacCamp = make(map[int][]int64)
		for camp := 1; camp <= 3; camp++ {
			redPacMgr.RedPacCamp[camp] = make([]int64, 0)
		}
		redPacMgr.RedPacUnion = make(map[int][]int64)
		redPacMgr.RedPacCampLock = new(sync.RWMutex)
		redPacMgr.RedPacUnionLock = new(sync.RWMutex)
		redPacMgr.MaxKeyId = int64(10000000 * GetServer().Con.ServerId)
	}

	return redPacMgr
}

type RedPacParam struct {
	Uid      int64  //! 发起人uid
	UnionId  int    //! 军团Id: 红包信息=0 说明是军团红包
	Camp     int    //! 国家Id: 红包信息=0, 说明是国家红包
	RedType  int    //! 红包类型: 1国家 2军团 3个人红包池
	UName    string //! 发起人姓名
	People   int    //! 人数
	TotalNum int    //! 金额
	ItemId   int    //! 道具Id
	Msg      string //! 留言
	IconId   int    //! 头像Id
}

// 创建红包
func (self *RedPacMgr) CreateRedPac(param *RedPacParam) *San_GloRedPac {
	atomic.AddInt64(&self.MaxKeyId, 1)
	redPac := &San_GloRedPac{
		RedId:     atomic.LoadInt64(&self.MaxKeyId),
		UnionId:   param.UnionId,
		Camp:      param.Camp,
		RedType:   param.RedType,
		EndTime:   TimeServer().Unix() + RedPacTime,
		Uid:       param.Uid,
		UName:     param.UName,
		People:    param.People,
		TotalNum:  param.TotalNum,
		ItemId:    param.ItemId,
		Msg:       param.Msg,
		Iconid:    param.IconId,
		StartTime: TimeServer().Unix(),
		Status:    make(map[int64]*RedStatus),
		Num:       RandRedPac(param.TotalNum, param.People),
	}

	redPac.Encode()
	// 插入数据库
	lastId := InsertTable(self.GetTable(), redPac, 0, false)
	if lastId != 0 {
		redPac.KeyId = lastId
	}
	redPac.Init(self.GetTable(), redPac, false)

	self.RedPacLock.Lock()
	self.RedPacMap[redPac.RedId] = redPac
	self.RedPacLock.Unlock()

	return redPac
}

// 全国广播红包
func (self *RedPacMgr) Send2Camp(camp int, keyId int64) {
	var msg S2C_RedPac
	msg.Cid = "redpac"
	msg.KeyId = keyId
	GetPlayerMgr().BroadCastMsgToCamp(camp, msg.Cid, HF_JtoB(&msg))
}

// 军团广播红包
func (self *RedPacMgr) Send2Union(union *San_Union, keyId int64) {
	var msg S2C_RedPac
	msg.Cid = "redpac"
	msg.KeyId = keyId
	union.BroadCastMsg(msg.Cid, HF_JtoB(&msg))
}

// 查询红包状态
func (self *RedPacMgr) GetRedPac(uid int64, keyId int64) (*San_GloRedPac, bool, bool) {
	self.RedPacLock.RLock()
	defer self.RedPacLock.RUnlock()

	redPac, ok := self.RedPacMap[keyId]
	if !ok {
		return nil, false, false
	}

	status := redPac.Status
	_, isTaken := status[uid]
	allTaken := (len(redPac.Status) == redPac.People)

	return redPac, isTaken, allTaken
}

// 分配红包
func (self *RedPacMgr) AllocateRedPac(keyId int64, uid int64, name string) *RedStatus {
	self.RedPacLock.Lock()
	defer self.RedPacLock.Unlock()

	redPac, ok := self.RedPacMap[keyId]
	if !ok {
		return nil
	}

	allocateNum := len(redPac.Status)
	if allocateNum >= redPac.People {
		return nil
	}

	if allocateNum < 0 || allocateNum >= len(redPac.Num) {
		return nil
	}

	num := redPac.Num[allocateNum]
	now := TimeServer().Unix()
	pRedStatus := &RedStatus{
		Uid:       uid,
		Uname:     name,
		Num:       num,
		TimeStamp: now,
		Status:    RedPacTaken,
		ItemId:    redPac.ItemId,
	}

	redPac.Status[uid] = pRedStatus
	if allocateNum+1 == redPac.People {
		redPac.EndTime = now + RedPacLast
	}

	return pRedStatus
}

// 增加国家红包
func (self *RedPacMgr) addCampRed(camp int, keyId int64) {
	self.RedPacCampLock.Lock()
	defer self.RedPacCampLock.Unlock()
	self.RedPacCamp[camp] = append(self.RedPacCamp[camp], keyId)
}

// 删除国家红包
func (self *RedPacMgr) removeCampRed(camp int, keyIds map[int64]bool) {
	self.RedPacCampLock.Lock()
	defer self.RedPacCampLock.Unlock()

	campSlice := self.RedPacCamp[camp]
	var leftKey []int64
	mapKey := make(map[int64]struct{})
	for keyId := range keyIds {
		mapKey[keyId] = struct{}{}
	}

	for index := range campSlice {
		keyId := campSlice[index]
		if _, ok := mapKey[keyId]; !ok {
			leftKey = append(leftKey, keyId)
		}
	}

	if len(leftKey) > 0 {
		self.RedPacCamp[camp] = leftKey
	} else {
		self.RedPacCamp[camp] = make([]int64, 0)
	}
	LogDebug("unionSlice remove, now leftKey:", leftKey)

}

// 增加军团红包
func (self *RedPacMgr) addUnionRed(unionId int, keyId int64) {
	self.RedPacUnionLock.Lock()
	defer self.RedPacUnionLock.Unlock()
	self.RedPacUnion[unionId] = append(self.RedPacUnion[unionId], keyId)
}

// 删除军团红包, 批量删除多个key
func (self *RedPacMgr) removeUnionRed(unionId int, keyIds map[int64]bool) {
	self.RedPacUnionLock.Lock()
	defer self.RedPacUnionLock.Unlock()

	unionSlice := self.RedPacUnion[unionId]
	var leftKey []int64
	mapKey := make(map[int64]struct{})
	for keyId := range keyIds {
		mapKey[keyId] = struct{}{}
	}

	for index := range unionSlice {
		keyId := unionSlice[index]
		if _, ok := mapKey[keyId]; !ok {
			leftKey = append(leftKey, keyId)
		}
	}

	if len(leftKey) > 0 {
		self.RedPacUnion[unionId] = leftKey
	} else {
		self.RedPacUnion[unionId] = make([]int64, 0)
	}
	LogDebug("unionSlice remove, now leftKey:", leftKey)
}

// 删除红包
func (self *RedPacMgr) removeRedPac(keyId int64) {
	self.RedPacLock.Lock()
	defer self.RedPacLock.Unlock()
	delete(self.RedPacMap, keyId)
}

// 增加时检查红包
func (self *RedPacMgr) checkUnionRedPac(unionId int) {
	unionKey := make(map[int64]bool)
	now := TimeServer().Unix()
	self.RedPacLock.RLock()
	for key, value := range self.RedPacMap {
		if value.EndTime > now {
			continue
		}

		if value.RedType != RedUnionType && value.RedType != RedPoolType {
			continue
		}

		if value.RedType == RedActType && value.UnionId <= 0 {
			continue
		}

		unionKey[key] = true
	}
	self.RedPacLock.RUnlock()

	self.RedPacUnionLock.RLock()
	redSlice := self.RedPacUnion[unionId]
	curUnionNum := len(redSlice)
	if curUnionNum >= MaxUnionRedNum {
		keyId := redSlice[curUnionNum-1]
		unionKey[keyId] = true
	}
	self.RedPacUnionLock.RUnlock()

	for key := range unionKey {
		self.removeRedPac(key)
	}

	if len(unionKey) > 0 {
		self.removeUnionRed(unionId, unionKey)
	}
}

// 获取需要删除的国家红包Id
func (self *RedPacMgr) checkCampRedPac(camp int) {
	campKey := make(map[int64]bool)
	now := TimeServer().Unix()
	self.RedPacLock.RLock()
	for key, value := range self.RedPacMap {
		if value.EndTime > now {
			continue
		}

		if value.RedType != RedCampType {
			continue
		}

		if value.RedType == RedActType && value.Camp <= 0 {
			continue
		}

		campKey[key] = true
	}
	self.RedPacLock.RUnlock()

	self.RedPacCampLock.RLock()
	curCampNum := len(self.RedPacCamp[camp])
	if curCampNum >= MaxCampRedNum {
		campKey[self.RedPacCamp[camp][curCampNum-1]] = true
	}
	self.RedPacCampLock.RUnlock()

	for key := range campKey {
		self.removeRedPac(key)
	}

	if len(campKey) > 0 {
		self.removeCampRed(camp, campKey)
	}
}

func (self *RedPacMgr) newRedPacShow(redPac *San_GloRedPac) *RedPacShow {
	duration := redPac.EndTime - TimeServer().Unix()
	if duration < 0 {
		duration = 0
	}
	return &RedPacShow{
		KeyId:     redPac.RedId,
		Uname:     redPac.UName,
		ItemId:    redPac.ItemId,
		ItemNum:   redPac.TotalNum,
		TimeStamp: redPac.StartTime,
		Status:    0,
		Duration:  duration,
		Uid:       redPac.Uid,
	}
}

// 查询国家红包信息
func (self *RedPacMgr) getCampRedPacShow(camp int) map[int64]*RedPacShow {
	var campKey []int64

	self.RedPacCampLock.RLock()
	campSlice := self.RedPacCamp[camp]
	campKey = append(campKey, campSlice...)
	self.RedPacCampLock.RUnlock()

	res := make(map[int64]*RedPacShow)
	self.RedPacLock.RLock()
	for index := range campKey {
		keyId := campKey[index]
		redPac, ok := self.RedPacMap[keyId]
		if !ok {
			LogError("camp redpac:", keyId, "is nil.")
			continue
		}
		if redPac.EndTime <= TimeServer().Unix() {
			continue
		}
		res[redPac.RedId] = self.newRedPacShow(redPac)
	}
	self.RedPacLock.RUnlock()

	return res
}

// 查询军团红包信息
func (self *RedPacMgr) getUnionRedPacShow(unionId int) map[int64]*RedPacShow {
	var unionkey []int64
	self.RedPacUnionLock.RLock()
	unionSlice := self.RedPacUnion[unionId]
	unionkey = append(unionkey, unionSlice...)
	self.RedPacUnionLock.RUnlock()

	res := make(map[int64]*RedPacShow)
	self.RedPacLock.RLock()
	for index := range unionkey {
		keyId := unionkey[index]
		redPac, ok := self.RedPacMap[keyId]
		if !ok {
			LogError("union redpac:", keyId, "is nil, unionkey = ", unionkey)
			continue
		}
		if redPac.EndTime <= TimeServer().Unix() {
			continue
		}
		res[redPac.RedId] = self.newRedPacShow(redPac)
	}
	self.RedPacLock.RUnlock()

	return res
}

// 全局保存
func (self *RedPacMgr) Save() {
	self.RedPacLock.RLock()
	defer self.RedPacLock.RUnlock()

	for _, value := range self.RedPacMap {
		value.Encode()
		value.Update(false)
	}
}

// 从数据库拉取数据,删除过期的红包
func (self *RedPacMgr) GetData() {
	var dbData San_GloRedPac
	now := TimeServer().Unix()
	sql := fmt.Sprintf("select * from `%s` where `endtime` >= %d", self.GetTable(), now)
	res := GetServer().DBUser.GetAllData(sql, &dbData)
	for i := 0; i < len(res); i++ {
		data := res[i].(*San_GloRedPac)
		data.Init(self.GetTable(), data, false)
		data.Decode()
		if data.EndTime-now > 72000 {
			data.EndTime = now + 7200
		}
		self.RedPacLock.Lock()
		self.RedPacMap[data.RedId] = data
		self.RedPacLock.Unlock()
		if data.RedType == RedCampType {
			self.addCampRed(data.Camp, data.RedId)
		} else if data.RedType == RedUnionType || data.RedType == RedPoolType {
			self.addUnionRed(data.UnionId, data.RedId)
		} else if data.RedType == RedActType && data.Camp > 0 {
			self.addCampRed(data.Camp, data.RedId)
		} else if data.RedType == RedActType && data.UnionId > 0 {
			self.addUnionRed(data.UnionId, data.RedId)
		}

		if redPacMgr.MaxKeyId < data.RedId {
			redPacMgr.MaxKeyId = data.RedId
		}
	}

	//self.checkRedPac()
}

func (self *RedPacMgr) GetTable() string {
	return "san_redpac"
}

// 国家系统消息
func (self *RedPacMgr) sendSystemChat(uname string, camp int, itemId int) {
	itemConfig := GetCsvMgr().GetItemConfig(itemId)
	if itemConfig != nil {
		str := fmt.Sprintf(GetCsvMgr().GetText("STR_RED_PAC_CHAT"), HF_GetColorByCamp(camp), uname, GetCsvMgr().GetItemName(itemId))
		GetServer().sendCampChat(camp, str)
	}
}

// 军团系统消息
func (self *RedPacMgr) sendUnionChat(sanUnion *San_Union, camp int, uname string, itemId int) {
	itemConfig := GetCsvMgr().GetItemConfig(itemId)
	if itemConfig != nil {
		str := fmt.Sprintf(GetCsvMgr().GetText("STR_RED_PAC_CHAT"), HF_GetColorByCamp(camp), uname, GetCsvMgr().GetItemName(itemId))
		GetServer().sendUnionChat(sanUnion, str)
	}
}

// 开启定时器自动发放红包12,18
// 是否已经发放存放在redis里面,每天进行更新
// key: actRedPac, 服务器年月日小时, 服务器检查启动12:00:01~12:00:59, 18:00:01~18:00:59
// HMGET, HMSET
// actRedPac1: time; actRedPac2: time
func (self *RedPacMgr) checkRedPac() {
	key := "actRedpac"
	filed1 := "act12"
	filed2 := "act18"

	v, err := GetRedisMgr().HGetAll(key)
	if err != nil {
		LogDebug("checkRedPac:", err.Error())
		return
	}

	value := map[string]interface{}{
		filed1: "1234",
		filed2: "5678",
	}

	if len(v) <= 0 {
		err := GetRedisMgr().HMSet(key, value)
		if err != nil {
			LogDebug(err.Error())
		}
	} else {
		LogDebug(fmt.Sprintf("checkRedPac:+%v", v))
	}
}

// 增加红包活动检查, 一分钟检查一下,差不多12,18点的时候发送红包
func (self *RedPacMgr) OnTimer() {
	ticker := time.NewTicker(time.Minute * 1)
	for {
		<-ticker.C
		self.checkRedHour()
	}

	ticker.Stop()
}

// 系统红包keyId = guid
func (self *RedPacMgr) checkRedHour() {
	now := TimeServer()
	// 获取活动
	if now.Hour() == 12 && now.Minute() == 0 {
		self.checkRedAct(12)
	}

	if now.Hour() == 18 && now.Minute() == 0 {
		self.checkRedAct(18)
	}
}

func (self *RedPacMgr) isHourOk(hour int) bool {
	return hour == 12 || hour == 18
}

// 分三种情况发送红包
func (self *RedPacMgr) checkRedAct(hour int) {
	activity := GetActivityMgr().GetActivity(ActTypeRed)
	if activity == nil {
		return
	}

	startTime := activity.getActTime()
	endTime := startTime + int64(activity.info.Continued)
	now := TimeServer().Unix()
	if now < startTime || now > endTime {
		return
	}

	taskType := activity.getTaskType()
	if taskType < 1 || taskType > 2 {
		return
	}

	pacSlice := activity.getN4()
	if len(pacSlice) <= 0 {
		return
	}

	timeInfo := pacSlice[3]
	if timeInfo == 1 {
		if hour != 12 {
			return
		}
	} else if timeInfo == 2 {
		if hour != 18 {
			return
		}
	} else if timeInfo == 3 {
		if !self.isHourOk(hour) {
			return
		}
	}

	txt := activity.getTxt()
	redPacParam := &RedPacParam{
		Uid:      0,
		UnionId:  0,
		Camp:     0,
		RedType:  RedActType,
		UName:    GetCsvMgr().GetText("STR_SYS_TEXT"),
		People:   pacSlice[0],
		TotalNum: pacSlice[2],
		ItemId:   pacSlice[1],
		Msg:      txt,
		IconId:   1003,
	}

	if taskType == 1 { // 国家红包
		for camp := 1; camp <= 3; camp++ { //
			redPacParam.Camp = camp
			redPac := GetRedPacMgr().CreateRedPac(redPacParam)
			GetRedPacMgr().checkCampRedPac(camp)
			GetRedPacMgr().addCampRed(camp, redPac.RedId)
			GetRedPacMgr().Send2Camp(camp, redPac.RedId)
			//GetRedPacMgr().sendSystemChat("系统", camp, pacSlice[1])
		}
	} else if taskType == 2 { // 军团红包
		//GetUnionMgr().sendRedPac(redPacParam)
	}
}

// 查询红包状态
func (self *RedPacMgr) CheckRedPacOk(keyId []int64) map[int64]bool {
	self.RedPacLock.RLock()
	defer self.RedPacLock.RUnlock()
	removeIds := make(map[int64]bool)
	for index := range keyId {
		value := keyId[index]
		_, ok := self.RedPacMap[value]
		if !ok {
			removeIds[value] = true
		}
	}
	return removeIds
}
