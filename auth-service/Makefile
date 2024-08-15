generate-proto:
	protoc -I api api/proto/auth.proto \
	--go_out=./api/gen/ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./api/gen/ \
	--go-grpc_opt=paths=source_relative

run-local:
	go run ./cmd/main.go --config=./config/local.yaml
