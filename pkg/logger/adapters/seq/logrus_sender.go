package seq

import (
	"github.com/nullseed/logruseq"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap/zapcore"
	"io"
)

type LogrusSender struct {
	lg *logrus.Logger
}

func NewLogrusSender(endPoint, api string, consoleOutputEnable bool) *LogrusSender {
	lg := logrus.New()
	if consoleOutputEnable == false {
		lg.SetOutput(io.Discard) // 不輸出到控制台
	}
	lg.SetFormatter(&logrus.JSONFormatter{})
	lg.AddHook(logruseq.NewSeqHook(endPoint, logruseq.OptionAPIKey(api)))
	return &LogrusSender{lg: lg}
}

func (s *LogrusSender) SendLevel(level zapcore.Level, msg string, fields logrus.Fields) error {
	entry := s.lg.WithFields(fields)
	entry.Log(mapLevel(level), msg)
	return nil
}
