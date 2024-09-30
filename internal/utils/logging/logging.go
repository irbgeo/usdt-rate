package logging

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type errorArguments interface {
	Args() map[string]any
}

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

// Init initializes the logger.
func Init(logLevel string) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	level, err := zapcore.ParseLevel(logLevel)
	if err == nil {
		config.Level.SetLevel(level)
	}

	logger, _ = config.Build()
	sugar = logger.Sugar()
}

// SetLogLevel sets the log level.
func SetLogLevel(logLevel string) error {
	if logLevel == "" {
		return nil
	}

	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	logger.Core().Enabled(level)
	return nil
}

// Error formats and prints an error.
func Error(err error, args ...interface{}) {
	fields := parseFields(err, args...)
	sugar.Errorw("", append(fields, zap.Error(err))...)
}

// Info formats and prints an info line.
func Info(message string, args ...interface{}) {
	fields := parseFields(nil, args...)
	sugar.Infow(message, fields...)
}

func parseFields(err error, args ...interface{}) []interface{} {
	if len(args)%2 != 0 {
		return nil
	}

	fields := make([]interface{}, 0, len(args)+2)

	if raw, ok := err.(errorArguments); ok {
		for k, v := range raw.Args() {
			fields = append(fields, k, v)
		}
	}

	for i := 0; i < len(args); i += 2 {
		key := args[i].(string) // nolint:errcheck
		value := args[i+1]

		if key == "wasted_time" {
			switch v := value.(type) {
			case time.Duration:
				fields = append(fields, key, v.Seconds())
			case float64, int, int64:
				fields = append(fields, key, v)
			default:
				fields = append(fields, key, value)
			}
		} else {
			fields = append(fields, key, value)
		}
	}

	if err != nil {
		fields = append(fields, "trace", trace(err))
	}

	return fields
}

func trace(x interface{}) string {
	if raw, ok := x.(error); ok {
		err, ok := errors.Wrap(raw, "").(stackTracer)
		if !ok {
			return ""
		}

		st := err.StackTrace()
		if len(st) > 2 {
			a := fmt.Sprintf("%+v", st[2])
			a = strings.Replace(a, "\n\t", " at ", 1)
			return a
		}
	}
	return ""
}
