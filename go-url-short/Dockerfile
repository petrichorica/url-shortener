FROM golang:1.21.0

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

EXPOSE 3000

CMD ["./main"]