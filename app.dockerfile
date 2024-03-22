# 開発用
FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./
COPY .air.toml ./

RUN go install github.com/cosmtrek/air@latest

CMD ["air"]

# 本番用
# FROM golang:1.19

# WORKDIR /app

# COPY go.mod go.sum ./

# RUN go install github.com/cosmtrek/air@latest

# CMD ["go", "run", "main.go"]