# API List (with curl examples)

## Public

### Đăng nhập
```
curl -X POST http://localhost:8082/api/login -H "Content-Type: application/json" -d '{"username":"admin","password":"1"}'
```

## User (JWT required)

### Tạo user (admin only)
```
curl -X POST http://localhost:8082/api/users/create  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"username":"newuser","password":"123","full_name":"New User","email":"new@example.com"}'
```

### Đổi mật khẩu
```
curl -X POST http://localhost:8082/api/users/change-password -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"username":"user","new_password":"newpass"}'
```

### Cập nhật user
```
curl -X POST http://localhost:8082/api/users/update -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"username":"user","full_name":"User Name","email":"user@example.com"}'
```

### Lấy danh sách user (admin only)
```
curl -X GET http://localhost:8082/api/users -H "Authorization: Bearer $TOKEN"
```

### Cập nhật thông tin cá nhân
```
curl -X POST http://localhost:8082/api/users/update-info -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"username":"user","full_name":"User Name","email":"user@example.com"}'
```

### Xóa user (admin only)
```
curl -X DELETE http://localhost:8082/api/users/delete -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"username":"user"}'
```

## Client (JWT required)

### Lấy tất cả client
```
curl -X GET http://localhost:8082/api/clients -H "Authorization: Bearer $TOKEN"
```

### Lấy client theo agent_id
```
curl -X GET http://localhost:8082/api/clients/<agent_id> -H "Authorization: Bearer $TOKEN"
```

### Lấy client theo client_id
```
curl -X GET http://localhost:8082/api/clients/by-id/<client_id> -H "Authorization: Bearer $TOKEN"
```

### Xóa client theo agent_id (admin only)
```
curl -X POST http://localhost:8082/api/clients/delete-agentid -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"agent_id":"..."}'
```

### Xóa client theo client_id (admin only)
```
curl -X POST http://localhost:8082/api/clients/delete-clientid -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"client_id":"..."}'
```

### Gán user cho client theo agent_id (admin only)
```
curl -X POST http://localhost:8082/api/clients/assign-agentid -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"agent_id":"...","username":"..."}'
```

### Gán user cho client theo client_id (admin only)
```
curl -X POST http://localhost:8082/api/clients/assign-clientid -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"client_id":"...","username":"..."}'
```

### Lấy OTP của client theo agent_id
```
curl -X GET http://localhost:8082/api/clients/<agent_id>/otp -H "Authorization: Bearer $TOKEN"
```

### Lấy OTP của thiết bị mình quản lý
```
curl -X GET http://localhost:8082/api/clients/my-otp -H "Authorization: Bearer $TOKEN"
```

## Log (JWT required)

### Lấy log archive (admin only)
```
curl -X GET http://localhost:8082/api/logs/archive -H "Authorization: Bearer $TOKEN"
```

### Lấy log thiết bị của mình
```
curl -X GET http://localhost:8082/api/logs/my-device -H "Authorization: Bearer $TOKEN"
```
