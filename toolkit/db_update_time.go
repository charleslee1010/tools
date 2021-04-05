package toolkit

import (
    "time"
)

type DBUpdateTime struct {
    // time when load app configuration data
	LastUpdateTime string
}



func NewDBUpdateTime() (*DBUpdateTime) {
	return &DBUpdateTime {
		 LastUpdateTime: "2018-04-23 12:24:51",
	}
}

func (this*DBUpdateTime) BuildSql(sql string) string {
	// form sql string
	sql = sql + ` where updatetime > "` + this.LastUpdateTime + `"`
	return sql
}

func (this*DBUpdateTime) Update() {
	// record the current time
	this.LastUpdateTime = time.Now().Format("2006-01-02 15:04:05")
}