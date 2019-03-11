FROM golang:1.13-alpine3.10 as builder

WORKDIR /usr/src/app

RUN apk add --no-cache ca-certificates curl gcc git libc-dev \
  && curl -o /go/bin/grpc_health_probe -sSL https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.2.1/grpc_health_probe-linux-amd64 \
  && chmod +x /go/bin/grpc_health_probe \
  && go get -u golang.org/x/lint/golint

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build and install a static binary, stripping DWARF debugging information and
# preventing the generation of the Go symbol table.
RUN GOOS=linux GOARCH=amd64 go install -a -ldflags '-w -s -linkmode external -extldflags "-static"' ./cmd/comment

FROM scratch

COPY --from=builder /go/bin/comment /comment
COPY --from=builder /go/bin/grpc_health_probe /grpc_health_probe
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/comment"]
