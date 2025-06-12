package core

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

// 颜色
const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

// LogFormatter 日志格式化
type LogFormatter struct{}

// Format 实现Formatter(entry *logrus.Entry) ([]byte, error)接口
func (t *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	//根据不同的level去展示颜色
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	//自定义日期格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	prefix := "hsy"
	if entry.HasCaller() {
		//自定义文件路径
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		//自定义输出格式
		fmt.Fprintf(b,
			"%s %s  \033[%dm%-8s\033[0m %s %s %s\n",
			prefix,
			timestamp,
			levelColor,
			entry.Level,
			fileVal,
			funcVal,
			entry.Message)
	} else {
		fmt.Fprintf(b,
			"%s %s \033[%dm%-8s\033[0m %s\n",
			prefix,
			timestamp,
			levelColor,
			entry.Level,
			entry.Message)
	}
	return b.Bytes(), nil
}

// InitLogger 初始化log
func InitLogger() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)          //设置输出类型
	log.SetReportCaller(true)         //开启返回函数名和行号
	log.SetFormatter(&LogFormatter{}) //设置自己定义的Formatter
	log.SetLevel(logrus.DebugLevel)   //设置最低的Level

	logrus.SetOutput(os.Stdout)          //设置输出类型
	logrus.SetReportCaller(true)         //开启返回函数名和行号
	logrus.SetFormatter(&LogFormatter{}) //设置自己定义的Formatter
	//logLevel, err := logrus.ParseLevel(logrus.InfoLevel)
	//if err != nil {
	//	logrus.Warnf("日志级别设置错误 %s %s", logLevel, err)
	//	logrus.Warnf("设置默认日志级别 warn")
	//	logLevel = logrus.WarnLevel
	//}
	logrus.SetLevel(logrus.InfoLevel)

	fileDate := time.Now().Format("2006-01-02")

	logPath := "logs"
	appName := "node"

	//创建目录
	err := os.MkdirAll(fmt.Sprintf("%s/%s", logPath, fileDate), os.ModePerm)
	if err != nil {
		logrus.Fatalf(err.Error())
		return nil
	}

	filename := fmt.Sprintf("%s/%s/%s.log", logPath, fileDate, appName)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		logrus.Fatalf(err.Error())
		return nil
	}
	fileHook := FileDateHook{file, logPath, fileDate, appName}
	log.AddHook(&fileHook)
	logrus.AddHook(&fileHook)
	return log
}

// FileDateHook 按照时间分割的hook
type FileDateHook struct {
	file     *os.File
	logPath  string
	fileDate string //判断日期切换目录
	appName  string
}

func (hook FileDateHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
func (hook FileDateHook) Fire(entry *logrus.Entry) error {
	timer := entry.Time.Format("2006-01-02")
	line, _ := entry.String()

	// 单独存一份错误的日志
	if entry.Level == logrus.ErrorLevel {
		filename := fmt.Sprintf("%s/%s/err.log", hook.logPath, timer)
		file, _ := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		file.Write([]byte(line))
	}

	if hook.fileDate == timer {
		hook.file.Write([]byte(line))
		return nil
	}
	// 时间不等
	hook.file.Close()
	os.MkdirAll(fmt.Sprintf("%s/%s", hook.logPath, timer), os.ModePerm)
	filename := fmt.Sprintf("%s/%s/%s.log", hook.logPath, timer, hook.appName)

	hook.file, _ = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	hook.fileDate = timer
	hook.file.Write([]byte(line))
	return nil
}

// LoggerEsHook 将日志分发到es数据库
type LoggerEsHook struct {
}
