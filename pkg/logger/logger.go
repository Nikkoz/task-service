package logger

import (
	"fmt"
	"strings"

	"github.com/Nikkoz/task-service/pkg/context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

func new(isProduction bool, l string) (*Logger, error) {
	var config zap.Config

	if isProduction {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	level := level(l)
	config.Level = zap.NewAtomicLevelAt(level)

	newLogger, err := config.Build(zap.AddCallerSkip(2))
	if err != nil {
		return nil, err
	}

	newLogger.Info("Set LOG_LEVEL", zap.Stringer("level", level))

	log = &Logger{logger: newLogger}

	return log, nil
}

func level(level string) zapcore.Level {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case "ERR", "ERROR":
		return zapcore.ErrorLevel
	case "WARN", "WARNING":
		return zapcore.WarnLevel
	case "INFO":
		return zapcore.InfoLevel
	case "DEBUG":
		return zapcore.DebugLevel
	case "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l *Logger) getContextFields(ctx context.Context) []zap.Field {
	return []zap.Field{zap.String("id", ctx.ID())}
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) DebugWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, l.getContextFields(ctx)...)

	l.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) InfoWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, l.getContextFields(ctx)...)

	l.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) WarnWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, l.getContextFields(ctx)...)

	l.Warn(msg, fields...)
}

func (l *Logger) Error(msg interface{}, fields ...zap.Field) {
	if msg == nil {
		return
	}

	switch v := msg.(type) {
	case string:
		l.logger.Error(v, fields...)
	case error:
		l.logger.Error(v.Error(), fields...)
	case fmt.Stringer:
		l.logger.Error(v.String(), fields...)
	default:
		l.logger.Error(fmt.Sprintf("%v", v), fields...)
	}
}

func (l *Logger) ErrorWithContext(ctx context.Context, err error, fields ...zap.Field) error {
	fields = append(fields, l.getContextFields(ctx)...)

	l.Error(err, fields...)

	return err
}

func (l *Logger) Fatal(msg interface{}, fields ...zap.Field) {
	if msg == nil {
		return
	}

	switch v := msg.(type) {
	case string:
		l.logger.Fatal(v, fields...)
	case error:
		l.logger.Fatal(v.Error(), fields...)
	case fmt.Stringer:
		l.logger.Fatal(v.String(), fields...)
	default:
		l.logger.Fatal(fmt.Sprintf("%v", v), fields...)
	}
}

func (l *Logger) FatalWithContext(ctx context.Context, err error, fields ...zap.Field) error {
	fields = append(fields, l.getContextFields(ctx)...)

	l.Fatal(err, fields...)

	return err
}
