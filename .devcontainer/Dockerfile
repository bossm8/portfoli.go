FROM golang:latest

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
        git \
        vim \
        webp \
        bash-completion

RUN useradd -ms /bin/bash port
USER port