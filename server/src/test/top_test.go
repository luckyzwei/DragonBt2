package test

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJsonIgnore(t *testing.T)  {
	type Js_ActTop struct {
		Uid       int64  `json:"uid"`        //! 玩家Id
		Uname     string `json:"uname"`      //! 玩家姓名
		Iconid    int    `json:"iconid"`     //! 玩家头像
		Level     int    `json:"level"`      //! 玩家等级
		Camp      int    `json:"camp"`       //! 阵营
		Vip       int    `json:"vip"`        //! vip等级
		Num       int    `json:"num"`        //! 排行榜数值
		UnionName string `json:"union-name"` //! 军团名字
		LastRank  int    `json:"last-rank"`          //! 原有排名
	}
	top := &Js_ActTop{}
	top.UnionName = "军团名字"
	b, _ := json.Marshal(top)
	fmt.Printf("%s", string(b))
}
