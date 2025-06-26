package main

import (
	"encoding/json"
	"fmt"
	"gou-pc/internal/agent"
	"gou-pc/internal/config"
	"gou-pc/internal/logutil"
	"os"
	"time"
)

func main() {
	cfg := config.DefaultClientConfig()
	if err := logutil.Init(cfg.LogFile, logutil.DEBUG); err != nil {
		fmt.Printf("Could not open log file: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) >= 2 {
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

	a := &agent.Agent{}
	if err := a.Connect(cfg.ServerAddr, 10*time.Second); err != nil {
		logutil.Error("failed to connect: %v", err)
		os.Exit(1)
	}
	defer a.Close()

	if needRegister {
		clientID, agentID, err := agent.RegisterAgent(a, cfg.ConfigFile)
		if err != nil {
			logutil.Error("register error: %v", err)
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
			if err := a.Send(otpMsg); err != nil {
				logutil.Error("Send OTP error: %v", err)
				return err
			}
			resp, err := a.Receive()
			if err != nil {
				logutil.Error("Receive OTP error: %v", err)
				return err
			}
			logutil.Info("Received OTP response: %+v", resp)
			if m, ok := resp.Data.(map[string]interface{}); ok {
				if otp, ok := m["otp"].(string); ok {
					logutil.Info("Parsed OTP: %s", otp)
					logutil.Info("Received OTP: {agent_id:%s, otp:%s}", m["agent_id"], m["otp"])
					otpChan <- otp
					return nil
				}
			}
			logutil.Error("OTP response parse failed: %+v", resp.Data)
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
			_ = a.Send(helloMsg)
			<-ticker.C
		}
	}()

	// Sau khi demo xong mới bắt đầu gửi log song song
	go a.WatchLogAndSend(cfg.EventLog, cfg.Interval, clientInfo.AgentID)

	select {}
}
