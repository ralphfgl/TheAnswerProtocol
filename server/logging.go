package main

import (
	"fmt"
	"sync"
	"time"
)

type LogLevel string

const (
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Event     string                 `json:"event"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

type Logger struct {
	Mu sync.Mutex
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) logEntry(
	level LogLevel,
	event string,
	message string,
	data map[string]interface{},
) {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     string(level),
		Event:     event,
		Message:   message,
		Data:      data,
	}

	fmt.Printf(
		"[%s] [%s] %s: %s\n",
		entry.Timestamp,
		entry.Level,
		event,
		message,
	)
}

func (l *Logger) Info(event, message string, data map[string]interface{}) {
	l.logEntry(LevelInfo, event, message, data)
}

func (l *Logger) Warn(event, message string, data map[string]interface{}) {
	l.logEntry(LevelWarn, event, message, data)
}

func (l *Logger) Error(event, message string, data map[string]interface{}) {
	l.logEntry(LevelError, event, message, data)
}
