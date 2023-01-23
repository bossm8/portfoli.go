.devcontainer/.installed:
	/bin/bash .devcontainer/download.sh
	touch .devcontainer/.installed

setup: .devcontainer/.installed

dist: setup
	test -d _build || mkdir _build
	go build -o _build/portfolio portfolio.go

test:
	staticcheck ./...
	go test ./...

run: setup
	go run portfolio.go

.PHONY: test dist setup run