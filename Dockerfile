FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git && adduser -D -g '' gopher && apk add -U --no-cache ca-certificates
WORKDIR /app/
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o tweetsvg main.go

FROM scratch

ENV ACCESS_TOKEN=
ENV ACCESS_TOKEN_SECRET=
ENV CONSUMER_KEY=
ENV CONSUMER_SECRET=

WORKDIR /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/tweetsvg /app/tweetsvg
USER gopher

ENTRYPOINT ["/app/tweetsvg"]
