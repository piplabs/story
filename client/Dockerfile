FROM scratch

# Install ca-certificates (for https to EVMs)
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Create /tmp directory (default cometBFT --temp-dir)
COPY --from=alpine:latest /tmp /tmp

COPY iliad /app

# Mount home directory at /iliad
VOLUME ["/iliad"]

ENTRYPOINT ["/app"]
CMD ["run"]
