package main


import (
	"fmt"
	//"log"
	//"github.com/garyburd/redigo/redis"
	"reflect"
	"strings"
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

//! 更新数据
func (self *DataUpdate) Update(sql bool) {
	if !self.init {
		return
	}

	if self.isRedis {
		c := GetServer().GetRedisConn()
		defer c.Close()

		valueKey := self.newData.Field(0).Int()

		_, err := c.Do("SETEX", fmt.Sprintf("%s_%d", self.tableName, valueKey), 864000, HF_JtoB(self.newData.Interface()))
		if err != nil {
			LogError("redis fail! query:", "SET", ",err:", err)
		}

		if !sql {
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
			//! 转成小写
			//updateQuery = strings.ToLower(updateQuery)
			GetServer().SqlSet(updateQuery)
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
		//! 转成小写
		//updateQuery = strings.ToLower(updateQuery)
		//! 执行
		//self.dbServer.Exec(updateQuery)
		GetServer().SqlSet(updateQuery)
	}
}

////! 插入,含首字段,必须先调用init
//func (self *DataUpdate) Insert() {
//	tableList := ""
//	valueList := ""
//	for i := 0; i < self.baseData.NumField(); i++ {
//		switch self.baseData.Field(i).Kind() {
//		case reflect.Int64:
//			valueList += fmt.Sprintf("%d", self.newData.Field(i).Int())
//			valueList += ","
//		case reflect.Int:
//			valueList += fmt.Sprintf("%d", self.newData.Field(i).Int())
//			valueList += ","
//		case reflect.Int8:
//			valueList += fmt.Sprintf("%d", self.newData.Field(i).Int())
//			valueList += ","
//		case reflect.String:
//			valueList += fmt.Sprintf("'%s'", self.newData.Field(i).String())
//			valueList += ","
//		case reflect.Float32:
//			valueList += fmt.Sprintf("%f", self.newData.Field(i).Float())
//			valueList += ","
//		case reflect.Float64:
//			valueList += fmt.Sprintf("%f", self.newData.Field(i).Float())
//			valueList += ","
//		default:
//			continue
//		}

//		rowName := strings.ToLower(self.baseData.Type().Field(i).Name)
//		tableList += fmt.Sprintf("`%s`", rowName)
//		tableList += ","
//	}

//	if valueList != "" {
//		updateQuery := fmt.Sprintf("insert into `%s`(%s) values(%s)", self.tableName, tableList, valueList)
//		//! 去掉多余逗号
//		updateQuery = strings.Replace(updateQuery, ",)", ")", 2)
//		//! 执行
//		self.dbServer.Exec(updateQuery)
//	}
//}

//! 插入
func InsertTable(table string /*表名*/, data interface{} /*数据*/, index int /*从第几个字段开始*/, isredis bool) int64 {
	if isredis {
		c := GetServer().GetRedisConn()
		defer c.Close()

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
