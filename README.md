# Go日志工具
本工具在于简化日志在控制台的输出或输出到文件中
## 使用方法
将整个logger文件夹放到`GO_PATH/src/`下，在main文件头部`import "logger"`

## 日志级别
```go
const (
	DEBUG LogLevel = iota
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)
```

## 控制台输出日志
```go
//参数：日志的级别
newLog := logger.NewLog(logger.INFO) 
//参数1：打印的级别，大于等于这个级别的日志才会被打印  
//参数2-N：打印信息
newLog.RecordLog(logger.INFO,"id:",1,"name:","张三") 

输出结果：[日期][日志级别][调用方法:调用文件:行号][打印信息]
示例：[2020-10-04 17:38:11] [INFO] [main:emptyInterface.go:11] [id: 1 name: 张三]
```

## 输出日志到指定目录

```go

//参数1：日志级别 LogLevel
//参数2：输出目录 string  
//参数3：文件名  string
//参数4：文件的大小最大值 int64 (如果日志内容大于这个值会备份原来的文件并新创建一个文件)
log := logger.NewFileLogger(logger.ERROR, "C:\\Go\\src\\", "test.log", 1) 

//参数1：输出的级别，大于等于这个级别的日志才会被输出
//参数2-N：输出内容 interface{}
log.RecordFileLog(logger.INFO, "测试日志")
```
