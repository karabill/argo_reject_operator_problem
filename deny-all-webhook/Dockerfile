FROM golang:1.15-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download -x

COPY main.go main.go

RUN CGO_ENABLED=0 go build -v -ldflags="-s -w" -o /server main.go

ENTRYPOINT [ "/server" ]
