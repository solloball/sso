all:
	go run cmd/migrator/migrator.go --storage-path=./storage/sso.db --migrations-path=./migrations

