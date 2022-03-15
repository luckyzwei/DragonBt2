package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"encoding/json"
	"runtime/debug"
	"net/http"
	"strings"
)

func HF_JtoB(v interface{}) []byte {
	s, err := json.Marshal(v)
	if err != nil {
		LogError("HF_JtoB err:", string(debug.Stack()))
	}
	return s
}


//! 得到ip
func HF_GetHttpIP(req *http.Request) string {
	ip := req.Header.Get("Remote_addr")
	if ip == "" {
		ip = req.RemoteAddr
	}
	return strings.Split(ip, ":")[0]
}

// 重写生成连接池方法
func NewPool(ip string, db int, auth string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   12000, // max number of connections
		IdleTimeout: 180 * time.Second,
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
