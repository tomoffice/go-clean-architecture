package config

// Config 是應用程式的設定結構
type Config struct {
	Env      string         `envconfig:"ENV"    yaml:"env"      validate:"required"`
	Server   ServerConfig   `envconfig:"-"      yaml:"server"   validate:"required"`
	Database DatabaseConfig `envconfig:"-"      yaml:"database" validate:"required"`
	Auth     AuthConfig     `envconfig:"-"      yaml:"auth"     validate:"required"`
}

type ServerConfig struct {
	HTTP HTTPConfig `envconfig:"-"    yaml:"http"    validate:"required"`
}

type HTTPConfig struct {
	Host string `envconfig:"SERVER_HTTP_HOST"    yaml:"host"    validate:"required"`
	Port string `envconfig:"SERVER_HTTP_PORT"    yaml:"port"    validate:"required"`
}

type DatabaseConfig struct {
	DSN string `envconfig:"DB_DSN"    yaml:"dsn"    validate:"required"`
}

type AuthConfig struct {
	JWT JWTConfig `envconfig:"-"    yaml:"jwt"    validate:"required"`
}
type JWTConfig struct {
	Algorithm string `envconfig:"JWT_ALGORITHM"    yaml:"algorithm" validate:"required"`
	Secret    string `envconfig:"JWT_SECRET"      yaml:"secret"    validate:"required"`
	Expire    int    `envconfig:"JWT_EXPIRE"      yaml:"expire"    validate:"required"`
}
