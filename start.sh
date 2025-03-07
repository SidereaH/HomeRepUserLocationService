#!/bin/bash

# Запуск TimescaleDB в фоновом режиме
docker-entrypoint.sh postgres &

# Ждём, пока TimescaleDB запустится
echo "Waiting for TimescaleDB to start..."
while ! pg_isready -U user -d location_service -h localhost -p 5432; do
  sleep 1
done

# Создаём таблицу, если её нет
psql -U user -d location_service -c "
CREATE TABLE IF NOT EXISTS user_locations (
    user_id BIGINT NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    time TIMESTAMPTZ DEFAULT NOW()
);

SELECT create_hypertable('user_locations', 'time');
"

# Запуск Go-приложения
echo "Starting Go application..."
./location-service