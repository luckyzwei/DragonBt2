package game

//! 配置
type Config struct {
	ServerId   int    `json:"serverid"`   //! 服
	ServerName string `json:"servername"` //! 服务器名称
	GameName   string `json:"gamename"`   //! 游戏名字
	Host       string `json:"wshost"`     //! 服务器
	OpenTime   string `json:"opentime"`   //! 开服时间
	ServerVer  int    `json:"serverver"`  //! 服务器版本
	AdminCode  string `json:"admincode"`  //! 后台token
	GameID     string `json:"gameid"`     //! 经分接入Id
	AppKey     string `json:"appkey"`     //! appkey
	GM         int64  `json:"gm"`         //! gmID

	NetworkCon   *NetworkConfig   `json:"network"`   //! 网络相关
	ServerExtCon *ServerExtConfig `json:"serverext"` //! 服务端组件
	DBCon        *DatabaseConfig  `json:"database"`  //! 数据库
	LogCon       *LoggerConfig    `json:"log"`       //! 日志配置
	AppCon       []AppConfig      `json:"appconfig"` //! 应用配置
}

type LoggerConfig struct {
	MaxFileSize int64 `json:"maxfilesize"` //! 文件长度
	MaxFileNum  int   `json:"maxfilenum"`  //! 文件数量
	LogLevel    int   `json:"loglevel"`    //! 日志等级
	LogConsole  bool  `json:"logconsole"`  //! 是否输出控制台
	SDK         bool  `json:"sdk"`         //! 是否输出SDK日志
}

type NetworkConfig struct {
	BlackIP   []string `json:"blackip"`   //! 黑名单ip
	WhiteID   []int64  `json:"whiteid"`   //! 白名单id
	MaxPlayer int      `json:"maxplayer"` //! 最高在线玩家
	MsgFilter bool     `json:"msgfilter"` //! 过滤消息ID
}

type ServerExtConfig struct {
	IsMaster        bool   `json:"ismaster"`        //! 是否中心服务器
	MasterSvr       string `json:"master"`          //! 跨服服务器
	UpRecord        int    `json:"uprecord"`        //! 是否数据上报
	NumRecord       string `json:"numrecord"`       //! 上报人数url
	Cache           string `json:"cachehost"`       //! 中转服务器
	UpdateTime      int    `json:"updatetime"`      //! 跨服更新时间
	MasterWarnClose int    `json:"masterwarnclose"` //! 是否开启中心服状态变化通知
}

type DatabaseConfig struct {
	DBUser    string `json:"dbuser"`    //! 游戏数据库
	DBLog     string `json:"dblog"`     //! 日志数据库
	DBId      int    `json:"dbid"`      //! dbid
	MaxDBConn int    `json:"maxdbconn"` //! 最大数据库连接数
	Redis     string `json:"redis"`     //! redis地址
	RedisDB   int    `json:"redisdb"`   //! redis db编号
	RedisAuth string `json:"redisauth"` //! redis认证
}

type AppConfig struct {
	AppId    string `json:"appid"`    //! SDK Id
	GameId   string `json:"gameid"`   //! 游戏Id
	AppToken string `json:"apptoken"` //! 经分Token
	AppleId  string `json:"appleid"`  //! 苹果Id
	AppKey   string `json:"appkey"`   //! 登录Token
}

func (self *Config) GetSDKLogGameId(creator string) string {
	return self.GameID
	/*
	switch creator {
	case "sdk_abjuice_ios":
		return SDK_EXT_GAMEID_303
	default:
		return self.GameID
	}
	 */
}
