# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /gobox ./cmd

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /gobox /gobox

EXPOSE 8080

CMD ["/gobox"]
