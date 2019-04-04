//high level log wrapper, so it can output different log based on level
package mlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type logger struct {
	_log         *log.Logger
	level        LogLevel
	highlighting bool

	dailyRolling bool
	hourRolling  bool

	fileName  string
	logSuffix string
	fd        *os.File

	// support set depth for runtime.Caller
	depth     int

	lock sync.Mutex
}

func New() *logger {
	return Newlogger(os.Stderr, "")
}

func Newlogger(w io.Writer, prefix string) *logger {
	return &logger{_log: log.New(w, prefix, LstdFlags), level: LOG_LEVEL_ALL, highlighting: true}
}

func (l *logger) Clone() *logger {
	p := &logger{}
	p._log = l._log
	p.level = l.level
	p.highlighting = l.highlighting

	p.dailyRolling = l.dailyRolling
	p.hourRolling = l.hourRolling
	p.fileName = l.fileName
	p.logSuffix = l.logSuffix
	p.fd = l.fd
	p.lock = l.lock
	return p
}

func (l *logger) GetSysLog() *log.Logger {
	return l._log
}

func (l *logger) SetHighlighting(highlighting bool) {
	l.highlighting = highlighting
}

func (l *logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *logger) SetDepth(depth int) {
	l.depth = depth
}

func (l *logger) SetLevelByString(level string) {
	l.level = StringToLogLevel(level)
}

func (l *logger) SetRotateByDay() {
	l.dailyRolling = true
	l.logSuffix = genDayTime(time.Now())
}

func (l *logger) SetRotateByHour() {
	l.hourRolling = true
	l.logSuffix = genHourTime(time.Now())
}

func (l *logger) rotate() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	var suffix string
	if l.dailyRolling {
		suffix = genDayTime(time.Now())
	} else if l.hourRolling {
		suffix = genHourTime(time.Now())
	} else {
		return nil
	}

	// Notice: if suffix is not equal to l.LogSuffix, then rotate
	if suffix != l.logSuffix {
		// if split by hour, save the files to the diectory named by curent date
		err := l.doRotate(l.getRotatePath(l.fileName), suffix)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *logger) getRotatePath(filename string) string {
	if l.hourRolling == false {
		return ""
	}

	logDir := filepath.Dir(filename)
	if strings.HasSuffix(logDir, "/") == false {
		logDir += "/"
	}
	rotatePath := logDir + genDayTime(time.Now()) + "/"

	// check for rotatePath
	if isDirExists(rotatePath) == false {
		os.Mkdir(rotatePath, os.ModePerm)
	}
	return rotatePath
}

func (l *logger) doRotate(rotatePath, suffix string) error {
	// Notice: Not check error, is this ok?
	l.fd.Close()

	lastFileName := rotatePath + filepath.Base(l.fileName) + "." + l.logSuffix
	err := os.Rename(l.fileName, lastFileName)
	if err != nil {
		return err
	}

	err = l.SetOutputByName(l.fileName)
	if err != nil {
		return err
	}

	l.logSuffix = suffix

	return nil
}

func (l *logger) SetOutput(out io.Writer) {
	l._log = log.New(out, l._log.Prefix(), l._log.Flags())
}

func (l *logger) SetOutputByMultiWriter(path string, writers ...io.Writer) error {
	var err error
	var f *os.File
	if path != "" {
		f, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Fatal(err)
		}
		writers = append(writers, f)

		l.fileName = path
		l.fd = f
	}

	if len(writers) > 1 {
		w := io.MultiWriter(writers...)
		l.SetOutput(w)
	} else if len(writers) == 1 {
		l.SetOutput(writers[0])
	}
	return err
}

func (l *logger) SetOutputByName(path string) error {
	l.fileName = path
	path = l.getRotatePath(path) + filepath.Base(path)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		l.fileName = ""
		log.Fatal(err)
	}
	l.SetOutput(f)

	// l.fileName = path
	l.fd = f

	return err
}

func (l *logger) log(t LogType, v ...interface{}) {
	if l.level|LogLevel(t) != l.level {
		return
	}

	err := l.rotate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	v1 := make([]interface{}, len(v)+2)
	logStr, logColor := LogTypeToString(t)
	if l.highlighting {
		v1[0] = "\033" + logColor + "m[" + logStr + "]"
		copy(v1[1:], v)
		v1[len(v)+1] = "\033[0m"
	} else {
		v1[0] = "[" + logStr + "]"
		copy(v1[1:], v)
		v1[len(v)+1] = ""
	}

	s := fmt.Sprint(v1...)
	sdirFile := GetDirFileInfo(l.depth)
	sLvl, scolor := GetLvlColorStr(t, s)
	if f, bok := _colorMap[logStr]; bok && l.fileName == "" {
		if len(v) > 1 {
			scolor = f("%v", v1...)
		}
		if len(v) == 1 {
			scolor = f("%v", scolor)
		}
	}
	s = GetOutPutInfo(sdirFile, sLvl, scolor)
	l._log.Output(l.depth, s)
	//OutPut(l, 4, s)
}

func (l *logger) logf(t LogType, format string, v ...interface{}) {
	if l.level|LogLevel(t) != l.level {
		return
	}

	err := l.rotate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	logStr, logColor := LogTypeToString(t)
	var s string
	if l.highlighting {
		s = "\033" + logColor + "m[" + logStr + "] " + fmt.Sprintf(format, v...) + "\033[0m"
	} else {
		s = "[" + logStr + "] " + fmt.Sprintf(format, v...)
	}

	sdirFile := GetDirFileInfo(l.depth)
	sLvl, scolor := GetLvlColorStr(t, s)
	if f, bok := _colorMap[logStr]; bok && l.fileName == "" {
		scolor = f("%v", scolor)
	}
	s = GetOutPutInfo(sdirFile, sLvl, scolor)
	l._log.Output(l.depth, s)
	//OutPut(l, 4, s)
}

func (l *logger) Fatal(v ...interface{}) {
	l.log(LOG_FATAL, v...)
	os.Exit(-1)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.logf(LOG_FATAL, format, v...)
	os.Exit(-1)
}

func (l *logger) Error(v ...interface{}) {
	l.log(LOG_ERROR, v...)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.logf(LOG_ERROR, format, v...)
}

func (l *logger) Warning(v ...interface{}) {
	l.log(LOG_WARNING, v...)
}

func (l *logger) Warningf(format string, v ...interface{}) {
	l.logf(LOG_WARNING, format, v...)
}

func (l *logger) Debug(v ...interface{}) {
	l.log(LOG_DEBUG, v...)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.logf(LOG_DEBUG, format, v...)
}

func (l *logger) Info(v ...interface{}) {
	l.log(LOG_INFO, v...)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.logf(LOG_INFO, format, v...)
}

func (l *logger) Enter(count int) {
	var strPrint string = "[info]"
	for i := 0; i < count; i++ {
		strPrint += "\n"
	}
	l._log.Output(4, strPrint)
	//OutPut(l, 4, strPrint)
}

func (l *logger) LogPrint(t LogType, v ...interface{}) {
	l.log(t, v...)
}

func (l *logger) LogPrintf(t LogType, format string, v ...interface{}) {
	l.logf(t, format, v...)
}
