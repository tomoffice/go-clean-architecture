package seq

import (
	"github.com/nullseed/logruseq"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap/zapcore"
)

// NewLogrusSender 建立安裝 SeqHook 的 Logrus Sender
func NewLogrusSender(cfg Config) Sender {
	lg := logrus.New()
	lg.SetFormatter(&logrus.JSONFormatter{})
	lg.AddHook(logruseq.NewSeqHook(cfg.Endpoint, logruseq.OptionAPIKey(cfg.APIKey)))
	return &logrusSender{lg: lg}
}

type logrusSender struct {
	lg *logrus.Logger
}

func (s *logrusSender) SendLevel(level zapcore.Level, msg string, fields logrus.Fields) error {
	entry := s.lg.WithFields(fields)
	entry.Log(mapLevel(level), msg)
	return nil
}
