package agent

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Microsoft/go-winio"
)

// var ipcListener net.Listener

// StartIPCListener mở named pipe IPC cho client
// requestOTP: hàm gửi yêu cầu OTP lên server, nhận channel otp để trả về
// serverConn: kết nối TCP tới server (có thể nil nếu chưa kết nối)
func StartIPCListener(requestOTP func(chan<- string) error, serverConn net.Conn) {
	log.Println("Bắt đầu thực thi hàm StartIPCListener.")
	pipePath := `\\.\pipe\MySecretServicePipe`
	_ = os.Remove(pipePath)
	config := &winio.PipeConfig{
		SecurityDescriptor: "D:P(A;;GA;;;WD)",
	}
	ipcListener, err := winio.ListenPipe(pipePath, config)
	if err != nil {
		log.Printf("Không thể lắng nghe trên named pipe: %v", err)
		return
	}
	defer ipcListener.Close()
	log.Printf("Trình nghe IPC đã bắt đầu thành công tại %s", pipePath)
	for {
		conn, err := ipcListener.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") || err == io.ErrClosedPipe {
				log.Println("Trình nghe IPC đã dừng.")
				return
			}
			log.Printf("Không thể chấp nhận kết nối IPC: %v", err)
			continue
		}
		// Tạo channel otp riêng cho từng kết nối
		otpChan := make(chan string, 1)
		go handleIPCConnection(conn, func() error { return requestOTP(otpChan) }, otpChan, serverConn)
	}
}

// handleIPCConnection xử lý một kết nối IPC đến.
func handleIPCConnection(conn net.Conn, requestOTP func() error, otpResponseChan <-chan string, serverConn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			log.Printf("Không thể đọc từ kết nối IPC: %v", err)
		}
		return
	}
	processedRequest := strings.TrimSpace(strings.ReplaceAll(string(buf[:n]), "\x00", ""))
	if processedRequest == "GET_SECRET" {
		log.Println("Yêu cầu 'GET_SECRET' hợp lệ. Đang yêu cầu OTP mới từ server...")
		if serverConn == nil {
			log.Println("Không thể yêu cầu OTP: Không có kết nối đến server.")
			conn.Write([]byte("ERROR: Not connected to server"))
			return
		}
		if err := requestOTP(); err != nil {
			log.Printf("Lỗi khi gửi yêu cầu OTP đến server: %v", err)
			conn.Write([]byte("ERROR: Failed to request OTP from server"))
			return
		}
		select {
		case otp := <-otpResponseChan:
			log.Println("Đã nhận được OTP, đang gửi cho client IPC.")
			conn.Write([]byte(otp))
		case <-time.After(10 * time.Second):
			log.Println("Lỗi: Hết thời gian chờ phản hồi OTP từ server.")
			conn.Write([]byte("ERROR: Timeout waiting for OTP from server"))
		}
	} else {
		log.Printf("Yêu cầu không xác định: '%s'", processedRequest)
		conn.Write([]byte("ERROR: Unknown request"))
	}
}
