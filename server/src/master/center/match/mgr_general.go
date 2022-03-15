package match

import (
	"encoding/json"
	"fmt"
	"master/db"
	"master/utils"
	"sort"
	"sync"
	"time"
)

type Js_GeneralUser struct {
	Uid     int64  `json:"uid"`
	SvrId   int    `json:"svrid"`
	SvrName string `json:"svrname"`
	UName   string `json:"uname"`
	Level   int    `json:"level"`
	Vip     int    `json:"vip"`
	Icon    int    `json:"icon"`
	Point   int    `json:"point"`
	Rank    int    `json:"rank"`
	KeyId   int    `json:"keyid"`
	Time    int64  `json:"time"`
}

type Js_GeneralUserDB struct {
	Uid   int64  `json:"uid"`
	KeyId int    `json:"keyid"`
	SvrId int    `json:"svrid"`
	Info  string `json:"info"`

	info *Js_GeneralUser
	db.DataUpdate
}

type GeneralRecord struct {
	Uid        int64  `json:"uid"`
	SvrId      int    `json:"svrid"`
	UName      string `json:"uname"`
	Item       int    `json:"item"`
	Num        int    `json:"num"`
	Time       int64  `json:"time"`
	RecordType int    `json:"recordtype"`
}

type GeneralUserArr []*Js_GeneralUser

func (s GeneralUserArr) Len() int      { return len(s) }
func (s GeneralUserArr) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s GeneralUserArr) Less(i, j int) bool {
	if s[i].Point == s[j].Point {
		return s[i].Time < s[j].Time
	}
	return s[i].Point > s[j].Point
}

type GeneralInfo struct {
	KeyId          int    `json:"keyid"`
	GeneralUserTop string `json:"generalusertop"`

	Mu                    *sync.RWMutex
	generalUserTop        map[int64]*Js_GeneralUser
	generalUserTopNodeArr GeneralUserArr
	generalRecord         []*GeneralRecord
	db_list               map[int64]*Js_GeneralUserDB
}

type GeneralMgr struct {
	Locker      *sync.RWMutex
	GeneralInfo map[int]*GeneralInfo
}

var generalMgr *GeneralMgr = nil

func GetGeneralMgr() *GeneralMgr {
	if generalMgr == nil {
		generalMgr = new(GeneralMgr)
		generalMgr.GeneralInfo = make(map[int]*GeneralInfo)
		generalMgr.Locker = new(sync.RWMutex)
	}

	return generalMgr
}

func (self *Js_GeneralUserDB) Encode() {
	self.Info = utils.HF_JtoA(self.info)
}

func (self *Js_GeneralUserDB) Decode() {
	json.Unmarshal([]byte(self.Info), &self.info)
}

// 存储数据库
func (self *GeneralMgr) OnSave() {
	self.Locker.Lock()
	defer self.Locker.Unlock()
	for _, v := range self.GeneralInfo {
		v.Save()
	}
}

func (self *GeneralInfo) Save() {
	self.Mu.Lock()
	defer self.Mu.Unlock()
	for _, v := range self.db_list {
		v.Encode()
		v.Update(true, false)
	}
}

func (self *GeneralMgr) GetAllData() {
	self.LoadGeneral()
}

func (self *GeneralMgr) LoadGeneral() {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	queryStr := fmt.Sprintf("select uid,keyid,svrid,info from `tbl_generalex`;")
	var msg Js_GeneralUserDB
	res := db.GetDBMgr().DBUser.GetAllData(queryStr, &msg)

	for i := 0; i < len(res); i++ {
		data := res[i].(*Js_GeneralUserDB)
		if data.KeyId > 0 {

			_, ok := self.GeneralInfo[data.KeyId]
			if !ok {
				self.GeneralInfo[data.KeyId] = self.NewGeneralInfo(data.KeyId)
			}

			data.Decode()
			if data.info == nil {
				continue
			}
			data.Init("tbl_generalex", data, true)
			self.GeneralInfo[data.KeyId].generalUserTop[data.Uid] = data.info
			self.GeneralInfo[data.KeyId].db_list[data.Uid] = data
		}
	}

	for _, v := range self.GeneralInfo {
		v.MakeArr()
	}
}

func (self *GeneralMgr) NewGeneralInfo(KeyId int) *GeneralInfo {
	data := new(GeneralInfo)
	data.KeyId = KeyId
	data.Mu = new(sync.RWMutex)
	data.generalUserTop = make(map[int64]*Js_GeneralUser)
	data.generalUserTopNodeArr = make([]*Js_GeneralUser, 0)
	data.generalRecord = make([]*GeneralRecord, 0)
	data.db_list = make(map[int64]*Js_GeneralUserDB, 0)
	return data
}

func (self *GeneralInfo) MakeArr() {
	self.Mu.Lock()
	defer self.Mu.Unlock()

	self.generalUserTopNodeArr = make([]*Js_GeneralUser, 0)
	for _, v := range self.generalUserTop {
		self.generalUserTopNodeArr = append(self.generalUserTopNodeArr, v)
	}
	sort.Sort(self.generalUserTopNodeArr)

	for i := 0; i < len(self.generalUserTopNodeArr); i++ {
		self.generalUserTopNodeArr[i].Rank = i + 1
	}
}

func (self *GeneralMgr) UpdatePoint(req *RPC_GeneralActionReq, res *RPC_GeneralActionRes) {
	if req.SelfInfo == nil {
		res.RetCode = RETCODE_DATA_ERROR
		return
	}

	self.Locker.Lock()
	info, ok := self.GeneralInfo[req.SelfInfo.KeyId]
	if !ok {
		info = self.NewGeneralInfo(req.SelfInfo.KeyId)
		self.GeneralInfo[info.KeyId] = info
	}
	self.Locker.Unlock()

	info.UpdatePoint(req, res)
}

func (self *GeneralInfo) UpdatePoint(req *RPC_GeneralActionReq, res *RPC_GeneralActionRes) {
	if req.SelfInfo == nil {
		res.RetCode = RETCODE_DATA_ERROR
		return
	}
	self.Mu.Lock()
	defer self.Mu.Unlock()

	info, ok := self.generalUserTop[req.SelfInfo.Uid]
	if ok {
		info.UName = req.SelfInfo.UName
		info.Level = req.SelfInfo.Level
		info.Vip = req.SelfInfo.Vip
		info.Icon = req.SelfInfo.Icon
		info.SvrName = req.SelfInfo.SvrName

		if info.Point > req.SelfInfo.Point {
			return
		}
		info.Point = req.SelfInfo.Point
		info.Time = time.Now().Unix()
	} else {
		info = req.SelfInfo
		info.Rank = len(self.generalUserTopNodeArr) + 1
		self.generalUserTopNodeArr = append(self.generalUserTopNodeArr, info)
		self.generalUserTop[req.SelfInfo.Uid] = info

		dbData := new(Js_GeneralUserDB)
		dbData.info = req.SelfInfo
		dbData.Uid = req.SelfInfo.Uid
		dbData.KeyId = self.KeyId
		dbData.SvrId = req.SelfInfo.SvrId
		dbData.Encode()
		self.db_list[req.Uid] = dbData
		db.InsertTable("tbl_generalex", dbData, 0, true)
		dbData.Init("tbl_generalex", dbData, true)
	}

	for i := info.Rank - 2; i >= 0; i-- {
		if info.Point > self.generalUserTopNodeArr[i].Point {
			self.generalUserTopNodeArr[i].Rank++
			info.Rank--
			self.generalUserTopNodeArr.Swap(info.Rank-1, self.generalUserTopNodeArr[i].Rank-1)
		} else {
			break
		}
	}

	self.generalRecord = append(self.generalRecord, req.GeneralRecord...)
	size := len(self.generalRecord)
	if size > RECORD_MAX {
		self.generalRecord = self.generalRecord[size-RECORD_MAX:]
	}

	res.RetCode = RETCODE_OK
	if len(self.generalUserTopNodeArr) > 50 {
		res.RankInfo = self.generalUserTopNodeArr[:50]
	} else {
		res.RankInfo = self.generalUserTopNodeArr
	}
	res.SelfInfo = self.generalUserTop[req.Uid]
	res.GeneralRecord = self.generalRecord
	return
}

func (self *GeneralMgr) GetAllRank(req *RPC_GeneralActionReqAll, res *RPC_GeneralActionRes) {
	self.Locker.RLock()
	info, ok := self.GeneralInfo[req.KeyId]
	self.Locker.RUnlock()
	if ok {
		info.GetAllRank(req, res)
	}
}

func (self *GeneralInfo) GetAllRank(req *RPC_GeneralActionReqAll, res *RPC_GeneralActionRes) {
	res.RetCode = RETCODE_OK

	self.Mu.RLock()
	defer self.Mu.RUnlock()

	res.RankInfo = self.generalUserTopNodeArr
	res.GeneralRecord = self.generalRecord
}
