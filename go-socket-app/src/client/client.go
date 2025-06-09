package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go-socket-app/src/common"
	"go-socket-app/src/config"
	"go-socket-app/src/logger"
	"net"
	"os"
	"strings"

	"github.com/google/uuid"
)

func main() {
	configPath := "src/config/config.json" // Đường dẫn file config đã sửa lại cho đúng
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	logFile := "client.log"
	logLevel := logger.DEBUG
	if cfg.LogFile != "" {
		logFile = cfg.LogFile
	}
	if cfg.LogLevel != "" {
		switch cfg.LogLevel {
		case "DEBUG":
			logLevel = logger.DEBUG
		case "INFO":
			logLevel = logger.INFO
		case "WARN":
			logLevel = logger.WARN
		case "ERROR":
			logLevel = logger.ERROR
		}
	}

	log, err := logger.NewMultiLogger(logLevel, logFile)
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	address := net.JoinHostPort(cfg.Address, cfg.Port)
	var conn net.Conn
	if cfg.UseTLS {
		log.Info("Using TLS for connection")
		tlsConfig := &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify}
		conn, err = tls.Dial("tcp", address, tlsConfig)
		if err != nil {
			log.Error("Error connecting to server via TLS: %v", err)
			os.Exit(1)
		}
	} else {
		conn, err = net.Dial("tcp", address)
		if err != nil {
			log.Error("Error connecting to server: %v", err)
			os.Exit(1)
		}
	}
	defer conn.Close()

	var key string
	keyFile := "secret.key"
	var clientID string
	// Đọc clientID và key từ file (nếu có)
	if data, err := os.ReadFile(keyFile); err == nil {
		parts := strings.Split(string(data), ":")
		if len(parts) == 2 {
			clientID = parts[0]
			key = parts[1]
			log.Info("Loaded clientID=%s, key from %s", clientID, keyFile)
		}
	}
	if clientID == "" {
		clientID = uuid.New().String()
	}
	if key == "" {
		// Nếu chưa có key, gửi HELLO để đăng ký
		helloMsg := common.Message{Type: common.HELLO, Content: "Xin chào", ClientID: clientID}
		err := SendMessageConn(conn, helloMsg, log, cfg, key)
		if err != nil {
			log.Error("Error sending HELLO: %v", err)
			return
		}
		resp, err := ReceiveMessageConn(conn, log, cfg, key)
		if err != nil {
			log.Error("Error receiving HELLO response: %v", err)
			return
		}
		if resp.Type == common.HELLO && resp.Key != "" {
			key = resp.Key
			os.WriteFile(keyFile, []byte(clientID+":"+key), 0600)
			log.Info("Received and saved key from server: %s", key)
		} else {
			log.Error("Failed to register key: %s", resp.Content)
			return
		}
	}

	// AUTH: Xác thực key một lần
	authMsg := common.Message{Type: common.AUTH, Key: key, ClientID: clientID}
	err = SendMessageConn(conn, authMsg, log, cfg, key)
	if err != nil {
		log.Error("Error sending AUTH: %v", err)
		return
	}
	authResp, err := ReceiveMessageConn(conn, log, cfg, key)
	if err != nil {
		log.Error("Error receiving AUTH response: %v", err)
		return
	}
	if authResp.Type == common.AUTH && authResp.Content == "Auth success" {
		log.Info("Auth success!")
	} else {
		log.Error("Auth failed: %s", authResp.Content)
		return
	}

	// MSG: Chỉ đợi phản hồi từ server, không tự động gửi MSG liên tục
	for {
		// Nhận tin nhắn từ server
		log.Info("Waiting for messages from server...")
		msgResp, err := ReceiveMessageConn(conn, log, cfg, key)
		if err != nil {
			log.Error("Error receiving message from server: %v", err)
			return
		}
		log.Info("Server message: %s", msgResp.Content)

		// Xử lý lệnh CMD từ server (nếu có)
		if msgResp.Type == common.CMD {
			cmdContent := msgResp.Content
			if cfg.UseMessageEncryption {
				dec, err := common.Decrypt([]byte(key)[:16], cmdContent)
				if err == nil {
					log.Info("[DECRYPTED CMD] %s", dec)
					cmdContent = dec
				} else {
					log.Warn("Failed to decrypt CMD: %v", err)
				}
			}
			log.Info("Received CMD from server: %s", cmdContent)
			if cmdContent == "PING" {
				log.Info("PONG")
			}
			// Có thể thực thi lệnh hệ thống, gọi hàm, ... ở đây
		}

		// Goroutine luôn nhận và xử lý message từ server
		go func() {
			for {
				msgResp, err := ReceiveMessageConn(conn, log, cfg, key)
				if err != nil {
					log.Error("Error receiving message from server: %v", err)
					os.Exit(1)
				}
				log.Info("Server message: %s", msgResp.Content)
				if msgResp.Type == common.CMD {
					cmdContent := msgResp.Content
					if cfg.UseMessageEncryption {
						dec, err := common.Decrypt([]byte(key)[:16], cmdContent)
						if err == nil {
							log.Info("[DECRYPTED CMD] %s", dec)
							cmdContent = dec
						} else {
							log.Warn("Failed to decrypt CMD: %v", err)
						}
					}
					log.Info("Received CMD from server: %s", cmdContent)
					if cmdContent == "PING" {
						log.Info("PONG")
					}
				}
			}
		}()

		// Chỉ giữ client chạy, không nhập MSG từ bàn phím nữa
		select {}
	}
}

// Debug: log chi tiết message gửi đi
func debugSend(log *logger.Logger, msg common.Message) {
	log.Debug("[DEBUG SEND] %+v", msg)
}

// Debug: log chi tiết message nhận về
func debugRecv(log *logger.Logger, msg common.Message) {
	log.Debug("[DEBUG RECV] %+v", msg)
}

func SendMessageConn(conn net.Conn, msg common.Message, log *logger.Logger, cfg *config.Config, key string) error {
	if cfg.UseMessageEncryption {
		if msg.Content != "" {
			enc, err := common.Encrypt([]byte(key)[:16], msg.Content)
			if err != nil {
				log.Error("Error encrypting message content: %v", err)
				return err
			}
			log.Debug("[AUTO ENCRYPT CONTENT] %s => %s", msg.Content, enc)
			msg.Content = enc
		}
		if msg.Key != "" {
			encKey, err := common.Encrypt([]byte(key)[:16], msg.Key)
			if err != nil {
				log.Error("Error encrypting message key: %v", err)
				return err
			}
			log.Debug("[AUTO ENCRYPT KEY] %s => %s", msg.Key, encKey)
			msg.Key = encKey
		}
	}
	debugSend(log, msg)
	encoder := json.NewEncoder(conn)
	return encoder.Encode(msg)
}

func ReceiveMessageConn(conn net.Conn, log *logger.Logger, cfg *config.Config, key string) (common.Message, error) {
	var msg common.Message
	debugRecv(log, msg)
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&msg)
	if cfg.UseMessageEncryption {
		if msg.Content != "" {
			dec, err := common.Decrypt([]byte(key)[:16], msg.Content)
			if err == nil {
				log.Debug("[AUTO DECRYPT CONTENT] %s => %s", msg.Content, dec)
				msg.Content = dec
			}
		}
		if msg.Key != "" {
			decKey, err := common.Decrypt([]byte(key)[:16], msg.Key)
			if err == nil {
				log.Debug("[AUTO DECRYPT KEY] %s => %s", msg.Key, decKey)
				msg.Key = decKey
			}
		}
	}
	return msg, err
}
