export PATH := $(PWD)/build:$(PATH)

start-server: build run-start-server

run-start-server:
	muon server start

build: build-cli

build-cli:
	go build -o build/muon ./cmd/cli

go-get-muto:
	go get github.com/SSripilaipong/muto@$(shell curl -s https://api.github.com/repos/SSripilaipong/muto/commits/main | jq -r .sha)

go-get-common:
	go get github.com/SSripilaipong/go-common@$(shell curl -s https://api.github.com/repos/SSripilaipong/go-common/commits/main | jq -r .sha)
