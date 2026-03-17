FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app ./cmd/app

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app /app
COPY internal/config/config.yaml /config.yaml
EXPOSE 50051
ENTRYPOINT ["/app"]
