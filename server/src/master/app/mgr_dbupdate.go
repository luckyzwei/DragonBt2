package app

import (
	"errors"
	"fmt"
	"master/db"
	"master/utils"
	"os"
	"strings"
)

// 数据库字段检查
var s_mgrdbupdate *DBUpdateMgr = nil

type DBUpdateMgr struct {
	tableCheck            []string    //! 检查表是否存在，不存在则Fatal
	logTableCheck         []string    //! 日志表检查，不存在则Fatal
	fieldCheck            [][2]string //! 表的字段检查
	createTableStatements []string    //! 直接创建表，不存在创建
	addField              [][3]string //! 增加字段
	modifyField           [][4]string
	dropField             [][3]string
}

func (self *DBUpdateMgr) initData() {
	self.tableCheck = []string{}

	self.fieldCheck = [][2]string{}

	self.logTableCheck = []string{}

	self.addField = [][3]string{}

	self.modifyField = [][4]string{}

	self.createTableStatements = []string{}

	self.dropField = [][3]string{}
}

func (self *DBUpdateMgr) createTables() {
	for _, stmt := range self.createTableStatements {
		_, _, res := db.GetDBMgr().DBUser.Exec(stmt)
		if !res {
			utils.LogError("创建table失败")
			os.Exit(1)
		}
	}
}

// 检查字段是否存在
func (self *DBUpdateMgr) CheckMysql() {
	self.initData()
	var checkErr error
	for _, filedName := range self.fieldCheck {
		checkErr = self.CheckFiled(filedName[0], filedName[1])
		if checkErr != nil {
			utils.LogError(checkErr.Error())
			os.Exit(1)
		}
	}

	for index := range self.tableCheck {
		checkErr = self.CheckTable(self.tableCheck[index])
		if checkErr != nil {
			utils.LogError(checkErr.Error())
			os.Exit(1)
		}
	}

	for index := range self.logTableCheck {
		checkErr = self.CheckLogTable(self.logTableCheck[index])
		if checkErr != nil {
			utils.LogError(checkErr.Error())
			os.Exit(1)
		}
	}

	self.createTables()
	self.checkAddField()
	self.modifyColumnType()
	self.checkDropField()
}

func (self *DBUpdateMgr) CheckFiled(tableName string, filedName string) error {
	sql := "SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'  AND COLUMN_NAME = '%s'"
	dbName := self.getDbName()
	if dbName == "" {
		utils.LogDebug("dbName is empty!")
		return nil
	}
	sqlStr := fmt.Sprintf(sql, dbName, tableName, filedName)
	res := db.GetDBMgr().DBUser.Query(sqlStr)
	//LogDebug("sqlStr:", sqlStr, ", CheckFiledParam:", res)
	if res {
		return nil
	}
	return errors.New(tableName + " has no filed:" + filedName)
}

func (self *DBUpdateMgr) CheckTable(tableName string) error {
	sql := "SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'"
	dbName := self.getDbName()
	if dbName == "" {
		utils.LogDebug("dbName is empty!")
		return nil
	}
	sqlStr := fmt.Sprintf(sql, dbName, tableName)
	res := db.GetDBMgr().DBUser.Query(sqlStr)
	//LogDebug("sqlStr:", sqlStr, ", CheckFiledParam:", res)
	if res {
		return nil
	}
	return errors.New(tableName + " not exists!")
}

func (self *DBUpdateMgr) getDbName() string {
	sqlSplit1 := strings.Split(GetMasterApp().Conf.DBConf.DBUser, "?")
	if len(sqlSplit1) < 1 {
		return ""
	}
	sqlSplit2 := strings.Split(sqlSplit1[0], "/")
	if len(sqlSplit2) < 2 {
		return ""
	}
	return sqlSplit2[1]
}



func (self *DBUpdateMgr) getDbLogName() string {
	sqlSplit1 := strings.Split(GetMasterApp().Conf.DBConf.DBLog, "?")
	if len(sqlSplit1) < 1 {
		return ""
	}
	sqlSplit2 := strings.Split(sqlSplit1[0], "/")
	if len(sqlSplit2) < 2 {
		return ""
	}
	return sqlSplit2[1]
}

func (self *DBUpdateMgr) CheckLogTable(tableName string) error {
	sql := "SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'"
	dbName := self.getDbLogName()
	if dbName == "" {
		utils.LogDebug("dbName is empty!")
		return nil
	}
	sqlStr := fmt.Sprintf(sql, dbName, tableName)
	res := db.GetDBMgr().DBLog.Query(sqlStr)
	if res {
		return nil
	}
	return errors.New(tableName + " not exists!")
}

func (self *DBUpdateMgr) checkAddField() {
	var checkErr error
	for _, stmt := range self.addField {
		checkErr = self.CheckFiled(stmt[0], stmt[1])
		if checkErr != nil { // 没有才插入
			_, _, res := db.GetDBMgr().DBUser.Exec(stmt[2])
			if !res {
				utils.LogError("增加字段失败", stmt[2])
				os.Exit(1)
			}
		}
	}

}

func (self *DBUpdateMgr) checkDropField() {
	var checkErr error
	for _, stmt := range self.dropField {
		checkErr = self.CheckFiled(stmt[0], stmt[1])
		if checkErr == nil { // 没有才插入
			_, _, res := db.GetDBMgr().DBUser.Exec(stmt[2])
			if !res {
				utils.LogError("删除字段失败", stmt[2])
				os.Exit(1)
			}
		}
	}

}

// ALTER TABLE `san_usertreasure` MODIFY COLUMN `info`  mediumtext CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '宝物信息' AFTER `washfreetimes`;
func (self *DBUpdateMgr) modifyColumnType() {
	for _, stmt := range self.modifyField {
		if len(stmt) != 4 {
			continue
		}

		isStmtOk := true
		for _, v := range stmt {
			if v == "" {
				isStmtOk = false
				break
			}
		}

		if !isStmtOk {
			break
		}

		sql := "SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'  AND COLUMN_NAME = '%s'"
		dbName := self.getDbName()
		if dbName == "" {
			utils.LogError("dbName is empty!")
			os.Exit(1)
		}

		sqlStr := fmt.Sprintf(sql, dbName, stmt[0], stmt[1])
		res, fieldName := db.GetDBMgr().DBUser.QueryColomn(sqlStr)
		if res && fieldName == stmt[2] {
			_, _, res := db.GetDBMgr().DBUser.Exec(stmt[3])
			if !res {
				utils.LogError("修改字段类型失败", stmt[3])
				os.Exit(1)
			} else {
				utils.LogDebug("修改字段类型成功!")
			}
		}
	}

}
