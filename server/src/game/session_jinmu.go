package game

import (
	//"bytes"
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
	//"net/url"
	//"time"
)

type JS_SDKLogin_Jinmu struct {
	Service string `json:"service"`
	AppID   int    `json:"appid"`
	Data    string `json:"data"`
	Sign    string `json:"sign"`
}

type JS_SDKLoginData_Jinmu struct {
	Sid      string `json:"sid"`
	Username string `json:"username"`
}

type JS_SDKData_Jinmu struct {
	UserID    int    `json:"userID"`
	Username  string `json:"username"`
	ChannelId int    `json:"channelID"`
}

type JS_SDKBody_Jinmu struct {
	State int `json:"state"`
	//State JS_SDKState    `json:"state"`
	Data    JS_SDKData_Jinmu `json:"data"`
	Message string           `json:"message"`
}

// 安卓登录-金木登录
func (self *Session) SDKReg_Jinmu(token string, serverid int, username string, ctrl string) *San_Account {
	log.Println(token)

	//appid := "12"
	appkey := "e0b49324c7be91f08ea8287cb3e83cfa"
	//appsecret := "8115a4a5d830bbe641bbd54b180fc28c"

	str := "userID=" + username + "token=" + token + appkey
	h := md5.New()
	h.Write([]byte(str))
	md5str := fmt.Sprintf("%x", h.Sum(nil))
	log.Println(str, md5str)

	//var m JS_SDKLogin
	//m.Id = TimeServer().Unix()
	//m.Game.Id = GetServer().Con.GameID
	//m.Data.Token = token
	//m.Sign = md5str
	//body := bytes.NewBuffer(HF_JtoB(&m))
	url := "http://119.23.161.29:8082/u8server/user/verifyAccount?userID=" + username + "&token=" + token + "&sign=" + md5str
	LogDebug("req url:", url)
	res, err := http.Get(url)
	if err != nil {
		log.Println(err, url)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_APPLICATION_ERROR"))
		return nil
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_APPLICATION_ERROR"))
		return nil
	}
	log.Println("result=", string(result))

	var ret JS_SDKBody_Jinmu
	err = json.Unmarshal(result, &ret)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_APPLICATION_ERROR"))
		return nil
	}

	switch ret.State {
	case 0:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_REQUEST_PARAMETER_ERROR"))
		return nil
	case 1:
		LogDebug("金木登入成功：", result)
	case 12:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_REQUEST_PARAMETER_ERROR"))
		return nil
	default:
		self.SendErrInfo("err", ret.Message)
		return nil
	}

	account := ret.Data.Username
	password := "CYCYCY"

	if account == "" { //! 游客登录
		b := make([]byte, 48)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
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
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_PASSWORD"))
			return nil
		}
	} else { //! 插入
		_account.Account = account
		_account.Password = password
		_account.Creator = "jinmu-android"
		_account.Channelid = fmt.Sprintf("%d", ret.Data.ChannelId)
		_account.ServerId = serverid
		_account.Time = TimeServer().Unix()
		_account.Uid = InsertTable("san_account", &_account, 1, false)
		//_account.Uid = GetServer().GetRedisInc("san_account")
		if _account.Uid <= 0 {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return nil
		}
		//value := HF_JtoB(&_account)
		//c := GetServer().GetRedisConn()
		//defer c.Close()
		//_, err := c.Do("MSET", fmt.Sprintf("%s_%d", "san_account", _account.Uid), value,
		//	fmt.Sprintf("%s_%s", "san_account", _account.Account), value)
		//if err != nil {
		//	LogError("redis set err:", "san_account", ",", string(value))
		//	self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
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

//----------------------------------------------------------------
//-数果接入，同类型，归属到一个文件

type JS_SDKLogin_Shuguo struct {
	Service string `json:"service"`
	AppID   int    `json:"appid"`
	Data    string `json:"data"`
	Sign    string `json:"sign"`
}

type JS_SDKLoginData_Shuguo struct {
	Sid      string `json:"sid"`
	Username string `json:"username"`
}

type JS_SDKData_Shuguo struct {
	UserID    int    `json:"userID"`
	Username  string `json:"username"`
	ChannelId int    `json:"channelID"`
}

type JS_SDKBody_Shuguo struct {
	State int `json:"state"`
	//State JS_SDKState    `json:"state"`
	Data JS_SDKData_Shuguo `json:"data"`
}

// 安卓登录-金木登录
func (self *Session) SDKReg_Shuguo(token string, serverid int, username string, ctrl string) *San_Account {
	log.Println(token)

	//appid := "12"
	//appkey := "55a63944ce55cab8f13e22ecd3e4088c"
	appsecret := "58b3ebcf89ea57639760dfc6a23d03a6"

	str := "userID=" + username + "token=" + token + appsecret
	h := md5.New()
	h.Write([]byte(str))
	md5str := fmt.Sprintf("%x", h.Sum(nil))
	log.Println(str, md5str)

	//var m JS_SDKLogin
	//m.Id = TimeServer().Unix()
	//m.Game.Id = GetServer().Con.GameID
	//m.Data.Token = token
	//m.Sign = md5str
	//body := bytes.NewBuffer(HF_JtoB(&m))
	url := "http://www.hnshuguo.com:8080/u8server/user/verifyAccount?userID=" + username + "&token=" + token + "&sign=" + md5str
	LogDebug("req url:", url)
	res, err := http.Get(url)
	if err != nil {
		log.Println(err, url)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_APPLICATION_ERROR"))
		return nil
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_APPLICATION_ERROR"))
		return nil
	}
	log.Println("result=", string(result))

	var ret JS_SDKBody_Jinmu
	err = json.Unmarshal(result, &ret)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_APPLICATION_ERROR"))
		return nil
	}

	switch ret.State {
	case 0:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_REQUEST_PARAMETER_ERROR"))
		return nil
	case 1:
		LogDebug("数果登入成功：", result)
	case 12:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_REQUEST_PARAMETER_ERROR"))
		return nil
	default:
		self.SendErrInfo("err", ret.Message)
		return nil
	}

	account := ret.Data.Username
	password := "CYCYCY"

	if account == "" { //! 游客登录
		b := make([]byte, 48)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
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
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_PASSWORD"))
			return nil
		}
	} else { //! 插入
		_account.Account = account
		_account.Password = password

		_account.Creator = ctrl
		_account.Channelid = fmt.Sprintf("%d", ret.Data.ChannelId)
		_account.ServerId = serverid
		_account.Time = TimeServer().Unix()
		_account.Uid = InsertTable("san_account", &_account, 1, false)
		//_account.Uid = GetServer().GetRedisInc("san_account")
		if _account.Uid <= 0 {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
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
		//	self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
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