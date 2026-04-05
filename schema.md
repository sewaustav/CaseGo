Вот краткое описание проекта для фронтенд-разработки:

---

# CaseGo — Платформа для прокачки soft skills через кейс-интервью

## Архитектура: 4 микросервиса

| Сервис | Язык | Порт | Назначение |
|--------|------|------|-----------|
| **Auth** | Python/FastAPI | 8000 | Аутентификация, JWT, OAuth Google |
| **Profile** | Go/Gin | 8080 | Профили пользователей, поиск |
| **CaseGo** | Go/Gin | 8081 | Кейсы, диалоги с ИИ |
| **CaseProfile** | Go/Gin | 8082 | Статистика и результаты |

**Авторизация:** Bearer JWT (RS256). Токен содержит `user_id`, `user_role` (0=Admin, 1=User, 2=Creator).

---

## Auth Service `/api/v1`

### Пользовательские эндпоинты

| Метод | Путь | Описание |
|-------|------|---------|
| POST | `/auth/register` | Регистрация `{username, email, password}` → `{id, username, email}` |
| POST | `/auth/token` | Логин form-data `{username, password}` → `{access_token, refresh_token, token_type, expires_in}` |
| POST | `/auth/refresh` | Обновить токен `{refresh_token}` → `{access_token, ...}` |
| POST | `/auth/auth/google` | Google OAuth `{id_token}` → `{access_token, ...}` |
| GET | `/users/me` | Текущий пользователь → `{id, username, email}` |
| GET | `/health/` | Healthcheck |

---

## Profile Service `/profile/api/v1`

### Пользовательские эндпоинты

#### Профиль
| Метод | Путь | Описание |
|-------|------|---------|
| POST | `/profile` | Создать профиль |
| GET | `/profile` | Получить свой профиль |
| PUT | `/profile` | Полное обновление профиля |
| PATCH | `/profile` | Частичное обновление |
| DELETE | `/profile` | Мягкое удаление (деактивация) |

**Тело создания/обновления профиля:**
```json
{
  "info": {
    "avatar": "url",
    "username": "min3-max30",
    "name": "string",
    "surname": "string",
    "patronymic": "string?",
    "city": "string?",
    "age": 14-120,
    "sex": 0 | 1,
    "description": "max500",
    "profession": "string?"
  },
  "social_links": [{"type": "telegram", "url": "https://..."}],
  "purposes": [{"purpose": "min5chars"}]
}
```

#### Соцсети
| Метод | Путь | Описание |
|-------|------|---------|
| POST | `/profile/social` | Добавить ссылки (массив) |
| PUT | `/profile/social/:id` | Обновить ссылку |
| DELETE | `/profile/social/:id` | Удалить ссылку |

#### Цели
| Метод | Путь | Описание |
|-------|------|---------|
| POST | `/profile/purpose` | Добавить цели (массив) |
| PUT | `/profile/purpose/:id` | Обновить цель |
| DELETE | `/profile/purpose/:id` | Удалить цель |

#### Профессии/категории
| Метод | Путь | Описание |
|-------|------|---------|
| POST | `/profile/profession` | Добавить профессию (массив `{profession_id}`) |
| GET | `/profile/profession` | Получить свои профессии |
| PUT | `/profile/profession/:id` | Изменить профессию |
| DELETE | `/profile/profession/:id` | Удалить профессию |

#### Поиск
| Метод | Путь | Query-параметры |
|-------|------|----------------|
| GET | `/search` | `profession_id, profession, min_age, max_age, city, sex, limit, page, order_by, order_direction` |
| GET | `/search/fio` | `name, surname, patronymic, limit, page` |

#### Категории (публичные)
| Метод | Путь | Описание |
|-------|------|---------|
| GET | `/profession_categories` | Все категории |
| GET | `/profession_categories/:id` | Категория по ID |
| GET | `/profession_categories/parent/:id` | Подкатегории |

### Админские эндпоинты

| Метод | Путь | Описание |
|-------|------|---------|
| GET | `/profile/:id` | Просмотр любого профиля по ID |
| DELETE | `/profile/:id` | Полное (hard) удаление профиля |
| POST | `/profession_categories` | Создать категорию `{name, parent_id?}` |

---

## CaseGo Service `/api/v1/case_go`

### Пользовательские эндпоинты

| Метод | Путь | Описание |
|-------|------|---------|
| GET | `/cases` | Список кейсов `{limit, page, topic?, category?}` |
| GET | `/cases/:caseID` | Кейс по ID |
| POST | `/cases` | Начать диалог по кейсу `{case_id}` → возвращает кейс с первым вопросом |
| POST | `/dialogs/:dialogID/complete` | Завершить диалог → результаты |
| GET | `/dialogs/:dialogID` | Получить диалог с историей |
| GET | `/users/:userID/dialogs` | Список диалогов пользователя (свои) `{limit, page}` |

**Тело шага диалога (интеракция):**
```json
{
  "dialog_id": 1,
  "step": 3,
  "question": "string",
  "answer": "string"
}
```

### Админские / Creator эндпоинты

| Метод | Путь | Описание |
|-------|------|---------|
| POST | `/case` | Создать кейс вручную или через ИИ (`{prompt}` или `{topic, category, description, first_question}`) |
| PUT | `/case/:caseID` | Обновить кейс |
| DELETE | `/case/:caseID` | Удалить кейс (только Admin) |
| GET | `/users/:userID/dialogs` | Просмотр диалогов любого пользователя (Admin) |

---

## CaseProfile Service `/api/v1/case_go`

### Пользовательские эндпоинты

| Метод | Путь | Описание |
|-------|------|---------|
| GET | `/profile` | Своя статистика навыков |
| GET | `/history` | История результатов `?from=2026-01-01` |

**Ответ профиля навыков:**
```json
{
  "user_id": 1,
  "total_cases": 42,
  "assertiveness": 0.8,
  "empathy": 0.7,
  "clarity_communication": 0.6,
  "resistance": 0.5,
  "eloquence": 0.4,
  "initiative": 0.3
}
```

### Админские эндпоинты

| Метод | Путь | Описание |
|-------|------|---------|
| GET | `/admin/profile` | Профиль пользователя по `?user_id=` или `?id=` |
| GET | `/admin/history` | История пользователя по `?user_id=` |
| DELETE | `/admin/result/:id` | Удалить результат |

---

## Модели данных (ключевые)

### Profile
```json
{
  "id": 1,
  "user_id": 1,
  "avatar": "url",
  "is_active": true,
  "username": "string",
  "name": "string",
  "surname": "string",
  "patronymic": "string|null",
  "city": "string|null",
  "age": 25,
  "sex": 0,
  "description": "string",
  "profession": "string|null",
  "case_count": 10,
  "created_at": "iso8601",
  "updated_at": "iso8601"
}
```

### UserProfile (полный)
```json
{
  "UsrProfile": { /* Profile */ },
  "UsrPurposes": [{"id": 1, "user_id": 1, "purpose": "string"}],
  "UsrSocials": [{"id": 1, "user_id": 1, "type": "telegram", "url": "https://..."}]
}
```

### Case
```json
{
  "id": 1,
  "topic": "string",
  "category": 1,
  "description": "string",
  "first_question": "string",
  "is_generated": false,
  "creator": 1,
  "created_at": "iso8601"
}
```

---

## Коды ошибок

| Код | Ситуация |
|-----|---------|
| 400 | Невалидные данные |
| 401 | Нет/истёк токен |
| 403 | Нет прав (не владелец / не Admin) |
| 404 | Не найдено / профиль деактивирован |
| 409 | Конфликт (username/email уже занят) — тело: `{"error": "Conflict", "field": "username", "message": "..."}` |
| 500 | Внутренняя ошибка |

---

## Роли

| Значение | Роль | Возможности |
|---------|------|-----------|
| 0 | Admin | Всё, включая hard-delete, просмотр любых профилей |
| 1 | User | CRUD своего профиля, прохождение кейсов |
| 2 | Creator | Создание и редактирование кейсов |