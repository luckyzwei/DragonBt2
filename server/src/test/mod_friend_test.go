package test

import (
	"encoding/json"
	"errors"
	. "game"
	"reflect"
	"sync"
	"testing"
)

func TestFriendMap(t *testing.T) {
	friend := new(sync.Map)

	for i := 1; i < 11; i++ {
		info := &JS_FriendNode{Uid: int64(i), Vip: i}
		friend.Store(int64(i), info)
	}

	friend.Range(func(key, value interface{}) bool {
		info, ok := value.(*JS_FriendNode)
		if !ok {
			return true
		}

		info.Vip = 15
		return true
	})

	friend.Range(func(key, value interface{}) bool {
		_, ok := value.(*JS_FriendNode)
		if !ok {
			return true
		}
		//fmt.Printf("%#v\n", info)
		return true
	})

	var data string

	var friendInfo []*JS_FriendNode

	friend.Range(func(key, value interface{}) bool {
		info, ok := value.(*JS_FriendNode)
		if !ok {
			return true
		}

		friendInfo = append(friendInfo, info)
		return true
	})

	v, ok := json.Marshal(&friendInfo)
	if ok == nil {
		data = string(v)
	}

	var friendInfoOut []*JS_FriendNode
	json.Unmarshal([]byte(data), &friendInfoOut)
	friend = new(sync.Map)

	for k, v := range friendInfoOut {
		friend.LoadOrStore(k, v)
	}

	friend.Range(func(key, value interface{}) bool {
		_, ok := value.(*JS_FriendNode)
		if !ok {
			return true
		}
		//fmt.Printf("%#v\n", info)
		return true
	})

	if reflect.DeepEqual(friendInfoOut, friendInfo) {
		t.Error(errors.New("friend data is not equal!"))
	}
}
