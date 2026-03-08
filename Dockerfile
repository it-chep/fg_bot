FROM golang:1.24.4-alpine AS builder

# Устанавливаем системные зависимости для компиляции
RUN apk add \
    ffmpeg-dev \
    make \
    gcc \
    musl-dev \
    pkgconfig

WORKDIR /app

COPY .env ./

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Компилируем приложение с оптимизациями
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/main \
    ./cmd/main.go

# Stage 2: Runtime
FROM alpine:3.21

# Устанавливаем только необходимые runtime зависимости
RUN apk add \
    ffmpeg \
    ca-certificates \
    tzdata

WORKDIR /app

# Копируем бинарник из builder stage
COPY --from=builder /app/main /app/main
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/k8s /app/k8s

EXPOSE 8080 7002

# Запускаем скомпилированный бинарник
CMD ["/app/main"]