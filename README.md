# Распределённый вычислитель арифметических выражений
Этот проект представляет собой веб-сервис для вычисления арифметических выражений. Пользователь отправляет выражение через HTTP-запрос, а сервис возвращает результат. 
## Back-end часть
### Состоит из 2-ух элементов:
- Сервер, который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения. (Оркестратор)
- Вычислитель, который может получить от оркестратора задачу, выполнить его и вернуть серверу результат. (Агент)
## Как пользоваться проектом
1. Установите переменные среды
Для линукса/macOs:
```bash
   export TIME_ADDITION_MS=1000
   export TIME_SUBTRACTION_MS=1000
   export TIME_MULTIPLICATIONS_MS=2000
   export TIME_DIVISIONS_MS=2000
   export COMPUTING_POWER=4
   
```
Для Windows в командной строке cmd:
```bash
   set TIME_ADDITION_MS=1000
   set TIME_SUBTRACTION_MS=1000
   set TIME_MULTIPLICATIONS_MS=2000
   set TIME_DIVISIONS_MS=2000
   set COMPUTING_POWER=4
```
2. Установите зависимости:
```
go mod download
```
3. Запустите агент и оркестратор в разных терминалах
Запуск оркестратора:
```
go run ./cmd/orchestrator/main.go
```
Запуск агента:
```
go run ./cmd/agent/main.go
```
## Примеры использования
### Успешный запрос
```
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```
### Ошибка: неверное выражение
```
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+*2"
}'
```
### Получение списка выражений
```
curl --location 'localhost:8080/api/v1/expressions'
```
### Получение выражения по ID
```
curl --location 'localhost:8080/api/v1/expressions/id'
```
