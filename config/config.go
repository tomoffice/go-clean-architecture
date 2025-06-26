package config

// Config 是應用程式的設定結構
type Config struct {
	Env      string         `envconfig:"ENV" yaml:"env" validate:"required"`
	Server   ServerConfig   `envconfig:"-"   yaml:"server" validate:"required"`
	Database DatabaseConfig `envconfig:"-"   yaml:"database" validate:"required"`
}

type ServerConfig struct {
	HTTP HTTPServerConfig `envconfig:"-" yaml:"http" validate:"required"`
}

type HTTPServerConfig struct {
	Host string `envconfig:"SERVER_HTTP_HOST" yaml:"host" validate:"required"`
	Port string `envconfig:"SERVER_HTTP_PORT" yaml:"port" validate:"required"`
}

type DatabaseConfig struct {
	DSN string `envconfig:"DB_DSN" yaml:"dsn" validate:"required"`
}
