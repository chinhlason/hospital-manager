FROM golang:1.22

WORKDIR /room

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV listen=:8881
ENV dbURL=postgresql://sonnvt:sonnvt@localhost:5432/demo
ENV userURL=authservice:8880

CMD go run ./cmd/main.go -httpAddr=${listen} -postgresURL=${dbURL} -userAddr=${userURL}

EXPOSE 8881
