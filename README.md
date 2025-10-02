# Инкогнито (Секретный Гость)

Приложение позволяет клиентам Островка просматривать и участвовать в розыгрышах на проведение экспертизы отеля. Результат экспертизы - отчет с фотографиями, который участник заполняет в личном кабинете. За качественные отчеты клиент получает рейтинг и промокоды на бронирование в островке, за некачественные - теряет рейтинг. Рейтинг влияет вероятность выигрыша его заявки в следующих отборах.

- ссылка на видео-скринкаст: 

## Быстрый старт

```bash
git clone https://github.com/ostrovok-hackathon-2025/afrikanskie-petushki.git
cd afrikanskie-petushki
cp .env.example .env
docker compose up --build
открыть http://localhost:8080
```

## Зависимости и переменные окружения

**Рекомендованные ресурсы:**
- 4 CPU cores
- 16GB RAM

**Зависимости:**
- PostgreSQL (в `docker-compose.yml` как `db`)
- MinIO (в `docker-compose.yml` как `minio`)

**Переменные окружения**
- `CONFIG_PATH` - путь к файлу конфигурации для бэкенд-приложения
- `MINIO_ROOT_USER` - имя пользователя в S3
- `MINIO_ROOT_PASSWORD` - пароль в S3
- `POSTGRES_USER` - имя пользователя в PostgreSQL
- `POSTGRES_PASSWORD` - пароль в PostgreSQL
- `POSTGRES_DB` - имя базы данных PostgreSQL
- `NEXTAUTH_SECRET` - секретный ключ для NextAuth
- `NEXTAUTH_URL` - эндпоинт API NextAuth

## Сидирование

```bash
make create_test_data
# или
docker exec -i postgres_secret_guest psql -U admin -d secret-guest -f /var/lib/testdata/fill-all.sql
```

**Обратите внимание:** в тестовых данных создается розыгрыш лота на Moscow Grand Hotel (это будет единственный розыгрыш для Москвы). Его итоги будут объявлены через 5 минут после генерации. Успейте подать заявку.  

## Маршруты/доступ

- `/` — UI
- `/api/health` — 200 OK, JSON `{ "status": "ok" }`

**Тестовые пользователи:**

Клиент островка:
- Логин: chicherin
- Пароль: Chicherin12345!

Админ:
- Логин: root
- Пароль: Root12345!