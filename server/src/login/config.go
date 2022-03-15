package main

//! 配置
type Config struct {
	Host   string          `json:"wshost"`    //! 服务器
	LogCon *LoggerConfig   `json:"log"`       //! 日志配置
	DBCon  *DatabaseConfig `json:"database"`  //! 数据库
	GameID string          `json:"gameid"`    //! 经分接入Id
	AppKey string          `json:"appkey"`    //! appkey
	AppCon []AppConfig     `json:"appconfig"` //! 应用配置
}

type AppConfig struct {
	AppId    string `json:"appid"`    //! SDK Id
	GameId   string `json:"gameid"`   //! 游戏Id
	AppToken string `json:"apptoken"` //! 经分Token
	AppleId  string `json:"appleid"`  //! 苹果Id
	AppKey   string `json:"appkey"`   //! 登录Token
}

type DatabaseConfig struct {
	DBUser    string `json:"dbuser"`    //! 游戏数据库
	MaxDBConn int    `json:"maxdbconn"` //! 最大数据库连接数
	Redis     string `json:"redis"`     //! redis地址
	RedisDB   int    `json:"redisdb"`   //! redis db编号
	RedisAuth string `json:"redisauth"` //! redis认证
}

type LoggerConfig struct {
	MaxFileSize int64 `json:"maxfilesize"` //! 文件长度
	MaxFileNum  int   `json:"maxfilenum"`  //! 文件数量
	LogLevel    int   `json:"loglevel"`    //! 日志等级
	LogConsole  bool  `json:"logconsole"`  //! 是否输出控制台
}

// 获取安卓服对应的gameId
func (self *Config) GetAndroidGameId(appId string) string {
	return self.GameID
}

// 获取Ios服对应的gameId
func (self *Config) GetIosGameId(appId string) string {
	if len(self.AppCon) <= 0 {
		LogError("len(self.AppConfig) <= 0")
		return "10000008"
	}

	for _, v := range self.AppCon {
		if v.AppId == appId {
			return v.GameId
		}
	}
	return self.AppCon[0].GameId
}

func (self *Config) GetGameIdByAppId(appId string) string {
	if appId == "" {
		return self.GetAndroidGameId(appId)
	} else {
		return self.GetIosGameId(appId)
	}
}

// Ios专用
func (self *Config) GetAppKeyByAppId(appId string) string {
	if len(self.AppCon) <= 0 {
		LogError("len(self.AppConfig) <= 0")
		return "6d0f09494f4e2935ceb971f4048db965"
	}

	for _, v := range self.AppCon {
		if v.AppId == appId {
			return v.AppKey
		}
	}
	return self.AppCon[0].AppKey
}
