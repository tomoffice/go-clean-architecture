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
}

func NewSeqCore(sender Sender, enc zapcore.Encoder, min zapcore.Level) zapcore.Core {
	return &seqCore{
		sender: sender,
		enc:    enc,
		min:    min,
	}
}
func (c *seqCore) Enabled(l zapcore.Level) bool {
	return l >= c.min
}
func (c *seqCore) With(_ []zapcore.Field) zapcore.Core {
	// 無狀態實現，直接返回自身
	return c
}
func (c *seqCore) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(e.Level) {
		return ce.AddCore(e, c)
	}
	return ce
}
func (c *seqCore) Write(e zapcore.Entry, fields []zapcore.Field) error {
	lF := make(logrus.Fields, len(fields))
	for _, field := range fields {
		if v := extractValue(field); v != nil {
			lF[field.Key] = v
		}
	}
	return c.sender.SendLevel(e.Level, e.Message, lF)
}
func (c *seqCore) Sync() error { return nil }
