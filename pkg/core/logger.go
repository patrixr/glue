package core

import (
	"io"
	"os"

	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
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

	Stdout GlueWriter
	Stderr GlueWriter
}

func (gl *GlueLogger) Loud() {
	gl.Stderr.Loud = true
	gl.Stdout.Loud = true
}

func (gl *GlueLogger) Quiet() {
	gl.Stderr.Loud = false
	gl.Stdout.Loud = false
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
	writer := CreateGlueWriter(os.Stdout)
	writerErr := CreateGlueWriter(os.Stderr)

	options := log.Options{
		ReportTimestamp: true,
	}

	errLog := log.NewWithOptions(writerErr, options)
	stdLog := log.NewWithOptions(writer, options)

	stdLog.SetColorProfile(termenv.TrueColor)
	errLog.SetColorProfile(termenv.TrueColor)

	return &GlueLogger{
		errLog, stdLog,
		writer, writerErr,
	}
}

type GlueWriter struct {
	Loud      bool
	OutWriter io.Writer
}

func CreateGlueWriter(std io.Writer) GlueWriter {
	return GlueWriter{
		OutWriter: std,
		Loud:      true,
	}
}

func (gw GlueWriter) Write(p []byte) (n int, err error) {
	if gw.Loud {
		return gw.OutWriter.Write(p)
	}
	return 0, nil
}
