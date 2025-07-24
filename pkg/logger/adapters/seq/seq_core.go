package seq

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/zap/zapcore"
)

// seqCore 實作 zapcore.Core，橋接到 Sender
type seqCore struct {
	sender Sender
	enc    zapcore.Encoder
	min    zapcore.Level
	fields []zapcore.Field
}

func NewSeqCore(sender Sender, enc zapcore.Encoder, min zapcore.Level) zapcore.Core {
	return &seqCore{
		sender: sender,
		enc:    enc,
		min:    min,
		fields: make([]zapcore.Field, 0),
	}
}
func (c *seqCore) Enabled(l zapcore.Level) bool {
	return l >= c.min
}
func (c *seqCore) With(fields []zapcore.Field) zapcore.Core {
	clone := &seqCore{
		sender: c.sender,
		enc:    c.enc.Clone(),
		min:    c.min,
		fields: make([]zapcore.Field, len(c.fields)+len(fields)),
	}
	copy(clone.fields, c.fields)
	copy(clone.fields[len(c.fields):], fields)
	return clone
}
func (c *seqCore) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(e.Level) {
		return ce.AddCore(e, c)
	}
	return ce
}
func (c *seqCore) Write(e zapcore.Entry, fields []zapcore.Field) error {
	allFields := make([]zapcore.Field, 0, len(c.fields)+len(fields))
	allFields = append(allFields, c.fields...)
	allFields = append(allFields, fields...)

	lF := make(logrus.Fields, len(allFields))
	for _, field := range allFields {
		if v := extractValue(field); v != nil {
			lF[field.Key] = v
		}
	}
	return c.sender.SendLevel(e.Level, e.Message, lF)
}
func (c *seqCore) Sync() error { return nil }
