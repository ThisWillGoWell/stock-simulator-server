FROM golang:1.13.3-alpine3.10
WORKDIR /go/src/github.com/ThisWillGoWell/stock-simulator-server
COPY main.go main.go
COPY src src
COPY vendor vendor
RUN GOOS=linux GOARCH=amd64  go build -o app .
WORKDIR /bin
RUN cp /go/src/github.com/ThisWillGoWell/stock-simulator-server/app .
COPY config config
RUN chmod +x app
RUN ls -la
EXPOSE 8000
CMD ["app"]