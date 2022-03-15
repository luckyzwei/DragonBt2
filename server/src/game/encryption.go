package game

import (
	"fmt"
	"strconv"
	"strings"
)

func encryptData(code_data string, callback_key string) string {
	dataArr := []rune(code_data)
	keyArr := []byte(callback_key)
	keyLen := len(keyArr)

	var tmpList []int

	for index, value := range dataArr {
		base := int(value)
		dataString := base + int(0xFF&keyArr[index%keyLen])
		tmpList = append(tmpList, dataString)
	}

	var str string

	for _, value := range tmpList {
		str += "@" + fmt.Sprintf("%d", value)
	}
	return str
}

func decryptData(nt_data string, callback_key string) string {
	strLen := len(nt_data)
	newData := []rune(nt_data)
	resultData := string(newData[1:strLen])
	dataArr := strings.Split(resultData, "@")
	keyArr := []byte(callback_key)
	keyLen := len(keyArr)

	var tmpList []int

	for index, value := range dataArr {
		base, _ := strconv.Atoi(value)
		dataString := base - int(0xFF&keyArr[index%keyLen])
		tmpList = append(tmpList, dataString)
	}

	var str string

	for _, val := range tmpList {
		str += string(rune(val))
	}
	return str
}
