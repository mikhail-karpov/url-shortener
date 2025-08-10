FROM golang:1.24-alpine AS builder

WORKDIR /opt/app
COPY go.mod go.sum ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/main.go

FROM scratch
WORKDIR /opt/app
COPY --from=builder /opt/app/server server

EXPOSE 8080
ENTRYPOINT ["./server"]