FROM golang:1.25-alpine AS builder

LABEL maintainer="Raffael Ramalho"
LABEL description="API de Empacotamento de Pedidos"

RUN apk add --no-cache git

WORKDIR /app

# dependências primeiro 
COPY go.mod go.sum ./
RUN go mod download

# código
COPY . .

#  Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

# Build o
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o empacotador-api .

# piada engraçada docker compose = docker com mc pose 
#  Runtime 
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Usuário não-root 
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

#  apenas binário e docs
COPY --from=builder /app/empacotador-api .
COPY --from=builder /app/docs ./docs

RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./empacotador-api"]