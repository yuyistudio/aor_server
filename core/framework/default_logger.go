package framework

import (
	"fmt"
	"log"
)

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

type DefaultLogger struct {
}

func (l *DefaultLogger) Debug(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Output(3, s)
}

func (l *DefaultLogger) Info(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Output(3, s)
}

func (l *DefaultLogger) Error(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Output(3, s)
}
