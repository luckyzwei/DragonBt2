package game

import (
	"encoding/json"
	"sync"
	//"time"
)

type San_Fight struct {
	Id   int
	Info string

	info JS_Fight

	DataUpdate

	Locker *sync.RWMutex
}

type JS_Fight struct {
	Act  string `json:"act"`
	Def  string `json:"def"`
	Time int64  `json:"time"`
	Id   int64  `json:"id"`
}

type FightCountMgr struct {
	Sql_Fight       []*JS_Fight
	Sql_FightResult map[int64]*JS_FightSetData

	Id int64

	Locker   *sync.RWMutex
	LockerId *sync.RWMutex
}

//结果
type JS_FightCountHeroData struct {
	Hp    []int `json:"hp"`
	Maxhp []int `json:"maxhp"`
	Rage  []int `json:"rage`
}

type JS_FightSetData struct {
	Herodata []JS_FightCountHeroData `json:"herodata"`
	State    int                     `json:"state"` //0是还没开始处理 2是处理中 3处理结束
	Star     bool                    `json:"star"`
}

func (self *San_Fight) Decode() { //! 将数据库数据写入data
	json.Unmarshal([]byte(self.Info), &self.info)
}

func (self *San_Fight) Encode() { //! 将data数据写入数据库
	self.Info = HF_JtoA(&self.info)
}

var fightcountsingleton *FightCountMgr = nil

//! public
func GetFightCountMgr() *FightCountMgr {
	if fightcountsingleton == nil {
		fightcountsingleton = new(FightCountMgr)
		fightcountsingleton.Locker = new(sync.RWMutex)
		fightcountsingleton.Sql_Fight = make([]*JS_Fight, 0)
		fightcountsingleton.Sql_FightResult = make(map[int64]*JS_FightSetData)

		fightcountsingleton.LockerId = new(sync.RWMutex)

		str := `{"uname":"robot","soldiertype2":[[0,0,0],[0,0,0],[0,0,0]],"befight":0,"heroinfo":[{"fight":0,"skilllevel5":0,"levels":10,"color":1,"skilllevel2":0,"heroid":10021,"skilllevel1":0,"fervor1":0,"fervor2":0,"fervor3":0,"skilllevel3":0,"skilllevel4":0,"fervor4":0,"stars":1,"soldiercolor":0,"skilllevel6":0},{"fight":0,"skilllevel5":0,"levels":10,"color":1,"skilllevel2":0,"heroid":40061,"skilllevel1":0,"fervor1":0,"fervor2":0,"fervor3":0,"skilllevel3":0,"skilllevel4":0,"fervor4":0,"stars":1,"soldiercolor":0,"skilllevel6":0},{"fight":0,"skilllevel5":0,"levels":10,"color":1,"skilllevel2":0,"heroid":40081,"skilllevel1":0,"fervor1":0,"fervor2":0,"fervor3":0,"skilllevel3":0,"skilllevel4":0,"fervor4":0,"stars":1,"soldiercolor":0,"skilllevel6":0}],"morale":60,"deffight":1247,"iconid":1002,"soldiertype1":[[0,0,0],[0,0,0],[0,0,0]],"uid":1,"heroparam":[{"param":[56,42,14,4,3,1,3648,282,139.2,42,55.2,27417,200,0,0,0,0,0,0,0,0,0,0,0,0,2000,0,0,0,0,0,0],"heroid":10021},{"param":[14,14,56,1,1,4,2528,42,55.2,282,139.2,15735,170,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"heroid":40061},{"param":[14,28,42,1,2,3,3290,42,56.2,218.8,112.2,15376,200,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"heroid":40081}],"rankid":32673,"soldiertype3":[[0,0,0],[0,0,0],[0,0,0]],"level":10,"defhero":[1002,4006,4008]}`

		for i := 0; i < 50; i++ {
			fightcountsingleton.ProcData(str, str)
		}

	}

	return fightcountsingleton
}

func (self *FightCountMgr) Save() {
	/*self.Locker.Lock()
	defer self.Locker.Unlock()

	for i := 0; i < len(self.Sql_Fight); i++ {
		self.Sql_Fight[i].Update()
	}*/
}

func (self *FightCountMgr) GetData() {

	/*self.Locker.Lock()
	defer self.Locker.Unlock()
	var dip Sql_Fight
	sql := fmt.Sprintf("select * from `san_mail`")
	res := GetServer().DBUser.GetAllData(sql, &dip)

	for i := 0; i < len(res); i++ {
		data := res[i].(*Sql_Fight)
		data.Init("san_mail", data)
		data.Locker = new(sync.RWMutex)

		//data.info.Mailid = data.Id
		self.Sql_Fight = append(self.Sql_Fight, data)
	}*/
}

func (self *FightCountMgr) GetId() int64 {

	self.LockerId.Lock()
	defer self.LockerId.Unlock()
	self.Id++
	return self.Id
}

//计算客户端调用
func (self *FightCountMgr) GetOne() *JS_Fight {

	self.Locker.Lock()
	defer self.Locker.Unlock()

	if len(self.Sql_Fight) == 0 {
		return nil
	}
	ret := self.Sql_Fight[0]

	self.Sql_FightResult[ret.Id].State = 2

	if len(self.Sql_Fight) == 1 {
		self.Sql_Fight = make([]*JS_Fight, 0)
	} else {
		self.Sql_Fight = self.Sql_Fight[1:len(self.Sql_Fight)]
	}

	return ret

}

// 等于0表示计算成功，1表示找不到id的数据 2表示处理中 3没有开始处理
func (self *FightCountMgr) GetResult(id int64) (int, *JS_FightSetData) {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	result, ok := self.Sql_FightResult[id]

	if ok {
		delete(self.Sql_FightResult, id)

		if result.State == 0 {
			return 3, nil
		}

		if result.State == 2 {
			return 2, nil
		}

		if result.State == 3 {
			return 0, result
		} else {
			return -1, nil
		}

	}

	return 1, nil
}

func (self *FightCountMgr) ProcData(act string, def string) {

	self.Locker.Lock()
	defer self.Locker.Unlock()

	data := new(JS_Fight)
	data.Act = act
	data.Def = def
	data.Time = TimeServer().Unix()
	data.Id = self.GetId()

	self.Sql_Fight = append(self.Sql_Fight, data)

	result := new(JS_FightSetData)

	self.Sql_FightResult[data.Id] = result

}

func (self *FightCountMgr) SetRetData(id int64, ret string) {
	self.Locker.Lock()
	defer self.Locker.Unlock()

	LogDebug(ret)
	data := new(JS_FightSetData)
	json.Unmarshal([]byte(ret), data)
	LogDebug(data)

	_, ok := self.Sql_FightResult[id]

	if ok {
		self.Sql_FightResult[id] = data
		self.Sql_FightResult[id].State = 3
	}
}
