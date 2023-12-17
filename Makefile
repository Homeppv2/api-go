vprotoc:
	protoc -I. --go-grpc_out=. --go_out=. pkg/proto/$(file).proto

lint:
	golangci-lint run --config=./.golangci.yml

migrate-new:
	migrate create -ext sql -dir db/migration -seq $(name)

migrate-up:
	migrate -path db/migration \
	-database "postgresql://root:pass@127.0.0.1:5432/api?sslmode=disable" \
	-verbose up

migrate-down:
	migrate -path db/migration \
  -database "postgresql://root:pass@127.0.0.1:5432/api?sslmode=disable" \
  -verbose down

up:
	docker compose up -d