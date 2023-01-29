# ---------------------------------build-------------------------------------- #
FROM golang:1.19

ENV CGO_ENABLED=0

WORKDIR /home/portfoli.go

RUN apt-get update && \
    apt-get install -y wget unzip 

COPY . .

RUN /bin/bash .devcontainer/download.sh

RUN mkdir -p www/ && \
    mv public www/public && \
    mv templates www/templates && \
    mv examples/configs www/configs

RUN go build -o portfoli-go portfoli.go && \
    chmod +x portfoli-go && ./portfoli-go --help && \
    chmod +x docker/static.sh

# ---------------------------------package------------------------------------ #
FROM alpine:latest

RUN mkdir -p /var/www/portfoli.go && \
    adduser --system \
            --home /var/www/portfoli.go \
            --shell /bin/sh \
            --ingroup www-data \
            portfoligo

COPY --from=0 /home/portfoli.go/portfoli-go /usr/sbin/portfoli-go
COPY --from=0 /home/portfoli.go/www /var/www/portfoli.go
COPY --from=0 /home/portfoli.go/docker/static.sh /usr/sbin/portfoli-go-static.sh

RUN chown -R portfoligo:www-data /var/www/portfoli.go

USER portfoligo

WORKDIR /var/www/portfoli.go

EXPOSE "8080"

ENTRYPOINT ["portfoli-go"]

CMD ["-srv.address", "0.0.0.0", "-config.dir", "/var/www/portfoli.go/configs"]