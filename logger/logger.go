package logger

import (
	"log"
	"strings"
	"sync"
)

type LogLevel int

const (
	LDebug = iota
	LError
	LWarning
	LInfo
	LNone
)

var lvl LogLevel = LInfo
var mux sync.Mutex = sync.Mutex{}

func (l LogLevel) Name() string {
	switch l {
	case LDebug:
		return "Debug"
	case LError:
		return "Error"
	case LWarning:
		return "Warning"
	case LInfo:
		return "Info"
	default:
		return "None"
	}
}

func NewLogLevelFromString(in string) LogLevel {
	switch strings.ToLower(in) {
	case "debug":
		return LDebug
	case "error":
		return LError
	case "warning":
		return LWarning
	case "info":
		return LInfo
	default:
		return LNone
	}
}

func SetLogLevel(l LogLevel) {
	mux.Lock()
	defer mux.Unlock()
	lvl = l
}

func GetLogLevel() LogLevel {
	mux.Lock()
	defer mux.Unlock()
	return lvl
}

func Debugln(v ...interface{}) {
	if GetLogLevel() <= LDebug {
		log.Println(v...)
	}
}

func Debugf(format string, args ...interface{}) {
	if GetLogLevel() <= LDebug {
		log.Printf(format, args...)
	}
}

func Errorln(v ...interface{}) {
	if GetLogLevel() <= LError {
		log.Println(v...)
	}
}

func Errorf(format string, args ...interface{}) {
	if GetLogLevel() <= LError {
		log.Printf(format, args...)
	}
}

func Warningln(v ...interface{}) {
	if GetLogLevel() <= LWarning {
		log.Println(v...)
	}
}

func Warningf(format string, args ...interface{}) {
	if GetLogLevel() <= LWarning {
		log.Printf(format, args...)
	}
}

func Infoln(v ...interface{}) {
	if GetLogLevel() <= LInfo {
		log.Println(v...)
	}
}

func Infof(format string, args ...interface{}) {
	if GetLogLevel() <= LInfo {
		log.Printf(format, args...)
	}
}

func Println(v ...interface{}) {
	log.Println(v...)
}

func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Fatalln(v ...interface{}) {
	log.Fatalln(v...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Panicln(v ...interface{}) {
	log.Panicln(v...)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}
