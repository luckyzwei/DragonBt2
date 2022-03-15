package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"master/center"
	"master/center/tower"
	"master/center/union"
	"master/core"
	"master/db"
	"master/utils"
	"net"
	"net/http"
	"sync"
)

//! 中心服 主app
type MasterApp struct {
	Conf       core.Config     //! 配置
	Shutdown   bool            //! 是否关闭
	WaitGroup  *sync.WaitGroup //! 同步阻塞
	WorldLevel int             //! 世界等级
	init       bool            //! 是否初始化
}

//! 单例模式
var s_MasterApp *MasterApp

func GetMasterApp() *MasterApp {
	if s_MasterApp == nil {
		s_MasterApp = new(MasterApp)
		s_MasterApp.Shutdown = false
		s_MasterApp.WaitGroup = new(sync.WaitGroup)
		//! 全局接口
		core.MasterApp = s_MasterApp
	}

	return s_MasterApp

}

//! 开启服务
func (self *MasterApp) StartService() {
	//! 初始化
	if self.init == false {
		self.Init()
		self.init = true
	}

	//! 启动服务，函数是阻塞模式，后面不会执行
	center.GetCenterApp().StartService()

	//! 注册TCP，端口
	handle, err := net.Listen("tcp", self.Conf.Host)
	if err != nil {
		log.Fatalln("listen rpc fatal error: ", err)
	}

	log.Println("Start Service...", self.Conf.Host)

	/*var battleInfo game.BattleInfo
	battleInfo.Id = 10001
	battleInfo.LevelID = 2
	battleInfo.Time = 102
	battleInfo.UserInfo[0] = new(game.BattleUserInfo)
	battleInfo.UserInfo[1] = new(game.BattleUserInfo)

	str := `{\"Id\":938165852,\"Result\":2,\"Fight\":[{\"rankid\":0,\"uid\":19000892,\"uname\":\"\xe7\x9f\xb3\xe5\xbf\x83\xe7\xa5\x9e\xe5\x89\x91\",\"union\":\"\xe5\x85\x89\xe6\x98\x8e\xe8\xbf\x9c\xe6\x96\xb9\",\"iconid\":1003,\"camp\":1,\"level\":20,\"vip\":2,\"defhero\":[5,7,1,3,8],\"heroinfo\":[{\"heroid\":3002,\"herokeyid\":5,\"color\":1,\"stars\":4,\"levels\":53,\"soldiercolor\":0,\"soldierid\":0,\"skilllevel1\":0,\"skilllevel2\":0,\"skilllevel3\":0,\"skilllevel4\":1,\"skilllevel5\":0,\"skilllevel6\":0,\"fervor1\":0,\"fervor2\":0,\"fervor3\":0,\"fervor4\":0,\"fight\":651953,\"armsskill\":[{\"id\":300204,\"level\":1},{\"id\":300201,\"level\":2},{\"id\":300202,\"level\":1}],\"talentskill\":[],\"army_id\":0,\"maintalent\":0,\"heroquality\":4,\"heroartifactid\":0,\"heroartifactlv\":0,\"heroexclusivelv\":0,\"skin\":0,\"exclusiveunlock\":0,\"isarmy\":0},{\"heroid\":4015,\"herokeyid\":7,\"color\":1,\"stars\":4,\"levels\":22,\"soldiercolor\":0,\"soldierid\":0,\"skilllevel1\":0,\"skilllevel2\":0,\"skilllevel3\":0,\"skilllevel4\":1,\"skilllevel5\":0,\"skilllevel6\":0,\"fervor1\":0,\"fervor2\":0,\"fervor3\":0,\"fervor4\":0,\"fight\":300284,\"armsskill\":[{\"id\":401504,\"level\":1},{\"id\":401501,\"level\":2}],\"talentskill\":[],\"army_id\":0,\"maintalent\":0,\"heroquality\":4,\"heroartifactid\":0,\"heroartifactlv\":0,\"heroexclusivelv\":0,\"skin\":0,\"exclusiveunlock\":0,\"isarmy\":0},{\"heroid\":1012,\"herokeyid\":1,\"color\":1,\"stars\":2,\"levels\":22,\"soldiercolor\":0,\"soldierid\":0,\"skilllevel1\":0,\"skilllevel2\":0,\"skilllevel3\":0,\"skilllevel4\":1,\"skilllevel5\":0,\"skilllevel6\":0,\"fervor1\":0,\"fervor2\":0,\"fervor3\":0,\"fervor4\":0,\"fight\":227183,\"armsskill\":[{\"id\":101204,\"level\":1},{\"id\":101201,\"level\":2}],\"talentskill\":[],\"army_id\":0,\"maintalent\":0,\"heroquality\":2,\"heroartifactid\":0,\"heroartifactlv\":0,\"heroexclusivelv\":0,\"skin\":0,\"exclusiveunlock\":0,\"isarmy\":0},{\"heroid\":1015,\"herokeyid\":3,\"color\":1,\"stars\":2,\"levels\":24,\"soldiercolor\":0,\"soldierid\":0,\"skilllevel1\":0,\"skilllevel2\":0,\"skilllevel3\":0,\"skilllevel4\":1,\"skilllevel5\":0,\"skilllevel6\":0,\"fervor1\":0,\"fervor2\":0,\"fervor3\":0,\"fervor4\":0,\"fight\":191154,\"armsskill\":[{\"id\":101504,\"level\":1},{\"id\":101501,\"level\":2}],\"talentskill\":[],\"army_id\":0,\"maintalent\":0,\"heroquality\":2,\"heroartifactid\":0,\"heroartifactlv\":0,\"heroexclusivelv\":0,\"skin\":0,\"exclusiveunlock\":0,\"isarmy\":0},{\"heroid\":5006,\"herokeyid\":8,\"color\":1,\"stars\":4,\"levels\":21,\"soldiercolor\":0,\"soldierid\":0,\"skilllevel1\":0,\"skilllevel2\":0,\"skilllevel3\":0,\"skilllevel4\":1,\"skilllevel5\":0,\"skilllevel6\":0,\"fervor1\":0,\"fervor2\":0,\"fervor3\":0,\"fervor4\":0,\"fight\":313018,\"armsskill\":[{\"id\":500604,\"level\":1},{\"id\":500601,\"level\":2}],\"talentskill\":[],\"army_id\":0,\"maintalent\":0,\"heroquality\":4,\"heroartifactid\":0,\"heroartifactlv\":0,\"heroexclusivelv\":0,\"skin\":0,\"exclusiveunlock\":0,\"isarmy\":0}],\"heroparam\":[{\"heroid\":3002,\"param\":[0,11894.300000000001,891.0000000000001,177,0,0,200,0,0,0,0,0,0,0,0,0,0,0,0,0,80,10000,250,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"hp\":11894.300000000001,\"energy\":0,\"pos\":0,\"ext\":null},{\"heroid\":4015,\"param\":[0,5934.500000000001,379.50000000000006,61,0,0,200,0,0,0,0,0,0,0,0,0,0,0,0,0,0,10000,250,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"hp\":5934.500000000001,\"energy\":0,\"pos\":0,\"ext\":null},{\"heroid\":1012,\"param\":[0,4039.2000000000003,236.50000000000003,64,0,0,200,0,0,0,0,0,0,0,0,0,0,0,0,70,0,10000,250,0,0,0,0,0,0,0,0,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"hp\":4039.2000000000003,\"energy\":0,\"pos\":0,\"ext\":null},{\"heroid\":1015,\"param\":[0,2554.2000000000003,320.1,45,0,0,200,0,0,0,0,0,0,0,0,0,0,0,0,0,0,10000,250,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"hp\":2554.2000000000003,\"energy\":0,\"pos\":0,\"ext\":null},{\"heroid\":5006,\"param\":[0,6091.8,447.70000000000005,70,0,0,200,0,0,0,0,0,0,0,0,0,0,0,0,0,40,10000,250,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"hp\":6091.8,\"energy\":0,\"pos\":0,\"ext\":null}],\"deffight\":1683592,\"fightteam\":0,\"fightteampos\":{\"fightpos\":[5,7,1,3,8],\"hydraid\":0},\"portrait\":1000},{\"rankid\":0,\"uid\":19000745,\"uname\":\"\xe6\x99\xb4\xe6\x99\xb4\xe5\xb0\x8f\xe5\x85\xac\xe4\xb8\xbb\",\"union\":\"\xe5\x85\xbb\xe8\x80\x81\xe9\x99\xa2\",\"iconid\":1003,\"camp\":1,\"level\":27,\"vip\":2,\"defhero\":[6,101,69,25,45],\"heroinfo\":[{\"heroid\":4008,\"herokeyid\":6,\"color\":1,\"stars\":6,\"levels\":62,\"soldiercolor\":0,\"soldierid\":0,\"skilllevel1\":0,\"skilllevel2\":0,\"skilllevel3\":0,\"skilllevel4\":1,\"skilllevel5\":0,\"skilllevel6\":0,\"fervor1\":0,\"fervor2\":0,\"fervor3\":0,\"fervor4\":0,\"fight\":1311533,\"armsskill\":[{\"id\":400804,\"level\":1},{\"id\":400801,\"level\":2},{\"id\":400802,\"level\":1},{\"id\":400803,\"level\":1}],\"talentskill\":[],\"army_id\":0,\"maintalent\":0,\"heroquality\":6,\"heroartifactid\":0,\"heroartifactlv\":0,\"heroexclusivelv\":0,\"skin\":0,\"exclusiveunlock\":0,\"isarmy\":0},{\"heroid\":3002,\"herokeyid\":101,\"color\":1,\"stars\":4,\"levels\":71,\"soldiercolor\":0,\"soldierid\":0,\"skilllevel1\":0,\"skilllevel2\":0,\"skilllevel3\":0,\"skilllevel4\":1,\"skilllevel5\":0,\"skilllevel6\":0,\"fervor1\":0,\"fervor2\":0,\"fervor3\":0,\"fervor4\":0,\"fight\":955774,\"armsskill\":[{\"id\":300204,\"level\":1},{\"id\":300201,\"level\":2},{\"id\":300202,\"level\":1},{\"id\":300203,\"level\":1}],\"talentskill\":[],\"army_id\":0,\"maintalent\":0,\"heroquality\":4,\"heroartifactid\":0,\"heroartifactlv\":0,\"heroexclusivelv\":0,\"skin\":0,\"exclusiveunlock\":0,\"isarmy\":0},{\"heroid\":1008,\"herokeyid\":69,\"color\":1,\"stars\":4,\"levels\":61,\"soldiercolor\":0,\"soldierid\":0,\"skilllevel1\":0,\"skilllevel2\":0,\"skilllevel3\":0,\"skilllevel4\":1,\"skilllevel5\":0,\"skilllevel6\":0,\"fervor1\":0,\"fervor2\":0,\"fervor3\":0,\"fervor4\":0,\"fight\":820062,\"armsskill\":[{\"id\":100804,\"level\":1},{\"id\":100801,\"level\":2},{\"id\":100802,\"level\":1},{\"id\":100803,\"level\":1}],\"talentskill\":[],\"army_id\":0,\"maintalent\":0,\"heroquality\":4,\"heroartifactid\":0,\"heroartifactlv\":0,\"heroexclusivelv\":0,\"skin\":0,\"exclusiveunlock\":0,\"isarmy\":0},{\"heroid\":5006,\"herokeyid\":25,\"color\":1,\"stars\":5,\"levels\":64,\"soldiercolor\":0,\"soldierid\":0,\"skilllevel1\":0,\"skilllevel2\":0,\"skilllevel3\":0,\"skilllevel4\":1,\"skilllevel5\":0,\"skilllevel6\":0,\"fervor1\":0,\"fervor2\":0,\"fervor3\":0,\"fervor4\":0,\"fight\":1027640,\"armsskill\":[{\"id\":500604,\"level\":1},{\"id\":500601,\"level\":2},{\"id\":500602,\"level\":1},{\"id\":500603,\"level\":1}],\"talentskill\":[],\"army_id\":0,\"maintalent\":0,\"heroquality\":5,\"heroartifactid\":0,\"heroartifactlv\":0,\"heroexclusivelv\":0,\"skin\":0,\"exclusiveunlock\":0,\"isarmy\":0},{\"heroid\":1010,\"herokeyid\":45,\"color\":1,\"stars\":4,\"levels\":62,\"soldiercolor\":0,\"soldierid\":0,\"skilllevel1\":0,\"skilllevel2\":0,\"skilllevel3\":0,\"skilllevel4\":1,\"skilllevel5\":0,\"skilllevel6\":0,\"fervor1\":0,\"fervor2\":0,\"fervor3\":0,\"fervor4\":0,\"fight\":854840,\"armsskill\":[{\"id\":101004,\"level\":1},{\"id\":101001,\"level\":2},{\"id\":101002,\"level\":1},{\"id\":101003,\"level\":1}],\"talentskill\":[],\"army_id\":0,\"maintalent\":0,\"heroquality\":4,\"heroartifactid\":0,\"heroartifactlv\":0,\"heroexclusivelv\":0,\"skin\":0,\"exclusiveunlock\":0,\"isarmy\":0}],\"heroparam\":[{\"heroid\":4008,\"param\":[0,26149.38,1650.0000000000002,381,0,0,200,0,0,0,0,0,0,0,0,0,0,0,0,410,50,10000,250,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"hp\":26149.38,\"energy\":0,\"pos\":0,\"ext\":null},{\"heroid\":3002,\"param\":[0,17581.300000000003,1323.3000000000002,246,0,0,200,0,0,0,0,0,0,0,0,0,0,0,0,0,40,10000,250,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"hp\":17581.300000000003,\"energy\":0,\"pos\":0,\"ext\":null},{\"heroid\":1008,\"param\":[0,16519.800000000003,994.4000000000001,197,0,0,200,0,0,0,0,0,0,0,0,0,0,0,0,0,40,10000,250,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"hp\":16519.800000000003,\"energy\":0,\"pos\":0,\"ext\":null},{\"heroid\":5006,\"param\":[0,19790.100000000002,1483.9,273,0,0,200,0,0,0,0,0,0,0,0,0,0,0,0,0,0,10000,250,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"hp\":19790.100000000002,\"energy\":0,\"pos\":0,\"ext\":null},{\"heroid\":1010,\"param\":[0,16523.100000000002,1091.2,209,0,0,200,0,0,0,0,0,0,0,0,0,0,0,0,0,40,10000,250,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"hp\":16523.100000000002,\"energy\":0,\"pos\":0,\"ext\":null}],\"deffight\":4969849,\"fightteam\":0,\"fightteampos\":{\"fightpos\":[6,101,69,25,45],\"hydraid\":0},\"portrait\":1000}],\"Info\":[[{\"id\":3002,\"hp\":0,\"rage\":370,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":4015,\"hp\":0,\"rage\":478,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":1012,\"hp\":0,\"rage\":0,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":1015,\"hp\":0,\"rage\":952,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":5006,\"hp\":0,\"rage\":756,\"damage\":0,\"takedamage\":0,\"healing\":0}],[{\"id\":4008,\"hp\":9818,\"rage\":817,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":3002,\"hp\":9850,\"rage\":536,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":1008,\"hp\":9572,\"rage\":483,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":5006,\"hp\":9957,\"rage\":902,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":1010,\"hp\":9824,\"rage\":508,\"damage\":0,\"takedamage\":0,\"healing\":0}]],\"Random\":1609381658,\"Time\":1609766641,\"ResultDetail\":{\"fightid\":938165852,\"info\":[[{\"id\":3002,\"hp\":0,\"rage\":370,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":4015,\"hp\":0,\"rage\":478,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":1012,\"hp\":0,\"rage\":0,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":1015,\"hp\":0,\"rage\":952,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":5006,\"hp\":0,\"rage\":756,\"damage\":0,\"takedamage\":0,\"healing\":0}],[{\"id\":4008,\"hp\":9818,\"rage\":817,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":3002,\"hp\":9850,\"rage\":536,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":1008,\"hp\":9572,\"rage\":483,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":5006,\"hp\":9957,\"rage\":902,\"damage\":0,\"takedamage\":0,\"healing\":0},{\"id\":1010,\"hp\":9824,\"rage\":508,\"damage\":0,\"takedamage\":0,\"healing\":0}]],\"time\":7500,\"winner\":1},\"CityId\":0,\"SecKill\":0,\"TeamA\":0,\"TeamB\":0,\"Group\":0,\"IsSet\":true}`


	//encodeStr1 := utils.HF_JtoB(&battleInfo)
	str = string(bytes.Replace([]byte(str), []byte(`\\`), []byte(`\`), -1))
	str = string(bytes.Replace([]byte(str), []byte(`\"`), []byte(`"`), -1))
	str = string(bytes.Replace([]byte(str), []byte(`\x`), []byte(``), -1))

	utils.LogDebug("encode Str:", len(str), string([]byte(str)))

	encodeStr := utils.HF_CompressAndBase64([]byte(str))
	utils.LogDebug("encode Str:", len(encodeStr), encodeStr)

	decodeBytes := utils.HF_Base64AndDecompress(encodeStr)
	utils.LogDebug("decode Str:", len(decodeBytes), string(decodeBytes))

	var battleInfo1 game.FightResult
	json.Unmarshal(decodeBytes, &battleInfo1)

	json.Unmarshal([]byte(str), &battleInfo1)*/

	log.Fatal(http.Serve(handle, nil))
}

//! 关闭服务
//! 停止服务器
func (self *MasterApp) StopService() {

	//! 关闭中心服逻辑
	center.GetCenterApp().StopService()
}

//! 初始化系统
func (self *MasterApp) Init() {
	//! 初始化配置
	self.InitConfig()

	//! 连接数据库
	utils.LogDebug("连接数据库...")
	self.ConnectDB()

	//! 启动常驻服务
	utils.LogDebug("启动逻辑处理...")
	go center.GetCenterApp().OnTimer()
	go center.GetCenterApp().OnLogicEvent()

	go utils.GetLoggerCheckMgr().Run()

	go union.GetUnionMgr().Run()

	//! 迁移爬塔的战斗数据
	go tower.GetTowerMgr().RunMigrate()

	//! 注册服务
	utils.LogDebug("注册并启动服务...")
	center.GetCenterApp().RegisterService()
}

func (self *MasterApp) ConnectDB() {
	//! 初始化redis模块
	//! 连接redis
	ret := db.GetRedisMgr().Init(self.Conf.DBConf.Redis, self.Conf.DBConf.RedisDB, self.Conf.DBConf.RedisAuth)
	if ret == false {
		log.Fatal("redis init err...")
		return
	}

	//! 初始化数据模块
	db.GetDBMgr().Init(self.Conf.DBConf.DBUser, self.Conf.DBConf.DBLog)

	//! 检查数据库
	GetSqlMgr().CheckMysql()
}

func (self *MasterApp) Close() {
	self.Shutdown = true

	//! 数据库关闭OK
	utils.LogInfo("db close...")
	db.GetDBMgr().Close()

	//! 等待处理结果，结束则完成
	utils.LogInfo("master close...")
	self.WaitGroup.Wait()

	//! 服务器停止
	utils.LogFatal("server shutdown")
}

func (self *MasterApp) IsClosed() bool {
	return self.Shutdown
}

func (self *MasterApp) GetConfig() *core.Config {
	return &self.Conf
}

func (self *MasterApp) GetPlayerOnline(serverId int) int {
	return 0
}

//! 载入配置文件
func (self *MasterApp) InitConfig() {
	configFile, err := ioutil.ReadFile("./config.json") ///尝试打开配置文件
	if err != nil {
		log.Fatal("config err 1")
	}
	err = json.Unmarshal(configFile, &self.Conf)
	if err != nil {
		log.Fatal("server InitConfig err:", err.Error())
	}

	if self.Conf.ServerVer == 0 {
		self.Conf.ServerVer = 1000
	}

	utils.GetLogMgr().SetLevel(self.Conf.LogConf.LogLevel, self.Conf.LogConf.LogConsole)

}

func (self *MasterApp) Wait() {
	self.WaitGroup.Add(1)
}

func (self *MasterApp) Done() {
	self.WaitGroup.Done()
}

//! 获取服务器世界等级
func (self *MasterApp) GetWorldLevel(refresh bool) int {
	return db.GetDBMgr().GetWorldLevel(refresh)
}

//! 开服天数，开服当天5点开始计算
func (self *MasterApp) GetOpenTime() int {
	return 0
}

//! 开服时间，时间戳
func (self *MasterApp) GetOpenServer() int64 {
	return 0
}
