run:
	go run portfolio.go

dist:
	go build -o portfolio portfolio.go

test:
	staticcheck ./...
	go test ./...