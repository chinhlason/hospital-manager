FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV listen=:8880
ENV dbURL=scylladb:9042
ENV redisHost=localhost:6379

CMD go run ./cmd/main.go -httpAddr=${listen} -scyllaHost=${dbURL} -redisHost=${redisHost}

EXPOSE 8880
