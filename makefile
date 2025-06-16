ENVIRONMENT=dev
DB_FILE=/tmp/$(ENVIRONMENT).sqlite
MIGRATE_SQL=$(CURDIR)/migrations/$(ENVIRONMENT)/sqlite_init.sql
SEED_SQL=$(CURDIR)/migrations/$(ENVIRONMENT)/sqlite_seed.sql

ent:
	go install entgo.io/ent/cmd/ent@latest

ent-generate:
	ent generate ./internal/infrastructure/db/ent/schema

tree:
	tree -I 'ent*' | pbcopy

.PHONY: sqlite-migrate sqlite-seed sqlite-reset

sqlite-migrate:
	@sqlite3 $(DB_FILE) < $(MIGRATE_SQL)

sqlite-seed:
	@sqlite3 $(DB_FILE) < $(SEED_SQL)

sqlite-reset:
	@rm -f $(DB_FILE)
	@$(MAKE) sqlite-migrate
	@$(MAKE) sqlite-seed