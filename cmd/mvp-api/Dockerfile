FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/mvp-api ./cmd/mvp-api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/mvp-api .

ENV APP_PORT=4004 \
   APP_HOST=0.0.0.0

EXPOSE 4004
CMD ["./mvp-api"]


