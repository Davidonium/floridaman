.PHONY: build
build:
	go build -o build/api cmd/api.go
	go build -o build/readreddit cmd/readreddit.go

.PHONY: runapi
api: build
	./build/api

.PHONY: readreddit
readreddit: build
	./build/readreddit

.PHONY: test
test:
	go test ./...

.PHONY: get
get:
	go get ./...

