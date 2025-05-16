package logging

import (
	"encoding/json"
	"fmt"
	"github.com/nullseed/logruseq"
	logrus_ "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

const (
	TimeKey      = "@t"
	LevelKey     = "@l"
	MessageKey   = "@mt"
	UnknownLevel = "unknown"
)

var logLevelMap = map[string]logrus_.Level{
	"panic":   logrus_.PanicLevel,
	"fatal":   logrus_.FatalLevel,
	"error":   logrus_.ErrorLevel,
	"warn":    logrus_.WarnLevel,
	"warning": logrus_.WarnLevel,
	"info":    logrus_.InfoLevel,
	"debug":   logrus_.DebugLevel,
	"trace":   logrus_.TraceLevel,
}

type SeqHookAdapter struct {
	Hook *logruseq.SeqHook
}

func NewSeqHookAdapter(hook *logruseq.SeqHook) *SeqHookAdapter {
	return &SeqHookAdapter{Hook: hook}
}
func (s *SeqHookAdapter) Write(p []byte) (n int, err error) {
	var logEntry map[string]interface{}
	if err := json.Unmarshal(p, &logEntry); err != nil {
		log.Println("Error unmarshalling log entry:", err)
		return 0, err
	}

	// Parse timestamp
	timeStr, ok := logEntry[TimeKey].(string)
	if !ok {
		log.Println("Timestamp is not a string")
		return 0, fmt.Errorf("timestamp is not a string")
	}

	parsedTime, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		log.Println("Error parsing timestamp:", err)
		return 0, err
	}

	// Parse log level
	levelStr, ok := logEntry[LevelKey].(string)
	if !ok {
		log.Println("Log level is not a string")
		return 0, fmt.Errorf("log level is not a string")
	}

	logLevel, err := parseLogLevel(levelStr)
	if err != nil {
		log.Println("Error parsing log level:", err)
		return 0, err
	}
	// Copy all fields except the reserved ones
	logEntryCP := make(map[string]interface{}, len(logEntry))
	for key, value := range logEntry {
		logEntryCP[key] = value
	}
	// Remove reserved fields
	keysToRemove := []string{TimeKey, LevelKey, MessageKey}
	for _, key := range keysToRemove {
		delete(logEntryCP, key)
	}
	// Add emmitTime field
	logEntryCP["emmitTime"] = parsedTime
	logEntryCP["emmitGroupId"] = logEntry["emmitGroupId"]
	// If data is nil or empty delete data
	data, ok := logEntryCP["data"]
	if !ok || data == nil {
		delete(logEntryCP, "data")
	}
	entry := &logrus_.Entry{
		Data:    logEntryCP,
		Time:    time.Now(),
		Level:   logLevel,
		Message: logEntry[MessageKey].(string),
	}

	err = s.Hook.Fire(entry)
	if err != nil {
		log.Println("Error firing log entry to Seq:", err)
	}
	return len(p), nil
}

func parseLogLevel(levelStr string) (logrus_.Level, error) {
	level, exists := logLevelMap[levelStr]
	if !exists {
		return logrus_.InfoLevel, fmt.Errorf("unknown log level: %s", levelStr)
	}
	return level, nil
}

type SeqLoggerBase struct {
	Logger *zap.Logger
}

func NewSeqLoggerBase(level zapcore.Level, seqURL string, token string) (ILoggerBase, error) {
	hook := logruseq.NewSeqHook(seqURL, logruseq.OptionAPIKey(token))
	adapter := NewSeqHookAdapter(hook)
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:       "@t",
		LevelKey:      "@l",
		NameKey:       "Logger",
		MessageKey:    "@mt",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: zapcore.TimeEncoder(func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format(time.RFC3339Nano))
		}),
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	encoder := zapcore.NewJSONEncoder(encoderCfg)
	core := zapcore.NewCore(encoder, zapcore.AddSync(adapter), level)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(3))
	return &SeqLoggerBase{Logger: zapLogger}, nil
}
func (l *SeqLoggerBase) BInfo(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *SeqLoggerBase) BError(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

func (l *SeqLoggerBase) BDebug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

func (l *SeqLoggerBase) BWarn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

func (l *SeqLoggerBase) IsLevelEnabled(level zapcore.Level) bool {
	return l.Logger.Core().Enabled(level)
}

type SeqLogger struct {
	SeqLoggerBase ILoggerBase
	AppName       string
	AppIdentity   string
	Identity      []zap.Field
}

func NewSeqLogger(level zapcore.Level, seqURL, token, appName, appIdentity string) (ILogger, error) {
	base, err := NewSeqLoggerBase(level, seqURL, token)
	if err != nil {
		return nil, err
	}
	return &SeqLogger{
		SeqLoggerBase: base,
		AppName:       appName,
		AppIdentity:   appIdentity,
	}, nil
}

func (l *SeqLogger) Info(msg string, data any, fields ...zap.Field) {
	commonFields := []zap.Field{
		zap.String("identity", l.AppIdentity),
		zap.String("application", l.AppName),
		zap.Any("data", data),
	}
	allFields := append(commonFields, fields...)
	l.SeqLoggerBase.BInfo(msg, allFields...)
}

func (l *SeqLogger) Error(msg string, err error, data any, fields ...zap.Field) {
	commonFields := []zap.Field{
		zap.String("identity", l.AppIdentity),
		zap.String("application", l.AppName),
		zap.Error(err),
		zap.Any("data", data),
	}
	allFields := append(commonFields, fields...)
	l.SeqLoggerBase.BError(msg, allFields...)
}

func (l *SeqLogger) Debug(msg string, data any, fields ...zap.Field) {
	commonFields := []zap.Field{
		zap.String("identity", l.AppIdentity),
		zap.String("application", l.AppName),
		zap.Any("data", data),
	}
	allFields := append(commonFields, fields...)
	l.SeqLoggerBase.BDebug(msg, allFields...)
}

func (l *SeqLogger) Warn(msg string, data any, fields ...zap.Field) {
	commonFields := []zap.Field{
		zap.String("identity", l.AppIdentity),
		zap.String("application", l.AppName),
		zap.Any("data", data),
	}
	allFields := append(commonFields, fields...)
	l.SeqLoggerBase.BWarn(msg, allFields...)
}

func (l *SeqLogger) Write(p []byte) (n int, err error) {
	l.Debug(string(p), nil)
	return len(p), nil
}

func (l *SeqLogger) With(fields ...zap.Field) ILogger {
	newLogger := l.SeqLoggerBase
	newLogger.(*SeqLoggerBase).Logger = newLogger.(*SeqLoggerBase).Logger.With(fields...)
	return &SeqLogger{
		SeqLoggerBase: l.SeqLoggerBase,
		AppName:       l.AppName,
		AppIdentity:   l.AppIdentity,
		Identity:      fields,
	}
}

func (l *SeqLogger) Clone() ILogger {
	clonedLoggerBase := &GCPLoggerBase{
		Logger: l.SeqLoggerBase.(*SeqLoggerBase).Logger.With(), // 這裡不帶任何 fields，所以會返回一個全新的 Logger 實例
	}
	return &SeqLogger{
		SeqLoggerBase: clonedLoggerBase,
		AppName:       l.AppName,
		AppIdentity:   l.AppIdentity,
		Identity:      append([]zap.Field{}, l.Identity...),
	}
}
