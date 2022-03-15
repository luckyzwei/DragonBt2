package game

import (
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

type JS_SDKLogin_YunKe struct {
	ChannelId   int    `json:"channel_id"`
	GameId      int    `json:"game_id"`
	PlayerId    int    `json:"player_id"`
	PlayerToken string `json:"player_token"`
	Sign        string `json:"sign"`
}

type JS_SDKLoginData_YunKe struct {
	Sid      string `json:"sid"`
	Username string `json:"username"`
}

type JS_SDKData_YunKe struct {
	Uid string `json:"uid"`
}

type JS_SDKBody_YunKe struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}


type JS_SDKBody_YunKe_IOS struct {
	Status string `json:"code"`
	Msg    string `json:"msg"`
}

//func u2s(form string) (to string, err error) {
//	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
//	if err != nil {
//		return
//	}
//	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
//		binary.Read(br, binary.BigEndian, &r)
//		to += string(r)
//	}
//	return
//}

// 安卓登录-云客登录
func (self *Session) SDKReg_YunKe(token string, serverid int, username string, ctrl string) *San_Account {
	log.Println(token)
	//channelId := 1236
	channelId := 1236
	gameId := 258
	//gameKey := "d25a3761568c2b9d6bcf370f4d3f27f4"
	gameKey := "d25a3761568c2b9d6bcf370f4d3f27f4"

	str := fmt.Sprintf("channel_id=%d&game_id=%d&player_id=%d&game_key=%s&player_token=%s",
		channelId, gameId, HF_Atoi(username), gameKey, token)
	h := md5.New()
	h.Write([]byte(str))
	md5str := fmt.Sprintf("%x", h.Sum(nil))
	log.Println("yunke", md5str, str)

	var data JS_SDKLoginData_YunKe
	data.Username = username
	data.Sid = token

	var m JS_SDKLogin_YunKe
	m.GameId = gameId
	m.PlayerId = HF_Atoi(username)
	m.ChannelId = channelId
	m.PlayerToken = token
	m.Sign = md5str

	postvalue := url.Values{
		"channel_id":   {fmt.Sprintf("%d", channelId)},
		"game_id":      {fmt.Sprintf("%d", gameId)},
		"player_id":    {fmt.Sprintf("%d", HF_Atoi(username))},
		"player_token": {token},
		"sign":         {md5str},
	}

	url := "http://usdk.api.kokoyou.com/cp/cp/check"
	res, err := http.PostForm(url, postvalue)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_YUNKE_IDENTIFICATION_ERROR"))
		return nil
	}

	result, err := ioutil.ReadAll(res.Body)
	if res.Body != nil {
		defer res.Body.Close()
	}
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_YUNKE_IDENTIFICATION_ERROR"))
		return nil
	}
	log.Println("result=", string(result))

	var ret JS_SDKBody_YunKe
	err = json.Unmarshal(result, &ret)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_YUNKE_IDENTIFICATION_ERROR"))
		return nil
	}

	switch ret.Status {
	case "1":
		LogInfo("登录云客SDK成功:", username, token)
	default:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_YUNKE_REQUEST_PARAMETER_ERROR"))
		return nil
	}

	h1 := md5.New()
	LogDebug("New YunKe User:", username, token)
	h1.Write([]byte(username))
	account := fmt.Sprintf("%x", h1.Sum(nil))
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
	if _account.Uid > 0 {
		if password != _account.Password { //! 密码错误
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_PASSWORD"))
			return nil
		}
	} else { //! 插入
		_account.Account = account
		_account.Password = password
		_account.UserId = username
		_account.Creator = ctrl
		_account.Channelid = fmt.Sprintf("%d", channelId)
		_account.ServerId = serverid
		_account.Time = TimeServer().Unix()
		_account.Uid = InsertTable("san_account", &_account, 1, false)
		if _account.Uid <= 0 {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return nil
		}

		sql := fmt.Sprintf("select * from `san_account` where `account` = '%s' and serverid =%d", account, serverid)
		GetServer().DBUser.GetOneData(sql, &_account, "", 0)

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

// IOS 登录-云客登录
func (self *Session) SDKReg_YunKe_IOS(token string, serverid int, username string, ctrl string) *San_Account {
	log.Println(token)
	channelId := 1236
	gameId := 80099
	gameKey := "oanelm4vx4vuiqh1wxwco5v110quvdc4"

	str := fmt.Sprintf("gameid=%d&sessionid=%s&uid=%s", gameId, username, token) + gameKey
	h := md5.New()
	h.Write([]byte(str))
	md5str := fmt.Sprintf("%x", h.Sum(nil))
	log.Println("yunke", md5str, str)

	var data JS_SDKLoginData_YunKe
	data.Username = username
	data.Sid = token

	var m JS_SDKLogin_YunKe
	m.GameId = gameId
	m.PlayerId = HF_Atoi(username)
	m.ChannelId = channelId
	m.PlayerToken = token
	m.Sign = md5str

	url := "https://mfgamesdk.mafengkj.com/loginvalid?" + fmt.Sprintf("gameid=%d&sessionid=%s&uid=%s&sign=%s", gameId, username, token, md5str)
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_YUNKE_IDENTIFICATION_ERROR"))
		return nil
	}

	result, err := ioutil.ReadAll(res.Body)
	if res.Body != nil {
		defer res.Body.Close()
	}
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_YUNKE_IDENTIFICATION_ERROR"))
		return nil
	}
	log.Println("result=", string(result))

	var ret JS_SDKBody_YunKe_IOS
	err = json.Unmarshal(result, &ret)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_YUNKE_IDENTIFICATION_ERROR"))
		return nil
	}

	switch ret.Status {
	case "0":
		LogInfo("登录云客SDK成功:", username, token)
	default:
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_YUNKE_REQUEST_PARAMETER_ERROR"))
		return nil
	}

	h1 := md5.New()
	LogDebug("New YunKe User:", username, token)
	h1.Write([]byte(token))
	account := fmt.Sprintf("%x", h1.Sum(nil))
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
	if _account.Uid > 0 {
		if password != _account.Password { //! 密码错误
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_PASSWORD"))
			return nil
		}
	} else { //! 插入
		_account.Account = account
		_account.Password = password
		_account.UserId = username
		_account.Creator = ctrl
		_account.Channelid = fmt.Sprintf("%d", channelId)
		_account.ServerId = serverid
		_account.Time = TimeServer().Unix()
		_account.Uid = InsertTable("san_account", &_account, 1, false)
		if _account.Uid <= 0 {
			self.SendErrInfo("err", GetCsvMgr().GetText("STR_ERROR"))
			return nil
		}

		sql := fmt.Sprintf("select * from `san_account` where `account` = '%s' and serverid =%d", account, serverid)
		GetServer().DBUser.GetOneData(sql, &_account, "", 0)

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
