package network

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"master/utils"
	"time"
)

//! 解消息
func HF_EncodeMsg(msg []byte) (string, []byte, bool) {
	data := &MsgBase{}
	err := proto.Unmarshal(msg, data)
	if err != nil {
		log.Println("utils.HF_EncodeMsg:", err)
		return "", []byte(""), false
	}

	return data.GetMsghead(), data.GetMsgdata(), true
}

//! 加密消息
func HF_DecodeMsg(msghead string, msgdata []byte) []byte {
	timeStart := time.Now()
	data := &MsgBase{
		Msghead: proto.String(msghead),
		Msgtime: proto.Int64(time.Now().Unix()),
		Msgsign: proto.String("111"),
		Msgdata: []byte(msgdata),
	}

	msg, err := proto.Marshal(data)
	if err != nil {
		log.Println("utils.HF_DecodeMsg:", err)
		return []byte("")
	}
	var buf bytes.Buffer
	compressor := zlib.NewWriter(&buf)
	compressor.Write(msg)
	compressor.Close()

	size := []byte(fmt.Sprintf("%d ", len(msg)))
	newmsg := []byte("")
	newmsg = append(newmsg, size...)
	newmsg = append(newmsg, buf.Bytes()...)

	utils.LogDebug("message zlib cost: ", len(msgdata), buf.Len(), time.Since(timeStart), ", cid:", msghead)
	return newmsg
}
