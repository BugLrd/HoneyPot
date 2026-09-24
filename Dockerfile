# Dockerfile

FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go build -o /honeypot-dashboard cmd/server/main.go

EXPOSE 8080

CMD ["/honeypot-dashboard"]
