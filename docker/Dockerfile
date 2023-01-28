FROM golang:1.19

WORKDIR portfoli.go
COPY . 
RUN make build

FROM alpine:latest

COPY --from=0 portfoli.go/_build/portfoli.go \
    public/ \
    templates \