# API Documentation

## Authentication

### POST `/api/login`
- **Input:**
  ```json
  { "username": "string", "password": "string" }
  ```
- **Output:**
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "token": "<jwt>",
      "user": { ...user info... }
    }
  }
  ```

---

## User APIs

### POST `/api/users/create` (admin only)
- **Input:**
  ```json
  { "username": "string", "password": "string", "full_name": "string", "email": "string" }
  ```
- **Output:**
  ```json
  { "code": 0, "message": "success", "data": { "user": { ...user info... } } }
  ```

### POST `/api/users/delete` (admin only)
- **Input:**
  ```json
  { "user_id": "string" }
  ```
- **Output:**
  ```json
  { "code": 0, "message": "user deleted successfully" }
  ```

### POST `/api/users/change-password`
- **Input:**
  ```json
  { "user_id": "string", "new_password": "string" }
  ```
- **Output:**
  ```json
  { "code": 0, "message": "password changed successfully" }
  ```

### POST `/api/users/update`
- **Input:**
  ```json
  { ...user object... }
  ```
- **Output:**
  ```json
  { "code": 0, "message": "user updated successfully" }
  ```

### POST `/api/users/update-info`
- **Input:**
  ```json
  { "username": "string", "full_name": "string (optional)", "email": "string (optional)" }
  ```
- **Output:**
  ```json
  { "code": 0, "message": "user info updated successfully" }
  ```

### GET `/api/users` (admin only)
- **Output:**
  ```json
  { "code": 0, "message": "success", "data": [ { ...user info... } ] }
  ```

---

## Client APIs

### GET `/api/clients`
- **Output:**
  ```json
  { "code": 0, "message": "success", "data": { "clients": [ ... ] } }
  ```

### GET `/api/clients/:agent_id`
- **Output:**
  ```json
  { "code": 0, "message": "success", "data": { ...client info... } }
  ```

### GET `/api/clients/by-user/:user_id`
- **Output:**
  ```json
  { "code": 0, "message": "success", "data": [ ...client list... ] }
  ```

### POST `/api/clients/delete` (admin only)
- **Input:**
  ```json
  { "agent_id": "string" }
  ```
- **Output:**
  ```json
  { "code": 0, "message": "client deleted successfully" }
  ```

### POST `/api/clients/assign-user` (admin only)
- **Input:**
  ```json
  { "agent_id": "string", "username": "string" }
  ```
- **Output:**
  ```json
  { "code": 0, "message": "user assigned to client successfully" }
  ```

### GET `/api/clients/:agent_id/otp`
- **Output:**
  ```json
  { "code": 0, "message": "success", "data": { "otp": "string" } }
  ```

### GET `/api/clients/my-otp`
- **Output:**
  ```json
  { "code": 0, "message": "success", "data": { "otp": "string" } }
  ```

---

## Log APIs

### GET `/api/logs/archive` (admin only)
- **Output:**
  ```json
  { "code": 0, "message": "success", "data": [ ...log list... ] }
  ```

### GET `/api/logs/my-device`
- **Output:**
  ```json
  { "code": 0, "message": "success", "data": [ ...log list... ] }
  ```

---

## User Object Example
```json
{
  "id": "string",
  "username": "string",
  "full_name": "string",
  "email": "string",
  "role": "user|admin",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

## Notes
- Tất cả API (trừ /login) đều cần JWT Bearer token ở header: `Authorization: Bearer <token>`
- Các API gắn (admin only) chỉ cho user có role admin.
- Các trường input/output có thể bổ sung thêm tuỳ logic thực tế.
