package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type LogLevel string

const (
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

type Logger struct {
	mu sync.Mutex
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) log(
	level LogLevel,
	format string,
	args ...any,
) {
	l.mu.Lock()
	defer l.mu.Unlock()
	timestamp := time.Now().Format("15:04:05")
	message := fmt.Sprintf(format, args...)
	var levelText string
	switch level {
	case LevelInfo:
		levelText = "\033[34mINFO \033[0m"
	case LevelWarn:
		levelText = "\033[33mWARN \033[0m"
	case LevelError:
		levelText = "\033[31mERROR\033[0m"
	}
	output := fmt.Sprintf(
		"\033[2m%s\033[0m %s %s\n",
		timestamp,
		levelText,
		message,
	)
	if level == LevelError {
		fmt.Fprint(os.Stderr, output)
	} else {
		fmt.Print(output)
	}
}

func (l *Logger) Info(format string, args ...any) {
	l.log(LevelInfo, format, args...)
}

func (l *Logger) Warn(format string, args ...any) {
	l.log(LevelWarn, format, args...)
}

func (l *Logger) Error(format string, args ...any) {
	l.log(LevelError, format, args...)
}
