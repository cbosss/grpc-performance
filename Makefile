.PHONY: proto-build proto-compile

proto-build:
	docker build . -t proto-compiler --target proto

compile: proto-build
	docker run --rm -v $(shell pwd)/proto:/proto proto-compiler protoc \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    echo.proto
