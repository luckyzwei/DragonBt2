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
	//"time"
)

type JS_SDKLogin_MYX struct {
	Service string `json:"service"`
	AppID   int    `json:"appid"`
	Data    string `json:"data"`
	Sign    string `json:"sign"`
}

type JS_SDKLoginData_MYX struct {
	Sid      string `json:"sid"`
	Username string `json:"username"`
}

type JS_SDKData_MYX struct {
	Uid string `json:"uid"`
}

type JS_SDKBody_MYX struct {
	Id    string         `json:"id"`
	State JS_SDKState    `json:"state"`
	Data  JS_SDKData_MYX `json:"data"`
}

// 安卓登录-锦游登录
func (self *Session) SDKReg_MYX(token string, serverid int, username string, ctrl string) *San_Account {
	log.Println(token)

	appkey := "0a7357b0c3d8023c920150aed32b1754"
	appid := 1000037

	//str := "token=" + token + GetServer().Con.AppKey
	str := fmt.Sprintf("%dsdk.game.checkentersid=%s&username=%s%s", appid, token, username, appkey)
	l := url.QueryEscape(str)
	h := md5.New()
	h.Write([]byte(l))
	md5str := fmt.Sprintf("%x", h.Sum(nil))
	log.Println(md5str, str, l)

	var data JS_SDKLoginData_MYX
	data.Username = username
	data.Sid = token

	var m JS_SDKLogin_MYX
	m.AppID = appid
	m.Service = "sdk.game.checkenter"
	m.Data = HF_JtoA(data)
	//m.Id = TimeServer().Unix()
	//m.Game.Id = GetServer().Con.GameID
	//m.Data.Token = token
	m.Sign = md5str

	postvalue := url.Values{
		"service": {"sdk.game.checkenter"},
		"appid":   {fmt.Sprintf("%d", appid)},
		"data":    {HF_JtoA(data)},
		"sign":    {md5str},
	}

	//body := bytes.NewBuffer(HF_JtoB(&m))
	url := "http://mp.gzjykj.com/index.php"
	res, err := http.PostForm(url, postvalue)
	if err != nil {
		log.Println(err)
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

	var ret JS_SDKBody_MYX
	err = json.Unmarshal(result, &ret)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_APPLICATION_ERROR"))
		return nil
	}

	switch ret.State.Code {
	case 0:
		self.SendErrInfo("err", ret.State.Msg)
		return nil
	case 4000000:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_REQUEST_PARAMETER_ERROR"))
		return nil
	case 4000001:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_BUSINESS_PARAMETER_ERROR"))
		return nil
	case 5000000:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_THE_NETWORK_IS_BUSY_PLEASE"))
		return nil
	case 5000003:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_THE_SYSTEM_IS_BUSY_PLEASE"))
		return nil
	case 4001001:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_INVALID_SACRED_MARK"))
		return nil
	case 4001003:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_MYX_THE_GREAT_SAINT_LOGO_HAS"))
		return nil
	}

	h1 := md5.New()
	LogDebug("New MYX User:", ret.Data.Uid+appkey)
	h1.Write([]byte(ret.Data.Uid + appkey))
	account := fmt.Sprintf("%x", h1.Sum(nil))
	//account := ret.Data.Uid
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
		_account.Creator = "myx"
		_account.Channelid = fmt.Sprintf("%d", appid)
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