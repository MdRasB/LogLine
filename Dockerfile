# ---------- Build stage ----------
FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /logline \
    ./cmd/api


# ---------- Runtime stage ----------
FROM alpine:3.24

WORKDIR /app

RUN apk add --no-cache ca-certificates

RUN addgroup -S logline && \
    adduser -S -G logline logline

COPY --from=builder /logline /app/logline
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/web /app/web

RUN chown -R logline:logline /app

USER logline

EXPOSE 8080

CMD ["/app/logline"]
