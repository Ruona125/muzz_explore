FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY proto ./proto
COPY internal ./internal
COPY cmd ./cmd

RUN go mod download

RUN CGO_ENABLED=0 go build -o explore-service ./cmd/explore-service

FROM alpine:3.18

RUN adduser -D explore
WORKDIR /home/explore

COPY --from=builder /app/explore-service .
COPY .env .env

USER explore

EXPOSE 50051

ENTRYPOINT ["./explore-service"]