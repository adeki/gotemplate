package logger

import (
	"github.com/doloopwhile/logrusltsv"
	"github.com/sirupsen/logrus"
)

func detectFormatter(format LogFormat) logrus.Formatter {
	switch format {
	case JSONFormat:
		return &logrus.JSONFormatter{}
	case LTSVFormat:
		return &logrusltsv.Formatter{}
	case TextFormat:
		return &logrus.TextFormatter{}
	default:
		return &logrus.TextFormatter{}
	}
}

func newLogger(cnfs configs) Logger {
	logger := logrus.New()
	logger.SetFormatter(detectFormatter(cnfs.format))
	logger.SetOutput(cnfs.output)
	lvl, _ := logrus.ParseLevel(cnfs.level.String())
	logger.SetLevel(lvl)

	return &logrusLogger{
		logger: logger,
	}
}

type logrusLogger struct {
	logger *logrus.Logger
}

func (l *logrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}
func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *logrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}
func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *logrusLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}
func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *logrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}
func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *logrusLogger) Print(args ...interface{}) {
	l.logger.Print(args...)
}
func (l *logrusLogger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *logrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}
func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *logrusLogger) WithFields(fields Fields) Logger {
	return &logrusEntry{
		entry: l.logger.WithFields(logrus.Fields(fields)),
	}
}

type logrusEntry struct {
	entry *logrus.Entry
}

func (e *logrusEntry) Debug(args ...interface{}) {
	e.entry.Debug(args...)
}
func (e *logrusEntry) Debugf(format string, args ...interface{}) {
	e.entry.Debugf(format, args...)
}

func (e *logrusEntry) Error(args ...interface{}) {
	e.entry.Error(args...)
}
func (e *logrusEntry) Errorf(format string, args ...interface{}) {
	e.entry.Errorf(format, args...)
}

func (e *logrusEntry) Fatal(args ...interface{}) {
	e.entry.Fatal(args...)
}
func (e *logrusEntry) Fatalf(format string, args ...interface{}) {
	e.entry.Fatalf(format, args...)
}

func (e *logrusEntry) Info(args ...interface{}) {
	e.entry.Info(args...)
}
func (e *logrusEntry) Infof(format string, args ...interface{}) {
	e.entry.Infof(format, args...)
}

func (e *logrusEntry) Print(args ...interface{}) {
	e.entry.Print(args...)
}
func (e *logrusEntry) Printf(format string, args ...interface{}) {
	e.entry.Printf(format, args...)
}

func (e *logrusEntry) Warn(args ...interface{}) {
	e.entry.Warn(args...)
}
func (e *logrusEntry) Warnf(format string, args ...interface{}) {
	e.entry.Warnf(format, args...)
}

func (e *logrusEntry) WithFields(fields Fields) Logger {
	return &logrusEntry{
		entry: e.entry.WithFields(logrus.Fields(fields)),
	}
}
