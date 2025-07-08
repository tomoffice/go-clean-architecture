# Module Clean - Go Clean Architecture

基於 Clean Architecture 的 Go 專案，實作分層架構模組化設計。

## 目錄

- [專案結構](#專案結構)
- [技術棧](#技術棧)
- [開始使用](#開始使用)
  - [環境需求](#環境需求)
  - [安裝步驟](#安裝步驟)
  - [設定](#設定)
- [開發指南](#開發指南)
  - [架構概覽](#架構概覽)
  - [新增模組](#新增模組)
  - [測試](#測試)
- [指令說明](#指令說明)

## 專案結構

```
module-clean/
├── cmd/                      # 應用程式進入點
│   └── api-server/          # API 伺服器
│       └── main.go
├── config/                   # 設定檔管理
│   ├── config.go            # 設定結構定義
│   ├── config.yaml          # 設定檔
│   └── loader.go            # 設定載入邏輯
├── internal/                 # 私有應用程式程式碼
│   ├── bootstrap/           # 應用程式初始化
│   │   └── bootstrap.go
│   ├── framework/           # 框架相關實作
│   │   ├── database/        # 資料庫連線管理
│   │   │   └── sqlx/        # SQLx 實作
│   │   └── http/            # HTTP 框架整合
│   │       └── gin/         # Gin 框架實作
│   ├── modules/             # 業務模組
│   │   ├── member/          # 會員模組 (範例)
│   │   │   ├── entity/              # 企業規則層
│   │   │   ├── usecase/             # 應用業務規則層
│   │   │   ├── interface_adapter/   # 介面配接層
│   │   │   │   ├── controller/      # 控制器
│   │   │   │   ├── gateway/         # 閘道器（外部服務整合）
│   │   │   │   ├── presenter/       # 呈現器
│   │   │   │   └── router/          # 路由設定
│   │   │   └── framework/           # 框架與驅動層
│   │   │       └── persistence/     # 持久化實作
│   │   └── modules.go       # 模組介面定義
│   └── shared/              # 共用元件
│       ├── apperror/        # 應用程式錯誤
│       ├── errorcode/       # 錯誤碼定義
│       ├── errordefs/       # 錯誤定義
│       ├── enum/            # 列舉定義
│       ├── pagination/      # 分頁功能
│       └── viewmodel/       # 視圖模型
├── pkg/                     # 公用函式庫
│   └── logger/             # 日誌套件 (獨立模組)
├── migrations/             # 資料庫遷移腳本
├── seed/                   # 種子資料
├── data/                   # 本地資料存放
└── docs/                   # 文件
```

## 技術棧

**核心框架**
- Go 1.23+
- Gin (HTTP 框架)
- SQLite (預設資料庫)

**資料層**
- SQLx (SQL 工具包)
- golang-migrate (資料庫遷移)
- Ent (可選 ORM)

**HTTP 層**
- Gin binding (請求綁定)
- DTO 轉換
- 錯誤映射

**中間件**
- CORS 支援
- 認證中間件
- 錯誤處理

**設定管理**
- YAML 設定檔
- 環境變數覆蓋
- 多環境支援

**日誌系統**
- 結構化日誌
- 多種輸出適配器 (Console, GCP, Seq)
- 可配置日誌等級

**測試工具**
- 內建測試框架
- Mock 生成
- 單元測試覆蓋

## 開始使用

### 環境需求

- Go 1.23+
- SQLite3
- Make
- migrate CLI

### 安裝步驟

1. **複製專案**
   ```bash
   git clone <repository-url>
   cd module-clean
   ```

2. **安裝依賴**
   ```bash
   go mod download
   ```

3. **設定環境**
   ```bash
   cp .env.example .env
   # 編輯 .env 檔案設定資料庫路徑等
   ```

4. **執行資料庫遷移**
   ```bash
   make db-migrate
   ```

5. **載入種子資料**
   ```bash
   make db-seed
   ```

6. **啟動服務**
   ```bash
   go run cmd/api-server/main.go
   ```

   服務預設啟動在 `http://localhost:8080`

### 設定

設定優先順序：CLI 參數 > 環境變數 > config.yaml > 預設值

- `config/config.yaml` - 主要設定檔
- `.env` - 環境變數
- 支援 development/staging/production 環境

## 開發指南

### 架構概覽

專案嚴格遵循 Clean Architecture 四層架構：

1. **Entities** (`entity/`)
   - 核心業務邏輯與規則
   - 完全獨立於外部依賴

2. **Use Cases** (`usecase/`)
   - 應用程式特定業務邏輯
   - 定義 port 介面
   - 協調 entities 與外部服務

3. **Interface Adapters** (`interface_adapter/`)
   - `controller/` - HTTP 請求處理
   - `gateway/` - 外部服務整合
   - `presenter/` - 回應格式化
   - `dto/` - 資料傳輸物件
   - `mapper/` - 資料轉換

4. **Framework & Drivers** (`framework/`)
   - 資料庫實作 (`persistence/`)
   - HTTP 框架整合
   - 外部工具整合

### 新增模組

1. **建立模組結構**
   ```bash
   internal/modules/your_module/
   ├── entity/
   │   ├── your_entity.go
   │   └── errors.go
   ├── usecase/
   │   ├── your_usecase.go
   │   ├── errors.go
   │   ├── port/
   │   ├── inputmodel/
   │   └── mock/
   ├── interface_adapter/
   │   ├── controller/
   │   ├── dto/
   │   ├── gateway/
   │   ├── mapper/
   │   ├── presenter/
   │   └── router/
   └── framework/
       └── persistence/
   ```

2. **實作步驟**
   - 定義 entity 與核心業務邏輯
   - 建立 usecase 與 port 介面
   - 實作 controller 處理 HTTP 請求
   - 建立 gateway 處理外部整合
   - 實作 persistence 層處理資料存取

3. **模組註冊**
   - 更新 `modules.go` ��冊新模組
   - 設定依賴注入
   - 註冊路由

### 測試

- **執行全部測試**
  ```bash
  go test ./...
  ```

- **測試特定模組**
  ```bash
  go test ./internal/modules/member/...
  ```

- **執行測試並查看覆蓋率**
  ```bash
  go test -cover ./...
  ```

- **產生 mock 檔案**
  ```bash
  make generate
  ```

## 指令說明

| 指令 | 說明 |
|------|------|
| `make help` | 顯示所有可用指令 |
| `make db-migrate` | 執行資料庫遷移 |
| `make db-seed` | 載入種子資料到資料庫 |
| `make db-reset` | 重置資料庫 (清空+遷移+種子) |
| `make ent-generate` | 產生 Ent ORM 程式碼 |
| `make tree` | 顯示專案目錄結構 |
| `make clean` | 清理暫存檔案 |
