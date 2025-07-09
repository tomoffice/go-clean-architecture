//go:build integration
// +build integration

package integration

import (
	"testing"
	"time"

	"github.com/tomoffice/logger"
	"github.com/tomoffice/logger/zap_adapters/seq"
)

func TestSendToSeq(t *testing.T) {
	cfg := seq.Config{
		Endpoint: "http://localhost:5341",
		APIKey:   "",
		Level:    logger.InfoLevel,
	}
	lg, err := seq.NewLogger(cfg)
	if err != nil {
		t.Fatalf("Seq logger 建立失敗: %v", err)
	}
	defer lg.Sync()

	lg.Info("ping", logger.NewField("env", "integration"))
	lg.Info("string test", logger.NewField("key1", "my string"))
	lg.Info("int test", logger.NewField("key2", 12345))
	lg.Info("bool test", logger.NewField("key3", true))
	lg.Info("float test", logger.NewField("key4", 3.14159))
	lg.Info("duration test", logger.NewField("key5", 2*time.Second))
	lg.Info("time test", logger.NewField("key6", time.Now().Format(time.RFC3339Nano)))
	lg.Info("object test", logger.NewField("key7", map[string]interface{}{
		"nestedKey1": "nestedValue1",
		"nestedKey2": 42,
	}))
	lg.Debug("debug not sent", logger.NewField("debugKey", "debugValue")) // Debug level 不會被發送

	// 給 Hook 一點時間 flush；正式專案可改用 WaitGroup 或 Hook 回傳確認
	//time.Sleep(500 * time.Millisecond)
}
