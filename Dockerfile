FROM golang:1.22.5

WORKDIR /app

ARG GITHUB_TOKEN

RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"

ENV GOPRIVATE=github.com/timewise-team/timewise-models
ENV GONOSUMDB=github.com/timewise-team/timewise-models

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/

RUN go get -u github.com/swaggo/swag/cmd/swag
RUN swag init --parseDependency

RUN CGO_ENABLED=0 GOOS=linux go build -o timewise-api .

EXPOSE 8080

CMD ["./timewise-api"]