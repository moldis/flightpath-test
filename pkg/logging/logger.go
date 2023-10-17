package logging

import (
	"artemb/flights-path/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewLogger(config *config.Logging) (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	err := level.UnmarshalText([]byte(config.Level))
	if err != nil {
		return nil, err
	}

	cw := zapcore.Lock(os.Stdout)
	je := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "init_timestamp",
		LevelKey:       "log_level",
		NameKey:        "log_name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	zapCore := zapcore.NewCore(je, cw, level)

	logger := zap.New(
		zapCore,
		zap.AddCaller(),
		zap.AddStacktrace(zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= zapcore.ErrorLevel
		})),
	)

	return logger, nil
}
