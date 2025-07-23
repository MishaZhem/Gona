FROM golang:1.24.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o supportchat ./cmd/main.go

CMD ["./supportchat"]