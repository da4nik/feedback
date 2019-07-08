package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

// LoggerOpts - parameters to configure logger
type LoggerOpts struct {
}

var logger = logrus.New()

// InitLogger - initializes logger
func InitLogger(opts LoggerOpts) {
	logger.Level = logrus.DebugLevel
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

// Errorf logs error message
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Infof logs info message to log
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Debugf logs info message to log
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}
