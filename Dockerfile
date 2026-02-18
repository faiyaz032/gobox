
FROM golang:1.25-alpine AS development


RUN go install github.com/air-verse/air@latest


WORKDIR /app


COPY go.mod go.sum ./


RUN go mod download


COPY . .


EXPOSE 8010


CMD ["air", "-c", ".air.toml"]


FROM golang:1.25-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./


RUN go mod download


COPY . .

# build app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
FROM alpine:latest AS production
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8010
CMD ["./main"]
