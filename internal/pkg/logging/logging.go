// Package logging provides an abstraction layer for logging functionality.
package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Interface wraps the logrus.FieldLogger interface.methods.
type Interface interface {
	logrus.FieldLogger
	WrappedLogger() *logrus.Logger
}

// Logger is a struct to wrap a *logrus.Logger.
type Logger struct {
	*logrus.Logger
}

// New returns a new *Logger instance.
func New() *Logger {
	return &Logger{
		Logger: &logrus.Logger{
			Formatter: &logrus.JSONFormatter{},
			Out:       os.Stdout,
			Level:     logrus.InfoLevel,
		},
	}
}

// WrappedLogger returns the wrapped *logrus.Logger instance.
func (log *Logger) WrappedLogger() *logrus.Logger {
	return log.Logger
}
