FROM golang:1.22.5

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/

RUN CGO_ENABLED=0 GOOS=linux go build -o timewise-api .

EXPOSE 8080

CMD ["./timewise-api"]