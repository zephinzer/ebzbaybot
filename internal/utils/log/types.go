package log

import "github.com/sirupsen/logrus"

type SimpleLog func(...interface{})
type FormattedLog func(string, ...interface{})

type Format uint

const (
	LogFormatText Format = 0
	LogFormatJSON Format = 1
)

type Level uint

const (
	LogLevelTrace Level = 0
	LogLevelDebug Level = 1
	LogLevelInfo  Level = 2
	LogLevelWarn  Level = 3
	LogLevelError Level = 4
)

var logLevelsMap = map[Level]logrus.Level{
	LogLevelTrace: logrus.TraceLevel,
	LogLevelDebug: logrus.DebugLevel,
	LogLevelInfo:  logrus.InfoLevel,
	LogLevelWarn:  logrus.WarnLevel,
	LogLevelError: logrus.ErrorLevel,
}
