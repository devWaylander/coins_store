
# Coins Store

## Доступные команды

Для получения списка всех доступных команд выполните команду в корне проекта:

```bash
make help
```

## Управление миграциями БД

Для осуществления ручного управления миграциями БД выполните следующие шаги:

1. Выполните команду для установки зависимостей:

   ```bash
   make installDeps
   ```

2. Ознакомьтесь с инструментом **dbmate** для работы с миграциями.

## Тестирование 
Для проверки покрытия unit тестами слоя **use-case** выполните команду:

```bash
go test -cover ./internal/service
```

Для запуска тестов необходимо запустить контейнер с тестовой БД и миграциями

```bash
docker compose up postgresdbtest migrationstest
```

Для запуска unit тестов auth выполните команду:

```bash
go test -timeout 30s -run ^Test_middleware_LoginWithPass$ github.com/devWaylander/coins_store/internal/middleware/auth -count=1 -v
```

Для запуска unit тестов usecase выполните команду:

```bash
go test -timeout 30s -run ^Test_middleware_LoginWithPass$ github.com/devWaylander/coins_store/internal/middleware/auth -count=1 -v
```

Для запуска e2e тестов API выполните команду:

```bash
go test -timeout 30s -run ^TestE2eIntegrationTestSuite$ github.com/devWaylander/coins_store/internal/tests -count=1 -v
```

Либо воспользуйтесь встроенным плагином `Testing` для VSCode, в таком случае можно будет просмотреть дерево тестов и запустить их

## Настройка и запуск проекта

### Создание файла конфигурации

Для любого типа запуска необходимо создать файл `.env` в корне проекта и заполнить его в соответствии с шаблоном `.example.env`

### Локальный запуск проекта

Для запуска проекта локально выполните следующие шаги:

1. Выполните команду для поднятия контейнера с БД и миграциями:

   ```bash
   docker compose up postgresdb migrations
   ```

2. Откройте новую консоль и выполните команду для запуска приложения:

   ```bash
   cd cmd && go run main.go
   ```

3. Откройте новую консоль и выполните команду для запуска Swagger UI:

   ```bash
   make swaggerui
   ```

4. Откройте браузер и введите в адресную строку:

   ```
   http://localhost:5440
   ```

5. Проект будет развёрнут локально, с БД в контейнере и автоматическим применением миграций. SwaggerUI будет доступен для работы с API.

### Запуск проекта в контейнере

Для запуска проекта в контейнере выполните следующие шаги:

1. Выполните команду для поднятия всех контейнеров:

   ```bash
   docker compose up
   ```

2. Откройте новую консоль и выполните команду для запуска Swagger UI:

   ```bash
   make swaggerui
   ```

3. Откройте браузер и введите в адресную строку:

   ```
   http://localhost:5440
   ```

4. Проект будет развёрнут в контейнере, с БД в контейнере и автоматическим применением миграций. SwaggerUI будет доступен для работы с API.

## Требования к данным

### Username

- Должен содержать только латинские строчные и заглавные буквы, а также цифры.
- Должен быть не длиннее 64 символов.

Пример: `user1`

### Password

- Должен содержать минимум 8 символов.
- Должен включать латинские строчные и заглавные буквы.
- Должен содержать хотя бы один спецсимвол.
- Должен быть не менее 8 символов.

Пример: `Test123@`

## Секция вопросов

### Нагрузочное тестирование

- **Репозиторий с нагрузочным тестированием:**  
  [https://github.com/devWaylander/coins_store_load_testing](https://github.com/devWaylander/coins_store_load_testing)

- **Результаты нагрузочного тестирования:**  
  [Отчет о тестировании в формате HTML](https://github.com/devWaylander/coins_store_load_testing/blob/main/report.html)

**Примечание:**  
Я не до конца уверен в том, что библиотека Locust правильно агрегирует данные, но это был самый быстрый способ провести нагрузочное тестирование. Поэтому прошу воспринимать результаты скорее как примерные.

