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
	GOOS=darwin go build -o /dev/null ./...
	GOOS=windows go build -o /dev/null ./...
	GOOS=linux go build -o /dev/null ./...

.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy
