FROM golang:1.23.3

WORKDIR /app

COPY . . 

RUN cd ./cmd/api && go build -o main .

CMD ["./cmd/api/main"]