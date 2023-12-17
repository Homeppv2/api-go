# API

## Краткое описание
Основной сервис API

## Настройки сервиса

```cgo
// Основная структура настроек
Config struct {
    // Структура настроек логера
    Logger   Logger
    // Структура настроек HTTP-сервера
    HTTP     HTTP
    // Структура настроек БД PostgreSQL
    Postgres Postgres
    // Структура настроек Redis
    Redis    Redis
    // Структура настроек авторизации
    Auth     Auth
}

// Структура настроек логера
Logger struct {
    // Уровень логов
    Level slog.Level
}

// Структура настроек HTTP-сервера
HTTP struct {
    // Имя хоста
    Host string
    // Номер порта
    Port string
}

// Структура настроек БД PostgreSQL
Postgres struct {
    // Строка соединения
    ConnString      string
    // Максимальное количество соединений
    MaxOpenConns    int
    // Максимальный период соедиения в режиме Lifetime
    ConnMaxLifetime time.Duration
    // Число соединений, которое разрешается иметь в состоянии Idle (т.е. открытых TCP-соединений, которые в данный момент не используются)
    MaxIdleConns    int
    // Максимальный период соедиения в режиме Idle
    ConnMaxIdleTime time.Duration
    // Флаг включения/выключения автомиграций
    AutoMigrate     bool
    // Путь к файлам миграций
    MigrationsPath  string
}

// Структура настроек Redis
Redis struct {
    // Имя хоста
    Host     string
    // Имя порта
    Port     string
    // Пароль
    Password string
    // Название БД
    DB       int
}

// Структура настроек авторизации
Auth struct {
    // Структура настроек токенов авторизации
    Token Token
}

// Структура настроек токенов авторизации
Token struct {
    // Время, которое "живет" токен
    ExpiresIn time.Duration
}
```


## Описание API `/api`

---

### Auth `/auth`
Сервис авторизации
- **login()** `POST /auth/login `

Метод авторизует пользователя в системе

- **logout()** `POST /auth/logout`

Метод разлогинивает пользователя в системе

---

### Users `/users`
Сервис пользователей
- **register()** `POST /users/register `

Метод регистрирует пользователя в системе

- **getByID()** `GET /users/id/:id`

Метод возвращает пользователя по идентификатору

- **getByEmail()** `GET /users/email/:email`

Метод возвращает пользователя по почте

---

### Controllers `/controllers`
Сервис контроллеров
- **create()** `POST /controllers/create `

Метод создает контроллер с уникальным ключом HwKey

- **getByID()** `GET /controllers/id/:id`

Метод возвращает контроллер по идентификатору

- **getByHwKey()** `GET /controllers/hw-key/:hwKey`

Метод возвращает контроллер по уникальному ключу HwKey

- **getByIsUsedBy()** `GET /controllers/is-used-by/:isUsedBy `

Метод возвращает контроллеры по идентификатору пользователя, которому они принадлежат

- **updateIsUsed()** `PATCH /controllers/is-used-by/`

Метод меняет у определенного контроллера пользователя, который использует контроллер

- **delete()** `DELETE /controllers/:id`

Метод удаляет контроллер по идентификатору

---

## Changelog

### 1.0.0
- Создан сервис