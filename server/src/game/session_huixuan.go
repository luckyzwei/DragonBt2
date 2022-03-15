package game

import (
	//"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	//"time"
)

type JS_SDKLogin_HuiXuan struct {
	Service string `json:"service"`
	AppID   int    `json:"appid"`
	Data    string `json:"data"`
	Sign    string `json:"sign"`
}

type JS_SDKLoginData_HuiXuan struct {
	Sid      string `json:"sid"`
	Username string `json:"username"`
}

type JS_SDKData_HuiXuan struct {
	Uid string `json:"uid"`
}

type JS_SDKBody_HuiXuan struct {
	Status    int    `json:"status"`
	Errror    string `json:"error"`
	ErrorDesc string `json:"error_description"`
	Data      string `json:"data"`
}

// 安卓登录-慧选登录
func (self *Session) SDKReg_HuiXuan(token string, serverid int, username string, ctrl string) *San_Account {
	log.Println(token)

	//! 传递参数
	arrParam := strings.Split(token, "||")
	if (len(arrParam) < 5) {
		self.SendErrInfo("err", "Param Error.")
		return nil
	}

	//agentName := arrParam[0]
	clientId := arrParam[2]
	playerToken := arrParam[3]
	appid := arrParam[1]

	//appkey := "1ca0b94650f930bbbd9cf75e94610f13"
	//appid := 132
	appsecret := "8dca8de93d524b6d789260a60c375b12"
	if ctrl == "sdk_huixuan_ios" {
		appsecret = "54aca9eec11d7132335d269180e5d1da"
	} else if ctrl == "sdk_huixuan_new_ios" {
		appsecret = "65f060588358bd19434889c585e3a39c"
	}
	timeNow := TimeServer().Unix()
	//str := "token=" + token + GetServer().Con.AppKey
	str := fmt.Sprintf("client_id=%s&open_id=%s&session_token=%s&timestamp=%d",
		appid, clientId, playerToken, timeNow)
	l := str //url.QueryEscape(str)
	l = l + "|" + appsecret
	h := md5.New()
	h.Write([]byte(l))
	md5str := fmt.Sprintf("%x", h.Sum(nil))
	log.Println(md5str, str, l)

	//var data JS_SDKLoginData_MYX
	//data.Username = username
	//data.Sid = token

	//var m JS_SDKLogin_MYX
	//m.AppID = appid
	//m.Service = "sdk.game.checkenter"
	//m.Data = HF_JtoA(data)
	////m.Id = TimeServer().Unix()
	////m.Game.Id = GetServer().Con.GameID
	////m.Data.Token = token
	//m.Sign = md5str

	postvalue := url.Values{
		"client_id":     {appid},
		"open_id":       {clientId},
		"session_token": {playerToken},
		"sign":          {md5str},
		"from":          {arrParam[4]},
		"timestamp":     {fmt.Sprintf("%d", timeNow)},
	}

	LogInfo(fmt.Sprintf("%v", postvalue))

	//body := bytes.NewBuffer(HF_JtoB(&m))
	url := "http://open.douyouzhiyu.com/v2/users/check_token"
	res, err := http.PostForm(url, postvalue)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", "sdk错误")
		return nil
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", "sdk错误")
		return nil
	}
	log.Println("result=", string(result))

	var ret JS_SDKBody_HuiXuan
	err = json.Unmarshal(result, &ret)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", "sdk错误")
		return nil
	}

	switch ret.Status {
	case 0:
		self.SendErrInfo("err", ret.Errror)
		return nil
	case 200:
		LogInfo("HuiXuan Login OK: ", clientId, username, token)
	case 4000000:
		self.SendErrInfo("err", "请求参数错误")
		return nil
	case 4000001:
		self.SendErrInfo("err", "业务参数错误")
		return nil
	case 5000000:
		self.SendErrInfo("err", "网络繁忙，请稍后重试")
		return nil
	case 5000003:
		self.SendErrInfo("err", "系统繁忙，请稍后重试")
		return nil
	case 4001001:
		self.SendErrInfo("err", "无效的大圣token")
		return nil
	case 4001003:
		self.SendErrInfo("err", "大圣token已过期")
		return nil
	default:
		self.SendErrInfo("err", "请求参数错误")
		return nil
	}

	h1 := md5.New()
	LogDebug("New HuiXuan User:", username)
	h1.Write([]byte(username))
	account := fmt.Sprintf("%x", h1.Sum(nil))
	//account := ret.Data.Uid
	password := "CYCYCY"

	if account == "" { //! 游客登录
		b := make([]byte, 48)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			self.SendErrInfo("err", "ERROR")
			return nil
		}
		h := md5.New()
		h.Write([]byte(base64.URLEncoding.EncodeToString(b)))
		account = hex.EncodeToString(h.Sum(nil))
	}

	var _account San_Account
	sql := fmt.Sprintf("select * from `san_account` where `account` = '%s' and serverid =%d", account, serverid)
	GetServer().DBUser.GetOneData(sql, &_account, "", 0)
	//c := GetServer().GetRedisConn()
	//defer c.Close()

	//value, _ := redis.Bytes(c.Do("GET", fmt.Sprintf("%s_%s", "san_account", account)))
	//json.Unmarshal(value, &_account)

	if _account.Uid > 0 {
		if password != _account.Password { //! 密码错误
			self.SendErrInfo("err", "ERROR")
			return nil
		}
	} else { //! 插入
		_account.Account = account
		_account.Password = password
		_account.Creator = ctrl
		_account.Channelid = appid
		_account.ServerId = serverid
		_account.Time = TimeServer().Unix()
		_account.Uid = InsertTable("san_account", &_account, 1, false)
		//_account.Uid = GetServer().GetRedisInc("san_account")
		if _account.Uid <= 0 {
			self.SendErrInfo("err", "ERROR")
			return nil
		}

		sql := fmt.Sprintf("select * from `san_account` where `account` = '%s' and serverid =%d", account, serverid)
		GetServer().DBUser.GetOneData(sql, &_account, "", 0)
		//value := HF_JtoB(&_account)
		//c := GetServer().GetRedisConn()
		//defer c.Close()
		//_, err := c.Do("MSET", fmt.Sprintf("%s_%d", "san_account", _account.Uid), value,
		//	fmt.Sprintf("%s_%s", "san_account", _account.Account), value)
		//if err != nil {
		//	LogError("redis set err:", "san_account", ",", string(value))
		//	self.SendErrInfo("err", STR_ERR)
		//	return nil
		//}

		var msg S2C_Reg
		msg.Cid = "reg"
		msg.Uid = _account.Uid
		msg.Account = _account.Account
		msg.Password = _account.Password
		msg.Creator = _account.Creator
		self.SendMsg("1", HF_JtoB(&msg))
	}

	return &_account
}

// 安卓登录-慧选登录
func (self *Session) SDKReg_HuiXuan_IOS(token string, serverid int, username string, ctrl string) *San_Account {
	log.Println(token)

	//! 传递参数
	arrParam := strings.Split(token, "||")
	if (len(arrParam) < 4) {
		self.SendErrInfo("err", "Param Error.")
		return nil
	}

	//agentName := arrParam[0]
	clientId := arrParam[2]
	playerToken := arrParam[3]
	appid := arrParam[1]

	//appkey := "1ca0b94650f930bbbd9cf75e94610f13"
	//appid := 132
	appsecret := "8dca8de93d524b6d789260a60c375b12"
	if ctrl == "sdk_huixuan_ios" {
		appsecret = "54aca9eec11d7132335d269180e5d1da"
	}
	timeNow := TimeServer().Unix()
	//str := "token=" + token + GetServer().Con.AppKey
	str := fmt.Sprintf("client_id=%s&open_id=%s&session_token=%s&timestamp=%d",
		appid, clientId, playerToken, timeNow)
	l := str //url.QueryEscape(str)
	l = l + "|" + appsecret
	h := md5.New()
	h.Write([]byte(l))
	md5str := fmt.Sprintf("%x", h.Sum(nil))
	log.Println(md5str, str, l)

	//var data JS_SDKLoginData_MYX
	//data.Username = username
	//data.Sid = token

	//var m JS_SDKLogin_MYX
	//m.AppID = appid
	//m.Service = "sdk.game.checkenter"
	//m.Data = HF_JtoA(data)
	////m.Id = TimeServer().Unix()
	////m.Game.Id = GetServer().Con.GameID
	////m.Data.Token = token
	//m.Sign = md5str

	postvalue := url.Values{
		"client_id":     {appid},
		"open_id":       {clientId},
		"session_token": {playerToken},
		"sign":          {md5str},
		"timestamp":     {fmt.Sprintf("%d", timeNow)},
	}

	LogInfo(fmt.Sprintf("%v", postvalue))

	//body := bytes.NewBuffer(HF_JtoB(&m))
	url := "http://open.douyouzhiyu.com/v2/users/check_token"
	res, err := http.PostForm(url, postvalue)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", "sdk错误")
		return nil
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", "sdk错误")
		return nil
	}
	log.Println("result=", string(result))

	var ret JS_SDKBody_HuiXuan
	err = json.Unmarshal(result, &ret)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", "sdk错误")
		return nil
	}

	switch ret.Status {
	case 0:
		self.SendErrInfo("err", ret.Errror)
		return nil
	case 200:
		LogInfo("HuiXuan Login OK: ", clientId, username, token)
	case 4000000:
		self.SendErrInfo("err", "请求参数错误")
		return nil
	case 4000001:
		self.SendErrInfo("err", "业务参数错误")
		return nil
	case 5000000:
		self.SendErrInfo("err", "网络繁忙，请稍后重试")
		return nil
	case 5000003:
		self.SendErrInfo("err", "系统繁忙，请稍后重试")
		return nil
	case 4001001:
		self.SendErrInfo("err", "无效的大圣token")
		return nil
	case 4001003:
		self.SendErrInfo("err", "大圣token已过期")
		return nil
	default:
		self.SendErrInfo("err", "请求参数错误")
		return nil
	}

	h1 := md5.New()
	LogDebug("New HuiXuan User:", username)
	h1.Write([]byte(username))
	account := fmt.Sprintf("%x", h1.Sum(nil))
	//account := ret.Data.Uid
	password := "CYCYCY"

	if account == "" { //! 游客登录
		b := make([]byte, 48)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			self.SendErrInfo("err", "ERROR")
			return nil
		}
		h := md5.New()
		h.Write([]byte(base64.URLEncoding.EncodeToString(b)))
		account = hex.EncodeToString(h.Sum(nil))
	}

	var _account San_Account
	sql := fmt.Sprintf("select * from `san_account` where `account` = '%s' and serverid =%d", account, serverid)
	GetServer().DBUser.GetOneData(sql, &_account, "", 0)
	//c := GetServer().GetRedisConn()
	//defer c.Close()

	//value, _ := redis.Bytes(c.Do("GET", fmt.Sprintf("%s_%s", "san_account", account)))
	//json.Unmarshal(value, &_account)

	if _account.Uid > 0 {
		if password != _account.Password { //! 密码错误
			self.SendErrInfo("err", "ERROR")
			return nil
		}
	} else { //! 插入
		_account.Account = account
		_account.Password = password
		_account.Creator = ctrl
		_account.Channelid = fmt.Sprintf("%d", appid)
		_account.ServerId = serverid
		_account.Time = TimeServer().Unix()
		_account.Uid = InsertTable("san_account", &_account, 1, false)
		//_account.Uid = GetServer().GetRedisInc("san_account")
		if _account.Uid <= 0 {
			self.SendErrInfo("err", "ERROR")
			return nil
		}

		sql := fmt.Sprintf("select * from `san_account` where `account` = '%s' and serverid =%d", account, serverid)
		GetServer().DBUser.GetOneData(sql, &_account, "", 0)
		//value := HF_JtoB(&_account)
		//c := GetServer().GetRedisConn()
		//defer c.Close()
		//_, err := c.Do("MSET", fmt.Sprintf("%s_%d", "san_account", _account.Uid), value,
		//	fmt.Sprintf("%s_%s", "san_account", _account.Account), value)
		//if err != nil {
		//	LogError("redis set err:", "san_account", ",", string(value))
		//	self.SendErrInfo("err", STR_ERR)
		//	return nil
		//}

		var msg S2C_Reg
		msg.Cid = "reg"
		msg.Uid = _account.Uid
		msg.Account = _account.Account
		msg.Password = _account.Password
		msg.Creator = _account.Creator
		self.SendMsg("1", HF_JtoB(&msg))
	}

	return &_account
}
