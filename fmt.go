package mlog

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

type SysLog struct {
	Module string
}

func (p *SysLog) print(level string, fprintf func(string, ...interface{}), datas ...interface{}) {
	if getLogLevel(level) > getLogLevel(Level) {
		return
	}
	format := fmt.Sprintf("[%v] [%v] %v", level, p.Module, stackTrace(3)) + " %v"
	fprintf(format, datas...)
}

func (p *SysLog) printf(level string, fprintf func(string, ...interface{}), format string, datas ...interface{}) {
	if getLogLevel(level) > getLogLevel(Level) {
		return
	}
	format = fmt.Sprintf("[%v] [%v] %v ", level, p.Module, stackTrace(3)) + format
	fprintf(format, datas...)
}

func (p *SysLog) Info(datas ...interface{}) { p.print("Info", log.Printf, datas...) }
func (p *SysLog) Infof(format string, datas ...interface{}) {
	p.printf("Info", log.Printf, format, datas...)
}
func (p *SysLog) Debug(datas ...interface{}) { p.print("Debug", log.Printf, datas...) }
func (p *SysLog) Debugf(format string, datas ...interface{}) {
	p.printf("Debug", log.Printf, format, datas...)
}
func (p *SysLog) Warning(datas ...interface{}) { p.print("Warning", log.Printf, datas...) }
func (p *SysLog) Warningf(format string, datas ...interface{}) {
	p.printf("Warning", log.Printf, format, datas...)
}
func (p *SysLog) Error(datas ...interface{}) { p.print("Error", log.Printf, datas...) }
func (p *SysLog) Errorf(format string, datas ...interface{}) {
	p.printf("Error", log.Printf, format, datas...)
}

func stackTrace(depth int) string {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		return ""
	}
	strFileLine := "[" + filepath.Base(file) + fmt.Sprintf(":%d", line) + "]"
	return strFileLine
}

func getLogLevel(level string) LogLevel {
	switch level {
	case "Info":
		return LOG_LEVEL_INFO
	case "Debug":
		return LOG_LEVEL_DEBUG
	case "Error":
		return LOG_LEVEL_WARN
	}
	return 0
}
