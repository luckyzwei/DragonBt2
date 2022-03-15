package db

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

const (
	expireOption    = "EX"
	notExistsOption = "NX"
	matchOption     = "MATCH"
	countOption     = "COUNT"

	setCommand          = "SET"
	delCommand          = "DEL"
	getCommand          = "GET"
	keysCommand         = "KEYS"
	pingCommand         = "PING"
	echoCommand         = "ECHO"
	infoCommand         = "INFO"
	hSetCommand         = "HSET"
	hGetCommand         = "HGET"
	hmSetCommand        = "HMSET"
	hDelCommand         = "HDEL"
	hLenCommand         = "HLEN"
	hKeysCommand        = "HKEYS"
	scanCommand         = "SCAN"
	hScanCommand        = "HSCAN"
	getRangeCommand     = "GETRANGE"
	setRangeCommand     = "SETRANGE"
	expireCommand       = "EXPIRE"
	existsCommand       = "EXISTS"
	hExistsCommand      = "HEXISTS"
	hGetAllCommand      = "HGETALL"
	incrByCommand       = "INCRBY"
	incrByFloatCommand  = "INCRBYFLOAT"
	hIncrByCommand      = "HINCRBY"
	hIncrByFloatCommand = "HINCRBYFLOAT"
	TTL                 = "TTL"
)

const (
	Max_Redis_Idle_Conn   = 50
	Max_Redis_Active_Conn = 1000
	Max_Redis_Idle_Time   = 180
)

// redis 管理器
type RedisMgr struct {
	RedisPool *redis.Pool //! redis 连接池
	Host      string      //! IP+端口，127.0.0.1:6379
	Index     int         //! redis db index
	Auth      string      //! redis 密码
	InitFlag  bool        //! 是否初始化
}

var s_redismgr *RedisMgr = nil

func GetRedisMgr() *RedisMgr {
	if s_redismgr == nil {
		s_redismgr = new(RedisMgr)
		s_redismgr.InitFlag = false
	}

	return s_redismgr
}

//! 初始化 Redis 连接池
func (self *RedisMgr) Init(host string, index int, auth string) bool {
	self.Host = host
	self.Index = index
	self.Auth = auth

	self.RedisPool = NewPool(self.Host, self.Index, self.Auth)
	if self.RedisPool == nil {
		return false
	} else {
		self.InitFlag = true
		return true
	}
}

func (self *RedisMgr) IsInit() bool {
	return self.InitFlag
}

//! 获得 Redis 连接池
func (self *RedisMgr) GetRedisConn() redis.Conn {
	return self.RedisPool.Get()
}

// 重写生成连接池方法
func NewPool(ip string, db int, auth string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     Max_Redis_Idle_Conn,
		MaxActive:   Max_Redis_Active_Conn, // max number of connections
		IdleTimeout: Max_Redis_Idle_Time * time.Second,
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

// Ping pings redis
func (self *RedisMgr) Ping() (string, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.String(conn.Do(pingCommand))
}

// Echo echoes the message
func (self *RedisMgr) Echo(message string) (string, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.String(conn.Do(echoCommand, message))
}

// Info returns redis information and statistics
func (self *RedisMgr) Info() (string, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.String(conn.Do(infoCommand))
}

// Scan incrementally iterate over keys
func (self *RedisMgr) Scan(startIndex int64, pattern string) (int64, []string, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	results, err := redis.Values(conn.Do(scanCommand, startIndex, matchOption, pattern))
	if err != nil {
		return 0, nil, err
	}
	return parseScanResults(results)
}

// Append to a key's value
func (self *RedisMgr) Append(key string, value string) (int64, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.Int64(conn.Do("APPEND", key, value))
}

// GetRange to get a key's value's range
func (self *RedisMgr) GetRange(key string, start int, end int) (string, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.String(conn.Do(getRangeCommand, key, start, end))
}

// SetRange to set a key's value's range
func (self *RedisMgr) SetRange(key string, start int, value string) (int64, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.Int64(conn.Do(setRangeCommand, key, start, value))
}

// Expire sets a key's timeout in seconds
func (self *RedisMgr) Expire(key string, timeout int64) (bool, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	count, err := redis.Int64(conn.Do(expireCommand, key, timeout))
	return count > 0, err
}

// Set sets a key/value pair
func (self *RedisMgr) Set(key string, value string) (bool, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return toBool(conn.Do(setCommand, key, value))
}

// SetNx sets a key/value pair if the key does not exist
func (self *RedisMgr) SetNx(key string, value string) (bool, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return toBool(conn.Do(setCommand, key, value, notExistsOption))
}

// SetEx sets a key/value pair with a timeout in seconds
func (self *RedisMgr) SetEx(key string, value string, timeout int64) (bool, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return toBool(conn.Do(setCommand, key, value, expireOption, timeout))
}

// Get retrieves a key's value
func (self *RedisMgr) Get(key string) (string, bool, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return toString(conn.Do(getCommand, key))
}

// Exists checks how many keys exist
func (self *RedisMgr) Exists(keys ...string) (bool, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	interfaces := make([]interface{}, len(keys))
	for i, key := range keys {
		interfaces[i] = key
	}
	count, err := redis.Int64(conn.Do(existsCommand, interfaces...))
	return count > 0, err
}

// Del deletes keys
func (self *RedisMgr) Del(keys ...string) (int64, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	interfaces := make([]interface{}, len(keys))
	for i, key := range keys {
		interfaces[i] = key
	}
	return redis.Int64(conn.Do(delCommand, interfaces...))
}

// Keys retrieves keys that match a pattern
func (self *RedisMgr) Keys(pattern string) ([]string, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.Strings(conn.Do(keysCommand, pattern))
}

// Incr increments the key's value
func (self *RedisMgr) Incr(key string) (int64, error) {
	return self.IncrBy(key, 1)
}

// IncrBy increments the key's value by the increment provided
func (self *RedisMgr) IncrBy(key string, increment int64) (int64, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.Int64(conn.Do(incrByCommand, key, increment))
}

// IncrByFloat increments the key's value by the increment provided
func (self *RedisMgr) IncrByFloat(key string, increment float64) (float64, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.Float64(conn.Do(incrByFloatCommand, key, increment))
}

// Decr decrements the key's value
func (self *RedisMgr) Decr(key string) (int64, error) {
	return self.IncrBy(key, -1)
}

// DecrBy decrements the key's value by the decrement provided
func (self *RedisMgr) DecrBy(key string, decrement int64) (int64, error) {
	return self.IncrBy(key, -decrement)
}

// DecrByFloat decrements the key's value by the decrement provided
func (self *RedisMgr) DecrByFloat(key string, decrement float64) (float64, error) {
	return self.IncrByFloat(key, -decrement)
}

// HScan incrementally iterate over key's fields and values
func (self *RedisMgr) HScan(key string, startIndex int64, pattern string, count int) (int64, []string, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	results, err := redis.Values(conn.Do(hScanCommand, key, startIndex, matchOption, pattern, countOption, count))
	if err != nil {
		return 0, nil, err
	}
	return parseScanResults(results)
}

// HSet sets a key's field/value pair
func (self *RedisMgr) HSet(key string, field string, value string) (bool, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	code, err := redis.Int(conn.Do(hSetCommand, key, field, value))
	return code > 0, err
}

// HMSet sets a key's field/value pair map
func (self *RedisMgr) HMSet(key string, item map[string]interface{}) error {
	conn := self.GetRedisConn()
	defer conn.Close()
	reply, err := conn.Do(hmSetCommand, redis.Args{}.Add(key).AddFlat(item)...)
	if err != nil {
		return err
	}

	if reply != "OK" {
		return fmt.Errorf("reply string is wrong!: %s", reply)
	}

	return nil
}

// HKeys retrieves a hash's keys
func (self *RedisMgr) HKeys(key string) ([]string, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.Strings(conn.Do(hKeysCommand, key))
}

// HExists determine's a key's field's existence
func (self *RedisMgr) HExists(key string, field string) (bool, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.Bool(conn.Do(hExistsCommand, key, field))
}

// HExists determine's a key's field's existence
func (self *RedisMgr) HLen(key string) (int, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.Int(conn.Do(hLenCommand, key))
}

// HGet retrieves a key's field's value
func (self *RedisMgr) HGet(key string, field string) (string, bool, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return toString(conn.Do(hGetCommand, key, field))
}

// HGetAll retrieves the key
func (self *RedisMgr) HGetAll(key string) (map[string]string, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.StringMap(conn.Do(hGetAllCommand, key))
}

// HDel deletes a key's fields
func (self *RedisMgr) HDel(key string, fields ...string) (int64, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	interfaces := make([]interface{}, len(fields)+1)
	interfaces[0] = key
	for i, key := range fields {
		interfaces[i+1] = key
	}
	return redis.Int64(conn.Do(hDelCommand, interfaces...))
}

// HIncr increments the key's field's value
func (self *RedisMgr) HIncr(key string, field string) (int64, error) {
	return self.HIncrBy(key, field, 1)
}

// HIncrBy increments the key's field's value by the increment provided
func (self *RedisMgr) HIncrBy(key string, field string, increment int64) (int64, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.Int64(conn.Do(hIncrByCommand, key, field, increment))
}

// HIncrByFloat increments the key's field's value by the increment provided
func (self *RedisMgr) HIncrByFloat(key string, field string, increment float64) (float64, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.Float64(conn.Do(hIncrByFloatCommand, key, field, increment))
}

// HDecr decrements the key's field's value
func (self *RedisMgr) HDecr(key string, field string) (int64, error) {
	return self.HIncrBy(key, field, -1)
}

// HDecrBy decrements the key's field's value by the decrement provided
func (self *RedisMgr) HDecrBy(key string, field string, decrement int64) (int64, error) {
	return self.HIncrBy(key, field, -decrement)
}

// HDecrByFloat decrements the key's field's value by the decrement provided
func (self *RedisMgr) HDecrByFloat(key string, field string, decrement float64) (float64, error) {
	return self.HIncrByFloat(key, field, -decrement)
}

func toError(reply interface{}, err error) error {
	_, _, e := toString(reply, err)
	return e
}

func toBool(reply interface{}, err error) (bool, error) {
	_, ok, e := toString(reply, err)
	return ok, e
}

func toString(reply interface{}, err error) (string, bool, error) {
	result, e := redis.String(reply, err)
	if e == redis.ErrNil {
		return result, false, nil
	}
	if e != nil {
		return result, false, e
	}
	return result, true, nil
}

func parseScanResults(results []interface{}) (int64, []string, error) {
	if len(results) != 2 {
		return 0, []string{}, nil
	}

	cursorIndex, err := strconv.ParseInt(string(results[0].([]byte)), 10, 64)
	if err != nil {
		return 0, nil, err
	}

	keyInterfaces := results[1].([]interface{})
	keys := make([]string, len(keyInterfaces))
	for index, keyInterface := range keyInterfaces {
		keys[index] = string(keyInterface.([]byte))
	}
	return cursorIndex, keys, nil
}

func (self *RedisMgr) GetTTL(key string) (int, error) {
	conn := self.GetRedisConn()
	defer conn.Close()

	return redis.Int(conn.Do(TTL, key))
}
