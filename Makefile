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

static:
	test -d ${PWD}/dist && rm -rf ${PWD}/dist 
	mkdir ${PWD}/dist
	go run portfoli.go -config.dir ${PWD}/configs -static
	/bin/bash .devcontainer/download.sh > /dev/null
	cp -r ${PWD}/public ${PWD}/dist/static
	mv ${PWD}/dist/static/favicon.ico ${PWD}/dist
	

.PHONY: test dist setup run static