package core

import (
	"os"

	"github.com/charmbracelet/log"
)

type Logger interface {
	Info(msg interface{}, keyvals ...interface{})
	Debug(msg interface{}, keyvals ...interface{})
	Warn(msg interface{}, keyvals ...interface{})
	Error(msg interface{}, keyvals ...interface{})
}

type GlueLogger struct {
	errLog *log.Logger
	stdLog *log.Logger
}

func (gl *GlueLogger) Info(msg interface{}, keyvals ...interface{}) {
	gl.stdLog.Info(msg, keyvals...)
}
func (gl *GlueLogger) Debug(msg interface{}, keyvals ...interface{}) {
	gl.stdLog.Debug(msg, keyvals...)
}
func (gl *GlueLogger) Warn(msg interface{}, keyvals ...interface{}) {
	gl.errLog.Warn(msg, keyvals...)
}
func (gl *GlueLogger) Error(msg interface{}, keyvals ...interface{}) {
	gl.errLog.Error(msg, keyvals...)
}

func CreateLogger() Logger {
	options := log.Options{
		ReportTimestamp: true,
	}

	errLog := log.NewWithOptions(os.Stderr, options)
	stdLog := log.NewWithOptions(os.Stdout, options)

	return &GlueLogger{
		errLog, stdLog,
	}
}