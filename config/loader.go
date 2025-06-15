package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func LoadConfig() *Config {
	cfg := &Config{}
	// 1. 載入 YAML
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("讀取 YAML 設定失敗: %v", err)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Fatalf("解析 YAML 失敗: %v", err)
	}

	// 2. 載入 .env 覆寫
	_ = godotenv.Load("config/.env")
	envOverrides(cfg)

	return cfg
}

func envOverrides(cfg *Config) {
	if val := os.Getenv("ENV"); val != "" {
		cfg.Env = val
		fmt.Printf("環境變數 ENV: %s\n", cfg.Env)
	}
	if val := os.Getenv("SERVER_HOST"); val != "" {
		cfg.Server.HTTP.Host = val
		fmt.Printf("環境變數 SERVER_HOST: %s\n", cfg.Server.HTTP.Host)
	}
	if val := os.Getenv("SERVER_PORT"); val != "" {
		cfg.Server.HTTP.Port = val
		fmt.Printf("環境變數 SERVER_PORT: %s\n", cfg.Server.HTTP.Port)
	}
	if val := os.Getenv("DB_DSN"); val != "" {
		cfg.Database.DSN = val
		fmt.Printf("環境變數 DB_DSN: %s\n", cfg.Database.DSN)
	}
}
