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

lint-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.5.0

# Regenerate the Ent typed client from ent/schema after editing a schema file.
ent-generate:
	go tool ent generate ./ent/schema

tidy:
	go mod tidy
