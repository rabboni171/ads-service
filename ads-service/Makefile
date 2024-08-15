
cover:
	go test -short -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

swag_init:
	swag init -g internal/controllers/server.go