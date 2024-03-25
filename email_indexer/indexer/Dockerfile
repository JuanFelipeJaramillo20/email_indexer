FROM golang:1.22.1-bookworm

RUN mkdir app

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o indexer .

CMD ["/app/indexer"]