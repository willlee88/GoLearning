# syntax=docker/dockerfile:1
FROM golang:1.22-bookworm AS build
WORKDIR /src
COPY services/runner/go.mod ./
COPY services/runner/*.go ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/runner .

# Need full Go toolchain at runtime to execute user programs
FROM golang:1.22-bookworm
RUN useradd -m -u 10001 runner \
 && mkdir -p /tmp/golearning-runs \
 && chown -R runner:runner /tmp/golearning-runs \
 && mkdir -p /home/runner/go/cache \
 && chown -R runner:runner /home/runner \
 && apt-get update \
 && apt-get install -y --no-install-recommends wget \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /home/runner
COPY --from=build /out/runner /usr/local/bin/runner

ENV ADDR=:8091
ENV HOME=/home/runner
ENV GOPATH=/home/runner/go
ENV GOCACHE=/home/runner/go/cache
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPROXY=off
ENV GOSUMDB=off
ENV GOTOOLCHAIN=local

# Warm stdlib compile once as root-owned cache then chown
USER root
RUN mkdir -p /tmp/warm \
 && printf '%s\n' 'package main' 'import "fmt"' 'func main(){fmt.Println("ok")}' > /tmp/warm/main.go \
 && printf '%s\n' 'module warm' '' 'go 1.22' > /tmp/warm/go.mod \
 && cd /tmp/warm && CGO_ENABLED=0 GOPROXY=off GOSUMDB=off go run . \
 && rm -rf /tmp/warm \
 && chown -R runner:runner /home/runner

USER runner
EXPOSE 8091

HEALTHCHECK --interval=30s --timeout=3s --start-period=20s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8091/healthz >/dev/null || exit 1

ENTRYPOINT ["/usr/local/bin/runner"]

