package hlog

import (
	"io"
	"log"
	"os"
)

var (
	logger FullLogger = &defaultLogger{
		stdlog: log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds),
		depth:  4,
	}
)

// SetOutPut sets the output of default logger.
func SetOutput(w io.Writer) {
	logger.SetOutput(w)
}

// SetLevel sets the level of logs below which logs will not be output.
func SetLevel(lv Level) {
	logger.SetLevel(lv)
}

// DefaultLogger return the default logger.
func DefaultLogger() FullLogger {
	return logger
}

// SetLogger sets the default logger.
func SetLogger(v FullLogger) {
	logger = v
}
