package logging

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel string

const (
	InfoLevel  LogLevel = "info"
	ErrorLevel LogLevel = "error"
	DebugLevel LogLevel = "debug"
)

// New returns a configured Zap logger suitable for
// container application which need to log structured
// messages to sdout.
func New(level LogLevel) (*zap.Logger, error) {
	logLevel, err := zap.ParseAtomicLevel(string(level))
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %s: %v", level, err)
	}

	logConf := zap.NewProductionConfig()
	logConf.OutputPaths = []string{"stdout"}
	logConf.EncoderConfig.MessageKey = "message"
	logConf.EncoderConfig.TimeKey = "time"
	logConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConf.Level = logLevel
	logConf.EncoderConfig.StacktraceKey = ""
	return logConf.Build()
}
