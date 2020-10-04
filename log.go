package logger

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

const (
	DEBUG LogLevel = iota
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

type LogLevel uint16


type Logger struct {
	LogLevel
}

func parseLevel(l LogLevel) string {
	switch l {
	case 0:
		return "DEBUG"
	case 1:
		return "TRACE"
	case 2:
		return "INFO"
	case 3:
		return "WARNING"
	case 4:
		return "ERROR"
	default:
		return "FATAL"
	}
}

func NewLog(level LogLevel) Logger {
	return Logger{
		level,
	}
}

func (level Logger) RecordLog(printLevel LogLevel, msg ...interface{}) {
	if level.LogLevel >= printLevel {
		t := time.Now().Format("2006-01-02 15:04:05")
		funcName, fileName, lineNo := getInfo(2)
		fmt.Printf("[%s] [%s] [%s:%s:%d] %v\n", t, parseLevel(level.LogLevel), funcName, fileName, lineNo, msg)
	}
}

func getInfo(skip int) (funcName, fileName string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		fmt.Printf("runtime.Caller() failed\n")
		return
	}
	funcName = strings.Split(runtime.FuncForPC(pc).Name(), ".")[1]
	fileName = path.Base(file)
	return
}

