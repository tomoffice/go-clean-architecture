pkg/logger/
├── go.mod                          # 可獨立成 go module
├── logger.go                       # 抽象介面 Logger / Field
├── field.go                        # Field 工具函式（String, Int, Error...）
├── level.go                        # 日誌等級常數定義（InfoLevel 等）
├── errors.go                       # Logger 專用錯誤（如 ErrNoAdapters）
├── config.go                       # LoggerConfig 結構（供 factory 使用）
├── factory/
│   └── factory.go                  # 依 config 組合 adapter，產出 Logger 實體
├── adapters/
│   ├── console/
│   │   └── console_logger.go       # consoleLogger 實作（stdout）
│   ├── gcp/
│   │   └── gcp_logger.go           # gcpLogger 實作（Cloud Logging）
│   ├── seq/
│   │   └── seq_logger.go           # seqLogger 實作（Seq 接收器）
│   └── multi/
│       └── multi_logger.go        # multiLogger：包多個 Logger 的組合器