package logger

import (
	"fmt"
	"os"
	"path"
	"time"
)

var (
	//日志通道最大值
	MaxSize = 50000
	//后台线程数
	threadNums = 1
)

type FileLogger struct {
	Level       LogLevel
	filePath    string //保存路径
	fileName    string //文件名
	fileObj     *os.File
	errFileObj  *os.File
	maxFileSize int64
	logChan     chan *logMsg
}

type logMsg struct {
	Level     LogLevel
	msg       interface{}
	funcName  string
	fileName  string
	timestamp string
	line      int
}

func NewFileLogger(levelStr LogLevel, fp, fn string, maxSize int64) *FileLogger {
	fileLogger := FileLogger{
		Level:       levelStr,
		filePath:    fp,
		fileName:    fn,
		maxFileSize: maxSize,
		logChan:     make(chan *logMsg, MaxSize),
	}
	err := fileLogger.initFile()
	if err != nil {
		panic(err)
	}
	return &fileLogger
}

func (f *FileLogger) Close() {
	f.fileObj.Close()
	f.errFileObj.Close()
}

func (f *FileLogger) initFile() error {

	fullFileName := path.Join(f.filePath, f.fileName)
	file, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file failed,err:%v\n", err)
		return err
	}

	errFile, err := os.OpenFile(fullFileName+".err", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Printf("open err log file failed,err:%v\n", err)
		return err
	}

	defer f.Close()

	f.fileObj = file
	f.errFileObj = errFile

	//开启后台日志
	for i := 0; i < threadNums; i++ {
		go f.writeInBack()
	}

	return nil
}

func (f *FileLogger) writeInBack() {
	isBig, err := checkSize(f.fileObj, f.maxFileSize)
	if err != nil {
		fmt.Printf("check size error err:%v\n", err)
		return
	}
	if isBig {
		f.fileObj = f.splitFile(f.fileObj)
	}
	select {
	case logTmp := <-f.logChan:
		_, _ = fmt.Fprintf(f.fileObj, "[%s] [%s] [%s:%s:%d] %v\n", logTmp.timestamp, parseLevel(f.Level), logTmp.funcName, logTmp.fileName, logTmp.line, logTmp.msg)
		if f.Level >= ERROR {
			if isBig {
				f.errFileObj = f.splitFile(f.errFileObj)
			}
			_, _ = fmt.Fprintf(f.errFileObj, "[%s] [%s] [%s:%s:%d] %v\n", logTmp.timestamp, parseLevel(f.Level), logTmp.funcName, logTmp.fileName, logTmp.line, logTmp.msg)
		}
	default: //取不到日志休息500毫秒
		time.Sleep(time.Millisecond * 500)
	}

}

func (f *FileLogger) RecordFileLog(printLevel LogLevel, msg ...interface{}) {
	if f.Level >= printLevel {

		funcName, fileName, lineNo := getInfo(2)

		//先把日志发送到通道中
		logTmp := &logMsg{
			Level:     printLevel,
			msg:       msg,
			funcName:  funcName,
			fileName:  fileName,
			timestamp: time.Now().Format("2006-01-02 15:04:05"),
			line:      lineNo,
		}

		select {
		case f.logChan <- logTmp:
		default: //保证业务代码顺畅执行
		}

	}
}

func checkSize(file *os.File, maxSize int64) (bool, error) {
	stat, err := file.Stat()
	if err != nil {
		return false, err
	}
	return stat.Size() > maxSize, nil
}

func (f *FileLogger) splitFile(file *os.File) *os.File {
	file.Close()
	nowStr := time.Now().Format("20060102150405000")
	logName := path.Join(f.filePath, f.fileName)
	newLogName := fmt.Sprintf("%s.bak%s", logName, nowStr)
	_ = os.Rename(logName, newLogName)
	file, err := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open new log file failed.err:%v\n", err)
		return nil
	}
	return file
}
