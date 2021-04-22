FROM golang:1.15.6 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o todo_service
EXPOSE 8080
ENTRYPOINT ["/app/todo_service"]
