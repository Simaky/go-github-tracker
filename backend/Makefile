SVC := go-github-tracker-ms
NAME := GoGithubTracker

# The minimum Go version is pinned by the `go` directive in go.mod; the default
# GOTOOLCHAIN=auto downloads it if your local toolchain is older.

run:
	go run ./$(SVC).go

build:
	CGO_ENABLED=0 go build -ldflags "-X 'main.appVersion=$$(git rev-parse --short HEAD)'" -o $(SVC) ./$(SVC).go

test:
	go test ./... -count=1 -race

lint:
	golangci-lint run

# Build from source with the local Go toolchain so the linter is built with a
# Go >= the version this module targets (prebuilt releases lag and refuse to
# lint a newer-targeted module).
lint-install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0

# Regenerate the Ent typed client from ent/schema after editing a schema file.
ent-generate:
	go tool ent generate ./ent/schema

tidy:
	go mod tidy
