FROM golang:1.23-bullseye AS builder
WORKDIR /src

# Copy shared library first
COPY shared/ ./shared/

# Copy service files
COPY ra/go.mod ra/go.sum ./ra/
WORKDIR /src/ra
RUN go mod download

WORKDIR /src
COPY ra/ ./ra/
WORKDIR /src/ra
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/ra ./cmd/ra

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/ra /usr/local/bin/ra
COPY ra/config/ /config/
EXPOSE 8080 9090
ENTRYPOINT ["/usr/local/bin/ra"]
