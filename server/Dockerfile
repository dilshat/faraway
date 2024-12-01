FROM golang:1.21-alpine as builder

WORKDIR /app

COPY server/ /app/server/
COPY pkg/ /app/pkg/

RUN go mod init dilshat/faraway
RUN go mod tidy

RUN GOARCH=amd64 GOOS=linux go build -o server ./server
RUN chmod +x ./server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
