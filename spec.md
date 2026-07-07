# Спецификация веб-приложения "Мини анкета". 

Бекенд — go 1.26, фронтенд — статический HTML + CSS + vanilla JS

## Backend
Написан на go 1.26. В качестве сервера берется нативная библиотека net/http.

Спецификация API описана с использованием OpenAPI (генерация через oapi-codegen). Лежит в папке `api/`.

### API методы

- `GET /questions` — возвращает список вопросов анкеты. Код ответа: `200`
- `POST /answers` — принимает ответы пользователя и сохраняет их в памяти. Код ответа: `201` / `400`

### Вопросы анкеты (хардкод)

Всего 5 вопросов разных типов:

| № | Тип | Текст | Варианты |
|---|-----|-------|----------|
| 1 | text | Как вас зовут? | — |
| 2 | single_choice | Какой ваш любимый цвет? | Красный, Синий, Зелёный, Жёлтый, Другой |
| 3 | multiple_choice | Какими языками программирования вы владеете? | Go, Python, JavaScript, Java, C++, Rust, Другой |
| 4 | single_choice | Сколько лет вы занимаетесь программированием? | Меньше года, 1–3 года, 3–5 лет, 5–10 лет, Больше 10 лет |
| 5 | text | Что бы вы хотели улучшить в нашем продукте? | — |

### Хранение данных

Ответы сохраняются в памяти (in-memory) в slice `[]api.AnswerRecord`, защищённый `sync.Mutex`. Данные теряются при перезапуске сервера.

### Тестирование

Оба метода покрыты юнит-тестами (`internal/handler/handler_test.go`, 8 тестов). Запуск: `go test -v ./...`.

## Ручное тестирование API бекенда

```bash
# 1. Получить список вопросов
curl -v http://localhost:8080/questions

# 2. Отправить ответы (текстовый ответ)
curl -v -X POST http://localhost:8080/answers \
  -H "Content-Type: application/json" \
  -d '{
    "respondent": "Иван",
    "answers": [
      {"question_id": "550e8400-e29b-41d4-a716-446655440000", "value": "Синий"}
    ]
  }'

# 3. Отправить ответы (множественный выбор)
curl -v -X POST http://localhost:8080/answers \
  -H "Content-Type: application/json" \
  -d '{
    "answers": [
      {"question_id": "550e8400-e29b-41d4-a716-446655440000", "value": ["Красный", "Зелёный"]}
    ]
  }'

# 4. Отправить ответы (без респондента, несколько вопросов)
curl -v -X POST http://localhost:8080/answers \
  -H "Content-Type: application/json" \
  -d '{
    "answers": [
      {"question_id": "550e8400-e29b-41d4-a716-446655440000", "value": "Текстовый ответ"},
      {"question_id": "550e8400-e29b-41d4-a716-446655440001", "value": "Вариант А"}
    ]
  }'
```

## Frontend

Индексная страница (`/`) — статический HTML, который загружает вопросы через API и отправляет ответы.

- `/static/css/style.css` — стили
- `/static/js/app.js` — скрипты

Статика вшита в бинарник через `//go:embed` (пакет `embed`), читается из памяти при запуске.

## Docker

Минималистичный образ на основе `scratch` (двухэтапная сборка).

```dockerfile
FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/server ./cmd/server

FROM scratch
COPY --from=build /app/server /server
EXPOSE 8080
CMD ["/server"]
```

Сборка и запуск:

```bash
docker build -t mini-questionnaire .
docker run -p 8080:8080 mini-questionnaire
```

Статика вшита в бинарник через `//go:embed`, дополнительные слои не нужны.

## Отладка

Конфигурация для VS Code — `.vscode/launch.json` (тип `go`, режим `debug`, точка входа `cmd/server`). Брейкпоинты расставляются в редакторе на строках `internal/handler/handler.go:57` (начало `GetQuestions`) и `:87` (`s.answers = append`).