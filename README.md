# Gou-PC: Golang Learning Project

## 1. Giới thiệu
Đây là dự án cá nhân nhằm thực hành, làm quen với Golang thông qua việc xây dựng một hệ thống quản lý agent/client, user, log và xác thực bảo mật. Dự án tập trung vào kiến trúc sạch (clean architecture), tách biệt các tầng, dễ mở rộng, dễ test, mô phỏng workflow thực tế của một hệ thống quản lý thiết bị phân tán.

## 2. Kiến trúc tổng thể
```
+-----------------------------+        +-------------------+         +-------------------+         +-------------------+
|  Ứng dụng hệ thống (client) |<------>|    Agent (Client) | <-----> |   TCP Server      | <-----> |   RESTful API     |
|  - Lấy OTP qua IPC          |  IPC   |-------------------|   TLS   |-------------------|   HTTP  |-------------------|
+-----------------------------+        | - Đăng ký         |         | - Quản lý agent   |         | - Quản lý user    |
                                       | - Gửi log         |         | - Sinh OTP động   |         | - Quản lý client  |
                                       | - Nhận OTP        |         | - Lưu log         |         | - Quản lý log     |
                                       | - IPC (Windows)   |         | - Mapping agentID |         | - Xác thực JWT    |
                                       +-------------------+         +-------------------+         +-------------------+
```

## 3. Cấu trúc thư mục
```
internal/
├── agent/           # Toàn bộ logic agent: TCP, IPC, đăng ký, gửi log, nhận OTP
│   ├── agent.go
│   ├── manager_agent.go
│   └── ipc.go
├── api/
│   ├── handler/     # Xử lý request/response REST API
│   ├── service/     # Logic nghiệp vụ (user, client, log, OTP)
│   ├── repository/  # Truy xuất dữ liệu (file/json)
│   ├── middleware/  # JWT, logging, CORS
│   ├── model/       # Định nghĩa struct dữ liệu
│   └── response/    # Chuẩn hóa response API
├── config/          # Định nghĩa, load cấu hình server/client
├── crypto/otp.go    # Sinh OTP động chuẩn TOTP
├── tcpserver/       # TCP server nhận/gửi dữ liệu agent
cmd/
└── server/main.go   # Entry point server, khởi tạo config, inject, chạy API & TCP
etc/
├── manager_client.json # Dữ liệu client/agent
├── users.json          # Dữ liệu user
├── archive.log         # Log thu thập
└── ...
```

## 4. Agent (Client)
- **Đăng ký:** Gửi device info lên server, nhận agentID/clientID, lưu vào file cấu hình.
- **Gửi log:** Theo dõi file log, gửi dòng mới lên server qua TCP.
- **Nhận OTP:** Gửi yêu cầu OTP lên server, nhận về mã OTP động (không lưu secret).
- **IPC (Windows):** Mở named pipe, cho phép ứng dụng khác lấy OTP qua IPC.
- **Quản lý agentID/clientID:** Đồng bộ, không trùng lặp, mapping rõ ràng.

## 5. TCP Server
- Lắng nghe kết nối agent qua TLS.
- Xác thực, mapping agentID <-> clientID.
- Nhận log, lưu log, trả OTP động cho agent.
- Giao tiếp thread-safe, đồng bộ dữ liệu agent.

## 6. RESTful API (Gin)
- **Xác thực:** Đăng nhập trả JWT, mọi API (trừ login) đều yêu cầu JWT.
- **User:** CRUD, đổi mật khẩu, cập nhật info, phân quyền.
- **Client/Agent:** CRUD, gán user, lấy theo agentID/userID.
- **OTP:** Sinh OTP động (TOTP) từ clientID/agentID, không lưu secret.
- **Log:** Lấy log archive, log theo thiết bị.
- **Middleware:** JWT, role-based access, logging, CORS.

## 7. Cấu hình
- `internal/config/config.go`: Định nghĩa đường dẫn file, cổng, JWT secret, thời gian sống JWT...
- Dễ dàng mở rộng để load từ file hoặc biến môi trường.

## 8. Hướng dẫn build, run, test
### Yêu cầu
- Go >= 1.21
- Các package đã liệt kê trong go.mod

### Build & Run
```sh
# Cài dependency
go mod tidy
# Build server
go build -o gou-pc-server ./cmd/server
# Chạy server
./gou-pc-server
```

### Test API
- Xem file `api_test_examples.md` để biết cách test API bằng curl/PowerShell.
- Tất cả API (trừ login) đều yêu cầu JWT.

### Đóng gói
- Cấu hình trong `internal/config/config.go` và thư mục `etc/`.
- Build cross-platform dễ dàng với Go.

## 9. Tài liệu API (ví dụ)
- Đăng nhập: `POST /api/login` (username, password)
- CRUD user: `/api/users/*`
- CRUD client: `/api/clients/*`
- Sinh OTP: `/api/clients/:agent_id/otp`, `/api/clients/my-otp`
- Lấy log: `/api/logs/*`
- Xem chi tiết trong code hoặc file test mẫu.

---
Dự án này phục vụ mục đích học tập, thực hành Golang và clean architecture. Mọi góp ý, thắc mắc vui lòng liên hệ .
