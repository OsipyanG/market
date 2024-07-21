# Микросервис склада

## Обзор

Микросервис склада состоит из трех основных сервисов:

1. **Сервис каталога (Catalog Service)**
2. **Сервис склада (Warehouse Service)**
3. **Сервис управления складом (WarehouseAdmin Service)**

## Сервисы и методы RPC

### Сервис каталога (Catalog Service)

Сервис `Catalog` отвечает за получение каталога продуктов.

#### Методы RPC

- **GetCatalog**: Получает полный каталог продуктов с возможностью пагинации.
  - **Запрос**: `GetCatalogRequest`
  - **Ответ**: `GetCatalogResponse`

#### Сообщения

- **GetCatalogRequest**: Запрос на получение каталога продуктов.
  - `int32 offset`: Смещение для пагинации.
  - `int32 limit`: Лимит количества возвращаемых продуктов.
- **GetCatalogResponse**: Ответ на запрос каталога продуктов.
  - `repeated Product products`: Список продуктов.

### Сервис склада (Warehouse Service)

Сервис `Warehouse` отвечает за управление резервированием продуктов на складе.

#### Методы RPC

- **ReserveProducts**: Резервирует указанные продукты.
  - **Запрос**: `ReserveProductsRequest`
  - **Ответ**: `google.protobuf.Empty`
- **FreeReservedProducts**: Освобождает зарезервированные продукты.
  - **Запрос**: `FreeReservedProductsRequest`
  - **Ответ**: `google.protobuf.Empty`
- **DeleteReservedProducts**: Удаляет зарезервированные продукты.
  - **Запрос**: `DeleteReservedProductsRequest`
  - **Ответ**: `google.protobuf.Empty`

Примечание: освобождение необходимо в случае отмены заказа, а удаление, когда заказ успешно доставлен.

#### Сообщения

- **ReserveProductsRequest**: Запрос на резервирование продуктов.
  - `jwt.JWTClaims jwt_claims`: Данные из JWT-токена для аутентификации.
  - `repeated ProductQuantity products`: Список продуктов для резервирования.
- **FreeReservedProductsRequest**: Запрос на освобождение зарезервированных продуктов.
  - `jwt.JWTClaims jwt_claims`: Данные из JWT-токена для аутентификации.
  - `repeated ProductQuantity products`: Список продуктов для освобождения.
- **DeleteReservedProductsRequest**: Запрос на удаление зарезервированных продуктов.
  - `jwt.JWTClaims jwt_claims`: Данные из JWT-токена для аутентификации.
  - `repeated ProductQuantity products`: Список продуктов для удаления.

### Сервис управления складом (WarehouseAdmin Service)

Сервис `WarehouseAdmin` отвечает за управление продуктами на складе.

#### Методы RPC

- **CreateProduct**: Создает новый продукт.
  - **Запрос**: `CreateProductRequest`
  - **Ответ**: `google.protobuf.Empty`
- **AddProduct**: Добавляет количество к существующему продукту.
  - **Запрос**: `AddProductRequest`
  - **Ответ**: `google.protobuf.Empty`
- **DeleteProduct**: Удаляет количество из существующего продукта.
  - **Запрос**: `DeleteProductRequest`
  - **Ответ**: `google.protobuf.Empty`

#### Сообщения

- **CreateProductRequest**: Запрос на создание нового продукта.
  - `jwt.JWTClaims jwt_claims`: Данные из JWT-токена для аутентификации.
  - `NewProduct new_product`: Информация о новом продукте.
- **AddProductRequest**: Запрос на добавление количества к продукту.
  - `jwt.JWTClaims jwt_claims`: Данные из JWT-токена для аутентификации.
  - `ProductQuantity product`: Информация о продукте и количестве для добавления.
- **DeleteProductRequest**: Запрос на удаление количества из продукта.
  - `jwt.JWTClaims jwt_claims`: Данные из JWT-токена для аутентификации.
  - `ProductQuantity product`: Информация о продукте и количестве для удаления.

### Общие сообщения

- **Product**: Представление продукта.

  - `string product_id`: Идентификатор продукта.
  - `string name`: Название продукта.
  - `string description`: Описание продукта.
  - `int64 available`: Доступность продукта.
  - `int64 quantity`: Количество продукта.
  - `int64 price`: Цена продукта.

- **NewProduct**: Информация о новом продукте.

  - `string name`: Название продукта.
  - `string description`: Описание продукта.
  - `int64 quantity`: Количество продукта.
  - `int64 price`: Цена продукта.

- **ProductQuantity**: Количество продукта.
  - `string product_id`: Идентификатор продукта.
  - `int64 quantity`: Количество продукта.
