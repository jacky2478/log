package mlog

var (
	Level = "Info"

	GetLog func(string) ILog
)

func init() {
	GetLog = func(module string) ILog {
		return &SysLog{Module: module}
	}
}

func UseSys() {
	GetLog = func(module string) ILog {
		return &SysLog{Module: module}
	}
}

func UseMlog() {
	GetLog = func(module string) ILog {
		return &MLog{Module: module}
	}
}

type ILog interface {
	Info(datas ...interface{})
	Infof(format string, datas ...interface{})

	Debug(datas ...interface{})
	Debugf(format string, datas ...interface{})

	Warning(datas ...interface{})
	Warningf(format string, datas ...interface{})

	Error(datas ...interface{})
	Errorf(format string, datas ...interface{})
}
