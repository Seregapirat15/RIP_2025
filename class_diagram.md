# Диаграмма классов бэкенда системы управления услугами

## Mermaid диаграмма классов

```mermaid
classDiagram
    %% Основные модели данных
    class User {
        +ID: int
        +Username: string
        +Email: string
        +Password: string
        +FullName: string
        +Role: string
        +CreatedAt: time.Time
        +UpdatedAt: time.Time
        +Register()
        +Login()
        +Update()
        +Delete()
    }

    class Service {
        +ID: int
        +Name: string
        +Description: string
        +Price: float64
        +Category: string
        +ImageURL: string
        +CreatedAt: time.Time
        +UpdatedAt: time.Time
        +Create()
        +Update()
        +Delete()
        +UploadImage()
    }

    class Order {
        +ID: int
        +UserID: int
        +Status: string
        +Title: string
        +Description: string
        +TotalCost: float64
        +DeliveryDate: time.Time
        +CreatedAt: time.Time
        +FormedAt: time.Time
        +CompletedAt: time.Time
        +CreatorID: int
        +ModeratorID: int
        +Create()
        +Update()
        +Form()
        +Complete()
        +Delete()
    }

    class OrderService {
        +ID: int
        +OrderID: int
        +ServiceID: int
        +Quantity: int
        +Price: float64
        +CreatedAt: time.Time
        +AddToOrder()
        +UpdateQuantity()
        +RemoveFromOrder()
    }

    %% API Handlers
    class APIHandlers {
        +GetServicesHandler()
        +GetServiceHandler()
        +CreateServiceHandler()
        +UpdateServiceHandler()
        +DeleteServiceHandler()
        +UploadServiceImageHandler()
        +GetCartIconHandler()
        +GetOrdersHandler()
        +GetOrderHandler()
        +UpdateOrderHandler()
        +FormOrderHandler()
        +CompleteOrderHandler()
        +DeleteOrderHandler()
        +AddServiceToOrderHandler()
        +DeleteOrderServiceHandler()
        +UpdateOrderServiceHandler()
        +RegisterUserHandler()
        +GetUserHandler()
        +UpdateUserHandler()
        +LoginUserHandler()
        +LogoutUserHandler()
    }

    %% Database Queries
    class ServiceQueries {
        +GetServices(filters)
        +GetService(id)
        +CreateService(service)
        +UpdateService(id, service)
        +DeleteService(id)
        +UploadServiceImage(id, image)
    }

    class OrderQueries {
        +GetOrders(filters)
        +GetOrder(id)
        +CreateOrder(order)
        +UpdateOrder(id, order)
        +FormOrder(id)
        +CompleteOrder(id, action)
        +DeleteOrder(id)
        +GetCartIcon(userID)
    }

    class OrderServiceQueries {
        +AddServiceToOrder(orderID, serviceID)
        +UpdateOrderService(orderID, serviceID, data)
        +DeleteOrderService(orderID, serviceID)
        +GetOrderServices(orderID)
    }

    class UserQueries {
        +RegisterUser(user)
        +GetUser(id)
        +UpdateUser(id, user)
        +LoginUser(credentials)
        +LogoutUser(id)
    }

    %% Database Connection
    class DatabaseConnection {
        +DB: *sql.DB
        +Connect()
        +Close()
        +Ping()
    }

    %% MinIO Integration
    class MinIOClient {
        +Endpoint: string
        +AccessKey: string
        +SecretKey: string
        +Bucket: string
        +UploadFile(file)
        +DeleteFile(filename)
        +GetFileURL(filename)
    }

    %% Configuration
    class Config {
        +DatabaseURL: string
        +MinIOEndpoint: string
        +MinIOAccessKey: string
        +MinIOSecretKey: string
        +MinIOBucket: string
        +ServerPort: string
        +Load()
    }

    %% Relationships
    User ||--o{ Order : creates
    User ||--o{ Order : moderates
    Service ||--o{ OrderService : included_in
    Order ||--o{ OrderService : contains
    
    APIHandlers --> ServiceQueries : uses
    APIHandlers --> OrderQueries : uses
    APIHandlers --> OrderServiceQueries : uses
    APIHandlers --> UserQueries : uses
    APIHandlers --> MinIOClient : uses
    
    ServiceQueries --> DatabaseConnection : uses
    OrderQueries --> DatabaseConnection : uses
    OrderServiceQueries --> DatabaseConnection : uses
    UserQueries --> DatabaseConnection : uses
    
    DatabaseConnection --> Config : uses
    MinIOClient --> Config : uses
```

## Описание доменов методов по URL

### 1. Домен услуг (/api/services)
- **GET /api/services** - получение списка услуг с фильтрацией
- **GET /api/services/{id}** - получение одной услуги
- **POST /api/services** - создание новой услуги
- **PUT /api/services/{id}** - обновление услуги
- **DELETE /api/services/{id}** - удаление услуги
- **POST /api/services/{id}/image** - загрузка изображения услуги

### 2. Домен заявок (/api/orders)
- **GET /api/orders/cart** - получение иконки корзины
- **GET /api/orders** - получение списка заявок с фильтрацией
- **GET /api/orders/{id}** - получение одной заявки с услугами
- **PUT /api/orders/{id}** - изменение заявки
- **PUT /api/orders/{id}/form** - формирование заявки
- **PUT /api/orders/{id}/complete** - завершение/отклонение заявки
- **DELETE /api/orders/{id}** - удаление заявки

### 3. Домен связей заявка-услуга (/api/orders/services)
- **POST /api/orders/services** - добавление услуги в заявку
- **PUT /api/orders/{order_id}/services/{service_id}** - изменение связи м-м
- **DELETE /api/orders/{order_id}/services/{service_id}** - удаление связи м-м

### 4. Домен пользователей (/api/users)
- **POST /api/users/register** - регистрация пользователя
- **GET /api/users/me** - получение данных пользователя
- **PUT /api/users/me** - обновление пользователя
- **POST /api/users/login** - аутентификация
- **POST /api/users/logout** - деавторизация

## Модели и таблицы БД

### Таблица users
- id (PRIMARY KEY)
- username (UNIQUE)
- email (UNIQUE)
- password_hash
- full_name
- role
- created_at
- updated_at

### Таблица services
- id (PRIMARY KEY)
- name
- description
- price
- category
- image_url
- created_at
- updated_at

### Таблица orders
- id (PRIMARY KEY)
- user_id (FOREIGN KEY -> users.id)
- status
- title
- description
- total_cost
- delivery_date
- created_at
- formed_at
- completed_at
- creator_id (FOREIGN KEY -> users.id)
- moderator_id (FOREIGN KEY -> users.id)

### Таблица order_services (связь м-м)
- id (PRIMARY KEY)
- order_id (FOREIGN KEY -> orders.id)
- service_id (FOREIGN KEY -> services.id)
- quantity
- price
- created_at

## Связи между моделями

### 1. Методы используют разные модели
- `GetOrdersHandler` использует `Order` и `User` модели
- `GetOrderHandler` использует `Order`, `OrderService` и `Service` модели
- `CompleteOrderHandler` использует `Order`, `OrderService` и `Service` модели

### 2. Модели используют другие модели
- `Order` модель связана с `User` через `CreatorID` и `ModeratorID`
- `OrderService` модель связана с `Order` и `Service`

### 3. Модели используют несколько таблиц
- `OrderService` модель использует таблицы `orders`, `services` и `order_services`
- `GetOrderHandler` использует JOIN запросы для получения заявки с услугами
- `CompleteOrderHandler` использует транзакции для обновления нескольких таблиц

## Статусы заявок и переходы

### Статусы:
- **черновик** - создана, но не сформирована
- **сформированная** - сформирована, ожидает модерации
- **завершенная** - одобрена модератором
- **отклоненная** - отклонена модератором
- **удаленная** - удалена создателем

### Переходы статусов:
- **создатель**: черновик → сформированная, черновик → удаленная
- **модератор**: сформированная → завершенная, сформированная → отклоненная

## Бизнес-логика

### При завершении заявки:
1. Вычисляется общая стоимость заказа
2. Рассчитывается дата доставки (в течение месяца)
3. Обновляется статус заявки
4. Проставляется модератор и дата завершения

### При формировании заявки:
1. Проверяются обязательные поля
2. Обновляется статус на "сформированная"
3. Проставляется дата формирования
