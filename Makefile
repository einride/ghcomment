.PHONY: all
all: \
	go-test \
	go-build \
	go-mod-tidy

.PHONY: go-test
go-test:
	go test -race -timeout 30s ./...

.PHONY: go-build
go-build:
	GOOS=darwin go build ./...
	GOOS=windows go build ./...
	GOOS=linux go build ./...

.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy
