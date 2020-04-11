.PHONY: build
build:
	@go build -o build/api cmd/api/main.go
	@go build -o build/readreddit cmd/cli/main.go

.PHONY: build_linux
build_linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api_linux cmd/api/main.go
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/readreddit_linux cmd/cli/main.go

clean:
	@go clean
	@rm -rf build

api: build
	@./build/api

readreddit: build
	@./build/readreddit

.PHONY: test
test:
	@go test ./...

.PHONY: get
get:
	@go get ./...

deploy: build_linux
	@ansible-playbook -i infra/ansible/hosts infra/ansible/floridaman-deploy.yml

.PHONY: format
format:
	@goimports -local github.com/davidonium/floridaman -w .

.PHONY: lint
lint:
	@goimports -local github.com/davidonium/floridaman -l .