FROM golang:1.24 AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN go mod tidy

WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o tgnotify

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/tgnotify .
ENTRYPOINT [ "./tgnotify" ]
