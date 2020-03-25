package logger

import (
	"log"
	"strings"
	"sync"
)

// LogLevel defines the enum of available log levels
type LogLevel int

const (
	// LDebug most verbose log level
	LDebug = iota

	// LError very verbose, shows all messages except debug
	LError

	// LWarning verbose, shows all messages except debud and error
	LWarning

	// LInfo only shows info messages
	LInfo

	// LNone almost no output at all
	LNone
)

var lvl LogLevel = LInfo
var mux sync.Mutex = sync.Mutex{}

// Name returns the name of the loglevel
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

// NewLogLevelFromString initializes the loglevel from it's name string (case insensitive)
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

// SetLogLevel configures the log level to use
func SetLogLevel(l LogLevel) {
	mux.Lock()
	defer mux.Unlock()
	lvl = l
}

// GetLogLevel retrieves the log level to use
func GetLogLevel() LogLevel {
	mux.Lock()
	defer mux.Unlock()
	return lvl
}

// Debugln ...
func Debugln(v ...interface{}) {
	if GetLogLevel() <= LDebug {
		log.Println(v...)
	}
}

// Debugf ...
func Debugf(format string, args ...interface{}) {
	if GetLogLevel() <= LDebug {
		log.Printf(format, args...)
	}
}

// Errorln ...
func Errorln(v ...interface{}) {
	if GetLogLevel() <= LError {
		log.Println(v...)
	}
}

// Errorf ...
func Errorf(format string, args ...interface{}) {
	if GetLogLevel() <= LError {
		log.Printf(format, args...)
	}
}

// Warningln ...
func Warningln(v ...interface{}) {
	if GetLogLevel() <= LWarning {
		log.Println(v...)
	}
}

// Warningf ...
func Warningf(format string, args ...interface{}) {
	if GetLogLevel() <= LWarning {
		log.Printf(format, args...)
	}
}

// Infoln ...
func Infoln(v ...interface{}) {
	if GetLogLevel() <= LInfo {
		log.Println(v...)
	}
}

// Infof ...
func Infof(format string, args ...interface{}) {
	if GetLogLevel() <= LInfo {
		log.Printf(format, args...)
	}
}

// Println ...
func Println(v ...interface{}) {
	log.Println(v...)
}

// Printf ...
func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Fatalln ...
func Fatalln(v ...interface{}) {
	log.Fatalln(v...)
}

// Fatalf ...
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Panicln ...
func Panicln(v ...interface{}) {
	log.Panicln(v...)
}

// Panicf ...
func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}
