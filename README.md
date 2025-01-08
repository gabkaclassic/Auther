# Auth Service

Этот проект представляет собой сервис для аутентификации, реализованный на Go, с использованием PostgreSQL в качестве базы данных.

## Конфигурация

Перед запуском убедитесь, что файл `configs/config.yaml` содержит правильные параметры конфигурации. Пример:
```yaml
server:
  port: 5000
  host: localhost
db:
  dialect: postgresql
  host: localhost
  port: 5432
  Username: USER
  password: PASSWORD
  database: auth_service
jwt:
  secret: secret
  expiration: 3600
  refresh_secret: refresh-secret
  refresh_expiration: 86400
admin:
  tokens:
    - admin_token_1
```

## Сборка и запуск

### Сборка Docker-образа

1. Соберите образ:
   ```bash
   docker build -t auth_service .
   ```

2. Запустите контейнер:
   ```bash
   docker run -d --name auth_service_app -p 5000:5000 --env CONFIG_PATH=/app/configs/config.yaml auth_service
   ```

### Использование Docker Compose

1. Убедитесь, что файл `docker-compose.yml` настроен корректно.
2. Запустите сервисы:
   ```bash
   docker-compose up --build
   ```

3. Для остановки контейнеров:
   ```bash
   docker-compose down
   ```

### Локальная разработка

1. Установите зависимости:
   ```bash
   go mod download
   ```

2. Соберите приложение:
   ```bash
   make build_app
   ```

3. Запустите тесты:
   ```bash
   make test
   ```

4. Очистите бинарные файлы:
   ```bash
   make clean
   ```

## Описание сервисов

- **app**: Основное приложение, запускаемое на Go. Порт по умолчанию — `5000`.
- **db**: PostgreSQL, порт по умолчанию — `5432`.

## Примечания

- Все переменные окружения, необходимые для работы приложения, задаются в файле `configs/config.yaml`.
- Если конфигурация изменяется, убедитесь, что изменения синхронизированы с монтируемым файлом в контейнере.

## Полезные команды

- Просмотр логов приложения:
  ```bash
  docker logs auth_service_app
  ```

- Перезапуск сервиса:
  ```bash
  docker-compose restart
  ```

- Удаление контейнеров и данных:
  ```bash
  docker-compose down -v
  ```