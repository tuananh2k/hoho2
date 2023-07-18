package infrastructure

import (
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: service xử lý log cho framework
*/

type logData struct {
	message string
	data    map[string]interface{}
}
type symperLog struct {
	errorLogger *os.File
	warnLogger  *os.File
	infoLogger  *os.File
}
type STrace struct {
	File string
	Line int
	Func string
}

var l *symperLog

// NewLogger creates a logger.
func NewLogger() {
	l = &symperLog{
		errorLogger: newLogrus("./../symper_error_log.log"),
		infoLogger:  newLogrus("./../symper_log.log"),
		warnLogger:  newLogrus("./../symper_warn_log.log"),
	}
}

func newLogrus(path string) *os.File {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Println(err)
		return nil
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	return file
}

func Warn(message string, data map[string]interface{}) {
	logData := new(logData)
	logData.message = message
	logData.data = data
	logrus.SetOutput(l.warnLogger)
	logrus.WithFields(logData.data).Warn(logData.message)
}
func Info(message string, data map[string]interface{}) {
	logData := new(logData)
	logData.message = message
	logData.data = data
	logrus.SetOutput(l.infoLogger)
	logrus.WithFields(logData.data).Info(logData.message)
}
func Error(message string, data map[string]interface{}) {
	logData := new(logData)
	logData.message = message
	logData.data = data
	logrus.SetOutput(l.errorLogger)
	logrus.WithFields(logData.data).Error(logData.message)
}

func Trace() *STrace {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return &STrace{
		File: frame.File,
		Func: frame.Function,
		Line: frame.Line,
	}
}
