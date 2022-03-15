package game

import (
	"encoding/json"
	"fmt"
	"master/utils"
	"reflect"
	"strings"
	"time"
)

//! 数据更新器
//! 约束必须有int型的tableKey且字段均为小写
type DataUpdate struct {
	baseData  reflect.Value //! 原始数据
	newData   reflect.Value //! 新数据
	tableName string        //! 表名
	dbServer  *DBServer     //! 数据库
	init      bool
	isRedis   bool //! redis
}

//! 初始化
func (self *DataUpdate) Init(tableName string, data interface{}, isredis bool) {
	self.newData = reflect.ValueOf(data).Elem()
	self.baseData = reflect.New(self.newData.Type()).Elem()
	self.baseData.Set(self.newData)
	self.tableName = tableName
	self.dbServer = GetServer().DBUser
	self.init = true
	self.isRedis = isredis
}

//! 更新数据,sql=false表示也更新redis
func (self *DataUpdate) Update(sql bool) {
	if !self.init {
		return
	}

	if self.isRedis {
		c := GetServer().GetRedisConn()
		defer c.Close()

		valueKey := self.newData.Field(0).Int()

		redisWriteOk, mySqlWriteOk := false, false
		_, err := c.Do("SETEX", fmt.Sprintf("%s_%d", self.tableName, valueKey), 864000, HF_JtoB(self.newData.Interface()))
		if err != nil {
			LogError("redis fail! query:", "SET", ",err:", err)
		} else {
			redisWriteOk = true
		}

		//if !sql {
		//	return
		//}

		valueList := ""
		//! 跳过tableKey
		for i := 1; i < self.baseData.NumField(); i++ {
			baseInt, newInt := int64(0), int64(0)
			baseStr, newStr := "", ""
			baseFloat, newFloat := float64(0.0), float64(0.0)

			//! 类型不同
			if self.baseData.Field(i).Type() != self.newData.Field(i).Type() {
				continue
			}

			switch self.baseData.Field(i).Kind() {
			case reflect.Int64:
				baseInt = self.baseData.Field(i).Int()
				newInt = self.newData.Field(i).Int()
			case reflect.Int:
				baseInt = self.baseData.Field(i).Int()
				newInt = self.newData.Field(i).Int()
			case reflect.Int8:
				baseInt = self.baseData.Field(i).Int()
				newInt = self.newData.Field(i).Int()
			case reflect.String:
				baseStr = self.baseData.Field(i).String()
				newStr = self.newData.Field(i).String()
			case reflect.Float32:
				baseFloat = self.baseData.Field(i).Float()
				newFloat = self.newData.Field(i).Float()
			case reflect.Float64:
				baseFloat = self.baseData.Field(i).Float()
				newFloat = self.newData.Field(i).Float()
			default:
				continue
			}

			rowName := strings.ToLower(self.baseData.Type().Field(i).Name)

			if baseInt != newInt {
				valueList += fmt.Sprintf("`%s`=%d,", rowName, int(newInt))
				self.baseData.Field(i).SetInt(newInt)
			} else if baseStr != newStr {
				valueList += fmt.Sprintf("`%s`='%s',", rowName, newStr)
				self.baseData.Field(i).SetString(newStr)
			} else if baseFloat != newFloat {
				valueList += fmt.Sprintf("`%s`=%f,", rowName, newFloat)
				self.baseData.Field(i).SetFloat(newFloat)
			}
		}

		if valueList != "" {
			valueKey := self.baseData.Field(0).Int()
			tableKey := strings.ToLower(self.baseData.Type().Field(0).Name)

			updateQuery := fmt.Sprintf("update `%s` set %s where `%s`=%d limit 1", self.tableName, valueList, tableKey, valueKey)
			//! 去掉多余逗号
			updateQuery = strings.Replace(updateQuery, ", where", " where", 1)
			GetServer().SqlSet(updateQuery)
			mySqlWriteOk = true
		}

		if !redisWriteOk && mySqlWriteOk {
			LogError("data may be error, !redisWriteOk && mySqlWriteOk")
		} else if !redisWriteOk && !mySqlWriteOk {
			LogError("data may be error, !redisWriteOk && !mySqlWriteOk")
		}

		return
	}

	valueList := ""
	//! 跳过tableKey
	for i := 1; i < self.baseData.NumField(); i++ {
		baseInt, newInt := int64(0), int64(0)
		baseStr, newStr := "", ""
		baseFloat, newFloat := float64(0.0), float64(0.0)

		//! 类型不同
		if self.baseData.Field(i).Type() != self.newData.Field(i).Type() {
			continue
		}

		switch self.baseData.Field(i).Kind() {
		case reflect.Int64:
			baseInt = self.baseData.Field(i).Int()
			newInt = self.newData.Field(i).Int()
		case reflect.Int:
			baseInt = self.baseData.Field(i).Int()
			newInt = self.newData.Field(i).Int()
		case reflect.Int8:
			baseInt = self.baseData.Field(i).Int()
			newInt = self.newData.Field(i).Int()
		case reflect.String:
			baseStr = self.baseData.Field(i).String()
			newStr = self.newData.Field(i).String()
		case reflect.Float32:
			baseFloat = self.baseData.Field(i).Float()
			newFloat = self.newData.Field(i).Float()
		case reflect.Float64:
			baseFloat = self.baseData.Field(i).Float()
			newFloat = self.newData.Field(i).Float()
		default:
			continue
		}

		rowName := strings.ToLower(self.baseData.Type().Field(i).Name)

		if baseInt != newInt {
			valueList += fmt.Sprintf("`%s`=%d,", rowName, int(newInt))
			self.baseData.Field(i).SetInt(newInt)
		} else if baseStr != newStr {
			valueList += fmt.Sprintf("`%s`='%s',", rowName, newStr)
			self.baseData.Field(i).SetString(newStr)
		} else if baseFloat != newFloat {
			valueList += fmt.Sprintf("`%s`=%f,", rowName, newFloat)
			self.baseData.Field(i).SetFloat(newFloat)
		}
	}

	if valueList != "" {
		valueKey := self.baseData.Field(0).Int()
		tableKey := strings.ToLower(self.baseData.Type().Field(0).Name)

		updateQuery := fmt.Sprintf("update `%s` set %s where `%s`=%d limit 1", self.tableName, valueList, tableKey, valueKey)
		//! 去掉多余逗号
		updateQuery = strings.Replace(updateQuery, ", where", " where", 1)
		GetServer().SqlBaseSet(updateQuery)
	}
}

//! isredis=true,插入redis以及mysql数据库
//! isredis=false, 插入mysql数据库
func InsertTable(table string /*表名*/, data interface{} /*数据*/, index int /*从第几个字段开始*/, isredis bool) int64 {
	if isredis {
		c := GetServer().GetRedisConn()
		defer c.Close()

		// 插入redis数据库
		valueKey := reflect.ValueOf(data).Elem().Field(0).Int()
		_, err := c.Do("SETEX", fmt.Sprintf("%s_%d", table, valueKey), 864000, HF_JtoB(data))
		if err != nil {
			LogError("redis fail! query:", "SET", ",err:", err)
		}

		edata := reflect.ValueOf(data).Elem()

		tableList := ""
		valueList := ""
		for i := index; i < edata.NumField(); i++ {
			switch edata.Field(i).Kind() {
			case reflect.Int64:
				valueList += fmt.Sprintf("%d", edata.Field(i).Int())
				valueList += ","
			case reflect.Int:
				valueList += fmt.Sprintf("%d", edata.Field(i).Int())
				valueList += ","
			case reflect.Int8:
				valueList += fmt.Sprintf("%d", edata.Field(i).Int())
				valueList += ","
			case reflect.String:
				valueList += fmt.Sprintf("'%s'", edata.Field(i).String())
				valueList += ","
			case reflect.Float32:
				valueList += fmt.Sprintf("%f", edata.Field(i).Float())
				valueList += ","
			case reflect.Float64:
				valueList += fmt.Sprintf("%f", edata.Field(i).Float())
				valueList += ","
			default:
				continue
			}

			rowName := strings.ToLower(edata.Type().Field(i).Name)
			tableList += fmt.Sprintf("`%s`", rowName)
			tableList += ","
		}

		//LogError(fmt.Sprintf("%s------------------valueList=%s", table, valueList))
		// 插入mysql数据库
		if valueList != "" {
			updateQuery := fmt.Sprintf("insert into `%s`(%s) values(%s)", table, tableList, valueList)
			//! 去掉多余逗号
			updateQuery = strings.Replace(updateQuery, ",)", ")", 2)
			GetServer().SqlSet(updateQuery)
		}

		return valueKey
	}

	edata := reflect.ValueOf(data).Elem()

	tableList := ""
	valueList := ""
	for i := index; i < edata.NumField(); i++ {
		switch edata.Field(i).Kind() {
		case reflect.Int64:
			valueList += fmt.Sprintf("%d", edata.Field(i).Int())
			valueList += ","
		case reflect.Int:
			valueList += fmt.Sprintf("%d", edata.Field(i).Int())
			valueList += ","
		case reflect.Int8:
			valueList += fmt.Sprintf("%d", edata.Field(i).Int())
			valueList += ","
		case reflect.String:
			valueList += fmt.Sprintf("'%s'", edata.Field(i).String())
			valueList += ","
		case reflect.Float32:
			valueList += fmt.Sprintf("%f", edata.Field(i).Float())
			valueList += ","
		case reflect.Float64:
			valueList += fmt.Sprintf("%f", edata.Field(i).Float())
			valueList += ","
		default:
			continue
		}

		rowName := strings.ToLower(edata.Type().Field(i).Name)
		tableList += fmt.Sprintf("`%s`", rowName)
		tableList += ","
	}

	if valueList != "" {
		updateQuery := fmt.Sprintf("insert into `%s`(%s) values(%s)", table, tableList, valueList)
		//! 去掉多余逗号
		updateQuery = strings.Replace(updateQuery, ",)", ")", 2)
		//! 执行
		lastid, _, _ := GetServer().DBUser.Exec(updateQuery)

		return lastid
	}

	return 0
}

//! 插入
func InsertTableBatch(table string /*表名*/, data []interface{} /*数据*/, index int /*从第几个字段开始*/, db *DBServer) int64 {
	if len(data) == 0 {
		return 0
	}
	edata := reflect.ValueOf(data[0]).Elem()

	tableList := ""
	valueList := ""
	for i := index; i < edata.NumField(); i++ {
		switch edata.Field(i).Kind() {
		case reflect.Int64:
			valueList += fmt.Sprintf("%d", edata.Field(i).Int())
			valueList += ","
		case reflect.Int:
			valueList += fmt.Sprintf("%d", edata.Field(i).Int())
			valueList += ","
		case reflect.Int8:
			valueList += fmt.Sprintf("%d", edata.Field(i).Int())
			valueList += ","
		case reflect.String:
			valueList += fmt.Sprintf("'%s'", edata.Field(i).String())
			valueList += ","
		case reflect.Float32:
			valueList += fmt.Sprintf("%f", edata.Field(i).Float())
			valueList += ","
		case reflect.Float64:
			valueList += fmt.Sprintf("%f", edata.Field(i).Float())
			valueList += ","
		default:
			continue
		}

		rowName := strings.ToLower(edata.Type().Field(i).Name)
		tableList += fmt.Sprintf("`%s`", rowName)
		tableList += ","
	}

	if valueList != "" {
		updateQuery := fmt.Sprintf("insert into `%s`(%s) values(%s)", table, tableList, valueList)
		//! 去掉多余逗号
		updateQuery = strings.Replace(updateQuery, ",)", ")", 2)
		//! 执行
		lastid, _, _ := GetServer().DBLog.Exec(updateQuery)

		return lastid
	}

	return 0
}

//! 插入
func InsertLogTable(table string /*表名*/, data interface{} /*数据*/, index int /*从第几个字段开始*/) int64 {
	edata := reflect.ValueOf(data).Elem()

	tableList := ""
	valueList := ""
	for i := index; i < edata.NumField(); i++ {
		switch edata.Field(i).Kind() {
		case reflect.Int64:
			valueList += fmt.Sprintf("%d", edata.Field(i).Int())
			valueList += ","
		case reflect.Int:
			valueList += fmt.Sprintf("%d", edata.Field(i).Int())
			valueList += ","
		case reflect.Int8:
			valueList += fmt.Sprintf("%d", edata.Field(i).Int())
			valueList += ","
		case reflect.String:
			valueList += fmt.Sprintf("'%s'", edata.Field(i).String())
			valueList += ","
		case reflect.Float32:
			valueList += fmt.Sprintf("%f", edata.Field(i).Float())
			valueList += ","
		case reflect.Float64:
			valueList += fmt.Sprintf("%f", edata.Field(i).Float())
			valueList += ","
		default:
			continue
		}

		rowName := strings.ToLower(edata.Type().Field(i).Name)
		tableList += fmt.Sprintf("`%s`", rowName)
		tableList += ","
	}

	if valueList != "" {
		updateQuery := fmt.Sprintf("insert into `%s`(%s) values(%s)", table, tableList, valueList)
		//! 去掉多余逗号
		updateQuery = strings.Replace(updateQuery, ",)", ")", 2)
		//! 执行
		lastid, _, _ := GetServer().DBLog.Exec(updateQuery)

		return lastid
	}

	return 0
}

//! 删除
func DeleteTable(table string /*表名*/, data interface{} /*数据*/, index []int /*第几个字段作为主键*/) {
	edata := reflect.ValueOf(data).Elem()

	valueList := ""
	for i := 0; i < len(index); i++ {
		rowName := strings.ToLower(edata.Type().Field(index[i]).Name)

		switch edata.Field(index[i]).Kind() {
		case reflect.Int64:
			valueList += fmt.Sprintf("`%s` = %d", rowName, edata.Field(index[i]).Int())
		case reflect.Int:
			valueList += fmt.Sprintf("`%s` = %d", rowName, edata.Field(index[i]).Int())
		case reflect.Int8:
			valueList += fmt.Sprintf("`%s` = %d", rowName, edata.Field(index[i]).Int())
		case reflect.String:
			valueList += fmt.Sprintf("`%s` = '%s'", rowName, edata.Field(index[i]).String())
		default:
			continue
		}

		if i < len(index)-1 {
			valueList += " and "
		}
	}

	if valueList != "" {
		updateQuery := fmt.Sprintf("delete from `%s` where %s", table, valueList)
		//! 执行
		GetServer().DBUser.Exec(updateQuery)
	}
}

//! 更新数据
func (self *DataUpdate) UpdateEx(field string, fieldVal int) {
	if !self.init {
		return
	}

	valueList := ""
	//! 跳过tableKey
	for i := 1; i < self.baseData.NumField(); i++ {
		baseInt, newInt := int64(0), int64(0)
		baseStr, newStr := "", ""
		baseFloat, newFloat := float64(0.0), float64(0.0)

		//! 类型不同
		if self.baseData.Field(i).Type() != self.newData.Field(i).Type() {
			continue
		}

		switch self.baseData.Field(i).Kind() {
		case reflect.Int64:
			baseInt = self.baseData.Field(i).Int()
			newInt = self.newData.Field(i).Int()
		case reflect.Int:
			baseInt = self.baseData.Field(i).Int()
			newInt = self.newData.Field(i).Int()
		case reflect.Int8:
			baseInt = self.baseData.Field(i).Int()
			newInt = self.newData.Field(i).Int()
		case reflect.String:
			baseStr = self.baseData.Field(i).String()
			newStr = self.newData.Field(i).String()
		case reflect.Float32:
			baseFloat = self.baseData.Field(i).Float()
			newFloat = self.newData.Field(i).Float()
		case reflect.Float64:
			baseFloat = self.baseData.Field(i).Float()
			newFloat = self.newData.Field(i).Float()
		default:
			continue
		}

		rowName := strings.ToLower(self.baseData.Type().Field(i).Name)

		if baseInt != newInt {
			valueList += fmt.Sprintf("`%s`=%d,", rowName, int(newInt))
			self.baseData.Field(i).SetInt(newInt)
		} else if baseStr != newStr {
			valueList += fmt.Sprintf("`%s`='%s',", rowName, newStr)
			self.baseData.Field(i).SetString(newStr)
		} else if baseFloat != newFloat {
			valueList += fmt.Sprintf("`%s`=%f,", rowName, newFloat)
			self.baseData.Field(i).SetFloat(newFloat)
		}
	}

	if valueList != "" {
		valueKey := self.baseData.Field(0).Int()
		tableKey := strings.ToLower(self.baseData.Type().Field(0).Name)

		updateQuery := fmt.Sprintf("update `%s` set %s where `%s`=%d and `%s` = %d limit 1", self.tableName, valueList, tableKey, valueKey, field, fieldVal)
		//! 去掉多余逗号
		updateQuery = strings.Replace(updateQuery, ", where", " where", 1)
		GetServer().SqlBaseSet(updateQuery)
	}
}

// 插入redis[国战, 矿战战报], HMSET
func HMSetRedis(tableName string, fightId int64, data interface{}, timeout int64) error {
	c := GetServer().GetRedisConn()
	defer c.Close()

	// 查询ttl
	ttl, err := GetRedisMgr().GetTTL(tableName)
	if err != nil {
		LogError(err.Error())
		return err
	}

	// 设置ttl
	if ttl <= 0 || ttl > 86400 {
		_, err := GetRedisMgr().Expire(tableName, timeout)
		if err != nil {
			LogError(err.Error())
			return err
		}
	}

	item := map[string]interface{}{
		fmt.Sprintf("%d", fightId): HF_JtoA(data),
	}
	GetRedisMgr().HMSet(tableName, item)

	return nil
}

// 获取redis数据
func HGetRedis(tableName string, field string) (string, bool, error) {
	return GetRedisMgr().HGet(tableName, field)
}

// 获取redis数据HMGet
func GetAllRedis(tableName string) (map[string]string, error) {
	items, err := GetRedisMgr().HGetAll(tableName)
	if err != nil {
		LogError(err.Error())
		return nil, err
	}
	return items, err
}

// 单条数据正常过期
func HMSetRedisEx(tableName string, fightId int64, data interface{}, timeout int64) error {
	c := GetServer().GetRedisConn()
	defer c.Close()

	_, err := c.Do("SETEX", fmt.Sprintf("%s_%d", tableName, fightId), timeout, utils.HF_JtoB(data))
	if err != nil {
		utils.LogError("redis fail! query:", "SET", ",err:", err)
	}

	return err
}

// 获取redis数据
func HGetRedisEx(tableName string, id int64, field string) (string, bool, error) {
	tableKey := fmt.Sprintf("%s_%d", tableName, id)
	return GetRedisMgr().Get(tableKey)
}

//infoName san_towerbattleinfo
//recordName san_towerbattlerecord
//recordType core.BATTLE_TYPE_TOWER
// tableName tbl_crossarenarecord
func MigrateDataOne(infoName, recordName, tableName string, recordType int, cursor int64) (int64, int) {
	var info BattleInfo
	var record BattleRecord
	cursor, dataSlice, err := GetRedisMgr().HScan(infoName, cursor, "*", 20) //c.Do("hscan", infoName, 1, "match *", "count 1")
	if err != nil {
		//! 读取配置
		//self.migrateOK = true
		return 0, 0
	}

	lenArr := len(dataSlice) / 2

	if lenArr > 0 {
		for i := 0; i < lenArr; i++ {
			json.Unmarshal([]byte(dataSlice[i*2+1]), &info)

			v1, ret, err := HGetRedis(recordName, fmt.Sprintf("%d", info.Id))
			if err != nil {
				//self.migrateOK = true
				return 0, 0
			}

			//! 获取成功
			if ret {
				json.Unmarshal([]byte(v1), &record)
			}

			var db_battleInfo JS_CrossArenaBattleInfo
			sql := fmt.Sprintf("select * from `%s` where fightid=%d limit 1;", tableName, info.Id)
			ret1 := GetServer().DBUser.GetOneData(sql, &db_battleInfo, "", 0)
			if ret1 == true && db_battleInfo.Id > 0 { //! 获取成功
				LogInfo("已存在，更新：", cursor, dataSlice[0], info.Id)
				db_battleInfo.BattleInfo = HF_CompressAndBase64(HF_JtoB(&info))
				db_battleInfo.BattleRecord = HF_CompressAndBase64(HF_JtoB(&record))
				db_battleInfo.UpdateTime = time.Now().Unix()
				db_battleInfo.Update(true)
			} else {
				db_battleInfo.FightId = info.Id
				db_battleInfo.RecordType = recordType
				db_battleInfo.BattleInfo = HF_CompressAndBase64(HF_JtoB(&info))
				db_battleInfo.BattleRecord = HF_CompressAndBase64(HF_JtoB(&record))
				db_battleInfo.UpdateTime = time.Now().Unix()
				InsertTable(tableName, &db_battleInfo, 0, false)
			}
		}
	}

	return cursor, lenArr
}
