package logger

import "io"

type configs struct {
	level  LogLevel
	format LogFormat
	output io.Writer
}

type Config func(*configs)

func WithLogLevel(lvl LogLevel) Config {
	return func(cnfs *configs) {
		cnfs.level = lvl
	}
}

func WithLogFormat(format LogFormat) Config {
	return func(cnfs *configs) {
		cnfs.format = format
	}
}

func WithOutput(out io.Writer) Config {
	return func(cnfs *configs) {
		cnfs.output = out
	}
}
