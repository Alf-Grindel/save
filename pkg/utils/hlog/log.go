package hlog

import (
	"context"
	"fmt"
	"io"
)

// Level defines the priority of a log message.
type Level int

// The levels of logs
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var strs = []string{
	"[Debug] ",
	"[Info] ",
	"[Warn] ",
	"[Error] ",
	"[Fatal] ",
}

// FormatLogger is a logger interface that output logs with a format.
type FormatLogger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

// Logger is a logger interface that provides logging function with levels.
type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
}

// CtxLogger is a logger interface that accepts a context argument and output
// logs with a format.
type CtxLogger interface {
	CtxDebugf(ctx context.Context, format string, v ...interface{})
	CtxInfof(ctx context.Context, format string, v ...interface{})
	CtxWarnf(ctx context.Context, format string, v ...interface{})
	CtxErrorf(ctx context.Context, format string, v ...interface{})
	CtxFatalf(ctx context.Context, format string, v ...interface{})
}

// Control provides methods to config a logger.
type Control interface {
	SetLevel(Level)
	SetOutput(io.Writer)
}

// FullLogger is the conbination of Logger, FormatLogger, CtxLogger and Control.
type FullLogger interface {
	Logger
	FormatLogger
	CtxLogger
	Control
}

func (lv Level) toString() string {
	if lv >= LevelDebug && lv <= LevelFatal {
		return strs[lv]
	}
	return fmt.Sprintf("[?%d]", lv)
}
