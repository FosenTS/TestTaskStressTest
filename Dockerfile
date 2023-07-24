FROM golang:alpine

WORKDIR /build

COPY .env .
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o main main.go

CMD ["./main"]