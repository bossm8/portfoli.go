.devcontainer/.installed:
	/bin/bash .devcontainer/download.sh
	touch .devcontainer/.installed

setup: .devcontainer/.installed

dist: setup
	test -d _build || mkdir _build
	go build -o _build/portfoligo portfoli.go

test:
	staticcheck ./...
	go test ./...

run: setup
	go run portfoli.go -config.dir ${PWD}/configs

static: setup
	go run portfoli.go -config.dir ${PWD}/configs -static

.PHONY: test dist setup run