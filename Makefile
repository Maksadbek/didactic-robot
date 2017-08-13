all:
	$(MAKE) deps
	$(MAKE) install	

deps:
	go get ./...

install:
	go install ./cmd/drobot

build:
	go build ./cmd/drobot/main.go

lint:
	go vet ./...

test-all:
	go test ./...

test:
	go test -short ./...
