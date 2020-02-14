.PHONY: build
build:
	go build -o build/api cmd/api.go
	go build -o build/readreddit cmd/readreddit.go

build_linux: build/api_linux build/readreddit_linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api_linux cmd/api.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/readreddit_linux cmd/readreddit.go

api: build
	./build/api

readreddit: build
	./build/readreddit

.PHONY: test
test:
	go test ./...

.PHONY: get
get:
	go get ./...

deploy: build_linux
	ansible-playbook -i infra/ansible/hosts infra/ansible/floridaman-deploy.yml