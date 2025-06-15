package config

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Env      string         `yaml:"env"`
	Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
	HTTP HTTPServerConfig `yaml:"http"`
}

type HTTPServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}
