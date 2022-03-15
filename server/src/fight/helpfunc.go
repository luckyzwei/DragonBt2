package main

import (
	"encoding/json"
	"runtime/debug"
)

func HF_JtoB(v interface{}) []byte {
	s, err := json.Marshal(v)
	if err != nil {
		LogError("HF_JtoB err:", string(debug.Stack()))
	}
	return s
}
