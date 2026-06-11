# syntax=docker/dockerfile:1
FROM golang:alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG APP_VERSION=docker
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags "-X 'main.appVersion=${APP_VERSION}'" -o /out/go-github-tracker-ms ./go-github-tracker-ms.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build /out/go-github-tracker-ms /app/go-github-tracker-ms
EXPOSE 12010
ENTRYPOINT ["/app/go-github-tracker-ms"]
