package game

import (
	//"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"strings"

	//"encoding/json"
	"fmt"
	//"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	//"io/ioutil"
	"log"
	//"net/http"
	//"net/url"
	//"time"
)

type JS_SDKLogin_Common struct {
	Service string `json:"service"`
	AppID   int    `json:"appid"`
	Data    string `json:"data"`
	Sign    string `json:"sign"`
}

type JS_SDKLoginData_Common struct {
	Sid      string `json:"sid"`
	Username string `json:"username"`
}

type JS_SDKData_Common struct {
	Uid string `json:"uid"`
}

type JS_SDKBody_Common struct {
	Id    string         `json:"id"`
	State JS_SDKState    `json:"state"`
	Data  JS_SDKData_MZY `json:"data"`
}

// 安卓登录-拇指游玩登录
func (self *Session) SDKReg_Common(token string, serverid int, username string, ctrl string) *San_Account {
	log.Println("mzy login:", username, token, serverid)

	loginKey := "184506d927d68ee7fd41bc1b876ad5b6"

	arrParam := strings.Split(token, "||")
	if len(arrParam) < 4 {
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_YUNKE_REQUEST_PARAMETER_ERROR"))
		return nil
	}

	channelid := 0
	userId := ""
	if len(arrParam) >= 7 {
		appid := arrParam[4]
		userId = arrParam[1]
		channelid = HF_Atoi(arrParam[6])
		if HF_Atoi(appid) == 90690 {
			loginKey = "184506d927d68ee7fd41bc1bbui5b6"
		} else if HF_Atoi(appid) == 90689 {
			loginKey = "32f31af59c3b7a097dd4f648d4ef9da0"
		} else if HF_Atoi(appid) == 90791 {
			loginKey = "fa12c9b27524f30ea534728f7a81fae4"
		} else if HF_Atoi(appid) == 90798 {
			loginKey = "9379c38c61acce532452c0940dfebd81"
		} else if HF_Atoi(appid) == 90816 {
			loginKey = "e07815ca9dffc9b8d360efa95ee556aa"
		} else if HF_Atoi(appid) == 90820 {
			loginKey = "6744fcd6a22c52fa9eb43986283c1953"
		} else if HF_Atoi(appid) == 90875 {
			loginKey = "15f51f026eb0b2eca877fc29ee12c68f"
		} else if HF_Atoi(appid) == 90900 {
			loginKey = "e9099e39495de470f07852e4875f569b"
		} else if HF_Atoi(appid) == 90904 {
			loginKey = "9878549498f62d4b8005ae295dc90b2e"
		} else if HF_Atoi(appid) == 90929 {
			loginKey = "69172266f61ead02d60cade8e0037373"
		} else if HF_Atoi(appid) == 90935 {
			loginKey = "637c5584110aafed6ed677f0337570fb"
		} else if HF_Atoi(appid) == 90937 {
			loginKey = "b633e016c4f764cb525f3c66d5f2643f"
		} else if HF_Atoi(appid) == 90940 {
			loginKey = "8618359185c9292bcdc0012e3674662d"
		} else if HF_Atoi(appid) == 90944 {
			loginKey = "8c87635e3e5ec8a39159a879da829170"
		}
	}

	//test string : sdk_common||13298ig1165914||1585899790||287d6c38bc2a982c4f167318bb86fd2b
	//appkey := "0a7357b0c3d8023c920150aed32b1754"
	//appid := 1000037
	//
	str := arrParam[1] + arrParam[2] + loginKey
	//str := fmt.Sprintf("%dsdk.game.checkentersid=%s&username=%s%s", appid, token, username, appkey)
	//l := url.QueryEscape(str)
	h := md5.New()
	h.Write([]byte(str))
	md5str := fmt.Sprintf("%x", h.Sum(nil))
	log.Println(md5str, str, str)

	if md5str != arrParam[3] {
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_YUNKE_REQUEST_PARAMETER_ERROR"))
		return nil
	}

	h1 := md5.New()
	LogInfo("New quxuan User:", arrParam[1], arrParam[2])
	h1.Write([]byte(arrParam[1]))
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
		_account.UserId = userId
		_account.Creator = ctrl
		_account.Channelid = fmt.Sprintf("%d", channelid)
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
