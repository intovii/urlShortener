psql:
	echo STORAGE_TYPE=PSQL>.env
	docker compose --profile psql up --build -d

in_memory:
	echo STORAGE_TYPE=InMemo>.env
	docker compose --profile in_memory up --build -d

tests:
	go test -cover -v ./internal/usecase

remove:
	docker compose down url_shortener db_urls in_memory -v
