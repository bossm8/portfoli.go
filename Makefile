.devcontainer/.installed:
	/bin/bash .devcontainer/download.sh
	touch .devcontainer/.installed

setup: .devcontainer/.installed

build: setup
	test -d ${PWD}/_build || mkdir ${PWD}/_build
	go build -o ${PWD}/_build/portfoli.go portfoli.go

test:
	staticcheck ./...
	go test ./...

run: setup
	go run portfoli.go -config.dir ${PWD}/configs

dist:
	test -d ${PWD}/dist && rm -rf ${PWD}/dist 
	mkdir ${PWD}/dist
	go run portfoli.go -config.dir ${PWD}/configs -static
	/bin/bash .devcontainer/download.sh > /dev/null
	cp -r ${PWD}/public ${PWD}/dist/static
	mv ${PWD}/dist/static/favicon.ico ${PWD}/dist

clean:
	rm -rf ${PWD}/.devcontainer/.installed ${PWD}/dist ${PWD}/_build
	

.PHONY: test setup run build dist clean