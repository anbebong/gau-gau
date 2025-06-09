package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go-socket-app/src/common"
	"go-socket-app/src/config"
	"go-socket-app/src/logger"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const keyStoreFile = "server_keys.json"

// KeyStore lưu key và trạng thái xác thực cho client
// map[clientID] = {Key: ..., Authed: ...}
type KeyStoreEntry struct {
	Key    string `json:"key"`
	Authed bool   `json:"authed"`
}

var (
	// clientKeys   = make(map[string]string) // map[clientAddr]key
	// clientAuth   = make(map[string]bool)   // map[clientAddr]isAuthed
	// clientKeysMu sync.Mutex
	keyStore = make(map[string]KeyStoreEntry)
	// Map lưu CMD chờ gửi cho từng client
	pendingCMDs   = make(map[string]string)
	pendingCMDsMu sync.Mutex
)

func main() {
	// Đọc cấu hình logger từ file config (ví dụ: logFile, logLevel)
	configPath := "src/config/config.json"
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	logFile := "server.log"
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
	addr := fmt.Sprintf("%s:%s", cfg.Address, cfg.Port)
	var ln net.Listener
	if cfg.UseTLS {
		log.Info("Using TLS for server")
		cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
		if err != nil {
			log.Error("Failed to load certificate or key: %v", err)
			os.Exit(1)
		}
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
		ln, err = tls.Listen("tcp", addr, tlsConfig)
		if err != nil {
			log.Error("Error starting TLS server: %v", err)
			os.Exit(1)
		}
	} else {
		ln, err = net.Listen("tcp", addr)
		if err != nil {
			log.Error("Error starting server: %v", err)
			os.Exit(1)
		}
	}
	defer ln.Close()
	log.Info("Server listening on %s", addr)
	currentConfig = cfg

	LoadKeyStore() // Tải key từ file JSON
	defer SaveKeyStore()

	// Goroutine cho phép nhập lệnh CMD từ bàn phím và gửi tới tất cả client đang xác thực
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Nhập clientID và CMD (clientID:CMD): ")
			if !scanner.Scan() {
				break
			}
			input := scanner.Text()
			if input == "" {
				continue
			}
			// Phân tách clientID và CMD
			sep := strings.Index(input, ":")
			if sep < 0 {
				fmt.Println("Sai định dạng. Nhập theo: clientID:CMD")
				continue
			}
			clientID := input[:sep]
			cmd := input[sep+1:]
			entry, ok := keyStore[clientID]
			if !ok || !entry.Authed {
				fmt.Printf("ClientID %s chưa xác thực hoặc không tồn tại.\n", clientID)
				continue
			}
			pendingCMDsMu.Lock()
			pendingCMDs[clientID] = cmd
			pendingCMDsMu.Unlock()
			fmt.Printf("Đã gửi CMD '%s' cho clientID %s\n", cmd, clientID)
		}
	}()

	// Goroutine tự động spam CMD cho mọi client đã xác thực mỗi 5 giây
	go func() {
		for {
			log.Info("Spamming PING command to all authenticated clients")
			for clientID, entry := range keyStore {
				if entry.Authed {
					pendingCMDsMu.Lock()
					pendingCMDs[clientID] = "PING"
					pendingCMDsMu.Unlock()
				}
			}
			time.Sleep(10 * time.Second)
		}
	}()

	// Goroutine IPC socket nội bộ cho admin/process khác gửi lệnh (Unix domain socket)
	go func() {
		sockPath := "./admin.sock"
		// Xóa file socket cũ nếu có
		os.Remove(sockPath)
		ln, err := net.Listen("unix", sockPath)
		if err != nil {
			log.Error("IPC unix socket listen error: %v", err)
			return
		}
		log.Info("IPC admin unix socket listening on %s", sockPath)
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Warn("IPC accept error: %v", err)
				continue
			}
			go func(c net.Conn) {
				defer c.Close()
				var req struct {
					ClientID string `json:"clientID"`
					CMD      string `json:"cmd"`
				}
				dec := json.NewDecoder(c)
				if err := dec.Decode(&req); err != nil {
					log.Warn("IPC decode error: %v", err)
					return
				}
				pendingCMDsMu.Lock()
				pendingCMDs[req.ClientID] = req.CMD
				pendingCMDsMu.Unlock()
				log.Info("[IPC] Đã nhận CMD '%s' cho clientID %s", req.CMD, req.ClientID)
				c.Write([]byte("OK\n"))
			}(conn)
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Warn("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn, log)
	}
}

func handleConnection(conn net.Conn, log *logger.Logger) {
	defer conn.Close()
	cfg := getConfig()
	var clientID string
	var authed bool

	// Goroutine gửi CMD chủ động cho client
	go func() {
		for {
			if clientID != "" && authed {
				pendingCMDsMu.Lock()
				cmd, hasCMD := pendingCMDs[clientID]
				if hasCMD {
					cmdContent := cmd
					msgType := common.CMD
					if cfg.UseMessageEncryption && msgType != common.HELLO && cmdContent != "" {
						entry, ok := keyStore[clientID]
						if ok {
							enc, err := common.Encrypt([]byte(entry.Key)[:16], cmdContent)
							if err == nil {
								cmdContent = enc
								log.Info("Encrypted CMD for %s: %s", clientID, cmdContent)
							}
						}
					}
					resp := common.Message{Type: msgType, Content: cmdContent}
					err := common.SendMessageConn(conn, resp)
					if err == nil {
						log.Info("Sent CMD '%s' to %s", cmd, clientID)
						delete(pendingCMDs, clientID)
					} else {
						log.Warn("Error writing CMD to connection: %v", err)
					}
				}
				pendingCMDsMu.Unlock()
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		msg, err := common.ReceiveMessageConn(conn)
		if err != nil {
			log.Warn("Error reading from connection: %v", err)
			break
		}
		log.Debug("Received message: %+v", msg)
		clientID = msg.ClientID
		if cfg.UseMessageEncryption {
			entry, ok := keyStore[clientID]
			if ok {
				if msg.Content != "" {
					dec, err := common.Decrypt([]byte(entry.Key)[:16], msg.Content)
					if err == nil {
						log.Debug("[AUTO DECRYPT CONTENT] %s => %s", msg.Content, dec)
						msg.Content = dec
					}
				}
				if msg.Key != "" {
					decKey, err := common.Decrypt([]byte(entry.Key)[:16], msg.Key)
					if err == nil {
						log.Debug("[AUTO DECRYPT KEY] %s => %s", msg.Key, decKey)
						msg.Key = decKey
					}
				}
			}
		}
		var resp common.Message
		if msg.Type == common.HELLO {
			key, err := common.GenerateRandomKey()
			if err != nil {
				resp = common.Message{Type: common.MSG, Content: "Error generating key"}
			} else {
				keyStore[clientID] = KeyStoreEntry{Key: key, Authed: false}
				SaveKeyStore()
				resp = common.Message{Type: common.HELLO, Content: "Key registered", Key: key}
				log.Info("Registered key for %s: %s", clientID, key)
			}
		} else if msg.Type == common.AUTH {
			entry, ok := keyStore[clientID]
			if ok && msg.Key == entry.Key {
				entry.Authed = true
				keyStore[clientID] = entry
				SaveKeyStore()
				resp = common.Message{Type: common.AUTH, Content: "Auth success"}
				authed = true
				log.Info("Client %s authenticated", clientID)
			} else {
				resp = common.Message{Type: common.AUTH, Content: "Auth failed"}
				log.Warn("Client %s failed auth", clientID)
			}
		} else if msg.Type == common.MSG {
			entry, ok := keyStore[clientID]
			authed = ok && entry.Authed
			if ok && authed && cfg.UseMessageEncryption {
				decrypted, err := common.Decrypt([]byte(entry.Key)[:16], msg.Content)
				if err != nil {
					resp = common.Message{Type: common.MSG, Content: "Decrypt failed"}
					log.Warn("Decrypt failed for %s", clientID)
				} else {
					resp = common.Message{Type: common.MSG, Content: "Received: " + decrypted}
					log.Info("Received encrypted MSG from %s: %s", clientID, decrypted)
				}
			} else {
				resp = common.Message{Type: common.MSG, Content: "Not authed or no key"}
				log.Warn("Client %s not authed or no key", clientID)
			}
		} else {
			resp = common.Message{Type: common.MSG, Content: "Unknown message type"}
		}
		err = sendMessageToClient(conn, resp, cfg, keyStore[clientID].Key, log)
		if err != nil {
			log.Warn("Error writing to connection: %v", err)
			break
		}
	}
}

// Helper để lấy config hiện tại (nếu cần dùng trong handleConnection)
var currentConfig *config.Config

func getConfig() *config.Config {
	return currentConfig
}

func LoadKeyStore() {
	f, err := os.Open(keyStoreFile)
	if err != nil {
		return // file chưa tồn tại
	}
	defer f.Close()
	json.NewDecoder(f).Decode(&keyStore)
}

func SaveKeyStore() {
	f, err := os.Create(keyStoreFile)
	if err != nil {
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(keyStore)
}

func GetKeyEntry(clientID string) (KeyStoreEntry, bool) {
	entry, ok := keyStore[clientID]
	return entry, ok
}

func SetKeyEntry(clientID, key string, authed bool) {
	keyStore[clientID] = KeyStoreEntry{Key: key, Authed: authed}
}

// Helper để gửi message về client (tự động mã hóa nếu cần)
func sendMessageToClient(conn net.Conn, msg common.Message, cfg *config.Config, key string, log *logger.Logger) error {
	if cfg.UseMessageEncryption {
		if msg.Content != "" {
			enc, err := common.Encrypt([]byte(key)[:16], msg.Content)
			if err == nil {
				log.Debug("[AUTO ENCRYPT CONTENT] %s => %s", msg.Content, enc)
				msg.Content = enc
			}
		}
		if msg.Key != "" {
			encKey, err := common.Encrypt([]byte(key)[:16], msg.Key)
			if err == nil {
				log.Debug("[AUTO ENCRYPT KEY] %s => %s", msg.Key, encKey)
				msg.Key = encKey
			}
		}
	}
	return common.SendMessageConn(conn, msg)
}
