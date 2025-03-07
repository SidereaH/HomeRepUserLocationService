# Этап 1: Сборка Go-приложения
FROM golang:1.23 AS builder

WORKDIR /app

# Копируем файлы go.mod и go.sum перед тем, как загружать зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o location-service ./server/server.go

# Этап 2: Создаём минимальный контейнер для запуска Go-сервера
FROM gcr.io/distroless/base-debian11

WORKDIR /app

# Копируем бинарник из builder-а
COPY --from=builder /app/location-service .

# Открываем порт gRPC
EXPOSE 50051

# Запускаем Go-приложение
CMD ["/app/location-service"]
