package core

//! 配置
type Config struct {
	ServerId   int    `json:"serverid"`   //! 服
	ServerName string `json:"servername"` //! 服务器名称
	Host       string `json:"wshost"`     //! 服务器
	OpenTime   string `json:"opentime"`   //! 开服时间
	ServerVer  int    `json:"serverver"`  //! 服务器版本

	DBConf  *DatabaseConfig `json:"database"` //! 数据库
	LogConf *LoggerConfig   `json:"log"`
}

type LoggerConfig struct {
	MaxFileSize int64 `json:"maxfilesize"` //! 文件长度
	MaxFileNum  int   `json:"maxfilenum"`  //! 文件数量
	LogLevel    int   `json:"loglevel"`    //! 日志等级
	LogConsole  bool  `json:"logconsole"`  //! 是否输出控制台
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
