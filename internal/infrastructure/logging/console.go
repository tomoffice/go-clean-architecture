package logging

import (
	"fmt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ConsoleLoggerBase struct {
	Logger *zap.Logger
}

func NewConsoleBase(level zapcore.Level) (ILoggerBase, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.Development = false
	cfg.Level = zap.NewAtomicLevelAt(level) // 直接設定日誌級別
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	zapLogger, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(3))
	if err != nil {
		return nil, err
	}
	return &ConsoleLoggerBase{Logger: zapLogger}, nil
}
func (l *ConsoleLoggerBase) BInfo(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}
func (l *ConsoleLoggerBase) BError(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}
func (l *ConsoleLoggerBase) BDebug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}
func (l *ConsoleLoggerBase) BWarn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}
func (l *ConsoleLoggerBase) IsLevelEnabled(level zapcore.Level) bool {
	return l.Logger.Core().Enabled(level)
}

type ConsoleLogger struct {
	ConsoleLoggerBase ILoggerBase
}

func NewConsoleLogger(level zapcore.Level) (ILogger, error) {
	base, err := NewConsoleBase(level)
	if err != nil {
		return nil, err
	}
	return &ConsoleLogger{
		ConsoleLoggerBase: base,
	}, nil
}
func (l *ConsoleLogger) Info(msg string, data any, fields ...zap.Field) {

	if l.ConsoleLoggerBase.IsLevelEnabled(zapcore.InfoLevel) {
		fields = append(fields, zap.Any("data", data))
		l.ConsoleLoggerBase.BInfo(msg, fields...)
	}
}
func (l *ConsoleLogger) Error(msg string, err error, data any, fields ...zap.Field) {
	if l.ConsoleLoggerBase.IsLevelEnabled(zapcore.ErrorLevel) {
		fields = append(fields, zap.Any("err", err))
		fields = append(fields, zap.Any("data", data))
		l.ConsoleLoggerBase.BError(msg, fields...)
	}
}
func (l *ConsoleLogger) Debug(msg string, data any, fields ...zap.Field) {
	if l.ConsoleLoggerBase.IsLevelEnabled(zapcore.DebugLevel) {
		fields = append(fields, zap.Any("data", data))
		l.ConsoleLoggerBase.BDebug(msg, fields...)
	}
}
func (l *ConsoleLogger) Warn(msg string, data any, fields ...zap.Field) {
	if l.ConsoleLoggerBase.IsLevelEnabled(zapcore.WarnLevel) {
		fields = append(fields, zap.Any("data", data))
		l.ConsoleLoggerBase.BWarn(msg, fields...)
	}
}
func (l *ConsoleLogger) Write(p []byte) (n int, err error) {
	l.Debug(string(p), nil)
	return len(p), nil
}
func (l *ConsoleLogger) With(fields ...zap.Field) ILogger {
	// 複製當前的 Logger 並增加 fields
	newLogger := l.ConsoleLoggerBase
	newLogger.(*ConsoleLoggerBase).Logger = newLogger.(*ConsoleLoggerBase).Logger.With(fields...)
	return &ConsoleLogger{
		ConsoleLoggerBase: newLogger,
	}
}
func (l *ConsoleLogger) Clone() ILogger {
	return &ConsoleLogger{
		ConsoleLoggerBase: l.ConsoleLoggerBase,
	}
}

type CliLoggerBase struct {
	Logger *zap.Logger
}

func NewCliLoggerBase(level zapcore.Level) (ILoggerBase, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level) // 直接設定日誌級別
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	zapLogger, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(3))
	if err != nil {
		return nil, err
	}

	return &CliLoggerBase{Logger: zapLogger}, nil
}
func (l *CliLoggerBase) BInfo(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}
func (l *CliLoggerBase) BError(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}
func (l *CliLoggerBase) BDebug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}
func (l *CliLoggerBase) BWarn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}
func (l *CliLoggerBase) IsLevelEnabled(level zapcore.Level) bool {
	return l.Logger.Core().Enabled(level)
}

type CliLogger struct {
	CliLoggerBase ILoggerBase
}

func NewCliLogger(level zapcore.Level) (ILogger, error) {
	base, err := NewCliLoggerBase(level)
	if err != nil {
		return nil, err
	}
	return &CliLogger{
		CliLoggerBase: base,
	}, nil
}
func (l *CliLogger) Info(msg string, data any, fields ...zap.Field) {
	if l.CliLoggerBase.IsLevelEnabled(zapcore.InfoLevel) {
		color.Green(msg)
	}
}
func (l *CliLogger) Error(msg string, err error, data any, fields ...zap.Field) {
	if l.CliLoggerBase.IsLevelEnabled(zapcore.ErrorLevel) {
		color.Red(msg)
	}
}
func (l *CliLogger) Debug(msg string, data any, fields ...zap.Field) {
	if l.CliLoggerBase.IsLevelEnabled(zapcore.DebugLevel) {
		color.Blue(msg)
	}
}
func (l *CliLogger) Warn(msg string, data any, fields ...zap.Field) {
	if l.CliLoggerBase.IsLevelEnabled(zapcore.WarnLevel) {
		color.Yellow(msg)
	}
}
func (l *CliLogger) Write(p []byte) (n int, err error) {
	//l.Debug(string(p), nil)
	fmt.Println(string(p))
	return len(p), nil
}
func (l *CliLogger) With(fields ...zap.Field) ILogger {
	// 複製當前的 Logger 並增加 fields
	newLogger := l.CliLoggerBase
	newLogger.(*CliLoggerBase).Logger = newLogger.(*CliLoggerBase).Logger.With(fields...)
	return &CliLogger{
		CliLoggerBase: newLogger,
	}
}
func (l *CliLogger) Clone() ILogger {
	return &CliLogger{
		CliLoggerBase: l.CliLoggerBase,
	}
}
