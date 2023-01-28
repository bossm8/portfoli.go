PATHARGS=\
	-config.dir ${PWD}/examples/configs \
	-static.dir ${PWD}/public \
	-templates.dir ${PWD}/templates
 
build: setup _build/portfoli-go

_build/portfoli-go: **/*.go
	test -d _build || mkdir _build
	go build -o _build/portfoli-go portfoli.go

setup: .devcontainer/.installed

.devcontainer/.installed: .devcontainer/download.sh
	/bin/bash .devcontainer/download.sh > /dev/null
	touch .devcontainer/.installed

test:
	staticcheck ./...
	go test ./...

run: setup
	go run portfoli.go -verbose $(PATHARGS)

dist: setup
	(rm -rf dist || true) && mkdir dist
	go run portfoli.go -verbose -dist -dist.dir ${PWD}/dist $(PATHARGS)
	cp -r public dist/static
	mv dist/static/favicon.ico dist

docker:
	docker build . -t portfoligo:latest -f docker/Dockerfile

clean:
	rm -rf .devcontainer/.installed dist _build

.PHONY: test setup run build dist docker clean