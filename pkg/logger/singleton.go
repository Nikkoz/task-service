package logger

import (
	"github.com/Nikkoz/task-service/pkg/context"
	"go.uber.org/zap"
)

var log *Logger

func New(isProduction bool, level string) {
	if log == nil {
		newLogger, err := new(isProduction, level)
		if err != nil {
			panic(err)
		}

		log = newLogger
	}
}

func getLog() *Logger {
	return log
}

func Debug(msg string, fields ...zap.Field) {
	if l := getLog(); l != nil {
		l.Debug(msg, fields...)
	}
}

func DebugWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	if l := getLog(); l != nil {
		l.DebugWithContext(ctx, msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if l := getLog(); l != nil {
		l.Info(msg, fields...)
	}
}

func InfoWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	if l := getLog(); l != nil {
		l.InfoWithContext(ctx, msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if l := getLog(); l != nil {
		l.Warn(msg, fields...)
	}
}

func WarnWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	if l := getLog(); l != nil {
		l.WarnWithContext(ctx, msg, fields...)
	}
}

func Error(msg interface{}, fields ...zap.Field) {
	if l := getLog(); l != nil {
		l.Error(msg, fields...)
	}
}

func ErrorWithContext(ctx context.Context, err error, fields ...zap.Field) error {
	if l := getLog(); l != nil {
		return l.ErrorWithContext(ctx, err, fields...)
	}
	return err
}

func Fatal(msg interface{}, fields ...zap.Field) {
	if l := getLog(); l != nil {
		l.Fatal(msg, fields...)
	}
}

func FatalWithContext(ctx context.Context, err error, fields ...zap.Field) error {
	if l := getLog(); l != nil {
		return l.FatalWithContext(ctx, err, fields...)
	}
	return err
}
