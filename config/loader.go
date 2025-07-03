package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	defaultEnv     = "default"
	configFilePath = "config/config.yaml"
	envPrefix      = ""
)

// Load 載入組態，優先序：struct 預設值 → YAML → ENV → 驗證
func Load() (*Config, error) {
	loadDotEnv()

	cfg := createDefaultConfig()
	fmt.Println("載入組態...")
	fmt.Printf("初始組態: %+v\n", cfg)

	if err := loadYAMLConfig(&cfg); err != nil {
		fmt.Printf("警告：%v\n", err)
	}

	if err := loadEnvConfig(&cfg); err != nil {
		fmt.Printf("警告：%v\n", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	fmt.Printf("組態載入成功: %+v\n", cfg)
	return &cfg, nil
}

func loadDotEnv() {
	_ = godotenv.Load() // .env 可有可無，載入不到也沒關係
}

func createDefaultConfig() Config {
	// 預設值只能設定非必要欄位，必要欄位不設定預設值，讓驗證能夠攔截缺漏
	return Config{
		Env: defaultEnv, // 預設環境為 "default"
		// 例如：LogLevel: "info",
		// 其餘關鍵欄位不預設，讓驗證攔下缺漏
	}
}

func loadYAMLConfig(cfg *Config) error {
	yamlData, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("找不到 %s，僅用 ENV", configFilePath)
	}

	if err := yaml.Unmarshal(yamlData, cfg); err != nil {
		return fmt.Errorf("YAML 解析失敗，僅用 ENV: %w", err)
	}

	fmt.Println("YAML組態覆蓋成功")
	fmt.Printf("YAML組態: %+v\n", *cfg)
	return nil
}

func loadEnvConfig(cfg *Config) error {
	if err := envconfig.Process(envPrefix, cfg); err != nil {
		return fmt.Errorf("環境變數處理失敗，可能有些欄位未覆蓋: %w", err)
	}

	fmt.Println("環境變數組態覆蓋成功")
	fmt.Printf("環境變數組態: %+v\n", *cfg)
	return nil
}

func validateConfig(cfg *Config) error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("設定結構驗證失敗: %w", err)
	}
	return nil
}
