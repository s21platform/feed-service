#!/bin/bash

# Применим миграции
# ВАЖНО! Для локального запуска вызвать из корня проекта, предварительно подгрузив переменные окружения в терминал командой  `set -a; source <путь к .env файлу>; set +a`
# Установка goose: go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir ./migrations postgres "user=$FEED_SERVICE_POSTGRES_USER password=$FEED_SERVICE__POSTGRES_PASSWORD dbname=$FEED_SERVICE__POSTGRES_DB host=$FEED_SERVICE__POSTGRES_HOST port=$FEED_SERVICE__POSTGRES_PORT sslmode=disable" down