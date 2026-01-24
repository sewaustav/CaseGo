api schema

profile - post patch put delete get
socials - post put delete
purposes - post put delete

POST /profile   →   Создание профиля пользователя
Авторизация: требуется Bearer JWT
Тело запроса (JSON):

{
  "info": {                  // обязательный объект
    "avatar": "string",      // обязательное, URL аватарки
    "username": "string",    // обязательное, 3–30 символов
    "name": "string",        // обязательное, имя
    "surname": "string",     // обязательное, фамилия
    "patronymic": "string?", // необязательное, отчество (может отсутствовать)
    "email": "string",       // обязательное, валидный email
    "phone_number": "string?", // необязательное, формат E.164 (+79161234567)
    "sex": 0 | 1 | null,     // необязательное, 1 = женский, 0 = мужской
    "description": "string?", // необязательное, до 500 символов
    "profession": "string?"  // необязательное, профессия / должность
  },

  "social_links": [          // необязательный массив, может быть пустым
    {
      "type": "string",      // обязательное, например: "telegram", "instagram", "vk", "github"...
      "url":  "string"       // обязательное, валидный URL
    },
    ...
  ],

  "purposes": [              // обязательный массив, минимум 1 элемент
    {
      "purpose": "string"    // обязательное, минимум 5 символов
    },
    ...
  ]
}

GET /profile   →   получить профиль
Авторизация: требуется Bearer JWT
Параметры в URL: отсутствуют (user_id берётся из токена)

GET /profile/{id} - запрос по profile id(только фдмин)



PATCH /profile  →   Частичное обновление

Авторизация: требуется Bearer JWT
Тело запроса: JSON с полями, которые нужно обновить (все поля опциональные)

{
  "avatar":       "string?",      // новый URL аватарки
  "username":     "string?",      // 3–30 символов, если передан
  "name":         "string?",      // имя
  "surname":      "string?",      // фамилия
  "patronymic":   "string?",      // отчество (null → удалить, если разрешено)
  "email":        "string?",      // валидный email
  "phone_number": "string?",      // E.164 формат
  "sex":          0 | 1 | null?,  // 0 = женский, 1 = мужской, null = не указан
  "description":  "string?",      // до 500 символов
  "profession":   "string?"       // профессия / должность
}

PUT /profile  →   Полное обновление

Авторизация: требуется Bearer JWT

Тело запроса: JSON — полный объект профиля (почти как в CreateProfileRequest → info, но без social_links и purposes)
Все поля обязательные
Если какое-то необязательное поле нужно убрать/обнулить — null или пустую строку "".

{
  "avatar":       "string",       // ОБЯЗАТЕЛЬНО — URL аватарки
  "username":     "string",       // ОБЯЗАТЕЛЬНО — 3–30 символов
  "name":         "string",       // ОБЯЗАТЕЛЬНО — имя
  "surname":      "string",       // ОБЯЗАТЕЛЬНО — фамилия
  "patronymic":   "string?",      // необязательно, можно null или отсутствовать
  "email":        "string",       // ОБЯЗАТЕЛЬНО — валидный email
  "phone_number": "string?",      // необязательно, E.164 формат, можно null
  "sex":          0 | 1 | null?,  // необязательно, можно null
  "description":  "string",       // можно пустую строку "", max 500 символов
  "profession":   "string?"       // необязательно, можно null
}


DELETE /profile   →   Удаление собственного профиля текущим пользователем

Авторизация: требуется Bearer JWT 

Параметры в URL: отсутствуют (ID пользователя берётся исключительно из JWT)

DELETE /profile/?id= || /profile/?user_id=  →   Полное и окончательное удаление любого профиля (только для админов)

Авторизация: требуется Bearer JWT (только админ)
query - user_id или id профиля


POST /me/social-links   →   Добавление одной или нескольких социальных ссылок
Авторизация: требуется Bearer JWT
Тело запроса (можно массив или один объект — оба варианта поддерживаются):

Вариант 1: одна ссылка
{
  "type": "string",     // обязательно, например: "telegram", "instagram", "github", "linkedin", "twitter", "vk"...
  "url":  "string"      // обязательно, валидный URL
}

Вариант 2: сразу несколько (массив)
{
  "social_links": [
    {
      "type": "telegram",
      "url": "https://t.me/sewaustaff"
    },
    {
      "type": "github",
      "url": "https://github.com/sewaustav"
    }
  ]
}

POST /me/purposes   →   Добавление одной или нескольких целей
Авторизация: требуется Bearer JWT
Тело запроса:

Вариант 1: одна цель
{
  "purpose": "string"    // обязательно, минимум 5 символов
}

Вариант 2: несколько сразу
{
  "purposes": [
    { "purpose": "попрогать" },
    { "purpose": "подрочить" },
    { "purpose": "сыграть каточку в рояль" }
  ]
}