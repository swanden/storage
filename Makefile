build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

integration-tests:
	docker-compose run --rm storage go run ./test/integration

unit-tests:
	docker-compose run --rm storage go test -v -count=1 ./...

build-proto:
	protoc \
	-I=./api/proto storage.proto \
    --go_out=./api/go \
    --go_opt=module=storage/api \
    --go-grpc_out=./api/go \
    --go-grpc_opt=module=storage/api && \
    protoc-go-inject-tag -input="./api/go/*.pb.go"