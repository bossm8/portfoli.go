FROM golang:latest

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
        git \
        vim \
        bash-completion

RUN useradd -ms /bin/bash port
USER port