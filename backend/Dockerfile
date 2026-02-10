FROM golang:1.22-bookworm AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/server ./cmd/server

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
# Install mdbook (adjust version if needed)
ARG MDBOOK_VERSION=0.4.40
RUN curl -fsSL -o /tmp/mdbook.tar.gz https://github.com/rust-lang/mdBook/releases/download/v${MDBOOK_VERSION}/mdbook-v${MDBOOK_VERSION}-x86_64-unknown-linux-gnu.tar.gz \
    && tar -xzf /tmp/mdbook.tar.gz -C /usr/local/bin \
    && rm /tmp/mdbook.tar.gz
WORKDIR /app
COPY --from=builder /bin/server /app/server
COPY books /data/books
ENV BOOKS_ROOT=/data/books
ENV BOOKS_BUILD_ROOT=/data/build
EXPOSE 8080
CMD ["/app/server"]
