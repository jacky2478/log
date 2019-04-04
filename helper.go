package mlog

import (
    "os"
    "fmt"
    "runtime"
    "path/filepath"
    "strings"
    "time"
)

func StringToLogLevel(level string) LogLevel {
    switch level {
    case "fatal":
        return LOG_LEVEL_FATAL
    case "error":
        return LOG_LEVEL_ERROR
    case "warn":
        return LOG_LEVEL_WARN
    case "warning":
        return LOG_LEVEL_WARN
    case "debug":
        return LOG_LEVEL_DEBUG
    case "info":
        return LOG_LEVEL_INFO
    }
    return LOG_LEVEL_ALL
}

func LogLevelToString(level LogLevel) string {
    switch level {
    case LOG_LEVEL_FATAL:
        return "fatal"
    case LOG_LEVEL_ERROR:
        return "error"
    case LOG_LEVEL_WARN:
        return "warn"
    //case LOG_LEVEL_WARN:
    //	return "warning"
    case LOG_LEVEL_DEBUG:
        return "debug"
    case LOG_LEVEL_INFO:
        return "info"
    }
    return "unknownLevel"
}

func LogTypeToString(t LogType) (string, string) {
    switch t {
    case LOG_FATAL:
        return "fatal", "[0;31"
    case LOG_ERROR:
        return "error", "[0;31"
    case LOG_WARNING:
        return "warning", "[0;33"
    case LOG_DEBUG:
        return "debug", "[0;36"
    case LOG_INFO:
        return "info", "[0;37"
    }
    return "unknown", "[0;37"
}

func genDayTime(t time.Time) string {
    return t.Format(FORMAT_TIME_DAY)
}

func genHourTime(t time.Time) string {
    return t.Format(FORMAT_TIME_HOUR)
}

func GetOutPutInfo(strDirFile, strLevel, strColor string) string {
    return strLevel + " " + strDirFile + strColor
}

func GetDirFileInfo(depth int) string {
    _, file, line, ok := runtime.Caller(depth)
    if !ok {
        return ""
    }

    path := filepath.Dir(file)
    dirs := strings.Split(path, "/")
    strParent := "[" + dirs[len(dirs) - 1] + "]"
    strFileLine := "[" + filepath.Base(file) + fmt.Sprintf(":%d", line) + "]"
    return strings.Join([]string{strParent, strFileLine}, " ")
}

func GetLvlColorStr(t LogType, msg string) (strLevel, strColor string) {
    strFlag, _ := LogTypeToString(t)
    strLevel = "[" + strFlag + "]"
    strColor = strings.TrimPrefix(msg, strLevel)
    return
}

func isDirExists(name string) bool {
    f, err := os.Stat(name)
    if err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return f.IsDir()
}