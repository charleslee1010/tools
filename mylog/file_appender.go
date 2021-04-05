package mylog

import (
	//	"log"
	"fmt"
	toolkit "github.com/charles/toolkit"
	"os"
	"time"

	"strings"
	//	"runtime"
	"errors"
	"path"
	//	"strconv"
)

// log severity
const (
	DEBUG = iota
	NOTICE
	INFO
	WARNING
	ERROR
	CRITICAL
	TRACE
	FA_EXIT // indicate the log thread to exit
)

const (
	MASK_DEBUG = 1
	MASK_NOTICE = 2
	MASK_INFO = 4
	MASK_WARNING = 8
	MASK_ERROR = 16
	MASK_CRITICAL = 32
	MASK_TRACE = 64
)

var severityString []string = []string{
	"DEBUG",
	"NOTICE",
	"INFO",
	"WARNING",
	"ERROR",
	"CRIT",
	"TRACE",
}

// except DEBUG
const LOG_MASK_DEFAULT int = 126  
const LOG_MASK_ALL int = 127

type FileAppender struct {
	logFilePath    string
	logFile        *os.File
	rotateSize     int64
	rotateInterval int64
	rotateNumber   int
	nextIndex      int
	logChan        chan *LogItem

	logFormat       string
	fileNamePattern string

	logTimer *time.Timer

	// severity level
//	level int
	
	logMask     int
	// trace on/off
	trace bool

	elogFormat string
	wlogFormat string
	ilogFormat string
	dlogFormat string
	tlogFormat string
	nlogFormat string

	chanFull int
}

type LogItem struct {
	t0       time.Time
	severity uint
	info     string

	file string
	line string
}

const (
	MIN_ROTATE_INTERVAL = 5
	LOG_CHANN_SIZE      = 5000
)

//
func NewFileAppender(prop *toolkit.Properties, logger string) (*FileAppender, error) {

	fa := &FileAppender{
		logFilePath:     prop.GetPropertyString(logger+".filePath", "stderr"),
		rotateSize:      int64(prop.GetPropertyInt(logger+".maxSize", 100)) * 1e6,
		rotateInterval:  int64(prop.GetPropertyInt(logger+".rotateInterval", 24*60)),
		rotateNumber:    prop.GetPropertyInt(logger+".maxNumber", 0),
		logChan:         make(chan *LogItem, LOG_CHANN_SIZE),
		logFormat:       prop.GetPropertyString(logger+".logFormat", ""),
		fileNamePattern: prop.GetPropertyString(logger+".fileNamePattern", "default.%Y-%M-%D_%h-%m-%s.log"),
//		level:           INFO,
		trace:           true,
		logMask:         LOG_MASK_DEFAULT,

		elogFormat: prop.GetPropertyString(logger+".ErrorLogFormat", ""),
		wlogFormat: prop.GetPropertyString(logger+".WarningLogFormat", ""),
		ilogFormat: prop.GetPropertyString(logger+".InfoLogFormat", ""),
		tlogFormat: prop.GetPropertyString(logger+".TraceLogFormat", ""),
		nlogFormat: prop.GetPropertyString(logger+".NoticeLogFormat", ""),
		dlogFormat: prop.GetPropertyString(logger+".DebugLogFormat", ""),
	}

	// rotateInterval must be multiple of 10 minutes
	fa.rotateInterval = (fa.rotateInterval + MIN_ROTATE_INTERVAL - 1) / MIN_ROTATE_INTERVAL * MIN_ROTATE_INTERVAL * 60

	if fa.rotateNumber == 0 && fa.fileNamePattern == "" {
		return nil, errors.New(logger + " [maxNumber or fileNamePattern] is not configured properly")
	}

	if f, err := openLogFile(fa.logFilePath); err != nil {
		return nil, err
	} else {
		fa.logFile = f
		if f == os.Stdout || f == os.Stderr {
			fa.rotateInterval = 0
			fa.rotateNumber = 0
		}
	}

	// find the next index to be used
	if fa.rotateNumber > 0 {
		lidx := fa.findLatestIndex()
		fa.nextIndex = fa.next(lidx)
	}

	if fa.rotateInterval > 0 {
		fa.startLogTimer()
		//		fmt.Info("log switch interval is set to %d", fa.rotateInterval)
	} else {
		// make a chann which never expires
		fa.logTimer = &time.Timer{
			C: make(chan time.Time, 1),
		}
		//		fmt.Info("log switch interval is not set")
	}

	go fa.logWriter()

	return fa, nil
}

func (fa *FileAppender) GetMask() int {
	return fa.logMask
}

func (fa *FileAppender) SetTrace(onoff bool) {
	if onoff {
		fa.logMask = fa.logMask|MASK_TRACE
	} else {
		fa.logMask = fa.logMask&(LOG_MASK_ALL^MASK_TRACE)
	}
}

func (fa *FileAppender) SetLevel(lvl uint) {
	fa.logMask = fa.logMask|(1 << lvl)
}

func (fa *FileAppender) Enable(lvl uint) {
	fa.logMask = fa.logMask|(1 << lvl)
}

func (fa *FileAppender) Disable(lvl uint) {
	fa.logMask = fa.logMask&(LOG_MASK_ALL^(1<<lvl))
}

func (fa *FileAppender) IsEnabled(lvl uint) bool {
	 return (fa.logMask|(1 << lvl)) > 0
}

func (fa *FileAppender) ChanFull() int {
	return fa.chanFull
}

func (fa *FileAppender) startLogTimer() {
	// start timer
	// time in seconds of midnight
	t0 := time.Now()
	t := toolkit.SecondsFromMidnight(t0)

	d := fa.rotateInterval - t%fa.rotateInterval
	if d <= 5 {
		d += fa.rotateInterval
	}
	fa.logTimer = time.NewTimer(time.Duration(d) * time.Second)

	//		fmt.Println("[Info]" + toolkit.GetTimeStamp(t0) + fmt.Sprintf("switch after %d seconds\n", d))

}

func (fa *FileAppender) findLatestIndex() int {
	// find the latest log file
	var ts int64
	ts = 0
	latest := -1

	for i := 0; i < fa.rotateNumber; i++ {
		f := fmt.Sprintf("%s.%d", fa.logFilePath, i)
		if finfo, err := os.Stat(f); err == nil {
			fts := finfo.ModTime().UnixNano()
			if fts > ts {
				ts = fts
				latest = i
			}
		}
	}

	return latest
}

func openLogFile(path string) (*os.File, error) {

	s := strings.ToLower(path)
	if s == "stderr" {
		return os.Stderr, nil
	}

	if s == "stdout" {
		return os.Stdout, nil
	}

	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

}

func (fa *FileAppender) logWriter() {
	for true {
		select {

		case log := <-fa.logChan:
			fa.writeLogItem(log)
			// switch by size
			if fa.logFile != nil && fa.logFile != os.Stderr && fa.logFile != os.Stdout {
				if fi, err := fa.logFile.Stat(); err == nil {
					if fa.rotateSize > 0 && fi.Size() >= fa.rotateSize {
						fmt.Println("size = ", fi.Size())
						fa.switchLogFile()
					}
				}
			}
		case <-fa.logTimer.C:
			// switch file, start new ticker
			fa.switchLogFile()
			fa.startLogTimer()
		}
	}

	fa.logFile.Close()
}

func (fa *FileAppender) getLogFormat(sev uint) (logf string) {
	logf = fa.logFormat

	if sev == ERROR && fa.elogFormat != "" {
		logf = fa.elogFormat
	} else if sev == WARNING && fa.wlogFormat != "" {
		logf = fa.wlogFormat
	} else if sev == INFO && fa.ilogFormat != "" {
		logf = fa.ilogFormat
	} else if sev == DEBUG && fa.dlogFormat != "" {
		logf = fa.dlogFormat
	} else if sev == NOTICE && fa.nlogFormat != "" {
		logf = fa.nlogFormat
	} else if sev == TRACE && fa.tlogFormat != "" {
		logf = fa.tlogFormat
	}
	return
}

func getLogSeverityString(sev uint) string {
	if sev >= 0 && sev <= TRACE {
		return severityString[sev]
	}
	return ""
}

// write a log item by logFormat
//
func (fa *FileAppender) writeLogItem(log *LogItem) {

	ts := toolkit.GetTimeStamp(log.t0)

	if fa.logFile == nil {
		// output to console
		fmt.Println(ts, log.severity, log.file, log.line, log.info)
		return
	}

	if fa.logFormat == "" {
		// no logformat specified
		fa.logFile.WriteString(log.file)
		fa.logFile.WriteString(" ")
		fa.logFile.WriteString(log.line)
		fa.logFile.WriteString(" ")
		fa.logFile.WriteString(log.info)
		fa.logFile.WriteString("\n")
		return
	}

	s := fa.getLogFormat(log.severity)

	for true {
		idx := strings.Index(s, "%")
		if idx > 0 {
			fa.logFile.WriteString(s[0:idx])
			s = s[idx:]
		} else if idx == 0 {
			if strings.Index(s, "%V") == 0 {
				fa.logFile.WriteString(getLogSeverityString(log.severity))
				s = s[2:]
			} else if strings.Index(s, "%T") == 0 {
				fa.logFile.WriteString(ts)
				s = s[2:]
			} else if strings.Index(s, "%I") == 0 {
				fa.logFile.WriteString(log.info)
				s = s[2:]
			} else if strings.Index(s, "%F") == 0 {
				fa.logFile.WriteString(path.Base(log.file))
				s = s[2:]
			} else if strings.Index(s, "%L") == 0 {
				fa.logFile.WriteString(log.line)
				s = s[2:]
			} else {
				// skip the %
				s = s[1:]
			}
		} else {
			// idx < 0
			fa.logFile.WriteString(s)
			break
		}
	}

	fa.logFile.WriteString("\n")
}

func (fa *FileAppender) switchLogFile() {
	t0 := time.Now()
	fmt.Println("[Info]" + toolkit.GetTimeStamp(t0) + " log file is closed\n")
	fa.logFile.Close()

	// get new log file name
	newFilePath := fa.getNewFilePath(t0)
	fmt.Println("new file path=", newFilePath)
	// rename to new file
	if err := os.Rename(fa.logFilePath, newFilePath); err != nil {
		fmt.Println("can not rename file", fa.logFilePath, newFilePath, err.Error())
	} else {
		fmt.Println("rename file", fa.logFilePath, " to:", newFilePath)
	}
	//	fmt.Println("log file ", newpath, "is generated")
	// reopen new file
	fa.logFile, _ = openLogFile(fa.logFilePath)
	fmt.Println("[Info]" + toolkit.GetTimeStamp(t0) + " log file is created\n")

}

func (fa *FileAppender) getNewFilePath(t time.Time) string {
	var newpath string
	if fa.rotateNumber > 0 {
		newpath = fmt.Sprintf("%s.%d", fa.logFilePath, fa.nextIndex)
		fa.nextIndex = fa.next(fa.nextIndex)
	} else {
		// rotateNumber = 0, filename is generated according to fileNamePattern
		dir := path.Dir(fa.logFilePath)
		fn := toolkit.FormatFileName(fa.fileNamePattern, t)
		fmt.Println("filepath=", fa.logFilePath, ",dir=", dir, ", fn=", fn)
		newpath = path.Join(dir, fn)
	}

	os.Remove(newpath)
	return newpath
}

func (fa *FileAppender) next(idx int) int {
	idx++
	if idx >= fa.rotateNumber {
		return 0
	}
	return idx
}

func (fa *FileAppender) Error(format string, v ...interface{}) {
	if (MASK_ERROR & fa.logMask) == 0 {
		return
	}

	//	_, file, lno, _ := runtime.Caller(2)
	//	line := strconv.Itoa(lno)

	file := ""
	line := ""
	if len(fa.logChan) >= LOG_CHANN_SIZE {
		fa.chanFull++
		return
	}

	s := fmt.Sprintf(format, v...)
	fa.logChan <- &LogItem{time.Now(), ERROR, s, file, line}
}

func (fa *FileAppender) Print(v ...interface{}) {
	if (MASK_INFO & fa.logMask) == 0 {
		return
	}

	file := ""
	line := ""
	if len(fa.logChan) >= LOG_CHANN_SIZE {
		fa.chanFull++
		return
	}
	s := fmt.Sprint(v...)
	fa.logChan <- &LogItem{time.Now(), INFO, s, file, line}
}

func (fa *FileAppender) Info(format string, v ...interface{}) {
	if (MASK_INFO & fa.logMask) == 0 {
		return
	}

	//	_, file, lno, _ := runtime.Caller(2)
	//	line := strconv.Itoa(lno)

	file := ""
	line := ""
	if len(fa.logChan) >= LOG_CHANN_SIZE {
		fa.chanFull++
		return
	}
	s := fmt.Sprintf(format, v...)
	fa.logChan <- &LogItem{time.Now(), INFO, s, file, line}
}

func (fa *FileAppender) Trace(format string, v ...interface{}) {
	if (MASK_TRACE & fa.logMask) == 0 {
		return
	}
	//	_, file, lno, _ := runtime.Caller(2)
	//	line := strconv.Itoa(lno)

	file := ""
	line := ""
	if len(fa.logChan) >= LOG_CHANN_SIZE {
		fa.chanFull++
		return
	}

	s := fmt.Sprintf(format, v...)
	fa.logChan <- &LogItem{time.Now(), TRACE, s, file, line}
}

func (fa *FileAppender) Critical(format string, v ...interface{}) {

	file := ""
	line := ""
	if len(fa.logChan) >= LOG_CHANN_SIZE {
		fa.chanFull++
		return
	}

	s := fmt.Sprintf(format, v...)
	fa.logChan <- &LogItem{time.Now(), CRITICAL, s, file, line}
}

func (fa *FileAppender) Debug(format string, v ...interface{}) {
	if (MASK_DEBUG & fa.logMask) == 0 {
		return
	}
	//	_, file, lno, _ := runtime.Caller(2)
	//	line := strconv.Itoa(lno)
	file := ""
	line := ""
	if len(fa.logChan) >= LOG_CHANN_SIZE {
		fa.chanFull++
		return
	}

	s := fmt.Sprintf(format, v...)
	fa.logChan <- &LogItem{time.Now(), DEBUG, s, file, line}
}

func (fa *FileAppender) Warning(format string, v ...interface{}) {
	if (MASK_WARNING & fa.logMask) == 0 {
		return
	}
	//	_, file, lno, _ := runtime.Caller(2)
	//	line := strconv.Itoa(lno)

	file := ""
	line := ""
	if len(fa.logChan) >= LOG_CHANN_SIZE {
		fa.chanFull++
		return
	}

	s := fmt.Sprintf(format, v...)
	fa.logChan <- &LogItem{time.Now(), WARNING, s, file, line}
}

func (fa *FileAppender) Notice(format string, v ...interface{}) {
	if (MASK_NOTICE & fa.logMask) == 0 {
		return
	}
	//	_, file, lno, _ := runtime.Caller(2)
	//	line := strconv.Itoa(lno)
	file := ""
	line := ""
	if len(fa.logChan) >= LOG_CHANN_SIZE {
		fa.chanFull++
		return
	}

	s := fmt.Sprintf(format, v...)
	fa.logChan <- &LogItem{time.Now(), NOTICE, s, file, line}
}
func (fa *FileAppender) Close() {
	fa.logChan <- &LogItem{time.Now(), FA_EXIT, "", "", ""}
}
