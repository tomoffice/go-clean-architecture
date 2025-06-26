package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
)

// Load 載入組態，優先序：struct 預設值 → YAML → ENV → 驗證
func Load() (*Config, error) {
	_ = godotenv.Load() // .env 可有可無，載入不到也沒關係

	// 預設值只能設定非必要欄位，必要欄位不設定預設值，讓驗證能夠攔截缺漏
	cfg := Config{
		Env: "default", // 預設環境為 "default"
		// 例如：LogLevel: "info",
		// 其餘關鍵欄位不預設，讓驗證攔下缺漏
	}
	fmt.Println("載入組態...")
	fmt.Printf("初始組態: %+v\n", cfg)
	// 讀取 YAML（不存在不致命，僅提醒）
	if data, err := os.ReadFile("config/config.yaml"); err != nil {
		fmt.Printf("警告：找不到 config/config.yaml，僅用 ENV\n")
	} else {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			fmt.Printf("警告：YAML 解析失敗，僅用 ENV: %v\n", err)
		} else {
			fmt.Println("YAML組態覆蓋成功")
			fmt.Printf("YAML組態: %+v\n", cfg)
		}
	}

	// 環境變數覆蓋（有就覆蓋，無則略過）
	if err := envconfig.Process("", &cfg); err != nil {
		fmt.Printf("警告：環境變數處理失敗，可能有些欄位未覆蓋: %v\n", err)
	} else {
		fmt.Println("環境變數組態覆蓋成功")
		fmt.Printf("環境變數組態: %+v\n", cfg)
	}

	// 驗證所有 required 欄位
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("設定結構驗證失敗: %w", err)
	}

	fmt.Printf("組態載入成功: %+v\n", cfg)
	return &cfg, nil
}
