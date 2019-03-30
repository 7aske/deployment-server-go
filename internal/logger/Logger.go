package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
)

const timeFormat string = "2006-01-02"
const (
	LOG_DEPLOYER string = "deployer"
	LOG_SERVER   string = "server"
)

type Logger struct {
	logger      log.Logger
	currentDate string
	loggerType  string
}

func NewLogger(ltype string) *Logger {
	logger := &Logger{}
	logger.loggerType = ltype
	logger.updateLogger()
	return logger
}
func (l *Logger) Log(message string) {
	if l.currentDate != time.Now().Format(timeFormat) {
		l.updateLogger()
	}
	l.logger.Println(message)
}

func (l *Logger) updateLogger() {
	cwd, _ := os.Getwd()
	l.currentDate = time.Now().Format(timeFormat)
	logFP := path.Join(cwd, "logs", l.loggerType, l.currentDate+".log")
	err := os.MkdirAll(path.Dir(logFP), 0776)
	if err != nil {
		fmt.Println(err)
	}
	fi, err := os.OpenFile(logFP, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	mw := io.MultiWriter(os.Stdout, fi)
	l.logger = *log.New(mw, "", log.LstdFlags)
}
