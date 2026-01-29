# ====== variables ======
PROTO_DIR=proto

# ====== commands ======
.PHONY: generate build run test docker-build docker-up

generate:
	protoc \
		-I proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=module=github.com/maximfill/go-pet-backend \
		--go-grpc_opt=module=github.com/maximfill/go-pet-backend \
		proto/*.proto


build:
	go build ./...

run:
	go run ./cmd/api/main.go

test:
	go test ./...

docker-build:
	docker compose build

docker-up:
	docker compose up

generate-postman:
	protoc \
  -I proto \
  --go_out=./internal/transport/grpc \
  --go-grpc_out=./internal/transport/grpc \
  --go_opt=module=github.com/maximfill/go-pet-backend \
  --go-grpc_opt=module=github.com/maximfill/go-pet-backend \
  proto/*.proto
