FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o stress-test-cli .

FROM scratch

COPY --from=builder /app/stress-test-cli /usr/bin/stress-test-cli
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["stress-test-cli"]
