package test

import (
	"fmt"
	"game"
	"log"
	"testing"
	"time"
)

func InitRedis()  {
	redisClient := game.NewPool("127.0.0.1:6379", 13, "123456")
	if redisClient == nil {
		log.Println("redis err")
		return
	}
	game.GetServer().Redis = redisClient
}

// 测试redis插入HMSet
func TestHashFiledSet(t *testing.T)  {
	// 连接redis
	//! 连接redis
	InitRedis()

	var item = make(map[string]interface{})
	item["10001"] = "fight_result_1"
	item["10002"] = "fight_result_1"
	game.GetRedisMgr().HMSet("fight_record", item)
	items, err := game.GetRedisMgr().HGetAll("fight_record")
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Println(items)
	// 设置键10秒过期
	//flag, err := game.GetRedisMgr().Expire("fight_record", 10)
	//if err != nil {
	//	log.Println(err.Error())
	//	return
	//}
	//
	//if flag {
	//	fmt.Println("flag is ok")
	//}
	time.Sleep(time.Second * 2)
	ttl, err := game.GetRedisMgr().GetTTL("fight_record")
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Println("ttl:", ttl)
}

// 测试redis插入HGet
func TestHashFiledGet(t *testing.T)  {
	// 连接redis
	//! 连接redis
	InitRedis()
	value,flag,err := game.GetRedisMgr().HGet("fight_record", "10001")
	fmt.Println("value:", value, ",flag:", flag, ",err:", err)
}
