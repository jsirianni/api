package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a configured Zap logger suitable for
// container application which need to log structured
// messages to sdout.
func New() (*zap.Logger, error) {
	logConf := zap.NewProductionConfig()
	logConf.OutputPaths = []string{"stdout"}
	logConf.EncoderConfig.MessageKey = "message"
	logConf.EncoderConfig.TimeKey = "time"
	logConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return logConf.Build()
}
