# Go Auth Training

![Go](https://img.shields.io/badge/Go-1.20-blue) ![Gin](https://img.shields.io/badge/Gin-framework-green) ![JWT](https://img.shields.io/badge/JWT-auth-orange)

Простой сервис аутентификации и регистрации пользователей на Go с использованием **чистой архитектуры** и **JWT**.

---

## Описание

Проект демонстрирует:

- Чистую архитектуру (Clean Architecture) с четким разделением слоев (Handler, UseCase, Repository, Entity).
- Регистрацию и вход пользователей с хешированием паролей.
- Генерацию и валидацию JWT токенов (access + refresh tokens).
- Middleware для проверки и обновления токенов.
- Настроенный CORS для фронтенда.
- Использование Gin в качестве HTTP-фреймворка.

---

## Структура проекта

```plaintext
/internal
  /adapter      # Адаптеры (например, JWT, DB)
/entity       # Сущности домена (User и др.)
/usecases     # Бизнес-логика (UseCase слои)
/http         # HTTP обработчики и middleware
/pkg          # Утилиты (хеширование паролей, экстракция токенов и др.)
```
Установка и запуск

    Клонируй репозиторий:

```plaintext
git clone https://github.com/yourusername/golang-auth-training.git
cd golang-auth-training
```
    Установи зависимости:
```plaintext
go mod tidy
```
    Создай файл .env с настройками JWT и базы данных:
```plaintext
JWT_SECRET=your_jwt_secret_key
DB_DSN=your_database_dsn
```
    Запусти сервер:
```plaintext
go run cmd/main.go
```

API
Регистрация

    POST /register

    Тело запроса:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```
Ответ:
```json
{
  "id": 1,
  "email": "user@example.com"
}
```
Вход (Login)

    POST /login

    Тело запроса:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```
Ответ:
```json
{
  "access_token": "<jwt_access_token>",
  "refresh_token": "<jwt_refresh_token>"
}
```

Middleware

    JWT middleware проверяет access_token из заголовка Authorization.

    Если access_token истек, пытается обновить его с помощью refresh_token из заголовка X-Refresh-Token.

    При успешном обновлении возвращает новые токены в заголовках X-New-Access-Token и X-New-Refresh-Token.

CORS

Настроен с помощью github.com/gin-contrib/cors:

    Разрешены методы: GET, POST, PUT, DELETE.

    Разрешены заголовки: Origin, Authorization, Content-Type, X-Refresh-Token.

    Разрешены креденшелы (cookies и т.д.).

    Экспонируются заголовки новых токенов.

Технологии

    Go

    Gin

    JWT

    bcrypt (для хеширования паролей)

    Viper (опционально, для конфигураций)

    (Добавь сюда свою БД или ORM, если есть)