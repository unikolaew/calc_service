# Распределённый вычислитель арифметических выражений

## Описание проекта
Этот проект реализует распределенную систему для выполнения арифметических вычислений с длительным временем обработки. Он состоит из:
- **Оркестратора** (сервер, управляющий вычислениями)
- **Агента** (вычислительный модуль, выполняющий отдельные операции)

### Основные принципы работы
- Оркестратор принимает выражение, разбивает его на задачи и управляет их выполнением.
- Агенты получают задачи, вычисляют результат и отправляют обратно на оркестратор.
- Пользователь периодически запрашивает результат вычисления.
- Операции выполняются "очень долго", поэтому система позволяет масштабировать вычислительные мощности.

## Запуск системы
### Требования
- **Go** >= 1.18

0. Перейти в директорию проекта:
   ```sh
   cd calc_service
   ```

1. Запустить оркестратор:
   ```sh
   go run orchestrator/cmd/main.go
   ```
2. Запустить агента:
   ```sh
   go run agent/cmd/main.go
   ```

## API эндпоинты
### 1. Добавление выражения
```sh
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{"expression": "2+2*2"}'
```
Ответ:
```json
{
    "id": "123e4567-e89b-12d3-a456-426614174000"
}
```

### 2. Получение списка выражений
```sh
curl --location 'localhost/api/v1/expressions'
```
Ответ:
```json
{
    "expressions": [
        {"id": "123e4567-e89b-12d3-a456-426614174000", "status": "pending", "result": null}
    ]
}
```

### 3. Получение выражения по ID
```sh
curl --location 'localhost/api/v1/expressions/123e4567-e89b-12d3-a456-426614174000'
```
Ответ:
```json
{
    "expression": {"id": "123e4567-e89b-12d3-a456-426614174000", "status": "done", "result": 6}
}
```

### 4. Получение задачи агентом
```sh
curl --location 'localhost/internal/task'
```
Ответ:
```json
{
    "task": {"id": 1, "arg1": 2, "arg2": 2, "operation": "multiplication", "operation_time": 5000}
}
```

### 5. Прием результата обработки данных.
```sh
curl --location 'localhost/internal/task' \
--header 'Content-Type: application/json' \
--data '{
  "id": 1,
  "result": 2.5
}'
```

## Переменные окружения
- `TIME_ADDITION_MS` - время сложения (мс)
- `TIME_SUBTRACTION_MS` - время вычитания (мс)
- `TIME_MULTIPLICATIONS_MS` - время умножения (мс)
- `TIME_DIVISIONS_MS` - время деления (мс)
- `COMPUTING_POWER` - количество горутин у агента

## Тестирование
Запуск тестов:
```sh
go test ./...
```
