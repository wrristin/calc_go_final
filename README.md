# Распределенный вычислитель арифметических выражений. Финал
Проект представляет из себя распределённый вычислитель арифметических выражений. Для регистрации пользователь отправляет запрос POST /api/v1/register { "login": , "password": }
В ответ получает 200+OK (в случае успеха)
В противном случае - ошибка.
Для входа пользователь отправляет запрос POST /api/v1/login { "login": , "password": }
В ответ получает 200+OK и JWT токен для последующей авторизации.
## Запуск проекта:
- Для начала вам нужно клонировать репозиторий.
```
git clone https://github.com/wrristin/calc_go_final.git
cd calc_service
```
- Далее вам нужно установить зависимости.
Windows: Скачайте protoc (выберите protoc-*.zip), распакуйте и добавьте bin/ в PATH.
- Плагины для Go:
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```
- После этого вам потребуется установить Установить компилятор C (требуется для SQLite)
Windows: Установите TDM-GCC x86_64 (выберите 64-битную версию).
Примечание: Обязательно добавьте путь к 64-битному GCC в переменную окружения PATH:
Например: C:\Program Files\mingw64\bin или C:\TDM-GCC-64\bin.
- После этого проверьте версию GCC
```
gcc --version
```
- Запустите оркестратор
Включите CGO (обязательно для Windows)
```
$env:CGO_ENABLED = "1"  (Вписывайте данную команду в PowerShell)
go run cmd/orchestrator/main.go
```
- Запустите агент
```
$env:CGO_ENABLED = "1" 
go run cmd/agent/main.go
```
- Регистрация пользвателя:
```
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"login": "user", "password": "pass"}'
  ```
  - Авторизация (получение токена):
```
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"login": "user", "password": "pass"}'
  ```
  - Добавить выражение:
```
  curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Authorization: Bearer ВАШ_ТОКЕН" \
  -H "Content-Type: application/json" \
  -d '{"expression": "2+2*2"}'
  ```
  - Проверить статус выражения:
 ```
curl http://localhost:8080/api/v1/expressions \
  -H "Authorization: Bearer ВАШ_ТОКЕН"
  ```
