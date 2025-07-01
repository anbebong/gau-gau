package main

import (
	"encoding/json"
	"fmt"
	"gou-pc/internal/agent"
	"gou-pc/internal/config"
	"gou-pc/internal/logutil"
	"os"
	"time"

	"github.com/kardianos/service"
)

// Đóng gói toàn bộ logic cũ vào hàm mainLogic
func mainLogic() {
	cfg := config.DefaultClientConfig()
	if err := logutil.InitCoreLogger(cfg.LogFile, logutil.DEBUG); err != nil {
		fmt.Printf("Could not open log file: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) >= 2 && os.Args[1] != "install-service" && os.Args[1] != "uninstall-service" {
		cfg.ServerAddr = os.Args[1]
	}

	clientInfo := struct {
		ClientID string `json:"client_id"`
		AgentID  string `json:"agent_id"`
	}{}
	needRegister := false
	if data, err := os.ReadFile(cfg.ConfigFile); err == nil {
		_ = json.Unmarshal(data, &clientInfo)
		if clientInfo.ClientID == "" || clientInfo.AgentID == "" {
			needRegister = true
		}
	} else {
		needRegister = true
	}

	var a *agent.Agent
	for {
		a = &agent.Agent{}
		if err := a.Connect(cfg.ServerAddr, 10*time.Second); err != nil {
			logutil.CoreError("failed to connect: %v", err)
			logutil.CoreInfo("Retrying connect after 10s...")
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}
	defer a.Close()

	if needRegister {
		clientID, agentID, err := agent.RegisterAgent(a, cfg.ConfigFile)
		if err != nil {
			logutil.CoreError("register error: %v", err)
			os.Exit(1)
		}
		clientInfo.ClientID = clientID
		clientInfo.AgentID = agentID
		fmt.Println("Đăng ký thành công, đã lưu client_id và agent_id!")
	}

	fmt.Printf("ClientID: %s, AgentID: %s\n", clientInfo.ClientID, clientInfo.AgentID)

	// IPC: truyền hàm requestOTP nhận channel otp riêng cho từng kết nối
	go agent.StartIPCListener(
		func(otpChan chan<- string) error {
			otpMsg := agent.Message{
				Type: agent.TypeRequestOTP,
				Data: agent.AgentMessageData{
					AgentID: clientInfo.AgentID,
					Payload: nil,
				},
			}
			resp, err := a.Request(otpMsg, 10*time.Second)
			if err != nil {
				logutil.CoreError("Request OTP error: %v", err)
				return err
			}
			logutil.CoreInfo("Received OTP response: %+v", resp)
			// Kiểm tra kiểu dữ liệu trả về
			if m, ok := resp.Data.(map[string]interface{}); ok {
				if otp, ok := m["otp"]; ok {
					switch v := otp.(type) {
					case string:
						logutil.CoreInfo("Parsed OTP (string): %s", v)
						otpChan <- v
						return nil
					case float64:
						logutil.CoreInfo("Parsed OTP (float64): %.0f", v)
						otpChan <- fmt.Sprintf("%.0f", v)
						return nil
					default:
						logutil.CoreError("OTP type unknown: %T, value: %+v", v, v)
					}
				} else {
					logutil.CoreError("OTP field not found in response: %+v", m)
				}
			} else if s, ok := resp.Data.(string); ok {
				// Nếu trả về là string (trường hợp server trả về lỗi dạng chuỗi)
				logutil.CoreError("OTP response error string: %s", s)
				return fmt.Errorf("OTP response error: %s", s)
			} else {
				logutil.CoreError("OTP response parse failed: %+v", resp.Data)
			}
			return fmt.Errorf("no otp in response")
		},
		a.Conn,
	)

	// Gửi hello định kỳ 10s
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			helloMsg := agent.Message{
				Type: agent.TypeHello,
				Data: agent.AgentMessageData{
					AgentID: clientInfo.AgentID,
					Payload: nil,
				},
			}
			_, _ = a.Request(helloMsg, 10*time.Second)
			<-ticker.C
		}
	}()

	// Gửi log: dùng agent chính, không tạo agent riêng, mọi log đều gửi qua a.Request
	go func() {
		logPath := cfg.EventLog
		offsetPath := cfg.OffsetFile
		if !agent.IsAbsPath(offsetPath) {
			cwd, _ := os.Getwd()
			offsetPath = cwd + string(os.PathSeparator) + offsetPath
		}
		var lastSize int64 = 0
		if b, err := os.ReadFile(offsetPath); err == nil {
			fmt.Sscanf(string(b), "%d", &lastSize)
		}
		for {
			file, err := os.Open(logPath)
			if err != nil {
				logutil.CoreError("WatchLogAndSend: open log file error: %v", err)
				time.Sleep(cfg.Interval)
				continue
			}
			stat, err := file.Stat()
			if err != nil {
				file.Close()
				logutil.CoreError("WatchLogAndSend: stat log file error: %v", err)
				time.Sleep(cfg.Interval)
				continue
			}
			if stat.Size() < lastSize {
				lastSize = 0
			}
			if stat.Size() > lastSize {
				file.Seek(lastSize, 0)
				buf := make([]byte, stat.Size()-lastSize)
				_, err := file.Read(buf)
				if err == nil {
					lines := agent.SplitLines(string(buf))
					for _, line := range lines {
						if line == "" {
							continue
						}
						logMsg := agent.LogData{Message: line}
						msgData := agent.AgentMessageData{AgentID: clientInfo.AgentID, Payload: logMsg}
						msg := agent.Message{Type: agent.TypeLog, Data: msgData}
						_, _ = a.Request(msg, 10*time.Second)
					}
				}
				lastSize = stat.Size()
				os.WriteFile(offsetPath, []byte(fmt.Sprintf("%d", lastSize)), 0644)
			}
			file.Close()
			time.Sleep(cfg.Interval)
		}
	}()

	select {}
}

type program struct{}

func (p *program) Start(s service.Service) error {
	go func() {
		// Nếu chạy như service, gọi hàm mainLogic
		mainLogic()
	}()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Xử lý khi dừng service nếu cần
	return nil
}

// Định nghĩa tên service chuẩn, không dấu, không khoảng trắng
const serviceName = "GouPcClientSvc"
const serviceDisplayName = "Gou PC Client Service"
const serviceDescription = "Client agent for Gou PC running as a Windows service."

func main() {
	cfg := &service.Config{
		Name:        serviceName,
		DisplayName: serviceDisplayName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, cfg)
	if err != nil {
		fmt.Println("Create service failed:", err)
		return
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install-service":
			err = s.Install()
			if err != nil {
				logutil.CoreError("Install service failed: %v", err)
				fmt.Println("Install service failed:", err)
				// fmt.Println("Nếu lỗi quyền, hãy chạy terminal với quyền Administrator!")
			} else {
				logutil.CoreInfo("Service installed successfully!")
				fmt.Println("Service installed successfully!")
			}
			return
		case "uninstall-service":
			err = s.Uninstall()
			if err != nil {
				logutil.CoreError("Uninstall service failed: %v", err)
				fmt.Println("Uninstall service failed:", err)
				// fmt.Println("Nếu lỗi quyền, hãy chạy terminal với quyền Administrator!")
			} else {
				logutil.CoreInfo("Service uninstalled successfully!")
				fmt.Println("Service uninstalled successfully!")
			}
			return
		}
	}

	// Nếu không phải lệnh service thì chạy service logic
	err = s.Run()
	if err != nil {
		logutil.CoreError("Run service failed: %v", err)
		fmt.Println("Run service failed:", err)
	}
}
