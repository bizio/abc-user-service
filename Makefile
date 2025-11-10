PKGS := $(shell go list ./... | grep -e domain -e application)

# Build the project
binary:
	CGO_ENABLED=0 GOOS=linux go build -o bin/abc-user-service ./cmd/server/main.go

docker-compose-up:
	docker compose -f build/development/docker-compose.yml up --build

docker-compose-down:
	docker compose -f build/development/docker-compose.yml down

# Run the server
run-local: build
	scripts/run_local.sh
	
# Run tests
test: gomod 
	mkdir -p out
	export PLAYBACK=true; \
	go test -short -coverprofile out/cover.out $(PKGS)
	go tool cover -html=out/cover.out -o out/cover.html
	go tool cover -func=out/cover.out

lint:
	golangci-lint run cmd/... pkg/... internal/...

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# add missing and remove unused modules
gomod:
	go mod tidy -v

# Generate mocks
mocks: gomod
	rm -fR mocks/*
	go generate ./...

# Update mocks
update-mocks:
	# ensure mockery is updated
	go install github.com/vektra/mockery/v2@latest

	mocks

# Run tests with mocks
test-mocks: update-mocks test

# Clean up build artifacts
clean:
	rm -rf bin/*

install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@v2.10.2

install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

swagger: install-swag
	swag init -g cmd/server/main.go

