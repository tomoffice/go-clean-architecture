package seq

//go:generate mockgen -source=sender.go -destination=mock/mock_sender.go -package=mock
import (
	"github.com/sirupsen/logrus"
	"go.uber.org/zap/zapcore"
)

// Sender 負責把日誌送到 Seq 的抽象介面
type Sender interface {
	SendLevel(level zapcore.Level, msg string, fields logrus.Fields) error
}
