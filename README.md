# Module Clean - Clean Architecture Go 專案

這是一個遵循 Clean Architecture 原則的 Go 語言專案，展示了如何使用物件導向設計、依賴注入和模組化架構來建立可維護、可測試的應用程式。

## 📋 目錄

- [特色](#特色)
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

## ✨ 特色

- 🏗️ **Clean Architecture** - 嚴格遵循 Clean Architecture 原則
- 💉 **依賴注入** - 所有元件依賴抽象介面，而非具體實作
- 📦 **模組化設計** - 高內聚低耦合的模組架構
- 🧪 **高度可測試** - 透過介面隔離實現單元測試
- 🔄 **資料庫遷移** - 使用 migrate 工具管理資料庫版本
- 📝 **程式碼風格** - 遵循 [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

## 📂 專案結構

```
module-clean/
├── cmd/                      # 應用程式進入點
│   └── api-server/          # API 伺服器執行檔
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
│   │   └── http/            # HTTP 框架整合
│   ├── modules/             # 業務模組
│   │   ├── member/          # 會員模組 (範例)
│   │   │   ├── entity/              # 企業規則層
│   │   │   ├── usecase/             # 應用業務規則層
│   │   │   ├── interface_adapter/   # 介面配接層
│   │   │   │   ├── controller/      # 控制器
│   │   │   │   ├── gateway/         # 閘道器（外部服務整合）
│   │   │   │   ├── presenter/       # 呈現器
│   │   │   │   └── router/          # 路由設定
│   │   │   └── driver/              # 框架與驅動層
│   │   │       └── persistence/     # 持久化實作
│   │   └── modules.go       # 模組介面定義
│   └── shared/              # 共用元件
│       ├── apperror/        # 應用程式錯誤
│       ├── errorcode/       # 錯誤碼定義
│       └── pagination/      # 分頁功能
├── pkg/                      # 公開的程式庫
│   └── logger/              # 日誌套件
├── migrations/               # 資料庫遷移檔案
├── seed/                     # 資料庫種子資料
├── data/                     # 本地資料儲存
├── makefile                  # 建構自動化
├── go.mod                    # Go 模組定義
└── go.sum                    # Go 模組相依性鎖定
```

## 🛠 技術棧

- **語言**: Go 1.23+
- **Web 框架**: [Gin](https://github.com/gin-gonic/gin)
- **資料庫**: SQLite (可替換為 MySQL/PostgreSQL)
- **ORM/查詢建構器**: 
  - [sqlx](https://github.com/jmoiron/sqlx) - SQL 查詢
  - [Ent](https://entgo.io/) - ORM (可選)
- **驗證**: [validator/v10](https://github.com/go-playground/validator)
- **設定管理**: 
  - [envconfig](https://github.com/kelseyhightower/envconfig)
  - [godotenv](https://github.com/joho/godotenv)
  - YAML 設定檔
- **測試**: 
  - [testify](https://github.com/stretchr/testify)
  - [gomock](https://github.com/golang/mock)

## 🚀 開始使用

### 環境需求

- Go 1.23 或更高版本
- SQLite3
- Make
- [migrate](https://github.com/golang-migrate/migrate) CLI 工具

### 安裝步驟

1. **複製專案**
   ```bash
   git clone <repository-url>
   cd module-clean
   ```

2. **安裝相依套件**
   ```bash
   go mod download
   ```

3. **設定環境變數**
   ```bash
   cp .env.example .env
   # 編輯 .env 檔案設定您的環境變數
   ```

4. **執行資料庫遷移**
   ```bash
   make db-migrate
   ```

5. **載入種子資料（可選）**
   ```bash
   make db-seed
   ```

6. **執行應用程式**
   ```bash
   go run cmd/api-server/main.go
   ```

### 設定

應用程式使用分層設定系統：

1. **預設值** - 在程式碼中定義
2. **設定檔** - `config/config.yaml`
3. **環境變數** - 透過 `.env` 檔案或系統環境變數
4. **命令列參數** - 執行時期覆寫

設定優先順序：命令列 > 環境變數 > 設定檔 > 預設值

## 👨‍💻 開發指南

### 架構概覽

本專案嚴格遵循 Clean Architecture 的同心圓架構：

1. **Entities (實體層)** - 企業業務規則
   - 位置：`internal/modules/*/entity/`
   - 包含核心業務邏輯和規則

2. **Use Cases (用例層)** - 應用業務規則
   - 位置：`internal/modules/*/usecase/`
   - 協調實體和外部系統的互動

3. **Interface Adapters (介面配接層)** - 資料轉換
   - Controllers：處理輸入
   - Gateways：外部服務整合
   - Presenters：格式化輸出
   - 位置：`internal/modules/*/interface_adapter/`

4. **Frameworks & Drivers (框架與驅動層)** - 外部工具
   - 位置：`internal/modules/*/driver/`
   - 資料庫、Web 框架等具體實作

### 新增模組

1. **建立模組目錄結構**
   ```bash
   mkdir -p internal/modules/your_module/{entity,usecase,interface_adapter,driver}
   ```

2. **定義實體**
   ```go
   // internal/modules/your_module/entity/your_entity.go
   package entity

   type YourEntity struct {
       ID   string
       Name string
   }
   ```

3. **實作用例**
   ```go
   // internal/modules/your_module/usecase/your_usecase.go
   package usecase

   type YourUseCase struct {
       repo YourRepository
   }
   ```

4. **建立介面配接器**
   - Controller
   - Gateway
   - Presenter
   - Router

5. **註冊模組**
   ```go
   // internal/modules/your_module/your_module.go
   package your_module

   type YourModule struct {
       // 依賴注入
   }

   func (m *YourModule) Assemble() {
       // 組裝元件
   }
   ```

### 測試

執行所有測試：
```bash
go test ./...
```

執行特定模組測試：
```bash
make test-member
```

產生 Mock：
```bash
make generate
```

## 📝 指令說明

| 指令 | 說明 |
|------|------|
| `make help` | 顯示所有可用指令 |
| `make db-migrate` | 執行資料庫遷移 |
| `make db-seed` | 載入種子資料 |
| `make db-reset` | 重置資料庫（清除、遷移、種子） |
| `make clean` | 清除產生的檔案 |
| `make generate` | 產生 Mock 和其他自動產生的程式碼 |
| `make tree` | 複製目錄結構到剪貼簿 |

## 📖 API 文件

### 會員管理 API

#### 取得會員列表
```http
GET /members?page=1&pageSize=10
```

#### 取得單一會員
```http
GET /members/:id
```

#### 透過 Email 取得會員
```http
GET /members/email/:email
```

#### 建立會員
```http
POST /members
Content-Type: application/json

{
  "email": "user@example.com",
  "name": "User Name"
}
```

#### 更新會員資料
```http
PATCH /members/:id
Content-Type: application/json

{
  "name": "Updated Name",
  "phone": "+886912345678"
}
```

#### 更新會員信箱
```http
PATCH /members/:id/email
Content-Type: application/json

{
  "email": "newemail@example.com"
}
```

#### 更新會員密碼
```http
PATCH /members/:id/password
Content-Type: application/json

{
  "old_password": "oldPassword123",
  "new_password": "newPassword456"
}
```

#### 刪除會員
```http
DELETE /members/:id
```

