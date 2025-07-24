package seq

import (
	"github.com/nullseed/logruseq"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap/zapcore"
)

type logrusSender struct {
	lg *logrus.Logger
}

func NewLogrusSender(endPoint, api string) Sender {
	lg := logrus.New()
	lg.SetFormatter(&logrus.JSONFormatter{})
	lg.AddHook(logruseq.NewSeqHook(endPoint, logruseq.OptionAPIKey(api)))
	return &logrusSender{lg: lg}
}

func (s *logrusSender) SendLevel(level zapcore.Level, msg string, fields logrus.Fields) error {
	entry := s.lg.WithFields(fields)
	entry.Log(mapLevel(level), msg)
	return nil
}
