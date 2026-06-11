SVC := go-github-tracker-ms
NAME := GoGithubTracker

# This repo targets the latest Go (1.26) via the toolchain directive in go.mod.
# Toolchain auto-download verifies against the checksum database; if your machine
# has GOSUMDB=off it will refuse, so we re-enable it for these targets. Drop the
# GOSUMDB line if your environment already resolves the toolchain.
export GOTOOLCHAIN := go1.26.0
export GOSUMDB := sum.golang.org

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

# Regenerate the Ent client. Requires at least one schema under ent/schema.
ent-generate:
	go generate ./ent

tidy:
	go mod tidy
