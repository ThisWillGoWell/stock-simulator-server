FROM golang:1.7.3
WORKDIR /go/src/github.com/ThisWillGoWell/stock-simulator-server
COPY main.go .
COPY src .
COPY vendor .
RUN GOOS=linux go build cgo -o app .

FROM alpine:3.10.3
COPY --from=0 /go/src/github.com/ThisWillGoWell/stock-simulator-server/app .
COPY config .
RUN chmod +x app
EXPOSE 8000
CMD ["./app"]