run:
	go run portfolio.go

dist:
	test -d _build || mkdir _build
	go build -o _build/portfolio portfolio.go

test:
	staticcheck ./...
	go test ./...