# Go Socket App

Ứng dụng client-server mẫu sử dụng Go, hỗ trợ TCP/TLS, mã hóa message, cấu hình động qua file JSON.

## Cấu trúc dự án

```
go-socket-app/
├── bin/                # File thực thi client.exe, server.exe
├── src/
│   ├── client/         # Code client
│   ├── server/         # Code server
│   ├── common/         # Định nghĩa message, mã hóa, tiện ích
│   ├── config/         # Định nghĩa & đọc file config
│   └── logger/         # Ghi log đa cấp độ
├── server.crt          # TLS certificate (self-signed)
├── server.key          # TLS private key
├── config.json         # File cấu hình
├── Makefile            # Build nhanh
└── README.md
```

## Cấu hình (src/config/config.json)
```json
{
  "address": "localhost",
  "port": "15151",
  "logLevel": "DEBUG",
  "useTLS": true,
  "useMessageEncryption": true
}
```
- `useTLS`: true để bật TLS (yêu cầu server.crt, server.key)
- `useMessageEncryption`: true để mã hóa nội dung message (AES)

## Build
```powershell
# Build cả client và server
make
# Hoặc build riêng
make build-client
make build-server
```

## Chạy server
```powershell
cd go-socket-app
./bin/server.exe
```

## Chạy client
```powershell
cd go-socket-app
./bin/client.exe
```

## Tính năng
- Giao tiếp TCP hoặc TLS (tùy chọn)
- Định nghĩa nhiều loại bản tin: HELLO, MSG, CMD, AUTH...
- Mã hóa nội dung message bằng AES (tùy chọn)
- Ghi log đa cấp độ ra file và console
- Cấu hình động qua file JSON

## Ghi chú
- Để dùng TLS, cần có file server.crt và server.key hợp lệ.
- Key AES mẫu đang hardcode, nên thay bằng key thực tế cho môi trường production.
- Có thể mở rộng thêm loại bản tin, trường message, hoặc logic xử lý theo nhu cầu.

---