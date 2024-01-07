migrate:
	go run ./cmd/migrator --storage-path=./storage/texts.db --migrations-path=./migrations

run:
	go run ./cmd/texts --config=./config/config.yaml
