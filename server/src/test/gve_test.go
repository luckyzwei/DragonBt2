package test

import (
	"encoding/json"
	. "game"
	"testing"
)
func TestFixed(t *testing.T)  {
	t.Log(HF_ToFixed(float64(8)/float64(15000) * 100.0, 2))
}

// The "omitempty" option specifies that the field should be omitted
// from the encoding if the field has an empty value, defined as
// false, 0, a nil pointer, a nil interface value, and any empty array,
// slice, map, or string.

func TestJson(t *testing.T)  {
	msg := &S2C_GveAction{}
	msg.PlayerRank = make(map[int64]*PlayerGlory)
	msg.PlayerRank[int64(1)] = &PlayerGlory{}
	msg.PlayerRank[int64(2)] = &PlayerGlory{}

	body, err := json.Marshal(msg)
	if err != nil {
		return
	}
	t.Log(string(body))
}

func TestBuildJson(t *testing.T)  {
	info := make(map[int]*GveBuild)
	pBuild := NewGveBuildInfo(1,1)
	pBuild.CurNum = 100
	pBuild.TakenInfo.Store(int64(123),100)
	pBuild.TakenInfo.Store(int64(124),100)
	pBuild.Encode()
	info[1] = pBuild
	gveBuildInfo := HF_JtoA(info)
	t.Log(gveBuildInfo)
	var buildInfo = make(map[int]*GveBuild)
	err := json.Unmarshal([]byte(gveBuildInfo), &buildInfo)
	if err != nil {
		t.Error("序列化buildInfo错误:", err.Error())
	}

}
