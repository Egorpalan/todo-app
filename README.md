# Todo Application

Это приложение для управления задачами, позволяющее пользователям добавлять, удалять и обновлять задачи. Приложение использует SQLite в качестве базы данных и написано на языке Go.

## Основные функции

- Добавление новых задач.
- Удаление задач.
- Обновление существующих задач.
- Получение списка задач.
- Поддержка повторяющихся задач.

## Технологии

- Go (Golang)
- SQLite
- Docker
- HTML/CSS для фронтенда

## Установка и запуск

### Клонирование репозитория

```bash
git clone https://github.com/Egorpalan/todo-app.git
```

## Запуск с помощью Docker

```
docker build -t todo-app .

docker run -p 7540:7540 --name todo-app-container todo-app
```