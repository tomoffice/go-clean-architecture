//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"github.com/tomoffice/go-clean-architecture/pkg/logger/adapters/seq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func setupOpenTelemetry() func() error {
	tp := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	return func() error {
		return tp.Shutdown(context.Background())
	}
}

func createContextWithTrace() context.Context {
	tracer := otel.Tracer("my-service")
	ctx, _ := tracer.Start(context.Background(), "user-request",
		trace.WithAttributes(
			attribute.String("user.id", "12345"),
			attribute.String("request.path", "/api/users"),
		),
	)
	return ctx
}

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

func TestSeqWithOpenTelemetry(t *testing.T) {
	// 設置 OpenTelemetry
	cleanup := setupOpenTelemetry()
	defer cleanup()

	// 創建 logger
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

	tracer := otel.Tracer("test-service")

	// Main task
	mainCtx, mainSpan := tracer.Start(context.Background(), "main-task")
	lg.WithContext(mainCtx).Info("Main task started")

	// Subtask 1
	sub1Ctx, sub1Span := tracer.Start(mainCtx, "subtask-1")
	lg.WithContext(sub1Ctx).Info("Subtask 1 processing")
	sub1Span.End()

	// Subtask 2
	sub2Ctx, sub2Span := tracer.Start(mainCtx, "subtask-2")
	lg.WithContext(sub2Ctx).Info("Subtask 2 processing", logger.NewField("subtask", "data processing"))
	sub2Span.End()

	// Main task completed
	lg.WithContext(mainCtx).Info("Main task completed")
	mainSpan.End()
	// Main2 task
	main2Ctx, main2Span := tracer.Start(context.Background(), "main2-task")
	object := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	lg.WithContext(main2Ctx).Info("Main2 task started", logger.NewField("obj", object))
	// Main2 task completed
	lg.WithContext(main2Ctx).Warn("Main2 task completed")
	main2Span.End()

	t.Log("=== 測試完成 ===")
	t.Log("請檢查 Seq dashboard (http://localhost:5341) 查看日誌")
	t.Log("您應該會看到所有日誌都有相同的 trace_id，但不同的操作有不同的 span_id")
}
