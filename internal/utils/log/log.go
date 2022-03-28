package log

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetOutput(os.Stderr)
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&DefaultFormatter)
}

func SetDebugMode(on bool) {
	if on {
		Logger.SetLevel(logrus.TraceLevel)
	} else {
		Logger.SetLevel(logrus.InfoLevel)
	}
}

// Print writes the log to stdout
var Print,
	// Trace writes the log at the trace level
	Trace,
	// Debug writes the log at the debug level
	Debug,
	// Info writes the log at the info level
	Info,
	// Warn writes the log at the warn level
	Warn,
	// Error writes the log at the error level
	Error SimpleLog = print, Logger.Trace, Logger.Debug, Logger.Info, Logger.Warn, Logger.Error

// Printf writes the log to stdout with formatting
var Printf,
	// Tracef writes the log at the trace level with formatting
	Tracef,
	// Debugf writes the log at the debug level with formatting
	Debugf,
	// Infof writes the log at the info level with formatting
	Infof,
	// Warnf writes the log at the warn level with formatting
	Warnf,
	// Errorf writes the log at the error level with formatting
	Errorf FormattedLog = printf, Logger.Tracef, Logger.Debugf, Logger.Infof, Logger.Warnf, Logger.Errorf

func print(args ...interface{}) {
	fmt.Print(args...)
}

func printf(message string, args ...interface{}) {
	fmt.Printf(message, args...)
}

var DefaultFormatter = logrus.TextFormatter{
	CallerPrettyfier:       PrettifyCaller,
	FullTimestamp:          true,
	QuoteEmptyFields:       true,
	DisableSorting:         true,
	DisableLevelTruncation: false,
}
