include .env
LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-note-api

generate-note-api:
	mkdir -p pkg/note_v1
	protoc --proto_path api/note_v1 \
	--go_out=pkg/note_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/note_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/note_v1/note.proto

local-migration-status:
	$(LOCAL_BIN)/goose ${MIGRATION_DIR} postgres ${PG_DSN} status -v


local-migration-up:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v


local-migration-down:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v


#build:
#	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/grpc_server/main.go
#
#copy-to-server:
#	scp service_linux root@31.129.59.145:
#
#docker-build-and-push:
#	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/xeeetu/test-server:v0.0.1 .
#	docker login -u token -p CRgAAAAAKruUnuikRyAI4uy09y7gZTzEbEf6L4Eo cr.selcloud.ru/xeeetu
#	docker push cr.selcloud.ru/xeeetu/test-server:v0.0.1


