package game

import (
	"bytes"
	//"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	//"encoding/json"
	"io/ioutil"
	"net/http"
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

type JS_SDKBody_YouXiFan struct {
	Code int                      `json:"code"`
	Data JS_SDKBody_YouXiFan_Data `json:"data"`
}

type JS_SDKBody_YouXiFan_Data struct {
	UserInfo JS_SDKBody_YouXiFan_Data_UserInfo `json:"userInfo"`
	IdInfo   JS_SDKBody_YouXiFan_Data_IdInfo   `json:"idInfo"`
}

type JS_SDKBody_YouXiFan_Data_UserInfo struct {
	UserId   string `json:"userId"`
	UserName string `json:"username"`
}

type JS_SDKBody_YouXiFan_Data_IdInfo struct {
	Id           string `json:"id"`
	IdType       int    `json:"idType"`
	Age          int    `json:"age"`
	Birthday     string `json:"birthday"`
	OverSea      bool   `json:"oversea"`
	VerifyStatus int    `json:"verifyStatus"`
}

// 安卓登录-拇指游玩登录
func (self *Session) SDKReg_YouXiFan(token string, serverid int, username string, ctrl string) *San_Account {
	log.Println("youxifan login:", username, token, serverid)

	arrParam := strings.Split(token, "||")
	if len(arrParam) < 3 {
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_YUNKE_REQUEST_PARAMETER_ERROR"))
		return nil
	}

	jsonData := make(map[string]string)
	jsonData["token"] = arrParam[1]
	b, _ := json.Marshal(jsonData)
	url := "https://sdkapiv2.youxifan.com/sdkapi/members/checkToken"
	res, err := http.Post(url, "application/json", bytes.NewBuffer(b))
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

	var ret JS_SDKBody_YouXiFan
	err = json.Unmarshal(result, &ret)
	if err != nil {
		log.Println(err)
		self.SendErrInfo("err", GetCsvMgr().GetText("STR_SESSION_YUNKE_IDENTIFICATION_ERROR"))
		return nil
	}

	h1 := md5.New()
	LogInfo("New youxifan User:", arrParam[1], arrParam[2])
	h1.Write([]byte(ret.Data.UserInfo.UserId))
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
		_account.UserId = ret.Data.UserInfo.UserId
		_account.Creator = ctrl
		//_account.Channelid = fmt.Sprintf("%d", channelid)
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
