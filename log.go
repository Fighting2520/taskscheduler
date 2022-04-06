package taskscheduler

import (
	"log"
	"os"
)

type Logger interface {
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

type simpleLog struct {
}

func (sl *simpleLog) Infof(fmt string, args ...interface{}) {
	InfoLog.Printf(fmt, args...)
}

func (sl *simpleLog) Errorf(fmt string, args ...interface{}) {
	ErrLog.Println(fmt, args)
}

var (
	ErrLog  = log.New(os.Stderr, "[ERROR]", 0)
	InfoLog = log.New(os.Stdout, "[INFO]", 0)
)
