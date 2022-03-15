package test

import (
	"game"
	"log"
	"sync"
	"testing"
	"time"
)

// 矿战白盒测试
func TestMineDecode(t *testing.T) {
	minePlayer := game.NewMinePlayer(1, 1, 1001)
	fightTime := int(time.Now().Unix())
	minePlayer.FightRecord.Store(int64(1), fightTime)
	minePlayer.FightRecord.Store(int64(2), fightTime)
	minePlayer.FightRecord.Store(int64(3), fightTime)
	minePlayer.FightRecord.Store(int64(4), fightTime)
	minePlayer.Encode()
	t.Log(minePlayer.FightRecordInfo)
	minePlayer.FightRecord = new(sync.Map)
	minePlayer.Decode()
	t.Log(minePlayer.GetRecord())
}

func MockMinePlayer() *game.MinePlayer {
	minePlayer := game.NewMinePlayer(1, 1, 1001)
	fightTime := int(time.Now().Unix())
	minePlayer.FightRecord.Store(int64(1), fightTime)
	now := time.Now().Unix()
	minePlayer.CollectCd = now
	minePlayer.SpyCd = now
	minePlayer.FightCd = now
	minePlayer.AttackTimes = 10
	minePlayer.DefenceTimes = 2
	minePlayer.KillNum = 100
	minePlayer.DeadNum = 5
	minePlayer.CollectNum = 17
	minePlayer.Encode()
	minePlayer.MoveReason = game.MoveReasonMove
	return minePlayer
}

func ConnectDB() {
	server := game.GetServer()
	server.InitConfig()
	server.Redis = game.NewPool(server.Con.DBCon.Redis, server.Con.DBCon.RedisDB, server.Con.DBCon.RedisAuth)
	if server.Redis == nil {
		log.Fatal("redis err")
		return
	}
	server.DBUser.Init(server.Con.DBCon.DBUser)
}

// MinePlayer插入数据库, 更新, 查询
func TestMineInsert(t *testing.T) {
	ConnectDB()
	player := MockMinePlayer()
	player.Init("san_mineplayer", player, true)
	game.InsertTable("san_mineplayer", player, 0, false)
}

//Id               int64     // 主键
//NoticeTime       int64     // 提示时间
//MineStartTime    int64     // 下一次活动开始时间, 存盘
//MineEndTime      int64     // 下一次活动结束时间, 存盘
//BuffedMineId     int       // 双倍buff的矿点
//BuffedTimeCd     int64     // buff结束时间
//ScoreChanges     string    // 积分改变时间
//MineRankInfo     string    // 每次排行榜是否发送奖励, 记录当前活动的结束时间[需要每次存数据库]
//BuffCreatedInfo  string    // 双倍过了的Id
//MineRelationInfo string    // 矿点归属
//ScoreTimeChanges [3]int64  // 积分改变的时间
//MineRank         *sync.Map // 每次排行榜是否发送奖励, 记录当前活动的结束时间[需要每次存数据库]
//BuffCreated      *sync.Map // 双倍过了的Id
//MineRelation     *sync.Map // 矿点归属
func MockMineGlobal() *game.MineGlobal {
	mineGlobal := &game.MineGlobal{}
	mineGlobal.Id = 1
	now := time.Now().Unix()
	mineGlobal.NoticeTime = now
	mineGlobal.MineStartTime = now
	mineGlobal.MineEndTime = now
	mineGlobal.BuffedMineId = 1001
	mineGlobal.BuffedTimeCd = now
	mineGlobal.MineRank = new(sync.Map)
	mineGlobal.BuffCreated = new(sync.Map)
	mineGlobal.MineRelation = new(sync.Map)
	mineGlobal.ScoreTimeChanges[0] = now
	mineGlobal.ScoreTimeChanges[1] = now
	mineGlobal.ScoreTimeChanges[2] = now
	mineGlobal.MineRank.Store(now, 1)
	mineGlobal.BuffCreated.Store(1001, 1)
	mineGlobal.MineRelation.Store(1001, 1)
	mineGlobal.MineRelation.Store(2001, 2)
	mineGlobal.MineRelation.Store(3001, 3)

	return mineGlobal
}

func TestDecodeMineGlobal(t *testing.T) {
	mineGlobal := MockMineGlobal()
	mineGlobal.Encode()
	t.Log(mineGlobal.ScoreChanges, mineGlobal.MineRankInfo, mineGlobal.BuffCreatedInfo, mineGlobal.MineRelationInfo)
}

// MineGlobal插入数据库, 更新, 查询
func TestMineGlobalInsert(t *testing.T) {
	ConnectDB()
	info := MockMineGlobal()
	info.Init("san_mineglobal", info, true)
	info.Encode()
	game.InsertTable("san_mineglobal", info, 0, false)
}

// 获取数据库信息
func TestGetGlobalData(t *testing.T) {
	ConnectDB()
	game.GetMineMgr().GetGlobalData()
	mineGlobal := game.GetMineMgr().MineGlobal
	t.Log(mineGlobal.ScoreChanges, mineGlobal.MineRankInfo, mineGlobal.BuffCreatedInfo, mineGlobal.MineRelationInfo)
}
