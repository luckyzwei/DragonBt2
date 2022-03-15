package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type CsvNode map[string]string
type CsvAryNode []string

var csvmgrsingleton *CsvMgr = nil

//! public
func GetRobotCsvMgr() *CsvMgr {
	if csvmgrsingleton == nil {
		csvmgrsingleton = new(CsvMgr)
		csvmgrsingleton.InitStruct()
	}

	return csvmgrsingleton
}

type CsvMgr struct {
	Data        map[string]map[int]CsvNode
	Data2       map[string][]CsvNode
	DataAry     map[string][]CsvAryNode
	HeadAry     map[string][]string
	normal_data map[string]map[int]CsvNode
}

func (self *CsvMgr) InitStruct() {
	self.Data = make(map[string]map[int]CsvNode)
	self.Data2 = make(map[string][]CsvNode)
	self.DataAry = make(map[string][]CsvAryNode)
	self.HeadAry = make(map[string][]string)
	self.normal_data = make(map[string]map[int]CsvNode)
}

func (self *CsvMgr) InitData() {
	if IS_ROBOT {
		self.ReadData("robotnormal")
		self.ReadData2("robotfight")
	} else {
		self.ReadData2("changename")
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}

		dir_list, e := ioutil.ReadDir(dir)
		if e != nil {
			fmt.Println("read dir error")
			return
		}
		for i, v := range dir_list {
			if strings.Contains(v.Name(), "changename") == false && strings.Contains(v.Name(), ".csv") == true {
				self.ReadCsvData(v.Name())
				self.WritCsvData(v.Name())
				fmt.Println(i, "go over", v.Name())
			}

		}
		fmt.Println("over all")
	}
}

func (self *CsvMgr) Reload() {
	self.InitStruct()
	self.InitData()
}

func (self *CsvMgr) ReadData(name string) {
	_, ok := self.Data[name]
	if ok {
		log.Println("重复读入csv:", name)
		return
	}

	file, err := os.Open(name + ".csv")
	if err != nil {
		log.Fatalln("csv err1:", name, err)
		return
	}
	defer file.Close()

	header := make([]string, 0)
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln("csv err2:", name, err)
			return
		}

		if len(header) == 0 {
			header = record
		} else {
			id, err := strconv.Atoi(record[0])
			if err != nil {
				log.Fatalln("csv err3:", name, err)
				return
			}

			_, ok := self.Data[name]
			if !ok {
				self.Data[name] = make(map[int]CsvNode)
			}

			_, ok = self.Data[name][id]
			if !ok {
				self.Data[name][id] = make(CsvNode)
			}

			for i := 0; i < len(record); i++ {
				self.Data[name][id][header[i]] = record[i]
			}
		}
	}
}

func (self *CsvMgr) ReadCsvData(name string) {
	_, ok := self.DataAry[name]
	if ok {
		log.Println("重复读入csv:", name)
		return
	}

	if strings.Contains(name, ".csv") == false {
		name = name + ".csv"
	}

	file, err := os.Open(name)
	if err != nil {
		log.Fatalln("ReadCsvData err1:", name, err)
		return
	}
	defer file.Close()

	header := make([]string, 0)
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln("ReadCsvData err2:", name, err)
			return
		}
		//log.Println(record)
		if len(header) == 0 {
			header = record
			self.HeadAry[name] = header
		} else {
			_, err := strconv.Atoi(record[0])
			if err != nil {
				log.Fatalln("ReadCsvData err3:", name, err)
				return
			}

			_, ok := self.DataAry[name]
			if !ok {
				self.DataAry[name] = make([]CsvAryNode, 0)
			}

			var curary CsvAryNode
			curary = make(CsvAryNode, 0)

			for i := 0; i < len(record); i++ {
				curary = append(curary, record[i])
			}

			self.DataAry[name] = append(self.DataAry[name], curary)
		}
	}
}

func (self *CsvMgr) WritCsvData(name string) {

	if strings.Contains(name, ".csv") == false {
		name = name + ".csv"
	}

	file, err := os.OpenFile("outcsv/"+name, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("csv err1:", name, err)
		return
	}
	defer file.Close()
	nw := csv.NewWriter(file)
	var csv_data []CsvAryNode
	csv_data = self.DataAry[name]
	nw.Write(self.HeadAry[name])
	for i := 0; i < len(csv_data); i++ {
		wheader := make([]string, 0)
		for j := 0; j < len(csv_data[i]); j++ {
			csv_data[i][j] = self.getNewString(csv_data[i][j])
			wheader = append(wheader, csv_data[i][j])
		}
		nw.Write(wheader)
	}
	nw.Flush()
}

func (self *CsvMgr) ReadData2(name string) {
	_, ok := self.Data2[name]
	if ok {
		log.Println("重复读入csv:", name)
		return
	}

	file, err := os.Open(name + ".csv")
	if err != nil {
		log.Fatalln("csv err4:", name, err)
		return
	}
	defer file.Close()

	header := make([]string, 0)
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln("csv err5:", name, err)
			return
		}

		//log.Println(record)

		if len(header) == 0 {
			header = record
		} else {
			node := make(CsvNode)
			for i := 0; i < len(record); i++ {
				node[header[i]] = record[i]
			}
			self.Data2[name] = append(self.Data2[name], node)
		}
	}
}

func (self *CsvMgr) GetData2Int(table string, id int, field string) int {
	value, _ := strconv.Atoi(self.Data[table][id][field])
	return value
}

func (self *CsvMgr) GetData2String(table string, id int, field string) string {
	return self.Data[table][id][field]
}

func (self *CsvMgr) getNewString(curstr string) string {
	data := self.Data2["changename"]
	for _, value := range data {

		if curstr == value["curname"] {
			return value["newname"]
		} else if strings.Contains(curstr, value["curname"]) == true {
			return strings.Replace(curstr, value["curname"], value["newname"], -1)
		}
	}
	return curstr
}
