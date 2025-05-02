DB_FILE=identifier.sqlite
MIGRATE_SQL=migrations/sqlite_init.sql
SEED_SQL=migrations/sqlite_seed.sql
ent:
	go install entgo.io/ent/cmd/ent@latest
ent-generate:
	ent generate ./internal/infrastructure/db/ent/schema
tree:
	tree -I 'ent*' | pbcopy
sqlite-migrate:
	sqlite3 $(DB_FILE) < $(MIGRATE_SQL)

sqlite-seed:
	sqlite3 $(DB_FILE) < $(SEED_SQL)

