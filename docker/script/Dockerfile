FROM golang:1.22-bookworm AS builder

ENV GETH_VERSION=0.9.2
ENV STORY_VERSION=0.10.0

COPY script/init /source/init

RUN apt update \
 && apt install -y unzip

RUN cd /source/init \
 && go build -ldflags "-s -w" -o story-init main.go

RUN mkdir /source/story \
 && cd /source/story \
 && wget https://github.com/piplabs/story/archive/refs/tags/v$STORY_VERSION.zip -O source.zip \
 && unzip source.zip \
 && cd story-$STORY_VERSION \
 && go build -ldflags "-s -w" -o story ./client \
 && mv story /source/story/story

RUN mkdir /source/geth \
 && cd /source/geth \
 && wget https://github.com/piplabs/story-geth/archive/refs/tags/v$GETH_VERSION.zip -O source.zip \
 && unzip source.zip \
 && cd story-geth-$GETH_VERSION \
 && go build -ldflags "-s -w" -o geth ./cmd/geth \
 && mv geth /source/geth/geth

FROM debian:bookworm

COPY config /init-config
COPY --from=builder /source/init/story-init /source/geth/geth /source/story/story /usr/local/bin/

RUN apt update \
 && apt install -y openssl dnsutils

COPY script/init-node.sh /
RUN chmod +x /init-node.sh

ENTRYPOINT ["bash", "/init-node.sh"]
