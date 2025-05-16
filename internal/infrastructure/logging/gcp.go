package logging

import (
	"cloud.google.com/go/logging"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

var severityMap = map[string]logging.Severity{
	"DEFAULT":   logging.Default,
	"DEBUG":     logging.Debug,
	"INFO":      logging.Info,
	"NOTICE":    logging.Notice,
	"WARNING":   logging.Warning,
	"ERROR":     logging.Error,
	"CRITICAL":  logging.Critical,
	"ALERT":     logging.Alert,
	"EMERGENCY": logging.Emergency,
}
var loglevelMapToSeverityMap = map[zapcore.Level]logging.Severity{
	zapcore.DebugLevel: severityMap["DEBUG"],
	zapcore.InfoLevel:  severityMap["INFO"],
	zapcore.WarnLevel:  severityMap["WARNING"],
	zapcore.ErrorLevel: severityMap["ERROR"],
}
var severityToLogLevelMap = map[logging.Severity]zapcore.Level{
	severityMap["DEBUG"]:   zapcore.DebugLevel,
	severityMap["INFO"]:    zapcore.InfoLevel,
	severityMap["ERROR"]:   zapcore.ErrorLevel,
	severityMap["WARNING"]: zapcore.WarnLevel,
}

type GCPZapCore struct {
	zapcore.LevelEnabler
	gcpLogger *logging.Logger
	encoder   zapcore.Encoder
}

func NewGCPZapCore(envHostName string, level zapcore.Level, projectId, appName string, resType string, resLabelsProjectId string, resLabelsClusterName string, resLabelsLocation string, resLabelsNamespaceName string, resLabelsPodName string, minSeverity string, encoder zapcore.Encoder) (*GCPZapCore, error) {
	ctx := context.Background()
	client, err := logging.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}
	resource := &monitoredres.MonitoredResource{
		Type: resType,
		Labels: map[string]string{
			"cluster_name":   resLabelsClusterName,
			"location":       resLabelsLocation,
			"namespace_name": resLabelsNamespaceName,
			"project_id":     resLabelsProjectId,
			"pod_name":       envHostName,
		},
	}
	logger := client.Logger(appName, logging.CommonResource(resource))
	levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})

	return &GCPZapCore{
		LevelEnabler: levelEnabler,
		gcpLogger:    logger,
		encoder:      encoder,
	}, nil
}

func (c *GCPZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	buffer, err := c.encoder.EncodeEntry(entry, fields)
	if err != nil {
		return err
	}

	// 轉換為 JSON Payload
	var payload map[string]interface{}
	if err := json.Unmarshal(buffer.Bytes(), &payload); err != nil {
		return err
	}

	c.gcpLogger.Log(logging.Entry{
		Payload:   payload,
		Severity:  loglevelMapToSeverityMap[entry.Level],
		Timestamp: entry.Time,
	})
	return nil
}

func (c *GCPZapCore) Sync() error {
	return nil
}

func (c *GCPZapCore) With(fields []zapcore.Field) zapcore.Core {
	clone := *c
	clone.encoder = c.encoder.Clone()
	for _, f := range fields {
		f.AddTo(clone.encoder)
	}
	return &clone
}

func (c *GCPZapCore) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, c)
	}
	return checkedEntry
}

type GCPLoggerBase struct {
	Logger *zap.Logger
}

func NewGCPLoggerBase(envHostName, resType, resLabelsProjectId, resLabelsClusterName, resLabelsLocation, resLabelsNamespaceName, resLabelsPodName, minSeverity, projectId, appName string, level zapcore.Level) (ILoggerBase, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	core, err := NewGCPZapCore(envHostName, level, projectId, appName, resType, resLabelsProjectId, resLabelsClusterName, resLabelsLocation, resLabelsNamespaceName, resLabelsPodName, minSeverity, zapcore.NewJSONEncoder(encoderConfig))
	if err != nil {
		return nil, err
	}

	logger := zap.New(core, zap.AddCaller())
	return &GCPLoggerBase{Logger: logger}, nil
}

func (l *GCPLoggerBase) BInfo(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *GCPLoggerBase) BError(msg string, fields ...zap.Field) {

	l.Logger.Error(msg, fields...)
}

func (l *GCPLoggerBase) BDebug(msg string, fields ...zap.Field) {

	l.Logger.Debug(msg, fields...)
}

func (l *GCPLoggerBase) BWarn(msg string, fields ...zap.Field) {

	l.Logger.Warn(msg, fields...)
}

func (l *GCPLoggerBase) IsLevelEnabled(level zapcore.Level) bool {
	return l.Logger.Core().Enabled(level)
}

type GCPLogger struct {
	GCPLoggerBase ILoggerBase
	AppName       string
	AppIdentity   string
	Identity      []zap.Field
}

func NewGCPLogger(envHostName, resType, resLabelsProjectId, resLabelsClusterName, resLabelsLocation, resLabelsNamespaceName, resLabelsPodName, minSeverity, projectId, appName string) (ILogger, error) {
	logLevel := severityToLogLevelMap[severityMap[minSeverity]]

	base, err := NewGCPLoggerBase(envHostName, resType, resLabelsProjectId, resLabelsClusterName, resLabelsLocation, resLabelsNamespaceName, resLabelsPodName, minSeverity, projectId, appName, logLevel)
	if err != nil {
		return nil, err
	}
	return &GCPLogger{
		GCPLoggerBase: base,
		AppName:       appName,
		AppIdentity:   envHostName,
		Identity:      nil,
	}, nil
}

func (l *GCPLogger) Info(msg string, data any, fields ...zap.Field) {
	commonFields := []zap.Field{
		zap.String("appIdentity", l.AppIdentity),
		zap.String("application", l.AppName),
		zap.Any("data", data),
	}
	allFields := append(commonFields, fields...)
	l.GCPLoggerBase.BInfo(msg, allFields...)
}

func (l *GCPLogger) Error(msg string, err error, data any, fields ...zap.Field) {
	commonFields := []zap.Field{
		zap.String("appIdentity", l.AppIdentity),
		zap.String("application", l.AppName),
		zap.Error(err),
		zap.Any("data", data),
	}
	allFields := append(commonFields, fields...)
	l.GCPLoggerBase.BError(msg, allFields...)
}

func (l *GCPLogger) Debug(msg string, data any, fields ...zap.Field) {
	commonFields := []zap.Field{
		zap.String("appIdentity", l.AppIdentity),
		zap.String("application", l.AppName),
		zap.Any("data", data),
	}
	allFields := append(commonFields, fields...)
	l.GCPLoggerBase.BDebug(msg, allFields...)
}

func (l *GCPLogger) Warn(msg string, data any, fields ...zap.Field) {
	commonFields := []zap.Field{
		zap.String("appIdentity", l.AppIdentity),
		zap.String("application", l.AppName),
		zap.Any("data", data),
	}
	allFields := append(commonFields, fields...)
	l.GCPLoggerBase.BWarn(msg, allFields...)
}

func (l *GCPLogger) Write(p []byte) (n int, err error) {
	l.GCPLoggerBase.BDebug(string(p))
	return len(p), nil
}

func (l *GCPLogger) With(fields ...zap.Field) ILogger {
	// 複製當前的 Logger 並增加 fields
	newLogger := l.GCPLoggerBase
	newLogger.(*GCPLoggerBase).Logger = newLogger.(*GCPLoggerBase).Logger.With(fields...)
	return &GCPLogger{
		GCPLoggerBase: newLogger,
		AppName:       l.AppName,
		AppIdentity:   l.AppIdentity,
		Identity:      fields,
	}
}

func (l *GCPLogger) Clone() ILogger {
	clonedLoggerBase := &GCPLoggerBase{
		Logger: l.GCPLoggerBase.(*GCPLoggerBase).Logger.With(), // 這裡不帶任何 fields，所以會返回一個全新的 Logger 實例
	}
	return &GCPLogger{
		GCPLoggerBase: clonedLoggerBase,
		AppName:       l.AppName,
		AppIdentity:   l.AppIdentity,
		Identity:      append([]zap.Field{}, l.Identity...), // 深拷貝 Identity fields
	}
}
