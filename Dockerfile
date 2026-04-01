FROM golang:1.25.5 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o ./app ./cmd/app/main.go

FROM alpine:latest
COPY --from=builder /app/app /app
ENTRYPOINT ["/app"]

