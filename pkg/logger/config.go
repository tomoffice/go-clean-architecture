package logger


// Config 是 logger 套件的配置結構，用於初始化整個日誌系統
// 包含 OpenTelemetry 和各種 logger adapters 的配置
type Config struct {
	// OpenTelemetry 配置
	Telemetry TelemetryConfig `yaml:"telemetry" json:"telemetry"`

	// Logger Adapters 配置
	Console ConsoleConfig `yaml:"console" json:"console"`
	GCP     GCPConfig     `yaml:"gcp"     json:"gcp"`
	Seq     SeqConfig     `yaml:"seq"     json:"seq"`
}

// TelemetryConfig 定義 OpenTelemetry 的配置
type TelemetryConfig struct {
	Enabled     bool   `yaml:"enabled"      json:"enabled"      default:"true"`
	ServiceName string `yaml:"service_name" json:"service_name" default:"app"`
	Version     string `yaml:"version"      json:"version"      default:"1.0.0"`
	Environment string `yaml:"environment"  json:"environment"  default:"development"`
}

// ConsoleConfig 定義 Console Logger 的配置
type ConsoleConfig struct {
	Enabled bool   `yaml:"enabled" json:"enabled" default:"true"`
	Level   string `yaml:"level"   json:"level"   default:"info"`
	Format  string `yaml:"format"  json:"format"  default:"console"`
}

// GCPConfig 定義 GCP Cloud Logging 的配置
type GCPConfig struct {
	Enabled   bool   `yaml:"enabled"    json:"enabled"    default:"false"`
	ProjectID string `yaml:"project_id" json:"project_id"`
	Level     string `yaml:"level"      json:"level"      default:"info"`

	// 進階配置（可選）
	LogName      string            `yaml:"log_name"       json:"log_name"       default:"app-log"`
	ResourceType string            `yaml:"resource_type"  json:"resource_type"  default:"global"`
	Labels       map[string]string `yaml:"labels"         json:"labels"`
}

// SeqConfig 定義 Seq Logger 的配置
type SeqConfig struct {
	Enabled  bool   `yaml:"enabled"  json:"enabled"  default:"false"`
	Endpoint string `yaml:"endpoint" json:"endpoint" default:"http://localhost:5341"`
	APIKey   string `yaml:"api_key"  json:"api_key"`
	Level    string `yaml:"level"    json:"level"    default:"info"`
}


// NewConfig 創建新的預設配置
func NewConfig() Config {
	return Config{
		Telemetry: TelemetryConfig{
			Enabled:     false,
			ServiceName: "app",
			Version:     "1.0.0",
			Environment: "development",
		},
		Console: ConsoleConfig{
			Enabled: true,
			Level:   "info",
			Format:  "console",
		},
		GCP: GCPConfig{
			Enabled: false,
		},
		Seq: SeqConfig{
			Enabled: false,
		},
	}
}


