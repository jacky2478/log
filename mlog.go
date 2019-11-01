//Simple and efficient micro mlog library that supports features such as condition, color, and file split
package mlog

import (
	"io"
	"log"
	"runtime"

	"github.com/jacky2478/color"
)

const (
	Ldate         = log.Ldate
	Llongfile     = log.Llongfile
	Lmicroseconds = log.Lmicroseconds
	Lshortfile    = log.Lshortfile
	LstdFlags     = log.LstdFlags
	Ltime         = log.Ltime
)

type (
	LogLevel int
	LogType  int
)

const (
	LOG_FATAL   = LogType(0x1)
	LOG_ERROR   = LogType(0x2)
	LOG_WARNING = LogType(0x4)
	LOG_INFO    = LogType(0x8)
	LOG_DEBUG   = LogType(0x10)
)

const (
	LOG_LEVEL_NONE  = LogLevel(0x0)
	LOG_LEVEL_FATAL = LOG_LEVEL_NONE | LogLevel(LOG_FATAL)
	LOG_LEVEL_ERROR = LOG_LEVEL_FATAL | LogLevel(LOG_ERROR)
	LOG_LEVEL_WARN  = LOG_LEVEL_ERROR | LogLevel(LOG_WARNING)
	LOG_LEVEL_INFO  = LOG_LEVEL_WARN | LogLevel(LOG_INFO)
	LOG_LEVEL_DEBUG = LOG_LEVEL_INFO | LogLevel(LOG_DEBUG)
	LOG_LEVEL_ALL   = LOG_LEVEL_DEBUG
)

const FORMAT_TIME_DAY string = "20060102"
const FORMAT_TIME_HOUR string = "2006010215"

var _log *logger = New()
var _colorMap = make(map[string]func(string, ...interface{}) string, 0)
var _logMap = make(map[string]*logger, 0)

func init() {
	SetFlags(Ldate | Ltime | Lshortfile)
	SetHighlighting(runtime.GOOS != "windows")

	// support log with color
	SetColorByLogType(LOG_INFO, color.GreenString)
	SetColorByLogType(LOG_ERROR, color.RedString)
	SetColorByLogType(LOG_WARNING, color.YellowString)
	SetColorByLogType(LOG_FATAL, color.RedString)
	SetColorByLogType(LOG_DEBUG, color.BlueString)
}

func Logger() *log.Logger {
	return _log._log
}

// support regist logger by name
func RegistLog(key string, ptrLog *logger) {
	if key != "" && ptrLog != nil {
		_logMap[key] = ptrLog
	}
}

// get the regist logger by name
func GetRegistLog(key string) *logger {
	if ptrLog, bok := _logMap[key]; bok {
		return ptrLog
	}
	return _log
}

func SetLevel(level LogLevel) {
	_log.SetLevel(level)
}
func GetLogLevel() LogLevel {
	return _log.level
}

func SetDepth(depth int) {
	_log.SetDepth(depth)
}

func SetOutput(out io.Writer) {
	_log.SetOutput(out)
}

func SetOutputByName(path string) error {
	return _log.SetOutputByName(path)
}

func SetOutputByMultiWriter(path string, writers ...io.Writer) error {
	return _log.SetOutputByMultiWriter(path, writers...)
}

func SetFlags(flags int) {
	_log._log.SetFlags(flags)
}

func SetColorByLogType(t LogType, colorFunc func(string, ...interface{}) string) {
	strType, _ := LogTypeToString(t)
	if _, bok := _colorMap[strType]; bok {
		return
	}
	_colorMap[strType] = colorFunc
}

func Enter(count int) {
	_log.Enter(count)
}

func Info(v ...interface{}) {
	_log.Info(v...)
}

func Infof(format string, v ...interface{}) {
	_log.Infof(format, v...)
}

func Debug(v ...interface{}) {
	_log.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	_log.Debugf(format, v...)
}

func Warn(v ...interface{}) {
	_log.Warning(v...)
}

func Warnf(format string, v ...interface{}) {
	_log.Warningf(format, v...)
}

func Warning(v ...interface{}) {
	_log.Warning(v...)
}

func Warningf(format string, v ...interface{}) {
	_log.Warningf(format, v...)
}

func Error(v ...interface{}) {
	_log.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	_log.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	_log.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	_log.Fatalf(format, v...)
}

func SetLevelByString(level string) {
	_log.SetLevelByString(level)
}

func SetHighlighting(highlighting bool) {
	_log.SetHighlighting(highlighting)
}

func SetRotateByDay() {
	_log.SetRotateByDay()
}

func SetRotateByHour() {
	_log.SetRotateByHour()
}

func Clone() *logger {
	return _log.Clone()
}

func LogPrint(t LogType, v ...interface{}) {
	_log.LogPrint(t, v...)
}

func LogPrintf(t LogType, format string, v ...interface{}) {
	_log.LogPrintf(t, format, v...)
}

type MLog struct {
	Module string
}

func (p *MLog) print(level string, datas ...interface{}) {
	if getLogLevel(level) > getLogLevel(Level) {
		return
	}

	switch level {
	case "Info":
		Info(datas...)
	case "Debug":
		Debug(datas...)
	case "Errof":
		Error(datas...)
	case "Warning":
		Warning(datas...)
	}
}

func (p *MLog) printf(level string, format string, datas ...interface{}) {
	if getLogLevel(level) > getLogLevel(Level) {
		return
	}

	switch level {
	case "Info":
		Infof(format, datas...)
	case "Debug":
		Debugf(format, datas...)
	case "Errof":
		Errorf(format, datas...)
	case "Warning":
		Warningf(format, datas...)
	}
}

func (p *MLog) Info(datas ...interface{}) { p.print("Info", datas...) }
func (p *MLog) Infof(format string, datas ...interface{}) {
	p.printf("Info", format, datas...)
}
func (p *MLog) Debug(datas ...interface{}) { p.print("Debug", datas...) }
func (p *MLog) Debugf(format string, datas ...interface{}) {
	p.printf("Debug", format, datas...)
}
func (p *MLog) Warning(datas ...interface{}) { p.print("Warning", datas...) }
func (p *MLog) Warningf(format string, datas ...interface{}) {
	p.printf("Warning", format, datas...)
}
func (p *MLog) Error(datas ...interface{}) { p.print("Error", datas...) }
func (p *MLog) Errorf(format string, datas ...interface{}) {
	p.printf("Error", format, datas...)
}
