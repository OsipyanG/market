# Микросервис авторизации и аутентификации

## Обзор

Микросервис авторизации и аутентфикации состоит из двух сервисов:

1. **Сервис аутентификации (Auth Service)**
2. **Административный сервис аутентификации (AuthAdmin Service)**

Доступ к микросервису осуществляется с помощью gRPC

## Сервисы и методы gRPC

### Сервис аутентификации (Auth Service)

Сервис `Auth` отвечает за операции, связанные с аутентификацией и авторизацией пользователей, обновлением паролей и тд.

#### Методы RPC

- **NewUser**: Регистрирует нового пользователя. Автоматически входит в систему.
    - **Запрос**: `UserCredentials`
    - **Ответ**: `Tokens`

- **Login**: Выполняет вход пользователя в систему.
    - **Запрос**: `UserCredentials`
    - **Ответ**: `Tokens`

- **UpdateTokens**: Обновляет токены доступа и обновления.
    - **Запрос**: `RefreshToken`
    - **Ответ**: `Tokens`

- **UpdatePassword**: Обновляет пароль пользователя.
    - **Запрос**: `RequestUpdatePassword`
    - **Ответ**: `google.protobuf.Empty`

- **Logout**: Выполняет выход пользователя из системы.
    - **Запрос**: `RefreshToken`
    - **Ответ**: `google.protobuf.Empty`

- **GetJWTClaims**: Получает данные из JWT-токена.
    - **Запрос**: `AccessToken`
    - **Ответ**: `jwt.JWTClaims`

### proto Messages

- **UserCredentials**: Учетные данные пользователя.

```proto
message UserCredentials {
    string login = 1;  
    string password = 2; 
}
```

- **Tokens**: Токены доступа и обновления.

```proto
message Tokens {
    string access_token = 1;
    string refresh_token = 2;
}
```

- **RefreshToken**: Токен обновления.

```proto
message RefreshToken {
    string value = 1;
}
```

- **AccessToken**: Токен доступа.

```proto
message AccessToken {
    string value = 1;
}
```

- **RequestUpdatePassword**: Запрос на обновление пароля.

```proto
message RequestUpdatePassword {
    jwt.JWTClaims jwt_claims = 1;
    string old_password = 2;
    string new_password = 3;
}
```

### Административный сервис аутентификации (AuthAdmin Service)

Сервис `AuthAdmin` отвечает за административные операции, такие как удаление пользователей, установка уровней доступа и
получение списка пользователей с определенным уровнем доступа.

#### Методы RPC

- **DeleteUser**: Удаляет пользователя по UserID
    - **Запрос**: `RequestByUserID`
    - **Ответ**: `google.protobuf.Empty`

- **SetAccessLevel**: Устанавливает уровень доступа пользователя.
    - **Запрос**: `SetAccessLevelRequest`
    - **Ответ**: `google.protobuf.Empty`

- **GetAllUsersWithLevel**: Получает список всех пользователей с указанным уровнем доступа.
    - **Запрос**: `RequestByLevel`
    - **Ответ**: `UsersInfoResponse`

#### Сообщения

- **RequestByUserID**: Запрос на получение информации о пользователе по идентификатору.
  user_id - идентификатор пользователя (UUID), чьи данные необходимо получить.

```proto
message RequestByUserID {
    jwt.JWTClaims jwt_claims = 1;
    string user_id = 2;
}
```

- **SetAccessLevelRequest**: Запрос на установку уровня доступа пользователю.
  user_id - идентификатор пользователя (UUID), чей уровень доступа необходимо изменить.

```proto
message SetAccessLevelRequest {
    jwt.JWTClaims jwt_claims = 1;
    string user_id = 2; // Идентификатор пользователя (UUID).
    uint32 lvl = 3;
}
```

- **User**: Информация о пользователе.

```proto
message User {
    string id = 1; // UUID
    string login = 2;
    uint32 access_level = 3;
}
```

- **UsersInfoResponse**: Ответ с информацией о пользователях.

```proto
message UsersInfoResponse {
    repeated User users = 1;
}
```

- **RequestByLevel**: Запрос на получение пользователей по уровню доступа.

```proto
message RequestByLevel {
    jwt.JWTClaims jwt_claims = 1;
    uint32 lvl = 2;
}
```

### Уровни доступа

- **Admin - 1**
- **WarehouseWorker - 2**
- **Courier - 3**
- **Buyer - 4**
