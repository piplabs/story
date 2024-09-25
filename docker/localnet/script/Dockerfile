FROM ubuntu:22.04

RUN apt update \
 && apt install -y openssl dnsutils

COPY script/init-node.sh /
RUN chmod +x /init-node.sh

COPY binary/* /usr/local/bin

ENTRYPOINT ["bash", "/init-node.sh"]
