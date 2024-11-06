# Support setting various labels on the final image
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

# Build Geth in a stock Go builder container
FROM golang:1.22-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /story

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies are cached if the go.mod and go.sum files are not changed
RUN go mod download

ADD . /story/
RUN go build -o story ./client

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /story/story /usr/local/bin/

EXPOSE 8545 8546 30303 30303/udp

WORKDIR /root/.story/story
ENTRYPOINT ["/bin/sh", "-c", "story init --network $NETWORK && exec story run \"$@\"", "--"]

# Add some metadata labels to help programmatic image consumption
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

LABEL commit="$COMMIT" version="$VERSION" buildnum="$BUILDNUM" network="$NETWORK"
