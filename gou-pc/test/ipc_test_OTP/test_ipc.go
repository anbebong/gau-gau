package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Microsoft/go-winio"
)

func main() {
	pipePath := `\\.\pipe\MySecretServicePipe`
	log.Println("Đang kết nối đến IPC pipe:", pipePath)

	var conn io.ReadWriteCloser
	var err error

	// Thử kết nối trong vài giây
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

Loop:
	for {
		select {
		case <-timeout:
			log.Fatalf("Không thể kết nối đến pipe sau 5 giây. Đảm bảo client chính đang chạy và đã khởi tạo pipe. Lỗi: %v", err)
		case <-ticker.C:
			conn, err = winio.DialPipe(pipePath, nil)
			if err == nil {
				// Kết nối thành công
				break Loop
			}
		}
	}

	defer conn.Close()
	log.Println("Kết nối thành công. Đang gửi yêu cầu 'GET_SECRET'...")

	// Gửi yêu cầu
	_, err = conn.Write([]byte("GET_SECRET"))
	if err != nil {
		log.Fatalf("Gửi yêu cầu thất bại: %v", err)
	}

	log.Println("Đã gửi yêu cầu. Đang chờ phản hồi...")

	// Đọc phản hồi từ pipe
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		log.Fatalf("Đọc phản hồi thất bại: %v", err)
	}

	if n > 0 {
		response := string(buf[:n])
		fmt.Println("--- PHẢN HỒI TỪ CLIENT CHÍNH ---")
		fmt.Printf("OTP nhận được: %s\n", response)
		fmt.Println("----------------------------------")
	} else {
		log.Println("Không nhận được dữ liệu phản hồi.")
	}
}
