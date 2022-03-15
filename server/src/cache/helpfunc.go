package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const TIMEFORMAT = "2006-01-02 15:04:05" //! 时间格式化
const RESP_CODE_OK = 200

//! 加密消息
func HF_DecodeMsg(msghead string, msgdata []byte) []byte {
	data := &MsgBase{
		Msghead: proto.String(msghead),
		Msgtime: proto.Int64(time.Now().Unix()),
		Msgsign: proto.String("111"),
		Msgdata: []byte(msgdata),
	}

	msg, err := proto.Marshal(data)
	if err != nil {
		log.Println(err)
		return []byte("")
	}

	return msg
}

//! 克隆对象 dst为指针
func HF_DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

//! int取最小
func HF_MinInt(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

//! int取最大
func HF_MaxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

//! 得到一个随机数
func HF_GetRandom(num int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)
}

func HF_Atoi(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}

func HF_Atoi64(s string) int64 {
	num, _ := strconv.ParseInt(s, 10, 64)
	return num
}

func HF_JtoA(v interface{}) string {
	s, _ := json.Marshal(v)
	return string(s)
}

func HF_JtoB(v interface{}) []byte {
	s, _ := json.Marshal(v)
	return s
}

func HF_Atof(s string) float32 {
	num, _ := strconv.ParseFloat(s, 32)
	return float32(num)
}


func HF_GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}
