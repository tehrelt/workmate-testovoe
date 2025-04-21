# Test task

## Stack
- Go 1.24
- RabbitMQ (amqp091)
- PostgreSQL (pgxpool, migrate)
- Jaeger (opentelemetry)
- Prometheus
- Grafana

## Deploy

В корневой папке выполните команду:

```
docker compose up -d --build
```

При успешном запуске приложения API будет доступен на порту 4000
JaegerUI - http://localhost:16686
Prometheus - http://localhost:9090
Grafana - http://localhost:3000
RabbitMQ UI - http://localhost:15672
