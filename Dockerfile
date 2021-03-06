FROM golang:1.15.6 as builder
WORKDIR /app
COPY . .
EXPOSE 8080
ENTRYPOINT ["go", "run", "main.go"]
