package main

import (
	"encoding/binary"
	"fmt"
	"gou-pc/internal/crypto"
	"net"
	"sync"
	"time"
)

func buildHelloMsg(agentID string) ([]byte, error) {
	msg := fmt.Sprintf(`{"type":"hello","data":{"agent_id":"%s"}}`, agentID)
	encryptedMsg, err := crypto.Encrypt(msg)
	if err != nil {
		return nil, err
	}
	msgBytes := []byte(encryptedMsg)
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(msgBytes)))
	return append(lenBuf, msgBytes...), nil
}

type AgentStat struct {
	Success int
	Fail    int
}

func main() {
	const totalConn = 100                 // Số agent đồng thời muốn test
	const eps = 2                         // Số sự kiện mỗi giây cho mỗi agent
	const testDuration = 20 * time.Second // Thời gian test
	crypto.InitCipher()                   // Khởi tạo cipher nếu cần
	var wg sync.WaitGroup
	wg.Add(totalConn)

	messages := make([][]byte, totalConn)
	stats := make([]AgentStat, totalConn) // Lưu kết quả từng agent

	for i := 0; i < totalConn; i++ {
		agentID := fmt.Sprintf("test-agent-%d", i)
		msg, err := buildHelloMsg(agentID)
		if err != nil {
			fmt.Println("Lỗi mã hóa:", err)
			return
		}
		messages[i] = msg
	}

	interval := time.Second / time.Duration(eps)

	start := time.Now()
	for i := 0; i < totalConn; i++ {
		go func(idx int) {
			defer wg.Done()
			conn, err := net.Dial("tcp", "localhost:9000")
			if err != nil {
				fmt.Printf("Agent %d - Lỗi kết nối: %v\n", idx, err)
				return
			}
			defer conn.Close()
			success := 0
			fail := 0
			for time.Since(start) < testDuration {
				_, err := conn.Write(messages[idx])
				if err != nil {
					fail++
				} else {
					success++
				}
				time.Sleep(interval)
			}
			stats[idx].Success = success
			stats[idx].Fail = fail
		}(i)
	}
	wg.Wait()

	// Tổng hợp kết quả
	fmt.Println("=== Kết quả gửi của các agent ===")
	totalSuccess := 0
	totalFail := 0
	for i := 0; i < totalConn; i++ {
		fmt.Printf("Agent %d: thành công %d, thất bại %d\n", i, stats[i].Success, stats[i].Fail)
		totalSuccess += stats[i].Success
		totalFail += stats[i].Fail
	}
	fmt.Printf("Tổng: thành công %d, thất bại %d\n", totalSuccess, totalFail)
	fmt.Printf("Hoàn thành! Đã mở %d agent, mỗi agent gửi ~%d EPS, thời gian chạy: %v\n", totalConn, eps, testDuration)
}
