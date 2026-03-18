FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY . .
RUN go mod tidy

RUN go build -o /tmp/benthos ./cmd/benthos/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /tmp/benthos /benthos
ENTRYPOINT ["/benthos"]
