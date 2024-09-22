FROM golang:1.22-bookworm AS builder

ENV GETH_VERSION=0.9.2

RUN apt update \
 && apt install -y unzip

RUN mkdir -p /source/geth \
 && cd /source/geth \
 && wget https://github.com/piplabs/story-geth/archive/refs/tags/v$GETH_VERSION.zip -O source.zip \
 && unzip source.zip \
 && cd story-geth-$GETH_VERSION \
 && go build -ldflags "-s -w" -o geth ./cmd/geth \
 && mv geth /source/geth/geth

RUN mkdir -p /bin/cosmovisor \
 && cd /bin/cosmovisor \
 && wget https://story-geth-binaries.s3.us-west-1.amazonaws.com/cosmovisor/cosmovisor -O cosmovisor \
 && chmod +x cosmovisor

FROM debian:bookworm

RUN apt update \
 && apt install -y openssl dnsutils wget

COPY script/init-node.sh /
RUN chmod +x /init-node.sh

COPY --from=builder /source/geth/geth /bin/cosmovisor/cosmovisor /usr/local/bin/

ENTRYPOINT ["bash", "/init-node.sh"]
