package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

// Đường dẫn đến file log của client
const logFilePath = "../event.log"

var logLevels = []string{"INFO", "WARNING", "ERROR", "DEBUG"}
var logMessages = []string{
	"User logged in successfully",
	"User logged out",
	"User logged in with invalid credentials",
	"User registration successful",
	"User login failed password expired",
}

func main() {
	// Mở file log để ghi thêm vào (append)
	// Nếu file chưa tồn tại, nó sẽ được tạo
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Không thể mở file log '%s': %v", logFilePath, err)
	}
	defer file.Close()

	// Tạo một logger mới chỉ để ghi vào file
	fileLogger := log.New(file, "", 0)

	fmt.Printf("Bắt đầu tạo log ngẫu nhiên vào file: %s\n", logFilePath)
	fmt.Println("Nhấn Ctrl+C để dừng.")

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Vòng lặp vô hạn để tạo log
	for {
		// Chọn ngẫu nhiên một level và một message
		level := logLevels[rand.Intn(len(logLevels))]
		message := logMessages[rand.Intn(len(logMessages))]
		timestamp := time.Now().Format("2006-01-02 15:04:05")

		// Tạo dòng log hoàn chỉnh
		logLine := fmt.Sprintf("%s [%s] - %s", timestamp, level, message)

		// Ghi vào file
		fileLogger.Println(logLine)

		// In ra console để người dùng biết nó đang chạy
		fmt.Println(logLine)

		// Chờ một khoảng thời gian ngẫu nhiên trước khi ghi dòng tiếp theo
		// (ví dụ: từ 1 đến 5 giây)
		sleepDuration := time.Duration(rand.Intn(4)+1) * time.Second
		time.Sleep(sleepDuration)
	}
}
