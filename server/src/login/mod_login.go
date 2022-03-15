package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type ModLogin struct {
}

type San_Account struct {
	Uid       int64
	Account   string
	Password  string
	Creator   string
	Channelid string
	Time      int64
}

type JS_SDKData struct {
	Token string `json:"token"`
}

type JS_SDKGame struct {
	Id string `json:"id"`
}

type JS_SDKState struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type JS_SDKHData struct {
	UserId    string `json:"userId"`
	Creator   string `json:"creator"`
	ChannelId string `json:"channelid"`
}

type JS_SDKBody struct {
	Id    string      `json:"id"`
	State JS_SDKState `json:"state"`
	Data  JS_SDKHData `json:"data"`
}

type JS_SDKLogin struct {
	Id   int64      `json:"id"`
	Game JS_SDKGame `json:"game"`
	Data JS_SDKData `json:"data"`
	Sign string     `json:"sign"`
}

type JS_SDKIOSLogin struct {
	AppId     int    `json:"app_id"`
	MemId     int    `json:"mem_id"`
	UserToken string `json:"user_token"`
	Sign      string `json:"sign"`
}

type JS_SDKIOSBody struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

type S2C_Reg struct {
	Cid      string `json:"cid"`
	Uid      int64  `json:"uid"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Creator  string `json:"creator"`
}

func (self *ModLogin) GetName() string {
	return ModName_Login
}

func (self *ModLogin) Init() bool {
	GetServer().RegisterHandler(CMD_Login_Guest, self.LoginGuest)
	GetServer().RegisterHandler(CMD_Login_Sdk, self.LoginSdk)
	GetServer().RegisterHandler(CMD_Login_Ios, self.LoginIos)
	return true
}

func (self *ModLogin) Destory() {
}

func (self *ModLogin) Reg(account string, password string, creator string) (int, *San_Account) {
	if account == "" { //! 游客登录
		b := make([]byte, 48)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			return RET_PARAM_ERROR, nil
		}
		h := md5.New()
		h.Write([]byte(base64.URLEncoding.EncodeToString(b)))
		account = hex.EncodeToString(h.Sum(nil))
	}

	var _account San_Account
	sql := fmt.Sprintf("select * from `san_account` where `account` = '%s'", account)
	GetServer().DBUser.GetOneData(sql, &_account, "", 0)

	if _account.Uid > 0 {
		if password != _account.Password { //! 密码错误
			return RET_LOGIN_ERROR_PASSWORD, nil
		}
	} else { //! 插入
		_account.Account = account
		_account.Password = password
		_account.Creator = creator
		_account.Time = time.Now().Unix()
		_account.Uid = InsertTable("san_account", &_account, 1, false)
		//_account.Uid = GetServer().GetRedisInc("san_account")
		if _account.Uid <= 0 {
			return RET_PARAM_ERROR, nil
		}
	}

	return RET_OK, &_account
}

func (self *ModLogin) SDKReg(token string) (int, *San_Account) {
	log.Println(token)

	str := "token=" + token + GetServer().Con.AppKey
	h := md5.New()
	h.Write([]byte(str))
	md5str := fmt.Sprintf("%x", h.Sum(nil))
	log.Println(md5str)

	var m JS_SDKLogin
	m.Id = time.Now().Unix()
	m.Game.Id = GetServer().Con.GameID
	m.Data.Token = token
	m.Sign = md5str
	body := bytes.NewBuffer(HF_JtoB(&m))
	url := "http://account.flysdk.cn/gs/account.verifyToken?ver=1.0&df=json"
	res, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		log.Println(err)
		return RET_LOGIN_ERROR_SDK, nil
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Println(err)
		return RET_LOGIN_ERROR_SDK, nil
	}
	log.Println("result=", string(result))

	var ret JS_SDKBody
	err = json.Unmarshal(result, &ret)
	if err != nil {
		log.Println(err)
		return RET_LOGIN_ERROR_SDK, nil
	}

	switch ret.State.Code {
	case 4000000:
		return RET_PARAM_ERROR, nil
	case 4000001:
		return RET_PARAM_ERROR, nil
	case 5000000:
		return RET_LOGIN_ERROR_NETWORKBUSY_SDK, nil
	case 5000003:
		return RET_LOGIN_ERROR_SYSTEMBUSY_SDK, nil
	case 4001001:
		return RET_LOGIN_ERROR_TOKEN_INVALID_SDK, nil
	case 4001003:
		return RET_LOGIN_ERROR_TOKEN_TIMEOUT_SDK, nil
	}

	return self.Reg(ret.Data.UserId, "CYCYCY", ret.Data.Creator)
}

func (self *ModLogin) IOSReg(token string, appid string, memid string) (int, *San_Account) {

	str := "app_id=" + appid + "&mem_id=" + memid + "&user_token=" + token + "&app_key=" + GetServer().Con.GetAppKeyByAppId(appid)
	h := md5.New()
	h.Write([]byte(str))
	md5str := fmt.Sprintf("%x", h.Sum(nil))
	log.Println("sign = ", md5str)

	str1 := "app_id=" + appid + "&mem_id=" + memid + "&app_key=" + GetServer().Con.GetAppKeyByAppId(appid)
	h1 := md5.New()
	h1.Write([]byte(str1))
	md5str1 := fmt.Sprintf("%x", h1.Sum(nil))
	log.Println("sign1 = ", md5str1)

	//var m JS_SDKIOSLogin
	//m.AppId = HF_Atoi(appid)
	//m.MemId = memid
	//m.UserToken = token
	//m.Sign = md5str
	//body := bytes.NewBuffer(HF_JtoB(&m))
	check := "app_id=" + appid + "&mem_id=" + memid + "&user_token=" + token + "&sign=" + md5str

	//body := bytes.NewBuffer([]byte(""))
	log.Println("check = ", check)
	//url := "https://aliapi.1tsdk.com/api/v7/cp/user/check?" + check
	//res, err := http.Get(url)
	//if err != nil {
	//	log.Println(err)
	//	self.SendErrInfo("err", "sdk错误")
	//	return nil
	//}

	url := "https://aliapi.1tsdk.com/api/v7/cp/user/check?" + check
	body := bytes.NewBuffer([]byte(""))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(6 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*1)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Println(err)
		return RET_LOGIN_ERROR_SDK, nil
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return RET_LOGIN_ERROR_SDK, nil
	}

	if resp.StatusCode != 200 {
		log.Println(err)
		return RET_LOGIN_ERROR_SDK, nil
	}

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return RET_LOGIN_ERROR_SDK, nil
	}
	log.Println("result=", string(result))

	var ret JS_SDKIOSBody
	err = json.Unmarshal(result, &ret)
	if err != nil {
		log.Println(err)
		return RET_LOGIN_ERROR_SDK, nil
	}

	switch ret.Status {
	case "0":
		return RET_PARAM_ERROR, nil
	case "10":
		//self.SendErrInfo("err", "服务器内部错误")
		return RET_PARAM_ERROR, nil
	case "11":
		return RET_LOGIN_ERROR_IOS_APPID_SDK, nil
		//self.SendErrInfo("err", "app_id错误")
		//return nil
	case "12":
		return RET_LOGIN_ERROR_IOS_SIGNATURE_SDK, nil
		//self.SendErrInfo("err", "签名错误")
		//return nil
	case "13":
		return RET_LOGIN_ERROR_IOS_TOKEN_SDK, nil

		//self.SendErrInfo("err", "user_token错误")
		//return nil
	case "14":
		return RET_LOGIN_ERROR_IOS_TOKEN_TIMEOUT_SDK, nil
		//self.SendErrInfo("err", "user_token超时，表示用户登录授权已超时，需引导用户重新登录，并更新接口访问令牌。")
		//return nil
	}

	account := md5str
	account1 := md5str1
	password := "CYCYCY"

	if account == "" { //! 游客登录
		b := make([]byte, 48)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			return RET_PARAM_ERROR, nil
		}
		h := md5.New()
		h.Write([]byte(base64.URLEncoding.EncodeToString(b)))
		account = hex.EncodeToString(h.Sum(nil))
	}

	var _account San_Account
	sql := fmt.Sprintf("select * from `san_account` where `account` = '%s' or `account` = '%s'", account, account1)
	GetServer().DBUser.GetOneData(sql, &_account, "", 0)

	if _account.Uid > 0 {
		if password != _account.Password { //! 密码错误
			return RET_LOGIN_ERROR_PASSWORD, nil
		}

		_account.Account = account1
		updateQuery := fmt.Sprintf("update san_account set Account = '%s' where uid=%d limit 1", account1, _account.Uid)
		GetServer().SqlSet(updateQuery)
	} else { //! 插入
		_account.Account = account1
		_account.Password = password
		_account.Creator = "ios"
		_account.Time = time.Now().Unix()
		_account.Uid = InsertTable("san_account", &_account, 1, false)
		//_account.Uid = GetServer().GetRedisInc("san_account")
		if _account.Uid <= 0 {
			return RET_PARAM_ERROR, nil
		}
	}

	return RET_OK, &_account
}

func (self *ModLogin) LoginGuest(r *http.Request) interface{} {

	Account := r.FormValue("Account")
	Password := r.FormValue("Password")

	_, _account := self.Reg(Account, Password, "guest")

	var response S2C_Reg
	if _account != nil {
		response.Cid = "reg"
		response.Uid = _account.Uid
		response.Account = _account.Account
		response.Password = _account.Password
		response.Creator = _account.Creator
	} else {

	}
	return &response
}

func (self *ModLogin) LoginSdk(r *http.Request) interface{} {

	token := r.FormValue("password")
	_, _account := self.SDKReg(token)

	var response S2C_Reg
	if _account != nil {
		response.Cid = "reg"
		response.Uid = _account.Uid
		response.Account = _account.Account
		response.Password = _account.Password
		response.Creator = _account.Creator
	} else {

	}
	return &response
}

func (self *ModLogin) LoginIos(r *http.Request) interface{} {

	UserToken := r.FormValue("UserToken")
	AppId := r.FormValue("AppId")
	MemId := r.FormValue("MemId")

	_, _account := self.IOSReg(UserToken, AppId, MemId)

	var response S2C_Reg
	if _account != nil {
		response.Cid = "reg"
		response.Uid = _account.Uid
		response.Account = _account.Account
		response.Password = _account.Password
		response.Creator = _account.Creator
	} else {

	}
	return &response
}
