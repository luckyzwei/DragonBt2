package match

import (
	"encoding/json"
	"fmt"
	"master/db"
	"master/utils"
	"sort"
	"sync"
)

type JS_ConsumerTopServer struct {
	SvrId   int    `json:"svrid"`
	SvrName string `json:"svrname"`
	Rank    int    `json:"rank"`
	Point   int    `json:"point"`
	Kill    int    `json:"kill"`
	Step    int    `json:"step"`
}

type JS_ConsumerTopUserDB struct {
	Uid   int64  `json:"uid"`
	KeyId int    `json:"keyid"`
	SvrId int    `json:"svrid"`
	Info  string `json:"info"`

	info *JS_ConsumerTopUser
	db.DataUpdate
}

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
	KillAll  int    `json:"killall"`
}

type ConsumerTopUserArr []*JS_ConsumerTopUser

func (s ConsumerTopUserArr) Len() int      { return len(s) }
func (s ConsumerTopUserArr) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ConsumerTopUserArr) Less(i, j int) bool {
	return s[i].Point > s[j].Point
}

type ConsumerTopServerArr []*JS_ConsumerTopServer

func (s ConsumerTopServerArr) Len() int      { return len(s) }
func (s ConsumerTopServerArr) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ConsumerTopServerArr) Less(i, j int) bool {
	return s[i].Point > s[j].Point
}

type ConsumerTopInfo struct {
	KeyId       int    `json:"keyid"`
	ConsumerTop string `json:"generalusertop"`

	Mu                       *sync.RWMutex
	consumerTop              map[int64]*JS_ConsumerTopUser
	consumerTopNodeArr       ConsumerTopUserArr
	db_list                  map[int64]*JS_ConsumerTopUserDB //数据存储
	consumerTopServer        map[int]*JS_ConsumerTopServer
	consumerTopServerNodeArr ConsumerTopServerArr
}

type ConsumerTopMgr struct {
	Locker          *sync.RWMutex
	ConsumerTopInfo map[int]*ConsumerTopInfo
}

var consumerToMgr *ConsumerTopMgr = nil

func GetConsumerTopMgr() *ConsumerTopMgr {
	if consumerToMgr == nil {
		consumerToMgr = new(ConsumerTopMgr)
		consumerToMgr.ConsumerTopInfo = make(map[int]*ConsumerTopInfo)
		consumerToMgr.Locker = new(sync.RWMutex)
	}
	return consumerToMgr
}

func (self *JS_ConsumerTopUserDB) Encode() {
	self.Info = utils.HF_JtoA(self.info)
}

func (self *JS_ConsumerTopUserDB) Decode() {
	json.Unmarshal([]byte(self.Info), &self.info)
}

// 存储数据库
func (self *ConsumerTopMgr) OnSave() {
	for _, v := range self.ConsumerTopInfo {
		v.Save()
	}
}

func (self *ConsumerTopInfo) Save() {
	self.Mu.Lock()
	defer self.Mu.Unlock()
	for _, v := range self.db_list {
		v.Encode()
		v.Update(true, false)
	}
}

func (self *ConsumerTopMgr) NewConsumerTopInfo(KeyId int) *ConsumerTopInfo {
	data := new(ConsumerTopInfo)
	data.KeyId = KeyId
	data.Mu = new(sync.RWMutex)
	data.consumerTop = make(map[int64]*JS_ConsumerTopUser)
	data.consumerTopNodeArr = make([]*JS_ConsumerTopUser, 0)
	data.consumerTopServer = make(map[int]*JS_ConsumerTopServer, 0)
	data.consumerTopServerNodeArr = make([]*JS_ConsumerTopServer, 0)
	data.db_list = make(map[int64]*JS_ConsumerTopUserDB, 0)
	return data
}

func (self *ConsumerTopMgr) NewConsumerTopServer(serverId int, name string, keyId int) *JS_ConsumerTopServer {
	data := new(JS_ConsumerTopServer)
	data.SvrId = serverId
	data.SvrName = name
	data.Step = keyId
	return data
}

func (self *ConsumerTopMgr) GetAllData() {
	self.LoadConsumerTop()
}

func (self *ConsumerTopMgr) LoadConsumerTop() {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	queryStr := fmt.Sprintf("select uid,keyid,svrid,info from `tbl_consumertop`;")
	var msg JS_ConsumerTopUserDB
	res := db.GetDBMgr().DBUser.GetAllData(queryStr, &msg)

	for i := 0; i < len(res); i++ {
		data := res[i].(*JS_ConsumerTopUserDB)
		if data.KeyId <= 0 {
			continue
		}

		_, ok := self.ConsumerTopInfo[data.KeyId]
		if !ok {
			self.ConsumerTopInfo[data.KeyId] = self.NewConsumerTopInfo(data.KeyId)
		}

		if self.ConsumerTopInfo[data.KeyId] == nil {
			continue
		}
		data.Decode()
		if data.info == nil {
			continue
		}
		self.ConsumerTopInfo[data.KeyId].consumerTop[data.Uid] = data.info
		self.ConsumerTopInfo[data.KeyId].db_list[data.Uid] = data

		data.Init("tbl_consumertop", data, true)

		_, okServer := self.ConsumerTopInfo[data.KeyId].consumerTopServer[data.info.SvrId]
		if !okServer {
			self.ConsumerTopInfo[data.KeyId].consumerTopServer[data.info.SvrId] = self.NewConsumerTopServer(data.info.SvrId, data.info.SvrName, data.info.Step)
		}
		self.ConsumerTopInfo[data.KeyId].consumerTopServer[data.info.SvrId].Point += data.info.Point
		self.ConsumerTopInfo[data.KeyId].consumerTopServer[data.info.SvrId].Kill += data.info.KillAll
	}

	for _, v := range self.ConsumerTopInfo {
		v.MakeArr()
	}
}

func (self *ConsumerTopInfo) MakeArr() {
	self.Mu.Lock()
	defer self.Mu.Unlock()

	self.consumerTopNodeArr = make([]*JS_ConsumerTopUser, 0)
	for _, v := range self.consumerTop {
		self.consumerTopNodeArr = append(self.consumerTopNodeArr, v)
	}
	sort.Sort(self.consumerTopNodeArr)

	for i := 0; i < len(self.consumerTopNodeArr); i++ {
		self.consumerTopNodeArr[i].Rank = i + 1
	}

	self.consumerTopServerNodeArr = make([]*JS_ConsumerTopServer, 0)
	for _, v := range self.consumerTopServer {
		self.consumerTopServerNodeArr = append(self.consumerTopServerNodeArr, v)
	}
	sort.Sort(self.consumerTopServerNodeArr)

	for i := 0; i < len(self.consumerTopServerNodeArr); i++ {
		self.consumerTopServerNodeArr[i].Rank = i + 1
	}
}

func (self *ConsumerTopMgr) UpdatePoint(req *RPC_ConsumerTopActionReq, res *RPC_ConsumerTopActionRes) {
	if req.SelfInfo == nil {
		res.RetCode = RETCODE_DATA_ERROR
		return
	}
	self.Locker.Lock()
	info, ok := self.ConsumerTopInfo[req.SelfInfo.Step]
	if !ok {
		info = self.NewConsumerTopInfo(req.SelfInfo.Step)
		self.ConsumerTopInfo[info.KeyId] = info
	}
	self.Locker.Unlock()

	info.UpdatePoint(req, res)
}

func (self *ConsumerTopInfo) UpdatePoint(req *RPC_ConsumerTopActionReq, res *RPC_ConsumerTopActionRes) {
	if req.SelfInfo == nil {
		res.RetCode = RETCODE_DATA_ERROR
		return
	}
	self.Mu.Lock()
	defer self.Mu.Unlock()

	addpoint := 0
	kill := 0
	info, ok := self.consumerTop[req.SelfInfo.Uid]
	if ok {
		addpoint = req.SelfInfo.Point - info.Point
		info.SvrId = req.SelfInfo.SvrId
		info.SvrName = req.SelfInfo.SvrName
		info.UName = req.SelfInfo.UName
		info.Level = req.SelfInfo.Level
		info.Vip = req.SelfInfo.Vip
		info.Icon = req.SelfInfo.Icon
		info.Portrait = req.SelfInfo.Portrait
		info.Point = req.SelfInfo.Point
		info.KillAll += req.SelfInfo.Kill
		kill = req.SelfInfo.Kill
	} else {
		addpoint = req.SelfInfo.Point
		dbData := new(JS_ConsumerTopUserDB)
		dbData.info = req.SelfInfo
		dbData.Uid = req.SelfInfo.Uid
		dbData.KeyId = self.KeyId
		dbData.SvrId = req.SelfInfo.SvrId
		dbData.info.KillAll = req.SelfInfo.Kill
		kill = req.SelfInfo.Kill
		dbData.Encode()
		self.db_list[req.Uid] = dbData
		dbData.info.Rank = len(self.consumerTopNodeArr) + 1
		db.InsertTable("tbl_consumertop", dbData, 0, true)

		self.consumerTopNodeArr = append(self.consumerTopNodeArr, dbData.info)
		self.consumerTop[req.Uid] = req.SelfInfo
	}

	infoNow, ok := self.consumerTop[req.SelfInfo.Uid]

	for i := infoNow.Rank - 2; i >= 0; i-- {
		if infoNow.Point > self.consumerTopNodeArr[i].Point {
			self.consumerTopNodeArr[i].Rank++
			infoNow.Rank--
			self.consumerTopNodeArr.Swap(infoNow.Rank-1, self.consumerTopNodeArr[i].Rank-1)
		} else {
			break
		}
	}

	res.SelfInfo = infoNow
	if len(self.consumerTopNodeArr) >= 50 {
		res.TopUser = self.consumerTopNodeArr[0:50]
	} else {
		res.TopUser = self.consumerTopNodeArr
	}
	//计算服务器排行
	_, okServer := self.consumerTopServer[infoNow.SvrId]
	if !okServer {
		self.consumerTopServer[infoNow.SvrId] = GetConsumerTopMgr().NewConsumerTopServer(infoNow.SvrId, infoNow.SvrName, infoNow.Step)
		self.consumerTopServer[infoNow.SvrId].Rank = len(self.consumerTopServerNodeArr) + 1
		self.consumerTopServerNodeArr = append(self.consumerTopServerNodeArr, self.consumerTopServer[infoNow.SvrId])
	}
	self.consumerTopServer[infoNow.SvrId].Point += addpoint
	self.consumerTopServer[infoNow.SvrId].Kill += kill

	serverNow, ok := self.consumerTopServer[infoNow.SvrId]
	for i := serverNow.Rank - 2; i >= 0; i-- {
		if serverNow.Point > self.consumerTopServerNodeArr[i].Point {
			self.consumerTopServerNodeArr[i].Rank++
			serverNow.Rank--
			self.consumerTopServerNodeArr.Swap(serverNow.Rank-1, self.consumerTopServerNodeArr[i].Rank-1)
		} else {
			break
		}
	}

	res.TopSvr = self.consumerTopServerNodeArr
	return
}

func (self *ConsumerTopMgr) GetAllRank(req *RPC_ConsumerTopActionReqAll, res *RPC_ConsumerTopActionRes) {
	self.Locker.RLock()
	info, ok := self.ConsumerTopInfo[req.KeyId]
	self.Locker.RUnlock()
	if ok {
		info.GetAllRank(req, res)
	}
}

func (self *ConsumerTopInfo) GetAllRank(req *RPC_ConsumerTopActionReqAll, res *RPC_ConsumerTopActionRes) {
	res.RetCode = RETCODE_OK
	self.Mu.RLock()
	defer self.Mu.RUnlock()
	res.TopUser = self.consumerTopNodeArr
	res.TopSvr = self.consumerTopServerNodeArr
}
