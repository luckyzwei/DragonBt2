package main

//! 配置
type Config struct {
	Host       string         `json:"wshost"`  //! 服务器
	LogCon     *LoggerConfig  `json:"log"`     //! 日志配置
	ServersCon []ServerConfig `json:"servers"` //! 应用配置
}

type ServerConfig struct {
	Id   int    `json:"serverid"` //! 游戏Id
	Host string `json:"wshost"`   //! 经分Token
}

type LoggerConfig struct {
	MaxFileSize int64 `json:"maxfilesize"` //! 文件长度
	MaxFileNum  int   `json:"maxfilenum"`  //! 文件数量
	LogLevel    int   `json:"loglevel"`    //! 日志等级
	LogConsole  bool  `json:"logconsole"`  //! 是否输出控制台
}
