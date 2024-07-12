package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo  Level = iota //start at 0
	LevelError              //1
	LevelFatal              //2
	LevelOff                //3
)

func (l Level) String() string {
	switch l {
	case LevelInfo: //0
		return "INFO"
	case LevelError: //1
		return "ERROR"
	case LevelFatal: //2
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level //default is 0/int8
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}
func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}
	log_entry := struct { //initializing anonymous struct and assigning values below
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}
	if level >= LevelError {
		log_entry.Trace = string(debug.Stack())
	}
	var line []byte
	line, err := json.Marshal(log_entry)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log " + err.Error())
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	//return 'number of bytes written' and error if any
	return l.out.Write(append(line, '\n'))
}

func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
