FROM golang:1.24-alpine

WORKDIR /app
COPY . .
RUN go build -o /usr/local/bin/todo-syncer

ENTRYPOINT ["todo-syncer"]
