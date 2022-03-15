//! 数据库底层

package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"master/utils"
	"reflect"
	"strings"
	"time"
)

//! 常数定义
const (
	Max_Open_Conn = 100
	Max_Idle_Conn = 30
	Max_Life_Time = 28800
)

//! 数据库结构
type DBServer struct {
	m_db     *sql.DB //! db
	m_dbName string  //! 库名
}

//! 得到库名
func (self *DBServer) GetDBName() string {
	return self.m_dbName
}

//! 连接数据库
//! dsn root:Wnr*JS4*qUyy95ll@tcp(192.168.20.126:3306)/football_dynamic?charset=utf8&timeout=10s
func (self *DBServer) Init(dsn string) bool {
	self.parseDBName(dsn)

	db, err := sql.Open("mysql", dsn)
	db.SetMaxOpenConns(Max_Open_Conn)
	db.SetMaxIdleConns(Max_Idle_Conn)
	db.SetConnMaxLifetime(Max_Life_Time * time.Second)

	if err != nil {
		//db.Close()
		log.Fatalln("db open fail! err:%s dsn:%s", err.Error(), self.m_dbName)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatalln("db open ping fail!  err:%s dns:%s", err.Error(), self.m_dbName)
	}

	self.m_db = db

	//log.Println("db connect!", self.m_dbName)
	utils.LogDebug("db connect!", self.m_dbName)

	return true
}

func (self *DBServer) Close() {
	if self.m_db != nil {
		self.m_db.Close()
	}
}

func (self *DBServer) checkError(info string, sql string, err error) {
	log.Println(info, sql, ",err:", err)
	utils.LogError(info, sql, ",err:", err)
}

//! 执行语句
func (self *DBServer) Exec(query string, args ...interface{}) (int64, int64, bool) {
	sql := fmt.Sprintf(query, args...)
	result, err := self.m_db.Exec(sql)
	if err != nil {
		self.checkError("db exec fail! query:", sql, err)
		return 0, 0, false
	}

	LastInsertId := int64(0)
	LastInsertId, err = result.LastInsertId()
	if err != nil {
		self.checkError("db exec-LastInsertId fail! query:", sql, err)
	}

	RowsAffected := int64(0)
	RowsAffected, err = result.RowsAffected()
	if err != nil {
		self.checkError("db exec-RowsAffected fail! query:", sql, err)
	}
	return LastInsertId, RowsAffected, true
}

//! 得到一条数据
func (self *DBServer) GetOneData(query string, struc interface{}, table string, key int64) bool {
	if table != "" && key != 0 {
		c := GetRedisMgr().GetRedisConn()
		defer c.Close()

		v, err := redis.Bytes(c.Do("GET", fmt.Sprintf("%s_%d", table, key)))
		if err == nil {
			json.Unmarshal(v, struc)
			//log.Println("从redis读取", string(v))
		} else {
			log.Println("从数据库读取", query)
			rows, err := self.m_db.Query(query)
			defer rows.Close()
			if rows == nil || err != nil {
				log.Println("db GetOneData fail! query:", query, ",err:", err)
				utils.LogError("db GetOneData fail! query:", query, ",err:", err)
				return false
			}

			//! 得到反射
			s := reflect.ValueOf(struc).Elem()
			num := s.NumField()
			data := make([]interface{}, 0)
			for i := 0; i < num; i++ {
				ki := s.Field(i).Kind()
				if ki != reflect.Int && ki != reflect.Int64 && ki != reflect.Int8 && ki != reflect.String && ki != reflect.Float32 && ki != reflect.Float64 {
					continue
				}
				data = append(data, s.Field(i).Addr().Interface())
			}

			has := false
			for rows.Next() {
				has = true
				err = rows.Scan(data...)
				if err != nil {
					log.Println("db GetOneData-Scan fail! query:", query, ",err:", err)
					utils.LogError("db GetOneData-Scan fail! query:", query, ",err:", err)
					return false
				}
				break
			}
			//! 记录到redis
			if has {
				c.Do("SETEX", fmt.Sprintf("%s_%d", table, key), 864000, utils.HF_JtoB(struc))
			}
		}
		return true
	}

	rows, err := self.m_db.Query(query)
	if rows != nil {
		defer rows.Close()
	}
	if rows == nil || err != nil {
		log.Println("db GetOneData fail! query:", query, ",err:", err)
		utils.LogError("db GetOneData fail! query:", query, ",err:", err)
		return false
	}

	//! 得到反射
	s := reflect.ValueOf(struc).Elem()
	num := s.NumField()
	data := make([]interface{}, 0)
	for i := 0; i < num; i++ {
		ki := s.Field(i).Kind()
		if ki != reflect.Int && ki != reflect.Int64 && ki != reflect.Int8 && ki != reflect.String && ki != reflect.Float32 && ki != reflect.Float64 {
			continue
		}
		data = append(data, s.Field(i).Addr().Interface())
	}

	for rows.Next() {
		err = rows.Scan(data...)
		if err != nil {
			log.Println("db GetOneData-Scan fail! query:", query, ",err:", err)
			utils.LogError("db GetOneData-Scan fail! query:", query, ",err:", err)
			return false
		}
		break
	}

	return true
}

//! 得到多条数据
func (self *DBServer) GetAllData(query string, struc interface{}) []interface{} {
	rows, err := self.m_db.Query(query)
	if rows != nil {
		defer rows.Close()
	}
	if rows == nil || err != nil {
		log.Println("db GetAllData fail! query:", query, ",err:", err)
		utils.LogError("db GetAllData fail! query:", query, ",err:", err)
		return nil
	}

	//! 得到反射
	s := reflect.ValueOf(struc).Elem()
	num := s.NumField()
	data := make([]interface{}, 0)
	for i := 0; i < num; i++ {
		ki := s.Field(i).Kind()
		if ki != reflect.Int && ki != reflect.Int64 && ki != reflect.Int8 && ki != reflect.String && ki != reflect.Float32 && ki != reflect.Float64 {
			continue
		}
		data = append(data, s.Field(i).Addr().Interface())
	}

	result := make([]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(data...)
		if err != nil {
			log.Println("db GetAllData-Scan fail! query:", query, ",err:", err)
			utils.LogError("db GetAllData-Scan fail! query:", query, ",err:", err)
			return nil
		}
		newObj := reflect.New(reflect.TypeOf(struc).Elem()).Elem()
		newObj.Set(s)
		result = append(result, newObj.Addr().Interface())
	}

	return result
}

//! 解析库名
func (self *DBServer) parseDBName(dsn string) {
	name := dsn
	begin := strings.Index(name, "/")
	end := strings.Index(name, "?")
	self.m_dbName = name[begin+1 : end]
}

//! 检查是否含有字段
func (self *DBServer) Query(sql string) bool {
	rows, err := self.m_db.Query(sql)
	if rows != nil {
		defer rows.Close()
	}
	if rows == nil || err != nil {
		log.Println("Query fail! query:", sql, ",err:", err)
		utils.LogError("Query fail! query:", sql, ",err:", err)
		return false
	}

	num := 0
	for rows.Next() {
		num += 1
	}
	return num > 0
}

//! 得到多条数据, 支持tag="-"
func (self *DBServer) GetAllDataEx(query string, struc interface{}) []interface{} {
	rows, err := self.m_db.Query(query)
	if rows != nil {
		defer rows.Close()
	}
	if rows == nil || err != nil {
		log.Println("db GetAllData fail! query:", query, ",err:", err)
		utils.LogError("db GetAllData fail! query:", query, ",err:", err)
		return nil
	}

	//! 得到反射
	v := reflect.ValueOf(struc).Elem()
	data := make([]interface{}, 0)
	for i := 0; i < v.NumField(); i++ {
		// Get the field tag value
		tag := v.Type().Field(i).Tag.Get("json")

		// Skip if tag is not defined or ignored
		if tag == "" || strings.Contains(tag, "-") {
			continue
		}

		ki := v.Field(i).Kind()
		if ki != reflect.Int && ki != reflect.Int64 && ki != reflect.Int8 && ki != reflect.String && ki != reflect.Float32 && ki != reflect.Float64 {
			continue
		}
		data = append(data, v.Field(i).Addr().Interface())
	}

	result := make([]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(data...)
		if err != nil {
			log.Println("db GetAllDataEx-Scan fail! query:", query, ",err:", err)
			utils.LogError("db GetAllDataEx-Scan fail! query:", query, ",err:", err)
			return nil
		}
		newObj := reflect.New(reflect.TypeOf(struc).Elem()).Elem()
		newObj.Set(v)
		result = append(result, newObj.Addr().Interface())
	}

	return result
}

func (self *DBServer) QueryColomn(sql string) (bool, string) {
	rows, err := self.m_db.Query(sql)
	if rows != nil {
		defer rows.Close()
	}
	if rows == nil || err != nil {
		log.Println("Query fail! query:", sql, ",err:", err)
		utils.LogError("Query fail! query:", sql, ",err:", err)
		return false, ""
	}

	var fileName string
	for rows.Next() {
		err = rows.Scan(&fileName)
		if err != nil {
			log.Println("Query fail! query:", sql, ",err:", err)
			utils.LogError("Query fail! query:", sql, ",err:", err)
			return false, ""
		}
		break
	}

	if fileName == "" {
		return false, ""
	}

	return true, fileName
}
