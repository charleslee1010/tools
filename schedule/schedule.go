package schedule

import (
	"database/sql"
	"errors"
	"fmt"
	log "github.com/charles/mylog"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TimeRule struct {
	yMax, yMin int
	mMax, mMin int
	dMax, dMin int
	wMax, wMin int
	tMax, tMin int
}

type Branch struct {
	allowRules []*TimeRule
	fbRules    []*TimeRule
	appkey     string
}

type Schedule struct {
	m         sync.Mutex
	b         map[string]*Branch
	nBranches int
	nAllow    int
	nFb       int
}

var SCH_ERR_FORMAT error = errors.New("format error")
var SCH_ERR_RANGE error = errors.New("range error")
var SCH_ERR_CONVERT error = errors.New("strconv error")

var SCH_MISMATCH_Y error = errors.New("mismatch year")
var SCH_MISMATCH_M error = errors.New("mismatch month")
var SCH_MISMATCH_D error = errors.New("mismatch day")
var SCH_MISMATCH_W error = errors.New("mismatch weekday")
var SCH_MISMATCH_T error = errors.New("mismatch time")
var SCH_FORBIDDEN error = errors.New("sch forbidden")
var SCH_NOT_CONFIGURED error = errors.New("sch not configured")

var GbSchedule *Schedule = NewSchedule()

func NewSchedule() *Schedule {
	return &Schedule{
		b: make(map[string]*Branch, 0),
	}
}

func (this *TimeRule) match(y, m, d, w, t int) error {
	if (this.yMax == 0 && this.yMin == 0) || (this.yMax >= y && y >= this.yMin) {
	} else {
		return SCH_MISMATCH_Y
	}

	if (this.mMax == 0 && this.mMin == 0) || (this.mMax >= m && m >= this.mMin) {
	} else {
		return SCH_MISMATCH_M
	}

	if (this.dMax == 0 && this.dMin == 0) || (this.dMax >= d && d >= this.dMin) {
	} else {
		return SCH_MISMATCH_D
	}

	if (this.wMax == 0 && this.wMin == 0) || (this.wMax >= w && w >= this.wMin) {
	} else {
		return SCH_MISMATCH_W
	}

	if (this.tMax == 0 && this.tMin == 0) || (this.tMax >= t && t >= this.tMin) {
	} else {
		return SCH_MISMATCH_T
	}
	return nil
}

//
// return ivr string
// rule:   if not match, return error
//
func (this *Branch) match(sid string, y, m, d, w, h int) error {
	var last_err error
	for _, v := range this.allowRules {
		if err := v.match(y, m, d, w, h); err == nil {
			//			log.Debug("sid=%s, match allowRule：%+v, stop", sid, v)
			last_err = nil
			break
		} else {
			//			log.Debug("sid=%s, not match allowRule：%+v, err=%v, continue", sid, v, err)
			last_err = err
		}
	}

	if last_err != nil {
		return last_err
	}

	for _, v := range this.fbRules {
		if err := v.match(y, m, d, w, h); err == nil {
			//			log.Debug("sid=%s, match forbiddenRule：%+v, stop", sid, v)
			last_err = SCH_FORBIDDEN
			break
		} else {
			//			log.Debug("sid=%s, not match forbiddenRule：%+v, err=%v, continue", sid, v, err)
		}
	}

	//	log.Debug("sid=%s, return last err=%v", sid, last_err)

	return last_err
}

//func (this *Schedule) Reset() {
//	this.m.Lock()
//	defer this.m.Unlock()
//
//	this.b = make(map[string]*Branch, 0)
//}

func (this *Schedule) AddTimeRule(appkey string, stype string, y, m, d, w, h string) error {
	this.m.Lock()
	defer this.m.Unlock()

	b := this.findOrCreateBranch(appkey)

	var err error

	if stype == "allow" {
		b.allowRules, err = appendTimeRule(b.allowRules, y, m, d, w, h)
		if err == nil {
			this.nAllow++
		}
	} else {
		b.fbRules, err = appendTimeRule(b.fbRules, y, m, d, w, h)
		if err == nil {
			this.nFb++
		}
	}
	return err
}

func appendTimeRule(s []*TimeRule, y, m, d, w, tstr string) ([]*TimeRule, error) {
	var err error
	t := &TimeRule{}
	if t.yMin, t.yMax, err = convertTime(y, 2020, 2099); err == nil {
		if t.mMin, t.mMax, err = convertTime(m, 1, 12); err == nil {
			if t.dMin, t.dMax, err = convertTime(d, 1, 31); err == nil {
				if t.wMin, t.wMax, err = convertTime(w, 1, 7); err == nil {
					t.tMin, t.tMax, err = convertTime(tstr, 0, 1440)
				}
			}
		}
	}

	if err != nil {
		return s, err
	}

	if s == nil {
		s = make([]*TimeRule, 0)
	}

	return append(s, t), nil
}

func convertTime(s string, min int, max int) (start int, end int, err error) {
	if s == "" {
		return
	}

	tk := strings.Split(s, "-")
	//	log.Info("split produces:%v, length=%d", tk, len(tk))

	switch len(tk) {
	case 2:
		if tk[0] == "" || tk[1] == "" {
			return 0, 0, SCH_ERR_FORMAT
		}

		if start, err = strconv.Atoi(tk[0]); err == nil {
			end, err = strconv.Atoi(tk[1])
		}
	case 1:
		if tk[0] != "" {
			start, err = strconv.Atoi(tk[0])
			end = start
		}

	default:
		return 0, 0, SCH_ERR_FORMAT
	}
	//	log.Info("start=%d, end=%d, err=%v", start, end, err)

	if err != nil {
		log.Error("convert error:%v", err)
		err = SCH_ERR_CONVERT
		return
	}
	if start == 0 && end == 0 {
		return
	}

	if start >= min && start <= max && end >= min && end <= max && end >= start {
		return
	}
	err = SCH_ERR_RANGE
	return
}

func (this *Schedule) findOrCreateBranch(appkey string) *Branch {

	b, pres := this.b[appkey]
	if !pres {
		b = &Branch{
			allowRules: make([]*TimeRule, 0),
			fbRules:    make([]*TimeRule, 0),
			appkey:     appkey,
		}
		this.b[appkey] = b
		this.nBranches++
	}

	return b
}

func (this *Schedule) Match(sid string, ts time.Time, appkey string) error {
	y := ts.Year()
	m := int(ts.Month())
	d := ts.Day()
	w := int(ts.Weekday()) + 1
	h := ts.Hour()*60 + ts.Minute()

	//log.Debug("match %d-%d-%d %d %d appkey=%s", y, m, d, w, h, appkey)
	this.m.Lock()
	defer this.m.Unlock()

	// find appkey schedule
	if b, pres := this.b[appkey]; pres {
		return b.match(sid, y, m, d, w, h)
	} else {
		// use default if appkey specific settting is not present
		if b, pres := this.b["default"]; pres {
			return b.match(sid, y, m, d, w, h)
		}
	}
	// neither default nor app specific setting is present
	log.Debug("appkey not found %s", appkey)
	return SCH_NOT_CONFIGURED
}

func (this *Schedule) LoadCfg(sql string, db *sql.DB) (string, error) {

	sch := NewSchedule()

	if sql == "" {
		sql = "select appkey, y,m,d,w,t,stype from sch_time"
	}

	// load time rule
	rows, err := db.Query(sql)
	if err != nil {
		log.Error("can not query time db, err=%v", err)
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var appkey string
		var y, m, d, w, t, stype string

		if err := rows.Scan(&appkey, &y, &m, &d, &w, &t, &stype); err != nil {
			log.Error("can not scan rows, err=%v", err)
			return "", err
		}

		if err := sch.AddTimeRule(appkey, stype, y, m, d, w, t); err != nil {
			log.Error("invalid time rule, err=%v", err)
			return "", err
		}
	}

	if err := rows.Err(); err != nil {
		log.Error("rows.err, err=%v", err)
		return "", err
	}

	// we have data in sch, replace this with sch
	this.m.Lock()
	this.b = sch.b
	this.nAllow = sch.nAllow
	this.nFb = sch.nFb
	this.nBranches = sch.nBranches
	this.m.Unlock()

	info := fmt.Sprintf("nBranches=%d, nAllow=%d, nFb=%d",
		this.nBranches, this.nAllow, this.nFb)

	return info, nil
}

func (this *Schedule) Dump() {
	log.Info("dump schedule")
	log.Info("nBranch=%d, nAllow=%d, nFb=%d", this.nBranches, this.nAllow, this.nFb)

	for _, b := range this.b {
		log.Info("appkey[%s]", b.appkey)
		for _, v := range b.allowRules {
			log.Info("allow rule: %+v", v)
		}
		for _, v := range b.fbRules {
			log.Info("forbidden rule: %+v", v)
		}
	}
	log.Info("dump schedule ends")

}
