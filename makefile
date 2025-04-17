ent:
	go install entgo.io/ent/cmd/ent@latest
ent-generate:
	ent generate ./internal/infrastructure/db/ent/schema
