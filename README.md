# ReBalancer (Backend)

Этот репозиторий содержит backend-сервис для расчета ребалансировки портфеля активов, торгуемых на Московской бирже. В основе — использование API Т-Инвестиций через [официальный Go SDK](https://github.com/ssummers02/invest-api-go-sdk).

Проект предоставляет удобный инструмент для автоматизации расчетов, позволяя инвесторам поддерживать заданные пропорции активов или следовать за индексами, исключая необходимость ручного пересчета.

## ⚙️ Локальная установка и запуск

1. **Настройка конфигурации**  
   Укажите данные в файле `config/AccountData.go`.

2. **Запуск проекта**  
  ***Локальный запуск***
   Для запуска сервера локально выполните:
   ```powershell
   go run cmd/main.go
   ```
   Обратите внимание, что для успешного запуска проекта локально вы должны иметь запущенную базу данных PostgreSQL.
   ***Запуск через Docker***
  ```powershell
  docker build -t rebalancer_back .
  ```

  ```powershell
  docker run --name postgres -e POSTGRES_PASSWORD=your_db_password -p 5433:5432 -d postgres
  ```

  ```powershell
  docker run --name rebalancer_back -p 8080:8080 rebalancer_back
  ```

   
## 🧩 Технические детали

- **Язык программирования**: Go
- **API-интеграция**: Используется [SDK Т-Инвестиций](https://github.com/ssummers02/invest-api-go-sdk) для получения данных с Московской биржи.
- **База данных**: Встроенное хранилище пользователей

### 🔗 Доступные API-эндпоинты

- `POST /calc` — расчет ребалансировки портфеля  
- `POST /signup` — регистрация нового пользователя  
- `POST /finduser` — проверка существования пользователя  
- `POST /deleteuser` — удаление пользователя из базы данных

## 📌 Примеры запросов

**POST /calc**
```json
{
  "cur_allocation": [
    { "ticker": "FOR_TEST_1", "number": 10 },
    { "ticker": "FOR_TEST_0", "number": 5000 },
    { "ticker": "RUB", "number": 10000 }
  ],
  "target_allocation": [
    { "ticker": "FOR_TEST_1", "percent": 50 },
    { "ticker": "FOR_TEST_0", "percent": 50 }
  ]
}
```

**POST /signup**
```json
{
  "email": "kkhuzzyatov@gmail.com",
  "password": "1234"
}
```

**POST /finduser**
```json
{
  "email": "kkhuzzyatov@gmail.com",
  "password": "1234"
}
```

**POST /deleteuser**
```json
{
  "email": "kkhuzzyatov@gmail.com",
  "password": "1234"
}
```