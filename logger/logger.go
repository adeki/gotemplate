package logger

import (
	"os"
	"sync"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Print(args ...interface{})
	Printf(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	WithFields(fields Fields) Logger
}

type Fields map[string]interface{}

type LogLevel string

func (lvl LogLevel) String() string {
	return string(lvl)
}

type LogFormat string

func (lf LogFormat) String() string {
	return string(lf)
}

const (
	DebugLevel LogLevel = "debug"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	InfoLevel  LogLevel = "info"

	JSONFormat LogFormat = "json"
	LTSVFormat LogFormat = "ltsv"
	TextFormat LogFormat = "text"
)

var (
	logger Logger
	once   sync.Once
)

func Init(args ...Config) {
	once.Do(func() {
		logger = process(args...)
	})
}

func process(args ...Config) Logger {
	cnfs := configs{
		level:  InfoLevel,
		format: TextFormat,
		output: os.Stderr,
	}
	for _, c := range args {
		c(&cnfs)
	}
	return newLogger(cnfs)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Print(args ...interface{}) {
	logger.Print(args...)
}
func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func WithFields(fields Fields) Logger {
	return logger.WithFields(fields)
}
