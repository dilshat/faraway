FROM golang:1.21-alpine as builder

WORKDIR /app

COPY client/ /app/client/
COPY pkg/ /app/pkg/

RUN go mod init dilshat/faraway
RUN go mod tidy

RUN GOARCH=amd64 GOOS=linux go build -o client ./client
RUN chmod +x ./client

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/client .

CMD ["./client"]
