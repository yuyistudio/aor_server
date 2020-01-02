package log

type Logger interface {
	Debug(fmt string, args ...interface{})
	Info(fmt string, args ...interface{})
	Error(fmt string, args ...interface{})
}

var logger Logger

func SetLogger(newLogger Logger) {
	logger = newLogger
}

func HasLogger() bool {
	return logger != nil
}

func Debug(fmt string, args ...interface{}) {
	logger.Info(fmt, args...)
}

func Info(fmt string, args ...interface{}) {
	logger.Info(fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	logger.Error(fmt, args...)
}
