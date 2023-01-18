FROM golang:latest

RUN apt-get update && \
    apt-get upgrade -y

RUN useradd -ms /bin/bash port
USER port