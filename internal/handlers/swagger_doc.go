package handlers

// SwaggerDoc содержит Swagger документацию для API
const SwaggerDoc = `{
  "swagger": "2.0",
  "info": {
    "title": "Telescope Services API",
    "description": "API для управления телескопами и заявками на астрономические услуги",
    "version": "1.0.0"
  },
  "host": "localhost:8081",
  "basePath": "/api",
  "schemes": ["http"],
  "consumes": ["application/json"],
  "produces": ["application/json"],
  "securityDefinitions": {
    "Bearer": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header",
      "description": "JWT токен в формате: Bearer <token>"
    }
  },
  "paths": {
    "/services": {
      "get": {
        "tags": ["Услуги"],
        "summary": "Получение списка услуг",
        "description": "Получение списка доступных телескопов (публичный endpoint)",
        "responses": {
          "200": {
            "description": "Список услуг",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/Service"
              }
            }
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      },
      "post": {
        "tags": ["Услуги"],
        "summary": "Создание услуги",
        "description": "Создание новой услуги (только для модераторов)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "parameters": [
          {
            "name": "service",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/CreateServiceRequest"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Услуга создана",
            "schema": {
              "$ref": "#/definitions/Service"
            }
          },
          "401": {
            "description": "Не авторизован"
          },
          "403": {
            "description": "Недостаточно прав"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      }
    },
    "/services/{id}": {
      "get": {
        "tags": ["Услуги"],
        "summary": "Получение одной услуги",
        "description": "Получение информации об одной услуге (публичный endpoint)",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "type": "integer",
            "description": "ID услуги"
          }
        ],
        "responses": {
          "200": {
            "description": "Информация об услуге",
            "schema": {
              "$ref": "#/definitions/Service"
            }
          },
          "404": {
            "description": "Услуга не найдена"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      },
      "put": {
        "tags": ["Услуги"],
        "summary": "Обновление услуги",
        "description": "Обновление информации об услуге (только для модераторов)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "type": "integer",
            "description": "ID услуги"
          },
          {
            "name": "service",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/UpdateServiceRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Услуга обновлена",
            "schema": {
              "$ref": "#/definitions/Service"
            }
          },
          "401": {
            "description": "Не авторизован"
          },
          "403": {
            "description": "Недостаточно прав"
          },
          "404": {
            "description": "Услуга не найдена"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      },
      "delete": {
        "tags": ["Услуги"],
        "summary": "Удаление услуги",
        "description": "Удаление услуги (только для модераторов)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "type": "integer",
            "description": "ID услуги"
          }
        ],
        "responses": {
          "200": {
            "description": "Услуга удалена"
          },
          "401": {
            "description": "Не авторизован"
          },
          "403": {
            "description": "Недостаточно прав"
          },
          "404": {
            "description": "Услуга не найдена"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      }
    },
    "/services/{id}/image": {
      "post": {
        "tags": ["Услуги"],
        "summary": "Загрузка изображения услуги",
        "description": "Загрузка изображения для услуги (только для модераторов)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "type": "integer",
            "description": "ID услуги"
          },
          {
            "name": "file",
            "in": "formData",
            "required": true,
            "type": "file",
            "description": "Изображение услуги"
          }
        ],
        "responses": {
          "200": {
            "description": "Изображение загружено",
            "schema": {
              "type": "object",
              "properties": {
                "image_url": {
                  "type": "string",
                  "description": "URL загруженного изображения"
                }
              }
            }
          },
          "401": {
            "description": "Не авторизован"
          },
          "403": {
            "description": "Недостаточно прав"
          },
          "404": {
            "description": "Услуга не найдена"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      }
    },
    "/auth/register": {
      "post": {
        "tags": ["Авторизация"],
        "summary": "Регистрация пользователя",
        "description": "Создание нового пользователя в системе",
        "parameters": [
          {
            "name": "user",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/UserRegistration"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Пользователь успешно создан",
            "schema": {
              "$ref": "#/definitions/User"
            }
          },
          "400": {
            "description": "Неверные данные"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      }
    },
    "/auth/login": {
      "post": {
        "tags": ["Авторизация"],
        "summary": "Вход в систему",
        "description": "Аутентификация пользователя и получение JWT токена",
        "parameters": [
          {
            "name": "credentials",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/LoginRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Успешная аутентификация",
            "schema": {
              "$ref": "#/definitions/LoginResponse"
            }
          },
          "401": {
            "description": "Неверные учетные данные"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      }
    },
    "/auth/logout": {
      "post": {
        "tags": ["Авторизация"],
        "summary": "Выход из системы",
        "description": "Выход пользователя из системы",
        "responses": {
          "200": {
            "description": "Успешный выход"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      }
    },
    "/auth/me": {
      "get": {
        "tags": ["Авторизация"],
        "summary": "Получение данных пользователя",
        "description": "Получение данных текущего пользователя (требует авторизации)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "responses": {
          "200": {
            "description": "Данные пользователя",
            "schema": {
              "$ref": "#/definitions/User"
            }
          },
          "401": {
            "description": "Не авторизован"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      },
      "put": {
        "tags": ["Авторизация"],
        "summary": "Обновление данных пользователя",
        "description": "Обновление данных текущего пользователя (требует авторизации)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "parameters": [
          {
            "name": "user",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/UpdateUserRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Данные пользователя обновлены",
            "schema": {
              "$ref": "#/definitions/User"
            }
          },
          "401": {
            "description": "Не авторизован"
          },
          "400": {
            "description": "Неверные данные"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      }
    },
    "/orders": {
      "get": {
        "tags": ["Заявки"],
        "summary": "Получение списка заявок",
        "description": "Получение списка заявок пользователя (требует авторизации)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "responses": {
          "200": {
            "description": "Список заявок",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/Order"
              }
            }
          },
          "401": {
            "description": "Не авторизован"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      }
    },
    "/orders/{id}": {
      "get": {
        "tags": ["Заявки"],
        "summary": "Получение одной заявки",
        "description": "Получение информации об одной заявке (требует авторизации)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "type": "integer",
            "description": "ID заявки"
          }
        ],
        "responses": {
          "200": {
            "description": "Информация о заявке",
            "schema": {
              "$ref": "#/definitions/Order"
            }
          },
          "401": {
            "description": "Не авторизован"
          },
          "404": {
            "description": "Заявка не найдена"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      },
      "put": {
        "tags": ["Заявки"],
        "summary": "Обновление заявки",
        "description": "Обновление информации о заявке (требует авторизации)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "type": "integer",
            "description": "ID заявки"
          },
          {
            "name": "order",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/UpdateOrderRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Заявка обновлена",
            "schema": {
              "$ref": "#/definitions/Order"
            }
          },
          "401": {
            "description": "Не авторизован"
          },
          "404": {
            "description": "Заявка не найдена"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      },
      "delete": {
        "tags": ["Заявки"],
        "summary": "Удаление заявки",
        "description": "Удаление заявки (требует авторизации)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "type": "integer",
            "description": "ID заявки"
          }
        ],
        "responses": {
          "200": {
            "description": "Заявка удалена"
          },
          "401": {
            "description": "Не авторизован"
          },
          "404": {
            "description": "Заявка не найдена"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      }
    },
    "/orders/{id}/form": {
      "put": {
        "tags": ["Заявки"],
        "summary": "Формирование заявки",
        "description": "Формирование заявки из черновика (требует авторизации)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "type": "integer",
            "description": "ID заявки"
          }
        ],
        "responses": {
          "200": {
            "description": "Заявка сформирована",
            "schema": {
              "$ref": "#/definitions/Order"
            }
          },
          "401": {
            "description": "Не авторизован"
          },
          "404": {
            "description": "Заявка не найдена"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      }
    },
    "/orders/{id}/complete": {
      "put": {
        "tags": ["Заявки"],
        "summary": "Завершение заявки",
        "description": "Завершение заявки модератором (только для модераторов)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "type": "integer",
            "description": "ID заявки"
          },
          {
            "name": "action",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/CompleteOrderRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Заявка завершена",
            "schema": {
              "$ref": "#/definitions/Order"
            }
          },
          "401": {
            "description": "Не авторизован"
          },
          "403": {
            "description": "Недостаточно прав"
          },
          "404": {
            "description": "Заявка не найдена"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      }
    },
    "/orders/services": {
      "post": {
        "tags": ["Заявки"],
        "summary": "Добавление услуги в заявку",
        "description": "Добавление телескопа в черновую заявку (требует авторизации)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "parameters": [
          {
            "name": "service",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/AddServiceRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Услуга добавлена в заявку"
          },
          "401": {
            "description": "Не авторизован"
          },
          "409": {
            "description": "Услуга уже добавлена в заявку"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      }
    },
    "/orders/{order_id}/services/{service_id}": {
      "put": {
        "tags": ["Заявки"],
        "summary": "Обновление связи услуги с заявкой",
        "description": "Обновление дополнительной информации к услуге в заявке (требует авторизации)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "parameters": [
          {
            "name": "order_id",
            "in": "path",
            "required": true,
            "type": "integer",
            "description": "ID заявки"
          },
          {
            "name": "service_id",
            "in": "path",
            "required": true,
            "type": "integer",
            "description": "ID услуги"
          },
          {
            "name": "service",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/UpdateOrderServiceRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Связь обновлена"
          },
          "401": {
            "description": "Не авторизован"
          },
          "404": {
            "description": "Заявка или услуга не найдена"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      },
      "delete": {
        "tags": ["Заявки"],
        "summary": "Удаление услуги из заявки",
        "description": "Удаление услуги из заявки (требует авторизации)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "parameters": [
          {
            "name": "order_id",
            "in": "path",
            "required": true,
            "type": "integer",
            "description": "ID заявки"
          },
          {
            "name": "service_id",
            "in": "path",
            "required": true,
            "type": "integer",
            "description": "ID услуги"
          }
        ],
        "responses": {
          "200": {
            "description": "Услуга удалена из заявки"
          },
          "401": {
            "description": "Не авторизован"
          },
          "404": {
            "description": "Заявка или услуга не найдена"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      }
    },
    "/orders/cart": {
      "get": {
        "tags": ["Заявки"],
        "summary": "Получение корзины",
        "description": "Получение текущей корзины пользователя (требует авторизации)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "responses": {
          "200": {
            "description": "Корзина пользователя",
            "schema": {
              "$ref": "#/definitions/Cart"
            }
          },
          "401": {
            "description": "Не авторизован"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      }
    },
    "/admin/orders": {
      "get": {
        "tags": ["Администрирование"],
        "summary": "Получение всех заявок",
        "description": "Получение всех заявок для модераторов (только для модераторов)",
        "security": [
          {
            "Bearer": []
          }
        ],
        "responses": {
          "200": {
            "description": "Все заявки",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/Order"
              }
            }
          },
          "401": {
            "description": "Не авторизован"
          },
          "403": {
            "description": "Недостаточно прав"
          },
          "500": {
            "description": "Внутренняя ошибка сервера"
          }
        }
      }
    }
  },
  "definitions": {
    "Service": {
      "type": "object",
      "properties": {
        "id": {
          "type": "integer",
          "description": "ID услуги"
        },
        "name": {
          "type": "string",
          "description": "Название услуги"
        },
        "description": {
          "type": "string",
          "description": "Описание услуги"
        },
        "price": {
          "type": "number",
          "description": "Цена услуги"
        },
        "status": {
          "type": "string",
          "description": "Статус услуги"
        },
        "image_url": {
          "type": "string",
          "description": "URL изображения"
        },
        "created_at": {
          "type": "string",
          "format": "date-time",
          "description": "Дата создания"
        }
      }
    },
    "CreateServiceRequest": {
      "type": "object",
      "required": ["name", "description"],
      "properties": {
        "name": {
          "type": "string",
          "description": "Название услуги"
        },
        "description": {
          "type": "string",
          "description": "Описание услуги"
        },
        "price": {
          "type": "number",
          "description": "Цена услуги"
        }
      }
    },
    "UpdateServiceRequest": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string",
          "description": "Название услуги"
        },
        "description": {
          "type": "string",
          "description": "Описание услуги"
        },
        "price": {
          "type": "number",
          "description": "Цена услуги"
        }
      }
    },
    "UserRegistration": {
      "type": "object",
      "required": ["login", "password", "email", "first_name", "last_name"],
      "properties": {
        "login": {
          "type": "string",
          "description": "Логин пользователя"
        },
        "password": {
          "type": "string",
          "description": "Пароль пользователя"
        },
        "email": {
          "type": "string",
          "description": "Email пользователя"
        },
        "first_name": {
          "type": "string",
          "description": "Имя пользователя"
        },
        "last_name": {
          "type": "string",
          "description": "Фамилия пользователя"
        },
        "role": {
          "type": "string",
          "description": "Роль пользователя (user, moderator)",
          "default": "user"
        }
      }
    },
    "LoginRequest": {
      "type": "object",
      "required": ["login", "password"],
      "properties": {
        "login": {
          "type": "string",
          "description": "Логин пользователя"
        },
        "password": {
          "type": "string",
          "description": "Пароль пользователя"
        }
      }
    },
    "LoginResponse": {
      "type": "object",
      "properties": {
        "token": {
          "type": "string",
          "description": "JWT токен"
        },
        "user": {
          "$ref": "#/definitions/User"
        }
      }
    },
    "UpdateUserRequest": {
      "type": "object",
      "properties": {
        "login": {
          "type": "string",
          "description": "Логин пользователя"
        },
        "email": {
          "type": "string",
          "description": "Email пользователя"
        },
        "first_name": {
          "type": "string",
          "description": "Имя пользователя"
        },
        "last_name": {
          "type": "string",
          "description": "Фамилия пользователя"
        }
      }
    },
    "User": {
      "type": "object",
      "properties": {
        "id": {
          "type": "integer",
          "description": "ID пользователя"
        },
        "login": {
          "type": "string",
          "description": "Логин пользователя"
        },
        "email": {
          "type": "string",
          "description": "Email пользователя"
        },
        "first_name": {
          "type": "string",
          "description": "Имя пользователя"
        },
        "last_name": {
          "type": "string",
          "description": "Фамилия пользователя"
        },
        "role": {
          "type": "string",
          "description": "Роль пользователя"
        },
        "is_active": {
          "type": "boolean",
          "description": "Активен ли пользователь"
        },
        "created_at": {
          "type": "string",
          "format": "date-time",
          "description": "Дата создания"
        }
      }
    },
    "Order": {
      "type": "object",
      "properties": {
        "id": {
          "type": "integer",
          "description": "ID заявки"
        },
        "status": {
          "type": "string",
          "description": "Статус заявки"
        },
        "created_at": {
          "type": "string",
          "format": "date-time",
          "description": "Дата создания"
        },
        "creator_id": {
          "type": "integer",
          "description": "ID создателя"
        },
        "creator_login": {
          "type": "string",
          "description": "Логин создателя"
        },
        "formation_date": {
          "type": "string",
          "format": "date-time",
          "description": "Дата формирования"
        },
        "completion_date": {
          "type": "string",
          "format": "date-time",
          "description": "Дата завершения"
        },
        "services": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/OrderService"
          }
        }
      }
    },
    "UpdateOrderRequest": {
      "type": "object",
      "properties": {
        "exoplanet_name": {
          "type": "string",
          "description": "Название экзопланеты"
        },
        "star_mass": {
          "type": "number",
          "description": "Масса звезды"
        },
        "orbital_period": {
          "type": "number",
          "description": "Орбитальный период"
        },
        "velocity_amplitude": {
          "type": "number",
          "description": "Амплитуда скорости"
        },
        "inclination": {
          "type": "number",
          "description": "Наклонение орбиты"
        },
        "eccentricity": {
          "type": "number",
          "description": "Эксцентриситет орбиты"
        },
        "comment": {
          "type": "string",
          "description": "Комментарий"
        },
        "other_info": {
          "type": "string",
          "description": "Дополнительная информация"
        }
      }
    },
    "CompleteOrderRequest": {
      "type": "object",
      "required": ["action"],
      "properties": {
        "action": {
          "type": "string",
          "enum": ["complete", "reject"],
          "description": "Действие: complete - завершить, reject - отклонить"
        }
      }
    },
    "OrderService": {
      "type": "object",
      "properties": {
        "service_id": {
          "type": "integer",
          "description": "ID услуги"
        },
        "service_name": {
          "type": "string",
          "description": "Название услуги"
        },
        "added_at": {
          "type": "string",
          "format": "date-time",
          "description": "Дата добавления"
        },
        "description": {
          "type": "string",
          "description": "Дополнительное описание"
        }
      }
    },
    "AddServiceRequest": {
      "type": "object",
      "required": ["service_id"],
      "properties": {
        "service_id": {
          "type": "integer",
          "description": "ID услуги для добавления"
        },
        "exoplanet_name": {
          "type": "string",
          "description": "Название экзопланеты"
        },
        "star_mass": {
          "type": "number",
          "description": "Масса звезды"
        },
        "orbital_period": {
          "type": "number",
          "description": "Орбитальный период"
        },
        "velocity_amplitude": {
          "type": "number",
          "description": "Амплитуда скорости"
        },
        "inclination": {
          "type": "number",
          "description": "Наклонение орбиты"
        },
        "eccentricity": {
          "type": "number",
          "description": "Эксцентриситет орбиты"
        },
        "comment": {
          "type": "string",
          "description": "Комментарий"
        },
        "other_info": {
          "type": "string",
          "description": "Дополнительная информация"
        }
      }
    },
    "UpdateOrderServiceRequest": {
      "type": "object",
      "properties": {
        "description": {
          "type": "string",
          "description": "Дополнительное описание услуги"
        }
      }
    },
    "Cart": {
      "type": "object",
      "properties": {
        "order_id": {
          "type": "integer",
          "description": "ID заявки"
        },
        "services_count": {
          "type": "integer",
          "description": "Количество услуг в корзине"
        },
        "services": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/OrderService"
          }
        }
      }
    }
  }
}`