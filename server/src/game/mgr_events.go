package game

import (
	"strings"
)

var eventsMgr *EventsMgr

func GetEventsMgr() *EventsMgr {
	if eventsMgr == nil {
		eventsMgr = new(EventsMgr)
	}
	return eventsMgr
}

// 事件表
type LevelmapEventConfig struct {
	Id        int      `json:"id"`
	Nextid    int      `json:"nextid"`
	Ns        []string `json:"n"`
	Record    int      `json:"record"`
	EventType int      `json:"toucheffect"`
}

// 关卡事件表
type LevelmapThingConfig struct {
	Id         int   `json:"id"`
	Thingtype  int   `json:"thingtype"`
	Touchtype  int   `json:"touchtype"`
	Conditions []int `json:"condition"` // 填的事件Id
	Events     []int `json:"event"`
}

// 关卡章节表
type LevelmapChapterConfig struct {
	Chapterindex int    `json:"chapterindex"`
	MapId        string `json:"mapid"`
}

type LevelmapBoxConfig struct {
	Chapterid int   `json:"chapterid"`
	Items     []int `json:"item"`
	Nums      []int `json:"num"`
	Explore   int   `json:"explore"`
}

type LevelmapMailConfig struct {
	Id          int `json:"id"`
	Chapter     int `json:"chapter"`
	Thingid     int `json:"thingid"`
	Explore     int `json:"explore"`
	ChapterType int `json:"type"`
}

type EventsMgr struct {
	LevelmapEventConfig   map[int]*LevelmapEventConfig
	LevelmapChapterConfig []*LevelmapChapterConfig
	ThingChapterMap       map[int]int
	LevelmapBoxConfig     map[int]*LevelmapBoxConfig
	LevelmapMailConfig    []*LevelmapMailConfig
	LevelThingMap         map[int]map[int]*LevelmapThingConfig // 事件信息, chapter, thingId :  config
	TrigerThingMap        map[int][]*ThingEvent
	FileNames             []string
}

func (self *EventsMgr) LoadConfig() {
	self.LevelmapEventConfig = make(map[int]*LevelmapEventConfig)
	GetCsvUtilMgr().LoadCsv("Levelmap_Event", &self.LevelmapEventConfig)
	//litter.Dump(self.LevelmapEventConfig)
	GetCsvUtilMgr().LoadCsv("Levelmap_Chapter", &self.LevelmapChapterConfig)
	self.MakeFileNames()
	//litter.Dump(self.LevelmapChapterConfig)

	//litter.Dump(self.ThingChapterMap)
	self.LevelmapBoxConfig = make(map[int]*LevelmapBoxConfig)
	GetCsvUtilMgr().LoadCsv("Levelmap_Box", &self.LevelmapBoxConfig)
	GetCsvUtilMgr().LoadCsv("Levelmap_Mail", &self.LevelmapMailConfig)
	self.LoadThing()
}

func (self *EventsMgr) MakeFileNames() {
	for _, v := range self.LevelmapChapterConfig {
		self.FileNames = append(self.FileNames, strings.Trim(v.MapId, " "))
	}
}

type ThingEvent struct {
	ThingId int `json:"thing_id"`
	EventId int `json:"event_id"`
}

// 加载物件数据
func (self *EventsMgr) LoadThing() {
	// Map1001_Item.csv
	self.LevelThingMap = make(map[int]map[int]*LevelmapThingConfig)
	for i, fn := range self.FileNames {
		var info = make(map[int]*LevelmapThingConfig)
		if fn == "" {
			continue
		}
		GetCsvUtilMgr().LoadEventsCsv(fn, &info)
		//litter.Dump(info)
		self.LevelThingMap[i+1] = info
	}
	self.TrigerThingMap = make(map[int][]*ThingEvent)
	self.ThingChapterMap = make(map[int]int)
	for chapter, things := range self.LevelThingMap {
		for _, v2 := range things {
			for _, cond := range v2.Conditions {
				if cond == 0 {
					continue
				}
				thingId := cond / 100000
				eventId := cond % 100000
				self.TrigerThingMap[v2.Id] = append(self.TrigerThingMap[v2.Id], &ThingEvent{thingId, eventId})
			}
			self.ThingChapterMap[v2.Id] = chapter
		}
	}

	//litter.Dump(self.TrigerThingMap)
	//litter.Dump(self.LevelThingMap)
}

func (self *EventsMgr) getChapterThings(chapter int) map[int]*LevelmapThingConfig {
	data, ok := self.LevelThingMap[chapter]
	if !ok {
		return nil
	}

	return data
}

// 获取某个事件的前置条件
func (self *EventsMgr) getTrigger(thingId int) []*ThingEvent {
	data, ok := self.TrigerThingMap[thingId]
	if !ok {
		return nil
	}
	return data
}

// 获取事件配置
func (self *EventsMgr) GetThingConfig(chapter int, thingId int) *LevelmapThingConfig {
	chapterData, ok := self.LevelThingMap[chapter]
	if !ok {
		LogError("章节数据不存在, chapter = ", chapter)
		return nil
	}

	info, ok := chapterData[thingId]
	if !ok {
		LogError("章节数据不存在, 物件Id = ", thingId)
		return nil
	}

	return info
}
