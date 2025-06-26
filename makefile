# ------------------------------------------------------------------------------
#  Configuration
# ------------------------------------------------------------------------------

# 讀取 .env 檔案並將其變數匯出給後續指令使用
ifneq (,$(wildcard .env))
	include .env
	export
endif

# 若 .env 或外部未傳入 ENV，則預設為 development
ENV ?= development
ENVIRONMENT := $(ENV)

# 工具與路徑設定 (集中管理，方便修改)
MIGRATE_CLI     := migrate

# 目錄路徑 (直接使用 make 內建的 CURDIR)
DATA_DIR        := $(CURDIR)/data
MIGRATIONS_DIR  := $(CURDIR)/migrations
SEED_DIR        := $(CURDIR)/seed
ENT_SCHEMA_DIR  := $(CURDIR)/internal/modules/member/driver/persistence/ent/schema # 調整為你實際的 Ent Schema 路徑

# 資料庫設定
DB_FILE         := $(DATA_DIR)/$(ENVIRONMENT).sqlite
DB_URL          := sqlite3://$(DB_FILE)

# 取得所有 seed SQL 檔案，而非寫死單一檔案
SEED_FILES      := $(wildcard $(SEED_DIR)/*.sql)


# ------------------------------------------------------------------------------
#  Targets
# ------------------------------------------------------------------------------

.PHONY: all help ent-generate tree db-migrate db-seed db-reset clean

# 設定預設指令為 help
.DEFAULT_GOAL := help

# 自動產生說明指令
# ## @target-name: 說明文字
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

ent:
	go install entgo.io/ent/cmd/ent@latest

ent-generate:
	ent generate ./internal/infrastructure/db/ent/schema

## tree: 複製目錄結構到剪貼簿 (macOS)
tree:
	@tree -I 'ent*|data|*.sqlite' | pbcopy
	@echo "目錄結構已複製到剪貼簿（已排除 ent, data, sqlite 檔案）"

# 確保資料庫目錄存在 (這是一個隱藏的輔助 target)
$(DATA_DIR):
	@mkdir -p $(DATA_DIR)

## db:migrate: 執行資料庫遷移
db-migrate: $(DATA_DIR)
	@echo "=> Running migrations on $(DB_FILE)..."
	@$(MIGRATE_CLI) -database "$(DB_URL)" -path "$(MIGRATIONS_DIR)" up

## db:seed: 執行資料庫填充
db-seed:
	@if [ -z "$(SEED_FILES)" ]; then \
		echo "No seed files found in $(SEED_DIR). Skipping."; \
	else \
		echo "=> Seeding database $(DB_FILE)..."; \
		for f in $(SEED_FILES); do \
			echo "   - Applying $$f"; \
			sqlite3 "$(DB_FILE)" < "$$f"; \
		done; \
	fi

## db:reset: 重置資料庫 (清除 > 遷移 > 填充)
db-reset: clean db-migrate db-seed
	@echo "=> Database reset complete."

## clean: 清除產生的資料庫檔案
clean:
	@echo "=> Cleaning up generated files..."
	@rm -rf $(DATA_DIR)

## generate: 產生所有 go:generate 的程式碼 (例如 mocks)
generate:
	@echo "=> Generating code..."
	@go generate ./...