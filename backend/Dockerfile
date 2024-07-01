FROM golang:1.22.2-alpine AS builder
WORKDIR /app

RUN apk add --no-cache \
	gcc \
	musl-dev \
	ca-certificates

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server  -ldflags '-s -w -extldflags "-static"' ./cmd/web

FROM scratch
WORKDIR /app

LABEL maintainer="Said Rodrigues"
LABEL email="coderflemis@gmail.com"

COPY --from=builder /app/server ./
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/config.json ./

EXPOSE 42069

CMD ["/app/server"]
