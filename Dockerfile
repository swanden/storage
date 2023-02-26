FROM golang:1.20.1-alpine3.17
RUN apk add git
COPY . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/storage
CMD ["/bin/app"]