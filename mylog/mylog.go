package mylog

import (
	"fmt"
	toolkit "github.com/charles/toolkit"
)

type MyLogger struct {
	fa *FileAppender
}

var defaultLogger *MyLogger

//func GetDefaultLogger() *MyLogger {
//	return defaultLogger
//}

// init the default
func InitMyLog(prop *toolkit.Properties) error {

	if f, err := NewFileAppender(prop, "log"); err != nil {
		return err
	} else {
		defaultLogger = &MyLogger{
			fa: f,
		}
		return nil
	}
}

func NewMyLogger(prop *toolkit.Properties, loggerName string) (*MyLogger, error) {

	if f, err := NewFileAppender(prop, loggerName); err != nil {
		return nil, err
	} else {
		logger := &MyLogger{
			fa: f,
		}
		return logger, nil
	}
}
func SetLevel(level uint) {
	if defaultLogger != nil {
		defaultLogger.fa.SetLevel(level)
	}
}

func SetTrace(onoff bool) {
	if defaultLogger != nil {
		defaultLogger.fa.trace = onoff
	}
}

func Enable(alarm uint) {
	if defaultLogger != nil {
		defaultLogger.fa.Enable(alarm)
	}
}

func Disable(alarm uint) {
	if defaultLogger != nil {
		defaultLogger.fa.Disable(alarm)
	}
}
func ChanFull() int {
	return defaultLogger.fa.chanFull
}

func Error(format string, v ...interface{}) {
	if defaultLogger == nil {
		fmt.Printf(format+"\n", v...)
	} else {
		defaultLogger.fa.Error(format, v...)
	}
}

func Critical(format string, v ...interface{}) {
	if defaultLogger == nil {
		fmt.Printf(format+"\n", v...)
	} else {
		defaultLogger.fa.Critical(format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if defaultLogger == nil {
		fmt.Printf(format+"\n", v...)
	} else {
		defaultLogger.fa.Info(format, v...)
	}
}

func Trace(format string, v ...interface{}) {
	if defaultLogger == nil {
		fmt.Printf(format+"\n", v...)
	} else {
		defaultLogger.fa.Trace(format, v...)
	}
}

func Debug(format string, v ...interface{}) {

	if defaultLogger == nil {
		fmt.Printf(format+"\n", v...)
	} else {
		defaultLogger.fa.Debug(format, v...)
	}
}

func Warning(format string, v ...interface{}) {
	if defaultLogger == nil {
		fmt.Printf(format+"\n", v...)
	} else {
		defaultLogger.fa.Warning(format, v...)
	}
}

func IsEnabled(lvl uint) bool {

	if defaultLogger == nil {
		return true
	} else {
		return defaultLogger.fa.IsEnabled(lvl)
	}
}

func Notice(format string, v ...interface{}) {
	if defaultLogger == nil {
		fmt.Printf(format+"\n", v...)
	} else {
		defaultLogger.fa.Notice(format, v...)
	}
}

func GetMask() int {
	if defaultLogger == nil {
		return LOG_MASK_ALL
	} else {
		return defaultLogger.fa.GetMask()
	}
}

func Close() {
	if defaultLogger != nil {
		defaultLogger.fa.Close()
	}
}

func GetDefaultLogger() *MyLogger {
	return defaultLogger
}
func (ml *MyLogger) Critical(format string, v ...interface{}) {
	ml.fa.Critical(format, v...)
}
func (ml *MyLogger) Error(format string, v ...interface{}) {
	ml.fa.Error(format, v...)
}

func (ml *MyLogger) Warning(format string, v ...interface{}) {
	ml.fa.Warning(format, v...)
}

func (ml *MyLogger) Debug(format string, v ...interface{}) {
	ml.fa.Debug(format, v...)
}

func (ml *MyLogger) Info(format string, v ...interface{}) {
	ml.fa.Info(format, v...)
}
func (ml *MyLogger) Notice(format string, v ...interface{}) {
	ml.fa.Notice(format, v...)
}

func (ml *MyLogger) Print(v ...interface{}) {
	ml.fa.Print(v...)
}

func (ml *MyLogger) Trace(format string, v ...interface{}) {
	ml.fa.Trace(format, v...)
}

func (ml *MyLogger) Close() {
	ml.fa.Close()
}

func (ml *MyLogger) SetLevel(lvl uint) {
	ml.fa.SetLevel(lvl)
}

func (ml *MyLogger) SetTrace(onoff bool) {
	ml.fa.SetTrace(onoff)
}

func (ml *MyLogger) GetMask() int {
	return ml.fa.GetMask()
}

func (ml *MyLogger) IsEnabled(lvl uint) bool {
	return ml.fa.IsEnabled(lvl)
}
