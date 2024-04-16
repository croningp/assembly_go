package assembly

import (
	"log"
	"os"
)

// This file defines the logger struct and instantiates a global logger

type BuiltinLogger struct {
	logger *log.Logger
}

func NewBuiltinLogger() *BuiltinLogger {
	return &BuiltinLogger{logger: log.New(os.Stdout, "", 5)}
}

func (l *BuiltinLogger) Debug(args ...interface{}) {
	l.logger.Println(args...)
}

func (l *BuiltinLogger) Debugf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *BuiltinLogger) SetOutput(f *os.File){
	l.logger.SetOutput(f)
}

var Logger = NewBuiltinLogger()