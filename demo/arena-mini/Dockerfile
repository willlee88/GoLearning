# syntax=docker/dockerfile:1
# Build context: monorepo root
FROM golang:1.22-bookworm AS build

WORKDIR /src
COPY demo/arena-mini/server/ ./
# GOARCH defaults to the build machine (amd64 / arm64)
RUN go mod tidy \
 && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/arena ./cmd/server

COPY demo/arena-mini/web /out/web

# Non-root static binary (no libc / no shell)
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

COPY --from=build /out/arena /arena
COPY --from=build /out/web /web

ENV ADDR=:8080
ENV WEB_DIR=/web

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/arena"]
