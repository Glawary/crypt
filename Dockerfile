FROM golang:1.22.2

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8050

CMD ["./main"]