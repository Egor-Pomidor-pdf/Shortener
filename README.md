# URL Shortener

Сервис сокращения ссылок с аналитикой переходов.

## Стек

- **Go** — backend
- **PostgreSQL** — хранение ссылок
- **ClickHouse** — аналитика кликов

## Запуск

```bash
docker compose -f docker/docker-compose.yml up -d -- build


## API

### Создать короткую ссылку
```bash
POST /shorten
Content-Type: application/json

{
  "original_url": "https://example.com/very/long/url",
  "client_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Перейти по короткой ссылке
```bash
GET /s/{short_code}
```

### Получить аналитику
```bash
GET /analytics/{short_code}
```

## Конфигурация

Настройки в `config/.env`