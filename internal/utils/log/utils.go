package log

import (
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

const DefaultTimestampFormat = "2006-01-02 15:04:05 -0700"

var DefaultOutputTo = os.Stderr

type InitOpts struct {
	Format          Format
	Level           Level
	OutputTo        io.Writer
	TimestampFormat string
}

func (o InitOpts) GetOutputTo() io.Writer {
	if o.OutputTo == nil {
		return DefaultOutputTo
	}
	return o.OutputTo
}

func (o InitOpts) GetTimestampFormat() string {
	if len(o.TimestampFormat) == 0 {
		return DefaultTimestampFormat
	}
	return o.TimestampFormat
}

func Init(opts InitOpts) {
	Logger.SetLevel(logLevelsMap[opts.Level])
	Logger.SetOutput(opts.GetOutputTo())
	Logger.SetReportCaller(true)
	switch opts.Format {
	case LogFormatText:
		Logger.SetFormatter(&logrus.TextFormatter{
			CallerPrettyfier:       PrettifyCaller,
			TimestampFormat:        opts.GetTimestampFormat(),
			FullTimestamp:          true,
			QuoteEmptyFields:       true,
			DisableSorting:         true,
			DisableLevelTruncation: true,
		})
	case LogFormatJSON:
		Logger.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: PrettifyCaller,
			TimestampFormat:  opts.GetTimestampFormat(),
		})
	default:
		Logger.SetFormatter(&DefaultFormatter)
	}
}

func PrettifyCaller(frame *runtime.Frame) (function string, file string) {
	function = path.Base(frame.Function)
	file = path.Base(frame.File)
	return
}
