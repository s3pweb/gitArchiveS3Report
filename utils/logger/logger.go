package logger

import (
	"fmt"

	"github.com/fatih/color"
)

type Logger struct {
	name  string
	level string
}

func NewLogger(name string, level string) (*Logger, error) {
	return &Logger{name: name, level: level}, nil
}

func getLevel(level string) int {
	switch level {
	case "error":
		return 4
	case "warn":
		return 3
	case "info":
		return 2
	case "debug":
		return 1
	case "trace":
		return 0
	}
	return 2 // Default to info level
}

func (l *Logger) Error(format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	color.Red(l.name + ": " + msg)
	return nil
}

func (l *Logger) Warn(format string, a ...interface{}) error {
	if getLevel(l.level) > getLevel("warn") {
		return nil
	}

	msg := fmt.Sprintf(format, a...)
	color.Yellow(l.name + ": " + msg)
	return nil
}

func (l *Logger) Info(format string, a ...interface{}) error {
	if getLevel(l.level) > getLevel("info") {
		return nil
	}

	msg := fmt.Sprintf(format, a...)
	color.Blue(l.name + ": " + msg)
	return nil
}

func (l *Logger) Debug(format string, a ...interface{}) error {
	if getLevel(l.level) > getLevel("debug") {
		return nil
	}

	msg := fmt.Sprintf(format, a...)
	color.Cyan(l.name + ": " + msg)
	return nil
}

func (l *Logger) Trace(format string, a ...interface{}) error {
	if getLevel(l.level) > getLevel("trace") {
		return nil
	}

	msg := fmt.Sprintf(format, a...)
	color.Black(l.name + ": " + msg)
	return nil
}

func (l *Logger) Success(format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	color.Green(l.name + ": " + msg)
	return nil
}
