package config

// Config 是應用程式的設定結構
type Config struct {
	AppName  string         `envconfig:"APP_NAME" yaml:"app_name" validate:"required"`
	Env      string         `envconfig:"ENV"      yaml:"env"      validate:"required"`
	Server   ServerConfig   `envconfig:"-"        yaml:"server"   validate:"required"`
	Database DatabaseConfig `envconfig:"-"        yaml:"database" validate:"required"`
	Auth     AuthConfig     `envconfig:"-"        yaml:"auth"     validate:"required"`
	Logger   LoggerConfig   `envconfig:"-"        yaml:"logger"   validate:"required"`
	Tracer   TracerConfig   `envconfig:"-"        yaml:"tracer"   validate:"required"`
}

type ServerConfig struct {
	HTTP HTTPConfig `envconfig:"-" yaml:"http" validate:"required"`
}

type HTTPConfig struct {
	Host string `envconfig:"SERVER_HTTP_HOST" yaml:"host" validate:"required"`
	Port string `envconfig:"SERVER_HTTP_PORT" yaml:"port" validate:"required"`
}

type DatabaseConfig struct {
	DSN string `envconfig:"DB_DSN" yaml:"dsn" validate:"required"`
}

type AuthConfig struct {
	JWT JWTConfig `envconfig:"-" yaml:"jwt" validate:"required"`
}

type JWTConfig struct {
	Algorithm string `envconfig:"JWT_ALGORITHM" yaml:"algorithm" validate:"required"`
	Secret    string `envconfig:"JWT_SECRET" yaml:"secret" validate:"required"`
	Expire    int    `envconfig:"JWT_EXPIRE" yaml:"expire" validate:"required"`
}

// LoggerConfig 定義日誌配置
type LoggerConfig struct {
	Console ConsoleLoggerConfig `envconfig:"-" yaml:"console"`
	GCP     GCPLoggerConfig     `envconfig:"-" yaml:"gcp"`
	Seq     SeqLoggerConfig     `envconfig:"-" yaml:"seq"`
}

// ConsoleLoggerConfig 定義 Console Logger 配置
type ConsoleLoggerConfig struct {
	Enabled bool   `envconfig:"LOGGER_CONSOLE_ENABLED" yaml:"enabled" default:"true"`
	Level   string `envconfig:"LOGGER_CONSOLE_LEVEL"   yaml:"level"   default:"info"`
	Format  string `envconfig:"LOGGER_CONSOLE_FORMAT"  yaml:"format"  default:"console"`
}

// GCPLoggerConfig 定義 GCP Logger 配置
type GCPLoggerConfig struct {
	Enabled   bool   `envconfig:"LOGGER_GCP_ENABLED"    yaml:"enabled"    default:"false"`
	ProjectID string `envconfig:"LOGGER_GCP_PROJECT_ID" yaml:"project_id"`
	Level     string `envconfig:"LOGGER_GCP_LEVEL"      yaml:"level"      default:"info"`
}

// SeqLoggerConfig 定義 Seq Logger 配置
type SeqLoggerConfig struct {
	Enabled              bool   `envconfig:"LOGGER_SEQ_ENABLED"  yaml:"enabled"  default:"false"`
	Endpoint             string `envconfig:"LOGGER_SEQ_ENDPOINT" yaml:"endpoint" default:"http://localhost:5341"`
	APIKey               string `envconfig:"LOGGER_SEQ_API_KEY"  yaml:"api_key"`
	Level                string `envconfig:"LOGGER_SEQ_LEVEL"    yaml:"level"    default:"info"`
	ConsoleOutputEnabled bool   `envconfig:"LOGGER_SEQ_CONSOLE_OUTPUT_ENABLED" yaml:"console_output_enabled" default:"false"`
}
// TracerConfig 定義追蹤器配置（只保留實際使用的欄位）
type TracerConfig struct {
	Enabled     bool   `envconfig:"TRACER_ENABLED" yaml:"enabled" default:"true"`
	ServiceName string `envconfig:"TRACER_SERVICE_NAME" yaml:"service_name" default:"api-server"`
}