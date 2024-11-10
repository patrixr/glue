package core

import (
	"os"

	"github.com/charmbracelet/log"
)

type Logger interface {
	Info(msg interface{}, keyvals ...interface{})
	Infof(format string, args ...interface{})
	Debug(msg interface{}, keyvals ...interface{})
	Debugf(format string, args ...interface{})
	Warn(msg interface{}, keyvals ...interface{})
	Warnf(format string, args ...interface{})
	Error(msg interface{}, keyvals ...interface{})
	Errorf(format string, args ...interface{})
}

type GlueLogger struct {
	errLog *log.Logger
	stdLog *log.Logger

	Stdout *os.File
	Stderr *os.File
}

func (gl *GlueLogger) Info(msg interface{}, keyvals ...interface{}) {
	gl.stdLog.Info(msg, keyvals...)
}

func (gl *GlueLogger) Infof(format string, args ...interface{}) {
	gl.stdLog.Infof(format, args...)
}

func (gl *GlueLogger) Debug(msg interface{}, keyvals ...interface{}) {
	gl.stdLog.Debug(msg, keyvals...)
}

func (gl *GlueLogger) Debugf(format string, args ...interface{}) {
	gl.stdLog.Debugf(format, args...)
}

func (gl *GlueLogger) Warn(msg interface{}, keyvals ...interface{}) {
	gl.errLog.Warn(msg, keyvals...)
}

func (gl *GlueLogger) Warnf(format string, args ...interface{}) {
	gl.stdLog.Warnf(format, args...)
}

func (gl *GlueLogger) Error(msg interface{}, keyvals ...interface{}) {
	gl.errLog.Error(msg, keyvals...)
}

func (gl *GlueLogger) Errorf(format string, args ...interface{}) {
	gl.stdLog.Errorf(format, args...)
}

func CreateLogger() *GlueLogger {
	options := log.Options{
		ReportTimestamp: true,
	}

	errLog := log.NewWithOptions(os.Stderr, options)
	stdLog := log.NewWithOptions(os.Stdout, options)

	return &GlueLogger{
		errLog, stdLog,
		os.Stdout, os.Stderr,
	}
}
