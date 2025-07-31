package tracer

// Config 是 tracer 套件的配置結構，用於初始化追蹤系統
// 遵循與 logger.Config 相同的設計模式
type Config struct {
	// ServiceName 服務名稱，用於識別不同的服務
	ServiceName string `yaml:"service_name" json:"service_name" default:"unknown-service"`

	// Version 服務版本，用於追蹤不同版本的效能表現
	Version string `yaml:"version" json:"version" default:"1.0.0"`

	// Environment 運行環境 (development, staging, production)
	Environment string `yaml:"environment" json:"environment" default:"development"`

	// Enabled 是否啟用追蹤功能
	// 在開發或測試環境可能需要關閉以提升效能
	Enabled bool `yaml:"enabled" json:"enabled" default:"true"`
}

// DefaultConfig 返回預設配置
func DefaultConfig() Config {
	return Config{
		ServiceName: "unknown-service",
		Version:     "1.0.0",
		Environment: "development",
		Enabled:     true,
	}
}

// NewConfig 創建新的配置
func NewConfig(serviceName, version, environment string, enabled bool) Config {
	return Config{
		ServiceName: serviceName,
		Version:     version,
		Environment: environment,
		Enabled:     enabled,
	}
}
