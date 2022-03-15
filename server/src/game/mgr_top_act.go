package game

// 排行榜分国家
type Js_ActTop struct {
	Uid       int64  `json:"uid"`        //! 玩家Id
	Uname     string `json:"uname"`      //! 玩家姓名
	Iconid    int    `json:"iconid"`     //! 玩家头像
	Portrait  int    `json:"portrait"`   // 边框  20190412 by zy
	Level     int    `json:"level"`      //! 玩家等级
	Camp      int    `json:"camp"`       //! 阵营
	Fight     int64  `json:"fight"`      //! 战力
	Vip       int    `json:"vip"`        //! vip等级
	Num       int64  `json:"num"`        //! 排行榜数值
	UnionName string `json:"union_name"` //! 军团名字
	LastRank  int    `json:"-"`          //! 原有排名
	StartTime int64  `json:"starttime"`  //! 时间戳
}

// 排行榜分国家
type Js_ActTopLoad struct {
	Uid       int64  `json:"uid"`        //! 玩家Id
	Uname     string `json:"uname"`      //! 玩家姓名
	Iconid    int    `json:"iconid"`     //! 玩家头像
	Portrait  int    `json:"portrait"`   // 边框  20190412 by zy
	Level     int    `json:"level"`      //! 玩家等级
	Camp      int    `json:"camp"`       //! 阵营
	Fight     int64  `json:"fight"`      //! 战力
	Vip       int    `json:"vip"`        //! vip等级
	Nums      string `json:"num"`        //! 排行榜数值
	UnionName string `json:"union_name"` //! 军团名字
	LastRank  int    `json:"-"`          //! 原有排名
}

// 活动
type lstJsActTop []*Js_ActTop

func (s lstJsActTop) Len() int      { return len(s) }
func (s lstJsActTop) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s lstJsActTop) Less(i, j int) bool {
	if s[i].Num > s[j].Num { // 由大到小
		return true
	}

	if s[i].Num < s[j].Num {
		return false
	}

	if s[i].StartTime < s[j].StartTime {
		return true
	}

	if s[i].StartTime > s[j].StartTime {
		return false
	}

	if s[i].LastRank < s[j].LastRank { // 由大到小
		return true
	}

	if s[i].LastRank > s[j].LastRank {
		return false
	}

	if s[i].Uid < s[j].Uid { // 由小到大
		return true
	}

	if s[i].Uid > s[j].Uid {
		return false
	}

	return false
}
