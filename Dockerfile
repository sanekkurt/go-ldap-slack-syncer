FROM golang:1.17-alpine as builder

WORKDIR /src/app

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
        go build -ldflags="-w -s" -mod vendor -o /app ./cmd/main/...

# ---

FROM alpine:3.15

WORKDIR /

RUN apk add --no-cache ca-certificates

COPY --from=builder /app /app

ENTRYPOINT ["/app"]