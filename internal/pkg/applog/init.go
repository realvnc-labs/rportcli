package applog

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type CleanFormatter struct {
	TextFormatter logrus.TextFormatter
}

func (f *CleanFormatter) getLine(entry *logrus.Entry) (msg []byte) {
	msg = []byte(fmt.Sprintf("%s: %s\n", strings.ToUpper(entry.Level.String()), entry.Message))
	return msg
}

func (f *CleanFormatter) getDebugLine(entry *logrus.Entry) (msg []byte, err error) {
	msg, err = f.TextFormatter.Format(entry)
	return msg, err
}

// TODO: Not sure if this multiline support is required.
func (f *CleanFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	multiline, ok := entry.Data["multiline"]
	if ok {
		delete(entry.Data, "multiline")
	}
	var res []byte
	if entry.Level == logrus.FatalLevel || entry.Level == logrus.WarnLevel {
		res = f.getLine(entry)
	} else {
		var err error
		res, err = f.getDebugLine(entry)
		if err != nil {
			return nil, err
		}
	}
	if multiline, ok := multiline.(string); ok && multiline != "" {
		res = append(res, []byte(multiline)...)
	}
	return res, nil
}

func Init() (f *CleanFormatter) {
	f = &CleanFormatter{
		TextFormatter: logrus.TextFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000000",
			FullTimestamp:   true,
		},
	}

	logrus.SetFormatter(f)

	return f
}

func SetLogLevel(isVerbose bool) {
	logLevel := logrus.InfoLevel
	if isVerbose {
		logLevel = logrus.DebugLevel
	}

	logrus.SetLevel(logLevel)

	logrus.Debugf("will set log level to %v", logLevel)
}
